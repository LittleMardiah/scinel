package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LittleMardiah/scinel/internal/util"
)

func newTestRateLimiter(t *testing.T, requestsPerSec int) (*RateLimitMiddleware, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	rl := NewRateLimitMiddleware(rc, requestsPerSec)
	return rl, mr
}

func TestRateLimitMiddleware_getLimitKey(t *testing.T) {
	rl, _ := newTestRateLimiter(t, 1000)

	t.Run("api key header takes priority", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req.Header.Set("X-API-Key", "sk_test_123")
		assert.Equal(t, "sk_test_123", rl.getLimitKey(req))
	})

	t.Run("user context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req = req.WithContext(context.WithValue(req.Context(), UserContextKey, &UserInfo{UserID: "user-123"}))
		assert.Equal(t, "user:user-123", rl.getLimitKey(req))
	})

	t.Run("api key overrides user context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req.Header.Set("X-API-Key", "sk_test_abc")
		req = req.WithContext(context.WithValue(req.Context(), UserContextKey, &UserInfo{UserID: "user-123"}))
		assert.Equal(t, "sk_test_abc", rl.getLimitKey(req))
	})

	t.Run("ip fallback strips port", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		assert.Equal(t, "ip:192.168.1.1", rl.getLimitKey(req))
	})

	t.Run("x-forwarded-for single hop", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req.RemoteAddr = "10.10.10.10:5555"
		req.Header.Set("X-Forwarded-For", "10.0.0.1")
		assert.Equal(t, "ip:10.0.0.1", rl.getLimitKey(req))
	})

	t.Run("x-forwarded-for chain uses first hop only", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req.RemoteAddr = "10.10.10.10:5555"
		req.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
		assert.Equal(t, "ip:1.2.3.4", rl.getLimitKey(req))
	})

	t.Run("x-forwarded-for empty falls back to remote addr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		req.Header.Set("X-Forwarded-For", "")
		assert.Equal(t, "ip:192.168.1.1", rl.getLimitKey(req))
	})
}

func TestRateLimitMiddleware_checkLimit(t *testing.T) {
	t.Run("within limit allowed", func(t *testing.T) {
		rl, _ := newTestRateLimiter(t, 1000)
		allowed, remaining, retryAfter, err := rl.checkLimit(context.Background(), "test-user")
		assert.NoError(t, err)
		assert.True(t, allowed)
		assert.Equal(t, 999, remaining)
		assert.Equal(t, 0, retryAfter)
	})

	t.Run("multiple requests within limit", func(t *testing.T) {
		rl, _ := newTestRateLimiter(t, 1000)
		for i := 0; i < 5; i++ {
			allowed, _, _, err := rl.checkLimit(context.Background(), "burst")
			assert.NoError(t, err)
			assert.True(t, allowed)
		}
	})

	t.Run("over limit blocked with retryAfter", func(t *testing.T) {
		rl, _ := newTestRateLimiter(t, 2)
		// First two allowed
		allowed, _, _, err := rl.checkLimit(context.Background(), "limited-user")
		assert.True(t, allowed)
		assert.NoError(t, err)
		allowed, _, _, err = rl.checkLimit(context.Background(), "limited-user")
		assert.True(t, allowed)
		assert.NoError(t, err)
		// Third blocked
		allowed, remaining, retryAfter, err := rl.checkLimit(context.Background(), "limited-user")
		assert.NoError(t, err)
		assert.False(t, allowed)
		assert.Equal(t, 0, remaining)
		assert.GreaterOrEqual(t, retryAfter, 1)
	})
}

func TestRateLimitMiddleware_Apply(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("within limit passes through with headers", func(t *testing.T) {
		rl, _ := newTestRateLimiter(t, 1000)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		rr := httptest.NewRecorder()
		rl.Apply(okHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "1000", rr.Header().Get("X-RateLimit-Limit"))
		assert.NotEmpty(t, rr.Header().Get("X-RateLimit-Remaining"))
		assert.NotEmpty(t, rr.Header().Get("X-RateLimit-Reset"))
	})

	t.Run("over limit blocks with 429 and Retry-After", func(t *testing.T) {
		rl, _ := newTestRateLimiter(t, 1)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		req.Header.Set("X-API-Key", "sk_block")

		// Consume the single allowed request
		rr1 := httptest.NewRecorder()
		rl.Apply(okHandler).ServeHTTP(rr1, req)
		assert.Equal(t, http.StatusOK, rr1.Code)

		// Second request blocked
		rr2 := httptest.NewRecorder()
		rl.Apply(okHandler).ServeHTTP(rr2, req)
		assert.Equal(t, http.StatusTooManyRequests, rr2.Code)
		assert.NotEmpty(t, rr2.Header().Get("Retry-After"))
	})

	t.Run("redis down continues request", func(t *testing.T) {
		rl := NewRateLimitMiddleware(redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"}), 1000)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		rr := httptest.NewRecorder()
		rl.Apply(okHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestRateLimitMiddleware_New(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()
	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	rl := NewRateLimitMiddleware(rc, 500)
	assert.NotNil(t, rl)
	assert.Equal(t, 500, rl.requestsPerSec)
	assert.Equal(t, 1*time.Second, rl.windowSize)
	_ = util.GenerateRequestID
}
