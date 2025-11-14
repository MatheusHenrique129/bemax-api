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
	"github.com/MatheusHenrique129/bemax-api/pkg/cpf"
	"github.com/google/uuid"
)

var (
	ErrUserInactive         = errors.New("user is inactive")
	ErrCPFAlreadyExists     = errors.New("cpf already exists")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTooManyLoginAttempts = errors.New("too many login attempts, try again later")
)

type UserService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, apierrors.RestError)
	CreateUser(ctx context.Context, dataReq dto.UserRegisterRequest) (domain.User, apierrors.RestError)
	CreateUserOAuth(ctx context.Context, user domain.User) (domain.User, apierrors.RestError)
	AuthenticateUser(ctx context.Context, email, password string, ipAddress, userAgent string) (*domain.User, apierrors.RestError)
	UpdateLastLogin(ctx context.Context, user domain.User) apierrors.RestError

	GetTokenVersion(ctx context.Context, userID uuid.UUID) (int, apierrors.RestError)
	IncrementTokenVersion(ctx context.Context, userID uuid.UUID) apierrors.RestError
}

type userService struct {
	logger         ports.Logger
	userRepository ports.UserRepository
	roleService    RoleService
}

func (u userService) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, apierrors.RestError) {
	user, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mysql.ErrUserNotFound) {
			u.logger.Error(fmt.Sprintf("user not found with id %s.", id), err)
			return domain.User{}, apierrors.NewNotFoundRestError(err.Error())
		}

		u.logger.Error(fmt.Sprintf("error getting user with id %s.", id), err)
		return domain.User{}, apierrors.NewInternalServerRestError(fmt.Sprintf("error finding user by id: %s.", id), err)
	}

	roles, errRoles := u.roleService.GetUserRoles(ctx, user.ID)
	if errRoles != nil {
		return domain.User{}, errRoles
	}

	user.Roles = roles
	return user, nil
}

