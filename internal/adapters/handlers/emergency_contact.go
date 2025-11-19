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

type EmergencyContactHandler interface {
	CreateEmergencyContact(w http.ResponseWriter, r *http.Request)
	UpdateEmergencyContact(w http.ResponseWriter, r *http.Request)
	DeleteEmergencyContact(w http.ResponseWriter, r *http.Request)
	GetEmergencyContactByID(w http.ResponseWriter, r *http.Request)
	GetUserEmergencyContacts(w http.ResponseWriter, r *http.Request)
	SetPrimaryContact(w http.ResponseWriter, r *http.Request)
}

type emergencyContactHandler struct {
	logger                  ports.Logger
	emergencyContactService services.EmergencyContactService
}

func (h emergencyContactHandler) CreateEmergencyContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	var req dto.CreateEmergencyContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid request body"))
		return
	}

	contact, restErr := h.emergencyContactService.CreateEmergencyContact(ctx, claims.UserID, req)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	response := dto.EmergencyContactResponse{
		ID:           contact.ID,
		Name:         contact.Name,
		Relationship: contact.Relationship,
		Phone:        contact.Phone,
		Email:        contact.Email,
		Notes:        contact.Notes,
		IsPrimary:    contact.IsPrimary,
		IsActive:     contact.IsActive,
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h emergencyContactHandler) UpdateEmergencyContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	contactID, err := uuid.Parse(chi.URLParam(r, "contactID"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid contact_id"))
		return
	}

	var req dto.UpdateEmergencyContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid request body"))
		return
	}

	contact, restErr := h.emergencyContactService.UpdateEmergencyContact(ctx, claims.UserID, contactID, req)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(contact)
}

func (h emergencyContactHandler) DeleteEmergencyContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	contactID, err := uuid.Parse(chi.URLParam(r, "contactID"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid contact_id"))
		return
	}

	if restErr := h.emergencyContactService.DeleteEmergencyContact(ctx, claims.UserID, contactID); restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusNoContent)
}

func (h emergencyContactHandler) GetEmergencyContactByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	contactID, err := uuid.Parse(chi.URLParam(r, "contactID"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid contact_id"))
		return
	}

	contact, restErr := h.emergencyContactService.GetEmergencyContactByID(ctx, claims.UserID, contactID)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(contact)
}

func (h emergencyContactHandler) GetUserEmergencyContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	contacts, restErr := h.emergencyContactService.GetUserEmergencyContacts(ctx, claims.UserID)
	if restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(contacts)
}

func (h emergencyContactHandler) SetPrimaryContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	contactID, err := uuid.Parse(chi.URLParam(r, "contactID"))
	if err != nil {
		http_errors.ErrorHandler(w, apierrors.NewBadRequestRestError("invalid contact_id"))
		return
	}

	if restErr := h.emergencyContactService.SetPrimaryContact(ctx, claims.UserID, contactID); restErr != nil {
		http_errors.ErrorHandler(w, restErr)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "primary contact set successfully"})
}

func NewEmergencyContactHandler(logger ports.Logger, emergencyContactService services.EmergencyContactService) EmergencyContactHandler {
	return &emergencyContactHandler{
		logger:                  logger,
		emergencyContactService: emergencyContactService,
	}
}
