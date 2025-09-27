package services

import (
	"errors"

	"github.com/MatheusHenrique129/bemax-api/internal/core"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type AuthTokenService interface {
	ValidateAccessToken(accessToken string) (*auth.Claims, apierrors.RestError)
}

type authTokenService struct {
	logger  ports.Logger
	authJWT ports.AuthJWT
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
) AuthTokenService {
	return &authTokenService{
		logger:  logger,
		authJWT: authJWT,
	}
}
