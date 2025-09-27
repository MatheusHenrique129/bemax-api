package ports

import (
	"context"
	"database/sql"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
)

type FnWithTx func(ctx context.Context, tx *sql.Tx) error

type UserRepository interface {
	WithTransaction(ctx context.Context, fns ...FnWithTx) error
	FindByCPF(ctx context.Context, cpf string) (domain.User, error)
	Create(ctx context.Context, user domain.User) error
}
