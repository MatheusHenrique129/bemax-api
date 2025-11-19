package ports

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type ReminderRepository interface {
	Create(ctx context.Context, reminder *domain.Reminder) error
	Update(ctx context.Context, reminder *domain.Reminder) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Reminder, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Reminder, error)
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Reminder, error)
	FindUpcoming(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Reminder, error)
}
