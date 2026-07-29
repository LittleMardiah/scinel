package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAlgorithmConstants(t *testing.T) {
	assert.Equal(t, Algorithm("exponential_smoothing"), AlgorithmETS)
	assert.Equal(t, Algorithm("arima"), AlgorithmARIMA)
}
