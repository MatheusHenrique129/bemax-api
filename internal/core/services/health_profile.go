package services

import (
	"context"
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/google/uuid"
)

type HealthProfileService interface {
	GetOrCreateHealthProfile(ctx context.Context, userID uuid.UUID) (*domain.HealthProfile, apierrors.RestError)
	UpdateHealthProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateHealthProfileRequest) (*domain.HealthProfile, apierrors.RestError)
}

type healthProfileService struct {
	logger      ports.Logger
	profileRepo ports.HealthProfileRepository
}

func (s *healthProfileService) GetOrCreateHealthProfile(ctx context.Context, userID uuid.UUID) (*domain.HealthProfile, apierrors.RestError) {
	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err == nil && profile != nil {
		return profile, nil
	}

	newProfile := domain.NewHealthProfile(userID)
	if createErr := s.profileRepo.Create(ctx, newProfile); createErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create health profile for user %s", userID), createErr)
		return nil, apierrors.NewInternalServerRestError("failed to create health profile", createErr)
	}

	s.logger.Info(fmt.Sprintf("Health profile created for user: %s", userID))
	return newProfile, nil
}

func (s *healthProfileService) UpdateHealthProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateHealthProfileRequest) (*domain.HealthProfile, apierrors.RestError) {
	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("health profile not found for user: %s", userID), err)
		return nil, apierrors.NewNotFoundRestError("health profile not found")
	}

	profile.Update(req.BloodType, req.Height, req.Weight, req.Allergies, req.Medications, req.MedicalConditions, req.Notes)

	if updateErr := s.profileRepo.Update(ctx, profile); updateErr != nil {
		s.logger.Error(fmt.Sprintf("failed to update health profile for user: %s", userID), updateErr)
		return nil, apierrors.NewInternalServerRestError("failed to update health profile", updateErr)
	}

	s.logger.Info(fmt.Sprintf("Health profile updated for user: %s", userID))
	return profile, nil
}

func NewHealthProfileService(logger ports.Logger, profileRepo ports.HealthProfileRepository) HealthProfileService {
	return &healthProfileService{
		logger:      logger,
		profileRepo: profileRepo,
	}
}
