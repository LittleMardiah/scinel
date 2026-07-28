package util

import (
"testing"
"time"
"github.com/google/uuid"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

func TestJWTManager(t *testing.T) {
secretKey := "test-secret-key-1234567890"
accessExpiry := 1 * time.Hour
refreshExpiry := 24 * time.Hour
manager := NewJWTManager(secretKey, accessExpiry, refreshExpiry)
userID := uuid.New()
email := "test@example.com"
role := "admin"

t.Run("GenerateAccessToken", func(t *testing.T) {
token, err := manager.GenerateAccessToken(userID, email, role)
require.NoError(t, err)
assert.NotEmpty(t, token)
})
t.Run("VerifyAccessToken", func(t *testing.T) {
token, _ := manager.GenerateAccessToken(userID, email, role)
claims, err := manager.VerifyAccessToken(token)
require.NoError(t, err)
assert.Equal(t, userID, claims.UserID)
})
t.Run("InvalidToken", func(t *testing.T) {
_, err := manager.VerifyAccessToken("invalid")
assert.Error(t, err)
})
}

func TestExpiredToken(t *testing.T) {
secretKey := "test-secret"
manager := NewJWTManager(secretKey, -1*time.Hour, 24*time.Hour)
userID := uuid.New()
token, _ := manager.GenerateAccessToken(userID, "test@example.com", "admin")
_, err := manager.VerifyAccessToken(token)
assert.Error(t, err)
}
func TestRefreshToken(t *testing.T) {
	secretKey := "test-secret"
	manager := NewJWTManager(secretKey, 1*time.Hour, 24*time.Hour)
	userID := uuid.New()

	token, err := manager.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Error("refresh token is empty")
	}

	claims, err := manager.VerifyRefreshToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, claims.UserID)
	}
}



func TestVerifyAccessToken_EmptyToken(t *testing.T) {
	manager := NewJWTManager("secret", 1*time.Hour, 24*time.Hour)
	_, err := manager.VerifyAccessToken("")
	if err == nil {
		t.Error("expected error for empty token")
	}
}

func TestVerifyRefreshToken_EmptyToken(t *testing.T) {
	manager := NewJWTManager("secret", 1*time.Hour, 24*time.Hour)
	_, err := manager.VerifyRefreshToken("")
	if err == nil {
		t.Error("expected error for empty token")
	}
}



func TestVerifyAccessToken_InvalidFormat(t *testing.T) {
	manager := NewJWTManager("secret", 1*time.Hour, 24*time.Hour)
	// Token dengan 2 parts (invalid)
	_, err := manager.VerifyAccessToken("header.payload")
	if err == nil {
		t.Error("expected error for token with 2 parts")
	}
	
	// Token dengan 4 parts (invalid)
	_, err = manager.VerifyAccessToken("a.b.c.d")
	if err == nil {
		t.Error("expected error for token with 4 parts")
	}
	
	// Token dengan non-base64 characters
	_, err = manager.VerifyAccessToken("a.b.c")
	if err == nil {
		t.Error("expected error for malformed base64")
	}
}

func TestVerifyRefreshToken_InvalidFormat(t *testing.T) {
	manager := NewJWTManager("secret", 1*time.Hour, 24*time.Hour)
	// Token dengan 2 parts
	_, err := manager.VerifyRefreshToken("header.payload")
	if err == nil {
		t.Error("expected error for token with 2 parts")
	}
	
	// Token dengan 4 parts
	_, err = manager.VerifyRefreshToken("a.b.c.d")
	if err == nil {
		t.Error("expected error for token with 4 parts")
	}
	
	// Token dengan non-base64
	_, err = manager.VerifyRefreshToken("a.b.c")
	if err == nil {
		t.Error("expected error for malformed base64")
	}
}

func TestGenerateAccessToken_EmptySecret(t *testing.T) {
	manager := NewJWTManager("", 1*time.Hour, 24*time.Hour)
	userID := uuid.New()
	token, err := manager.GenerateAccessToken(userID, "test@example.com", "admin")
	// JWT tetap generate token walaupun secret kosong (tidak aman, tapi tidak error)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestGenerateRefreshToken_EmptySecret(t *testing.T) {
	manager := NewJWTManager("", 1*time.Hour, 24*time.Hour)
	userID := uuid.New()
	token, err := manager.GenerateRefreshToken(userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}
