package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
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

	logger.Info(fmt.Sprintf("🚀 Starting web server in Port %v", server.Addr))

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
	r.Use(CustomRecoverer(builder.Logger))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", builder.HealthHandler.Ping)

	r.Post("/auth/registry", builder.AuthHandler.RegistryUser)
	r.Post("/auth/login", builder.AuthHandler.Login)
	r.Post("/auth/firebase/login", builder.AuthHandler.LoginWithFirebase)
	r.Post("/auth/refresh", builder.AuthHandler.RefreshToken)

	r.With(builder.AuthMiddleware.AuthenticateRequest).Route("/", func(r chi.Router) {
		// Auth routes
		r.Post("/auth/logout", builder.AuthHandler.Logout)
		r.Post("/auth/logout-all", builder.AuthHandler.LogoutAllDevices)

		// Profile route
		r.Get("/me", builder.ProfileHandler.GetUserProfile)

		// Health Profile routes
		r.Route("/health-profile", func(r chi.Router) {
			r.Get("/", builder.HealthProfileHandler.GetHealthProfile)
			r.Put("/", builder.HealthProfileHandler.UpdateHealthProfile)
		})

		// Emergency Contacts routes
		r.Route("/emergency-contacts", func(r chi.Router) {
			r.Get("/", builder.EmergencyContactHandler.GetUserEmergencyContacts)
			r.Post("/", builder.EmergencyContactHandler.CreateEmergencyContact)
			r.Get("/{contactID}", builder.EmergencyContactHandler.GetEmergencyContactByID)
			r.Put("/{contactID}", builder.EmergencyContactHandler.UpdateEmergencyContact)
			r.Delete("/{contactID}", builder.EmergencyContactHandler.DeleteEmergencyContact)
			r.Post("/{contactID}/set-primary", builder.EmergencyContactHandler.SetPrimaryContact)
		})

		// Reminder Categories routes
		r.Route("/reminder-categories", func(r chi.Router) {
			r.Get("/", builder.ReminderCategoryHandler.GetCategoriesForUser)
			r.Post("/", builder.ReminderCategoryHandler.CreateUserCategory)
			r.Put("/{category_id}", builder.ReminderCategoryHandler.UpdateCategory)
			r.Delete("/{category_id}", builder.ReminderCategoryHandler.DeleteCategory)
		})

		// Reminders routes
		r.Route("/reminders", func(r chi.Router) {
			r.Post("/", builder.ReminderHandler.CreateReminder)
			r.Get("/", builder.ReminderHandler.GetUserReminders)
			r.Get("/active", builder.ReminderHandler.GetActiveReminders)
			r.Get("/upcoming", builder.ReminderHandler.GetUpcomingReminders)
			r.Get("/{reminder_id}", builder.ReminderHandler.GetReminderByID)
			r.Put("/{reminder_id}", builder.ReminderHandler.UpdateReminder)
			r.Delete("/{reminder_id}", builder.ReminderHandler.DeleteReminder)
			r.Post("/{reminder_id}/complete", builder.ReminderHandler.CompleteReminder)
			r.Post("/{reminder_id}/snooze", builder.ReminderHandler.SnoozeReminder)
		})

	})

	return r
}

// CustomRecoverer is a middleware that recovers from panics and logs them properly
func CustomRecoverer(logger ports.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					// Log the panic with stack trace
					logger.Error(fmt.Sprintf("PANIC RECOVERED: %v", rvr), fmt.Errorf("%v", rvr))

					// Ensure headers are written
					if !isResponseWritten(w) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						_ = json.NewEncoder(w).Encode(map[string]string{
							"error":   "internal_server_error",
							"message": "An unexpected error occurred",
						})
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Helper to check if response was already written
func isResponseWritten(w http.ResponseWriter) bool {
	// This is a heuristic; not perfect but helps
	// If you have access to a custom ResponseWriter, you can track this better
	return false // For now, always try to write error
}
