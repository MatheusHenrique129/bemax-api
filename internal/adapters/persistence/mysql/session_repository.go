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

type mysqlSessionRepository struct {
	BaseRepository
}

func (m mysqlSessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO active_sessions (
			id,
			user_id,
			session_id,
			last_access_token_jti,
			device_type,
			user_agent,
			ip_address,
			is_suspicious,
			risk_score,
			created_at, 
		    last_activity_at,
			last_refreshed_at,
		    expires_at,
			is_active
  		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	_, err := m.dbClient.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.SessionID,
		session.LastAccessTokenJTI,
		sql.NullString{String: string(session.DeviceType), Valid: session.DeviceType != ""},
		sql.NullString{String: session.UserAgent, Valid: session.UserAgent != ""},
		sql.NullString{String: session.IPAddress, Valid: session.IPAddress != ""},
		session.IsSuspicious,
		session.RiskScore,
		session.CreatedAt,
		session.LastActivityAt,
		session.LastRefreshedAt,
		session.ExpiresAt,
		session.IsActive,
	)

	if err != nil {
		m.logger.Error("error creating session", err)
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (m mysqlSessionRepository) UpdateSession(ctx context.Context, session *domain.Session) error {
	query := `
        UPDATE active_sessions 
        SET last_access_token_jti = ?,
            last_activity_at = ?,
            last_refreshed_at = ?,
            is_active = ?,
         	risk_score = ?,
		    is_suspicious = ?
        WHERE session_id = ?
`

	result, err := m.dbClient.ExecContext(ctx, query,
		session.LastAccessTokenJTI,
		session.LastActivityAt,
		session.LastRefreshedAt,
		session.IsActive,
		session.RiskScore,
		session.IsSuspicious,
		session.SessionID,
	)

	if err != nil {
		m.logger.Error("error updating session", err)
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("session not found: %s", session.SessionID)
	}

	return nil
}

func (m mysqlSessionRepository) FindBySessionID(ctx context.Context, sessionID string) (*domain.Session, error) {
	query := `
        SELECT 
            id,
            user_id,
            session_id,
            last_access_token_jti,
            device_type,
            user_agent,
            ip_address,
            is_suspicious,
            risk_score,
            created_at,
            last_activity_at,
            last_refreshed_at,
            expires_at,
            is_active
        FROM active_sessions
        WHERE session_id = ? AND is_active = true AND expires_at > NOW()
`

	var session domain.Session
	var deviceType, userAgent, ipAddress sql.NullString

	err := m.dbClient.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.SessionID,
		&session.LastAccessTokenJTI,
		&deviceType,
		&userAgent,
		&ipAddress,
		&session.IsSuspicious,
		&session.RiskScore,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.LastRefreshedAt,
		&session.ExpiresAt,
		&session.IsActive,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("session not found: %s", sessionID)
	case err != nil:
		m.logger.Error(fmt.Sprintf("error scanning session for ID %s", sessionID), err)
		return nil, fmt.Errorf("error finding session with ID %s: %w", sessionID, err)
	}

	if deviceType.Valid {
		session.DeviceType = domain.DeviceType(deviceType.String)
	}
	if userAgent.Valid {
		session.UserAgent = userAgent.String
	}
	if ipAddress.Valid {
		session.IPAddress = ipAddress.String
	}

	return &session, nil
}

func (m mysqlSessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	query := `
        SELECT 
            id,
            user_id,
            session_id,
            last_access_token_jti,
            device_type,
            user_agent,
            ip_address,
            is_suspicious,
            risk_score,
            created_at,
            last_activity_at,
            last_refreshed_at,
            expires_at,
            is_active
        FROM active_sessions
        WHERE id = ? AND is_active = true AND expires_at > NOW()
`

	var session domain.Session
	var deviceType, userAgent, ipAddress sql.NullString

	err := m.dbClient.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.SessionID,
		&session.LastAccessTokenJTI,
		&deviceType,
		&userAgent,
		&ipAddress,
		&session.IsSuspicious,
		&session.RiskScore,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.LastRefreshedAt,
		&session.ExpiresAt,
		&session.IsActive,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrSessionNotFound
	case err != nil:
		m.logger.Error(fmt.Sprintf("error scanning session for ID %s", id), err)
		return nil, fmt.Errorf("error finding session with ID %s: %w", id, err)
	}

	if deviceType.Valid {
		session.DeviceType = domain.DeviceType(deviceType.String)
	}
	if userAgent.Valid {
		session.UserAgent = userAgent.String
	}
	if ipAddress.Valid {
		session.IPAddress = ipAddress.String
	}

	return &session, nil
}

