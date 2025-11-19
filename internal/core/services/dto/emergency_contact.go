package dto

import (
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type CreateEmergencyContactRequest struct {
	Name         string                     `json:"name" validate:"required,min=3,max=255"`
	Relationship domain.ContactRelationship `json:"relationship" validate:"required"`
	Phone        string                     `json:"phone" validate:"required"`
	Email        string                     `json:"email,omitempty"`
	Address      string                     `json:"address,omitempty"`
	Notes        string                     `json:"notes,omitempty"`
	IsPrimary    bool                       `json:"is_primary"`
}

type UpdateEmergencyContactRequest struct {
	Name         string                     `json:"name,omitempty"`
	Relationship domain.ContactRelationship `json:"relationship,omitempty"`
	Phone        string                     `json:"phone,omitempty"`
	Email        string                     `json:"email,omitempty"`
	Address      string                     `json:"address,omitempty"`
	Notes        string                     `json:"notes,omitempty"`
	IsPrimary    *bool                      `json:"is_primary,omitempty"`
}

type EmergencyContactResponse struct {
	ID           uuid.UUID                  `json:"id"`
	Name         string                     `json:"name"`
	Relationship domain.ContactRelationship `json:"relationship"`
	Phone        string                     `json:"phone"`
	Email        string                     `json:"email,omitempty"`
	Address      string                     `json:"address,omitempty"`
	Notes        string                     `json:"notes,omitempty"`
	IsPrimary    bool                       `json:"is_primary"`
	IsActive     bool                       `json:"is_active"`
}

func (r *CreateEmergencyContactRequest) Validate() apierrors.RestError {
	if r.Name == "" {
		return apierrors.NewBadRequestRestError("name is required")
	}
	if r.Phone == "" {
		return apierrors.NewBadRequestRestError("phone is required")
	}
	if r.Relationship == "" {
		return apierrors.NewBadRequestRestError("relationship is required")
	}
	return nil
}
