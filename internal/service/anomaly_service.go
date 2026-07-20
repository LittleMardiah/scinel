package service

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
    "github.com/LittleMardiah/scinel/internal/repository"
    "github.com/LittleMardiah/scinel/pkg/timeseries"
)

// Logger interface for logging (placeholder)
type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
}

// WebhookPublisher interface for publishing events (placeholder)
type WebhookPublisher interface {
    PublishAsync(ctx context.Context, eventType string, data interface{})
}

// AnomalyService detects anomalies in inventory
type AnomalyService struct {
    repo      repository.AnomalyRepository
    invRepo   repository.InventoryRepository
    logger    Logger
    publisher WebhookPublisher
}

// NewAnomalyService creates a new anomaly service
func NewAnomalyService(
    repo repository.AnomalyRepository,
    invRepo repository.InventoryRepository,
    logger Logger,
    publisher WebhookPublisher,
) *AnomalyService {
    return &AnomalyService{
        repo:      repo,
        invRepo:   invRepo,
        logger:    logger,
        publisher: publisher,
    }
}

// DetectAnomalies runs anomaly detection after inventory transaction
// Triggered: After every inventory transaction (async via background job)
// Input: warehouse_id, sku_id
// Process: Get historical inventory, run Z-score detection, save if anomaly
func (as *AnomalyService) DetectAnomalies(ctx context.Context, warehouseID, skuID uuid.UUID) error {
    // 1. Get current inventory level
    inv, err := as.invRepo.GetBySKU(ctx, warehouseID, skuID)
    if err != nil || inv == nil {
        if err != nil {
            as.logger.Error("Failed to get inventory for anomaly detection", "error", err.Error())
        }
        return err
    }

    // 2. Get 90-day transaction history
    // TODO: Implement GetTransactionHistory in InventoryRepository
    // For now, placeholder - will be implemented in ROADMAP_2 Week 3.2
    transactions, err := as.getTransactionHistory(ctx, warehouseID, skuID, 90)
    if err != nil || len(transactions) < 10 {
        as.logger.Warn("Insufficient transaction history for anomaly detection", 
            "warehouse_id", warehouseID, 
            "sku_id", skuID,
            "count", len(transactions))
        return nil // Not enough data, skip
    }

    // 3. Convert to daily quantity changes
    dailyChanges := aggregateDailyChanges(transactions)
    if len(dailyChanges) < 5 {
        as.logger.Warn("Insufficient daily data for anomaly detection", 
            "warehouse_id", warehouseID, 
            "sku_id", skuID,
            "days", len(dailyChanges))
        return nil
    }

    // 4. Run anomaly detection
    result := timeseries.DetectAnomalyZScore(dailyChanges, float64(inv.Quantity))
    if !result.IsAnomaly {
        return nil // No anomaly detected
    }

    // 5. Save anomaly record
    now := time.Now().UTC()
    anomaly := &model.Anomaly{
        WarehouseID:       warehouseID,
        SKUID:             skuID,
        SKUCode:           inv.SKUCode,
        Type:              "unusual_stock_level",
        Severity:          model.Severity(result.Severity),
        Description:       result.Details,
        HistoricalAverage: result.HistoricalMean,
        DetectedValue:     float64(inv.Quantity),
        ZScore:            result.ZScore,
        Status:            model.StatusOpen,
        DetectedAt:        now,
        CreatedAt:         now,
        UpdatedAt:         now,
    }

    saved, err := as.repo.Create(ctx, anomaly)
    if err != nil {
        as.logger.Error("Failed to save anomaly", "error", err.Error())
        return err
    }

    as.logger.Info("Anomaly detected and saved", 
        "anomaly_id", saved.ID,
        "severity", saved.Severity,
        "sku_id", skuID,
        "warehouse_id", warehouseID)

    // 6. Publish webhook event
    as.publisher.PublishAsync(ctx, "anomaly.detected", saved)

    return nil
}

// getTransactionHistory is a placeholder - will be replaced with actual repository method
func (as *AnomalyService) getTransactionHistory(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]model.Transaction, error) {
    // TODO: Implement actual transaction history retrieval
    // This is a placeholder returning empty slice
    return []model.Transaction{}, nil
}

// aggregateDailyChanges converts transaction history to daily demand array
// Takes transactions, groups by day, sums "out" quantities
func aggregateDailyChanges(transactions []model.Transaction) []float64 {
    if len(transactions) == 0 {
        return []float64{}
    }

    // Group by day
    dailyMap := make(map[string]float64)
    for _, txn := range transactions {
        // Only consider "out" transactions (stock going out = demand)
        if txn.Type != model.TransactionOut {
            continue
        }
        day := txn.CreatedAt.Format("2006-01-02")
        dailyMap[day] += float64(txn.Quantity)
    }

    // Convert map to slice
    changes := make([]float64, 0, len(dailyMap))
    for _, val := range dailyMap {
        changes = append(changes, val)
    }

    // If no "out" transactions, use all transactions as fallback
    if len(changes) == 0 {
        for _, txn := range transactions {
            changes = append(changes, float64(txn.Quantity))
        }
    }

    return changes
}