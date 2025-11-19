package domain

import (
	"time"

	"github.com/google/uuid"
)

// ContactRelationship defines the relationship with the emergency contact
type ContactRelationship string

const (
	RelationshipSpouse    ContactRelationship = "spouse"
	RelationshipParent    ContactRelationship = "parent"
	RelationshipChild     ContactRelationship = "child"
	RelationshipSibling   ContactRelationship = "sibling"
	RelationshipFriend    ContactRelationship = "friend"
	RelationshipDoctor    ContactRelationship = "doctor"
	RelationshipCaregiver ContactRelationship = "caregiver"
	RelationshipOther     ContactRelationship = "other"
)

// EmergencyContact represents an emergency contact for a user
type EmergencyContact struct {
	ID           uuid.UUID           `json:"id"`
	UserID       uuid.UUID           `json:"user_id"`
	Name         string              `json:"name"`
	Relationship ContactRelationship `json:"relationship"`
	Phone        string              `json:"phone"`
	Email        string              `json:"email,omitempty"`
	Address      Address             `json:"address,omitempty"`
	Notes        string              `json:"notes,omitempty"`
	IsPrimary    bool                `json:"is_primary"`
	IsActive     bool                `json:"is_active"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// Update updates emergency contact information
func (e *EmergencyContact) Update(name, phone, email, notes string, address Address, relationship ContactRelationship, isPrimary bool) {
	if name != "" {
		e.Name = name
	}
	if phone != "" {
		e.Phone = phone
	}
	if email != "" {
		e.Email = email
	}
	e.Address = address

	if notes != "" {
		e.Notes = notes
	}
	if relationship != "" {
		e.Relationship = relationship
	}
	e.IsPrimary = isPrimary
	e.UpdatedAt = time.Now().UTC()
}

// Deactivate marks contact as inactive
func (e *EmergencyContact) Deactivate() {
	e.IsActive = false
	e.UpdatedAt = time.Now().UTC()
}

// Activate marks contact as active
func (e *EmergencyContact) Activate() {
	e.IsActive = true
	e.UpdatedAt = time.Now().UTC()
}

// SetAsPrimary sets this contact as the primary contact
func (e *EmergencyContact) SetAsPrimary() {
	e.IsPrimary = true
	e.UpdatedAt = time.Now().UTC()
}

// NewEmergencyContact creates a new emergency contact
func NewEmergencyContact(userID uuid.UUID, name, phone string, relationship ContactRelationship) *EmergencyContact {
	now := time.Now().UTC()
	return &EmergencyContact{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         name,
		Phone:        phone,
		Relationship: relationship,
		IsPrimary:    false,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
