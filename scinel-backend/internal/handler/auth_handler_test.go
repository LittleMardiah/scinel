package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/LittleMardiah/scinel/internal/config"
	"github.com/LittleMardiah/scinel/internal/model"
	"github.com/LittleMardiah/scinel/internal/util"
)

func newTestAuthHandler() (*AuthHandler, *MockUserRepository, *config.Config) {
	userRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:     "test-secret-key",
			AccessExpiry:  1 * time.Hour,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
	}
	jwtManager := util.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	return NewAuthHandler(userRepo, jwtManager, cfg), userRepo, cfg
}

func decodeResponse(t *testing.T, rr *httptest.ResponseRecorder) util.ResponseWrapper {
	t.Helper()
	var resp util.ResponseWrapper
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}

func TestAuthHandler_Register(t *testing.T) {
	now := time.Now().UTC()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()

		createdUser := &model.User{
			ID:        userID,
			Email:     "test@example.com",
			FullName:  "Test User",
			Role:      "viewer",
			Status:    "active",
			CreatedAt: now,
		}
		userRepo.On("CreateUser", mock.Anything, "test@example.com", mock.Anything, "Test User", "viewer").Return(createdUser, nil)

		body := `{"email":"test@example.com","password":"password123","full_name":"Test User"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		h.Register(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		resp := decodeResponse(t, rr)
		assert.Equal(t, "User registered successfully", resp.Message)
		userRepo.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader("{invalid"))
		rr := httptest.NewRecorder()
		h.Register(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing email", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		body := `{"email":"","password":"password123","full_name":"Test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Register(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("short password", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		body := `{"email":"a@b.com","password":"short","full_name":"Test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Register(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing full name", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		body := `{"email":"a@b.com","password":"password123","full_name":""}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Register(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("duplicate email returns conflict", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		userRepo.On("CreateUser", mock.Anything, "dup@example.com", mock.Anything, "Test", "viewer").Return((*model.User)(nil), errors.New("duplicate key value violates unique constraint"))
		body := `{"email":"dup@example.com","password":"password123","full_name":"Test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Register(rr, req)
		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("repository error returns internal error", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		userRepo.On("CreateUser", mock.Anything, "err@example.com", mock.Anything, "Test", "viewer").Return((*model.User)(nil), errors.New("db down"))
		body := `{"email":"err@example.com","password":"password123","full_name":"Test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Register(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	now := time.Now().UTC()
	userID := uuid.New()

	hashed, err := util.HashPassword("password123", 4)
	require.NoError(t, err)

	activeUser := &model.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: hashed,
		FullName:     "Test User",
		Role:         "admin",
		Status:       "active",
		CreatedAt:    now,
	}

	t.Run("success", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		userRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(activeUser, nil)
		userRepo.On("UpdateLastLogin", mock.Anything, userID).Return(nil)

		body := `{"email":"test@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Login(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		resp := decodeResponse(t, rr)
		assert.Equal(t, "Login successful", resp.Message)
		userRepo.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader("{"))
		rr := httptest.NewRecorder()
		h.Login(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing credentials", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		body := `{"email":"","password":""}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Login(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		userRepo.On("GetUserByEmail", mock.Anything, "nouser@example.com").Return((*model.User)(nil), errors.New("user not found"))
		body := `{"email":"nouser@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Login(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("wrong password", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		userRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(activeUser, nil)
		body := `{"email":"test@example.com","password":"wrongpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Login(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("suspended account", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		suspended := *activeUser
		suspended.Status = "suspended"
		userRepo.On("GetUserByEmail", mock.Anything, "susp@example.com").Return(&suspended, nil)
		body := `{"email":"susp@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Login(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("pending email verification", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		pending := *activeUser
		pending.Status = "pending_email_verification"
		userRepo.On("GetUserByEmail", mock.Anything, "pend@example.com").Return(&pending, nil)
		body := `{"email":"pend@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Login(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("invalid account status", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		invalid := *activeUser
		invalid.Status = "bogus"
		userRepo.On("GetUserByEmail", mock.Anything, "bogus@example.com").Return(&invalid, nil)
		body := `{"email":"bogus@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Login(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	now := time.Now().UTC()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		refreshToken, err := h.jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)

		activeUser := &model.User{
			ID:        userID,
			Email:     "test@example.com",
			FullName:  "Test User",
			Role:      "viewer",
			Status:    "active",
			CreatedAt: now,
		}
		userRepo.On("GetUserByID", mock.Anything, userID).Return(activeUser, nil)

		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		h.Refresh(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		resp := decodeResponse(t, rr)
		assert.Equal(t, "Token refreshed successfully", resp.Message)
		data := resp.Data.(map[string]interface{})
		assert.NotEmpty(t, data["access_token"])
		userRepo.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader("{"))
		rr := httptest.NewRecorder()
		h.Refresh(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing token", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		body := `{"refresh_token":""}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Refresh(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		h, _, _ := newTestAuthHandler()
		body := `{"refresh_token":"invalid-token"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Refresh(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		refreshToken, err := h.jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)
		userRepo.On("GetUserByID", mock.Anything, userID).Return((*model.User)(nil), errors.New("user not found"))
		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		h.Refresh(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("inactive user", func(t *testing.T) {
		h, userRepo, _ := newTestAuthHandler()
		refreshToken, err := h.jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)
		inactive := &model.User{ID: userID, Status: "suspended"}
		userRepo.On("GetUserByID", mock.Anything, userID).Return(inactive, nil)
		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		h.Refresh(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	userID := uuid.New()
	h, _, _ := newTestAuthHandler()

	t.Run("success", func(t *testing.T) {
		refreshToken, err := h.jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)
		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		h.Logout(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", strings.NewReader("{"))
		rr := httptest.NewRecorder()
		h.Logout(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing token", func(t *testing.T) {
		body := `{"refresh_token":""}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Logout(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid token still returns 200", func(t *testing.T) {
		body := `{"refresh_token":"invalid-token"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.Logout(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestAuthHandler_New(t *testing.T) {
	ctx := context.Background()
	userRepo := new(MockUserRepository)
	cfg := &config.Config{JWT: config.JWTConfig{SecretKey: "secret", AccessExpiry: time.Hour, RefreshExpiry: time.Hour}}
	jwtManager := util.NewJWTManager("secret", time.Hour, time.Hour)
	h := NewAuthHandler(userRepo, jwtManager, cfg)
	assert.NotNil(t, h)
	assert.NotNil(t, h.userRepo)
	assert.NotNil(t, h.jwtManager)
	assert.NotNil(t, h.cfg)
	_ = ctx
}
