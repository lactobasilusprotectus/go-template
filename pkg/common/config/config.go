package config

import (
	"fmt"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"log"
)

const (
	ENV = "ENV"

	LOC = "local"
)

var (
	Global GlobalConfig
)

// HttpConfig is the configuration for the HTTP server
type HttpConfig struct {
	Port    string `env:"HTTP_PORT"`
	TimeOut int    `env:"HTTP_TIMEOUT"`
	Env     string `env:"APP_ENV"`
}

// DatabaseConfig is the configuration for the database
type DatabaseConfig struct {
	Driver   string `env:"DB_DRIVER"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Database string `env:"DB_DATABASE"`

	MaxIdleConnections int `env:"DB_MAX_IDLE_CONNECTIONS"`
	MaxOpenConnections int `env:"DB_MAX_OPEN_CONNECTIONS"`
}

// RedisConfig is the configuration for the Redis server
type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Username string `env:"REDIS_USERNAME"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

// Config is the configuration for the application
type Config struct {
	Http     HttpConfig
	Database DatabaseConfig
	Redis    RedisConfig

	JwtSecretAccessToken  string `env:"JWT_SECRET_KEY_AT"`
	JwtSecretRefreshToken string `env:"JWT_SECRET_KEY_RT"`

	Title       string `env:"APP_TITLE"`
	Description string `env:"APP_DESCRIPTION"`
	URL         string `env:"APP_URL"`
}

// GlobalConfig is the configuration for the application
type GlobalConfig struct {
	GlobalTimeout int

	JwtSecretAccessToken  string
	JwtSecretRefreshToken string
}

// GetFilePath returns the path to the config file
func GetFilePath(env string) string {
	return fmt.Sprintf("etc/config/%s.env", env)
}

// Read ReadConfig reads the config file and returns the configuration
func Read(filepath string) (cfg Config, err error) {
	if err = godotenv.Load(filepath); err != nil {
		return Config{}, fmt.Errorf("error loading %s file", filepath)
	}

	if err = envdecode.StrictDecode(&cfg); err != nil {
		log.Fatalf("Error decoding config file: %s", err)
		return
	}

	// set to global vars
	Global.GlobalTimeout = cfg.Http.TimeOut
	Global.JwtSecretAccessToken = cfg.JwtSecretAccessToken
	Global.JwtSecretRefreshToken = cfg.JwtSecretRefreshToken

	return cfg, nil
}

// ResetGlobalConfig resets global configs to their default values.
func ResetGlobalConfig() {
	Global = GlobalConfig{}
}
