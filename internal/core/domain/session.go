package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	SessionID          string    `json:"session_id"`
	LastAccessTokenJTI string    `json:"last_access_token_jti"`
	DeviceInfo         string    `json:"device_info"`
	IPAddress          string    `json:"ip_address"`
	UserAgent          string    `json:"user_agent"`
	CreatedAt          time.Time `json:"created_at"`
	LastRefreshedAt    time.Time `json:"last_refreshed_at"`
	ExpiresAt          time.Time `json:"expires_at"`
	IsActive           bool      `json:"is_active"`
}

func NewSession(userID uuid.UUID, deviceInfo, ipAddress, userAgent string) *Session {
	now := time.Now().UTC()
	sessionID := uuid.New().String()

	return &Session{
		ID:              uuid.New(),
		UserID:          userID,
		SessionID:       sessionID,
		DeviceInfo:      deviceInfo,
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
		CreatedAt:       now,
		LastRefreshedAt: now,
		ExpiresAt:       now.Add(3 * 24 * time.Hour), // 3 days
		IsActive:        true,
	}
}

func (s *Session) UpdateLastAccessToken(tokenJTI string) {
	s.LastAccessTokenJTI = tokenJTI
	s.LastRefreshedAt = time.Now().UTC()
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Deactivate() {
	s.IsActive = false
}
