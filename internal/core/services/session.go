package services

import (
	"context"
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/google/uuid"
)

type SessionService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, deviceInfo, ipAddress, userAgent string) (*domain.Session, apierrors.RestError)
	ValidateSessionToken(ctx context.Context, sessionID, tokenJTI string) apierrors.RestError
	UpdateSessionToken(ctx context.Context, sessionID, tokenJTI string) apierrors.RestError
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, apierrors.RestError)
	TerminateSession(ctx context.Context, sessionID string) apierrors.RestError
	TerminateAllUserSessions(ctx context.Context, userID uuid.UUID) apierrors.RestError
}

type sessionService struct {
	logger      ports.Logger
	sessionRepo ports.SessionRepository
}

func (s *sessionService) CreateSession(ctx context.Context, userID uuid.UUID, deviceInfo, ipAddress, userAgent string) (*domain.Session, apierrors.RestError) {
	session := domain.NewSession(userID, deviceInfo, ipAddress, userAgent)

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		s.logger.Error("failed to create session", err)
		return nil, apierrors.NewInternalServerRestError("failed to create session", err)
	}

	s.logger.Info(fmt.Sprintf("session created for user %s: %s", userID, session.SessionID))
	return session, nil
}

func (s *sessionService) ValidateSessionToken(ctx context.Context, sessionID, tokenJTI string) apierrors.RestError {
	isLatest, err := s.sessionRepo.IsLatestAccessToken(ctx, sessionID, tokenJTI)
	if err != nil {
		s.logger.Error("error validating session token", err)
		return apierrors.NewInternalServerRestError("error validating session", err)
	}

	if !isLatest {
		s.logger.Warn(fmt.Sprintf("outdated access token used for session %s", sessionID))
		return apierrors.NewUnauthorizedRestError("access token has been replaced")
	}

	return nil
}

func (s *sessionService) UpdateSessionToken(ctx context.Context, sessionID, tokenJTI string) apierrors.RestError {
	if err := s.sessionRepo.UpdateLastAccessToken(ctx, sessionID, tokenJTI); err != nil {
		s.logger.Error("failed to update session token", err)
		return apierrors.NewInternalServerRestError("failed to update session", err)
	}

	return nil
}

func (s *sessionService) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, apierrors.RestError) {
	sessions, err := s.sessionRepo.FindActiveUserSessions(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user sessions", err)
		return nil, apierrors.NewInternalServerRestError("failed to get sessions", err)
	}

	return sessions, nil
}

func (s *sessionService) TerminateSession(ctx context.Context, sessionID string) apierrors.RestError {
	if err := s.sessionRepo.DeactivateSession(ctx, sessionID); err != nil {
		s.logger.Error("failed to terminate session", err)
		return apierrors.NewInternalServerRestError("failed to terminate session", err)
	}

	s.logger.Info(fmt.Sprintf("session terminated: %s", sessionID))
	return nil
}

func (s *sessionService) TerminateAllUserSessions(ctx context.Context, userID uuid.UUID) apierrors.RestError {
	if err := s.sessionRepo.DeactivateAllUserSessions(ctx, userID); err != nil {
		s.logger.Error("failed to terminate all user sessions", err)
		return apierrors.NewInternalServerRestError("failed to terminate sessions", err)
	}

	s.logger.Info(fmt.Sprintf("all sessions terminated for user: %s", userID))
	return nil
}

func NewSessionService(logger ports.Logger, sessionRepo ports.SessionRepository) SessionService {
	return &sessionService{
		logger:      logger,
		sessionRepo: sessionRepo,
	}
}
