package domain

import (
	"time"

	"github.com/google/uuid"
)

// CategoryScope defines if category is system or user-defined
type CategoryScope string

const (
	CategoryScopeSystem CategoryScope = "system" // Pre-defined by system
	CategoryScopeUser   CategoryScope = "user"   // Created by user
)

// ReminderCategory represents a category for reminders
type ReminderCategory struct {
	ID           uuid.UUID     `json:"id"`
	UserID       *uuid.UUID    `json:"user_id,omitempty"` // NULL for system categories
	Name         string        `json:"name"`
	NameKey      string        `json:"name_key"` // For i18n: "category.medication"
	Description  string        `json:"description,omitempty"`
	Icon         string        `json:"icon,omitempty"`  // Icon name or emoji
	Color        string        `json:"color,omitempty"` // Hex color: "#FF5733"
	Scope        CategoryScope `json:"scope"`
	DisplayOrder int           `json:"display_order"`
	IsActive     bool          `json:"is_active"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// IsSystemCategory checks if category is system-defined
func (c *ReminderCategory) IsSystemCategory() bool {
	return c.Scope == CategoryScopeSystem
}

// IsUserCategory checks if category is user-defined
func (c *ReminderCategory) IsUserCategory() bool {
	return c.Scope == CategoryScopeUser
}

// Update updates category information
func (c *ReminderCategory) Update(name, description, icon, color string, displayOrder int) {
	if name != "" {
		c.Name = name
	}
	if description != "" {
		c.Description = description
	}
	if icon != "" {
		c.Icon = icon
	}
	if color != "" {
		c.Color = color
	}
	if displayOrder > 0 {
		c.DisplayOrder = displayOrder
	}
	c.UpdatedAt = time.Now().UTC()
}

// Deactivate marks category as inactive
func (c *ReminderCategory) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now().UTC()
}

// Activate marks category as active
func (c *ReminderCategory) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now().UTC()
}

// NewSystemCategory creates a new system category
func NewSystemCategory(name, nameKey, description, icon, color string, displayOrder int) *ReminderCategory {
	now := time.Now().UTC()
	return &ReminderCategory{
		ID:           uuid.New(),
		Name:         name,
		NameKey:      nameKey,
		Description:  description,
		Icon:         icon,
		Color:        color,
		Scope:        CategoryScopeSystem,
		DisplayOrder: displayOrder,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewUserCategory creates a new user-defined category
func NewUserCategory(userID uuid.UUID, name, description, icon, color string) *ReminderCategory {
	now := time.Now().UTC()
	return &ReminderCategory{
		ID:           uuid.New(),
		UserID:       &userID,
		Name:         name,
		Description:  description,
		Icon:         icon,
		Color:        color,
		Scope:        CategoryScopeUser,
		DisplayOrder: 999,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
