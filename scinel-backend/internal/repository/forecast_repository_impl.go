package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/LittleMardiah/scinel/internal/database"
	"github.com/LittleMardiah/scinel/internal/model"
)

type forecastRepository struct {
	db *database.DB
}

func NewForecastRepository(db *database.DB) ForecastRepository {
	return &forecastRepository{db: db}
}

// forecastColumns is the shared SELECT column list for forecasts + JOIN sku_code.
const forecastColumns = `
	f.id, f.warehouse_id, f.sku_id,
	s.sku_code,
	f.forecast_date, f.predicted_demand, f.lower_bound, f.upper_bound,
	f.confidence, f.algorithm, f.seasonality_factor,
	f.generated_at, f.created_at
`

const forecastSelect = `
	SELECT ` + forecastColumns + `
	FROM forecasts f
	JOIN skus s ON f.sku_id = s.id
`

// scanForecast reads one forecast row (trend column omitted due to varchar vs *float64 model mismatch).
func scanForecast(row interface{ Scan(dest ...interface{}) error }) (model.Forecast, error) {
	var f model.Forecast
	var seasonalityFactor *float64
	err := row.Scan(
		&f.ID, &f.WarehouseID, &f.SKUID,
		&f.SKUCode,
		&f.ForecastDate, &f.PredictedDemand, &f.LowerBound, &f.UpperBound,
		&f.Confidence, &f.Algorithm, &seasonalityFactor,
		&f.GeneratedAt, &f.CreatedAt,
	)
	if err != nil {
		return f, err
	}
	f.SeasonalityFactor = seasonalityFactor
	return f, nil
}

func (r *forecastRepository) Create(ctx context.Context, forecast *model.Forecast) (*model.Forecast, error) {
	// Note: the `trend` column is VARCHAR(50) but the model field is *float64.
	// We intentionally omit it from INSERT to avoid a type mismatch error.
	// SeasonalityFactor (NUMERIC) maps correctly to *float64.
	query := `
		INSERT INTO forecasts (id, warehouse_id, sku_id, forecast_date, predicted_demand,
		                       lower_bound, upper_bound, confidence, algorithm,
		                       seasonality_factor, generated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, warehouse_id, sku_id,
		          (SELECT sku_code FROM skus WHERE skus.id = forecasts.sku_id) AS sku_code,
		          forecast_date, predicted_demand, lower_bound, upper_bound,
		          confidence, algorithm, seasonality_factor,
		          generated_at, created_at
	`
	created, err := scanForecast(r.db.QueryRowContext(ctx, query,
		forecast.ID, forecast.WarehouseID, forecast.SKUID,
		forecast.ForecastDate, forecast.PredictedDemand,
		forecast.LowerBound, forecast.UpperBound,
		forecast.Confidence, forecast.Algorithm,
		forecast.SeasonalityFactor, forecast.GeneratedAt, forecast.CreatedAt,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create forecast: %w", err)
	}
	return &created, nil
}

func (r *forecastRepository) GetByWarehouseAndSKU(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]model.Forecast, error) {
	query := forecastSelect + `
		WHERE f.warehouse_id = $1
		  AND f.sku_id = $2
		  AND f.forecast_date >= CURRENT_DATE
		  AND f.forecast_date <= CURRENT_DATE + $3::int
		ORDER BY f.forecast_date ASC
	`
	rows, err := r.db.QueryContext(ctx, query, warehouseID, skuID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query forecasts by warehouse+sku: %w", err)
	}
	defer rows.Close()

	var items []model.Forecast
	for rows.Next() {
		f, err := scanForecast(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan forecast: %w", err)
		}
		items = append(items, f)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

func (r *forecastRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Forecast, int64, error) {
	query := forecastSelect + `
		ORDER BY f.forecast_date DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query all forecasts: %w", err)
	}
	defer rows.Close()

	var items []model.Forecast
	for rows.Next() {
		f, err := scanForecast(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan forecast: %w", err)
		}
		items = append(items, f)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	var total int64
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM forecasts`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count forecasts: %w", err)
	}
	return items, total, nil
}

func (r *forecastRepository) GetAccuracy(ctx context.Context, warehouseID, skuID uuid.UUID, algorithm model.Algorithm) (*model.ForecastAccuracy, error) {
	query := `
		SELECT id, warehouse_id, sku_id, algorithm,
		       mape, rmse, mae, evaluation_date, created_at
		FROM forecast_accuracy
		WHERE warehouse_id = $1
		  AND sku_id = $2
		  AND algorithm = $3
		ORDER BY evaluation_date DESC
		LIMIT 1
	`
	var acc model.ForecastAccuracy
	err := r.db.QueryRowContext(ctx, query, warehouseID, skuID, string(algorithm)).Scan(
		&acc.ID, &acc.WarehouseID, &acc.SKUID, &acc.Algorithm,
		&acc.MAPE, &acc.RMSE, &acc.MAE, &acc.EvaluationDate, &acc.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get forecast accuracy: %w", err)
	}
	return &acc, nil
}

func (r *forecastRepository) UpdateAccuracy(ctx context.Context, accuracy *model.ForecastAccuracy) error {
	// Insert a new accuracy record (each evaluation creates a snapshot row).
	// The forecast_accuracy table has no unique constraint, so we always INSERT.
	query := `
		INSERT INTO forecast_accuracy (id, warehouse_id, sku_id, algorithm, mape, rmse, mae, evaluation_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		accuracy.ID, accuracy.WarehouseID, accuracy.SKUID,
		string(accuracy.Algorithm),
		accuracy.MAPE, accuracy.RMSE, accuracy.MAE,
		accuracy.EvaluationDate, accuracy.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert forecast accuracy: %w", err)
	}
	return nil
}
