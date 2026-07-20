package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/middleware"
    "github.com/LittleMardiah/scinel/internal/model"
    "github.com/LittleMardiah/scinel/internal/repository"
    "github.com/LittleMardiah/scinel/internal/service"
    "github.com/LittleMardiah/scinel/internal/util"
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
    warehouseIDStr := r.URL.Query().Get("warehouse_id")
    skuIDStr := r.URL.Query().Get("sku_id")
    daysStr := r.URL.Query().Get("days")
    limit, offset := util.ParsePagination(r)

    days := 30
    if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 90 {
        days = d
    }
    if limit > 100 {
        limit = 100
    }

    cacheKey := "forecasts"
    if warehouseIDStr != "" {
        cacheKey += ":" + warehouseIDStr
    }
    if skuIDStr != "" {
        cacheKey += ":" + skuIDStr
    }
    cacheKey += ":" + strconv.Itoa(days)

    cachedData, err := h.redisClient.Get(r.Context(), "cache:"+cacheKey).Result()
    if err == nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cachedData))
        return
    }

    forecasts := generatePlaceholderForecasts(warehouseIDStr, skuIDStr, days)
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
    response := util.BuildResponse(http.StatusOK, "OK", data, requestID)

    go func() {
        jsonData, _ := json.Marshal(response)
        _ = h.redisClient.Set(context.Background(), "cache:"+cacheKey, jsonData, 24*time.Hour).Err()
    }()

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

    go func() {
        time.Sleep(1 * time.Second)
    }()

    data := map[string]interface{}{
        "job_id":                     uuid.New().String(),
        "status":                     "processing",
        "estimated_duration_seconds": 5,
    }
    util.WriteResponse(w, http.StatusAccepted, "Forecast batch started", data, requestID)
}

// ==================== GET /api/v1/forecasts/accuracy ====================

func (h *ForecastHandler) GetAccuracy(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()
    accuracy := []accuracyResponse{
        {
            WarehouseID:    "wh-jakarta-001",
            SKUID:          "sku-nike-001",
            Algorithm:      "exponential_smoothing",
            MAPE:           12.5,
            RMSE:           45.3,
            MAE:            32.1,
            EvaluationDate: time.Now().UTC().Format("2006-01-02"),
        },
        {
            WarehouseID:    "wh-jakarta-001",
            SKUID:          "sku-nike-001",
            Algorithm:      "arima",
            MAPE:           14.2,
            RMSE:           48.7,
            MAE:            35.4,
            EvaluationDate: time.Now().UTC().Format("2006-01-02"),
        },
    }
    data := map[string]interface{}{"accuracy": accuracy}
    util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== HELPER ====================

func generatePlaceholderForecasts(warehouseIDStr, skuIDStr string, days int) []model.Forecast {
    warehouseID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
    if id, err := uuid.Parse(warehouseIDStr); err == nil {
        warehouseID = id
    }
    skuID, _ := uuid.Parse("660e8400-e29b-41d4-a716-446655440001")
    if id, err := uuid.Parse(skuIDStr); err == nil {
        skuID = id
    }
    forecasts := make([]model.Forecast, days)
    now := time.Now().UTC()
    for i := 0; i < days; i++ {
        demand := int32(100 + (i%7)*10 + (i%3)*5)
        forecasts[i] = model.Forecast{
            ID:              uuid.New(),
            WarehouseID:     warehouseID,
            SKUID:           skuID,
            ForecastDate:    now.AddDate(0, 0, i+1),
            PredictedDemand: demand,
            LowerBound:      demand - 20,
            UpperBound:      demand + 20,
            Confidence:      0.95,
            Algorithm:       model.AlgorithmETS,
            GeneratedAt:     now,
            CreatedAt:       now,
        }
    }
    return forecasts
}