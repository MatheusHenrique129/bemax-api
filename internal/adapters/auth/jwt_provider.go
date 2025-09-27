package auth

import (
	"errors"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	_defaultExpires = time.Hour * 12
)

type jwtAdapter struct {
	secretKey string
	expiresAt time.Duration
}

func (j *jwtAdapter) GenerateToken(userID uuid.UUID, email string, roles []domain.Role, ttl time.Duration) (dto.GetTokenResponse, apierrors.RestError) {
	if ttl <= 0 {
		ttl = _defaultExpires
	}

	claims := domain.NewTokenUserClaims(userID, email, core.TokenTypeBearer, roles, ttl)

	// TODO Change to SigningMethodRS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return dto.GetTokenResponse{}, apierrors.NewInternalServerRestError("error trying to sign the token", err)
	}

	return dto.GetTokenResponse{
		Token:     signed,
		Timestamp: time.Now().UTC().Local(),
		ExpireAt:  ttl,
	}, nil
}

func (j *jwtAdapter) ValidateToken(tokenString string) (*domain.Claims, apierrors.RestError) {
	if tokenString == "" {
		return nil, apierrors.NewUnauthorizedRestError("token must not be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// TODO change to SigningMethodRSA
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apierrors.NewBadRequestRestError("invalid token signing method")
		}

		return []byte(j.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apierrors.NewUnauthorizedRestError("token has expired")
		}
		return nil, apierrors.NewUnauthorizedRestError("invalid token")
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok || !token.Valid {
		return nil, apierrors.NewUnauthorizedRestError("invalid token")
	}

	return claims, nil
}

func NewJWTAdapter(secretKey string, expires time.Duration) ports.AuthJWT {
	return &jwtAdapter{
		secretKey: secretKey,
		expiresAt: expires,
	}
}
