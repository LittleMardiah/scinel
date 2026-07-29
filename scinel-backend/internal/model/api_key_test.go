package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAPIKey(t *testing.T) {
	key := &APIKey{Permissions: []string{"inventory:read"}}
	assert.Len(t, key.Permissions, 1)
}
