package timeseries

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectAnomalyZScore(t *testing.T) {
	historical := []float64{100, 102, 98, 101, 99, 100, 103, 97, 100, 102}

	t.Run("normal value", func(t *testing.T) {
		result := DetectAnomalyZScore(historical, 100)
		require.NotNil(t, result)
		assert.False(t, result.IsAnomaly)
		assert.Equal(t, "low", result.Severity)
	})

	t.Run("extreme high value", func(t *testing.T) {
		result := DetectAnomalyZScore(historical, 500)
		require.NotNil(t, result)
		assert.True(t, result.IsAnomaly)
		assert.Equal(t, "critical", result.Severity)
		assert.Greater(t, result.ZScore, 5.0)
	})

	t.Run("extreme low value", func(t *testing.T) {
		result := DetectAnomalyZScore(historical, 0)
		require.NotNil(t, result)
		assert.True(t, result.IsAnomaly)
		assert.Equal(t, "critical", result.Severity)
		assert.Less(t, result.ZScore, -5.0)
	})

	t.Run("insufficient data", func(t *testing.T) {
		small := []float64{100, 101}
		result := DetectAnomalyZScore(small, 150)
		require.NotNil(t, result)
		assert.False(t, result.IsAnomaly)
		assert.Equal(t, "Insufficient historical data", result.Details)
	})
}

func TestDetectAnomalyIQR(t *testing.T) {
	historical := []float64{10, 12, 11, 13, 14, 10, 12, 11, 13, 100}

	t.Run("normal value", func(t *testing.T) {
		result := DetectAnomalyIQR(historical, 12)
		require.NotNil(t, result)
		assert.False(t, result.IsAnomaly)
	})

	t.Run("outlier value", func(t *testing.T) {
		result := DetectAnomalyIQR(historical, 100)
		require.NotNil(t, result)
		assert.True(t, result.IsAnomaly)
		assert.Equal(t, "critical", result.Severity)
	})

	t.Run("insufficient data", func(t *testing.T) {
		small := []float64{10, 11}
		result := DetectAnomalyIQR(small, 20)
		require.NotNil(t, result)
		assert.False(t, result.IsAnomaly)
		assert.Equal(t, "Insufficient historical data", result.Details)
	})
}
