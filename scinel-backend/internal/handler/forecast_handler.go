package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/LittleMardiah/scinel/internal/middleware"
	"github.com/LittleMardiah/scinel/internal/repository"
	"github.com/LittleMardiah/scinel/internal/service"
	"github.com/LittleMardiah/scinel/internal/util"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type ForecastHandler struct {
	forecastRepo repository.ForecastRepository
	forecastSvc  *service.ForecastService
	redisClient  *redis.Client
}

func NewForecastHandler(
	forecastRepo repository.ForecastRepository,
	forecastSvc *service.ForecastService,
	redisClient *redis.Client,
) *ForecastHandler {
	return &ForecastHandler{
		forecastRepo: forecastRepo,
		forecastSvc:  forecastSvc,
		redisClient:  redisClient,
	}
}

type forecastResponse struct {
	ID              string    `json:"id"`
	WarehouseID     string    `json:"warehouse_id"`
	SKUID           string    `json:"sku_id"`
	ForecastDate    string    `json:"forecast_date"`
	PredictedDemand int32     `json:"predicted_demand"`
	LowerBound      int32     `json:"lower_bound"`
	UpperBound      int32     `json:"upper_bound"`
	Confidence      float64   `json:"confidence"`
	Algorithm       string    `json:"algorithm"`
	GeneratedAt     time.Time `json:"generated_at"`
}

type accuracyResponse struct {
	WarehouseID    string  `json:"warehouse_id"`
	SKUID          string  `json:"sku_id"`
	Algorithm      string  `json:"algorithm"`
	MAPE           float64 `json:"mape"`
	RMSE           float64 `json:"rmse"`
	MAE            float64 `json:"mae"`
	EvaluationDate string  `json:"evaluation_date"`
}

// ==================== GET /api/v1/forecasts ====================

func (h *ForecastHandler) QueryForecasts(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()
	warehouseID, ok := parseForecastID(w, r, "warehouse_id", requestID)
	if !ok {
		return
	}
	skuID, ok := parseForecastID(w, r, "sku_id", requestID)
	if !ok {
		return
	}

	daysStr := r.URL.Query().Get("days")
	limit, offset := util.ParsePagination(r)

	days := 30
	if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 90 {
		days = d
	}
	if limit > 100 {
		limit = 100
	}

	cacheKey := "forecasts:" + warehouseID.String() + ":" + skuID.String() + ":" + strconv.Itoa(days)
	if h.redisClient != nil {
		cachedData, err := h.redisClient.Get(r.Context(), "cache:"+cacheKey).Result()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedData))
			return
		}
	}

	if h.forecastSvc == nil {
		util.WriteError(w, http.StatusInternalServerError, "Forecast service unavailable", "", requestID)
		return
	}

	forecasts, err := h.forecastSvc.GetForecast(r.Context(), warehouseID, skuID, days)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve forecasts", err.Error(), requestID)
		return
	}

	items := make([]forecastResponse, len(forecasts))
	for i, f := range forecasts {
		items[i] = forecastResponse{
			ID:              f.ID.String(),
			WarehouseID:     f.WarehouseID.String(),
			SKUID:           f.SKUID.String(),
			ForecastDate:    f.ForecastDate.Format("2006-01-02"),
			PredictedDemand: f.PredictedDemand,
			LowerBound:      f.LowerBound,
			UpperBound:      f.UpperBound,
			Confidence:      f.Confidence,
			Algorithm:       string(f.Algorithm),
			GeneratedAt:     f.GeneratedAt,
		}
	}
	data := map[string]interface{}{
		"forecasts": items,
		"total":     len(items),
		"limit":     limit,
		"offset":    offset,
	}

	if h.redisClient != nil {
		response := util.BuildResponse(http.StatusOK, "OK", data, requestID)
		go func() {
			jsonData, err := json.Marshal(response)
			if err == nil {
				_ = h.redisClient.Set(context.Background(), "cache:"+cacheKey, jsonData, 24*time.Hour).Err()
			}
		}()
	}

	util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== POST /api/v1/forecasts/generate ====================

func (h *ForecastHandler) GenerateForecasts(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()
	userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
	if !ok || userInfo.Role != "admin" {
		util.WriteError(w, http.StatusForbidden, "Forbidden", "Only admin can trigger forecast generation", requestID)
		return
	}

	if h.forecastSvc == nil {
		util.WriteError(w, http.StatusInternalServerError, "Forecast service unavailable", "", requestID)
		return
	}
	if err := h.forecastSvc.GenerateForecastBatch(r.Context()); err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to generate forecasts", err.Error(), requestID)
		return
	}

	data := map[string]interface{}{
		"job_id": uuid.New().String(),
		"status": "completed",
	}
	util.WriteResponse(w, http.StatusAccepted, "Forecast batch completed", data, requestID)
}

// ==================== GET /api/v1/forecasts/accuracy ====================

func (h *ForecastHandler) GetAccuracy(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()
	warehouseID, ok := parseForecastID(w, r, "warehouse_id", requestID)
	if !ok {
		return
	}
	skuID, ok := parseForecastID(w, r, "sku_id", requestID)
	if !ok {
		return
	}

	if h.forecastSvc == nil {
		util.WriteError(w, http.StatusInternalServerError, "Forecast service unavailable", "", requestID)
		return
	}
	accuracyMetrics, err := h.forecastSvc.GetAccuracy(r.Context(), warehouseID, skuID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve forecast accuracy", err.Error(), requestID)
		return
	}

	accuracy := make([]accuracyResponse, len(accuracyMetrics))
	for i, metric := range accuracyMetrics {
		accuracy[i] = accuracyResponse{
			WarehouseID:    metric.WarehouseID.String(),
			SKUID:          metric.SKUID.String(),
			Algorithm:      string(metric.Algorithm),
			MAPE:           metricValue(metric.MAPE),
			RMSE:           metricValue(metric.RMSE),
			MAE:            metricValue(metric.MAE),
			EvaluationDate: metric.EvaluationDate.Format("2006-01-02"),
		}
	}
	data := map[string]interface{}{"accuracy": accuracy}
	util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== HELPER ====================

func parseForecastID(w http.ResponseWriter, r *http.Request, name, requestID string) (uuid.UUID, bool) {
	value := r.URL.Query().Get(name)
	if value == "" {
		util.WriteError(w, http.StatusBadRequest, "Missing required parameter", name+" is required", requestID)
		return uuid.Nil, false
	}

	id, err := uuid.Parse(value)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid "+name, err.Error(), requestID)
		return uuid.Nil, false
	}
	return id, true
}

func metricValue(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}
