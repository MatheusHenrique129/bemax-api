package main

import (
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/bootstrap"
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/MatheusHenrique129/bemax-api/pkg/config"
	"github.com/MatheusHenrique129/bemax-api/pkg/logger"
)

func main() {
	fmt.Println("initializing the application...")

	configAdapter := config.NewViperConfigAdapter()

	cfg, err := configAdapter.LoadConfiguration()
	if err != nil {
		fmt.Printf("error loading configuration: %s", err.Error())
	}

	logLevel := logger.LogLevel(cfg.LogLevel)
	loggerAdapter := logger.NewZapLoggerAdapter(logLevel)

	// Run the application
	loggerAdapter.Info("logger initialized")
	if err := runApp(loggerAdapter, cfg); err != nil {
		loggerAdapter.Fatal("failed to run application: %v", err)
	}
}

// runApp bootstraps all application components, sets up HTTP handlers,
// runs background jobs, and starts the HTTP server.
func runApp(vLogger ports.Logger, cfg ports.Configuration) error {
	// Load all core dependencies (use cases, repositories, handlers)
	appBuilder := bootstrap.BuildAppDependencies(vLogger, cfg)

	// Register all HTTP routes into the web application
	routes := bootstrap.RegisterRoutes(appBuilder)

	// Create web and run application
	bootstrap.CreateWebApplication(vLogger, cfg, routes)
	return nil
}
