package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrOAuthAccountNotFound = errors.New("oauth account not found")
)

type mysqlOAuthAccountRepository struct {
	BaseRepository
}

func (m *mysqlOAuthAccountRepository) Create(ctx context.Context, account *domain.OAuthAccount) error {
	query := `
		INSERT INTO oauth_accounts (
			id, user_id, provider, provider_uid, firebase_uid, provider_email, 
			provider_name, provider_picture, email_verified, expires_at, 
		    last_login_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	now := time.Now().UTC()
	_, err := m.dbClient.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.Provider,
		account.ProviderUID,
		account.FirebaseUID,
		account.ProviderEmail,
		account.ProviderName,
		account.ProviderPicture,
		account.EmailVerified,
		account.ExpiresAt,
		account.LastLoginAt,
		now,
		now,
	)

	if err != nil {
		m.logger.Error("failed to create OAuth account", err)
		return fmt.Errorf("failed to create oauth account: %w", err)
	}

	return nil
}

func (m *mysqlOAuthAccountRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.OAuthAccount, error) {
	query := `
		SELECT 
		    id, 
		    user_id,
		    provider,
		    provider_uid,
		    firebase_uid,
		    provider_email,
		    provider_name,
		    provider_picture,
		    email_verified,
		    expires_at,
		    last_login_at,
		    created_at,
		    updated_at
		FROM oauth_accounts
		WHERE firebase_uid = ?
`

	var account domain.OAuthAccount
	var providerPicture sql.NullString
	var lastLoginAt, expiresAt sql.NullTime

	err := m.dbClient.QueryRowContext(ctx, query, firebaseUID).Scan(
		&account.ID,
		&account.UserID,
		&account.Provider,
		&account.ProviderUID,
		&account.FirebaseUID,
		&account.ProviderEmail,
		&account.ProviderName,
		&providerPicture,
		&account.EmailVerified,
		&expiresAt,
		&lastLoginAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if providerPicture.Valid {
		account.ProviderPicture = providerPicture.String
	}
	if lastLoginAt.Valid {
		account.LastLoginAt = &lastLoginAt.Time
	}
	if expiresAt.Valid {
		account.ExpiresAt = &expiresAt.Time
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("error finding oauth account by firebase_uid. %s", firebaseUID), err)
		return nil, ErrOAuthAccountNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error finding oauth account by firebase_uid: %s", firebaseUID), err)
		return nil, fmt.Errorf("error finding oauth account: %w", err)
	}

	return &account, nil
}

func (m *mysqlOAuthAccountRepository) FindByProviderAndUID(ctx context.Context, provider domain.OAuthProvider, providerUID string) (*domain.OAuthAccount, error) {
	query := `
		SELECT 
		    id, 
		    user_id,
		    provider,
		    provider_uid,
		    firebase_uid,
		    provider_email,
		    provider_name,
		    provider_picture,
		    email_verified,
		    expires_at,
		    last_login_at,
		    created_at,
		    updated_at
		FROM oauth_accounts
		WHERE provider = ? AND provider_uid = ?
	`

	var account domain.OAuthAccount
	var providerPicture sql.NullString
	var lastLoginAt, expiresAt sql.NullTime

	err := m.dbClient.QueryRowContext(ctx, query, provider, providerUID).Scan(
		&account.ID,
		&account.UserID,
		&account.Provider,
		&account.ProviderUID,
		&account.FirebaseUID,
		&account.ProviderEmail,
		&account.ProviderName,
		&providerPicture,
		&account.EmailVerified,
		&expiresAt,
		&lastLoginAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("error finding oauth account by provider: %s, uid: %s", provider, providerUID), err)
		return nil, ErrOAuthAccountNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error finding oauth account by provider: %s, uid: %s", provider, providerUID), err)
		return nil, fmt.Errorf("error finding oauth account: %w", err)
	}

	if providerPicture.Valid {
		account.ProviderPicture = providerPicture.String
	}
	if lastLoginAt.Valid {
		account.LastLoginAt = &lastLoginAt.Time
	}
	if expiresAt.Valid {
		account.ExpiresAt = &expiresAt.Time
	}

	return &account, nil
}

func (m *mysqlOAuthAccountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.OAuthAccount, error) {
	res := make([]domain.OAuthAccount, 0)

	query := `
		SELECT 
			id, 
			user_id,
			provider,
			provider_uid,
			firebase_uid,
			provider_email,
			provider_name,
			provider_picture,
			email_verified,
			expires_at,
			last_login_at,
			created_at,
			updated_at
		FROM oauth_accounts
		WHERE user_id = ?
`

	rows, err := m.dbClient.QueryContext(ctx, query, userID)
	if err != nil {
		m.logger.Error(fmt.Sprintf("error finding oauth accounts for user: %s", userID), err)
		return nil, fmt.Errorf("error finding oauth accounts: %w", err)

	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		var account domain.OAuthAccount
		var providerPicture sql.NullString
		var lastLoginAt, expiresAt sql.NullTime

		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Provider,
			&account.ProviderUID,
			&account.FirebaseUID,
			&account.ProviderEmail,
			&account.ProviderName,
			&providerPicture,
			&account.EmailVerified,
			&expiresAt,
			&lastLoginAt,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			m.logger.Error("error scanning oauth account", err)
			return nil, err
		}

		if providerPicture.Valid {
			account.ProviderPicture = providerPicture.String
		}
		if lastLoginAt.Valid {
			account.LastLoginAt = &lastLoginAt.Time
		}
		if expiresAt.Valid {
			account.ExpiresAt = &expiresAt.Time
		}

		res = append(res, account)
	}

	return res, nil
}

func (m *mysqlOAuthAccountRepository) Update(ctx context.Context, account domain.OAuthAccount) error {
	query := `
		UPDATE oauth_accounts 
		SET provider_email = ?,
		    provider_name = ?,
		    provider_picture = ?,
		    email_verified = ?,
		    expires_at = ?,
		    last_login_at = ?,
		    updated_at = ?
		WHERE id = ?
	`

	res, err := m.dbClient.ExecContext(ctx, query,
		account.ProviderEmail,
		account.ProviderName,
		account.ProviderPicture,
		account.EmailVerified,
		account.ExpiresAt,
		account.LastLoginAt,
		time.Now().UTC(),
		account.ID,
	)

	if err != nil {
		m.logger.Error("error updating oauth account", err)
		return fmt.Errorf("failed to update oauth account: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrOAuthAccountNotFound
	}

	return nil
}

func (m *mysqlOAuthAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM oauth_accounts WHERE id = ?`

	res, err := m.dbClient.ExecContext(ctx, query, id)
	if err != nil {
		m.logger.Error(fmt.Sprintf("error deleting oauth account: %s", id), err)
		return fmt.Errorf("failed to delete oauth account: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrOAuthAccountNotFound
	}

	return nil
}

func (m *mysqlOAuthAccountRepository) DeleteByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider domain.OAuthProvider) error {
	query := `DELETE FROM oauth_accounts WHERE user_id = ? AND provider = ?`

	res, err := m.dbClient.ExecContext(ctx, query, userID, provider)
	if err != nil {
		m.logger.Error("failed to delete OAuth account by user ID and provider", err)
		return fmt.Errorf("failed to delete OAuth account by user ID and provider: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrOAuthAccountNotFound
	}

	return nil
}

func NewMysqlOAuthAccountRepository(logger ports.Logger, dbClient *sql.DB) ports.OAuthAccountRepository {
	return &mysqlOAuthAccountRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
