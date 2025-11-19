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

type EmergencyContactService interface {
	CreateEmergencyContact(ctx context.Context, userID uuid.UUID, req dto.CreateEmergencyContactRequest) (*domain.EmergencyContact, apierrors.RestError)
	UpdateEmergencyContact(ctx context.Context, userID, contactID uuid.UUID, req dto.UpdateEmergencyContactRequest) (*domain.EmergencyContact, apierrors.RestError)
	DeleteEmergencyContact(ctx context.Context, userID, contactID uuid.UUID) apierrors.RestError
	GetEmergencyContactByID(ctx context.Context, userID, contactID uuid.UUID) (*domain.EmergencyContact, apierrors.RestError)
	GetUserEmergencyContacts(ctx context.Context, userID uuid.UUID) ([]domain.EmergencyContact, apierrors.RestError)
	SetPrimaryContact(ctx context.Context, userID, contactID uuid.UUID) apierrors.RestError
}

type emergencyContactService struct {
	logger      ports.Logger
	contactRepo ports.EmergencyContactRepository
}

func (s *emergencyContactService) CreateEmergencyContact(ctx context.Context, userID uuid.UUID, req dto.CreateEmergencyContactRequest) (*domain.EmergencyContact, apierrors.RestError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	contact := domain.NewEmergencyContact(userID, req.Name, req.Phone, req.Relationship)
	contact.Email = req.Email
	contact.Notes = req.Notes

	if req.IsPrimary {
		if err := s.contactRepo.UnsetAllPrimaryForUser(ctx, userID); err != nil {
			s.logger.Error(fmt.Sprintf("failed to unset primary contacts for user %s", userID), err)
			return nil, apierrors.NewInternalServerRestError("failed to update primary contacts", err)
		}
		contact.SetAsPrimary()
	}

	if err := s.contactRepo.Create(ctx, contact); err != nil {
		s.logger.Error(fmt.Sprintf("failed to create emergency contact for user %s", userID), err)
		return nil, apierrors.NewInternalServerRestError("failed to create emergency contact", err)
	}

	s.logger.Info(fmt.Sprintf("Emergency contact created: %s for user: %s", contact.ID, userID))
	return contact, nil
}

func (s *emergencyContactService) UpdateEmergencyContact(ctx context.Context, userID, contactID uuid.UUID, req dto.UpdateEmergencyContactRequest) (*domain.EmergencyContact, apierrors.RestError) {
	contact, err := s.contactRepo.FindByID(ctx, contactID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("emergency contact not found: %s", contactID), err)
		return nil, apierrors.NewNotFoundRestError("emergency contact not found")
	}

	if contact.UserID != userID {
		s.logger.Warn(fmt.Sprintf("user %s attempted to update contact %s owned by %s", userID, contactID, contact.UserID))
		return nil, apierrors.NewForbiddenRestError("you can only update your own contacts")
	}

	if req.IsPrimary != nil && *req.IsPrimary {
		if err := s.contactRepo.UnsetAllPrimaryForUser(ctx, userID); err != nil {
			s.logger.Error(fmt.Sprintf("failed to unset primary contacts for user %s", userID), err)
			return nil, apierrors.NewInternalServerRestError("failed to update primary contacts", err)
		}
	}

	isPrimary := contact.IsPrimary
	if req.IsPrimary != nil {
		isPrimary = *req.IsPrimary
	}

	contact.Update(req.Name, req.Phone, req.Email, req.Notes, domain.Address{}, req.Relationship, isPrimary)

	if err := s.contactRepo.Update(ctx, contact); err != nil {
		s.logger.Error(fmt.Sprintf("failed to update emergency contact: %s", contactID), err)
		return nil, apierrors.NewInternalServerRestError("failed to update emergency contact", err)
	}

	s.logger.Info(fmt.Sprintf("Emergency contact updated: %s", contactID))
	return contact, nil
}

func (s *emergencyContactService) DeleteEmergencyContact(ctx context.Context, userID, contactID uuid.UUID) apierrors.RestError {
	contact, err := s.contactRepo.FindByID(ctx, contactID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("emergency contact not found: %s", contactID), err)
		return apierrors.NewNotFoundRestError("emergency contact not found")
	}

	if contact.UserID != userID {
		s.logger.Warn(fmt.Sprintf("user %s attempted to delete contact %s owned by %s", userID, contactID, contact.UserID))
		return apierrors.NewForbiddenRestError("you can only delete your own contacts")
	}

	if err := s.contactRepo.Delete(ctx, contactID); err != nil {
		s.logger.Error(fmt.Sprintf("failed to delete emergency contact: %s", contactID), err)
		return apierrors.NewInternalServerRestError("failed to delete emergency contact", err)
	}

	s.logger.Info(fmt.Sprintf("Emergency contact deleted: %s", contactID))
	return nil
}

func (s *emergencyContactService) GetEmergencyContactByID(ctx context.Context, userID, contactID uuid.UUID) (*domain.EmergencyContact, apierrors.RestError) {
	contact, err := s.contactRepo.FindByID(ctx, contactID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("emergency contact not found: %s", contactID), err)
		return nil, apierrors.NewNotFoundRestError("emergency contact not found")
	}

	if contact.UserID != userID {
		s.logger.Warn(fmt.Sprintf("user %s attempted to access contact %s owned by %s", userID, contactID, contact.UserID))
		return nil, apierrors.NewForbiddenRestError("you can only access your own contacts")
	}

	return contact, nil
}

func (s *emergencyContactService) GetUserEmergencyContacts(ctx context.Context, userID uuid.UUID) ([]domain.EmergencyContact, apierrors.RestError) {
	contacts, err := s.contactRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to retrieve emergency contacts for user: %s", userID), err)
		return nil, apierrors.NewInternalServerRestError("failed to retrieve emergency contacts", err)
	}

	return contacts, nil
}

func (s *emergencyContactService) SetPrimaryContact(ctx context.Context, userID, contactID uuid.UUID) apierrors.RestError {
	contact, err := s.contactRepo.FindByID(ctx, contactID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("emergency contact not found: %s", contactID), err)
		return apierrors.NewNotFoundRestError("emergency contact not found")
	}

	if contact.UserID != userID {
		s.logger.Warn(fmt.Sprintf("user %s attempted to set primary for contact %s owned by %s", userID, contactID, contact.UserID))
		return apierrors.NewForbiddenRestError("you can only modify your own contacts")
	}

	if err := s.contactRepo.UnsetAllPrimaryForUser(ctx, userID); err != nil {
		s.logger.Error(fmt.Sprintf("failed to unset primary contacts for user %s", userID), err)
		return apierrors.NewInternalServerRestError("failed to update primary contacts", err)
	}

	contact.SetAsPrimary()
	if err := s.contactRepo.Update(ctx, contact); err != nil {
		s.logger.Error(fmt.Sprintf("failed to set primary contact: %s", contactID), err)
		return apierrors.NewInternalServerRestError("failed to set primary contact", err)
	}

	s.logger.Info(fmt.Sprintf("Emergency contact set as primary: %s", contactID))
	return nil
}

func NewEmergencyContactService(logger ports.Logger, contactRepo ports.EmergencyContactRepository) EmergencyContactService {
	return &emergencyContactService{
		logger:      logger,
		contactRepo: contactRepo,
	}
}
