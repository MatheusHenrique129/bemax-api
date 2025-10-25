package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

type Token struct {
	ExpiresAt time.Time
	CreatedAt time.Time
	Token     string
	Type      TokenType
	ID        uuid.UUID
	UserID    uuid.UUID
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

func NewTokenUserClaims(
	userID uuid.UUID,
	email string,
	tokenType TokenType,
	roles []Role, ttl time.Duration,
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

func NewToken(userID uuid.UUID, token string, tokenType TokenType, ttl time.Duration) *Token {
	now := time.Now().UTC()
	return &Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		Type:      tokenType,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}
}
