#!/bin/bash

set -e

echo "📦 Generating Unit Tests for ScIntel (ROADMAP 1 & 2)"
echo "===================================================="

# ============================================
# PART 1: UTILITY TESTS
# ============================================
mkdir -p internal/util

cat > internal/util/hash_test.go << 'EOF'
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
EOF

cat > internal/util/jwt_test.go << 'EOF'
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
EOF

cat > internal/util/pagination_test.go << 'EOF'
package util

import (
	"net/http"
	"net/url"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name     string
		query    map[string]string
		lim, off int
	}{
		{"defaults", map[string]string{}, 20, 0},
		{"limit only", map[string]string{"limit": "50"}, 50, 0},
		{"offset only", map[string]string{"offset": "10"}, 20, 10},
		{"both", map[string]string{"limit": "30", "offset": "5"}, 30, 5},
		{"invalid", map[string]string{"limit": "abc"}, 20, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{URL: &url.URL{RawQuery: ""}}
			q := req.URL.Query()
			for k, v := range tt.query {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			limit, offset := ParsePagination(req)
			assert.Equal(t, tt.lim, limit)
			assert.Equal(t, tt.off, offset)
		})
	}
}
EOF

cat > internal/util/crypto_test.go << 'EOF'
package util

import (
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, HashSHA256(tt.input))
		})
	}
}
EOF

cat > internal/util/response_test.go << 'EOF'
package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
	w := httptest.NewRecorder()
	requestID := "req-123"
	WriteResponse(w, http.StatusOK, "Success", map[string]string{"id": "1"}, requestID)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp ResponseWrapper
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, requestID, resp.RequestID)
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	requestID := "req-456"
	WriteError(w, http.StatusBadRequest, "Bad Request", "Invalid", requestID)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGenerateRequestID(t *testing.T) {
	id1 := GenerateRequestID()
	id2 := GenerateRequestID()
	assert.NotEmpty(t, id1)
	assert.NotEqual(t, id1, id2)
}
EOF

echo "✅ Utility tests created."

# ============================================
# PART 2: MODEL TESTS
# ============================================
mkdir -p internal/model

cat > internal/model/user_test.go << 'EOF'
package model

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidRole(t *testing.T) {
	assert.True(t, ValidRole("admin"))
	assert.False(t, ValidRole("superadmin"))
}

func TestValidStatus(t *testing.T) {
	assert.True(t, ValidStatus("active"))
	assert.False(t, ValidStatus("deleted"))
}

func TestRoleHierarchy(t *testing.T) {
	assert.True(t, RoleHierarchy("admin", "viewer"))
	assert.False(t, RoleHierarchy("analyst", "manager"))
}

func TestUserToResponse(t *testing.T) {
	now := time.Now().UTC()
	user := &User{
		ID: uuid.New(),
		Email: "test@example.com",
		PasswordHash: "hash",
		FullName: "Test",
		Role: "admin",
		Status: "active",
		CreatedAt: now,
	}
	resp := user.ToResponse()
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Empty(t, resp.PasswordHash)
}
EOF

cat > internal/model/inventory_test.go << 'EOF'
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
EOF

cat > internal/model/transaction_test.go << 'EOF'
package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTransactionType(t *testing.T) {
	val, _ := TransactionIn.Value()
	assert.Equal(t, "in", val)
	var txn TransactionType
	txn.Scan("out")
	assert.Equal(t, TransactionOut, txn)
}
EOF

cat > internal/model/forecast_test.go << 'EOF'
package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAlgorithmConstants(t *testing.T) {
	assert.Equal(t, Algorithm("exponential_smoothing"), AlgorithmETS)
	assert.Equal(t, Algorithm("arima"), AlgorithmARIMA)
}
EOF

cat > internal/model/anomaly_test.go << 'EOF'
package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSeverityConstants(t *testing.T) {
	assert.Equal(t, Severity("low"), SeverityLow)
	assert.Equal(t, Severity("critical"), SeverityCritical)
}
EOF

