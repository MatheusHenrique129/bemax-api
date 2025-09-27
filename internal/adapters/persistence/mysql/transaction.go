package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
)

var (
	ErrDBRollback     = errors.New("unable to rollback tx")
	ErrDBPanicRecover = errors.New("panic recover")
)

type Transaction interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type BaseRepository struct {
	dbClient *sql.DB
	logger   ports.Logger
}

func (m *BaseRepository) WithTransaction(ctx context.Context, fns ...ports.FnWithTx) (e error) {
	var transaction *sql.Tx
	{
		var err error
		if transaction, err = m.dbClient.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault}); err != nil {
			m.logger.Error(fmt.Sprintf("Error creating database transaction %v", err), err)
			return err
		}
	}

	defer func() {
		if r := recover(); r != nil {
			switch rt := r.(type) {
			case error:
				e = rt
			default:
				e = fmt.Errorf("%w. %v", ErrDBPanicRecover, r)
			}
		}
		if e != nil {
			if err := transaction.Rollback(); err != nil {
				e = fmt.Errorf("%w. %v", ErrDBRollback, err)
			}
		}
	}()

	for _, fn := range fns {
		if err := fn(ctx, transaction); err != nil {
			return err
		}
	}

	return transaction.Commit()
}

func NewBaseRepository(dbClient *sql.DB, logger ports.Logger) BaseRepository {
	return BaseRepository{
		dbClient: dbClient,
		logger:   logger,
	}
}
