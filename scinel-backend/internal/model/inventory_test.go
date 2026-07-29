package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInventoryGetStatus(t *testing.T) {
	tests := []struct {
		qty, reorder int32
		expected     InventoryStatus
	}{
		{-5, 100, StatusCritical},
		{0, 100, StatusCritical},
		{50, 100, StatusLowStock},
		{100, 100, StatusLowStock},
		{150, 100, StatusNormal},
	}
	for _, tt := range tests {
		inv := &Inventory{Quantity: tt.qty, ReorderPoint: tt.reorder}
		assert.Equal(t, tt.expected, inv.GetStatus())
	}
}
