package dto

import (
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/pkg/cpf"
	"github.com/MatheusHenrique129/bemax-api/pkg/hash"
	"github.com/google/uuid"
)

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
		ID:        uuid.New(),
		Email:     user.Email,
		FullName:  user.FullName,
		Password:  hashPassword,
		CPF:       newCPF,
		Phone:     user.Phone,
		BirthDate: birth,
		Status:    domain.UserStatusActive,
	}, nil
}
