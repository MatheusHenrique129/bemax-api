package ports

import (
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/google/uuid"
)

type AuthJWT interface {
	GetTTL() time.Duration
	GenerateToken(userID uuid.UUID, email string, roles []auth.Role, ttl time.Duration) (dto.GetTokenResponse, apierrors.RestError)
	ValidateToken(tokenString string) (*auth.Claims, apierrors.RestError)
}
