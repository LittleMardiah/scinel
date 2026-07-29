package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
)

// ForecastRepository defines forecast data access operations
type ForecastRepository interface {
    // Create inserts forecast record
    Create(ctx context.Context, forecast *model.Forecast) (*model.Forecast, error)

    // GetByWarehouseAndSKU retrieves forecasts for a specific warehouse-SKU combination
    GetByWarehouseAndSKU(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]model.Forecast, error)

    // GetAll retrieves all recent forecasts
    GetAll(ctx context.Context, limit, offset int) ([]model.Forecast, int64, error)

    // GetAccuracy retrieves forecast accuracy metrics
    GetAccuracy(ctx context.Context, warehouseID, skuID uuid.UUID, algorithm model.Algorithm) (*model.ForecastAccuracy, error)

    // UpdateAccuracy records new accuracy metrics
    UpdateAccuracy(ctx context.Context, accuracy *model.ForecastAccuracy) error
}