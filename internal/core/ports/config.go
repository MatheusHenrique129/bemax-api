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
	Auth     Auth    `mapstructure:"auth" json:"auth"`
	Storage  Storage `mapstructure:"storage" json:"storage"`
}

type Server struct {
	Port              string `mapstructure:"port" json:"port"`
	AppIdleTimeoutMs  int    `mapstructure:"app_idle_timeout_ms" json:"app_idle_timeout_ms"`
	AppReadTimeoutMs  int    `mapstructure:"app_read_timeout_ms" json:"app_read_timeout_ms"`
	AppWriteTimeoutMs int    `mapstructure:"app_write_timeout_ms" json:"app_write_timeout_ms"`
}

type Auth struct {
	JWT      JWTConfig      `mapstructure:"jwt" json:"jwt"`
	Firebase FirebaseConfig `mapstructure:"firebase" json:"firebase"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret" json:"secret"`
	TTL    time.Duration `mapstructure:"ttl" json:"ttl"`
}

type FirebaseConfig struct {
	ProjectID       string `mapstructure:"project_id" json:"project_id"`
	CredentialsPath string `mapstructure:"credentials_path" json:"credentials_path"`
}

type Storage struct {
	MySQL MysqlConfig `mapstructure:"mysql" json:"mysql"`
}

type MysqlConfig struct {
	DriverName   string `mapstructure:"driver_name" json:"driver_name"`
	DBName       string `mapstructure:"db_name" json:"db_name"`
	HostName     string `mapstructure:"hostname" json:"host_name"`
	UserName     string `mapstructure:"user_name" json:"user_name"`
	UserPassword string `mapstructure:"user_password" json:"user_password"`
}

// Default values for configuration
const (
	DefaultLogLevel          = "debug"
	DefaultPort              = "8080"
	DefaultAppIdleTimeoutMs  = 70
	DefaultAppReadTimeoutMs  = 200
	DefaultAppWriteTimeoutMs = 200

	DefaultFirebaseProjectID       = "myapp"
	DefaultFirebaseCredentialsPath = "myapp"
)
