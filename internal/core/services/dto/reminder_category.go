package dto

import (
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Color       string `json:"color,omitempty"`
}

type UpdateCategoryRequest struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Color        string `json:"color,omitempty"`
	DisplayOrder int    `json:"display_order,omitempty"`
}

type CategoryResponse struct {
	ID           uuid.UUID            `json:"id"`
	Name         string               `json:"name"`
	NameKey      string               `json:"name_key,omitempty"`
	Description  string               `json:"description,omitempty"`
	Icon         string               `json:"icon,omitempty"`
	Color        string               `json:"color,omitempty"`
	Scope        domain.CategoryScope `json:"scope"`
	DisplayOrder int                  `json:"display_order"`
}

func (r *CreateCategoryRequest) Validate() apierrors.RestError {
	if r.Name == "" {
		return apierrors.NewBadRequestRestError("name is required")
	}
	return nil
}
