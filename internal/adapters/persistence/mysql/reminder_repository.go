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
	// ErrReminderNotFound is returned when a reminder is not found
	ErrReminderNotFound = errors.New("reminder not found")
)

type mysqlReminderRepository struct {
	BaseRepository
}

// Create inserts a new reminder into the database
func (m *mysqlReminderRepository) Create(ctx context.Context, reminder *domain.Reminder) error {
	query := `
		INSERT INTO reminders (
			id, user_id, category_id, title, description, status, frequency,
			start_date, end_date, reminder_at, next_occurrence, is_active, metadata,
		    created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Handle NULL fields properly
	var endDate sql.NullTime
	if reminder.EndDate != nil {
		endDate = sql.NullTime{Time: *reminder.EndDate, Valid: true}
	}

	var nextOccurrence sql.NullTime
	if reminder.NextOccurrence != nil {
		nextOccurrence = sql.NullTime{Time: *reminder.NextOccurrence, Valid: true}
	}

	var metadata sql.NullString
	if reminder.Metadata != "" {
		metadata = sql.NullString{String: reminder.Metadata, Valid: true}
	}

	result, err := m.dbClient.ExecContext(ctx, query,
		reminder.ID,
		reminder.UserID,
		reminder.CategoryID,
		reminder.Title,
		reminder.Description,
		reminder.Status,
		reminder.Frequency,
		reminder.StartDate,
		endDate,
		reminder.ReminderAt,
		nextOccurrence,
		reminder.IsActive,
		metadata,
		reminder.CreatedAt,
		reminder.UpdatedAt,
	)

	if err != nil {
		m.logger.Error("failed to create reminder", err)
		return fmt.Errorf("failed to create reminder: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		m.logger.Error("failed to get rows affected", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no rows affected during reminder creation")
	}

	return nil
}

// Update modifies an existing reminder in the database
func (m *mysqlReminderRepository) Update(ctx context.Context, reminder *domain.Reminder) error {
	query := `
		UPDATE reminders 
		SET title = ?, description = ?, status = ?, frequency = ?,
		    start_date = ?, end_date = ?, reminder_at = ?, next_occurrence = ?,
		    is_active = ?, metadata = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	// Handle NULL fields properly
	var endDate sql.NullTime
	if reminder.EndDate != nil {
		endDate = sql.NullTime{Time: *reminder.EndDate, Valid: true}
	}

	var nextOccurrence sql.NullTime
	if reminder.NextOccurrence != nil {
		nextOccurrence = sql.NullTime{Time: *reminder.NextOccurrence, Valid: true}
	}

	var metadata sql.NullString
	if reminder.Metadata != "" {
		metadata = sql.NullString{String: reminder.Metadata, Valid: true}
	}

	result, err := m.dbClient.ExecContext(ctx, query,
		reminder.Title,
		reminder.Description,
		reminder.Status,
		reminder.Frequency,
		reminder.StartDate,
		endDate,
		reminder.ReminderAt,
		nextOccurrence,
		reminder.IsActive,
		metadata,
		reminder.UpdatedAt,
		reminder.ID,
		reminder.UserID, // Security: ensure user owns the reminder
	)

	if err != nil {
		m.logger.Error("failed to update reminder", err)
		return fmt.Errorf("failed to update reminder: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		m.logger.Error("failed to get rows affected", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrReminderNotFound
	}

	return nil
}

// Delete removes a reminder from the database
func (m *mysqlReminderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reminders WHERE id = ?`

	result, err := m.dbClient.ExecContext(ctx, query, id)
	if err != nil {
		m.logger.Error("failed to delete reminder", err)
		return fmt.Errorf("failed to delete reminder: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		m.logger.Error("failed to get rows affected", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrReminderNotFound
	}

	return nil
}

// FindByID retrieves a reminder by its ID with category information
func (m *mysqlReminderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Reminder, error) {
	query := `
		SELECT 
			r.id, r.user_id, r.category_id, r.title, r.description, r.status,
		    r.frequency, r.start_date, r.end_date, r.reminder_at, r.next_occurrence,
		    r.is_active, r.metadata, r.created_at, r.updated_at,
		    c.id, c.user_id, c.name, c.name_key, c.description, c.icon, c.color, c.scope,
		    c.is_active, c.display_order, c.created_at, c.updated_at
		FROM reminders r
		LEFT JOIN reminder_categories c ON r.category_id = c.id
		WHERE r.id = ?
	`

	reminder, err := m.scanReminderWithCategory(m.dbClient.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReminderNotFound
		}
		m.logger.Error("failed to find reminder by ID", err)
		return nil, fmt.Errorf("failed to find reminder by ID: %w", err)
	}

	return reminder, nil
}

// FindByUserID retrieves all reminders for a specific user
func (m *mysqlReminderRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Reminder, error) {
	query := `
		SELECT 
			r.id, r.user_id, r.category_id, r.title, r.description, r.status,
		    r.frequency, r.start_date, r.end_date, r.reminder_at, r.next_occurrence,
		    r.is_active, r.metadata, r.created_at, r.updated_at,
		    c.id, c.user_id, c.name, c.name_key, c.description, c.icon, c.color, c.scope,
		    c.is_active, c.display_order, c.created_at, c.updated_at
		FROM reminders r
		LEFT JOIN reminder_categories c ON r.category_id = c.id
		WHERE r.user_id = ?
		ORDER BY r.reminder_at DESC
	`

	return m.queryReminders(ctx, query, userID)
}

// FindActiveByUserID retrieves all active reminders for a specific user
func (m *mysqlReminderRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Reminder, error) {
	query := `
		SELECT 
			r.id, r.user_id, r.category_id, r.title, r.description, r.status,
		    r.frequency, r.start_date, r.end_date, r.reminder_at, r.next_occurrence,
		    r.is_active, r.metadata, r.created_at, r.updated_at,
		    c.id, c.user_id, c.name, c.name_key, c.description, c.icon, c.color, c.scope,
		    c.is_active, c.display_order, c.created_at, c.updated_at
		FROM reminders r
		LEFT JOIN reminder_categories c ON r.category_id = c.id
		WHERE r.user_id = ? AND r.is_active = true AND r.status = 'active'
		ORDER BY r.reminder_at ASC
	`

	return m.queryReminders(ctx, query, userID)
}

// FindUpcoming retrieves upcoming active reminders for a user, limited by count
func (m *mysqlReminderRepository) FindUpcoming(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Reminder, error) {
	query := `
		SELECT 
			r.id, r.user_id, r.category_id, r.title, r.description, r.status,
		    r.frequency, r.start_date, r.end_date, r.reminder_at, r.next_occurrence,
		    r.is_active, r.metadata, r.created_at, r.updated_at,
		    c.id, c.user_id, c.name, c.name_key, c.description, c.icon, c.color, c.scope,
		    c.is_active, c.display_order, c.created_at, c.updated_at
		FROM reminders r
		LEFT JOIN reminder_categories c ON r.category_id = c.id
		WHERE r.user_id = ? 
		  AND r.is_active = true 
		  AND r.status = 'active'
		  AND r.next_occurrence > NOW()
		ORDER BY r.next_occurrence ASC
		LIMIT ?
	`

	return m.queryReminders(ctx, query, userID, limit)
}

// queryReminders executes a query and returns a slice of reminders
func (m *mysqlReminderRepository) queryReminders(ctx context.Context, query string, args ...interface{}) ([]domain.Reminder, error) {
	rows, err := m.dbClient.QueryContext(ctx, query, args...)
	if err != nil {
		m.logger.Error("failed to query reminders", err)
		return nil, fmt.Errorf("failed to query reminders: %w", err)
	}
	defer rows.Close()

	var reminders []domain.Reminder

	for rows.Next() {
		reminder, err := m.scanReminderWithCategory(rows)
		if err != nil {
			m.logger.Error("failed to scan reminder", err)
			continue // Skip invalid rows but continue processing
		}
		reminders = append(reminders, *reminder)
	}

	if err := rows.Err(); err != nil {
		m.logger.Error("error iterating over reminder rows", err)
		return nil, fmt.Errorf("error iterating over reminder rows: %w", err)
	}

	return reminders, nil
}

// scanReminderWithCategory scans a single reminder with its category from a row
func (m *mysqlReminderRepository) scanReminderWithCategory(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Reminder, error) {
	var reminder domain.Reminder
	var category domain.ReminderCategory

	// Handle NULL fields with sql.Null* types
	var descriptionNull, metadataNull sql.NullString
	var endDateNull, nextOccurrenceNull sql.NullTime
	var categoryUserIDNull, categoryDescriptionNull sql.NullString

	err := scanner.Scan(
		// Reminder fields (15 campos)
		&reminder.ID,
		&reminder.UserID,
		&reminder.CategoryID,
		&reminder.Title,
		&descriptionNull,
		&reminder.Status,
		&reminder.Frequency,
		&reminder.StartDate,
		&endDateNull,
		&reminder.ReminderAt,
		&nextOccurrenceNull,
		&reminder.IsActive,
		&metadataNull,
		&reminder.CreatedAt,
		&reminder.UpdatedAt,
		&category.ID,
		&categoryUserIDNull,
		&category.Name,
		&category.NameKey,
		&categoryDescriptionNull,
		&category.Icon,
		&category.Color,
		&category.Scope,
		&category.IsActive,
		&category.DisplayOrder,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan reminder: %w", err)
	}

	// Map NULL values to Go types - Reminder fields
	if descriptionNull.Valid {
		reminder.Description = descriptionNull.String
	}
	if endDateNull.Valid {
		reminder.EndDate = &endDateNull.Time
	}
	if nextOccurrenceNull.Valid {
		reminder.NextOccurrence = &nextOccurrenceNull.Time
	}
	if metadataNull.Valid {
		reminder.Metadata = metadataNull.String
	}

	// Map NULL values to Go types - Category fields
	if categoryUserIDNull.Valid {
		uid, err := uuid.Parse(categoryUserIDNull.String)
		if err != nil {
			m.logger.Error("failed to parse category user_id", err)
		} else {
			category.UserID = &uid
		}
	}
	if categoryDescriptionNull.Valid {
		category.Description = categoryDescriptionNull.String
	}

	// Attach category to reminder
	reminder.Category = &category

	return &reminder, nil
}

// NewMysqlReminderRepository creates a new instance of MySQL reminder repository
func NewMysqlReminderRepository(logger ports.Logger, dbClient *sql.DB) ports.ReminderRepository {
	return &mysqlReminderRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
