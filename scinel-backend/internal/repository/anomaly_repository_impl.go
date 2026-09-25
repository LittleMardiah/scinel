package repository

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/database"
    "github.com/LittleMardiah/scinel/internal/model"
)

// anomalyRepository implements AnomalyRepository interface
type anomalyRepository struct {
    db *database.DB
}

// NewAnomalyRepository creates a new anomaly repository instance
func NewAnomalyRepository(db *database.DB) AnomalyRepository {
    return &anomalyRepository{db: db}
}

// --- IMPLEMENTATIONS ---

// GetAll retrieves all anomalies with pagination
func (r *anomalyRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Anomaly, int64, error) {
    query := `
        SELECT 
            a.id, a.warehouse_id, a.sku_id, 
            s.sku_code,
            a.type, a.severity, a.description,
            a.historical_average, a.detected_value, a.z_score,
            a.status, a.resolution_notes, a.resolved_by, a.resolved_at,
            a.detected_at, a.created_at, a.updated_at
        FROM anomalies a
        JOIN skus s ON a.sku_id = s.id
        ORDER BY a.detected_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query anomalies: %w", err)
    }
    defer rows.Close()

    var items []model.Anomaly
    for rows.Next() {
        var anom model.Anomaly
        var resolutionNotes *string
        var resolvedBy *uuid.UUID
        var resolvedAt *time.Time

        err := rows.Scan(
            &anom.ID, &anom.WarehouseID, &anom.SKUID,
            &anom.SKUCode,
            &anom.Type, &anom.Severity, &anom.Description,
            &anom.HistoricalAverage, &anom.DetectedValue, &anom.ZScore,
            &anom.Status, &resolutionNotes, &resolvedBy, &resolvedAt,
            &anom.DetectedAt, &anom.CreatedAt, &anom.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan anomaly: %w", err)
        }
        anom.ResolutionNotes = resolutionNotes
        anom.ResolvedBy = resolvedBy
        anom.ResolvedAt = resolvedAt
        items = append(items, anom)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    // Count total
    var total int64
    countQuery := `SELECT COUNT(*) FROM anomalies`
    err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count anomalies: %w", err)
    }

    return items, total, nil
}

// GetByID retrieves single anomaly by ID
func (r *anomalyRepository) GetByID(ctx context.Context, anomalyID uuid.UUID) (*model.Anomaly, error) {
    query := `
        SELECT 
            a.id, a.warehouse_id, a.sku_id, 
            s.sku_code,
            a.type, a.severity, a.description,
            a.historical_average, a.detected_value, a.z_score,
            a.status, a.resolution_notes, a.resolved_by, a.resolved_at,
            a.detected_at, a.created_at, a.updated_at
        FROM anomalies a
        JOIN skus s ON a.sku_id = s.id
        WHERE a.id = $1
    `

    var anom model.Anomaly
    var resolutionNotes *string
    var resolvedBy *uuid.UUID
    var resolvedAt *time.Time

    err := r.db.QueryRowContext(ctx, query, anomalyID).Scan(
        &anom.ID, &anom.WarehouseID, &anom.SKUID,
        &anom.SKUCode,
        &anom.Type, &anom.Severity, &anom.Description,
        &anom.HistoricalAverage, &anom.DetectedValue, &anom.ZScore,
        &anom.Status, &resolutionNotes, &resolvedBy, &resolvedAt,
        &anom.DetectedAt, &anom.CreatedAt, &anom.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get anomaly by ID: %w", err)
    }

    anom.ResolutionNotes = resolutionNotes
    anom.ResolvedBy = resolvedBy
    anom.ResolvedAt = resolvedAt

    return &anom, nil
}

// GetOpen retrieves all open/unresolved anomalies
func (r *anomalyRepository) GetOpen(ctx context.Context, limit, offset int) ([]model.Anomaly, int64, error) {
    query := `
        SELECT 
            a.id, a.warehouse_id, a.sku_id, 
            s.sku_code,
            a.type, a.severity, a.description,
            a.historical_average, a.detected_value, a.z_score,
            a.status, a.resolution_notes, a.resolved_by, a.resolved_at,
            a.detected_at, a.created_at, a.updated_at
        FROM anomalies a
        JOIN skus s ON a.sku_id = s.id
        WHERE a.status = 'open'
        ORDER BY a.detected_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query open anomalies: %w", err)
    }
    defer rows.Close()

    var items []model.Anomaly
    for rows.Next() {
        var anom model.Anomaly
        var resolutionNotes *string
        var resolvedBy *uuid.UUID
        var resolvedAt *time.Time

        err := rows.Scan(
            &anom.ID, &anom.WarehouseID, &anom.SKUID,
            &anom.SKUCode,
            &anom.Type, &anom.Severity, &anom.Description,
            &anom.HistoricalAverage, &anom.DetectedValue, &anom.ZScore,
            &anom.Status, &resolutionNotes, &resolvedBy, &resolvedAt,
            &anom.DetectedAt, &anom.CreatedAt, &anom.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan anomaly: %w", err)
        }
        anom.ResolutionNotes = resolutionNotes
        anom.ResolvedBy = resolvedBy
        anom.ResolvedAt = resolvedAt
        items = append(items, anom)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    var total int64
    countQuery := `SELECT COUNT(*) FROM anomalies WHERE status = 'open'`
    err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count open anomalies: %w", err)
    }

    return items, total, nil
}

