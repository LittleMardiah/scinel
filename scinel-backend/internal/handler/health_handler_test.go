package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LittleMardiah/scinel/internal/database"
)

func newTestHealthHandler(t *testing.T) (*HealthHandler, sqlmock.Sqlmock, *miniredis.Miniredis) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })

	db := &database.DB{DB: sqlDB}

	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	return NewHealthHandler(db, rc), mock, mr
}

func TestHealthHandler_Health(t *testing.T) {
	h, _, _ := newTestHealthHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	h.Health(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	resp := decodeResponse(t, rr)
	assert.Equal(t, "OK", resp.Message)
	data := resp.Data.(map[string]interface{})
	assert.Equal(t, "healthy", data["status"])
}

func TestHealthHandler_Readiness(t *testing.T) {
	t.Run("all healthy", func(t *testing.T) {
		h, mock, _ := newTestHealthHandler(t)
		mock.ExpectPing()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/readiness", nil)
		rr := httptest.NewRecorder()
		h.Readiness(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		assert.Equal(t, "Ready", resp.Message)
		checks := resp.Data.(map[string]interface{})
		assert.Equal(t, true, checks["database"])
		assert.Equal(t, true, checks["cache"])
	})

	t.Run("database down", func(t *testing.T) {
		h, mock, _ := newTestHealthHandler(t)
		mock.ExpectPing().WillReturnError(assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/readiness", nil)
		rr := httptest.NewRecorder()
		h.Readiness(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
		resp := decodeResponse(t, rr)
		checks := resp.Data.(map[string]interface{})
		assert.Equal(t, false, checks["database"])
	})
}

func TestHealthHandler_Liveness(t *testing.T) {
	t.Run("all healthy", func(t *testing.T) {
		h, mock, _ := newTestHealthHandler(t)
		mock.ExpectPing()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/liveness", nil)
		rr := httptest.NewRecorder()
		h.Liveness(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		assert.Equal(t, "Alive", resp.Message)
		checks := resp.Data.(map[string]interface{})
		assert.Equal(t, true, checks["database_connected"])
		assert.Equal(t, true, checks["cache_connected"])
	})

	t.Run("database down", func(t *testing.T) {
		h, mock, _ := newTestHealthHandler(t)
		mock.ExpectPing().WillReturnError(assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/liveness", nil)
		rr := httptest.NewRecorder()
		h.Liveness(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	})
}
