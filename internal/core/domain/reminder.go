package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReminderFrequency defines how often a reminder repeats
type ReminderFrequency string

const (
	ReminderFrequencyOnce    ReminderFrequency = "once"
	ReminderFrequencyDaily   ReminderFrequency = "daily"
	ReminderFrequencyWeekly  ReminderFrequency = "weekly"
	ReminderFrequencyMonthly ReminderFrequency = "monthly"
	ReminderFrequencyCustom  ReminderFrequency = "custom"
)

// ReminderStatus defines the status of a reminder
type ReminderStatus string

const (
	ReminderStatusActive    ReminderStatus = "active"
	ReminderStatusCompleted ReminderStatus = "completed"
	ReminderStatusCancelled ReminderStatus = "cancelled"
	ReminderStatusSnoozed   ReminderStatus = "snoozed"
)

// Reminder represents a user reminder for medications, appointments, etc
type Reminder struct {
	ID             uuid.UUID         `json:"id"`
	UserID         uuid.UUID         `json:"user_id"`
	CategoryID     uuid.UUID         `json:"category_id"`
	Category       *ReminderCategory `json:"category,omitempty"`
	Title          string            `json:"title"`
	Description    string            `json:"description,omitempty"`
	Status         ReminderStatus    `json:"status"`
	Frequency      ReminderFrequency `json:"frequency"`
	StartDate      time.Time         `json:"start_date"`
	EndDate        *time.Time        `json:"end_date,omitempty"`
	ReminderAt     time.Time         `json:"reminder_at"`
	NextOccurrence *time.Time        `json:"next_occurrence,omitempty"`
	IsActive       bool              `json:"is_active"`
	Metadata       string            `json:"metadata,omitempty"` // JSON field for custom data
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// Update updates reminder information
func (r *Reminder) Update(title, description string, categoryID uuid.UUID, status ReminderStatus, reminderAt time.Time) {
	if title != "" {
		r.Title = title
	}
	if description != "" {
		r.Description = description
	}
	if categoryID != uuid.Nil {
		r.CategoryID = categoryID
	}
	if status != "" {
		r.Status = status
	}
	if !reminderAt.IsZero() {
		r.ReminderAt = reminderAt
		r.CalculateNextOccurrence()
	}
	r.UpdatedAt = time.Now().UTC()
}

// UpdateExtended updates reminder with frequency and dates (full update)
func (r *Reminder) UpdateExtended(
	title, description string,
	categoryID uuid.UUID,
	status ReminderStatus,
	frequency ReminderFrequency,
	startDate time.Time,
	endDate *time.Time,
	reminderAt time.Time,
	metadata string,
) {
	if title != "" {
		r.Title = title
	}
	if description != "" {
		r.Description = description
	}
	if categoryID != uuid.Nil {
		r.CategoryID = categoryID
	}
	if status != "" {
		r.Status = status
	}
	if frequency != "" {
		r.Frequency = frequency
	}
	if !startDate.IsZero() {
		r.StartDate = startDate
	}
	// EndDate can be explicitly set to nil to remove it
	r.EndDate = endDate

	if !reminderAt.IsZero() {
		r.ReminderAt = reminderAt
	}
	if metadata != "" {
		r.Metadata = metadata
	}

	// Recalculate next occurrence when frequency or dates change
	r.CalculateNextOccurrence()
	r.UpdatedAt = time.Now().UTC()
}

// Complete marks a reminder as completed
func (r *Reminder) Complete() {
	r.Status = ReminderStatusCompleted
	r.IsActive = false
	r.UpdatedAt = time.Now().UTC()
}

// CompleteAndScheduleNext marks as completed and calculates next occurrence for recurring reminders
func (r *Reminder) CompleteAndScheduleNext() {
	if r.Frequency == ReminderFrequencyOnce {
		r.Complete()
		return
	}

	// For recurring reminders, keep active and calculate next
	r.Status = ReminderStatusActive
	r.CalculateNextOccurrence()
	r.UpdatedAt = time.Now().UTC()
}

// Cancel marks a reminder as cancelled
func (r *Reminder) Cancel() {
	r.Status = ReminderStatusCancelled
	r.IsActive = false
	r.UpdatedAt = time.Now().UTC()
}

// Snooze postpones a reminder
func (r *Reminder) Snooze(duration time.Duration) {
	r.Status = ReminderStatusSnoozed
	newTime := time.Now().UTC().Add(duration)
	r.NextOccurrence = &newTime
	r.UpdatedAt = time.Now().UTC()
}

// Activate reactivates a reminder
func (r *Reminder) Activate() {
	r.Status = ReminderStatusActive
	r.IsActive = true
	r.UpdatedAt = time.Now().UTC()
}

// CalculateNextOccurrence calculates the next occurrence based on frequency
func (r *Reminder) CalculateNextOccurrence() {
	if r.Frequency == ReminderFrequencyOnce {
		r.NextOccurrence = nil
		return
	}

	current := r.ReminderAt
	if r.NextOccurrence != nil {
		current = *r.NextOccurrence
	}

	var next time.Time
	switch r.Frequency {
	case ReminderFrequencyDaily:
		next = current.AddDate(0, 0, 1)
	case ReminderFrequencyWeekly:
		next = current.AddDate(0, 0, 7)
	case ReminderFrequencyMonthly:
		next = current.AddDate(0, 1, 0)
	default:
		return
	}

	// Check if next is within end date
	if r.EndDate != nil && next.After(*r.EndDate) {
		r.IsActive = false
		r.NextOccurrence = nil
		return
	}

	r.NextOccurrence = &next
}

// NewReminder creates a new reminder
func NewReminder(userID, categoryID uuid.UUID, title, description string, frequency ReminderFrequency, startDate, reminderAt time.Time) *Reminder {
	now := time.Now().UTC()
	reminder := &Reminder{
		ID:             uuid.New(),
		UserID:         userID,
		CategoryID:     categoryID,
		Title:          title,
		Description:    description,
		Status:         ReminderStatusActive,
		Frequency:      frequency,
		StartDate:      startDate,
		ReminderAt:     reminderAt,
		NextOccurrence: &reminderAt,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Calculate next occurrence if recurring
	if frequency != ReminderFrequencyOnce {
		reminder.CalculateNextOccurrence()
	}

	return reminder
}
