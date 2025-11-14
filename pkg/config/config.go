package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ViperConfigAdapter implements the ConfigAdapter interface using Viper
type viperConfigAdapter struct {
	viper  *viper.Viper
	config *ports.Configuration
	mu     sync.RWMutex
}

// LoadConfiguration loads the configuration from various sources
func (c *viperConfigAdapter) LoadConfiguration() (ports.Configuration, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Environment variable configuration to override the config file
	c.viper.SetEnvPrefix("")
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.viper.AutomaticEnv()

	// Set configuration file properties
	c.viper.SetConfigName("config")
	c.viper.SetConfigType("json")

	// Add multiple paths where config file might be located
	c.viper.AddConfigPath(".")
	c.viper.AddConfigPath("./cmd/api")
	c.viper.AddConfigPath("./cmd/api/config.json")

	// Set default values
	c.setDefaults()

	// Try to read config file
	if err := c.viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return ports.Configuration{}, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, continue with defaults and env vars
	}

	// Unmarshal configuration
	var config ports.Configuration
	if err := c.viper.Unmarshal(&config); err != nil {
		return ports.Configuration{}, fmt.Errorf("error unmarshaling config: %w", err)
	}

	c.config = &config
	return config, nil
}

// Reload reloads the configuration from source
func (c *viperConfigAdapter) Reload() error {
	_, err := c.LoadConfiguration()
	return err
}

// setDefaults sets default configuration values
func (c *viperConfigAdapter) setDefaults() {
	// Log level
	c.viper.SetDefault("log_level", ports.DefaultLogLevel)

	// Server configurat ion
	c.viper.SetDefault("server.port", ports.DefaultPort)
	c.viper.SetDefault("server.app_idle_timeout_ms", ports.DefaultAppIdleTimeoutMs)
	c.viper.SetDefault("server.app_read_timeout_ms", ports.DefaultAppReadTimeoutMs)
	c.viper.SetDefault("server.app_write_timeout_ms", ports.DefaultAppWriteTimeoutMs)
	c.viper.SetDefault("auth.firebase.project_id", ports.DefaultFirebaseProjectID)
	c.viper.SetDefault("auth.firebase.credentials_path", ports.DefaultFirebaseCredentialsPath)
}

// WatchConfig enables configuration file watching for hot reload
func (c *viperConfigAdapter) WatchConfig(callback func()) {
	c.viper.WatchConfig()
	c.viper.OnConfigChange(func(e fsnotify.Event) {
		_ = c.Reload()
		if callback != nil {
			callback()
		}
	})
}

// NewViperConfigAdapter creates a new instance of ViperConfigRepository
func NewViperConfigAdapter() ports.Config {
	v := viper.New()

	return &viperConfigAdapter{
		viper: v,
	}
}
