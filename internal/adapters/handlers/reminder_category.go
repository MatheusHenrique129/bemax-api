package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MatheusHenrique129/bemax-api/internal/adapters/constants"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/handlers/middleware"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/MatheusHenrique129/bemax-api/pkg/http_errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ReminderCategoryHandler interface {
	CreateUserCategory(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoriesForUser(w http.ResponseWriter, r *http.Request)
}

type reminderCategoryHandler struct {
	logger          ports.Logger
	categoryService services.ReminderCategoryService
}

func (h reminderCategoryHandler) CreateUserCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid request body"))
		return
	}

	category, restErr := h.categoryService.CreateUserCategory(ctx, claims.UserID, req)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(category)
}

func (h reminderCategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	categoryID, err := uuid.Parse(chi.URLParam(r, "category_id"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid category_id"))
		return
	}

	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid request body"))
		return
	}

	category, restErr := h.categoryService.UpdateCategory(ctx, claims.UserID, categoryID, req)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(category)
}

func (h reminderCategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	categoryID, err := uuid.Parse(chi.URLParam(r, "category_id"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid category_id"))
		return
	}

	if restErr := h.categoryService.DeleteCategory(ctx, claims.UserID, categoryID); restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusNoContent)
}

func (h reminderCategoryHandler) GetCategoriesForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	categories, restErr := h.categoryService.GetCategoriesForUser(ctx, claims.UserID)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(categories)
}

func NewReminderCategoryHandler(logger ports.Logger, categoryService services.ReminderCategoryService) ReminderCategoryHandler {
	return &reminderCategoryHandler{
		logger:          logger,
		categoryService: categoryService,
	}
}
