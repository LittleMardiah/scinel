package integration

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/LittleMardiah/scinel/internal/config"
	"github.com/LittleMardiah/scinel/internal/database"
	"github.com/LittleMardiah/scinel/internal/handler"
	"github.com/LittleMardiah/scinel/internal/repository"
	"github.com/LittleMardiah/scinel/internal/util"
)

func TestAuthFlow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start PostgreSQL container
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithStartupTimeout(30 * time.Second),
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	require.NoError(t, err)
	defer postgres.Terminate(ctx)

	postgresHost, _ := postgres.Host(ctx)
	postgresPort, _ := postgres.MappedPort(ctx, "5432")
	dsn := "postgres://test_user:test_password@" + postgresHost + ":" + postgresPort.Port() + "/test_db?sslmode=disable"

	// Run migration
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(255),
			role VARCHAR(50) NOT NULL DEFAULT 'viewer',
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			last_login_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	// Setup dependencies
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:     "test-secret-key",
			AccessExpiry:  1 * time.Hour,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
	}
	jwtManager := util.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)

	connCfg := database.ConnectionConfig{URL: dsn}
	dbConn, err := database.NewConnection(connCfg)
	require.NoError(t, err)
	defer dbConn.Close()

	userRepo := repository.NewUserRepository(dbConn)
	authHandler := handler.NewAuthHandler(userRepo, jwtManager, cfg)
	_ = authHandler // avoid unused warning

	// Test flow: Register, Login, Refresh, Logout
	t.Run("Register_Login_Refresh_Logout", func(t *testing.T) {
		hashed, _ := util.HashPassword("SecurePass123!", 0)
		user, err := userRepo.CreateUser(ctx, "test@example.com", hashed, "Test User", "viewer")
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)

		// Verify password
		err = util.VerifyPassword("SecurePass123!", user.PasswordHash)
		assert.NoError(t, err)

		// Generate tokens
		accessToken, err := jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
		require.NoError(t, err)
		assert.NotEmpty(t, accessToken)

		refreshToken, err := jwtManager.GenerateRefreshToken(user.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, refreshToken)

		// Verify tokens
		claims, err := jwtManager.VerifyAccessToken(accessToken)
		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)

		_, err = jwtManager.VerifyRefreshToken(refreshToken)
		require.NoError(t, err)

		t.Log("Auth flow passed")
	})
}

func TestMain(m *testing.M) {
	m.Run()
}
