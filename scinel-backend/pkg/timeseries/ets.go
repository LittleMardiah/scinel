package timeseries

import (
    "fmt"
    "math"
)

// ETSForecast represents Exponential Smoothing forecast result
type ETSForecast struct {
    PredictedDemand   int32   // Rounded forecast
    LowerBound        int32   // 95% confidence interval
    UpperBound        int32   // 95% confidence interval
    Confidence        float64 // 0.95
    Trend             float64 // Trend component
    SeasonalityFactor float64 // Seasonal adjustment (if applicable)
    MAPE              float64 // Mean Absolute Percentage Error on validation set
}

// ExponentialSmoothingForecast performs exponential smoothing on time series
// Input: Historical daily demand (180 days minimum recommended)
// Output: 30-day forecast with confidence intervals
// Reference: TDD Section 6.3.1
func ExponentialSmoothingForecast(historicalDemand []int32, forecastDays int) (*ETSForecast, error) {
	if forecastDays <= 0 {
		return nil, fmt.Errorf("forecast days must be positive")
	}
    if len(historicalDemand) < 30 {
        return nil, fmt.Errorf("insufficient historical data: %d days (minimum 30 required)", len(historicalDemand))
    }

    // Initialize smoothing parameters (alpha, beta, gamma)
    // These are tuned for retail demand patterns
    alpha := 0.3  // Level smoothing (0.1-0.5 for trend-heavy data)
    beta := 0.1   // Trend smoothing (0.01-0.2 for gradual trends)
    gamma := 0.05 // Seasonal smoothing (for multi-year seasonality)

    // Initialize level, trend, seasonal components
    level := float64(historicalDemand[0])
    trend := (float64(historicalDemand[1]) - float64(historicalDemand[0])) / 2.0

    // Seasonal factors (optional, for multi-year patterns)
    seasonalPeriod := 7 // Weekly seasonality (or 365 for yearly)
    seasonalFactors := make([]float64, seasonalPeriod)
    for i := 0; i < seasonalPeriod && i < len(historicalDemand); i++ {
        seasonalFactors[i] = 1.0
    }

    // Apply exponential smoothing
    for i := 1; i < len(historicalDemand); i++ {
        demand := float64(historicalDemand[i])

        // Update level
        newLevel := alpha*demand + (1-alpha)*(level+trend)

        // Update trend
        newTrend := beta*(newLevel-level) + (1-beta)*trend

        // Update seasonal (if applicable)
        if seasonalPeriod > 0 && i >= seasonalPeriod {
            idx := i % seasonalPeriod
            seasonalFactors[idx] = gamma*(demand/level) + (1-gamma)*seasonalFactors[idx]
        }

        level = newLevel
        trend = newTrend
    }

    // Generate forecast
    forecast := make([]float64, forecastDays)
    for h := 1; h <= forecastDays; h++ {
        fh := level + float64(h)*trend
        if seasonalPeriod > 0 {
            idx := (len(historicalDemand) + h - 1) % seasonalPeriod
            fh *= seasonalFactors[idx]
        }
        forecast[h-1] = fh
    }

    // Calculate standard error for confidence intervals (95% ≈ ±1.96 * SE)
    residuals := calculateResiduals(historicalDemand, level, trend, seasonalFactors, seasonalPeriod)
    stdError := calculateStdError(residuals)

    // Average forecast for forecastDays period
    avgForecast := 0.0
    for _, f := range forecast {
        avgForecast += f
    }
    avgForecast /= float64(forecastDays)

    // Calculate MAPE on recent 30-day validation window
    mape := calculateMAPE(historicalDemand[len(historicalDemand)-30:], forecast[:30])

    return &ETSForecast{
        PredictedDemand:   int32(math.Round(avgForecast)),
        LowerBound:        int32(math.Max(0, avgForecast-1.96*stdError)),
        UpperBound:        int32(avgForecast + 1.96*stdError),
        Confidence:        0.95,
        Trend:             trend,
        SeasonalityFactor: 1.0, // Simplified for MVP
        MAPE:              mape,
    }, nil
}

// ==================== HELPER FUNCTIONS ====================

func calculateResiduals(demand []int32, level, trend float64, seasonal []float64, period int) []float64 {
    residuals := make([]float64, 0, len(demand))
    for i, d := range demand {
        predicted := level + float64(i)*trend
        if period > 0 {
            idx := i % period
            predicted *= seasonal[idx]
        }
        residuals = append(residuals, float64(d)-predicted)
    }
    return residuals
}

func calculateStdError(residuals []float64) float64 {
    if len(residuals) < 2 {
        return 1.0
    }
    mean := 0.0
    for _, r := range residuals {
        mean += r
    }
    mean /= float64(len(residuals))

    variance := 0.0
    for _, r := range residuals {
        variance += (r - mean) * (r - mean)
    }
    variance /= float64(len(residuals) - 1)
    return math.Sqrt(variance)
}

func calculateMAPE(actual []int32, predicted []float64) float64 {
    if len(actual) != len(predicted) {
        return 0.0
    }

    sumAPE := 0.0
    count := 0
    for i, a := range actual {
        if a > 0 {
            ape := math.Abs(float64(a)-predicted[i]) / float64(a)
            sumAPE += ape
            count++
        }
    }

    if count == 0 {
        return 0.0
    }
    return (sumAPE / float64(count)) * 100.0 // Percentage
}