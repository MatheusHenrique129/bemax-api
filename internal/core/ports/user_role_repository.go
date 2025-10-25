package ports

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
	FindRolesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Role, error)
}
