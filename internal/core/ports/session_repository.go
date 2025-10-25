package ports

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type SessionRepository interface {
	// Session management
	CreateSession(ctx context.Context, session *domain.Session) error
	UpdateSession(ctx context.Context, session *domain.Session) error
	FindBySessionID(ctx context.Context, sessionID string) (*domain.Session, error)
	FindActiveUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)

	// Access token validation
	IsLatestAccessToken(ctx context.Context, sessionID, tokenJTI string) (bool, error)
	UpdateLastAccessToken(ctx context.Context, sessionID, tokenJTI string) error

	// Session termination
	DeactivateSession(ctx context.Context, sessionID string) error
	DeactivateAllUserSessions(ctx context.Context, userID uuid.UUID) error

	// Cleanup
	DeleteExpiredSessions(ctx context.Context) error
}
