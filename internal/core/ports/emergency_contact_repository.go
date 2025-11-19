package ports

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type EmergencyContactRepository interface {
	Create(ctx context.Context, contact *domain.EmergencyContact) error
	Update(ctx context.Context, contact *domain.EmergencyContact) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.EmergencyContact, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.EmergencyContact, error)
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]domain.EmergencyContact, error)
	FindPrimaryByUserID(ctx context.Context, userID uuid.UUID) (*domain.EmergencyContact, error)
	UnsetAllPrimaryForUser(ctx context.Context, userID uuid.UUID) error
}
