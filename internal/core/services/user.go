package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/adapters/persistence/mysql"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
)

var (
	ErrCPFAlreadyExists = errors.New("cpf already exists")
)

type UserService interface {
	CreateUser(ctx context.Context, dataReq dto.UserRegisterRequest) (domain.User, apierrors.RestError)
}

type userService struct {
	logger         ports.Logger
	userRepository ports.UserRepository
	roleService    RoleService
}

func (u userService) CreateUser(ctx context.Context, userRegister dto.UserRegisterRequest) (domain.User, apierrors.RestError) {
	res, err := u.userRepository.FindByCPF(ctx, userRegister.CPF)
	if err != nil {
		if !errors.Is(err, mysql.ErrUserNotFound) {
			u.logger.Error(fmt.Sprintf("error getting user with cpf %s.", userRegister.CPF), err)
			return domain.User{}, apierrors.NewInternalServerRestError(fmt.Sprintf("error finding user by cpf: %s.", userRegister.CPF), err)
		}
	}

	if res.CPF != "" {
		var causes apierrors.CauseList
		causes = append(causes, apierrors.CauseList{
			dto.FieldError{
				Field:   "user_already_exists",
				Message: ErrCPFAlreadyExists.Error(),
			},
		})

		return domain.User{}, apierrors.NewConflictRestError(
			"User with this CPF already exists",
			causes,
		)
	}

	domainUser, errDomain := dto.NewUser(userRegister)
	if errDomain != nil {
		return domain.User{}, errDomain
	}

	txErr := u.userRepository.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err = u.userRepository.Create(ctx, domainUser); err != nil {
			return err
		}

		role, err := u.roleService.AssignRoleToUser(ctx, domainUser.ID, "ADMIN")
		if err != nil {
			return err
		}

		domainUser.Roles = append(domainUser.Roles, role)
		return nil
	})
	if txErr != nil {
		return domain.User{}, apierrors.NewInternalServerRestError("An error occurred while trying to create the user", txErr)
	}

	return domainUser, nil
}

func NewUserService(
	logger ports.Logger,
	userRepository ports.UserRepository,
	roleService RoleService,

) UserService {
	return &userService{
		logger:         logger,
		userRepository: userRepository,
		roleService:    roleService,
	}
}
