package bootstrap

import (
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/adapters/auth"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/handlers"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/handlers/middleware"
	"github.com/MatheusHenrique129/bemax-api/internal/adapters/persistence/mysql"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/internal/core/services"
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

func BuildAppDependencies(vLogger ports.Logger, cfg ports.Configuration) *AppBuilder {
	vLogger.Info("Creating database client & repositories")

	// Database Client
	dbClientAdapter, err := mysql.NewMysql(cfg.Storage.MySQL)
	if err != nil {
		vLogger.Fatal(err.Error())
	}

	// Repositories
	userRepository := mysql.NewMysqlUserRepository(vLogger, dbClientAdapter)
	roleRepositoryPort := mysql.NewMysqlRoleRepository(vLogger, dbClientAdapter)
	userRoleRepositoryPort := mysql.NewMysqlUserRoleRepository(vLogger, dbClientAdapter)

	// Adapters
	vLogger.Info(fmt.Sprintf("Creating JWT keys with secret: %s and ttl %v", cfg.Auth.JWT.Secret, cfg.Auth.JWT.TTL))
	jwtAdapter := auth.NewJWTAdapter(cfg.Auth.JWT.Secret, cfg.Auth.JWT.TTL)

	// Services
	authService := services.NewAuthTokenService(vLogger, jwtAdapter)
	roleService := services.NewRoleService(vLogger, roleRepositoryPort, userRoleRepositoryPort)
	userService := services.NewUserService(vLogger, userRepository, roleService)

	// Middlewares
	authMiddleware := middleware.NewAuthMiddleware(vLogger, jwtAdapter, authService)

	// Handlers
	vLogger.Info("Creating handlers")
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(vLogger, jwtAdapter, userService)

	return &AppBuilder{
		Logger:         vLogger,
		AuthMiddleware: authMiddleware,
		HealthHandler:  healthHandler,
		AuthHandler:    authHandler,
	}
}
