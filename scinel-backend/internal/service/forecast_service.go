package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/LittleMardiah/scinel/internal/model"
	"github.com/LittleMardiah/scinel/internal/repository"
	"github.com/LittleMardiah/scinel/pkg/timeseries"
)

// ServiceLogger is the shared logging contract used by all services.
// It is satisfied by *logger.Logger (slog-based) and by test mocks.
type ServiceLogger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

const (
	// forecastHistoryDays is how much historical demand we feed the models.
	forecastHistoryDays = 180
	// forecastHorizonDays is the number of daily points persisted per SKU.
	forecastHorizonDays = 30
)

// ForecastService handles forecast generation and management
type ForecastService struct {
	repo    repository.ForecastRepository
	invRepo repository.InventoryRepository
	txnRepo repository.TransactionRepository
	logger  ServiceLogger
}

// NewForecastService creates a new forecast service
func NewForecastService(
	repo repository.ForecastRepository,
	invRepo repository.InventoryRepository,
	txnRepo repository.TransactionRepository,
	logger ServiceLogger,
) *ForecastService {
	return &ForecastService{
		repo:    repo,
		invRepo: invRepo,
		txnRepo: txnRepo,
		logger:  logger,
	}
}

// GenerateForecastBatch generates forecasts for all SKU-warehouse combinations.
// Purpose: Batch job run daily (e.g. via cron). Runs sequentially to keep the
// implementation deterministic and testable.
func (fs *ForecastService) GenerateForecastBatch(ctx context.Context) error {
	start := time.Now()
	fs.logger.Info("Starting forecast batch job")

	inventory, total, err := fs.invRepo.GetAll(ctx, 10000, 0)
	if err != nil {
		fs.logger.Error("Failed to load inventory for forecast batch", "error", err.Error())
		return fmt.Errorf("failed to load inventory for forecast batch: %w", err)
	}
	if int64(len(inventory)) < total && total > 10000 {
		fs.logger.Warn("Forecast batch limited to first 10000 inventory rows", "total", total)
	}

	forecasted := 0
	errors := 0
	for _, inv := range inventory {
		if err := fs.GenerateForecast(ctx, inv.WarehouseID, inv.SKUID, ""); err != nil {
			errors++
			fs.logger.Error("Forecast generation failed",
				"warehouse_id", inv.WarehouseID,
				"sku_id", inv.SKUID,
				"error", err.Error())
			continue
		}
		forecasted++
	}

	fs.logger.Info("Forecast batch job complete", map[string]interface{}{
		"duration_ms":      time.Since(start).Milliseconds(),
		"total_processed":  len(inventory),
		"total_forecasted": forecasted,
		"total_errors":     errors,
	})
	return nil
}

// GenerateForecast generates a forecast for a single SKU-warehouse combination.
// algorithm empty string => auto-select (best MAPE between ETS and ARIMA).
func (fs *ForecastService) GenerateForecast(ctx context.Context, warehouseID, skuID uuid.UUID, algorithm model.Algorithm) error {
	// 1. Fetch historical demand
	historical, err := fs.fetchDailyDemand(ctx, warehouseID, skuID, forecastHistoryDays)
	if err != nil {
		return fmt.Errorf("failed to fetch historical demand: %w", err)
	}

	if len(historical) < 30 {
		fs.logger.Warn("Insufficient historical data for forecast",
			"warehouse_id", warehouseID, "sku_id", skuID, "days", len(historical))
		return nil
	}

	// 2. Run models
	ets, etsErr := timeseries.ExponentialSmoothingForecast(historical, forecastHorizonDays)
	arima, arimaErr := timeseries.AutoARIMAForecast(historical, forecastHorizonDays)

	var selectedAlgo model.Algorithm
	var pred int32
	var lower, upper int32
	var mape, rmse, mae, trend float64
	var seasonality *float64

	if algorithm == model.AlgorithmARIMA && arimaErr == nil {
		selectedAlgo = model.AlgorithmARIMA
		pred, lower, upper, mape, rmse, mae = arima.PredictedDemand, arima.LowerBound, arima.UpperBound, arima.MAPE, 0, 0
	} else if algorithm == model.AlgorithmETS && etsErr == nil {
		selectedAlgo = model.AlgorithmETS
		pred, lower, upper, mape, rmse, mae = ets.PredictedDemand, ets.LowerBound, ets.UpperBound, ets.MAPE, 0, 0
		trend = ets.Trend
		tmp := ets.SeasonalityFactor
		seasonality = &tmp
	} else {
		// Auto-select based on MAPE
		switch {
		case etsErr != nil && arimaErr != nil:
			fs.logger.Warn("Both forecast models failed",
				"warehouse_id", warehouseID, "sku_id", skuID,
				"ets", etsErr.Error(), "arima", arimaErr.Error())
			return fmt.Errorf("forecast models unavailable: %w", etsErr)
		case etsErr != nil:
			selectedAlgo = model.AlgorithmARIMA
			pred, lower, upper, mape = arima.PredictedDemand, arima.LowerBound, arima.UpperBound, arima.MAPE
		case arimaErr != nil:
			selectedAlgo = model.AlgorithmETS
			pred, lower, upper, mape = ets.PredictedDemand, ets.LowerBound, ets.UpperBound, ets.MAPE
			trend = ets.Trend
			tmp := ets.SeasonalityFactor
			seasonality = &tmp
		default:
			if ets.MAPE <= arima.MAPE {
				selectedAlgo = model.AlgorithmETS
				pred, lower, upper, mape = ets.PredictedDemand, ets.LowerBound, ets.UpperBound, ets.MAPE
				trend = ets.Trend
				tmp := ets.SeasonalityFactor
				seasonality = &tmp
			} else {
				selectedAlgo = model.AlgorithmARIMA
				pred, lower, upper, mape = arima.PredictedDemand, arima.LowerBound, arima.UpperBound, arima.MAPE
			}
		}
	}

	if mape > 0 {
		rmse, mae = estimateAccuracyMetrics(historical, pred)
	}

	// 3. Persist 30 daily forecast rows
	skuCode := ""
	if inv, err := fs.invRepo.GetBySKU(ctx, warehouseID, skuID); err == nil && inv != nil {
		skuCode = inv.SKUCode
	}

	now := time.Now().UTC()
	for i := 1; i <= forecastHorizonDays; i++ {
		dayAdjust := float64((i-1)%7) * 0.01 // mild weekly pattern within bounds
		dailyPred := int32(math.Round(float64(pred) * (1 + dayAdjust)))
		dailyLower := int32(math.Min(float64(dailyPred), float64(lower)))
		dailyUpper := int32(math.Max(float64(dailyPred), float64(upper)))

		f := &model.Forecast{
			ID:                uuid.New(),
			WarehouseID:       warehouseID,
			SKUID:             skuID,
			SKUCode:           skuCode,
			ForecastDate:      now.AddDate(0, 0, i).Truncate(24 * time.Hour),
			PredictedDemand:   dailyPred,
			LowerBound:        dailyLower,
			UpperBound:        dailyUpper,
			Confidence:        0.95,
			Algorithm:         selectedAlgo,
			Trend:             &trend,
			SeasonalityFactor: seasonality,
			GeneratedAt:       now,
			CreatedAt:         now,
		}
		if _, err := fs.repo.Create(ctx, f); err != nil {
			return fmt.Errorf("failed to persist forecast: %w", err)
		}
	}

	// 4. Persist accuracy metrics
	accuracy := &model.ForecastAccuracy{
		ID:             uuid.New(),
		WarehouseID:    warehouseID,
		SKUID:          skuID,
		Algorithm:      selectedAlgo,
		MAPE:           &mape,
		RMSE:           &rmse,
		MAE:            &mae,
		EvaluationDate: now,
		CreatedAt:      now,
	}
	if err := fs.repo.UpdateAccuracy(ctx, accuracy); err != nil {
		fs.logger.Warn("Failed to persist forecast accuracy",
			"warehouse_id", warehouseID, "sku_id", skuID, "error", err.Error())
	}

	fs.logger.Info("Forecast generated",
		"warehouse_id", warehouseID, "sku_id", skuID,
		"algorithm", selectedAlgo, "mape", mape, "predicted", pred)
	return nil
}

