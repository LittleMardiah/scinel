package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/LittleMardiah/scinel/internal/middleware"
	"github.com/LittleMardiah/scinel/internal/model"
	"github.com/LittleMardiah/scinel/internal/service"
)

type mockForecastRepository struct {
	mock.Mock
}

func (m *mockForecastRepository) Create(ctx context.Context, forecast *model.Forecast) (*model.Forecast, error) {
	args := m.Called(ctx, forecast)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Forecast), args.Error(1)
}

func (m *mockForecastRepository) GetByWarehouseAndSKU(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]model.Forecast, error) {
	args := m.Called(ctx, warehouseID, skuID, days)
	return args.Get(0).([]model.Forecast), args.Error(1)
}

func (m *mockForecastRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Forecast, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Forecast), args.Get(1).(int64), args.Error(2)
}

func (m *mockForecastRepository) GetAccuracy(ctx context.Context, warehouseID, skuID uuid.UUID, algorithm model.Algorithm) (*model.ForecastAccuracy, error) {
	args := m.Called(ctx, warehouseID, skuID, algorithm)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ForecastAccuracy), args.Error(1)
}

func (m *mockForecastRepository) UpdateAccuracy(ctx context.Context, accuracy *model.ForecastAccuracy) error {
	args := m.Called(ctx, accuracy)
	return args.Error(0)
}

type forecastTestLogger struct{}

func (*forecastTestLogger) Info(string, ...interface{})  {}
func (*forecastTestLogger) Error(string, ...interface{}) {}
func (*forecastTestLogger) Warn(string, ...interface{})  {}

func newTestForecastHandler(t *testing.T) (*ForecastHandler, *mockForecastRepository, *MockInventoryRepository) {
	t.Helper()
	forecastRepo := new(mockForecastRepository)
	inventoryRepo := new(MockInventoryRepository)
	transactionRepo := new(MockTransactionRepository)
	forecastSvc := service.NewForecastService(forecastRepo, inventoryRepo, transactionRepo, &forecastTestLogger{})
	return NewForecastHandler(forecastRepo, forecastSvc, nil), forecastRepo, inventoryRepo
}

