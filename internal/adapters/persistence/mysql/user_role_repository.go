package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func NewMysqlUserRoleRepository(logger ports.Logger, dbClient *sql.DB) ports.UserRoleRepository {
	return &mysqlUserRoleRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
