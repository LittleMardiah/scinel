package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSeverityConstants(t *testing.T) {
	assert.Equal(t, Severity("low"), SeverityLow)
	assert.Equal(t, Severity("critical"), SeverityCritical)
}
