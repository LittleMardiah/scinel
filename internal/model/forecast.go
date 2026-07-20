package model

import (
    "time"

    "github.com/google/uuid"
)

// Forecast represents a demand prediction for a SKU in a warehouse
type Forecast struct {
    ID                uuid.UUID  `db:"id"`
    WarehouseID       uuid.UUID  `db:"warehouse_id"`
    SKUID             uuid.UUID  `db:"sku_id"`
    SKUCode           string     `db:"sku_code"`
    ForecastDate      time.Time  `db:"forecast_date"`
    PredictedDemand   int32      `db:"predicted_demand"`
    LowerBound        int32      `db:"lower_bound"`          // 95% confidence interval
    UpperBound        int32      `db:"upper_bound"`          // 95% confidence interval
    Confidence        float64    `db:"confidence"`           // 0.95 for 95%
    Algorithm         Algorithm  `db:"algorithm"`            // exponential_smoothing, arima
    Trend             *float64   `db:"trend"`                // Trend coefficient
    SeasonalityFactor *float64   `db:"seasonality_factor"`   // Seasonal adjustment
    GeneratedAt       time.Time  `db:"generated_at"`         // When forecast was created
    CreatedAt         time.Time  `db:"created_at"`
}

// ForecastAccuracy tracks algorithm performance over time
type ForecastAccuracy struct {
    ID             uuid.UUID `db:"id"`
    WarehouseID    uuid.UUID `db:"warehouse_id"`
    SKUID          uuid.UUID `db:"sku_id"`
    Algorithm      Algorithm `db:"algorithm"`
    MAPE           *float64  `db:"mape"`        // Mean Absolute Percentage Error
    RMSE           *float64  `db:"rmse"`        // Root Mean Square Error
    MAE            *float64  `db:"mae"`         // Mean Absolute Error
    EvaluationDate time.Time `db:"evaluation_date"`
    CreatedAt      time.Time `db:"created_at"`
}

// Algorithm type for forecasting
type Algorithm string

const (
    AlgorithmETS   Algorithm = "exponential_smoothing"
    AlgorithmARIMA Algorithm = "arima"
)

// Trend classification
type TrendType string

const (
    TrendFlat       TrendType = "flat"
    TrendIncreasing TrendType = "increasing"
    TrendDecreasing TrendType = "decreasing"
)