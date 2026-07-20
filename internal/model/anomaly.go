package model

import (
    "time"

    "github.com/google/uuid"
)

// Anomaly represents an anomaly detection event
type Anomaly struct {
    ID                uuid.UUID     `db:"id"`
    WarehouseID       uuid.UUID     `db:"warehouse_id"`
    SKUID             uuid.UUID     `db:"sku_id"`
    SKUCode           string        `db:"sku_code"`
    Type              string        `db:"type"` // e.g., "unusual_stock_drop", "spike"
    Severity          Severity      `db:"severity"`
    Description       string        `db:"description"`
    HistoricalAverage float64       `db:"historical_average"`
    DetectedValue     float64       `db:"detected_value"`
    ZScore            float64       `db:"z_score"`
    Status            AnomalyStatus `db:"status"`
    ResolutionNotes   *string       `db:"resolution_notes"`
    ResolvedBy        *uuid.UUID    `db:"resolved_by"`
    ResolvedAt        *time.Time    `db:"resolved_at"`
    DetectedAt        time.Time     `db:"detected_at"`
    CreatedAt         time.Time     `db:"created_at"`
    UpdatedAt         time.Time     `db:"updated_at"`
}

// Severity levels
type Severity string

const (
    SeverityLow      Severity = "low"      // Z-score 2.0-3.0
    SeverityMedium   Severity = "medium"   // Z-score 3.0-4.0
    SeverityHigh     Severity = "high"     // Z-score 4.0-5.0
    SeverityCritical Severity = "critical" // Z-score >5.0
)

// AnomalyStatus
type AnomalyStatus string

const (
    StatusOpen         AnomalyStatus = "open"
    StatusAcknowledged AnomalyStatus = "acknowledged"
    StatusResolved     AnomalyStatus = "resolved"
)