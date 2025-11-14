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

type mysqlUserRepository struct {
	BaseRepository
}

func (m mysqlUserRepository) Create(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (
			id,
			email,
		    password_hash,
		    auth_provider,
		    full_name,
		    cpf,
		    birth_date,
		    phone,
		    profile_picture,
		    email_verified,
		    phone_verified,
		    profile_completed,
			status,
			token_version,
		    created_at,
		    updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	row, err := m.dbClient.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.AuthProvider,
		user.FullName,
		user.CPF,
		user.BirthDate,
		user.Phone,
		user.ProfilePicture,
		user.EmailVerified,
		user.PhoneVerified,
		user.ProfileCompleted,
		user.Status,
		user.TokenVersion,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	if err != nil {
		m.logger.Error("error inserting user from database", err)
		return fmt.Errorf("%s. %w", "failed to create user", err)
	}

	rows, err := row.RowsAffected()
	if err != nil || rows < 1 {
		m.logger.Error("no rows affected during create", err)
		return fmt.Errorf("%v. %w", ErrNoRowsAffected, err)
	}

	return nil
}

func (m mysqlUserRepository) Update(ctx context.Context, user domain.User) error {
	query := `
		UPDATE users 
		SET full_name = ?,
		    cpf = ?,
		    phone = ?,
		    birth_date = ?,
		    profile_picture = ?,
		    profile_completed = ?,
		    updated_at = ?
		WHERE id = ?
`

	now := time.Now().UTC()
	result, err := m.dbClient.ExecContext(ctx, query,
		user.FullName,
		sql.NullString{String: user.CPF, Valid: user.CPF != ""},
		sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		sql.NullTime{Time: *user.BirthDate, Valid: !user.BirthDate.IsZero()},
		sql.NullString{String: user.ProfilePicture, Valid: user.ProfilePicture != ""},
		user.ProfileCompleted,
		now,
		user.ID,
	)

	if err != nil {
		m.logger.Error("error updating user profile", err)
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows < 1 {
		m.logger.Error("no rows affected during profile update", err)
		return ErrNoRowsAffected
	}

	return nil
}

func (m mysqlUserRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	query := `
		SELECT 
		    id,
		    email,
		    password_hash,
		    auth_provider,
		    full_name,
		    cpf,
		    birth_date,
		    phone,
		    profile_picture,
		    email_verified,
		    phone_verified,
		    profile_completed,
		    last_login,
		    status,
		    token_version,
		    created_at,
		    updated_at
		FROM users
		WHERE id = ?
	`

	var res domain.User
	var cpfNull, phoneNull, profilePictureNull sql.NullString
	var birthDateNull, lastLoginNull sql.NullTime

	err := m.dbClient.QueryRowContext(ctx, query, id).Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.AuthProvider,
		&res.FullName,
		&cpfNull,
		&birthDateNull,
		&phoneNull,
		&profilePictureNull,
		&res.EmailVerified,
		&res.PhoneVerified,
		&res.ProfileCompleted,
		&lastLoginNull,
		&res.Status,
		&res.TokenVersion,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if cpfNull.Valid {
		res.CPF = cpfNull.String
	}
	if phoneNull.Valid {
		res.Phone = phoneNull.String
	}
	if profilePictureNull.Valid {
		res.ProfilePicture = profilePictureNull.String
	}
	if birthDateNull.Valid {
		res.BirthDate = &birthDateNull.Time
	}
	if lastLoginNull.Valid {
		res.LastLogin = &lastLoginNull.Time
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("error getting user with id %s from users", id), err)
		return domain.User{}, ErrUserNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error found user with id %s", id), err)
		return domain.User{}, fmt.Errorf("error finding user by ID: %s. %v", id, err)
	}

	return res, nil
}

func (m mysqlUserRepository) FindByCPF(ctx context.Context, cpf string) (domain.User, error) {
	query := `
		SELECT 
			id,
		    email,
		    password_hash,
		    auth_provider,
		    full_name,
		    cpf,
		    birth_date,
		    phone,
		    profile_picture,
		    email_verified,
		    phone_verified,
		    profile_completed,
		    last_login,
		    status,
		    token_version,
		    created_at,
		    updated_at
		FROM users
		WHERE cpf = ?
	`

	var res domain.User
	var cpfNull, phoneNull, profilePictureNull sql.NullString
	var birthDateNull, lastLoginNull sql.NullTime

	err := m.dbClient.QueryRowContext(ctx, query, cpf).Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.AuthProvider,
		&res.FullName,
		&cpfNull,
		&birthDateNull,
		&phoneNull,
		&profilePictureNull,
		&res.EmailVerified,
		&res.PhoneVerified,
		&res.ProfileCompleted,
		&lastLoginNull,
		&res.Status,
		&res.TokenVersion,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if cpfNull.Valid {
		res.CPF = cpfNull.String
	}
	if phoneNull.Valid {
		res.Phone = phoneNull.String
	}
	if profilePictureNull.Valid {
		res.ProfilePicture = profilePictureNull.String
	}
	if birthDateNull.Valid {
		res.BirthDate = &birthDateNull.Time
	}
	if lastLoginNull.Valid {
		res.LastLogin = &lastLoginNull.Time
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("error finding user with cpf %s.", cpf), err)
		return domain.User{}, ErrUserNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error found user with cpf %s", cpf), err)
		return domain.User{}, fmt.Errorf("error finding user by cpf: %s. %v", cpf, err)
	}

	return res, nil
}

