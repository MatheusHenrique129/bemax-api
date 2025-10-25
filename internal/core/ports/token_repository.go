package ports

import (
	"context"

	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type TokenRepository interface {
	Save(ctx context.Context, token *auth.Token) error
	FindByToken(ctx context.Context, refreshToken string) (*auth.Token, error)
	RevokeToken(ctx context.Context, tokenString string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
