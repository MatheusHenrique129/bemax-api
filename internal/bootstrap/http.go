package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MatheusHenrique129/bemax-backend/internal/core/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	_defaultWebServerPort = "8080"
	_shutdownTimeout      = 5000 * time.Millisecond
)

// CreateWebApplication builds and configures the HTTP server with timeouts and optional logging.
// It creates a new web application instance with configured read/write timeouts
// and conditionally enables request/response logging based on configuration settings.
//
// Parameters:
//   - logger: logger instance for HTTP request/response logging
//   - cfg: application configuration containing timeout settings and logging preferences
//
// Returns a configured application ready to register routes and handle HTTP requests,
// or an error if the application fails to initialize.
func CreateWebApplication(logger ports.Logger, cfg ports.Configuration, rt http.Handler) {
	port := cfg.Server.Port
	if port == "" {
		port = _defaultWebServerPort
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      rt,
		ReadTimeout:  time.Millisecond * time.Duration(cfg.Server.AppReadTimeoutMs),
		WriteTimeout: time.Millisecond * time.Duration(cfg.Server.AppWriteTimeoutMs),
		IdleTimeout:  time.Millisecond * time.Duration(cfg.Server.AppIdleTimeoutMs),
	}

	logger.Info(fmt.Sprintf("🚀 Starting web server in Port %v", cfg.Server.Port))

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Web Server Error", err)
		}
	}()

	awaitShutdown(server, logger, _shutdownTimeout)
}

func awaitShutdown(server *http.Server, logger ports.Logger, timeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("Received signal for shutdown: " + sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server Shutdown", err)
	}
	logger.Info("🫦 Server gracefully stopped")
}

// RegisterRoutes registers all HTTP handlers and defines the application's routing structure.
// It sets up two main route groups:
// - /api/v1: for internal stream and BigQueue consumer endpoints
// - /internal/v1: for internal member processing and job management endpoints
//
// The function maps HTTP routes to their corresponding handlers from the AppBuilder,
// establishing the complete HTTP API surface for the application.
//
// Parameters:
//   - app: the application instance to register routes on
//   - builder: the AppBuilder containing all initialized handlers
func RegisterRoutes(builder *AppBuilder) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", builder.HealthHandler.Ping)

	return r
}
