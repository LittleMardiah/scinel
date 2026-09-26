package service

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
    "github.com/LittleMardiah/scinel/internal/repository"
    "github.com/LittleMardiah/scinel/pkg/timeseries"
)

// WebhookPublisherInterface publishes domain events asynchronously.
// Implemented by *WebhookPublisher.
type WebhookPublisherInterface interface {
	PublishAsync(ctx context.Context, eventType string, data interface{})
}

const (
	// anomalyHistoryDays is the look-back window for historical stock levels.
	anomalyHistoryDays = 90
	// anomalyMinDailyPoints guards against detecting on too little data.
	anomalyMinDailyPoints = 5
)

// AnomalyService detects anomalies in inventory
type AnomalyService struct {
    repo      repository.AnomalyRepository
    invRepo   repository.InventoryRepository
    txnRepo   repository.TransactionRepository
    logger    ServiceLogger
    publisher WebhookPublisherInterface
    auditSvc  *AuditService
}

// NewAnomalyService creates a new anomaly service.
//
// auditSvc is variadic so existing call sites keep compiling; with no
// AuditService the ANOMALY_DETECTED audit entry is simply not written.
func NewAnomalyService(
    repo repository.AnomalyRepository,
    invRepo repository.InventoryRepository,
    txnRepo repository.TransactionRepository,
    logger ServiceLogger,
    publisher WebhookPublisherInterface,
    auditSvc ...*AuditService,
) *AnomalyService {
    as := &AnomalyService{
        repo:      repo,
        invRepo:   invRepo,
        txnRepo:   txnRepo,
        logger:    logger,
        publisher: publisher,
    }
    for _, a := range auditSvc {
        if a != nil {
            as.auditSvc = a
        }
    }
    return as
}

// DetectAnomalies runs anomaly detection after an inventory transaction.
// Triggered: after every transaction (async via background goroutine).
func (as *AnomalyService) DetectAnomalies(ctx context.Context, warehouseID, skuID uuid.UUID) error {
	// 1. Get current inventory level
	inv, err := as.invRepo.GetBySKU(ctx, warehouseID, skuID)
	if err != nil {
		as.logger.Error("Failed to get inventory for anomaly detection",
			"warehouse_id", warehouseID, "sku_id", skuID, "error", err.Error())
		return fmt.Errorf("failed to get inventory for anomaly detection: %w", err)
	}
	if inv == nil {
		as.logger.Warn("No inventory record for anomaly detection",
			"warehouse_id", warehouseID, "sku_id", skuID)
		return nil
	}

	// 2. Fetch historical daily changes (90 days)
	dailyChanges, err := as.fetchDailyChanges(ctx, warehouseID, skuID, anomalyHistoryDays)
	if err != nil {
		as.logger.Error("Failed to fetch transaction history for anomaly detection",
			"warehouse_id", warehouseID, "sku_id", skuID, "error", err.Error())
		return fmt.Errorf("failed to fetch transaction history: %w", err)
	}
	if len(dailyChanges) < anomalyMinDailyPoints {
		as.logger.Warn("Insufficient daily data for anomaly detection",
			"warehouse_id", warehouseID, "sku_id", skuID, "days", len(dailyChanges))
		return nil
	}

	// 3. Run detection (Z-score primary, IQR as secondary)
	result := timeseries.DetectAnomalyZScore(dailyChanges, float64(inv.Quantity))
	iqrResult := timeseries.DetectAnomalyIQR(dailyChanges, float64(inv.Quantity))
	if !result.IsAnomaly && !iqrResult.IsAnomaly {
		return nil // No anomaly detected
	}

	// Prefer the more severe signal for reporting.
	pick := result
	if !pick.IsAnomaly {
		pick = iqrResult
	} else if iqrResult.IsAnomaly && severityRank(iqrResult.Severity) > severityRank(pick.Severity) {
		pick = iqrResult
	}

	// 4. Save anomaly record
	now := time.Now().UTC()
	anomaly := &model.Anomaly{
		ID:                uuid.New(),
		WarehouseID:       warehouseID,
		SKUID:             skuID,
		SKUCode:           inv.SKUCode,
		Type:              "unusual_stock_level",
		Severity:          model.Severity(pick.Severity),
		Description:       pick.Details,
		HistoricalAverage: pick.HistoricalMean,
		DetectedValue:     float64(inv.Quantity),
		ZScore:            pick.ZScore,
		Status:            model.StatusOpen,
		DetectedAt:        now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	saved, err := as.repo.Create(ctx, anomaly)
	if err != nil {
		as.logger.Error("Failed to save anomaly",
			"warehouse_id", warehouseID, "sku_id", skuID, "error", err.Error())
		return fmt.Errorf("failed to save anomaly: %w", err)
	}

    as.logger.Info("Anomaly detected and saved",
        "anomaly_id", saved.ID,
        "severity", saved.Severity,
        "sku_id", skuID,
        "warehouse_id", warehouseID)

    // 4b. Audit: ANOMALY_DETECTED. System-initiated (no acting user), so
    // user_id stays NULL. Non-fatal by construction.
    if raw, err := json.Marshal(map[string]interface{}{
        "anomaly_id":         saved.ID.String(),
        "warehouse_id":       warehouseID.String(),
        "sku_id":             skuID.String(),
        "sku_code":           saved.SKUCode,
        "type":               saved.Type,
        "severity":           string(saved.Severity),
        "z_score":            saved.ZScore,
        "detected_value":     saved.DetectedValue,
        "historical_average": saved.HistoricalAverage,
    }); err == nil {
        as.auditSvc.Log(ctx, &model.AuditLog{
            Action:       model.AuditActionAnomalyDetected,
            ResourceType: nullableString(model.AuditResourceAnomaly),
            ResourceID:   nullableString(saved.ID.String()),
            Details:      raw,
        })
    }

    // 5. Publish webhook event (async)
    as.publisher.PublishAsync(ctx, "anomaly.detected", saved)
    return nil
}

// fetchDailyChanges retrieves daily net "out" quantities for a SKU-warehouse
// over the requested window, ordered ascending by date.
func (as *AnomalyService) fetchDailyChanges(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]float64, error) {
	from := time.Now().UTC().AddDate(0, 0, -days)
	warehouseIDPtr := warehouseID
	skuIDPtr := skuID
	txnType := model.TransactionOut
	filter := model.TransactionFilter{
		WarehouseID: &warehouseIDPtr,
		SKUID:       &skuIDPtr,
		Type:        &txnType,
		DateFrom:    &from,
	}
	txn, _, err := as.txnRepo.Query(ctx, filter, 10000, 0)
	if err != nil {
		return nil, err
	}
	return aggregateDailyChanges(txn), nil
}

// aggregateDailyChanges converts transaction history to daily "out" quantities.
func aggregateDailyChanges(transactions []model.Transaction) []float64 {
	if len(transactions) == 0 {
		return []float64{}
	}

	daily := make(map[time.Time]float64)
	for _, txn := range transactions {
		if txn.Type != model.TransactionOut {
			continue
		}
		day := txn.CreatedAt.Truncate(24 * time.Hour)
		daily[day] += float64(txn.Quantity)
	}

	out := make([]float64, 0, len(daily))
	for _, v := range daily {
		out = append(out, v)
	}
	return out
}

func severityRank(s string) int {
	switch model.Severity(s) {
	case model.SeverityCritical:
		return 4
	case model.SeverityHigh:
		return 3
	case model.SeverityMedium:
		return 2
	default:
		return 1
	}
}
