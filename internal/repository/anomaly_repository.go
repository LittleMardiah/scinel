package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
)

// AnomalyRepository defines anomaly data access operations
type AnomalyRepository interface {
    // Create inserts new anomaly record
    Create(ctx context.Context, anomaly *model.Anomaly) (*model.Anomaly, error)

    // GetByID retrieves single anomaly
    GetByID(ctx context.Context, anomalyID uuid.UUID) (*model.Anomaly, error)

    // GetAll retrieves anomalies with filters
    GetAll(ctx context.Context, limit, offset int) ([]model.Anomaly, int64, error)

    // GetOpen retrieves open/unresolved anomalies
    GetOpen(ctx context.Context, limit, offset int) ([]model.Anomaly, int64, error)

    // UpdateStatus updates anomaly status (acknowledge or resolve)
    UpdateStatus(ctx context.Context, anomalyID uuid.UUID, status model.AnomalyStatus, notes string, resolvedBy uuid.UUID) error

    // GetBySKU retrieves anomalies for a specific SKU
    GetBySKU(ctx context.Context, skuID uuid.UUID, limit, offset int) ([]model.Anomaly, int64, error)
}