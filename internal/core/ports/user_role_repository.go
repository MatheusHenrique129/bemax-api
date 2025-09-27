package ports

import (
	"context"

	"github.com/google/uuid"
)

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
}
