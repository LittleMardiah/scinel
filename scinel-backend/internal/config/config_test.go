package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for test
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENVIRONMENT", "test")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET_KEY")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("ENVIRONMENT")
	}()

	cfg, err := LoadConfig()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "test", cfg.Environment)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.Database.URL)
	assert.Equal(t, 10, cfg.Database.MinConnections)
	assert.Equal(t, 50, cfg.Database.MaxConnections)
	assert.Equal(t, 5*time.Second, cfg.Database.ConnectTimeout)
	assert.Equal(t, 30*time.Second, cfg.Database.IdleTimeout)
	assert.Equal(t, 5*time.Minute, cfg.Database.MaxLifetime)

	assert.Equal(t, "test-secret-key", cfg.JWT.SecretKey)
	assert.Equal(t, 1*time.Hour, cfg.JWT.AccessExpiry)
	assert.Equal(t, 7*24*time.Hour, cfg.JWT.RefreshExpiry)
	assert.Equal(t, 15*time.Minute, cfg.JWT.PasswordResetTTL)

	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, 15*time.Second, cfg.Server.ShutdownTimeout)

	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
	assert.Equal(t, "stdout", cfg.Logging.Output)
}

func TestLoadConfig_MissingDatabaseURL(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", "postgres://test")

	cfg, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "DATABASE_URL is required")
}

func TestLoadConfig_MissingJWTSecret(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Unsetenv("JWT_SECRET_KEY")
	defer os.Unsetenv("DATABASE_URL")

	cfg, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "JWT_SECRET_KEY is required")
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	val := getEnv("TEST_KEY", "default")
	assert.Equal(t, "test_value", val)

	val = getEnv("NON_EXISTENT", "default")
	assert.Equal(t, "default", val)
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	val := getEnvInt("TEST_INT", 0)
	assert.Equal(t, 123, val)

	val = getEnvInt("INVALID", 456)
	assert.Equal(t, 456, val)

	os.Setenv("TEST_INT", "not_number")
	val = getEnvInt("TEST_INT", 789)
	assert.Equal(t, 789, val)
}

func TestGetEnvDuration(t *testing.T) {
	os.Setenv("TEST_DURATION", "10s")
	defer os.Unsetenv("TEST_DURATION")

	val := getEnvDuration("TEST_DURATION", 5*time.Second)
	assert.Equal(t, 10*time.Second, val)

	val = getEnvDuration("INVALID", 5*time.Second)
	assert.Equal(t, 5*time.Second, val)

	os.Setenv("TEST_DURATION", "invalid")
	val = getEnvDuration("TEST_DURATION", 3*time.Second)
	assert.Equal(t, 3*time.Second, val)
}