func (u userService) CreateUser(ctx context.Context, userRegister dto.UserRegisterRequest) (domain.User, apierrors.RestError) {
	formatCPF := cpf.FormatCPF(userRegister.CPF)

	res, err := u.userRepository.FindByCPF(ctx, formatCPF)
	if err != nil {
		if !errors.Is(err, mysql.ErrUserNotFound) {
			u.logger.Error(fmt.Sprintf("error getting user with cpf %s.", userRegister.CPF), err)
			return domain.User{}, apierrors.NewInternalServerRestError(fmt.Sprintf("error finding user by cpf: %s.", userRegister.CPF), err)
		}
	}

	if res.CPF != "" {
		var causes apierrors.CauseList
		causes = append(causes, apierrors.CauseList{
			apierrors.FieldError{
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

func (u userService) CreateUserOAuth(ctx context.Context, user domain.User) (domain.User, apierrors.RestError) {
	res, err := u.userRepository.FindByEmail(ctx, user.Email)
	if err != nil {
		if !errors.Is(err, mysql.ErrUserNotFound) {
			u.logger.Error(fmt.Sprintf("error getting user with email %s.", user.Email), err)
			return domain.User{}, apierrors.NewInternalServerRestError(fmt.Sprintf("error finding user by email: %s.", user.Email), err)
		}
	}

	if res.Email != "" {
		var causes apierrors.CauseList
		causes = append(causes, apierrors.CauseList{
			apierrors.FieldError{
				Field:   "user_already_exists",
				Message: ErrEmailAlreadyExists.Error(),
			},
		})

		return domain.User{}, apierrors.NewConflictRestError(
			"User with this email already exists",
			causes,
		)
	}

	txErr := u.userRepository.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err = u.userRepository.Create(ctx, user); err != nil {
			return err
		}

		role, err := u.roleService.AssignRoleToUser(ctx, user.ID, "ADMIN")
		if err != nil {
			return err
		}

		user.Roles = append(user.Roles, role)
		return nil
	})
	if txErr != nil {
		return domain.User{}, apierrors.NewInternalServerRestError("An error occurred while trying to create the user", txErr)
	}

	return user, nil
}

func (u userService) AuthenticateUser(ctx context.Context, email, password string, ipAddress, userAgent string) (*domain.User, apierrors.RestError) {
	blocked, err := u.checkRateLimit(ctx, email, 5)
	if err != nil {
		return nil, apierrors.NewInternalServerRestError("error checking rate limit", err)
	}

	if blocked {
		if recordErr := u.recordFailedAttempt(ctx, email, ipAddress, userAgent); recordErr != nil {
			u.logger.Error("critical: failed to record failed attempt", recordErr)
			return nil, apierrors.NewInternalServerRestError("authentication system error", recordErr)
		}
		return nil, apierrors.NewTooManyRequestsRestError(ErrTooManyLoginAttempts.Error())
	}

	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error getting user with email %s.", email), err)
		if errors.Is(err, mysql.ErrUserNotFound) {
			if recordErr := u.recordFailedAttempt(ctx, email, ipAddress, userAgent); recordErr != nil {
				u.logger.Error("critical: failed to record failed attempt", recordErr)
				return nil, apierrors.NewInternalServerRestError("authentication system error", recordErr)
			}
			return nil, apierrors.NewNotFoundRestError(err.Error())
		}

		u.logger.Error(fmt.Sprintf("error finding user by email %s", email), err)
		return nil, apierrors.NewInternalServerRestError("error finding user", err)
	}

	if user.IsOAuthAuth() {
		if recordErr := u.recordFailedAttempt(ctx, email, ipAddress, userAgent); recordErr != nil {
			u.logger.Error("critical: failed to record failed attempt", recordErr)
		}
		return nil, apierrors.NewUnauthorizedRestError("this account uses OAuth authentication, please complete your details to login.")
	}

	roles, errRoles := u.roleService.GetUserRoles(ctx, user.ID)
	if errRoles != nil {
		return nil, errRoles
	}
	user.Roles = roles

	if !user.IsActive() {
		if recordErr := u.recordFailedAttempt(ctx, email, ipAddress, userAgent); recordErr != nil {
			u.logger.Error("critical: failed to record failed attempt", recordErr)
			return nil, apierrors.NewInternalServerRestError("authentication system error", recordErr)
		}

		return nil, apierrors.NewUnauthorizedRestError(ErrUserInactive.Error())
	}

	if err := user.CheckPassword(password); err != nil {
		if recordErr := u.recordFailedAttempt(ctx, email, ipAddress, userAgent); recordErr != nil {
			u.logger.Error("critical: failed to record failed attempt", recordErr)
			return nil, apierrors.NewInternalServerRestError("authentication system error", recordErr)
		}

		return nil, apierrors.NewUnauthorizedRestError(ErrInvalidCredentials.Error())
	}

	if err := u.userRepository.UpdateLastLogin(ctx, user.ID); err != nil {
		u.logger.Error(err.Error(), err)
		return nil, apierrors.NewInternalServerRestError("error update last login", err)
	}

	user.UpdateLastLogin()

	if recordErr := u.userRepository.RecordLoginAttempt(ctx, email, true, ipAddress, userAgent); recordErr != nil {
		u.logger.Error("critical: failed to record failed attempt", recordErr)
		return nil, apierrors.NewInternalServerRestError("authentication system error", recordErr)

	}

	return &user, nil
}

func (u userService) UpdateLastLogin(ctx context.Context, user domain.User) apierrors.RestError {
	if err := u.userRepository.UpdateLastLogin(ctx, user.ID); err != nil {
		u.logger.Error(err.Error(), err)
		return apierrors.NewInternalServerRestError("error update last login", err)
	}

	return nil

}

func (u userService) GetTokenVersion(ctx context.Context, userID uuid.UUID) (int, apierrors.RestError) {
	version, err := u.userRepository.GetTokenVersion(ctx, userID)
	if err != nil {
		if errors.Is(err, mysql.ErrUserNotFound) {
			return 0, apierrors.NewNotFoundRestError("user not found")
		}
		u.logger.Error(fmt.Sprintf("error getting token version for user %s", userID), err)
		return 0, apierrors.NewInternalServerRestError("error getting token version", err)
	}

	return version, nil
}

func (u userService) IncrementTokenVersion(ctx context.Context, userID uuid.UUID) apierrors.RestError {
	err := u.userRepository.IncrementTokenVersion(ctx, userID)
	if err != nil {
		if errors.Is(err, mysql.ErrUserNotFound) {
			return apierrors.NewNotFoundRestError("user not found")
		}
		u.logger.Error(fmt.Sprintf("error incrementing token version for user %s", userID), err)
		return apierrors.NewInternalServerRestError("error incrementing token version", err)
	}

	return nil
}

func (u userService) checkRateLimit(ctx context.Context, email string, minutes int) (bool, error) {
	attempts, err := u.userRepository.GetLoginAttempts(ctx, email, minutes)
	if err != nil && !errors.Is(err, mysql.ErrLoginNotFound) {
		return false, err
	}

	switch {
	case minutes <= 1 && attempts >= 3:
		return true, nil
	case minutes <= 5 && attempts >= 5:
		return true, nil
	case minutes <= 15 && attempts >= 10:
		return true, nil
	default:
		return false, nil
	}
}

func (u userService) recordFailedAttempt(ctx context.Context, email, ipAddress, userAgent string) error {
	// TODO - As a possible improvement, implement a queue using local SQLite that has no financial cost and avoids returning an error to the user if it fails to insert into the table.
	if err := u.userRepository.RecordLoginAttempt(ctx, email, false, ipAddress, userAgent); err != nil {
		u.logger.Error("failed to record failed login attempt", err)
		return fmt.Errorf("failed to record login attempt: %w", err)
	}
	return nil
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
