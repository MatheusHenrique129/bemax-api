package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
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
		    full_name,
		    cpf,
		    birth_date,
		    phone,
			status,
		    created_at,
		    updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	row, err := m.dbClient.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FullName,
		user.CPF,
		user.BirthDate,
		user.Phone,
		user.Status,
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

func (m mysqlUserRepository) FindByCPF(ctx context.Context, cpf string) (domain.User, error) {
	query := `
		SELECT 
			id,
		    email,
		    password_hash,
		    full_name,
		    cpf,
		    birth_date,
		    phone,
		    status,
		    created_at,
		    updated_at
		FROM users
		WHERE cpf = ?
	`

	var res domain.User
	err := m.dbClient.QueryRowContext(ctx, query, cpf).Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.FullName,
		&res.CPF,
		&res.BirthDate,
		&res.Phone,
		&res.Status,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

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

func NewMysqlUserRepository(logger ports.Logger, dbClient *sql.DB) ports.UserRepository {
	return &mysqlUserRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
