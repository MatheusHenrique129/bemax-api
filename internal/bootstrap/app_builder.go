package bootstrap

import (
	"github.com/MatheusHenrique129/bemax-backend/internal/adapters/handlers"
	"github.com/MatheusHenrique129/bemax-backend/internal/core/ports"
)

// AppBuilder holds all the initialized application dependencies and components.
// It serves as the main dependency container for the application, providing
// access to all configured services, handlers, use cases, and infrastructure components.
type AppBuilder struct {
	Logger ports.Logger

	HealthHandler handlers.HealthHandler
}

func BuildAppDependencies(vLogger ports.Logger, cfg ports.Configuration) *AppBuilder {
	vLogger.Info("Creating database client & repositories")

	// :: Handlers
	vLogger.Info("Creating handlers")
	healthHandler := handlers.NewHealthHandler()

	return &AppBuilder{
		Logger:        vLogger,
		HealthHandler: healthHandler,
	}
}
