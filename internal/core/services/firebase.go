package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/persistence/mysql"
	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services/dto"
	"github.com/google/uuid"
)

var (
	ErrInvalidFirebaseToken  = errors.New("invalid firebase token")
	ErrFirebaseUserNotFound  = errors.New("firebase user not found")
	ErrUnsupportedProvider   = errors.New("unsupported OAuth provider")
	ErrOAuthAccountNotLinked = errors.New("oauth account not linked to any user")
	ErrUserProfileIncomplete = errors.New("user profile incomplete, CPF required")
)

type FirebaseService interface {
	LoginWithFirebase(ctx context.Context, idToken, ipAddress, userAgent, deviceInfo string) (dto.FirebaseLoginResponse, apierrors.RestError)
	VerifyFirebaseToken(ctx context.Context, idToken string) (*firebaseAuth.Token, apierrors.RestError)
	GetUserByOAuthAccount(ctx context.Context, provider domain.OAuthProvider, providerUID string) (*domain.User, *domain.OAuthAccount, apierrors.RestError)
	LinkOAuthAccount(ctx context.Context, userID uuid.UUID, firebaseToken *firebaseAuth.Token) apierrors.RestError
}

type firebaseService struct {
	logger           ports.Logger
	firebaseClient   *firebaseAuth.Client
	config           ports.FirebaseConfig
	userService      UserService
	roleService      RoleService
	sessionService   SessionService
	authTokenService AuthTokenService
	oauthAccountRepo ports.OAuthAccountRepository
}

// VerifyFirebaseToken verifies a Firebase ID token and returns the token claims
func (f *firebaseService) VerifyFirebaseToken(ctx context.Context, idToken string) (*firebaseAuth.Token, apierrors.RestError) {
	if f.firebaseClient == nil {
		f.logger.Error("Firebase client not initialized", nil)
		return nil, apierrors.NewInternalServerRestError("Firebase authentication not configured", nil)
	}

	token, err := f.firebaseClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		f.logger.Error("Failed to verify Firebase token", err)
		return nil, apierrors.NewUnauthorizedRestError(err.Error())
	}

	return token, nil
}

// GetUserByOAuthAccount gets a user by OAuth account
func (f *firebaseService) GetUserByOAuthAccount(ctx context.Context, provider domain.OAuthProvider, providerUID string) (*domain.User, *domain.OAuthAccount, apierrors.RestError) {
	oauthAccount, err := f.oauthAccountRepo.FindByFirebaseUID(ctx, providerUID)
	if err != nil {
		if errors.Is(err, mysql.ErrOAuthAccountNotFound) {
			return nil, nil, apierrors.NewNotFoundRestError("oauth account not found")
		}
		f.logger.Error("Failed to find OAuth account", err)
		return nil, nil, apierrors.NewInternalServerRestError("failed to find oauth account", err)
	}

	user, userErr := f.userService.GetUserByID(ctx, oauthAccount.UserID)
	if userErr != nil {
		return nil, nil, userErr
	}

	return &user, oauthAccount, nil
}

// LinkOAuthAccount links an OAuth account to an existing user
func (f *firebaseService) LinkOAuthAccount(ctx context.Context, userID uuid.UUID, firebaseToken *firebaseAuth.Token) apierrors.RestError {
	// Extract Firebase UID
	firebaseUID := firebaseToken.UID
	if firebaseUID == "" {
		return apierrors.NewBadRequestRestError("firebase UID not found in token")
	}

	// Check if OAuth account already exists
	existingAccount, err := f.oauthAccountRepo.FindByProviderAndUID(ctx, domain.OAuthProviderGoogle, firebaseUID)
	if err == nil && existingAccount != nil {
		// Account already linked to another user
		if existingAccount.UserID != userID {
			return apierrors.NewBadRequestRestError("this OAuth account is already linked to another user")
		}
		// Already linked to this user
		return nil
	}

	// Create new OAuth account
	oauthAccount := domain.NewOAuthAccount(userID, firebaseToken.Claims)

	// Save OAuth account
	if err := f.oauthAccountRepo.Create(ctx, oauthAccount); err != nil {
		f.logger.Error("Failed to create OAuth account", err)
		return apierrors.NewInternalServerRestError("failed to link oauth account", err)
	}

	f.logger.Info(fmt.Sprintf("OAuth account linked successfully for user: %s", userID))
	return nil
}