// UpdateStatus updates anomaly status (acknowledge or resolve)
func (r *anomalyRepository) UpdateStatus(ctx context.Context, anomalyID uuid.UUID, status model.AnomalyStatus, notes string, resolvedBy uuid.UUID) error {
    query := `
        UPDATE anomalies 
        SET status = $1, 
            resolution_notes = $2,
            resolved_by = $3,
            resolved_at = CASE WHEN $1::varchar = 'resolved' THEN NOW() ELSE NULL END,
            updated_at = NOW()
        WHERE id = $4
    `
    result, err := r.db.ExecContext(ctx, query, status, notes, resolvedBy, anomalyID)
    if err != nil {
        return fmt.Errorf("failed to update anomaly status: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("anomaly not found")
    }
    return nil
}

// Create inserts new anomaly record
func (r *anomalyRepository) Create(ctx context.Context, anomaly *model.Anomaly) (*model.Anomaly, error) {
    query := `
        INSERT INTO anomalies (id, warehouse_id, sku_id, type, severity, description,
                               historical_average, detected_value, z_score, status,
                               resolution_notes, resolved_by, resolved_at,
                               detected_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        RETURNING
            id, warehouse_id, sku_id,
            (SELECT sku_code FROM skus WHERE skus.id = anomalies.sku_id) AS sku_code,
            type, severity, description,
            historical_average, detected_value, z_score,
            status, resolution_notes, resolved_by, resolved_at,
            detected_at, created_at, updated_at
    `
    var anom model.Anomaly
    var resolutionNotes *string
    var resolvedBy *uuid.UUID
    var resolvedAt *time.Time

    err := r.db.QueryRowContext(ctx, query,
        anomaly.ID, anomaly.WarehouseID, anomaly.SKUID,
        anomaly.Type, anomaly.Severity, anomaly.Description,
        anomaly.HistoricalAverage, anomaly.DetectedValue, anomaly.ZScore,
        anomaly.Status,
        anomaly.ResolutionNotes, anomaly.ResolvedBy, anomaly.ResolvedAt,
        anomaly.DetectedAt, anomaly.CreatedAt, anomaly.UpdatedAt,
    ).Scan(
        &anom.ID, &anom.WarehouseID, &anom.SKUID,
        &anom.SKUCode,
        &anom.Type, &anom.Severity, &anom.Description,
        &anom.HistoricalAverage, &anom.DetectedValue, &anom.ZScore,
        &anom.Status, &resolutionNotes, &resolvedBy, &resolvedAt,
        &anom.DetectedAt, &anom.CreatedAt, &anom.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create anomaly: %w", err)
    }
    anom.ResolutionNotes = resolutionNotes
    anom.ResolvedBy = resolvedBy
    anom.ResolvedAt = resolvedAt
    return &anom, nil
}

// GetBySKU retrieves anomalies for a specific SKU
func (r *anomalyRepository) GetBySKU(ctx context.Context, skuID uuid.UUID, limit, offset int) ([]model.Anomaly, int64, error) {
    query := `
        SELECT 
            a.id, a.warehouse_id, a.sku_id, 
            s.sku_code,
            a.type, a.severity, a.description,
            a.historical_average, a.detected_value, a.z_score,
            a.status, a.resolution_notes, a.resolved_by, a.resolved_at,
            a.detected_at, a.created_at, a.updated_at
        FROM anomalies a
        JOIN skus s ON a.sku_id = s.id
        WHERE a.sku_id = $1
        ORDER BY a.detected_at DESC
        LIMIT $2 OFFSET $3
    `

    rows, err := r.db.QueryContext(ctx, query, skuID, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query anomalies by SKU: %w", err)
    }
    defer rows.Close()

    var items []model.Anomaly
    for rows.Next() {
        var anom model.Anomaly
        var resolutionNotes *string
        var resolvedBy *uuid.UUID
        var resolvedAt *time.Time

        err := rows.Scan(
            &anom.ID, &anom.WarehouseID, &anom.SKUID,
            &anom.SKUCode,
            &anom.Type, &anom.Severity, &anom.Description,
            &anom.HistoricalAverage, &anom.DetectedValue, &anom.ZScore,
            &anom.Status, &resolutionNotes, &resolvedBy, &resolvedAt,
            &anom.DetectedAt, &anom.CreatedAt, &anom.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan anomaly by SKU: %w", err)
        }
        anom.ResolutionNotes = resolutionNotes
        anom.ResolvedBy = resolvedBy
        anom.ResolvedAt = resolvedAt
        items = append(items, anom)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    var total int64
    countQuery := `SELECT COUNT(*) FROM anomalies WHERE sku_id = $1`
    err = r.db.QueryRowContext(ctx, countQuery, skuID).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count anomalies by SKU: %w", err)
    }

    return items, total, nil
}