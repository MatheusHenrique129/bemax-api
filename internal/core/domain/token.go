package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

type Token struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	SessionID     uuid.UUID  `json:"session_id"`
	Token         string     `json:"token"`
	Type          TokenType  `json:"token_type"`
	IsRevoked     bool       `json:"is_revoked"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	RevokedReason string     `json:"revoked_reason,omitempty"`
	ExpiresAt     time.Time  `json:"expires_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Claims struct {
	jwt.RegisteredClaims
	Email        string    `json:"email"`
	TokenType    TokenType `json:"token_type"`
	Roles        []Role    `json:"roles"`
	UserID       uuid.UUID `json:"user_id"`
	TokenVersion int       `json:"token_version"`
	SessionID    string    `json:"session_id"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *Token) Revoke(reason string) {
	now := time.Now().UTC()
	t.IsRevoked = true
	t.RevokedAt = &now
	t.RevokedReason = reason
}

func (t *Token) IsValid() bool {
	return !t.IsRevoked && !t.IsExpired()
}

func NewTokenUserClaims(
	userID uuid.UUID,
	email string,
	tokenType TokenType,
	roles []Role,
	ttl time.Duration,
	tokenVersion int,
	sessionID string,
) *Claims {
	tokenID := uuid.New()
	now := time.Now().UTC()

	return &Claims{
		UserID:       userID,
		Email:        email,
		Roles:        roles,
		TokenType:    tokenType,
		TokenVersion: tokenVersion,
		SessionID:    sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
}

func NewToken(userID uuid.UUID, sessionID uuid.UUID, token string, tokenType TokenType, ttl time.Duration) *Token {
	now := time.Now().UTC()
	return &Token{
		ID:            uuid.New(),
		UserID:        userID,
		SessionID:     sessionID,
		Token:         token,
		Type:          tokenType,
		IsRevoked:     false,
		RevokedAt:     nil,
		RevokedReason: "",
		ExpiresAt:     now.Add(ttl),
		CreatedAt:     now,
	}
}
