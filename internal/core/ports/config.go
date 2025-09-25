package ports

// Config defines the contract for configuration management
type Config interface {
	LoadConfiguration() (Configuration, error)
	Reload() error
}

// Configuration represents the application configuration structure
type Configuration struct {
	LogLevel string       `mapstructure:"log_level" json:"log_level"`
	Server   ServerConfig `mapstructure:"server" json:"server"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port              string `mapstructure:"port" json:"port"`
	AppIdleTimeoutMs  int    `mapstructure:"app_idle_timeout_ms" json:"app_idle_timeout_ms"`
	AppReadTimeoutMs  int    `mapstructure:"app_read_timeout_ms" json:"app_read_timeout_ms"`
	AppWriteTimeoutMs int    `mapstructure:"app_write_timeout_ms" json:"app_write_timeout_ms"`
}

// Default values for configuration
const (
	DefaultLogLevel          = "debug"
	DefaultPort              = "8080"
	DefaultAppIdleTimeoutMs  = 70
	DefaultAppReadTimeoutMs  = 200
	DefaultAppWriteTimeoutMs = 200
)