cat > internal/model/webhook_test.go << 'EOF'
package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryStatus(t *testing.T) {
	assert.Equal(t, DeliveryStatus("pending"), DeliveryPending)
	assert.Equal(t, DeliveryStatus("success"), DeliverySuccess)
}
EOF

cat > internal/model/api_key_test.go << 'EOF'
package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAPIKey(t *testing.T) {
	key := &APIKey{Permissions: []string{"inventory:read"}}
	assert.Len(t, key.Permissions, 1)
}
EOF

echo "✅ Model tests created."

# ============================================
# PART 3: DATABASE, CONFIG, MIDDLEWARE (PLACEHOLDER)
# ============================================
mkdir -p internal/database internal/config internal/middleware

cat > internal/database/connection_test.go << 'EOF'
package database

import (
	"testing"
)

func TestConnectionConfig(t *testing.T) {
	cfg := ConnectionConfig{URL: "postgres://user:pass@localhost/db"}
	if cfg.URL == "" {
		t.Fatal("empty URL")
	}
}
EOF

cat > internal/database/transaction_test.go << 'EOF'
package database

import "testing"

func TestTxFunc(t *testing.T) {
	var fn TxFunc = func(ctx context.Context, tx *sql.Tx) error { return nil }
	if fn == nil {
		t.Fatal("fn is nil")
	}
}
EOF

cat > internal/config/config_test.go << 'EOF'
package config

import "testing"

func TestLoadConfig(t *testing.T) {
	// minimal placeholder
}
EOF

cat > internal/middleware/auth_test.go << 'EOF'
package middleware

import "testing"

func TestAuthMiddleware(t *testing.T) {
	// placeholder
}
EOF

cat > internal/middleware/rbac_test.go << 'EOF'
package middleware

import "testing"

func TestRBACMiddleware(t *testing.T) {
	// placeholder
}
EOF

cat > internal/middleware/rate_limit_test.go << 'EOF'
package middleware

import "testing"

func TestRateLimitMiddleware(t *testing.T) {
	// placeholder
}
EOF

echo "✅ Database, config, middleware tests created."

# ============================================
# PART 4: REPOSITORY & SERVICE (PLACEHOLDER)
# ============================================
mkdir -p internal/repository internal/service

for repo in user inventory transaction forecast anomaly webhook api_key; do
	cat > internal/repository/${repo}_repository_test.go << EOF
package repository

import "testing"

func Test${repo^}Repository(t *testing.T) {
	// placeholder
}
EOF
done

for svc in inventory transaction forecast anomaly webhook; do
	cat > internal/service/${svc}_service_test.go << EOF
package service

import "testing"

func Test${svc^}Service(t *testing.T) {
	// placeholder
}
EOF
done

echo "✅ Repository & Service tests created."

# ============================================
# PART 5: ALGORITHM & HANDLER (PLACEHOLDER)
# ============================================
mkdir -p pkg/timeseries internal/handler

for algo in ets arima anomaly; do
	cat > pkg/timeseries/${algo}_test.go << EOF
package timeseries

import "testing"

func Test${algo^}(t *testing.T) {
	// placeholder
}
EOF
done

for h in health inventory forecast anomaly webhook api_key; do
	cat > internal/handler/${h}_handler_test.go << EOF
package handler

import "testing"

func Test${h^}Handler(t *testing.T) {
	// placeholder
}
EOF
done

echo "✅ Algorithm & Handler tests created."

# ============================================
# FINALIZE
# ============================================
echo ""
echo "🔨 Installing testify..."
go get github.com/stretchr/testify

echo ""
echo "🚀 Running all unit tests (short mode)..."
go test -v ./... -short

echo ""
echo "📈 Coverage report..."
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | tail -20

echo ""
echo "✅ All unit test files generated successfully!"

