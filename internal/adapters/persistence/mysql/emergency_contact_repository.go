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
	ErrEmergencyContactNotFound = errors.New("emergency contact not found")
)

type mysqlEmergencyContactRepository struct {
	BaseRepository
}

func (m *mysqlEmergencyContactRepository) Create(ctx context.Context, contact *domain.EmergencyContact) error {
	query := `
		INSERT INTO emergency_contacts (
			id, user_id, name, relationship, phone, email, address, 
			notes, is_primary, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.dbClient.ExecContext(ctx, query,
		contact.ID,
		contact.UserID,
		contact.Name,
		contact.Relationship,
		contact.Phone,
		contact.Email,
		contact.Address,
		contact.Notes,
		contact.IsPrimary,
		contact.IsActive,
		contact.CreatedAt,
		contact.UpdatedAt,
	)

	if err != nil {
		m.logger.Error("failed to create emergency contact", err)
		return fmt.Errorf("failed to create emergency contact: %w", err)
	}

	return nil
}

func (m *mysqlEmergencyContactRepository) Update(ctx context.Context, contact *domain.EmergencyContact) error {
	query := `
		UPDATE emergency_contacts 
		SET name = ?, relationship = ?, phone = ?, email = ?, 
		    address = ?, notes = ?, is_primary = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := m.dbClient.ExecContext(ctx, query,
		contact.Name,
		contact.Relationship,
		contact.Phone,
		contact.Email,
		contact.Address,
		contact.Notes,
		contact.IsPrimary,
		contact.IsActive,
		contact.UpdatedAt,
		contact.ID,
	)

	if err != nil {
		m.logger.Error("failed to update emergency contact", err)
		return fmt.Errorf("failed to update emergency contact: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrEmergencyContactNotFound
	}

	return nil
}

func (m *mysqlEmergencyContactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM emergency_contacts WHERE id = ?`

	result, err := m.dbClient.ExecContext(ctx, query, id)
	if err != nil {
		m.logger.Error("failed to delete emergency contact", err)
		return fmt.Errorf("failed to delete emergency contact: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrEmergencyContactNotFound
	}

	return nil
}

func (m *mysqlEmergencyContactRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.EmergencyContact, error) {
	query := `
		SELECT id, user_id, name, relationship, phone, email, address, 
		       notes, is_primary, is_active, created_at, updated_at
		FROM emergency_contacts
		WHERE id = ?
	`

	contact, err := m.scanContact(m.dbClient.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmergencyContactNotFound
		}
		return nil, err
	}

	return contact, nil
}

func (m *mysqlEmergencyContactRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.EmergencyContact, error) {
	query := `
		SELECT id, user_id, name, relationship, phone, email, address, 
		       notes, is_primary, is_active, created_at, updated_at
		FROM emergency_contacts
		WHERE user_id = ?
		ORDER BY is_primary DESC, name
	`

	return m.queryContacts(ctx, query, userID)
}

func (m *mysqlEmergencyContactRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]domain.EmergencyContact, error) {
	query := `
		SELECT id, user_id, name, relationship, phone, email, address, 
		       notes, is_primary, is_active, created_at, updated_at
		FROM emergency_contacts
		WHERE user_id = ? AND is_active = true
		ORDER BY is_primary DESC, name
	`

	return m.queryContacts(ctx, query, userID)
}

func (m *mysqlEmergencyContactRepository) FindPrimaryByUserID(ctx context.Context, userID uuid.UUID) (*domain.EmergencyContact, error) {
	query := `
		SELECT id, user_id, name, relationship, phone, email, address, 
		       notes, is_primary, is_active, created_at, updated_at
		FROM emergency_contacts
		WHERE user_id = ? AND is_primary = true AND is_active = true
		LIMIT 1
	`

	contact, err := m.scanContact(m.dbClient.QueryRowContext(ctx, query, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmergencyContactNotFound
		}
		return nil, err
	}

	return contact, nil
}

func (m *mysqlEmergencyContactRepository) UnsetAllPrimaryForUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE emergency_contacts SET is_primary = false WHERE user_id = ?`

	_, err := m.dbClient.ExecContext(ctx, query, userID)
	if err != nil {
		m.logger.Error("failed to unset primary contacts", err)
		return fmt.Errorf("failed to unset primary contacts: %w", err)
	}

	return nil
}

func (m *mysqlEmergencyContactRepository) queryContacts(ctx context.Context, query string, args ...interface{}) ([]domain.EmergencyContact, error) {
	rows, err := m.dbClient.QueryContext(ctx, query, args...)
	if err != nil {
		m.logger.Error("failed to query emergency contacts", err)
		return nil, fmt.Errorf("failed to query emergency contacts: %w", err)
	}
	defer rows.Close()

	var contacts []domain.EmergencyContact

	for rows.Next() {
		contact, err := m.scanContact(rows)
		if err != nil {
			m.logger.Error("failed to scan emergency contact", err)
			continue
		}
		contacts = append(contacts, *contact)
	}

	return contacts, nil
}

func (m *mysqlEmergencyContactRepository) scanContact(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.EmergencyContact, error) {
	var contact domain.EmergencyContact

	err := scanner.Scan(
		&contact.ID,
		&contact.UserID,
		&contact.Name,
		&contact.Relationship,
		&contact.Phone,
		&contact.Email,
		&contact.Address,
		&contact.Notes,
		&contact.IsPrimary,
		&contact.IsActive,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func NewMysqlEmergencyContactRepository(logger ports.Logger, dbClient *sql.DB) ports.EmergencyContactRepository {
	return &mysqlEmergencyContactRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
