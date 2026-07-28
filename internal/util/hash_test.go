package util

import (
"testing"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
tests := []struct {
name     string
password string
cost     int
wantErr  bool
}{
{"default cost", "SecurePass123!", 0, false},
{"cost 10", "Test123!", 10, false},
{"cost 12", "AnotherPass456", 12, false},
{"empty password", "", 12, false},
}
for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
hash, err := HashPassword(tt.password, tt.cost)
if tt.wantErr {
assert.Error(t, err)
return
}
require.NoError(t, err)
assert.NotEmpty(t, hash)
err = VerifyPassword(tt.password, hash)
assert.NoError(t, err)
err = VerifyPassword("wrong", hash)
assert.Error(t, err)
})
}
}

func TestVerifyPassword(t *testing.T) {
password := "MySecretPass123!"
hash, err := HashPassword(password, 0)
require.NoError(t, err)

err = VerifyPassword(password, hash)
assert.NoError(t, err)

err = VerifyPassword("wrong", hash)
assert.Error(t, err)

err = VerifyPassword("", hash)
assert.Error(t, err)

err = VerifyPassword("pass", "invalidhash")
assert.Error(t, err)
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
err := VerifyPassword("pass", "invalid")
assert.Error(t, err)
}

func TestHashPassword_InvalidCost(t *testing.T) {
	// Bcrypt mengizinkan cost 4-31, cost di luar range akan di-clamp
	// Kita test bahwa fungsi tetap berjalan tanpa panic
	hash, err := HashPassword("pass", 3)
	if err != nil {
		t.Logf("cost 3 returned error: %v (bcrypt may clamp)", err)
	}
	if hash == "" && err == nil {
		t.Error("expected non-empty hash for cost 3")
	}
	
	hash2, err2 := HashPassword("pass", 32)
	if err2 != nil {
		t.Logf("cost 32 returned error: %v (bcrypt may clamp)", err2)
	}
	if hash2 == "" && err2 == nil {
		t.Error("expected non-empty hash for cost 32")
	}
}
