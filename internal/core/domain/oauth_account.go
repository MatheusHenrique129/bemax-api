package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type OAuthProvider string

const (
	OAuthProviderGoogle   OAuthProvider = "google"
	OAuthProviderFacebook OAuthProvider = "facebook"
	OAuthProviderApple    OAuthProvider = "apple"
	OAuthProviderGithub   OAuthProvider = "github"
	OAuthProviderTwitter  OAuthProvider = "twitter"
)

type OAuthAccount struct {
	ID              uuid.UUID     `json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	Provider        OAuthProvider `json:"provider"`
	ProviderUID     string        `json:"provider_uid"`
	FirebaseUID     string        `json:"firebase_uid"`
	ProviderEmail   string        `json:"provider_email"`
	ProviderName    string        `json:"provider_name"`
	ProviderPicture string        `json:"provider_picture"`
	EmailVerified   bool          `json:"email_verified"`
	ExpiresAt       *time.Time    `json:"expires_at"`
	LastLoginAt     *time.Time    `json:"last_login_at,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

func (o *OAuthAccount) UpdateLastLogin() {
	now := time.Now().UTC()
	o.LastLoginAt = &now
	o.UpdatedAt = now
}

func NormalizeFirebaseProvider(firebaseProvider string) OAuthProvider {
	provider := strings.TrimSuffix(firebaseProvider, ".com")

	switch provider {
	case "google":
		return OAuthProviderGoogle
	case "facebook":
		return OAuthProviderFacebook
	case "apple":
		return OAuthProviderApple
	case "github":
		return OAuthProviderGithub
	default:
		return OAuthProvider(provider)
	}
}

func IsValidProvider(provider string) bool {
	validProviders := map[string]bool{
		string(OAuthProviderGoogle):   true,
		string(OAuthProviderFacebook): true,
		string(OAuthProviderApple):    true,
		string(OAuthProviderGithub):   true,
	}
	return validProviders[provider]
}

func GetProviderName(provider OAuthProvider) string {
	switch provider {
	case OAuthProviderGoogle:
		return "Google"
	case OAuthProviderFacebook:
		return "Facebook"
	case OAuthProviderApple:
		return "Apple"
	case OAuthProviderGithub:
		return "GitHub"
	case OAuthProviderTwitter:
		return "Twitter"
	default:
		return "Unknown"
	}
}

func NewOAuthAccount(userID uuid.UUID, firebaseToken map[string]interface{}) *OAuthAccount {
	now := time.Now().UTC()

	firebase := firebaseToken["firebase"].(map[string]interface{})
	provider := firebase["sign_in_provider"].(string)

	identities := firebase["identities"].(map[string]interface{})
	var providerUID string
	if googleIDs, ok := identities["google.com"].([]interface{}); ok && len(googleIDs) > 0 {
		providerUID = googleIDs[0].(string)
	} else if fbIDs, ok := identities["facebook.com"].([]interface{}); ok && len(fbIDs) > 0 {
		providerUID = fbIDs[0].(string)
	}

	email := firebaseToken["email"].(string)
	name := ""
	if n, ok := firebaseToken["name"].(string); ok {
		name = n
	}
	picture := ""
	if p, ok := firebaseToken["picture"].(string); ok {
		picture = p
	}
	emailVerified := false
	if ev, ok := firebaseToken["email_verified"].(bool); ok {
		emailVerified = ev
	}

	return &OAuthAccount{
		ID:              uuid.New(),
		UserID:          userID,
		Provider:        OAuthProvider(provider),
		ProviderUID:     providerUID,
		FirebaseUID:     firebaseToken["user_id"].(string),
		ProviderEmail:   email,
		ProviderName:    name,
		ProviderPicture: picture,
		EmailVerified:   emailVerified,
		LastLoginAt:     &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
