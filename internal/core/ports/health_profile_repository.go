package ports

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type HealthProfileRepository interface {
	Create(ctx context.Context, profile *domain.HealthProfile) error
	Update(ctx context.Context, profile *domain.HealthProfile) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.HealthProfile, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}
