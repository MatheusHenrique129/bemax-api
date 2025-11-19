package dto

import (
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
)

// UserProfileResponse represents the complete user profile for mobile app
type UserProfileResponse struct {
	User              *domain.User               `json:"user"`
	HealthProfile     *HealthProfileResponse     `json:"health_profile,omitempty"`
	EmergencyContacts []EmergencyContactResponse `json:"emergency_contacts,omitempty"`
	Reminders         []ReminderResponse         `json:"reminders,omitempty"`
	Stats             *ProfileStats              `json:"stats,omitempty"`
}

// ProfileStats provides quick statistics for the user
type ProfileStats struct {
	TotalReminders    int `json:"total_reminders"`
	ActiveReminders   int `json:"active_reminders"`
	UpcomingReminders int `json:"upcoming_reminders"`
	TodayReminders    int `json:"today_reminders"`
}