func TestForecastHandler_QueryForecasts(t *testing.T) {
	warehouseID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	skuID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	forecastID := uuid.MustParse("33333333-3333-4333-8333-333333333333")
	generatedAt := time.Date(2026, 9, 25, 10, 0, 0, 0, time.UTC)
	forecast := model.Forecast{
		ID:              forecastID,
		WarehouseID:     warehouseID,
		SKUID:           skuID,
		ForecastDate:    time.Date(2026, 9, 26, 0, 0, 0, 0, time.UTC),
		PredictedDemand: 42,
		LowerBound:      30,
		UpperBound:      55,
		Confidence:      0.95,
		Algorithm:       model.AlgorithmETS,
		GeneratedAt:     generatedAt,
	}

	t.Run("returns forecasts from service", func(t *testing.T) {
		h, forecastRepo, _ := newTestForecastHandler(t)
		forecastRepo.On("GetByWarehouseAndSKU", mock.Anything, warehouseID, skuID, 7).Return([]model.Forecast{forecast}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/forecasts?warehouse_id="+warehouseID.String()+"&sku_id="+skuID.String()+"&days=7", nil)
		rr := httptest.NewRecorder()
		h.QueryForecasts(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data, ok := resp.Data.(map[string]interface{})
		require.True(t, ok)
		forecasts, ok := data["forecasts"].([]interface{})
		require.True(t, ok)
		require.Len(t, forecasts, 1)
		item, ok := forecasts[0].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, forecastID.String(), item["id"])
		assert.Equal(t, warehouseID.String(), item["warehouse_id"])
		assert.Equal(t, skuID.String(), item["sku_id"])
		assert.Equal(t, "2026-09-26", item["forecast_date"])
		assert.Equal(t, float64(42), item["predicted_demand"])
		assert.Equal(t, float64(30), item["lower_bound"])
		assert.Equal(t, float64(55), item["upper_bound"])
		assert.Equal(t, 0.95, item["confidence"])
		assert.Equal(t, string(model.AlgorithmETS), item["algorithm"])
		assert.Equal(t, generatedAt.Format(time.RFC3339Nano), item["generated_at"])
		assert.Equal(t, float64(1), data["total"])
		forecastRepo.AssertExpectations(t)
	})
}

func TestForecastHandler_GenerateForecasts(t *testing.T) {
	t.Run("non-admin forbidden", func(t *testing.T) {
		h, _, _ := newTestForecastHandler(t)
		user := &middleware.UserInfo{UserID: "11111111-1111-1111-1111-111111111111", Email: "u@example.com", Role: "viewer"}
		req := requestWithUser(http.MethodPost, "/api/v1/forecasts/generate", "", user)
		rr := httptest.NewRecorder()
		h.GenerateForecasts(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("admin invokes service", func(t *testing.T) {
		h, _, inventoryRepo := newTestForecastHandler(t)
		inventoryRepo.On("GetAll", mock.Anything, 10000, 0).Return([]model.Inventory{}, int64(0), nil).Once()
		user := &middleware.UserInfo{UserID: "11111111-1111-1111-1111-111111111111", Email: "a@example.com", Role: "admin"}
		req := requestWithUser(http.MethodPost, "/api/v1/forecasts/generate", "", user)
		rr := httptest.NewRecorder()
		h.GenerateForecasts(rr, req)

		require.Equal(t, http.StatusAccepted, rr.Code)
		resp := decodeResponse(t, rr)
		assert.Equal(t, "Forecast batch completed", resp.Message)
		data, ok := resp.Data.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "completed", data["status"])
		inventoryRepo.AssertCalled(t, "GetAll", mock.Anything, 10000, 0)
		inventoryRepo.AssertExpectations(t)
	})

	t.Run("missing user context forbidden", func(t *testing.T) {
		h, _, _ := newTestForecastHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/forecasts/generate", nil)
		rr := httptest.NewRecorder()
		h.GenerateForecasts(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}

func TestForecastHandler_GetAccuracy(t *testing.T) {
	warehouseID := uuid.MustParse("44444444-4444-4444-8444-444444444444")
	skuID := uuid.MustParse("55555555-5555-4555-8555-555555555555")
	evaluationDate := time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)
	etsMAPE, etsRMSE, etsMAE := 8.5, 13.25, 10.75
	arimaMAPE, arimaRMSE, arimaMAE := 9.5, 15.25, 12.75
	metrics := []model.ForecastAccuracy{
		{
			WarehouseID:    warehouseID,
			SKUID:          skuID,
			Algorithm:      model.AlgorithmETS,
			MAPE:           &etsMAPE,
			RMSE:           &etsRMSE,
			MAE:            &etsMAE,
			EvaluationDate: evaluationDate,
		},
		{
			WarehouseID:    warehouseID,
			SKUID:          skuID,
			Algorithm:      model.AlgorithmARIMA,
			MAPE:           &arimaMAPE,
			RMSE:           &arimaRMSE,
			MAE:            &arimaMAE,
			EvaluationDate: evaluationDate,
		},
	}

	h, forecastRepo, _ := newTestForecastHandler(t)
	forecastRepo.On("GetAccuracy", mock.Anything, warehouseID, skuID, model.AlgorithmETS).Return(&metrics[0], nil).Once()
	forecastRepo.On("GetAccuracy", mock.Anything, warehouseID, skuID, model.AlgorithmARIMA).Return(&metrics[1], nil).Once()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/forecasts/accuracy?warehouse_id="+warehouseID.String()+"&sku_id="+skuID.String(), nil)
	rr := httptest.NewRecorder()
	h.GetAccuracy(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	resp := decodeResponse(t, rr)
	data, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	accuracy, ok := data["accuracy"].([]interface{})
	require.True(t, ok)
	require.Len(t, accuracy, 2)

	ets, ok := accuracy[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, warehouseID.String(), ets["warehouse_id"])
	assert.Equal(t, skuID.String(), ets["sku_id"])
	assert.Equal(t, string(model.AlgorithmETS), ets["algorithm"])
	assert.Equal(t, 8.5, ets["mape"])
	assert.Equal(t, 13.25, ets["rmse"])
	assert.Equal(t, 10.75, ets["mae"])
	assert.Equal(t, "2026-09-25", ets["evaluation_date"])

	arima, ok := accuracy[1].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, warehouseID.String(), arima["warehouse_id"])
	assert.Equal(t, skuID.String(), arima["sku_id"])
	assert.Equal(t, string(model.AlgorithmARIMA), arima["algorithm"])
	assert.Equal(t, 9.5, arima["mape"])
	assert.Equal(t, 15.25, arima["rmse"])
	assert.Equal(t, 12.75, arima["mae"])
	assert.Equal(t, "2026-09-25", arima["evaluation_date"])
	forecastRepo.AssertExpectations(t)
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