// LoginWithFirebase authenticates a user using Firebase ID token
func (f *firebaseService) LoginWithFirebase(ctx context.Context, idToken, ipAddress, userAgent, deviceInfo string) (dto.FirebaseLoginResponse, apierrors.RestError) {
	token, err := f.VerifyFirebaseToken(ctx, idToken)
	if err != nil {
		f.logger.Error("Failed to verify Firebase token", err)
		return dto.FirebaseLoginResponse{}, err
	}

	firebaseUID := token.UID
	if firebaseUID == "" {
		return dto.FirebaseLoginResponse{}, apierrors.NewBadRequestRestError("firebase UID not found in token")
	}

	provider := domain.NormalizeFirebaseProvider(token.Firebase.SignInProvider)
	if !domain.IsValidProvider(string(provider)) {
		f.logger.Error(fmt.Sprintf("Unsupported OAuth provider: %s", provider), err)
		return dto.FirebaseLoginResponse{}, apierrors.NewBadRequestRestError(fmt.Sprintf("unsupported provider: %s", provider))
	}

	user, oauthAccount, userErr := f.GetUserByOAuthAccount(ctx, provider, firebaseUID)
	if userErr != nil && userErr.Status() == http.StatusNotFound {
		emailVerified := token.Claims["email_verified"].(bool)
		email := token.Claims["email"].(string)
		name := token.Claims["name"].(string)
		picture := ""
		if p, ok := token.Claims["picture"].(string); ok {
			picture = p
		}

		// Create new OAuth user
		newUser, err := domain.NewOAuthUser(email, name, emailVerified)
		if err != nil {
			return dto.FirebaseLoginResponse{}, err
		}

		if picture != "" {
			newUser.ProfilePicture = picture
		}

		if newUser, userErr = f.userService.CreateUserOAuth(ctx, newUser); userErr != nil {
			f.logger.Error("Failed to create OAuth user", userErr)
			return dto.FirebaseLoginResponse{}, apierrors.NewInternalServerRestError("failed to create user", userErr)
		}

		oauthAccount = domain.NewOAuthAccount(newUser.ID, token.Claims)
		oauthAccount.Provider = provider

		if createErr := f.oauthAccountRepo.Create(ctx, oauthAccount); createErr != nil {
			f.logger.Error("Failed to create OAuth account", createErr)
			return dto.FirebaseLoginResponse{}, apierrors.NewInternalServerRestError("failed to create OAuth account", createErr)
		}

		user = &newUser
	} else if userErr != nil {
		return dto.FirebaseLoginResponse{}, userErr
	}

	if !user.IsActive() {
		f.logger.Error(fmt.Sprintf("Inactive user attempted Firebase login: %s", user.Email), nil)
		return dto.FirebaseLoginResponse{}, apierrors.NewUnauthorizedRestError("user is inactive")
	}

	oauthAccount.ProviderEmail = token.Claims["email"].(string)
	if name, ok := token.Claims["name"].(string); ok {
		oauthAccount.ProviderName = name
	}
	if picture, ok := token.Claims["picture"].(string); ok {
		oauthAccount.ProviderPicture = picture
	}
	if verified, ok := token.Claims["email_verified"].(bool); ok {
		oauthAccount.EmailVerified = verified
	}

	oauthAccount.UpdateLastLogin()

	if updateErr := f.oauthAccountRepo.Update(ctx, *oauthAccount); updateErr != nil {
		f.logger.Error("Failed to update OAuth account", updateErr)
		// Don't fail login, just log
	}

	// Create session
	session, sessionErr := f.sessionService.CreateSession(ctx, user.ID, deviceInfo, ipAddress, userAgent)
	if sessionErr != nil {
		f.logger.Error("Failed to create session for Firebase login", errors.New(sessionErr.Message()))
		return dto.FirebaseLoginResponse{}, sessionErr
	}

	response, tokenErr := f.authTokenService.GenerateTokensForSession(ctx, user, session)
	if tokenErr != nil {
		f.logger.Error("Failed to generate tokens for Firebase login", tokenErr)
		return dto.FirebaseLoginResponse{}, tokenErr
	}

	f.logger.Info(fmt.Sprintf("Firebase login successful for user: %s", user.Email))
	return response, nil
}

func NewFirebaseService(
	logger ports.Logger,
	config ports.FirebaseConfig,
	firebaseClient *firebaseAuth.Client,
	userService UserService,
	roleService RoleService,
	sessionService SessionService,
	authTokenService AuthTokenService,
	oauthAccountRepo ports.OAuthAccountRepository,
) FirebaseService {
	return &firebaseService{
		logger:           logger,
		config:           config,
		firebaseClient:   firebaseClient,
		userService:      userService,
		roleService:      roleService,
		sessionService:   sessionService,
		authTokenService: authTokenService,
		oauthAccountRepo: oauthAccountRepo,
	}
}
