package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *Role) Update(description string) {
	if description != "" {
		r.Description = description
	}
	r.UpdatedAt = time.Now().UTC()
}

func NewRole(name, description string) *Role {
	now := time.Now().UTC()
	return &Role{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
