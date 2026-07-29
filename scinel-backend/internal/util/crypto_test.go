package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashSHA256(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"test", "test", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashSHA256(tt.input)
			assert.Equal(t, tt.want, hash)
		})
	}
}

func TestHashSHA256_LongInput(t *testing.T) {
	input := strings.Repeat("a", 1000)
	hash := HashSHA256(input)
	assert.Equal(t, 64, len(hash))
	assert.NotEmpty(t, hash)
}
