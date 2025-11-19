package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type mysqlReminderCategoryRepository struct {
	BaseRepository
}

func (m *mysqlReminderCategoryRepository) Create(ctx context.Context, category *domain.ReminderCategory) error {
	query := `
		INSERT INTO reminder_categories (
			id, user_id, name, name_key, description, icon, color, 
			scope, display_order, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.dbClient.ExecContext(ctx, query,
		category.ID,
		category.UserID,
		category.Name,
		category.NameKey,
		category.Description,
		category.Icon,
		category.Color,
		category.Scope,
		category.DisplayOrder,
		category.IsActive,
		category.CreatedAt,
		category.UpdatedAt,
	)

	if err != nil {
		m.logger.Error("failed to create reminder category", err)
		return fmt.Errorf("failed to create reminder category: %w", err)
	}

	return nil
}

func (m *mysqlReminderCategoryRepository) Update(ctx context.Context, category *domain.ReminderCategory) error {
	query := `
		UPDATE reminder_categories 
		SET name = ?, description = ?, icon = ?, color = ?, 
		    display_order = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := m.dbClient.ExecContext(ctx, query,
		category.Name,
		category.Description,
		category.Icon,
		category.Color,
		category.DisplayOrder,
		category.IsActive,
		category.UpdatedAt,
		category.ID,
	)

	if err != nil {
		m.logger.Error("failed to update reminder category", err)
		return fmt.Errorf("failed to update reminder category: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (m *mysqlReminderCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reminder_categories WHERE id = ?`

	result, err := m.dbClient.ExecContext(ctx, query, id)
	if err != nil {
		m.logger.Error("failed to delete reminder category", err)
		return fmt.Errorf("failed to delete reminder category: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (m *mysqlReminderCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ReminderCategory, error) {
	query := `
		SELECT id, user_id, name, name_key, description, icon, color, 
		       scope, display_order, is_active, created_at, updated_at
		FROM reminder_categories
		WHERE id = ?
	`

	var category domain.ReminderCategory
	var userID sql.NullString

	err := m.dbClient.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&userID,
		&category.Name,
		&category.NameKey,
		&category.Description,
		&category.Icon,
		&category.Color,
		&category.Scope,
		&category.DisplayOrder,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if userID.Valid {
		uid, _ := uuid.Parse(userID.String)
		category.UserID = &uid
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrCategoryNotFound
	case err != nil:
		m.logger.Error("failed to find category by ID", err)
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	return &category, nil
}

func (m *mysqlReminderCategoryRepository) FindAllActive(ctx context.Context) ([]domain.ReminderCategory, error) {
	query := `
		SELECT id, user_id, name, name_key, description, icon, color, 
		       scope, display_order, is_active, created_at, updated_at
		FROM reminder_categories
		WHERE is_active = true
		ORDER BY display_order, name
	`

	return m.queryCategories(ctx, query)
}

func (m *mysqlReminderCategoryRepository) FindSystemCategories(ctx context.Context) ([]domain.ReminderCategory, error) {
	query := `
		SELECT id, user_id, name, name_key, description, icon, color, 
		       scope, display_order, is_active, created_at, updated_at
		FROM reminder_categories
		WHERE scope = 'system' AND is_active = true
		ORDER BY display_order, name
	`

	return m.queryCategories(ctx, query)
}

func (m *mysqlReminderCategoryRepository) FindUserCategories(ctx context.Context, userID uuid.UUID) ([]domain.ReminderCategory, error) {
	query := `
		SELECT id, user_id, name, name_key, description, icon, color, 
		       scope, display_order, is_active, created_at, updated_at
		FROM reminder_categories
		WHERE user_id = ? AND is_active = true
		ORDER BY display_order, name
	`

	return m.queryCategories(ctx, query, userID)
}

func (m *mysqlReminderCategoryRepository) FindAllForUser(ctx context.Context, userID uuid.UUID) ([]domain.ReminderCategory, error) {
	query := `
		SELECT id, user_id, name, name_key, description, icon, color, 
		       scope, display_order, is_active, created_at, updated_at
		FROM reminder_categories
		WHERE (scope = 'system' OR user_id = ?) AND is_active = true
		ORDER BY display_order, name
	`

	return m.queryCategories(ctx, query, userID)
}

func (m *mysqlReminderCategoryRepository) queryCategories(ctx context.Context, query string, args ...interface{}) ([]domain.ReminderCategory, error) {
	rows, err := m.dbClient.QueryContext(ctx, query, args...)
	if err != nil {
		m.logger.Error("failed to query categories", err)
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.ReminderCategory

	for rows.Next() {
		var category domain.ReminderCategory
		var userID sql.NullString

		err := rows.Scan(
			&category.ID,
			&userID,
			&category.Name,
			&category.NameKey,
			&category.Description,
			&category.Icon,
			&category.Color,
			&category.Scope,
			&category.DisplayOrder,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			m.logger.Error("failed to scan category", err)
			continue
		}

		if userID.Valid {
			uid, _ := uuid.Parse(userID.String)
			category.UserID = &uid
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func NewMysqlReminderCategoryRepository(logger ports.Logger, dbClient *sql.DB) ports.ReminderCategoryRepository {
	return &mysqlReminderCategoryRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
