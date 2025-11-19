package ports

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type ReminderCategoryRepository interface {
	Create(ctx context.Context, category *domain.ReminderCategory) error
	Update(ctx context.Context, category *domain.ReminderCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.ReminderCategory, error)
	FindAllActive(ctx context.Context) ([]domain.ReminderCategory, error)
	FindSystemCategories(ctx context.Context) ([]domain.ReminderCategory, error)
	FindUserCategories(ctx context.Context, userID uuid.UUID) ([]domain.ReminderCategory, error)
	FindAllForUser(ctx context.Context, userID uuid.UUID) ([]domain.ReminderCategory, error)
}
