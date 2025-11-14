package ports

import (
	"context"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type OAuthAccountRepository interface {
	Create(ctx context.Context, account *domain.OAuthAccount) error
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.OAuthAccount, error)
	FindByProviderAndUID(ctx context.Context, provider domain.OAuthProvider, providerUID string) (*domain.OAuthAccount, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.OAuthAccount, error)
	Update(ctx context.Context, account domain.OAuthAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider domain.OAuthProvider) error
}
