package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/adapters/persistence/mysql"
	"github.com/MatheusHenrique129/bemax-api/internal/core"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/google/uuid"
)

const (
	_defaultRefreshTokenExpires = time.Hour * 24 * 3
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type AuthTokenService interface {
	Login(ctx context.Context, email, password string, ipAddress, userAgent, deviceInfo string) (string, string, time.Duration, apierrors.RestError)
	ValidateAccessToken(accessToken string) (*auth.Claims, apierrors.RestError)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, time.Duration, apierrors.RestError)

	Logout(ctx context.Context, refreshToken string) apierrors.RestError
	LogoutAllDevices(ctx context.Context, userID uuid.UUID) apierrors.RestError
}

type authTokenService struct {
	logger         ports.Logger
	authJWT        ports.AuthJWT
	tokenRepo      ports.TokenRepository
	userService    UserService
	roleService    RoleService
	sessionService SessionService
}

func (a *authTokenService) Login(ctx context.Context, email, password string, ipAddress, userAgent, deviceInfo string) (string, string, time.Duration, apierrors.RestError) {
	user, err := a.userService.AuthenticateUser(ctx, email, password, ipAddress, userAgent)
	if err != nil {
		return "", "", 0, err
	}

	session, err := a.sessionService.CreateSession(ctx, user.ID, deviceInfo, ipAddress, userAgent)
	if err != nil {
		return "", "", 0, err
	}

	accessToken, err := a.authJWT.GenerateToken(user.ID, user.Email, user.Roles, user.TokenVersion, session.SessionID)
	if err != nil {
		a.logger.Error("failed to generate access token", err)
		return "", "", 0, err
	}

	if err := a.sessionService.UpdateSessionToken(ctx, session.SessionID, accessToken.TokenJTI); err != nil {
		return "", "", 0, err
	}

	refreshTokenString := uuid.New().String()
	refreshToken := auth.NewToken(user.ID, refreshTokenString, core.TokenTypeRefresh, _defaultRefreshTokenExpires)

	if err := a.tokenRepo.Save(ctx, refreshToken); err != nil {
		a.logger.Error("error saving token in database", err)
		return "", "", 0, apierrors.NewInternalServerRestError("error saving token", err)
	}

	return accessToken.Token, refreshTokenString, accessToken.ExpireAt, nil
}

func (a *authTokenService) ValidateAccessToken(accessToken string) (*auth.Claims, apierrors.RestError) {
	claims, err := a.authJWT.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != core.TokenTypeBearer {
		return nil, apierrors.NewUnauthorizedRestError(ErrInvalidToken.Error())
	}

	currentVersion, err := a.userService.GetTokenVersion(context.Background(), claims.UserID)
	if err != nil {
		return nil, apierrors.NewInternalServerRestError("error checking token version", err)
	}

	if claims.TokenVersion < currentVersion {
		a.logger.Warn(fmt.Sprintf("outdated token used: token_version=%d, current_version=%d, user=%s",
			claims.TokenVersion, currentVersion, claims.UserID))
		return nil, apierrors.NewUnauthorizedRestError("token has been invalidated")
	}

	if err := a.sessionService.ValidateSessionToken(context.Background(), claims.SessionID, claims.ID); err != nil {
		a.logger.Error("error in validate session token", err)
		return nil, err
	}

	return claims, nil
}

