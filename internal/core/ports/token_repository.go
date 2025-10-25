package ports

import (
	"context"

	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
)

type TokenRepository interface {
	Save(ctx context.Context, token *auth.Token) error
}
