package domain

import (
	"fmt"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/pkg/hash"
	"github.com/google/uuid"
	"github.com/klassmann/cpfcnpj"
)

const (
	MinNameLen     = 3
	MinPasswordLen = 6

	UserStatusActive              Status = "active"
	UserStatusInactive            Status = "inactive"
	UserStatusBlocked             Status = "blocked"
	UserStatusPendingVerification Status = "pending_verification"

	AuthProviderLocal AuthProvider = "local"
	AuthProviderOAuth AuthProvider = "oauth"
)

type Status string
type AuthProvider string

type User struct {
	BirthDate        *time.Time   `json:"birth_date"`
	UpdatedAt        time.Time    `json:"updated_at"`
	CreatedAt        time.Time    `json:"created_at"`
	LastLogin        *time.Time   `json:"last_login"`
	Password         string       `json:"-"`
	Phone            string       `json:"phone"`
	Status           Status       `json:"status"`
	AuthProvider     AuthProvider `json:"auth_provider"`
	CPF              string       `json:"cpf"`
	Email            string       `json:"email"`
	FullName         string       `json:"full_name"`
	ProfilePicture   string       `json:"profile_picture,omitempty"`
	Addresses        []Address    `json:"addresses,omitempty"`
	Roles            []Role       `json:"roles"`
	EmailVerified    bool         `json:"email_verified"`
	PhoneVerified    bool         `json:"phone_verified"`
	ProfileCompleted bool         `json:"profile_completed"`
	ID               uuid.UUID    `json:"id"`
	TokenVersion     int          `json:"-"`
}

func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) Block() {
	u.Status = UserStatusBlocked
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) IsLocalAuth() bool {
	return u.AuthProvider == AuthProviderLocal
}

func (u *User) IsOAuthAuth() bool {
	return u.AuthProvider == AuthProviderOAuth
}

func (u *User) UpdateProfilePicture(pictureURL string) {
	u.ProfilePicture = pictureURL
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) CheckPassword(password string) error {
	//if u.IsOAuthAuth() {
	//	return fmt.Errorf("OAuth users don't have passwords")
	//}
	return hash.CheckPassword(password, u.Password)
}

func (u *User) UpdateLastLogin() {
	now := time.Now().UTC()
	u.LastLogin = &now
	u.UpdatedAt = now
}

func (u *User) MaskCPF() {
	plainCPF := cpfcnpj.Clean(u.CPF)

	maskedCPF := fmt.Sprintf("***.***.%s-%s",
		plainCPF[6:9],
		plainCPF[9:11],
	)

	u.CPF = maskedCPF
}

func (u *User) Update(email, fullName, phone, profilePicture string) {
	if email != "" {
		u.Email = email
	}
	if fullName != "" {
		u.FullName = fullName
	}
	if phone != "" {
		u.Phone = phone
	}
	if profilePicture != "" {
		u.ProfilePicture = profilePicture
	}
	u.UpdatedAt = time.Now().UTC()
}

// HasCPF checks if user has CPF registered
func (u *User) HasCPF() bool {
	return u.CPF != ""
}

// RequiresCPFCompletion checks if user needs to complete profile with CPF
func (u *User) RequiresCPFCompletion() bool {
	return !u.ProfileCompleted || !u.HasCPF()
}

// CompleteMandatoryProfile completes user profile with mandatory fields
func (u *User) CompleteMandatoryProfile(cpf, phone string, birthDate time.Time) error {
	if !cpfcnpj.ValidateCPF(cpf) {
		return fmt.Errorf("invalid CPF")
	}

	u.CPF = cpfcnpj.Clean(cpf)
	u.Phone = phone
	u.BirthDate = &birthDate
	u.ProfileCompleted = true
	u.UpdatedAt = time.Now().UTC()

	return nil
}

// VerifyEmail marks email as verified
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now().UTC()
}

// VerifyPhone marks phone as verified
func (u *User) VerifyPhone() {
	u.PhoneVerified = true
	u.UpdatedAt = time.Now().UTC()
}

// NewOAuthUser creates a new OAuth user (without password, CPF optional)
func NewOAuthUser(email, fullName string, emailVerified bool) (User, apierrors.RestError) {
	now := time.Now().UTC()

	// TODO implement Validators in email
	return User{
		ID:               uuid.New(),
		Email:            email,
		AuthProvider:     AuthProviderOAuth,
		FullName:         fullName,
		BirthDate:        nil,
		EmailVerified:    emailVerified,
		PhoneVerified:    false,
		ProfileCompleted: false,
		Status:           UserStatusActive,
		TokenVersion:     0,
		CreatedAt:        now,
		UpdatedAt:        now,
		Roles:            []Role{},
		Addresses:        []Address{},
	}, nil
}