func (a *authTokenService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, time.Duration, apierrors.RestError) {
	token, err := a.tokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, mysql.ErrTokenNotFound) {
			return "", "", 0, apierrors.NewUnauthorizedRestError("invalid refresh token")
		}
		a.logger.Error("error finding refresh token", err)
		return "", "", 0, apierrors.NewInternalServerRestError("error validating token", err)
	}

	if token.IsExpired() {
		return "", "", 0, apierrors.NewUnauthorizedRestError("refresh token expired")
	}

	user, err := a.userService.GetUserByID(ctx, token.UserID)
	if err != nil {
		return "", "", 0, apierrors.NewInternalServerRestError(err.Error(), err)
	}

	if !user.IsActive() {
		return "", "", 0, apierrors.NewUnauthorizedRestError(ErrUserInactive.Error())
	}

	roles, errRoles := a.roleService.GetUserRoles(ctx, user.ID)
	if errRoles != nil {
		return "", "", 0, errRoles
	}

	sessions, err := a.sessionService.GetUserSessions(ctx, user.ID)
	if err != nil {
		return "", "", 0, apierrors.NewUnauthorizedRestError(err.Error())
	}

	if len(sessions) == 0 {
		return "", "", 0, apierrors.NewUnauthorizedRestError("no active session found")
	}

	activeSession := sessions[0]

	newAccessToken, err := a.authJWT.GenerateToken(user.ID, user.Email, roles, user.TokenVersion, activeSession.SessionID)
	if err != nil {
		return "", "", 0, apierrors.NewInternalServerRestError("failed to generate access token", err)
	}

	if err := a.sessionService.UpdateSessionToken(ctx, activeSession.SessionID, newAccessToken.TokenJTI); err != nil {
		return "", "", 0, err
	}

	newRefreshTokenString := uuid.New().String()
	newRefreshToken := auth.NewToken(user.ID, newRefreshTokenString, core.TokenTypeRefresh, _defaultRefreshTokenExpires)

	if err := a.tokenRepo.Save(ctx, newRefreshToken); err != nil {
		return "", "", 0, apierrors.NewInternalServerRestError("error saving new refresh token", err)
	}

	if err := a.tokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		a.logger.Error("failed to revoke old refresh token", err)
		return "", "", 0, apierrors.NewInternalServerRestError("error to revoke old refresh token", err)
	}

	return newAccessToken.Token, newRefreshTokenString, newAccessToken.ExpireAt, nil
}

func (a *authTokenService) Logout(ctx context.Context, refreshToken string) apierrors.RestError {
	token, err := a.tokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, mysql.ErrTokenNotFound) {
			return apierrors.NewNotFoundRestError("refresh token not found")
		}
		a.logger.Error("error finding refresh token for logout", err)
		return apierrors.NewInternalServerRestError("error validating token", err)
	}

	if err := a.userService.IncrementTokenVersion(ctx, token.UserID); err != nil {
		a.logger.Error(fmt.Sprintf("failed to increment token version for user %s", token.UserID), err)
		return err
	}

	if err := a.tokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		a.logger.Error("failed to revoke refresh token", err)
		return apierrors.NewInternalServerRestError("failed to logout", err)
	}

	a.logger.Info(fmt.Sprintf("user %s logged out successfully", token.UserID))
	return nil
}

func (a *authTokenService) LogoutAllDevices(ctx context.Context, userID uuid.UUID) apierrors.RestError {
	if err := a.userService.IncrementTokenVersion(ctx, userID); err != nil {
		a.logger.Error(fmt.Sprintf("failed to increment token version for user %s", userID), err)
		return err
	}

	if err := a.sessionService.TerminateAllUserSessions(ctx, userID); err != nil {
		return err
	}

	if err := a.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		a.logger.Error(fmt.Sprintf("failed to revoke all tokens for user %s", userID), err)
		return apierrors.NewInternalServerRestError("failed to logout from all devices", err)
	}

	a.logger.Info(fmt.Sprintf("user %s logged out from all devices", userID))
	return nil
}

func NewAuthTokenService(
	logger ports.Logger,
	authJWT ports.AuthJWT,
	userService UserService,
	roleService RoleService,
	sessionService SessionService,
	tokenRepo ports.TokenRepository,
) AuthTokenService {
	return &authTokenService{
		logger:         logger,
		authJWT:        authJWT,
		userService:    userService,
		roleService:    roleService,
		sessionService: sessionService,
		tokenRepo:      tokenRepo,
	}
}
