package timeseries

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExponentialSmoothingForecast(t *testing.T) {
	// Buat historical data dengan 30 hari
	historical := make([]int32, 30)
	for i := 0; i < 30; i++ {
		historical[i] = int32(100 + i%10 + (i/7)*2)
	}
	forecastDays := 30

	t.Run("normal forecast", func(t *testing.T) {
		result, err := ExponentialSmoothingForecast(historical, forecastDays)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, result.PredictedDemand, int32(0))
		assert.Greater(t, result.LowerBound, int32(0))
		assert.Greater(t, result.UpperBound, result.LowerBound)
		assert.Equal(t, 0.95, result.Confidence)
		assert.GreaterOrEqual(t, result.MAPE, 0.0)
	})

	t.Run("insufficient data", func(t *testing.T) {
		small := []int32{100, 105, 110, 115, 120, 125, 130, 135, 140, 145}
		result, err := ExponentialSmoothingForecast(small, 30)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "insufficient historical data")
	})

	t.Run("forecastDays zero", func(t *testing.T) {
		result, err := ExponentialSmoothingForecast(historical, 0)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "forecast days must be positive")
	})
}

func TestCalculateMAPE(t *testing.T) {
	actual := []int32{100, 200, 300}
	predicted := []float64{110, 190, 310}
	mape := calculateMAPE(actual, predicted)
	expected := (10.0/100 + 10.0/200 + 10.0/300) / 3 * 100
	assert.InDelta(t, expected, mape, 0.01)

	actualZero := []int32{0, 200, 300}
	mapeZero := calculateMAPE(actualZero, predicted)
	expectedZero := (10.0/200 + 10.0/300) / 2 * 100
	assert.InDelta(t, expectedZero, mapeZero, 0.01)
}
