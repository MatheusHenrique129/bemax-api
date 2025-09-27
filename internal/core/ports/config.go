package ports

import "time"

// Config defines the contract for configuration management
type Config interface {
	LoadConfiguration() (Configuration, error)
	Reload() error
}

// Configuration represents the application configuration structure
type Configuration struct {
	LogLevel string  `mapstructure:"log_level" json:"log_level"`
	Server   Server  `mapstructure:"server" json:"server"`
	Auth     Auth    `mapstructure:"auth"`
	Storage  Storage `mapstructure:"storage"`
}

type Server struct {
	Port              string `mapstructure:"port" json:"port"`
	AppIdleTimeoutMs  int    `mapstructure:"app_idle_timeout_ms" json:"app_idle_timeout_ms"`
	AppReadTimeoutMs  int    `mapstructure:"app_read_timeout_ms" json:"app_read_timeout_ms"`
	AppWriteTimeoutMs int    `mapstructure:"app_write_timeout_ms" json:"app_write_timeout_ms"`
}

type Auth struct {
	JWT JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	TTL    time.Duration `mapstructure:"ttl"`
}

type Storage struct {
	MySQL MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	DriverName   string `mapstructure:"driver_name"`
	DBName       string `mapstructure:"db_name"`
	HostName     string `mapstructure:"hostname"`
	UserName     string `mapstructure:"user_name"`
	UserPassword string `mapstructure:"user_password"`
}

// Default values for configuration
const (
	DefaultLogLevel          = "debug"
	DefaultPort              = "8080"
	DefaultAppIdleTimeoutMs  = 70
	DefaultAppReadTimeoutMs  = 200
	DefaultAppWriteTimeoutMs = 200
)
