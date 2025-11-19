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
)

type HealthProfileHandler interface {
	GetHealthProfile(w http.ResponseWriter, r *http.Request)
	UpdateHealthProfile(w http.ResponseWriter, r *http.Request)
}

type healthProfileHandler struct {
	logger               ports.Logger
	healthProfileService services.HealthProfileService
}

func (h healthProfileHandler) GetHealthProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	// Service returns domain.HealthProfile
	profile, restErr := h.healthProfileService.GetOrCreateHealthProfile(ctx, claims.UserID)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	response := dto.HealthProfileResponse{
		ID:                profile.ID,
		BloodType:         profile.BloodType,
		Height:            profile.Height,
		Weight:            profile.Weight,
		Allergies:         profile.Allergies,
		Medications:       profile.Medications,
		MedicalConditions: profile.MedicalConditions,
		Notes:             profile.Notes,
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h healthProfileHandler) UpdateHealthProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	var req dto.UpdateHealthProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid request body"))
		return
	}

	// Service returns domain.HealthProfile
	profile, restErr := h.healthProfileService.UpdateHealthProfile(ctx, claims.UserID, req)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	response := dto.HealthProfileResponse{
		ID:                profile.ID,
		BloodType:         profile.BloodType,
		Height:            profile.Height,
		Weight:            profile.Weight,
		Allergies:         profile.Allergies,
		Medications:       profile.Medications,
		MedicalConditions: profile.MedicalConditions,
		Notes:             profile.Notes,
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func NewHealthProfileHandler(logger ports.Logger, healthProfileService services.HealthProfileService) HealthProfileHandler {
	return &healthProfileHandler{
		logger:               logger,
		healthProfileService: healthProfileService,
	}
}
