package repository

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/database"
    "github.com/LittleMardiah/scinel/internal/model"
)

// forecastRepository implements ForecastRepository interface
type forecastRepository struct {
    db *database.DB
}

// NewForecastRepository creates a new forecast repository instance
func NewForecastRepository(db *database.DB) ForecastRepository {
    return &forecastRepository{db: db}
}

// --- PLACEHOLDER IMPLEMENTATIONS ---
// All methods return "not implemented" errors for now.
// Will be implemented in ROADMAP_2 Week 3.3

func (r *forecastRepository) Create(ctx context.Context, forecast *model.Forecast) (*model.Forecast, error) {
    return nil, errors.New("Create forecast not implemented yet")
}

func (r *forecastRepository) GetByWarehouseAndSKU(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]model.Forecast, error) {
    return nil, errors.New("GetByWarehouseAndSKU not implemented yet")
}

func (r *forecastRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Forecast, int64, error) {
    return nil, 0, errors.New("GetAll forecast not implemented yet")
}

func (r *forecastRepository) GetAccuracy(ctx context.Context, warehouseID, skuID uuid.UUID, algorithm model.Algorithm) (*model.ForecastAccuracy, error) {
    return nil, errors.New("GetAccuracy not implemented yet")
}

func (r *forecastRepository) UpdateAccuracy(ctx context.Context, accuracy *model.ForecastAccuracy) error {
    return errors.New("UpdateAccuracy not implemented yet")
}