package domain

import (
	"time"

	"github.com/google/uuid"
)

// BloodType defines blood type options
type BloodType string

const (
	BloodTypeAPositive  BloodType = "A+"
	BloodTypeANegative  BloodType = "A-"
	BloodTypeBPositive  BloodType = "B+"
	BloodTypeBNegative  BloodType = "B-"
	BloodTypeABPositive BloodType = "AB+"
	BloodTypeABNegative BloodType = "AB-"
	BloodTypeOPositive  BloodType = "O+"
	BloodTypeONegative  BloodType = "O-"
	BloodTypeUnknown    BloodType = "unknown"
)

// HealthProfile represents user's health information
type HealthProfile struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	BloodType         BloodType `json:"blood_type"`
	Height            *float64  `json:"height,omitempty"`
	Weight            *float64  `json:"weight,omitempty"`
	Allergies         []string  `json:"allergies,omitempty"`
	Medications       []string  `json:"medications,omitempty"`
	MedicalConditions []string  `json:"medical_conditions,omitempty"`
	Notes             string    `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Update updates health profile information
func (h *HealthProfile) Update(bloodType BloodType, height, weight *float64, allergies, medications, medicalConditions []string, notes string) {
	if bloodType != "" {
		h.BloodType = bloodType
	}
	if height != nil {
		h.Height = height
	}
	if weight != nil {
		h.Weight = weight
	}
	if allergies != nil {
		h.Allergies = allergies
	}
	if medications != nil {
		h.Medications = medications
	}
	if medicalConditions != nil {
		h.MedicalConditions = medicalConditions
	}
	if notes != "" {
		h.Notes = notes
	}
	h.UpdatedAt = time.Now().UTC()
}

// NewHealthProfile creates a new health profile
func NewHealthProfile(userID uuid.UUID) *HealthProfile {
	now := time.Now().UTC()
	return &HealthProfile{
		ID:        uuid.New(),
		UserID:    userID,
		BloodType: BloodTypeUnknown,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
