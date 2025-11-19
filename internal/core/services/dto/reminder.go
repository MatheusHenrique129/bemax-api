package dto

import (
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type CreateReminderRequest struct {
	CategoryID  string                   `json:"category_id" validate:"required"`
	Title       string                   `json:"title" validate:"required,min=3,max=255"`
	Description string                   `json:"description,omitempty"`
	Frequency   domain.ReminderFrequency `json:"frequency" validate:"required"`
	StartDate   string                   `json:"start_date" validate:"required"`
	EndDate     string                   `json:"end_date,omitempty"`
	ReminderAt  string                   `json:"reminder_at" validate:"required"`
	Metadata    map[string]interface{}   `json:"metadata,omitempty"`
}

type UpdateReminderRequest struct {
	CategoryID  string                   `json:"category_id,omitempty"`
	Title       string                   `json:"title,omitempty"`
	Description string                   `json:"description,omitempty"`
	Status      domain.ReminderStatus    `json:"status,omitempty"`
	Frequency   domain.ReminderFrequency `json:"frequency,omitempty"`
	StartDate   string                   `json:"start_date,omitempty"`
	EndDate     string                   `json:"end_date,omitempty"`
	ReminderAt  string                   `json:"reminder_at,omitempty"`
	Metadata    map[string]interface{}   `json:"metadata,omitempty"`
}

type SnoozeReminderRequest struct {
	Duration int `json:"duration_minutes" validate:"required,min=5,max=1440"`
}

type ReminderResponse struct {
	ID             uuid.UUID                `json:"id"`
	Title          string                   `json:"title"`
	Description    string                   `json:"description,omitempty"`
	Category       *CategoryResponse        `json:"category,omitempty"`
	Status         domain.ReminderStatus    `json:"status"`
	Frequency      domain.ReminderFrequency `json:"frequency"`
	StartDate      time.Time                `json:"start_date"`
	EndDate        *time.Time               `json:"end_date,omitempty"`
	ReminderAt     time.Time                `json:"reminder_at"`
	NextOccurrence *time.Time               `json:"next_occurrence,omitempty"`
	IsActive       bool                     `json:"is_active"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

func (r *CreateReminderRequest) Validate() apierrors.RestError {
	if r.Title == "" {
		return apierrors.NewBadRequestRestError("title is required")
	}
	if r.CategoryID == "" {
		return apierrors.NewBadRequestRestError("category_id is required")
	}
	if r.Frequency == "" {
		return apierrors.NewBadRequestRestError("frequency is required")
	}
	if r.StartDate == "" {
		return apierrors.NewBadRequestRestError("start_date is required")
	}
	if r.ReminderAt == "" {
		return apierrors.NewBadRequestRestError("reminder_at is required")
	}
	return nil
}

func (r *SnoozeReminderRequest) Validate() apierrors.RestError {
	if r.Duration < 5 || r.Duration > 1440 {
		return apierrors.NewBadRequestRestError("duration must be between 5 and 1440 minutes")
	}
	return nil
}
