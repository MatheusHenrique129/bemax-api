package ports

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
)

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (domain.Role, error)
}
