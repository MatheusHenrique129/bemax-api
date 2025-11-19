package services

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/google/uuid"
)

type ReminderCategoryService interface {
	CreateUserCategory(ctx context.Context, userID uuid.UUID, req dto.CreateCategoryRequest) (*domain.ReminderCategory, apierrors.RestError)
	GetCategoriesForUser(ctx context.Context, userID uuid.UUID) ([]domain.ReminderCategory, apierrors.RestError)
	UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, req dto.UpdateCategoryRequest) (*domain.ReminderCategory, apierrors.RestError)
	DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) apierrors.RestError
}

type reminderCategoryService struct {
	logger       ports.Logger
	categoryRepo ports.ReminderCategoryRepository
}

func (r *reminderCategoryService) CreateUserCategory(ctx context.Context, userID uuid.UUID, req dto.CreateCategoryRequest) (*domain.ReminderCategory, apierrors.RestError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	category := domain.NewUserCategory(userID, req.Name, req.Description, req.Icon, req.Color)

	if err := r.categoryRepo.Create(ctx, category); err != nil {
		r.logger.Error("failed to create category", err)
		return nil, apierrors.NewInternalServerRestError("failed to create category", err)
	}

	return category, nil
}

func (r *reminderCategoryService) GetCategoriesForUser(ctx context.Context, userID uuid.UUID) ([]domain.ReminderCategory, apierrors.RestError) {
	categories, err := r.categoryRepo.FindAllForUser(ctx, userID)
	if err != nil {
		r.logger.Error("failed to get categories", err)
		return nil, apierrors.NewInternalServerRestError("failed to get categories", err)
	}

	return categories, nil
}

func (r *reminderCategoryService) UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, req dto.UpdateCategoryRequest) (*domain.ReminderCategory, apierrors.RestError) {
	category, err := r.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, apierrors.NewNotFoundRestError("category not found")
	}

	if category.UserID == nil || *category.UserID != userID {
		return nil, apierrors.NewForbiddenRestError("cannot modify system category or category from another user")
	}

	category.Update(req.Name, req.Description, req.Icon, req.Color, req.DisplayOrder)

	if updateErr := r.categoryRepo.Update(ctx, category); updateErr != nil {
		r.logger.Error("failed to update category", updateErr)
		return nil, apierrors.NewInternalServerRestError("failed to update category", updateErr)
	}

	return category, nil
}

func (r *reminderCategoryService) DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) apierrors.RestError {
	category, err := r.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return apierrors.NewNotFoundRestError("category not found")
	}

	if category.UserID == nil || *category.UserID != userID {
		return apierrors.NewForbiddenRestError("cannot delete system category or category from another user")
	}

	if deleteErr := r.categoryRepo.Delete(ctx, categoryID); deleteErr != nil {
		r.logger.Error("failed to delete category", deleteErr)
		return apierrors.NewInternalServerRestError("failed to delete category", deleteErr)
	}

	return nil
}

func NewReminderCategoryService(logger ports.Logger, categoryRepo ports.ReminderCategoryRepository) ReminderCategoryService {
	return &reminderCategoryService{
		logger:       logger,
		categoryRepo: categoryRepo,
	}
}
