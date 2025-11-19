package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

const (
	defaultLimit = 10
	maxLimit     = 100
)

type ReminderHandler interface {
	CreateReminder(w http.ResponseWriter, r *http.Request)
	UpdateReminder(w http.ResponseWriter, r *http.Request)
	DeleteReminder(w http.ResponseWriter, r *http.Request)
	GetReminderByID(w http.ResponseWriter, r *http.Request)
	GetUserReminders(w http.ResponseWriter, r *http.Request)
	GetActiveReminders(w http.ResponseWriter, r *http.Request)
	GetUpcomingReminders(w http.ResponseWriter, r *http.Request)
	CompleteReminder(w http.ResponseWriter, r *http.Request)
	SnoozeReminder(w http.ResponseWriter, r *http.Request)
}

type reminderHandler struct {
	logger          ports.Logger
	reminderService services.ReminderService
}

func (h reminderHandler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	var req dto.CreateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid request body"))
		return
	}

	reminder, restErr := h.reminderService.CreateReminder(ctx, claims.UserID, req)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(reminder)
}

func (h reminderHandler) UpdateReminder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	reminderID, err := uuid.Parse(chi.URLParam(r, "reminder_id"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid reminder_id"))
		return
	}

	var req dto.UpdateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid request body"))
		return
	}

	reminder, restErr := h.reminderService.UpdateReminder(ctx, claims.UserID, reminderID, req)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reminder)
}

func (h reminderHandler) DeleteReminder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	reminderID, err := uuid.Parse(chi.URLParam(r, "reminder_id"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid reminder_id"))
		return
	}

	if restErr := h.reminderService.DeleteReminder(ctx, claims.UserID, reminderID); restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusNoContent)
}

func (h reminderHandler) GetReminderByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	reminderID, err := uuid.Parse(chi.URLParam(r, "reminder_id"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid reminder_id"))
		return
	}

	reminder, restErr := h.reminderService.GetReminderByID(ctx, claims.UserID, reminderID)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reminder)
}

func (h reminderHandler) GetUserReminders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	reminders, restErr := h.reminderService.GetUserReminders(ctx, claims.UserID)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reminders)
}

func (h reminderHandler) GetActiveReminders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	reminders, restErr := h.reminderService.GetActiveReminders(ctx, claims.UserID)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reminders)
}

func (h reminderHandler) GetUpcomingReminders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	limit := defaultLimit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > maxLimit {
				limit = maxLimit
			}
		}
	}

	reminders, restErr := h.reminderService.GetUpcomingReminders(ctx, claims.UserID, limit)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reminders)
}

func (h reminderHandler) CompleteReminder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	reminderID, err := uuid.Parse(chi.URLParam(r, "reminder_id"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid reminder_id"))
		return
	}

	if restErr := h.reminderService.CompleteReminder(ctx, claims.UserID, reminderID); restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "reminder completed successfully"})
}

func (h reminderHandler) SnoozeReminder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	reminderID, err := uuid.Parse(chi.URLParam(r, "reminder_id"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid reminder_id"))
		return
	}

	var req dto.SnoozeReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid request body"))
		return
	}

	if restErr := req.Validate(); restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	duration := time.Duration(req.Duration) * time.Minute
	if restErr := h.reminderService.SnoozeReminder(ctx, claims.UserID, reminderID, duration); restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "reminder snoozed successfully"})
}

func NewReminderHandler(logger ports.Logger, reminderService services.ReminderService) ReminderHandler {
	return &reminderHandler{
		logger:          logger,
		reminderService: reminderService,
	}
}
