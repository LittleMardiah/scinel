package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRBACMiddleware(t *testing.T) {
	t.Run("role matches required", func(t *testing.T) {
		userInfo := &UserInfo{
			UserID: "user-123",
			Email:  "test@example.com",
			Role:   "admin",
		}
		ctx := context.WithValue(context.Background(), UserContextKey, userInfo)
		req := httptest.NewRequest("GET", "/api/v1/admin", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler := RBACMiddleware("admin", "manager")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("role mismatch returns 403", func(t *testing.T) {
		userInfo := &UserInfo{
			UserID: "user-123",
			Email:  "test@example.com",
			Role:   "viewer",
		}
		ctx := context.WithValue(context.Background(), UserContextKey, userInfo)
		req := httptest.NewRequest("GET", "/api/v1/admin", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler := RBACMiddleware("admin", "manager")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("no user context returns 401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin", nil)
		rr := httptest.NewRecorder()

		handler := RBACMiddleware("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
