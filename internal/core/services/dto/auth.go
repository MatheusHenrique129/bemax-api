package dto

import (
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto/validators"
	"github.com/MatheusHenrique129/bemax-api/pkg/cpf"
)

type UserRegisterRequest struct {
	Email     string `json:"email" binding:"required" validate:"required"`
	FullName  string `json:"full_name" binding:"required" validate:"required,min=3,max=100"`
	Password  string `json:"password" binding:"required" validate:"required,min=6,max=80"`
	CPF       string `json:"cpf" binding:"required" validate:"required"`
	Phone     string `json:"phone" binding:"required" validate:"required,min=10,max=14"`
	DateBirth string `json:"date_birth" binding:"required" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required,min=6,max=80"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" validate:"required"`
}

func (s UserRegisterRequest) Validate() (time.Time, apierrors.RestError) {
	var causes apierrors.CauseList

	if err := validators.ValidateEmail(s.Email); err != nil {
		causes = append(causes, apierrors.FieldError{
			Field:   "email",
			Message: err.Error(),
		})
	}

	if err := validators.ValidateName(s.FullName); err != nil {
		causes = append(causes, apierrors.FieldError{
			Field:   "full_name",
			Message: err.Error(),
		})
	}

	if err := cpf.ValidateCPF(s.CPF); err != nil {
		causes = append(causes, apierrors.FieldError{
			Field:   "cpf",
			Message: err.Error(),
		})
	}

	if err := validators.ValidatePhone(s.Phone); err != nil {
		causes = append(causes, apierrors.FieldError{
			Field:   "phone",
			Message: err.Error(),
		})
	}

	if err := validators.ValidatePassword(s.Password); err != nil {
		causes = append(causes, apierrors.FieldError{
			Field:   "password",
			Message: err.Error(),
		})
	}

	birth, err := validators.ValidateBirthDate(s.DateBirth)
	if err != nil {
		causes = append(causes, apierrors.FieldError{
			Field:   "date_birth",
			Message: err.Error(),
		})
	}

	if len(causes) > 0 {
		return time.Time{}, apierrors.NewBadRequestValidationRestError(
			"invalid user data",
			"validation_error",
			causes,
		)
	}

	return birth, nil
}

type UserRegisterResponse struct {
	Email     string    `json:"email" `
	FullName  string    `json:"full_name"`
	CPF       string    `json:"cpf"`
	Phone     string    `json:"phone"`
	DateBirth time.Time `json:"date_birth"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    time.Duration `json:"expires_in"`
}

type GetTokenResponse struct {
	Token     string        `json:"token"`
	TokenJTI  string        `json:"token_jti"`
	Timestamp time.Time     `json:"timestamp"`
	ExpireAt  time.Duration `json:"expire_at"`
}
