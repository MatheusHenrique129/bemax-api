package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
)

type mysqlTokenRepository struct {
	BaseRepository
}

func (m mysqlTokenRepository) Save(ctx context.Context, token *auth.Token) error {
	query := `
		INSERT INTO tokens (
			id,
			user_id,
		    token,
		    token_type,
			expires_at,
		    created_at
		) VALUES (?, ?, ?, ?, ?, ?)
`

	row, err := m.dbClient.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.Type,
		token.ExpiresAt,
		time.Now().UTC(),
	)

	if err != nil {
		m.logger.Error("error inserting token from database", err)
		return fmt.Errorf("%s. %w", "failed to save token", err)
	}

	rows, err := row.RowsAffected()
	if err != nil || rows < 1 {
		m.logger.Error("no rows affected during create", err)
		return fmt.Errorf("%v. %w", ErrNoRowsAffected, err)
	}

	return nil
}

func NewMysqlTokenRepository(logger ports.Logger, dbClient *sql.DB) ports.TokenRepository {
	return &mysqlTokenRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
