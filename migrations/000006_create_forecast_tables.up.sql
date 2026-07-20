CREATE TABLE IF NOT EXISTS forecasts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    sku_id UUID NOT NULL REFERENCES skus(id) ON DELETE CASCADE,
    forecast_date DATE NOT NULL,
    predicted_demand INT NOT NULL,
    lower_bound INT,
    upper_bound INT,
    confidence DECIMAL(3, 2),
    algorithm VARCHAR(50) NOT NULL CHECK (algorithm IN ('exponential_smoothing', 'arima')),
    trend VARCHAR(50),
    seasonality_factor DECIMAL(5, 2),
    generated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forecasts_warehouse_sku ON forecasts(warehouse_id, sku_id);
CREATE INDEX idx_forecasts_date ON forecasts(forecast_date);
CREATE INDEX idx_forecasts_generated_at ON forecasts(generated_at DESC);

CREATE TABLE IF NOT EXISTS forecast_accuracy (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    sku_id UUID NOT NULL REFERENCES skus(id) ON DELETE CASCADE,
    algorithm VARCHAR(50) NOT NULL,
    mape DECIMAL(5, 2),
    rmse DECIMAL(10, 2),
    mae DECIMAL(10, 2),
    evaluation_date DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forecast_accuracy_warehouse_sku ON forecast_accuracy(warehouse_id, sku_id);
CREATE INDEX idx_forecast_accuracy_algorithm ON forecast_accuracy(algorithm);