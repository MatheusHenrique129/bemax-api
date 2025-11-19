package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrHealthProfileNotFound = errors.New("health profile not found")
)

type mysqlHealthProfileRepository struct {
	BaseRepository
}

func (m *mysqlHealthProfileRepository) Create(ctx context.Context, profile *domain.HealthProfile) error {
	query := `
		INSERT INTO health_profiles (
			id, user_id, blood_type, height, weight, allergies, 
			medications, medical_conditions, notes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	allergiesJSON, _ := json.Marshal(profile.Allergies)
	medicationsJSON, _ := json.Marshal(profile.Medications)
	conditionsJSON, _ := json.Marshal(profile.MedicalConditions)

	_, err := m.dbClient.ExecContext(ctx, query,
		profile.ID,
		profile.UserID,
		profile.BloodType,
		profile.Height,
		profile.Weight,
		allergiesJSON,
		medicationsJSON,
		conditionsJSON,
		profile.Notes,
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	if err != nil {
		m.logger.Error("failed to create health profile", err)
		return fmt.Errorf("failed to create health profile: %w", err)
	}

	return nil
}

func (m *mysqlHealthProfileRepository) Update(ctx context.Context, profile *domain.HealthProfile) error {
	query := `
		UPDATE health_profiles 
		SET blood_type = ?, height = ?, weight = ?, allergies = ?, 
		    medications = ?, medical_conditions = ?, notes = ?, updated_at = ?
		WHERE user_id = ?
	`

	allergiesJSON, _ := json.Marshal(profile.Allergies)
	medicationsJSON, _ := json.Marshal(profile.Medications)
	conditionsJSON, _ := json.Marshal(profile.MedicalConditions)

	result, err := m.dbClient.ExecContext(ctx, query,
		profile.BloodType,
		profile.Height,
		profile.Weight,
		allergiesJSON,
		medicationsJSON,
		conditionsJSON,
		profile.Notes,
		profile.UpdatedAt,
		profile.UserID,
	)

	if err != nil {
		m.logger.Error("failed to update health profile", err)
		return fmt.Errorf("failed to update health profile: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrHealthProfileNotFound
	}

	return nil
}

func (m *mysqlHealthProfileRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.HealthProfile, error) {
	query := `
		SELECT id, user_id, blood_type, height, weight, allergies, 
		       medications, medical_conditions, notes, created_at, updated_at
		FROM health_profiles
		WHERE user_id = ?
	`

	var profile domain.HealthProfile
	var height, weight sql.NullFloat64
	var allergiesJSON, medicationsJSON, conditionsJSON sql.NullString

	err := m.dbClient.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.BloodType,
		&height,
		&weight,
		&allergiesJSON,
		&medicationsJSON,
		&conditionsJSON,
		&profile.Notes,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrHealthProfileNotFound
	case err != nil:
		m.logger.Error("failed to find health profile", err)
		return nil, fmt.Errorf("failed to find health profile: %w", err)
	}

	if height.Valid {
		profile.Height = &height.Float64
	}
	if weight.Valid {
		profile.Weight = &weight.Float64
	}
	if allergiesJSON.Valid {
		json.Unmarshal([]byte(allergiesJSON.String), &profile.Allergies)
	}
	if medicationsJSON.Valid {
		json.Unmarshal([]byte(medicationsJSON.String), &profile.Medications)
	}
	if conditionsJSON.Valid {
		json.Unmarshal([]byte(conditionsJSON.String), &profile.MedicalConditions)
	}

	return &profile, nil
}

func (m *mysqlHealthProfileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM health_profiles WHERE user_id = ?`

	result, err := m.dbClient.ExecContext(ctx, query, userID)
	if err != nil {
		m.logger.Error("failed to delete health profile", err)
		return fmt.Errorf("failed to delete health profile: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrHealthProfileNotFound
	}

	return nil
}

func NewMysqlHealthProfileRepository(logger ports.Logger, dbClient *sql.DB) ports.HealthProfileRepository {
	return &mysqlHealthProfileRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
