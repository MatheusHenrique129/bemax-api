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

type mysqlUserRoleRepository struct {
	BaseRepository
}

// AssignRoles insert user-role relations.
func (m mysqlUserRoleRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (
			user_id,
			role_id, 
			assigned_at
		) VALUES (?, ?, ?)
`

	row, err := m.dbClient.ExecContext(ctx, query, userID, roleID, time.Now().UTC())
	if err != nil {
		m.logger.Error("error inserting user role relation in database", err)
		return fmt.Errorf("%s. %w", "failed to assign user role", err)
	}

	rows, err := row.RowsAffected()
	if err != nil || rows < 1 {
		m.logger.Error("no rows affected during assign user role", err)
		return fmt.Errorf("%v. %w", ErrNoRowsAffected, err)
	}

	return nil
}

func (m mysqlUserRoleRepository) FindRolesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Role, error) {
	res := make([]domain.Role, 0)

	query := `
		SELECT 
			r.id,
		    r.name,
		    r.description,
		    r.created_at,
		    r.updated_at
		FROM user_roles ur
		INNER JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ?
`

	rows, err := m.dbClient.QueryContext(ctx, query, userID)
	if err != nil {
		m.logger.Error(fmt.Sprintf("error querying roles for userID %s", userID), err)
		return nil, fmt.Errorf("%s. %w", "failed to find roles for user", err)
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		r := domain.Role{}

		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.CreatedAt,
			&r.UpdatedAt,
		)

		switch {
		case errors.Is(err, sql.ErrNoRows):
			m.logger.Error(fmt.Sprintf("no roles found for userID %s", userID), err)
			return nil, ErrRolesForUserNotFound
		case err != nil:
			m.logger.Error(fmt.Sprintf("error scanning role for userID %s", userID), err)
			return nil, fmt.Errorf("error scanning role for user with ID %s. %w", userID, err)
		}

		res = append(res, r)
	}

	return res, nil
}

func NewMysqlUserRoleRepository(logger ports.Logger, dbClient *sql.DB) ports.UserRoleRepository {
	return &mysqlUserRoleRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
