package services

import (
	"context"
	"errors"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type AuthTokenService interface {
	Login(ctx context.Context, email, password string, ipAddress, userAgent string) (string, string, apierrors.RestError)
	ValidateAccessToken(accessToken string) (*auth.Claims, apierrors.RestError)
}

type authTokenService struct {
	logger      ports.Logger
	authJWT     ports.AuthJWT
	tokenRepo   ports.TokenRepository
	userService UserService
	roleService RoleService
}

func (a *authTokenService) Login(ctx context.Context, email, password string, ipAddress, userAgent string) (string, string, apierrors.RestError) {
	user, err := a.userService.AuthenticateUser(ctx, email, password, ipAddress, userAgent)
	if err != nil {
		return "", "", err
	}

	accessToken, err := a.authJWT.GenerateToken(user.ID, user.Email, user.Roles, time.Hour*12)
	if err != nil {
		a.logger.Error("failed to generate access token", err)
		return "", "", err
	}

	refreshTokenString := uuid.New().String()
	refreshToken := auth.NewToken(user.ID, refreshTokenString, "refresh", 12)

	if err := a.tokenRepo.Save(ctx, refreshToken); err != nil {
		a.logger.Error("error saving token in database", err)
		return "", "", apierrors.NewInternalServerRestError("error saving token", err)
	}

	return accessToken.Token, refreshTokenString, nil
}

func (a *authTokenService) ValidateAccessToken(accessToken string) (*auth.Claims, apierrors.RestError) {
	claims, err := a.authJWT.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != core.TokenTypeBearer {
		return nil, apierrors.NewUnauthorizedRestError(ErrInvalidToken.Error())
	}

	return claims, nil
}

func NewAuthTokenService(
	logger ports.Logger,
	authJWT ports.AuthJWT,
	userService UserService,
	tokenRepo ports.TokenRepository,
) AuthTokenService {
	return &authTokenService{
		logger:      logger,
		authJWT:     authJWT,
		userService: userService,
		tokenRepo:   tokenRepo,
	}
}
