package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/adapters/constants"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/handlers/middleware"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/MatheusHenrique129/bemax-api/pkg/http_errors"
)

type ProfileHandler interface {
	GetUserProfile(w http.ResponseWriter, r *http.Request)
}

type profileHandler struct {
	logger                  ports.Logger
	userService             services.UserService
	healthProfileService    services.HealthProfileService
	emergencyContactService services.EmergencyContactService
	reminderService         services.ReminderService
}

func (h profileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		http_errors.ErrorHandler(w, apierrors.NewUnauthorizedRestError("user not authenticated"))
		return
	}

	user, err := h.userService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		http_errors.ErrorHandler(w, err)
		return
	}

	healthProfile, errProfile := h.healthProfileService.GetOrCreateHealthProfile(ctx, claims.UserID)
	if errProfile != nil {
		http_errors.ErrorHandler(w, errProfile)
		return
	}

	emergencyContacts, errContact := h.emergencyContactService.GetUserEmergencyContacts(ctx, claims.UserID)
	if errContact != nil {
		http_errors.ErrorHandler(w, errContact)
		return
	}

	reminders, errReminders := h.reminderService.GetActiveReminders(ctx, claims.UserID)
	if errReminders != nil {
		http_errors.ErrorHandler(w, errReminders)
		return
	}

	allReminders, errAllReminders := h.reminderService.GetUserReminders(ctx, claims.UserID)
	if errAllReminders != nil {
		http_errors.ErrorHandler(w, errAllReminders)
		return
	}

	activeReminders, errActiveReminders := h.reminderService.GetActiveReminders(ctx, claims.UserID)
	if errActiveReminders != nil {
		http_errors.ErrorHandler(w, errActiveReminders)
		return
	}

	// Convert domain → DTO (Handler responsibility)
	var healthProfileResponse *dto.HealthProfileResponse
	if healthProfile != nil {
		healthProfileResponse = &dto.HealthProfileResponse{
			ID:                healthProfile.ID,
			BloodType:         healthProfile.BloodType,
			Height:            healthProfile.Height,
			Weight:            healthProfile.Weight,
			Allergies:         healthProfile.Allergies,
			Medications:       healthProfile.Medications,
			MedicalConditions: healthProfile.MedicalConditions,
			Notes:             healthProfile.Notes,
		}
	}

	// Convert emergency contacts to DTO
	var emergencyContactResponses []dto.EmergencyContactResponse
	for _, contact := range emergencyContacts {
		emergencyContactResponses = append(emergencyContactResponses, dto.EmergencyContactResponse{
			ID:           contact.ID,
			Name:         contact.Name,
			Relationship: contact.Relationship,
			Phone:        contact.Phone,
			Email:        contact.Email,
			Notes:        contact.Notes,
			IsPrimary:    contact.IsPrimary,
			IsActive:     contact.IsActive,
		})
	}

	// Convert reminders to DTO
	var reminderResponses []dto.ReminderResponse
	for _, reminder := range reminders {
		var categoryResponse *dto.CategoryResponse
		if reminder.Category != nil {
			categoryResponse = &dto.CategoryResponse{
				ID:           reminder.Category.ID,
				Name:         reminder.Category.Name,
				NameKey:      reminder.Category.NameKey,
				Description:  reminder.Category.Description,
				Icon:         reminder.Category.Icon,
				Color:        reminder.Category.Color,
				Scope:        reminder.Category.Scope,
				DisplayOrder: reminder.Category.DisplayOrder,
			}
		}

		reminderResponses = append(reminderResponses, dto.ReminderResponse{
			ID:             reminder.ID,
			Title:          reminder.Title,
			Description:    reminder.Description,
			Category:       categoryResponse,
			Status:         reminder.Status,
			Frequency:      reminder.Frequency,
			StartDate:      reminder.StartDate,
			EndDate:        reminder.EndDate,
			ReminderAt:     reminder.ReminderAt,
			NextOccurrence: reminder.NextOccurrence,
			IsActive:       reminder.IsActive,
			CreatedAt:      reminder.CreatedAt,
			UpdatedAt:      reminder.UpdatedAt,
		})
	}

	// Calculate statistics
	todayReminders := 0
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	for _, r := range activeReminders {
		if r.ReminderAt.After(startOfDay) && r.ReminderAt.Before(endOfDay) {
			todayReminders++
		}
	}

	stats := &dto.ProfileStats{
		TotalReminders:    len(allReminders),
		ActiveReminders:   len(activeReminders),
		UpcomingReminders: len(reminderResponses),
		TodayReminders:    todayReminders,
	}

	// Build response
	response := dto.UserProfileResponse{
		User:              &user,
		HealthProfile:     healthProfileResponse,
		EmergencyContacts: emergencyContactResponses,
		Reminders:         reminderResponses,
		Stats:             stats,
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func NewProfileHandler(
	logger ports.Logger,
	userService services.UserService,
	healthProfileService services.HealthProfileService,
	emergencyContactService services.EmergencyContactService,
	reminderService services.ReminderService,
) ProfileHandler {
	return &profileHandler{
		logger:                  logger,
		userService:             userService,
		healthProfileService:    healthProfileService,
		emergencyContactService: emergencyContactService,
		reminderService:         reminderService,
	}
}
