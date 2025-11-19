package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/adapters/persistence/mysql"
	"github.com/MatheusHenrique129/bemax-api/internal/core"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/google/uuid"
)

const (
	_defaultRefreshTokenExpires = time.Minute * 2
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type AuthTokenService interface {
	Login(ctx context.Context, email, password string, ipAddress, userAgent, deviceInfo string) (dto.LoginResponse, apierrors.RestError)
	ValidateAccessToken(accessToken string) (*domain.Claims, apierrors.RestError)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, int64, apierrors.RestError)
	GenerateTokensForSession(ctx context.Context, user *domain.User, session *domain.Session) (dto.FirebaseLoginResponse, apierrors.RestError)

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

func (a *authTokenService) Login(ctx context.Context, email, password string, ipAddress, userAgent, deviceInfo string) (dto.LoginResponse, apierrors.RestError) {
	user, err := a.userService.AuthenticateUser(ctx, email, password, ipAddress, userAgent)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	session, err := a.sessionService.CreateSession(ctx, user.ID, deviceInfo, ipAddress, userAgent)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	accessToken, err := a.authJWT.GenerateToken(user.ID, user.Email, user.Roles, user.TokenVersion, session.SessionID)
	if err != nil {
		a.logger.Error("failed to generate access token", err)
		return dto.LoginResponse{}, err
	}

	if err := a.sessionService.UpdateSessionToken(ctx, session.SessionID, accessToken.TokenJTI); err != nil {
		return dto.LoginResponse{}, err
	}

	refreshTokenString := uuid.New().String()
	refreshToken := domain.NewToken(user.ID, session.ID, refreshTokenString, core.TokenTypeRefresh, _defaultRefreshTokenExpires)

	if err := a.tokenRepo.Save(ctx, refreshToken); err != nil {
		a.logger.Error("error saving token in database", err)
		return dto.LoginResponse{}, apierrors.NewInternalServerRestError("error saving token", err)
	}

	return dto.LoginResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshTokenString,
		TokenType:    string(core.TokenTypeBearer),
		ExpiresIn:    accessToken.ExpireAt,
		User:         user,
	}, nil
}

func (a *authTokenService) ValidateAccessToken(accessToken string) (*domain.Claims, apierrors.RestError) {
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

func (a *authTokenService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, int64, apierrors.RestError) {
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

	sessions, err := a.sessionService.GetUserSessions(ctx, user.ID)
	if err != nil {
		return "", "", 0, apierrors.NewUnauthorizedRestError(err.Error())
	}

	if len(sessions) == 0 {
		return "", "", 0, apierrors.NewUnauthorizedRestError("no active session found")
	}

	activeSession := sessions[0]

	newAccessToken, err := a.authJWT.GenerateToken(user.ID, user.Email, user.Roles, user.TokenVersion, activeSession.SessionID)
	if err != nil {
		return "", "", 0, apierrors.NewInternalServerRestError("failed to generate access token", err)
	}

	if err := a.sessionService.UpdateSessionToken(ctx, activeSession.SessionID, newAccessToken.TokenJTI); err != nil {
		return "", "", 0, err
	}

	newRefreshTokenString := uuid.New().String()
	newRefreshToken := domain.NewToken(user.ID, activeSession.ID, newRefreshTokenString, core.TokenTypeRefresh, _defaultRefreshTokenExpires)

	if err := a.tokenRepo.Save(ctx, newRefreshToken); err != nil {
		return "", "", 0, apierrors.NewInternalServerRestError("error saving new refresh token", err)
	}

	if err := a.tokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		a.logger.Error("failed to revoke old refresh token", err)
		return "", "", 0, apierrors.NewInternalServerRestError("error to revoke old refresh token", err)
	}

	return newAccessToken.Token, newRefreshTokenString, newAccessToken.ExpireAt, nil
}

// GenerateTokensForSession generates access and refresh tokens for an already authenticated user with an existing session.
func (a *authTokenService) GenerateTokensForSession(ctx context.Context, user *domain.User, session *domain.Session) (dto.FirebaseLoginResponse, apierrors.RestError) {
	accessToken, err := a.authJWT.GenerateToken(user.ID, user.Email, user.Roles, user.TokenVersion, session.SessionID)
	if err != nil {
		a.logger.Error("failed to generate access token", err)
		return dto.FirebaseLoginResponse{}, err
	}

	// Update session with token JTI
	if err := a.sessionService.UpdateSessionToken(ctx, session.SessionID, accessToken.TokenJTI); err != nil {
		return dto.FirebaseLoginResponse{}, err
	}

	// Generate refresh token
	refreshTokenString := uuid.New().String()
	refreshToken := domain.NewToken(user.ID, session.ID, refreshTokenString, core.TokenTypeRefresh, _defaultRefreshTokenExpires)

	// Save refresh token to database
	if err := a.tokenRepo.Save(ctx, refreshToken); err != nil {
		a.logger.Error("error saving refresh token in database", err)
		return dto.FirebaseLoginResponse{}, apierrors.NewInternalServerRestError("error saving token", err)
	}

	return dto.FirebaseLoginResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshTokenString,
		TokenType:    string(core.TokenTypeBearer),
		ExpiresIn:    accessToken.ExpireAt,
		User:         user,
	}, nil
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

	if token.SessionID != uuid.Nil {
		session, errSession := a.sessionService.GetSessionByID(ctx, token.SessionID)
		if errSession == nil && session != nil {
			if errDeactivate := a.sessionService.TerminateSession(ctx, session.SessionID); errDeactivate != nil {
				a.logger.Error(fmt.Sprintf("failed to deactivate session for logout: %s", session.SessionID), errDeactivate)
				// Não falha o logout, apenas loga
			}
		}
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
