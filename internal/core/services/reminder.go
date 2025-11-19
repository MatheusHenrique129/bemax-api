package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/google/uuid"
)

type ReminderService interface {
	CreateReminder(ctx context.Context, userID uuid.UUID, req dto.CreateReminderRequest) (*domain.Reminder, apierrors.RestError)
	UpdateReminder(ctx context.Context, userID, reminderID uuid.UUID, req dto.UpdateReminderRequest) (*domain.Reminder, apierrors.RestError)
	DeleteReminder(ctx context.Context, userID, reminderID uuid.UUID) apierrors.RestError
	GetReminderByID(ctx context.Context, userID, reminderID uuid.UUID) (*domain.Reminder, apierrors.RestError)
	GetUserReminders(ctx context.Context, userID uuid.UUID) ([]domain.Reminder, apierrors.RestError)
	GetActiveReminders(ctx context.Context, userID uuid.UUID) ([]domain.Reminder, apierrors.RestError)
	GetUpcomingReminders(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Reminder, apierrors.RestError)
	CompleteReminder(ctx context.Context, userID, reminderID uuid.UUID) apierrors.RestError
	SnoozeReminder(ctx context.Context, userID, reminderID uuid.UUID, duration time.Duration) apierrors.RestError
}

type reminderService struct {
	logger       ports.Logger
	reminderRepo ports.ReminderRepository
}

func (r *reminderService) CreateReminder(ctx context.Context, userID uuid.UUID, req dto.CreateReminderRequest) (*domain.Reminder, apierrors.RestError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, apierrors.NewBadRequestRestError("invalid category_id format")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, apierrors.NewBadRequestRestError("invalid start_date format, use YYYY-MM-DD")
	}

	// Validate start date is not too far in the future
	if startDate.After(time.Now().UTC().AddDate(10, 0, 0)) {
		return nil, apierrors.NewBadRequestRestError("start_date cannot be more than 10 years in the future")
	}

	reminderAt, err := time.Parse(time.RFC3339, req.ReminderAt)
	if err != nil {
		return nil, apierrors.NewBadRequestRestError("invalid reminder_at format, use RFC3339")
	}

	if reminderAt.Before(time.Now().UTC()) {
		return nil, apierrors.NewBadRequestRestError("reminder_at must be in the future")
	}

	reminder := domain.NewReminder(userID, categoryID, req.Title, req.Description, req.Frequency, startDate, reminderAt)

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, apierrors.NewBadRequestRestError("invalid end_date format, use YYYY-MM-DD")
		}
		if endDate.Before(startDate) {
			return nil, apierrors.NewBadRequestRestError("end_date must be after start_date")
		}
		reminder.EndDate = &endDate
	}

	if req.Metadata != nil {
		metadataJSON, err := json.Marshal(req.Metadata)
		if err != nil {
			r.logger.Error("failed to marshal metadata", err)
			return nil, apierrors.NewBadRequestRestError("invalid metadata format")
		}
		reminder.Metadata = string(metadataJSON)
	}

	if err := r.reminderRepo.Create(ctx, reminder); err != nil {
		r.logger.Error("failed to create reminder", err)
		return nil, apierrors.NewInternalServerRestError("failed to create reminder", err)
	}

	return reminder, nil
}

func (r *reminderService) UpdateReminder(ctx context.Context, userID, reminderID uuid.UUID, req dto.UpdateReminderRequest) (*domain.Reminder, apierrors.RestError) {
	reminder, err := r.GetReminderByID(ctx, userID, reminderID)
	if err != nil {
		return nil, err
	}

	var categoryID uuid.UUID
	if req.CategoryID != "" {
		parsedCategoryID, parseErr := uuid.Parse(req.CategoryID)
		if parseErr != nil {
			return nil, apierrors.NewBadRequestRestError("invalid category_id format")
		}
		categoryID = parsedCategoryID
	}

	var reminderAt time.Time
	if req.ReminderAt != "" {
		parsedReminderAt, parseErr := time.Parse(time.RFC3339, req.ReminderAt)
		if parseErr != nil {
			return nil, apierrors.NewBadRequestRestError("invalid reminder_at format, use RFC3339")
		}
		// Validate it's in the future
		if parsedReminderAt.Before(time.Now().UTC()) {
			return nil, apierrors.NewBadRequestRestError("reminder_at must be in the future")
		}
		reminderAt = parsedReminderAt
	}

	var startDate time.Time
	if req.StartDate != "" {
		parsedStartDate, parseErr := time.Parse("2006-01-02", req.StartDate)
		if parseErr != nil {
			return nil, apierrors.NewBadRequestRestError("invalid start_date format, use YYYY-MM-DD")
		}
		if parsedStartDate.After(time.Now().UTC().AddDate(10, 0, 0)) {
			return nil, apierrors.NewBadRequestRestError("start_date cannot be more than 10 years in the future")
		}
		startDate = parsedStartDate
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsedEndDate, parseErr := time.Parse("2006-01-02", req.EndDate)
		if parseErr != nil {
			return nil, apierrors.NewBadRequestRestError("invalid end_date format, use YYYY-MM-DD")
		}
		// Validate end_date is after start_date
		compareDate := startDate
		if compareDate.IsZero() {
			compareDate = reminder.StartDate
		}
		if parsedEndDate.Before(compareDate) {
			return nil, apierrors.NewBadRequestRestError("end_date must be after start_date")
		}
		endDate = &parsedEndDate
	}

	var metadataStr string
	if req.Metadata != nil {
		metadataJSON, jsonErr := json.Marshal(req.Metadata)
		if jsonErr != nil {
			r.logger.Error("failed to marshal metadata", jsonErr)
			return nil, apierrors.NewBadRequestRestError("invalid metadata format")
		}
		metadataStr = string(metadataJSON)
	}

	if req.Frequency != "" || req.StartDate != "" || req.EndDate != "" || req.Metadata != nil {
		reminder.UpdateExtended(
			req.Title,
			req.Description,
			categoryID,
			req.Status,
			req.Frequency,
			startDate,
			endDate,
			reminderAt,
			metadataStr,
		)
	} else {
		// Use basic Update for simple changes
		reminder.Update(req.Title, req.Description, categoryID, req.Status, reminderAt)
	}

	if updateErr := r.reminderRepo.Update(ctx, reminder); updateErr != nil {
		r.logger.Error("failed to update reminder", updateErr)
		return nil, apierrors.NewInternalServerRestError("failed to update reminder", updateErr)
	}

	return reminder, nil
}

