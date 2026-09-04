package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LittleMardiah/scinel/internal/config"
	"github.com/LittleMardiah/scinel/internal/database"
	"github.com/LittleMardiah/scinel/internal/handler"
	"github.com/LittleMardiah/scinel/internal/repository"
	"github.com/LittleMardiah/scinel/internal/util"
)

// TestAuthFlow_Integration verifies the full auth flow (Register -> Login -> Refresh)
// using a sqlmock-backed database so it runs without external dependencies (Docker/Postgres).
func TestAuthFlow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// --- Setup sqlmock-backed database ---
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	userRepo := repository.NewUserRepository(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:     "test-secret-key",
			AccessExpiry:  1 * time.Hour,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
	}
	jwtManager := util.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)

	authHandler := handler.NewAuthHandler(userRepo, jwtManager, cfg)
	require.NotNil(t, authHandler)

	userID := uuid.New()

	t.Run("Register", func(t *testing.T) {
		// Mock the CreateUser query
		now := time.Now().UTC()
		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "full_name", "role", "status", "last_login_at", "created_at", "updated_at",
		}).AddRow(userID, "test@example.com", "hashed", "Test User", "viewer", "active", nil, now, now)

		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs("test@example.com", sqlmock.AnyArg(), "Test User", "viewer").
			WillReturnRows(rows)

		// Build request body
		body := `{"email":"test@example.com","password":"SecurePass123!","full_name":"Test User"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		authHandler.Register(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var resp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "User registered successfully", resp.Message)
	})

	t.Run("Login", func(t *testing.T) {
		hashed, err := util.HashPassword("SecurePass123!", 4)
		require.NoError(t, err)

		now := time.Now().UTC()
		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "full_name", "role", "status", "last_login_at", "created_at", "updated_at",
		}).AddRow(userID, "test@example.com", hashed, "Test User", "viewer", "active", nil, now, now)

		mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, status, last_login_at, created_at, updated_at\s+FROM users WHERE email = \$1`).
			WithArgs("test@example.com").
			WillReturnRows(rows)

		mock.ExpectExec(`UPDATE users SET last_login_at`).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		body := `{"email":"test@example.com","password":"SecurePass123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		authHandler.Login(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				TokenType    string `json:"token_type"`
			} `json:"data"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "Login successful", resp.Message)
		assert.NotEmpty(t, resp.Data.AccessToken)
		assert.NotEmpty(t, resp.Data.RefreshToken)
		assert.Equal(t, "Bearer", resp.Data.TokenType)

		// Verify the generated access token is valid
		claims, err := jwtManager.VerifyAccessToken(resp.Data.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, "test@example.com", claims.Email)
		assert.Equal(t, "viewer", claims.Role)
	})

	t.Run("Refresh", func(t *testing.T) {
		// Use a real refresh token for the known user
		refreshToken, err := jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)

		now := time.Now().UTC()
		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "full_name", "role", "status", "last_login_at", "created_at", "updated_at",
		}).AddRow(userID, "test@example.com", "hashed", "Test User", "viewer", "active", nil, now, now)

		mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, status, last_login_at, created_at, updated_at\s+FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		authHandler.Refresh(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				AccessToken string `json:"access_token"`
				TokenType   string `json:"token_type"`
			} `json:"data"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "Token refreshed successfully", resp.Message)
		assert.NotEmpty(t, resp.Data.AccessToken)

		claims, err := jwtManager.VerifyAccessToken(resp.Data.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})

	t.Run("Logout", func(t *testing.T) {
		refreshToken, err := jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)

		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		authHandler.Logout(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Login_invalid_password", func(t *testing.T) {
		hashed, err := util.HashPassword("SecurePass123!", 4)
		require.NoError(t, err)

		now := time.Now().UTC()
		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "full_name", "role", "status", "last_login_at", "created_at", "updated_at",
		}).AddRow(userID, "test@example.com", hashed, "Test User", "viewer", "active", nil, now, now)

		mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, status, last_login_at, created_at, updated_at\s+FROM users WHERE email = \$1`).
			WithArgs("test@example.com").
			WillReturnRows(rows)

		body := `{"email":"test@example.com","password":"WrongPass999!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		authHandler.Login(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Register_duplicate_email", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs("dup@example.com", sqlmock.AnyArg(), "Dup User", "viewer").
			WillReturnError(errors.New("duplicate key value violates unique constraint"))

		body := `{"email":"dup@example.com","password":"SecurePass123!","full_name":"Dup User"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		authHandler.Register(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled sqlmock expectations: %v", err)
	}
}
