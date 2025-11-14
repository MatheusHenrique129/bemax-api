package dto

import (
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/pkg/cpf"
	"github.com/MatheusHenrique129/bemax-api/pkg/hash"
	"github.com/google/uuid"
)

// NewUser creates a new local user (with password)
func NewUser(user UserRegisterRequest) (domain.User, apierrors.RestError) {
	birth, validationErr := user.Validate()
	if validationErr != nil {
		return domain.User{}, validationErr
	}

	hashPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, apierrors.NewInternalServerRestError(err.Error(), err)
	}

	newCPF := cpf.FormatCPF(user.CPF)

	return domain.User{
		ID:               uuid.New(),
		Email:            user.Email,
		FullName:         user.FullName,
		Password:         hashPassword,
		CPF:              newCPF,
		Phone:            user.Phone,
		BirthDate:        &birth,
		Status:           domain.UserStatusActive,
		AuthProvider:     domain.AuthProviderLocal,
		EmailVerified:    false,
		PhoneVerified:    false,
		ProfileCompleted: true,
		TokenVersion:     0,
	}, nil
}

// NewOAuthUser creates a new OAuth user (without password, CPF optional)
func NewOAuthUser(email, fullName string, emailVerified bool) (domain.User, apierrors.RestError) {
	// TODO implement Validators in email
	return domain.User{
		ID:               uuid.New(),
		Email:            email,
		AuthProvider:     domain.AuthProviderOAuth,
		FullName:         fullName,
		EmailVerified:    emailVerified,
		PhoneVerified:    false,
		ProfileCompleted: false,
		Status:           domain.UserStatusActive,
		TokenVersion:     0,
	}, nil
}
