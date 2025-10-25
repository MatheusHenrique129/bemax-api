package services

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/google/uuid"
)

type RoleService interface {
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) (domain.Role, apierrors.RestError)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]domain.Role, apierrors.RestError)
}

type roleService struct {
	logger             ports.Logger
	roleRepository     ports.RoleRepository
	userRoleRepository ports.UserRoleRepository
}

func (s *roleService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) (domain.Role, apierrors.RestError) {
	res, err := s.roleRepository.FindByName(ctx, roleName)
	if err != nil {
		return domain.Role{}, apierrors.NewInternalServerRestError(err.Error(), err)
	}

	err = s.userRoleRepository.AssignRole(ctx, userID, res.ID)
	if err != nil {
		return domain.Role{}, apierrors.NewInternalServerRestError(err.Error(), err)
	}

	return res, nil
}

func (s *roleService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]domain.Role, apierrors.RestError) {
	res, err := s.userRoleRepository.FindRolesByUserID(ctx, userID)
	if err != nil {
		return nil, apierrors.NewInternalServerRestError("error find roles for user", err)
	}

	return res, nil
}

func NewRoleService(
	logger ports.Logger,
	roleRepository ports.RoleRepository,
	userRoleRepository ports.UserRoleRepository,
) RoleService {
	return &roleService{
		logger:             logger,
		roleRepository:     roleRepository,
		userRoleRepository: userRoleRepository,
	}
}
