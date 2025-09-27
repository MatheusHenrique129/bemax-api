package ports

import (
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type AuthJWT interface {
	GenerateToken(userID uuid.UUID, email string, roles []auth.Role, ttl time.Duration) (string, apierrors.RestError)
	ValidateToken(tokenString string) (*auth.Claims, apierrors.RestError)
}
