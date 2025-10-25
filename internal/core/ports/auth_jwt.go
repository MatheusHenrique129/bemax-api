package ports

import (
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/google/uuid"
)

type AuthJWT interface {
	GenerateToken(userID uuid.UUID, email string, roles []auth.Role, tokenVersion int, sessionID string) (dto.GetTokenResponse, apierrors.RestError)
	ValidateToken(tokenString string) (*auth.Claims, apierrors.RestError)
}
