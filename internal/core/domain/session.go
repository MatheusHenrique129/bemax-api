package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type DeviceType string

const (
	DeviceTypeMobile  DeviceType = "mobile"
	DeviceTypeDesktop DeviceType = "desktop"
	DeviceTypeTablet  DeviceType = "tablet"
	DeviceTypeUnknown DeviceType = "unknown"
)

type Session struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"user_id"`
	SessionID          string     `json:"session_id"`
	LastAccessTokenJTI string     `json:"last_access_token_jti"`
	DeviceType         DeviceType `json:"device_type"`
	UserAgent          string     `json:"user_agent"`
	IPAddress          string     `json:"ip_address"`
	IsSuspicious       bool       `json:"is_suspicious"`
	RiskScore          int        `json:"risk_score"`
	CreatedAt          time.Time  `json:"created_at"`
	LastActivityAt     time.Time  `json:"last_activity_at"`
	LastRefreshedAt    time.Time  `json:"last_refreshed_at"`
	ExpiresAt          time.Time  `json:"expires_at"`
	IsActive           bool       `json:"is_active"`
}

func NewSession(userID uuid.UUID, deviceInfo, ipAddress, userAgent string) *Session {
	now := time.Now().UTC()
	sessionID := uuid.New().String()

	return &Session{
		ID:              uuid.New(),
		UserID:          userID,
		SessionID:       sessionID,
		DeviceType:      DetectDeviceType(userAgent),
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
		CreatedAt:       now,
		LastActivityAt:  now,
		LastRefreshedAt: now,
		IsSuspicious:    false,
		RiskScore:       0,
		ExpiresAt:       now.Add(3 * 24 * time.Hour), // 3 days TODO create method ttl in env
		IsActive:        true,
	}
}

func (s *Session) UpdateLastAccessToken(tokenJTI string) {
	s.LastAccessTokenJTI = tokenJTI
	s.LastActivityAt = time.Now().UTC()
	s.LastRefreshedAt = time.Now().UTC()
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Deactivate() {
	s.IsActive = false
}

func (s *Session) MarkSuspicious() {
	s.IsSuspicious = true
	s.RiskScore = 100
}

func DetectDeviceType(userAgent string) DeviceType {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone"):
		return DeviceTypeMobile
	case strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad"):
		return DeviceTypeTablet
	case strings.Contains(ua, "windows") || strings.Contains(ua, "macintosh") || strings.Contains(ua, "linux"):
		return DeviceTypeDesktop
	default:
		return DeviceTypeUnknown
	}
}
