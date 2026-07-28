package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryStatus(t *testing.T) {
	assert.Equal(t, DeliveryStatus("pending"), DeliveryPending)
	assert.Equal(t, DeliveryStatus("success"), DeliverySuccess)
}
