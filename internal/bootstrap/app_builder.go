package bootstrap

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/auth"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/handlers"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/handlers/middleware"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/persistence/mysql"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services"
	"google.golang.org/api/option"
)

// AppBuilder holds all the initialized application dependencies and components.
// It serves as the main dependency container for the application, providing
// access to all configured services, handlers, use cases, and infrastructure components.
type AppBuilder struct {
	Logger ports.Logger

	AuthMiddleware middleware.AuthMiddleware

	HealthHandler handlers.HealthHandler
	AuthHandler   handlers.AuthHandler
}

func BuildAppDependencies(vLogger ports.Logger, cfg ports.Configuration) (*AppBuilder, error) {
	vLogger.Info("Creating database client & repositories")

	firebaseClient, err := createFirebaseClient(cfg, vLogger)
	if err != nil {
		return nil, err
	}

	// Database Client
	dbClientAdapter, err := mysql.NewMysql(cfg.Storage.MySQL)
	if err != nil {
		vLogger.Fatal(err.Error())
	}

	// Repositories
	userRepositoryPort := mysql.NewMysqlUserRepository(vLogger, dbClientAdapter)
	roleRepositoryPort := mysql.NewMysqlRoleRepository(vLogger, dbClientAdapter)
	tokenRepositoryPort := mysql.NewMysqlTokenRepository(vLogger, dbClientAdapter)
	sessionRepositoryPort := mysql.NewMysqlSessionRepository(vLogger, dbClientAdapter)
	userRoleRepositoryPort := mysql.NewMysqlUserRoleRepository(vLogger, dbClientAdapter)
	oauthAccountRepositoryPort := mysql.NewMysqlOAuthAccountRepository(vLogger, dbClientAdapter)

	// Adapters
	vLogger.Info(fmt.Sprintf("Creating JWT keys with secret: %s and ttl %v", cfg.Auth.JWT.Secret, cfg.Auth.JWT.TTL))
	jwtAdapter := auth.NewJWTAdapter(vLogger, cfg.Auth.JWT.Secret, cfg.Auth.JWT.TTL)

	// Services
	roleService := services.NewRoleService(vLogger, roleRepositoryPort, userRoleRepositoryPort)
	userService := services.NewUserService(vLogger, userRepositoryPort, roleService)
	sessionService := services.NewSessionService(vLogger, sessionRepositoryPort)
	authTokenService := services.NewAuthTokenService(vLogger, jwtAdapter, userService, roleService, sessionService, tokenRepositoryPort)
	firebaseService := services.NewFirebaseService(vLogger, cfg.Auth.Firebase, firebaseClient, userService, roleService, sessionService, authTokenService, oauthAccountRepositoryPort)

	// Middlewares
	authMiddleware := middleware.NewAuthMiddleware(vLogger, jwtAdapter, authTokenService)

	// Handlers
	vLogger.Info("Creating handlers")
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(vLogger, jwtAdapter, userService, authTokenService, firebaseService)

	return &AppBuilder{
		Logger:         vLogger,
		AuthMiddleware: authMiddleware,
		HealthHandler:  healthHandler,
		AuthHandler:    authHandler,
	}, nil
}

// createFirebaseClient initializes the Firebase Authentication client.
// It creates a Firebase app using the project ID and optional credentials file path
// from the application configuration.
//
// Parameters:
//   - cfg: application configuration containing Firebase project ID and credentials path
//   - logger: logger instance for logging Firebase client initialization events
//
// Returns:
//   - *firebaseAuth.Client: configured Firebase Auth client for verifying ID tokens
//   - error: initialization error if Firebase app or auth client creation fails
//
// If CredentialsPath is provided in config, it will be used for authentication.
// Otherwise, the Firebase SDK will use Application Default Credentials (ADC) from the environment.
func createFirebaseClient(cfg ports.Configuration, logger ports.Logger) (*firebaseAuth.Client, error) {
	ctx := context.Background()
	firebaseConfig := cfg.Auth.Firebase

	var opt option.ClientOption
	if firebaseConfig.CredentialsPath != "" {
		opt = option.WithCredentialsFile(firebaseConfig.CredentialsPath)
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: firebaseConfig.ProjectID,
	}, opt)
	if err != nil {
		logger.Error("failed to initialize Firebase app", err)
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Error("failed to get Firebase Auth client", err)
		return nil, fmt.Errorf("failed to get Firebase Auth client: %w", err)
	}

	return authClient, nil
}