func (m mysqlUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT 
		    id,
		    email,
		    password_hash,
		    auth_provider,
		    full_name,
		    cpf,
		    birth_date,
		    phone,
		    profile_picture,
		    email_verified,
		    phone_verified,
		    profile_completed,
		    last_login,
		    status,
		    token_version,
		    created_at,
		    updated_at
		FROM users
		WHERE email = ?
`

	var res domain.User
	var cpfNull, phoneNull, profilePictureNull sql.NullString
	var birthDateNull, lastLoginNull sql.NullTime

	err := m.dbClient.QueryRowContext(ctx, query, email).Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.AuthProvider,
		&res.FullName,
		&cpfNull,
		&birthDateNull,
		&phoneNull,
		&profilePictureNull,
		&res.EmailVerified,
		&res.PhoneVerified,
		&res.ProfileCompleted,
		&lastLoginNull,
		&res.Status,
		&res.TokenVersion,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if cpfNull.Valid {
		res.CPF = cpfNull.String
	}
	if phoneNull.Valid {
		res.Phone = phoneNull.String
	}
	if profilePictureNull.Valid {
		res.ProfilePicture = profilePictureNull.String
	}
	if birthDateNull.Valid {
		res.BirthDate = &birthDateNull.Time
	}
	if lastLoginNull.Valid {
		res.LastLogin = &lastLoginNull.Time
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("error finding user with email %s.", email), err)
		return domain.User{}, ErrUserNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error found user with email %s", email), err)
		return domain.User{}, fmt.Errorf("error finding user by email: %s. %v", email, err)
	}

	return res, nil
}

func (m mysqlUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users SET
			last_login = ?
		WHERE id = ?
`

	result, err := m.dbClient.ExecContext(ctx, query, time.Now().UTC(), id.String())
	if err != nil {
		m.logger.Error("error updating last login for user", err)
		return ErrQuery
	}

	rows, err := result.RowsAffected()
	if err != nil || rows < 1 {
		m.logger.Error("no rows affected during last login update", err)
		return ErrNoRowsAffected
	}

	return nil
}

func (m mysqlUserRepository) GetLoginAttempts(ctx context.Context, email string, minutes int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM login_attempts
		WHERE email = ? 
		  AND success = false 
		  AND created_at > DATE_SUB(NOW(), INTERVAL ? MINUTE)
`

	var count int
	err := m.dbClient.QueryRowContext(ctx, query, email, minutes).Scan(&count)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("error finding login attempts with email %s.", email), err)
		return 0, ErrLoginNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error found login attempts with email %s", email), err)
		return 0, fmt.Errorf("error finding login attempts by email: %s. %v", email, err)
	}

	return count, nil
}

func (m mysqlUserRepository) RecordLoginAttempt(ctx context.Context, email string, success bool, ipAddress, userAgent string) error {
	query := `
		INSERT INTO login_attempts (
			id,
			email,
		    success,
		    ip_address,
		    user_agent,
		    created_at
		) VALUES (?, ?, ?, ?, ?, ?)
`
	row, err := m.dbClient.ExecContext(ctx, query,
		uuid.New().String(),
		email,
		success,
		ipAddress,
		userAgent,
		time.Now().UTC(),
	)

	if err != nil {
		m.logger.Error("error inserting record login attempt in database", err)
		return fmt.Errorf("%s. %w", "error to insert record login attempt", err)
	}

	rows, err := row.RowsAffected()
	if err != nil || rows < 1 {
		m.logger.Error("no rows affected during record login attempt", err)
		return fmt.Errorf("%s. %w", "no rows affected during record login attempt", err)
	}

	return nil
}

func (m mysqlUserRepository) GetTokenVersion(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT token_version FROM users 
        WHERE id = ?
`

	var version int
	err := m.dbClient.QueryRowContext(ctx, query, userID).Scan(&version)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("user not found for id: %s", userID), err)
		return 0, ErrUserNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error getting token version for user: %s", userID), err)
		return 0, fmt.Errorf("error getting token version: %w", err)
	}

	return version, nil
}

func (m mysqlUserRepository) IncrementTokenVersion(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users 
		SET token_version = token_version + 1
		WHERE id = ?
`

	result, err := m.dbClient.ExecContext(ctx, query, userID)
	if err != nil {
		m.logger.Error(fmt.Sprintf("error incrementing token version for user: %s", userID), err)
		return fmt.Errorf("error incrementing token version: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		m.logger.Error(fmt.Sprintf("no user found to increment token version, userID: %s", userID), nil)
		return ErrUserNotFound
	}

	m.logger.Info(fmt.Sprintf("token version incremented for user: %s", userID))
	return nil
}

func NewMysqlUserRepository(logger ports.Logger, dbClient *sql.DB) ports.UserRepository {
	return &mysqlUserRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
