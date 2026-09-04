package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LittleMardiah/scinel/internal/middleware"
)

func newTestForecastHandler(t *testing.T) (*ForecastHandler, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	return NewForecastHandler(nil, nil, rc), mr
}

func TestForecastHandler_QueryForecasts(t *testing.T) {
	t.Run("success generates placeholder forecasts", func(t *testing.T) {
		h, _ := newTestForecastHandler(t)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/forecasts?days=7", nil)
		rr := httptest.NewRecorder()
		h.QueryForecasts(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		forecasts := data["forecasts"].([]interface{})
		assert.Equal(t, 7, len(forecasts))
	})

	t.Run("returns cached data if present", func(t *testing.T) {
		h, mr := newTestForecastHandler(t)
		// Pre-populate cache with valid JSON response body
		mr.Set("cache:forecasts:30", `{"code":200,"message":"OK","data":{"forecasts":[]}}`)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/forecasts", nil)
		rr := httptest.NewRecorder()
		h.QueryForecasts(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Contains(t, rr.Body.String(), "forecasts")
	})
}

func TestForecastHandler_GenerateForecasts(t *testing.T) {
	t.Run("non-admin forbidden", func(t *testing.T) {
		h, _ := newTestForecastHandler(t)
		user := &middleware.UserInfo{UserID: "11111111-1111-1111-1111-111111111111", Email: "u@example.com", Role: "viewer"}
		req := requestWithUser(http.MethodPost, "/api/v1/forecasts/generate", "", user)
		rr := httptest.NewRecorder()
		h.GenerateForecasts(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("admin accepted", func(t *testing.T) {
		h, _ := newTestForecastHandler(t)
		user := &middleware.UserInfo{UserID: "11111111-1111-1111-1111-111111111111", Email: "a@example.com", Role: "admin"}
		req := requestWithUser(http.MethodPost, "/api/v1/forecasts/generate", "", user)
		rr := httptest.NewRecorder()
		h.GenerateForecasts(rr, req)
		assert.Equal(t, http.StatusAccepted, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, "processing", data["status"])
	})

	t.Run("missing user context forbidden", func(t *testing.T) {
		h, _ := newTestForecastHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/forecasts/generate", nil)
		rr := httptest.NewRecorder()
		h.GenerateForecasts(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}

func TestForecastHandler_GetAccuracy(t *testing.T) {
	h, _ := newTestForecastHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/forecasts/accuracy", nil)
	rr := httptest.NewRecorder()
	h.GetAccuracy(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	resp := decodeResponse(t, rr)
	data := resp.Data.(map[string]interface{})
	accuracy := data["accuracy"].([]interface{})
	assert.Equal(t, 2, len(accuracy))
}

func TestForecastHandler_New(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()
	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	h := NewForecastHandler(nil, nil, rc)
	assert.NotNil(t, h)
	assert.Equal(t, rc, h.redisClient)
}
