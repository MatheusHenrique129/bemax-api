package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
	ID          uuid.UUID
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