func (m mysqlSessionRepository) FindActiveUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	query := `
        SELECT 
        	id,
            user_id,
            session_id,
            last_access_token_jti,
            device_type,
            user_agent,
            ip_address,
            is_suspicious,
            risk_score,
            created_at,
            last_activity_at,
            last_refreshed_at,
            expires_at,
            is_active
        FROM active_sessions
        WHERE user_id = ? AND is_active = true AND expires_at > NOW()
        ORDER BY last_activity_at DESC
`

	rows, err := m.dbClient.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding user sessions: %w", err)
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var sessions []domain.Session
	for rows.Next() {
		var session domain.Session
		var deviceType, userAgent, ipAddress sql.NullString

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.SessionID,
			&session.LastAccessTokenJTI,
			&deviceType,
			&userAgent,
			&ipAddress,
			&session.IsSuspicious,
			&session.RiskScore,
			&session.CreatedAt,
			&session.LastActivityAt,
			&session.LastRefreshedAt,
			&session.ExpiresAt,
			&session.IsActive,
		)
		if err != nil {
			m.logger.Error("error scanning session row", err)
			continue
		}

		if deviceType.Valid {
			session.DeviceType = domain.DeviceType(deviceType.String)
		}
		if userAgent.Valid {
			session.UserAgent = userAgent.String
		}
		if ipAddress.Valid {
			session.IPAddress = ipAddress.String
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (m mysqlSessionRepository) IsLatestAccessToken(ctx context.Context, sessionID, tokenJTI string) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM active_sessions 
        WHERE session_id = ? 
          AND last_access_token_jti = ? 
          AND is_active = true 
          AND expires_at > NOW()
`

	var count int
	err := m.dbClient.QueryRowContext(ctx, query, sessionID, tokenJTI).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking latest access token: %w", err)
	}

	return count > 0, nil
}

func (m mysqlSessionRepository) UpdateLastAccessToken(ctx context.Context, sessionID, tokenJTI string) error {
	query := `
        UPDATE active_sessions 
        SET last_access_token_jti = ?,
            last_activity_at = NOW(),
            last_refreshed_at = NOW()
        WHERE session_id = ? AND is_active = true
`

	result, err := m.dbClient.ExecContext(ctx, query, tokenJTI, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update last access token: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("session not found or inactive: %s", sessionID)
	}

	return nil
}

func (m mysqlSessionRepository) DeactivateSession(ctx context.Context, sessionID string) error {
	query := `UPDATE active_sessions SET is_active = false WHERE session_id = ?`

	_, err := m.dbClient.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to deactivate session: %w", err)
	}

	return nil
}

func (m mysqlSessionRepository) DeactivateAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE active_sessions SET is_active = false WHERE user_id = ?`

	result, err := m.dbClient.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user sessions: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	m.logger.Info(fmt.Sprintf("deactivated %d sessions for user %s", rowsAffected, userID))

	return nil
}

func (m mysqlSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM active_sessions WHERE expires_at <= NOW() OR is_active = false`

	result, err := m.dbClient.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		m.logger.Info(fmt.Sprintf("deleted %d expired/inactive sessions", rowsAffected))
	}

	return nil
}

func (m mysqlSessionRepository) UpdateSessionRiskScore(ctx context.Context, sessionID string, riskScore int64, isSuspicious bool) error {
	query := `
		UPDATE active_sessions 
		SET risk_score = ?, is_suspicious = ?
		WHERE session_id = ?
	`

	result, err := m.dbClient.ExecContext(ctx, query, riskScore, isSuspicious, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session risk score: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

func NewMysqlSessionRepository(logger ports.Logger, dbClient *sql.DB) ports.SessionRepository {
	return &mysqlSessionRepository{
		BaseRepository: NewBaseRepository(dbClient, logger),
	}
}
