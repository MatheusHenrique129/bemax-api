package dto

import (
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type UpdateHealthProfileRequest struct {
	BloodType         domain.BloodType `json:"blood_type,omitempty"`
	Height            *float64         `json:"height,omitempty"`
	Weight            *float64         `json:"weight,omitempty"`
	Allergies         []string         `json:"allergies,omitempty"`
	Medications       []string         `json:"medications,omitempty"`
	MedicalConditions []string         `json:"medical_conditions,omitempty"`
	Notes             string           `json:"notes,omitempty"`
}

type HealthProfileResponse struct {
	ID                uuid.UUID        `json:"id"`
	BloodType         domain.BloodType `json:"blood_type"`
	Height            *float64         `json:"height,omitempty"`
	Weight            *float64         `json:"weight,omitempty"`
	Allergies         []string         `json:"allergies,omitempty"`
	Medications       []string         `json:"medications,omitempty"`
	MedicalConditions []string         `json:"medical_conditions,omitempty"`
	Notes             string           `json:"notes,omitempty"`
}
