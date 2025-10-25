package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core"
	auth "github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/google/uuid"
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

func (m mysqlTokenRepository) FindByToken(ctx context.Context, refreshToken string) (*auth.Token, error) {
	query := `
		SELECT 
			id,
			user_id,
			token,
			token_type,
			expires_at,
			created_at
		FROM tokens
		WHERE token = ? 
		  AND token_type = ?
		ORDER BY created_at DESC
		LIMIT 1
`

	var token auth.Token
	err := m.dbClient.QueryRowContext(ctx, query, refreshToken, core.TokenTypeRefresh).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.Type,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("token not found: %s", refreshToken), err)
		return nil, ErrTokenNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error finding token: %s", refreshToken), err)
		return nil, fmt.Errorf("error finding token: %w", err)
	}

	return &token, nil
}

func (m mysqlTokenRepository) RevokeToken(ctx context.Context, tokenString string) error {
	query := `DELETE FROM tokens
       WHERE token = ?
`

	result, err := m.dbClient.ExecContext(ctx, query, tokenString)
	if err != nil {
		m.logger.Error(fmt.Sprintf("error revoking token: %s", tokenString), err)
		return fmt.Errorf("error revoking token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking revoked rows: %w", err)
	}

	if rowsAffected == 0 {
		m.logger.Warn(fmt.Sprintf("no token found to revoke: %s", tokenString))
		return ErrTokenNotFound
	}

	m.logger.Info(fmt.Sprintf("token revoked successfully: %s", tokenString))
	return nil
}

func (m mysqlTokenRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM tokens 
       WHERE user_id = ? AND token_type = ?
`

	result, err := m.dbClient.ExecContext(ctx, query, userID, core.TokenTypeRefresh)
	if err != nil {
		m.logger.Error(fmt.Sprintf("error revoking all tokens for user: %s", userID), err)
		return fmt.Errorf("error revoking user tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking revoked rows: %w", err)
	}

	m.logger.Info(fmt.Sprintf("revoked %d tokens for user: %s", rowsAffected, userID))
	return nil
}

func (m mysqlTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM tokens 
       WHERE expires_at <= NOW()
`

	result, err := m.dbClient.ExecContext(ctx, query)
	if err != nil {
		m.logger.Error("error deleting expired tokens", err)
		return fmt.Errorf("error deleting expired tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}

	if rowsAffected > 0 {
		m.logger.Info(fmt.Sprintf("deleted %d expired tokens", rowsAffected))
	}

	return nil
}

func NewMysqlTokenRepository(logger ports.Logger, dbClient *sql.DB) ports.TokenRepository {
	return &mysqlTokenRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
