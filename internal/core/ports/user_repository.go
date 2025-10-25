package ports

import (
	"context"
	"database/sql"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/google/uuid"
)

type FnWithTx func(ctx context.Context, tx *sql.Tx) error

type UserRepository interface {
	WithTransaction(ctx context.Context, fns ...FnWithTx) error
	Create(ctx context.Context, user domain.User) error
	FindByCPF(ctx context.Context, cpf string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error

	GetLoginAttempts(ctx context.Context, email string, minutes int) (int, error)
	RecordLoginAttempt(ctx context.Context, email string, success bool, ipAddress, userAgent string) error
}
