package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
)

type mysqlRoleRepository struct {
	BaseRepository
}

func (m mysqlRoleRepository) FindByName(ctx context.Context, name string) (domain.Role, error) {
	query := `
		SELECT 
		    id,
		    name,
		    description,
		    created_at,
		    updated_at
		FROM roles
		WHERE name = ?
	`

	var res domain.Role
	err := m.dbClient.QueryRowContext(ctx, query, name).Scan(
		&res.ID,
		&res.Name,
		&res.Description,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		m.logger.Error(fmt.Sprintf("error finding role with name %s.", name), err)
		return domain.Role{}, ErrRoleNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error found role with name %s", name), err)
		return domain.Role{}, fmt.Errorf("error finding role by name: %s. %w", name, err)
	}

	return res, nil
}

func NewMysqlRoleRepository(logger ports.Logger, dbClient *sql.DB) ports.RoleRepository {
	return &mysqlRoleRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
