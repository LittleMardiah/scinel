package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.level)
			assert.NotNil(t, logger)
			assert.NotNil(t, logger.Logger)

			// Test that logging doesn't panic
			logger.Info("test message", "key", "value")
			logger.Error("test error", "error", "something")
			logger.Warn("test warning")
		})
	}
}
