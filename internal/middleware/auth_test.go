package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/LittleMardiah/scinel/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	secretKey := "test-secret-key"
	jwtManager := util.NewJWTManager(secretKey, 1*time.Hour, 24*time.Hour)

	userID := uuid.New()
	validToken, err := jwtManager.GenerateAccessToken(userID, "test@example.com", "admin")
	require.NoError(t, err)

	// Public paths
	publicPaths := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/api/v1/health",
		"/api/v1/readiness",
		"/api/v1/liveness",
		"/health",
		"/ready",
	}

	t.Run("public paths bypass auth", func(t *testing.T) {
		for _, path := range publicPaths {
			req := httptest.NewRequest("GET", path, nil)
			rr := httptest.NewRecorder()

			handler := AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code, "path %s should bypass auth", path)
		}
	})

	t.Run("protected route without token returns 401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/protected", nil)
		rr := httptest.NewRecorder()

		handler := AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("protected route with valid token returns 200 and sets user context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rr := httptest.NewRecorder()

		var userInfo *UserInfo
		handler := AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ui, ok := r.Context().Value(UserContextKey).(*UserInfo)
			if ok {
				userInfo = ui
			}
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, userInfo)
		assert.Equal(t, userID.String(), userInfo.UserID)
		assert.Equal(t, "test@example.com", userInfo.Email)
		assert.Equal(t, "admin", userInfo.Role)
	})

	t.Run("protected route with invalid token format returns 401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		rr := httptest.NewRecorder()

		handler := AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("protected route with expired token returns 401", func(t *testing.T) {
		expiredManager := util.NewJWTManager(secretKey, -1*time.Hour, 24*time.Hour)
		expiredToken, _ := expiredManager.GenerateAccessToken(userID, "test@example.com", "admin")

		req := httptest.NewRequest("GET", "/api/v1/protected", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		rr := httptest.NewRecorder()

		handler := AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