func (r *reminderService) GetUserReminders(ctx context.Context, userID uuid.UUID) ([]domain.Reminder, apierrors.RestError) {
	reminders, err := r.reminderRepo.FindByUserID(ctx, userID)
	if err != nil {
		r.logger.Error("failed to get user reminders", err)
		return nil, apierrors.NewInternalServerRestError("failed to get reminders", err)
	}
	return reminders, nil
}

func (r *reminderService) GetActiveReminders(ctx context.Context, userID uuid.UUID) ([]domain.Reminder, apierrors.RestError) {
	reminders, err := r.reminderRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		r.logger.Error("failed to get active reminders", err)
		return nil, apierrors.NewInternalServerRestError("failed to get active reminders", err)
	}
	return reminders, nil
}

func (r *reminderService) GetUpcomingReminders(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Reminder, apierrors.RestError) {
	reminders, err := r.reminderRepo.FindUpcoming(ctx, userID, limit)
	if err != nil {
		r.logger.Error("failed to get upcoming reminders", err)
		return nil, apierrors.NewInternalServerRestError("failed to get upcoming reminders", err)
	}
	return reminders, nil
}

func (r *reminderService) GetReminderByID(ctx context.Context, userID, reminderID uuid.UUID) (*domain.Reminder, apierrors.RestError) {
	reminder, err := r.reminderRepo.FindByID(ctx, reminderID)
	if err != nil {
		r.logger.Error(fmt.Sprintf("reminder not found: %s", reminderID), err)
		return nil, apierrors.NewNotFoundRestError("reminder not found")
	}

	if reminder.UserID != userID {
		r.logger.Warn(fmt.Sprintf("user %s attempted to access reminder %s owned by %s", userID, reminderID, reminder.UserID))
		return nil, apierrors.NewForbiddenRestError("access denied")
	}

	return reminder, nil
}

func (r *reminderService) CompleteReminder(ctx context.Context, userID, reminderID uuid.UUID) apierrors.RestError {
	reminder, err := r.GetReminderByID(ctx, userID, reminderID)
	if err != nil {
		return err
	}

	reminder.CompleteAndScheduleNext()

	if updateErr := r.reminderRepo.Update(ctx, reminder); updateErr != nil {
		r.logger.Error("failed to complete reminder", updateErr)
		return apierrors.NewInternalServerRestError("failed to complete reminder", updateErr)
	}

	return nil
}

func (r *reminderService) DeleteReminder(ctx context.Context, userID, reminderID uuid.UUID) apierrors.RestError {
	reminder, err := r.GetReminderByID(ctx, userID, reminderID)
	if err != nil {
		return err
	}

	if deleteErr := r.reminderRepo.Delete(ctx, reminder.ID); deleteErr != nil {
		r.logger.Error("failed to delete reminder", deleteErr)
		return apierrors.NewInternalServerRestError("failed to delete reminder", deleteErr)
	}

	return nil
}

func (r *reminderService) SnoozeReminder(ctx context.Context, userID, reminderID uuid.UUID, duration time.Duration) apierrors.RestError {
	reminder, err := r.GetReminderByID(ctx, userID, reminderID)
	if err != nil {
		return err
	}

	if duration < 5*time.Minute {
		return apierrors.NewBadRequestRestError("snooze duration must be at least 5 minutes")
	}
	if duration > 24*time.Hour {
		return apierrors.NewBadRequestRestError("snooze duration cannot exceed 24 hours")
	}

	reminder.Snooze(duration)

	if updateErr := r.reminderRepo.Update(ctx, reminder); updateErr != nil {
		r.logger.Error("failed to snooze reminder", updateErr)
		return apierrors.NewInternalServerRestError("failed to snooze reminder", updateErr)
	}

	return nil
}

func NewReminderService(logger ports.Logger, reminderRepo ports.ReminderRepository) ReminderService {
	return &reminderService{
		logger:       logger,
		reminderRepo: reminderRepo,
	}
}
