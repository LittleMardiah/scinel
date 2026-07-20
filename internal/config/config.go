package config

import (
"fmt"
"os"
"strconv"
"time"

"github.com/LittleMardiah/scinel/internal/database"
)

// Config holds all application configuration
type Config struct {
Database    database.ConnectionConfig
JWT         JWTConfig
Server      ServerConfig
Logging     LoggingConfig
Environment string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
SecretKey        string
AccessExpiry     time.Duration
RefreshExpiry    time.Duration
PasswordResetTTL time.Duration
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
Port            int
Host            string
ReadTimeout     time.Duration
WriteTimeout    time.Duration
ShutdownTimeout time.Duration
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
Level  string
Format string
Output string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
cfg := &Config{
Environment: getEnv("ENVIRONMENT", "development"),
Database: database.ConnectionConfig{
URL:             getEnv("DATABASE_URL", ""),
MinConnections:  getEnvInt("DB_MIN_CONNECTIONS", 10),
MaxConnections:  getEnvInt("DB_MAX_CONNECTIONS", 50),
ConnectTimeout:  getEnvDuration("DB_CONNECT_TIMEOUT", 5*time.Second),
IdleTimeout:     getEnvDuration("DB_IDLE_TIMEOUT", 30*time.Second),
MaxLifetime:     getEnvDuration("DB_MAX_LIFETIME", 5*time.Minute),
},
JWT: JWTConfig{
SecretKey:        getEnv("JWT_SECRET_KEY", ""),
AccessExpiry:     getEnvDuration("JWT_ACCESS_EXPIRY", 1*time.Hour),
RefreshExpiry:    getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
PasswordResetTTL: getEnvDuration("PASSWORD_RESET_TOKEN_EXPIRY", 15*time.Minute),
},
Server: ServerConfig{
Port:            getEnvInt("SERVER_PORT", 8080),
Host:            getEnv("SERVER_HOST", "0.0.0.0"),
ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 15*time.Second),
},
Logging: LoggingConfig{
Level:  getEnv("LOG_LEVEL", "info"),
Format: getEnv("LOG_FORMAT", "json"),
Output: getEnv("LOG_OUTPUT", "stdout"),
},
}

// Validate required config
if cfg.Database.URL == "" {
return nil, fmt.Errorf("DATABASE_URL is required")
}
if cfg.JWT.SecretKey == "" {
return nil, fmt.Errorf("JWT_SECRET_KEY is required")
}

return cfg, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
if value := os.Getenv(key); value != "" {
return value
}
return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
if value := os.Getenv(key); value != "" {
if intVal, err := strconv.Atoi(value); err == nil {
return intVal
}
}
return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
if value := os.Getenv(key); value != "" {
if duration, err := time.ParseDuration(value); err == nil {
return duration
}
}
return defaultValue
}