// GetForecast retrieves forecasts for a SKU-warehouse combination over `days`.
func (fs *ForecastService) GetForecast(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]model.Forecast, error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	return fs.repo.GetByWarehouseAndSKU(ctx, warehouseID, skuID, days)
}

// GetAccuracy retrieves the most recent accuracy metrics for both algorithms.
func (fs *ForecastService) GetAccuracy(ctx context.Context, warehouseID, skuID uuid.UUID) ([]model.ForecastAccuracy, error) {
	accs := make([]model.ForecastAccuracy, 0, 2)
	for _, algo := range []model.Algorithm{model.AlgorithmETS, model.AlgorithmARIMA} {
		acc, err := fs.repo.GetAccuracy(ctx, warehouseID, skuID, algo)
		if err != nil {
			return nil, fmt.Errorf("failed to get accuracy for %s: %w", algo, err)
		}
		if acc != nil {
			accs = append(accs, *acc)
		}
	}
	return accs, nil
}

// fetchDailyDemand aggregates transaction history into a daily demand series,
// sorted by date ascending.
func (fs *ForecastService) fetchDailyDemand(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]int32, error) {
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
	txn, _, err := fs.txnRepo.Query(ctx, filter, 10000, 0)
	if err != nil {
		return nil, err
	}
	return aggregateDailyDemand(txn), nil
}

// aggregateDailyDemand converts transaction history to a daily demand array,
// grouped by day (summing "out" quantities) and ordered ascending by date.
func aggregateDailyDemand(transactions []model.Transaction) []int32 {
	if len(transactions) == 0 {
		return []int32{}
	}

	daily := make(map[time.Time]int32)
	for _, txn := range transactions {
		if txn.Type != model.TransactionOut {
			continue
		}
		day := txn.CreatedAt.Truncate(24 * time.Hour)
		daily[day] += txn.Quantity
	}

	keys := make([]time.Time, 0, len(daily))
	for k := range daily {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })

	demand := make([]int32, len(keys))
	for i, k := range keys {
		demand[i] = daily[k]
	}
	return demand
}

// estimateAccuracyMetrics derives RMSE/MAE from the historical series and the
// predicted level. Used to populate the ForecastAccuracy record when the
// timeseries models do not expose RMSE/MAE directly.
func estimateAccuracyMetrics(historical []int32, predicted int32) (float64, float64) {
	if len(historical) < 2 {
		return 0, 0
	}
	var sumErr, sumSq float64
	for _, v := range historical {
		diff := float64(v) - float64(predicted)
		sumErr += math.Abs(diff)
		sumSq += diff * diff
	}
	mae := sumErr / float64(len(historical))
	rmse := math.Sqrt(sumSq / float64(len(historical)))
	return rmse, mae
}
