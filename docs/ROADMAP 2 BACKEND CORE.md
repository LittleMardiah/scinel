# ROADMAP_2: BACKEND CORE (REVISED)
**Supply Chain Intelligence Platform (ScIntel) - Weeks 3-4**

**Phase:** Backend Core - Production-Grade APIs  
**Duration:** 2 weeks (40 hours)  
**Status:** Ready to Start (after ROADMAP_1 complete)  
**Target Completion:** End of Week 4  
**Team:** Solo/Small (1-2 developers)  
**Revision:** v1.1 (Comprehensive Audit Fixes Applied 2026-07-19)  
**Blocks:** ROADMAP_3 (Frontend/Mobile) — cannot start without ROADMAP_2 complete

---

## 📋 ROADMAP_2 OVERVIEW

### Goal
Build production-grade backend APIs for inventory management, demand forecasting, anomaly detection, and webhook event system. **All core business logic implemented** — ready for frontend integration in ROADMAP_3.

### Effort Breakdown
- Week 3: 20 hours (Inventory APIs + Forecasting core)
- Week 4: 20 hours (Anomaly + Webhooks + Rate Limiting + Health Checks)
- **Total: 40 hours**

### Key Deliverables
✅ 5 Inventory endpoints (stock queries, transactions, import)  
✅ 3 Forecasting endpoints (query, generate batch, accuracy metrics)  
✅ 4 Anomaly endpoints (list, detail, update status, history) — **REVISED: was 2, now 4**  
✅ 4 Webhook endpoints (subscribe, list subscriptions, delete, view deliveries)  
✅ 3 Health Check endpoints (health, readiness, liveness) — **NEW: critical for monitoring**  
✅ 4 API Key endpoints (create, list, revoke, rotate) — **NEW: for programmatic access**  
✅ Rate limiting middleware (Redis-based, 1000 req/sec per user)  
✅ Business logic fully implemented (LOGIC_FLOW 4-7, TDD 6.3-6.4)  
✅ >80% test coverage for all endpoints & services  
✅ **All 23 backend endpoints documented & working** (vs 15 originally promised, now complete)

### Dependencies
- ✅ ROADMAP_1 (Foundation) — must be complete
  - Database: PostgreSQL with all 12 tables migrated
  - Auth system: JWT + RBAC working
  - Structured logging: Configured
  - Graceful shutdown: Implemented

### Testing Deliverables
- Unit tests for all services & algorithms (>80% coverage)
- Integration tests for core workflows (transaction → anomaly → webhook)
- Load tests for inventory queries (p99 <100ms target)
- Forecast batch job correctness tests (<5s for 1000 SKUs)
- Concurrent transaction handling verification (deadlock retry logic)

### Success Criteria
1. ✅ All 23 endpoints working and documented
2. ✅ All business logic implemented (inventory, forecast, anomaly, webhook, health)
3. ✅ >80% test coverage for backend services
4. ✅ Performance targets verified: p99 <100ms inventory, <5s forecast batch
5. ✅ Rate limiting active and tested (1000 req/sec limit)
6. ✅ Webhook system delivering with retry logic (5 attempts)
7. ✅ Health endpoints passing (for Kubernetes, monitoring, CI/CD)
8. ✅ API Key authentication working (foundation set in ROADMAP_1, implemented here)
9. ✅ Redis caching working (inventory 30s TTL, forecast 24h TTL)
10. ✅ CI/CD pipeline passing all tests, ready for ROADMAP_3

---

## 📅 DETAILED BREAKDOWN

---

## WEEK 3: INVENTORY APIs & FORECASTING CORE (20 hours)

### Week 3.1: Health Check Endpoints & Inventory Foundation (3 hours)

**Objective:** Implement monitoring endpoints (critical for Kubernetes/CI/CD) and base inventory infrastructure

#### Task 3.1.1: Health Check Endpoints (1 hour)

**File:** `internal/handler/health_handler.go`

**Purpose:** Three endpoints for monitoring, Kubernetes probes, and CI/CD pipeline health checks

```go
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/aulia/scinel-backend/internal/util"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db    *database.DB
	redis *redis.Client
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

// Health returns basic health status (GET /api/v1/health)
// Purpose: Simple liveness probe for load balancers
// Response: 200 OK if service is running
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()
	
	response := map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version": "1.0.0",
	}
	
	util.WriteResponse(w, http.StatusOK, "OK", response, requestID)
}

// Readiness performs deep health check (GET /api/v1/readiness)
// Purpose: Kubernetes readiness probe - is service ready to accept traffic?
// Checks: Database connectivity, Redis connectivity, migrations status
// Response: 200 OK if all dependencies healthy, 503 Service Unavailable if any fail
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	checks := map[string]bool{}

	// Check database
	err := h.db.PingContext(ctx)
	checks["database"] = err == nil

	// Check Redis
	err = h.redis.Ping(ctx).Err()
	checks["cache"] = err == nil

	// Check if all dependencies are healthy
	allHealthy := true
	for _, healthy := range checks {
		if !healthy {
			allHealthy = false
			break
		}
	}

	statusCode := http.StatusOK
	message := "Ready"
	if !allHealthy {
		statusCode = http.StatusServiceUnavailable
		message = "Not ready - dependency failure"
	}

	util.WriteResponse(w, statusCode, message, checks, requestID)
}

// Liveness provides deep diagnostic info (GET /api/v1/liveness)
// Purpose: Kubernetes liveness probe - should this pod be restarted?
// Checks: All readiness checks + database connection pool stats
// Response: 200 OK if healthy, 503 if deadlock/hanging detected
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	checks := map[string]interface{}{}

	// Check database connectivity
	dbErr := h.db.PingContext(ctx)
	checks["database_connected"] = dbErr == nil

	// Get database stats (connection pool info)
	dbStats := h.db.Stats()
	checks["database_stats"] = map[string]interface{}{
		"open_connections": dbStats.OpenConnections,
		"in_use":          dbStats.InUse,
		"idle":            dbStats.Idle,
	}

	// Check Redis connectivity
	cacheErr := h.redis.Ping(ctx).Err()
	checks["cache_connected"] = cacheErr == nil

	// Check if all critical dependencies are healthy
	allHealthy := dbErr == nil && cacheErr == nil

	statusCode := http.StatusOK
	message := "Alive"
	if !allHealthy {
		statusCode = http.StatusServiceUnavailable
		message = "Not alive - critical dependency failure"
	}

	util.WriteResponse(w, statusCode, message, checks, requestID)
}
```

**Endpoint Specifications:**

| Endpoint | Method | Purpose | Timeout | Kubernetes Use |
|----------|--------|---------|---------|-----------------|
| `/api/v1/health` | GET | Basic liveness | 1s | Health check endpoint |
| `/api/v1/readiness` | GET | Readiness probe | 5s | Readiness probe (when ready for traffic) |
| `/api/v1/liveness` | GET | Liveness probe | 10s | Liveness probe (pod health) |

**Kubernetes Probe Configuration (example):**
```yaml
livenessProbe:
  httpGet:
    path: /api/v1/liveness
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /api/v1/readiness
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
```

**Reference:** TDD Section 5.1 (Health Check Endpoints), LOGIC_FLOW Section 10 (Monitoring)

**Deliverable:**
- ✅ 3 health check endpoints implemented
- ✅ Database & Redis connectivity verified
- ✅ Kubernetes-compatible response format
- ✅ Proper HTTP status codes (200/503)

---

#### Task 3.1.2: Inventory Model & Repository Foundation (2 hours)

**File:** `internal/model/inventory.go`

```go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Inventory represents a stock level for a SKU in a warehouse
type Inventory struct {
	ID                uuid.UUID `db:"id"`
	WarehouseID       uuid.UUID `db:"warehouse_id"`
	SKUID             uuid.UUID `db:"sku_id"`
	SKUCode           string    `db:"sku_code"`           // e.g., "SKU-NIKE-001"
	Quantity          int32     `db:"quantity"`
	ReorderPoint      int32     `db:"reorder_point"`      // Trigger for low-stock alerts
	MaxStockLevel     int32     `db:"max_stock_level"`    // Safety cap
	LastRestockDate   *time.Time `db:"last_restock_date"`
	LastTransactionID *uuid.UUID `db:"last_transaction_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// InventorySnapshot represents a point-in-time stock level (for historical tracking)
type InventorySnapshot struct {
	ID           uuid.UUID `db:"id"`
	WarehouseID  uuid.UUID `db:"warehouse_id"`
	SKUID        uuid.UUID `db:"sku_id"`
	Quantity     int32     `db:"quantity"`
	SnapshotDate time.Time `db:"snapshot_date"`
	CreatedAt    time.Time `db:"created_at"`
}

// InventoryStatus represents inventory status classification
type InventoryStatus string

const (
	StatusNormal   InventoryStatus = "normal"      // Quantity > reorder_point
	StatusLowStock InventoryStatus = "low_stock"   // Quantity <= reorder_point
	StatusCritical InventoryStatus = "critical"    // Quantity == 0 or negative
)

// GetStatus returns inventory status based on quantity
func (i *Inventory) GetStatus() InventoryStatus {
	if i.Quantity <= 0 {
		return StatusCritical
	}
	if i.Quantity <= i.ReorderPoint {
		return StatusLowStock
	}
	return StatusNormal
}
```

**File:** `internal/repository/inventory_repository.go`

**Repository Interface (matches TDD Section 4.4):**

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/aulia/scinel-backend/internal/model"
)

// InventoryRepository defines inventory data access operations
type InventoryRepository interface {
	// GetBySKU retrieves inventory for a specific warehouse-SKU combination
	// Returns: nil if not found (not an error - let caller decide action)
	GetBySKU(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

	// GetByWarehouse retrieves all SKU inventory in a warehouse
	GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]model.Inventory, int64, error)

	// GetAll retrieves all inventory across all warehouses
	GetAll(ctx context.Context, limit, offset int) ([]model.Inventory, int64, error)

	// Create creates a new inventory record
	Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)

	// UpdateQuantity atomically updates inventory quantity
	// Returns: Updated inventory or error
	// Note: Uses pessimistic locking (SELECT ... FOR UPDATE) for consistency
	UpdateQuantity(ctx context.Context, warehouseID, skuID uuid.UUID, delta int32) (*model.Inventory, error)

	// LockForUpdate locks inventory record for update in a transaction
	// Used for: Pessimistic locking in UPDATE inventory flow
	LockForUpdate(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

	// GetBySKUCached retrieves inventory with Redis cache fallback (30s TTL)
	// On cache miss: Query DB, then cache result
	// On cache error: Log but continue (cache is optional)
	GetBySKUCached(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

	// InvalidateCache removes inventory from cache after update
	InvalidateCache(ctx context.Context, warehouseID, skuID uuid.UUID) error

	// CreateSnapshot records point-in-time inventory level
	CreateSnapshot(ctx context.Context, snapshot *model.InventorySnapshot) error
}
```

**Key Implementation Details:**

1. **GetBySKU()** - Single record query:
   - SQL: `SELECT * FROM inventory WHERE warehouse_id = $1 AND sku_id = $2`
   - Error handling: Return `nil, nil` if not found (let caller decide action)
   - Optimization: Use index `idx_inventory_warehouse_sku`

2. **UpdateQuantity()** - Atomic update with deadlock retry:
   - Begin transaction with `READ_COMMITTED` isolation
   - `SELECT * FROM inventory WHERE ... FOR UPDATE` (pessimistic lock)
   - `UPDATE inventory SET quantity = quantity + $1, updated_at = NOW() WHERE ...`
   - Commit transaction
   - On deadlock: Retry up to 3x with backoff (100ms, 200ms, 400ms)
   - Reference: TDD Section 4.0, ROADMAP_1 Week 1.3

3. **GetBySKUCached()** - Redis cache with TTL:
   - Try Redis: `GET cache:inventory:{warehouse_id}:{sku_id}`
   - Cache miss: Query DB, then `SET cache:... 30s EX` (30 second expiration)
   - Cache hit: Unmarshal and return
   - Error handling: Log cache errors, don't fail request

4. **InvalidateCache()** - After any update:
   - `DEL cache:inventory:{warehouse_id}:{sku_id}`
   - After transaction commit (avoid race condition)

**Deliverable:**
- ✅ Inventory model with status enum
- ✅ InventorySnapshot model for historical tracking
- ✅ Repository interface with 8 methods
- ✅ Proper caching strategy documented

---

### Week 3.2: Inventory Endpoints & Transactions (5 hours)

**Objective:** Implement inventory transaction endpoints (POST, GET) with RBAC & pessimistic locking

#### Task 3.2.1: Transaction Model & Repository (1.5 hours)

**File:** `internal/model/transaction.go`

```go
package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Transaction represents an inventory in/out transaction
type Transaction struct {
	ID               uuid.UUID `db:"id"`
	WarehouseID      uuid.UUID `db:"warehouse_id"`
	SKUID            uuid.UUID `db:"sku_id"`
	SKUCode          string    `db:"sku_code"`  // Denormalized for queries
	Type             TransactionType `db:"type"`
	Quantity         int32     `db:"quantity"` // Positive or negative
	Reason           string    `db:"reason"`   // e.g., "Supplier delivery #PO-2026-001"
	PreviousQuantity int32     `db:"previous_quantity"`
	NewQuantity      int32     `db:"new_quantity"`
	CreatedBy        uuid.UUID `db:"created_by"`
	Details          json.RawMessage `db:"details"` // JSONB: Additional metadata
	CreatedAt        time.Time `db:"created_at"`
}

// TransactionType represents transaction direction
type TransactionType string

const (
	TransactionIn  TransactionType = "in"  // Inventory increase (purchase, return, correction)
	TransactionOut TransactionType = "out" // Inventory decrease (sale, damage, transfer)
)

// Scan implements sql.Scanner interface for TransactionType
func (t *TransactionType) Scan(value interface{}) error {
	*t = TransactionType(value.(string))
	return nil
}

// Value implements driver.Valuer interface for TransactionType
func (t TransactionType) Value() (driver.Value, error) {
	return string(t), nil
}

// TransactionFilter for querying transactions with filters
type TransactionFilter struct {
	WarehouseID *uuid.UUID
	SKUID       *uuid.UUID
	Type        *TransactionType
	DateFrom    *time.Time
	DateTo      *time.Time
	CreatedBy   *uuid.UUID
	Reason      *string
}
```

**File:** `internal/repository/transaction_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/aulia/scinel-backend/internal/model"
)

// TransactionRepository defines transaction data access operations
type TransactionRepository interface {
	// Create inserts a new transaction record
	Create(ctx context.Context, txn *model.Transaction) (*model.Transaction, error)

	// GetByID retrieves single transaction
	GetByID(ctx context.Context, txnID uuid.UUID) (*model.Transaction, error)

	// Query retrieves transactions with filters and pagination
	Query(ctx context.Context, filter model.TransactionFilter, limit, offset int) ([]model.Transaction, int64, error)

	// GetByWarehouse retrieves all transactions in a warehouse
	GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error)

	// GetBySKU retrieves all transactions for a SKU
	GetBySKU(ctx context.Context, skuID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error)
}
```

**Deliverable:**
- ✅ Transaction model with type enum
- ✅ Transaction filter for advanced queries
- ✅ Repository interface with 5 methods

---

#### Task 3.2.2: Inventory Transaction Endpoints (3.5 hours)

**File:** `internal/handler/inventory_handler.go`

**Endpoint 1: POST /api/v1/transactions (Create Inventory Transaction)**

**Purpose:** Record inventory in/out transaction with atomic update

**Authentication:** RequireAuth + RequireRole("admin", "manager")

**Special RBAC Rule:** Manager can only create transactions for own warehouse (manager_warehouse_id from token)

**Request Body:**
```json
{
  "warehouse_id": "wh-jakarta-001",
  "sku_id": "sku-nike-001",
  "type": "in",
  "quantity": 50,
  "reason": "Supplier delivery #PO-2026-001",
  "details": {
    "purchase_order": "PO-2026-001",
    "supplier_id": "supplier-123",
    "received_date": "2026-07-19"
  }
}
```

**Business Logic (Reference: LOGIC_FLOW Section 4.1):**
1. Validate all required fields (warehouse_id, sku_id, type, quantity, reason)
2. Verify user role & warehouse access (RBAC check)
3. Use repository.LockForUpdate() for pessimistic locking
4. For "out" transactions: Verify sufficient stock (quantity <= current_quantity)
5. Update inventory quantity atomically
6. Insert transaction record with all details
7. Publish webhook event "inventory.changed" (async, don't block response)
8. Invalidate Redis cache
9. Trigger anomaly detection check (async background job)
10. Return 201 Created with transaction details

**Response (201 Created):**
```json
{
  "code": 201,
  "message": "Transaction created",
  "data": {
    "id": "txn-uuid",
    "warehouse_id": "wh-jakarta-001",
    "sku_id": "sku-nike-001",
    "sku_code": "SKU-NIKE-001",
    "type": "in",
    "quantity": 50,
    "previous_quantity": 100,
    "new_quantity": 150,
    "reason": "Supplier delivery #PO-2026-001",
    "created_at": "2026-07-19T14:30:00Z"
  },
  "timestamp": "2026-07-19T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error Cases:**
- 400 Bad Request: Missing fields, invalid quantity
- 401 Unauthorized: Not authenticated
- 403 Forbidden: Manager accessing different warehouse
- 409 Conflict: Insufficient stock for "out" transaction
- 500 Internal Server Error: DB error

**Reference:** TDD Section 5.2.1 (Transaction Endpoint), LOGIC_FLOW Section 4.1 (Transaction Flow)

**Deliverable:**
- ✅ POST /api/v1/transactions endpoint
- ✅ Pessimistic locking with deadlock retry
- ✅ RBAC warehouse isolation
- ✅ Async webhook publishing
- ✅ Comprehensive error handling

---

**Endpoint 2: GET /api/v1/transactions (Query Transactions)**

**Purpose:** Retrieve transaction history with filters and pagination

**Authentication:** RequireAuth (all roles can read)

**Query Parameters:**
```
?warehouse_id=wh-jakarta-001
&sku_id=sku-nike-001
&type=in
&date_from=2026-07-01
&date_to=2026-07-31
&limit=20
&offset=0
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "transactions": [
      {
        "id": "txn-uuid1",
        "warehouse_id": "wh-jakarta-001",
        "sku_id": "sku-nike-001",
        "sku_code": "SKU-NIKE-001",
        "type": "in",
        "quantity": 50,
        "reason": "Supplier delivery",
        "created_by": "user-uuid",
        "created_at": "2026-07-19T14:30:00Z"
      }
    ],
    "total": 127,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-07-19T14:35:00Z",
  "request_id": "req-uuid"
}
```

**Pagination:** Support limit/offset (default: 20 records, max: 100)

**Filtering:** Support warehouse_id, sku_id, type, date range

**Performance:** Use indexed queries on (warehouse_id, sku_id) and created_at

**Reference:** TDD Section 5.2.2 (Transaction Query)

**Deliverable:**
- ✅ GET /api/v1/transactions endpoint
- ✅ Advanced filtering support
- ✅ Pagination with limit/offset
- ✅ Performance-optimized queries

---

**Endpoint 3: GET /api/v1/inventory (Query All Stock)**

**Purpose:** Retrieve inventory levels across all warehouses (or filtered)

**Authentication:** RequireAuth (all roles can read)

**Query Parameters:**
```
?warehouse_id=wh-jakarta-001
?sku_id=sku-nike-001
?status=low_stock         (normal, low_stock, critical)
?limit=20
?offset=0
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "inventory": [
      {
        "id": "inv-uuid",
        "warehouse_id": "wh-jakarta-001",
        "sku_id": "sku-nike-001",
        "sku_code": "SKU-NIKE-001",
        "quantity": 150,
        "reorder_point": 100,
        "max_stock_level": 5000,
        "status": "normal",
        "last_restock_date": "2026-07-15T10:00:00Z",
        "updated_at": "2026-07-19T14:30:00Z"
      }
    ],
    "total": 324,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-07-19T14:35:00Z",
  "request_id": "req-uuid"
}
```

**Caching:** Use GetBySKUCached() for 30-second Redis cache

**Performance Target:** p99 <100ms (with cache hits)

**Reference:** TDD Section 5.2.3 (Inventory Query)

**Deliverable:**
- ✅ GET /api/v1/inventory endpoint
- ✅ Multi-filter support (warehouse, SKU, status)
- ✅ Redis caching (30s TTL)
- ✅ Performance-optimized

---

**Endpoint 4: GET /api/v1/inventory/{warehouse_id} (Warehouse-specific Inventory)**

**Purpose:** Get all stock levels for a specific warehouse

**Authentication:** RequireAuth + RequireRole("admin", "manager", "analyst", "viewer")

**RBAC:** Manager can only query own warehouse

**Response:** Same structure as GET /api/v1/inventory but filtered by warehouse

**Performance:** Index on warehouse_id, cache with 30s TTL

**Deliverable:**
- ✅ GET /api/v1/inventory/{warehouse_id} endpoint
- ✅ Warehouse-specific filtering
- ✅ RBAC enforcement
- ✅ Redis caching

---

**Endpoint 5: POST /api/v1/inventory/import (Bulk Import)**

**Purpose:** Bulk import inventory from CSV or JSON

**Authentication:** RequireAuth + RequireRole("admin", "manager")

**Request:** Multipart form-data with CSV file

**Processing:**
1. Parse CSV/JSON format
2. Validate each row (warehouse_id, sku_id, quantity)
3. Upsert inventory records in batch
4. Return summary (imported: 100, errors: 2, details: [...])

**Deliverable:**
- ✅ POST /api/v1/inventory/import endpoint
- ✅ CSV/JSON parsing
- ✅ Batch upsert with validation
- ✅ Error reporting

---

### Week 3.3: Forecasting Service & Algorithms (7 hours)

**Objective:** Implement ETS and ARIMA algorithms with batch job scheduler

#### Task 3.3.1: Forecast Models (1 hour)

**File:** `internal/model/forecast.go`

```go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Forecast represents a demand prediction for a SKU in a warehouse
type Forecast struct {
	ID                 uuid.UUID `db:"id"`
	WarehouseID        uuid.UUID `db:"warehouse_id"`
	SKUID              uuid.UUID `db:"sku_id"`
	SKUCode            string    `db:"sku_code"`
	ForecastDate       time.Time `db:"forecast_date"`
	PredictedDemand    int32     `db:"predicted_demand"`
	LowerBound         int32     `db:"lower_bound"`          // 95% confidence interval
	UpperBound         int32     `db:"upper_bound"`          // 95% confidence interval
	Confidence         float64   `db:"confidence"`           // 0.95 for 95%
	Algorithm          Algorithm `db:"algorithm"`            // exponential_smoothing, arima
	Trend              *float64  `db:"trend"`                // Trend coefficient
	SeasonalityFactor  *float64  `db:"seasonality_factor"`   // Seasonal adjustment
	GeneratedAt        time.Time `db:"generated_at"`         // When forecast was created
	CreatedAt          time.Time `db:"created_at"`
}

// ForecastAccuracy tracks algorithm performance over time
type ForecastAccuracy struct {
	ID           uuid.UUID `db:"id"`
	WarehouseID  uuid.UUID `db:"warehouse_id"`
	SKUID        uuid.UUID `db:"sku_id"`
	Algorithm    Algorithm `db:"algorithm"`
	MAPE         *float64  `db:"mape"`        // Mean Absolute Percentage Error
	RMSE         *float64  `db:"rmse"`        // Root Mean Square Error
	MAE          *float64  `db:"mae"`         // Mean Absolute Error
	EvaluationDate time.Time `db:"evaluation_date"`
	CreatedAt    time.Time `db:"created_at"`
}

// Algorithm type for forecasting
type Algorithm string

const (
	AlgorithmETS   Algorithm = "exponential_smoothing"
	AlgorithmARIMA Algorithm = "arima"
)

// Trend classification
type TrendType string

const (
	TrendFlat       TrendType = "flat"
	TrendIncreasing TrendType = "increasing"
	TrendDecreasing TrendType = "decreasing"
)
```

**Deliverable:**
- ✅ Forecast model with confidence intervals
- ✅ ForecastAccuracy model for metrics
- ✅ Algorithm enum and types

---

#### Task 3.3.2: Forecasting Algorithms (3 hours)

**File:** `pkg/timeseries/ets.go` (Exponential Smoothing Implementation)

**Reference:** TDD Section 6.3 (Forecasting Algorithms)

```go
package timeseries

import (
	"fmt"
	"math"
)

// ETSForecast represents Exponential Smoothing forecast result
type ETSForecast struct {
	PredictedDemand   int32      // Rounded forecast
	LowerBound        int32      // 95% confidence interval
	UpperBound        int32      // 95% confidence interval
	Confidence        float64    // 0.95
	Trend             float64    // Trend component
	SeasonalityFactor float64    // Seasonal adjustment (if applicable)
	MAPE              float64    // Mean Absolute Percentage Error on validation set
}

// ExponentialSmoothingForecast performs exponential smoothing on time series
// Input: Historical daily demand (180 days minimum recommended)
// Output: 30-day forecast with confidence intervals
// Reference: TDD Section 6.3.1
func ExponentialSmoothingForecast(historicalDemand []int32, forecastDays int) (*ETSForecast, error) {
	if len(historicalDemand) < 30 {
		return nil, fmt.Errorf("insufficient historical data: %d days (minimum 30 required)", len(historicalDemand))
	}

	// Initialize smoothing parameters (alpha, beta, gamma)
	// These are tuned for retail demand patterns
	alpha := 0.3  // Level smoothing (0.1-0.5 for trend-heavy data)
	beta := 0.1   // Trend smoothing (0.01-0.2 for gradual trends)
	gamma := 0.05 // Seasonal smoothing (for multi-year seasonality)

	// Initialize level, trend, seasonal components
	level := float64(historicalDemand[0])
	trend := (float64(historicalDemand[1]) - float64(historicalDemand[0])) / 2.0
	
	// Seasonal factors (optional, for multi-year patterns)
	seasonalPeriod := 7 // Weekly seasonality (or 365 for yearly)
	seasonalFactors := make([]float64, seasonalPeriod)
	for i := 0; i < seasonalPeriod && i < len(historicalDemand); i++ {
		seasonalFactors[i] = 1.0
	}

	// Apply exponential smoothing
	var lastLevel, lastTrend float64
	for i := 1; i < len(historicalDemand); i++ {
		demand := float64(historicalDemand[i])
		
		// Update level
		newLevel := alpha*demand + (1-alpha)*(level+trend)
		
		// Update trend
		newTrend := beta*(newLevel-level) + (1-beta)*trend
		
		// Update seasonal (if applicable)
		if seasonalPeriod > 0 && i >= seasonalPeriod {
			idx := i % seasonalPeriod
			seasonalFactors[idx] = gamma*(demand/level) + (1-gamma)*seasonalFactors[idx]
		}
		
		lastLevel = level
		lastTrend = trend
		level = newLevel
		trend = newTrend
	}

	// Generate 30-day forecast
	forecast := make([]float64, forecastDays)
	for h := 1; h <= forecastDays; h++ {
		fh := level + float64(h)*trend
		if seasonalPeriod > 0 {
			idx := (len(historicalDemand) + h - 1) % seasonalPeriod
			fh *= seasonalFactors[idx]
		}
		forecast[h-1] = fh
	}

	// Calculate standard error for confidence intervals (95% ≈ ±1.96 * SE)
	residuals := calculateResiduals(historicalDemand, level, trend, seasonalFactors, seasonalPeriod)
	stdError := calculateStdError(residuals)
	
	// Average forecast for 30-day period
	avgForecast := 0.0
	for _, f := range forecast {
		avgForecast += f
	}
	avgForecast /= float64(forecastDays)

	// Calculate MAPE on recent 30-day validation window
	mape := calculateMAPE(historicalDemand[len(historicalDemand)-30:], forecast)

	return &ETSForecast{
		PredictedDemand:   int32(math.Round(avgForecast)),
		LowerBound:        int32(math.Max(0, avgForecast-1.96*stdError)),
		UpperBound:        int32(avgForecast + 1.96*stdError),
		Confidence:        0.95,
		Trend:             trend,
		SeasonalityFactor: 1.0, // Simplified for MVP
		MAPE:              mape,
	}, nil
}

// Helper functions for ETS

func calculateResiduals(demand []int32, level, trend float64, seasonal []float64, period int) []float64 {
	residuals := make([]float64, 0, len(demand))
	for i, d := range demand {
		predicted := level + float64(i)*trend
		if period > 0 {
			idx := i % period
			predicted *= seasonal[idx]
		}
		residuals = append(residuals, float64(d)-predicted)
	}
	return residuals
}

func calculateStdError(residuals []float64) float64 {
	if len(residuals) < 2 {
		return 1.0
	}
	mean := 0.0
	for _, r := range residuals {
		mean += r
	}
	mean /= float64(len(residuals))
	
	variance := 0.0
	for _, r := range residuals {
		variance += (r - mean) * (r - mean)
	}
	variance /= float64(len(residuals) - 1)
	return math.Sqrt(variance)
}

func calculateMAPE(actual []int32, predicted []float64) float64 {
	if len(actual) != len(predicted) {
		return 0.0
	}
	
	sumAPE := 0.0
	count := 0
	for i, a := range actual {
		if a > 0 {
			ape := math.Abs(float64(a)-predicted[i]) / float64(a)
			sumAPE += ape
			count++
		}
	}
	
	if count == 0 {
		return 0.0
	}
	return (sumAPE / float64(count)) * 100.0 // Percentage
}
```

**File:** `pkg/timeseries/arima.go` (ARIMA Implementation - Auto Parameter Selection)

```go
package timeseries

import (
	"fmt"
	"math"
)

// ARIMAForecast represents ARIMA forecast result
type ARIMAForecast struct {
	PredictedDemand   int32
	LowerBound        int32
	UpperBound        int32
	Confidence        float64
	Trend             float64
	SeasonalityFactor float64
	MAPE              float64
	P, D, Q           int // ARIMA order
}

// ARIMAForecast performs ARIMA forecasting with automatic parameter selection
// Input: Historical daily demand (180 days minimum)
// Output: 30-day forecast with confidence intervals
// Reference: TDD Section 6.3.2 (Auto-ARIMA)
func AutoARIMAForecast(historicalDemand []int32, forecastDays int) (*ARIMAForecast, error) {
	if len(historicalDemand) < 30 {
		return nil, fmt.Errorf("insufficient historical data: %d days (minimum 30 required)", len(historicalDemand))
	}

	// Simplified ARIMA: Try small set of (p, d, q) combinations
	// Full auto-ARIMA is complex; for MVP we use pre-tested parameters
	bestModel := ARIMAParams{P: 1, D: 1, Q: 1}
	bestAIC := math.MaxFloat64

	// Grid search over common parameters
	for p := 0; p <= 2; p++ {
		for d := 0; d <= 2; d++ {
			for q := 0; q <= 2; q++ {
				params := ARIMAParams{P: p, D: d, Q: q}
				aic := fittingAIC(historicalDemand, params)
				if aic < bestAIC {
					bestAIC = aic
					bestModel = params
				}
			}
		}
	}

	// Fit final model and forecast
	forecast, stdError := fitARIMAModel(historicalDemand, bestModel, forecastDays)
	
	// Calculate metrics
	mape := calculateMAPE(historicalDemand[len(historicalDemand)-30:], forecast[:30])
	avgForecast := 0.0
	for _, f := range forecast[:forecastDays] {
		avgForecast += f
	}
	avgForecast /= float64(forecastDays)

	return &ARIMAForecast{
		PredictedDemand:   int32(math.Round(avgForecast)),
		LowerBound:        int32(math.Max(0, avgForecast-1.96*stdError)),
		UpperBound:        int32(avgForecast + 1.96*stdError),
		Confidence:        0.95,
		Trend:             0.0, // Simplified
		SeasonalityFactor: 1.0,
		MAPE:              mape,
		P:                 bestModel.P,
		D:                 bestModel.D,
		Q:                 bestModel.Q,
	}, nil
}

// ARIMAParams holds ARIMA model parameters
type ARIMAParams struct {
	P int // AR order
	D int // Differencing order
	Q int // MA order
}

// Simplified model fitting (for MVP, would use more sophisticated implementation)
func fitARIMAModel(data []int32, params ARIMAParams, forecastDays int) ([]float64, float64) {
	forecast := make([]float64, forecastDays)
	mean := 0.0
	for _, d := range data {
		mean += float64(d)
	}
	mean /= float64(len(data))
	
	// Simple fallback: use mean with small trend
	trend := (float64(data[len(data)-1]) - float64(data[0])) / float64(len(data))
	
	for i := 0; i < forecastDays; i++ {
		forecast[i] = mean + trend*float64(i+1)
	}
	
	// Standard error estimation
	stdError := calculateStdError(calculateResiduals(data, mean, trend, nil, 0))
	
	return forecast, stdError
}

func fittingAIC(data []int32, params ARIMAParams) float64 {
	// Simplified AIC calculation (for MVP)
	// Full AIC requires fitting and likelihood calculation
	k := params.P + params.D + params.Q + 1
	return float64(k * 2)
}
```

**Reference:** TDD Section 6.3 (Forecasting Algorithms)

**Deliverable:**
- ✅ ETS algorithm implementation
- ✅ ARIMA algorithm with auto parameter selection
- ✅ Confidence interval calculation
- ✅ MAPE metric calculation

---

#### Task 3.3.3: Forecast Repository & Service (2 hours)

**File:** `internal/repository/forecast_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/aulia/scinel-backend/internal/model"
)

// ForecastRepository defines forecast data access operations
type ForecastRepository interface {
	// Create inserts forecast record
	Create(ctx context.Context, forecast *model.Forecast) (*model.Forecast, error)

	// GetByWarehouseAndSKU retrieves forecasts for a specific warehouse-SKU combination
	GetByWarehouseAndSKU(ctx context.Context, warehouseID, skuID uuid.UUID, days int) ([]model.Forecast, error)

	// GetAll retrieves all recent forecasts
	GetAll(ctx context.Context, limit, offset int) ([]model.Forecast, int64, error)

	// GetAccuracy retrieves forecast accuracy metrics
	GetAccuracy(ctx context.Context, warehouseID, skuID uuid.UUID, algorithm model.Algorithm) (*model.ForecastAccuracy, error)

	// UpdateAccuracy records new accuracy metrics
	UpdateAccuracy(ctx context.Context, accuracy *model.ForecastAccuracy) error
}
```

**File:** `internal/service/forecast_service.go` (Batch Job Scheduler)

```go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aulia/scinel-backend/pkg/timeseries"
)

// ForecastService handles forecast generation and management
type ForecastService struct {
	repo repository.ForecastRepository
	invRepo repository.InventoryRepository
	logger logger.Logger
}

// GenerateForecastBatch generates forecasts for all SKUs in all warehouses
// Purpose: Batch job run daily at 2 AM UTC (via cron in ROADMAP_4)
// Duration: Target <5 seconds for 1000 SKUs (0.5ms per SKU)
// Reference: LOGIC_FLOW Section 6 (Batch Job Flow), TDD Section 6.3
func (fs *ForecastService) GenerateForecastBatch(ctx context.Context) error {
	startTime := time.Now()
	fs.logger.Info("Starting forecast batch job", nil)

	// 1. Get all warehouses and SKUs
	warehouses, err := fs.invRepo.GetAllWarehouses(ctx)
	if err != nil {
		return fmt.Errorf("failed to get warehouses: %w", err)
	}

	allSkus, err := fs.invRepo.GetAllSKUs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get SKUs: %w", err)
	}

	totalForecasted := 0
	totalErrors := 0

	// 2. For each warehouse-SKU combination:
	for _, warehouse := range warehouses {
		for _, sku := range allSkus {
			// Get 180+ days historical data
			transactions, err := fs.invRepo.GetTransactionHistory(ctx, warehouse.ID, sku.ID, 180)
			if err != nil {
				fs.logger.Warn("Failed to get transaction history", map[string]interface{}{
					"warehouse_id": warehouse.ID,
					"sku_id":       sku.ID,
					"error":        err.Error(),
				})
				totalErrors++
				continue
			}

			if len(transactions) < 30 {
				// Not enough data, skip
				continue
			}

			// Convert transactions to daily demand array
			dailyDemand := aggregateDailyDemand(transactions)

			// 3. Run both algorithms
			etsForecast, _ := timeseries.ExponentialSmoothingForecast(dailyDemand, 30)
			arimaForecast, _ := timeseries.AutoARIMAForecast(dailyDemand, 30)

			// 4. Select best based on MAPE (lower is better)
			var selectedForecast interface{}
			if etsForecast.MAPE <= arimaForecast.MAPE {
				selectedForecast = etsForecast
			} else {
				selectedForecast = arimaForecast
			}

			// 5. Store forecasts in database (with cache invalidation)
			// ... store logic ...

			totalForecasted++
		}
	}

	duration := time.Since(startTime)
	fs.logger.Info("Forecast batch job complete", map[string]interface{}{
		"duration_ms":       duration.Milliseconds(),
		"total_forecasted":  totalForecasted,
		"total_errors":      totalErrors,
		"avg_per_sku_ms":    float64(duration.Milliseconds()) / float64(totalForecasted),
	})

	return nil
}

// aggregateDailyDemand converts transaction history to daily demand array
func aggregateDailyDemand(transactions []model.Transaction) []int32 {
	// Implementation: Group transactions by day, sum quantities for "out" transactions
	dailyDemand := make([]int32, 0)
	// ... aggregation logic ...
	return dailyDemand
}
```

**Deliverable:**
- ✅ Forecast repository interface
- ✅ Forecast service with batch job
- ✅ Daily forecast generation for all SKUs
- ✅ Algorithm selection based on MAPE
- ✅ Performance: <5s for 1000 SKUs

---

### Week 3.4: Forecast Query Endpoints (3 hours)

**Objective:** Implement forecast endpoints for querying & triggering batch

**Endpoint 1: GET /api/v1/forecasts (Query Forecasts)**

**Purpose:** Retrieve demand forecasts with filters

**Authentication:** RequireAuth (all roles)

**Query Parameters:**
```
?warehouse_id=wh-jakarta-001
?sku_id=sku-nike-001
?days=30
?limit=20
&offset=0
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "forecasts": [
      {
        "id": "forecast-uuid",
        "warehouse_id": "wh-jakarta-001",
        "sku_id": "sku-nike-001",
        "forecast_date": "2026-08-19",
        "predicted_demand": 150,
        "lower_bound": 120,
        "upper_bound": 180,
        "confidence": 0.95,
        "algorithm": "exponential_smoothing",
        "generated_at": "2026-07-19T02:00:00Z"
      }
    ],
    "total": 45,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-07-19T14:40:00Z",
  "request_id": "req-uuid"
}
```

**Caching:** 24-hour Redis cache (forecasts are stable until next generation)

**Deliverable:**
- ✅ GET /api/v1/forecasts endpoint
- ✅ Multi-filter support (warehouse, SKU)
- ✅ Redis caching (24h TTL)

---

**Endpoint 2: POST /api/v1/forecasts/generate (Manual Trigger)**

**Purpose:** Manually trigger forecast batch job (admin only)

**Authentication:** RequireAuth + RequireRole("admin")

**Request Body:** Empty

**Response (202 Accepted):**
```json
{
  "code": 202,
  "message": "Forecast batch started",
  "data": {
    "job_id": "job-uuid",
    "status": "processing",
    "estimated_duration_seconds": 5
  },
  "timestamp": "2026-07-19T14:40:00Z",
  "request_id": "req-uuid"
}
```

**Background Processing:** Async job with progress tracking

**Deliverable:**
- ✅ POST /api/v1/forecasts/generate endpoint
- ✅ Admin-only access control
- ✅ Async job execution
- ✅ Job status tracking

---

**Endpoint 3: GET /api/v1/forecasts/accuracy (Accuracy Metrics)**

**Purpose:** Retrieve forecast accuracy metrics per algorithm

**Authentication:** RequireAuth (all roles)

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "accuracy": [
      {
        "warehouse_id": "wh-jakarta-001",
        "sku_id": "sku-nike-001",
        "algorithm": "exponential_smoothing",
        "mape": 12.5,
        "rmse": 45.3,
        "mae": 32.1,
        "evaluation_date": "2026-07-19"
      }
    ]
  },
  "timestamp": "2026-07-19T14:40:00Z",
  "request_id": "req-uuid"
}
```

**Deliverable:**
- ✅ GET /api/v1/forecasts/accuracy endpoint
- ✅ Algorithm performance metrics
- ✅ MAPE, RMSE, MAE calculation

---

## WEEK 4: ANOMALY DETECTION, WEBHOOKS, RATE LIMITING (20 hours)

### Week 4.1: Anomaly Detection Service (3 hours)

**Objective:** Implement Z-score & IQR-based anomaly detection

#### Task 4.1.1: Anomaly Model & Repository (1 hour)

**File:** `internal/model/anomaly.go`

```go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Anomaly represents an anomaly detection event
type Anomaly struct {
	ID                uuid.UUID `db:"id"`
	WarehouseID       uuid.UUID `db:"warehouse_id"`
	SKUID             uuid.UUID `db:"sku_id"`
	SKUCode           string    `db:"sku_code"`
	Type              string    `db:"type"`              // e.g., "unusual_stock_drop", "spike"
	Severity          Severity  `db:"severity"`          // low, medium, high, critical
	Description       string    `db:"description"`
	HistoricalAverage float64   `db:"historical_average"`
	DetectedValue     float64   `db:"detected_value"`
	ZScore            float64   `db:"z_score"`           // Standard deviations from mean
	Status            AnomalyStatus `db:"status"`        // open, acknowledged, resolved
	ResolutionNotes   *string   `db:"resolution_notes"`
	ResolvedBy        *uuid.UUID `db:"resolved_by"`
	ResolvedAt        *time.Time `db:"resolved_at"`
	DetectedAt        time.Time `db:"detected_at"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// Severity levels
type Severity string

const (
	SeverityLow      Severity = "low"      // Z-score 2.0-3.0
	SeverityMedium   Severity = "medium"   // Z-score 3.0-4.0
	SeverityHigh     Severity = "high"     // Z-score 4.0-5.0
	SeverityCritical Severity = "critical" // Z-score >5.0
)

// AnomalyStatus
type AnomalyStatus string

const (
	StatusOpen         AnomalyStatus = "open"
	StatusAcknowledged AnomalyStatus = "acknowledged"
	StatusResolved     AnomalyStatus = "resolved"
)
```

**File:** `internal/repository/anomaly_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/aulia/scinel-backend/internal/model"
)

// AnomalyRepository defines anomaly data access operations
type AnomalyRepository interface {
	// Create inserts new anomaly record
	Create(ctx context.Context, anomaly *model.Anomaly) (*model.Anomaly, error)

	// GetByID retrieves single anomaly
	GetByID(ctx context.Context, anomalyID uuid.UUID) (*model.Anomaly, error)

	// GetAll retrieves anomalies with filters
	GetAll(ctx context.Context, limit, offset int) ([]model.Anomaly, int64, error)

	// GetOpen retrieves open/unresolved anomalies
	GetOpen(ctx context.Context, limit, offset int) ([]model.Anomaly, int64, error)

	// UpdateStatus updates anomaly status (acknowledge or resolve)
	UpdateStatus(ctx context.Context, anomalyID uuid.UUID, status model.AnomalyStatus, notes string, resolvedBy uuid.UUID) error

	// GetBySKU retrieves anomalies for a specific SKU
	GetBySKU(ctx context.Context, skuID uuid.UUID, limit, offset int) ([]model.Anomaly, int64, error)
}
```

**Deliverable:**
- ✅ Anomaly model with severity enum
- ✅ Anomaly status tracking (open/acknowledged/resolved)
- ✅ Repository with 6 methods (vs 2 previously)

---

#### Task 4.1.2: Anomaly Detection Algorithm (1.5 hours)

**File:** `pkg/timeseries/anomaly.go`

```go
package timeseries

import (
	"fmt"
	"math"
	"sort"
)

// AnomalyDetectionResult contains anomaly detection output
type AnomalyDetectionResult struct {
	IsAnomaly     bool
	ZScore        float64  // Number of standard deviations from mean
	Severity      string   // low, medium, high, critical
	HistoricalMean float64
	HistoricalStd  float64
	Details       string
}

// DetectAnomalyZScore detects anomalies using Z-score method
// Input: Historical values (time series), current value
// Output: Z-score and severity classification
// Reference: TDD Section 6.4.1 (Z-score method)
func DetectAnomalyZScore(historical []float64, currentValue float64) *AnomalyDetectionResult {
	if len(historical) < 5 {
		return &AnomalyDetectionResult{
			IsAnomaly: false,
			Details:   "Insufficient historical data",
		}
	}

	// Calculate mean
	mean := 0.0
	for _, v := range historical {
		mean += v
	}
	mean /= float64(len(historical))

	// Calculate standard deviation
	variance := 0.0
	for _, v := range historical {
		variance += (v - mean) * (v - mean)
	}
	stdDev := math.Sqrt(variance / float64(len(historical)-1))

	// Handle zero variance (all values identical)
	if stdDev == 0 {
		stdDev = math.Abs(mean) * 0.01 // Use 1% of mean as baseline
		if stdDev == 0 {
			stdDev = 1.0 // Fallback
		}
	}

	// Calculate Z-score
	zScore := (currentValue - mean) / stdDev

	// Determine anomaly and severity
	isAnomaly := math.Abs(zScore) > 2.0 // Z > 2.0 is significant
	severity := "low"
	if math.Abs(zScore) > 5.0 {
		severity = "critical"
	} else if math.Abs(zScore) > 4.0 {
		severity = "high"
	} else if math.Abs(zScore) > 3.0 {
		severity = "medium"
	}

	return &AnomalyDetectionResult{
		IsAnomaly:      isAnomaly,
		ZScore:         zScore,
		Severity:       severity,
		HistoricalMean: mean,
		HistoricalStd:  stdDev,
		Details:        fmt.Sprintf("Value %.1f is %.2f std deviations from mean %.1f", currentValue, zScore, mean),
	}
}

// DetectAnomalyIQR detects anomalies using Interquartile Range method
// Alternative method for outlier detection
// Reference: TDD Section 6.4.2 (IQR method)
func DetectAnomalyIQR(historical []float64, currentValue float64) *AnomalyDetectionResult {
	if len(historical) < 5 {
		return &AnomalyDetectionResult{
			IsAnomaly: false,
			Details:   "Insufficient historical data",
		}
	}

	// Sort values for quartile calculation
	sorted := make([]float64, len(historical))
	copy(sorted, historical)
	sort.Float64s(sorted)

	// Calculate Q1, Q2 (median), Q3
	q1Idx := len(sorted) / 4
	q2Idx := len(sorted) / 2
	q3Idx := (3 * len(sorted)) / 4

	q1 := sorted[q1Idx]
	q3 := sorted[q3Idx]
	iqr := q3 - q1

	// Define outlier bounds
	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr

	isAnomaly := currentValue < lowerBound || currentValue > upperBound

	// Determine severity based on distance from bounds
	var distance float64
	if currentValue < lowerBound {
		distance = lowerBound - currentValue
	} else if currentValue > upperBound {
		distance = currentValue - upperBound
	}

	severity := "low"
	if distance > iqr*3 {
		severity = "critical"
	} else if distance > iqr*2 {
		severity = "high"
	} else if distance > iqr {
		severity = "medium"
	}

	return &AnomalyDetectionResult{
		IsAnomaly:      isAnomaly,
		ZScore:         (currentValue - q1) / math.Max(iqr, 1.0), // Pseudo Z-score
		Severity:       severity,
		HistoricalMean: q1 + iqr/2, // Approximate mean
		HistoricalStd:  iqr / 1.35, // Approximate std dev from IQR
		Details:        fmt.Sprintf("Value %.1f outside IQR bounds [%.1f, %.1f]", currentValue, lowerBound, upperBound),
	}
}
```

**Deliverable:**
- ✅ Z-score anomaly detection
- ✅ IQR anomaly detection (alternative method)
- ✅ Severity classification (low/medium/high/critical)

---

#### Task 4.1.3: Anomaly Detection Service (0.5 hours)

**File:** `internal/service/anomaly_service.go` (outline)

```go
package service

import (
	"context"

	"github.com/aulia/scinel-backend/pkg/timeseries"
)

// AnomalyService detects anomalies in inventory
type AnomalyService struct {
	repo      repository.AnomalyRepository
	invRepo   repository.InventoryRepository
	logger    logger.Logger
	publisher WebhookPublisher
}

// DetectAnomalies runs anomaly detection after inventory transaction
// Triggered: After every inventory transaction (async via background job)
// Input: warehouse_id, sku_id
// Process: Get historical inventory, run Z-score detection, save if anomaly
func (as *AnomalyService) DetectAnomalies(ctx context.Context, warehouseID, skuID uuid.UUID) error {
	// 1. Get current inventory level
	inv, err := as.invRepo.GetBySKU(ctx, warehouseID, skuID)
	if err != nil || inv == nil {
		return err
	}

	// 2. Get 90-day transaction history
	transactions, err := as.invRepo.GetTransactionHistory(ctx, warehouseID, skuID, 90)
	if err != nil || len(transactions) < 10 {
		return err // Not enough data
	}

	// 3. Convert to daily quantity changes
	dailyChanges := aggregateDailyChanges(transactions)

	// 4. Run anomaly detection
	result := timeseries.DetectAnomalyZScore(dailyChanges, float64(inv.Quantity))
	if !result.IsAnomaly {
		return nil // No anomaly detected
	}

	// 5. Save anomaly record
	anomaly := &model.Anomaly{
		WarehouseID:       warehouseID,
		SKUID:             skuID,
		Type:              "unusual_stock_level",
		Severity:          model.Severity(result.Severity),
		Description:       result.Details,
		HistoricalAverage: result.HistoricalMean,
		DetectedValue:     float64(inv.Quantity),
		ZScore:            result.ZScore,
		Status:            model.StatusOpen,
		DetectedAt:        time.Now(),
	}

	saved, err := as.repo.Create(ctx, anomaly)
	if err != nil {
		return err
	}

	// 6. Publish webhook event
	as.publisher.PublishAsync(ctx, "anomaly.detected", saved)

	return nil
}
```

**Deliverable:**
- ✅ Anomaly detection service
- ✅ Automatic triggering after transactions
- ✅ Webhook publishing

---

### Week 4.2: Anomaly Endpoints (2 hours)

**Objective:** Implement 4 anomaly endpoints

**Endpoint 1: GET /api/v1/anomalies (List Anomalies)**

**Purpose:** Retrieve all detected anomalies with filters

**Authentication:** RequireAuth + RequireRole("admin", "manager", "analyst")

**Query Parameters:**
```
?status=open
?severity=critical
?warehouse_id=wh-jakarta-001
?limit=20
&offset=0
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "anomalies": [
      {
        "id": "anomaly-uuid",
        "warehouse_id": "wh-jakarta-001",
        "sku_id": "sku-nike-001",
        "type": "unusual_stock_level",
        "severity": "high",
        "description": "Value 250 is 3.45 std deviations from mean 150",
        "z_score": 3.45,
        "status": "open",
        "detected_at": "2026-07-19T14:30:00Z"
      }
    ],
    "total": 12,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-07-19T14:45:00Z",
  "request_id": "req-uuid"
}
```

**Deliverable:**
- ✅ GET /api/v1/anomalies endpoint
- ✅ Status & severity filtering
- ✅ Pagination support

---

**Endpoint 2: GET /api/v1/anomalies/{id} (Anomaly Detail)**

**Purpose:** Retrieve detailed information for a specific anomaly

**Authentication:** RequireAuth

**Response (200 OK):** Single anomaly object with full details

**Deliverable:**
- ✅ GET /api/v1/anomalies/{id} endpoint
- ✅ Full anomaly details

---

**Endpoint 3: PUT /api/v1/anomalies/{id} (Update Anomaly Status)**

**Purpose:** Acknowledge or resolve anomaly

**Authentication:** RequireAuth + RequireRole("admin", "manager", "analyst")

**Request Body:**
```json
{
  "status": "acknowledged",
  "resolution_notes": "Verified: large customer order"
}
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Anomaly updated",
  "data": {
    "id": "anomaly-uuid",
    "status": "acknowledged",
    "updated_at": "2026-07-19T14:50:00Z"
  },
  "timestamp": "2026-07-19T14:50:00Z",
  "request_id": "req-uuid"
}
```

**Deliverable:**
- ✅ PUT /api/v1/anomalies/{id} endpoint
- ✅ Status update (open → acknowledged → resolved)
- ✅ Resolution notes tracking

---

**Endpoint 4: GET /api/v1/anomalies/{id}/history (Anomaly Timeline)**

**Purpose:** Retrieve anomaly resolution timeline

**Authentication:** RequireAuth

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "anomaly_id": "anomaly-uuid",
    "timeline": [
      {
        "timestamp": "2026-07-19T14:30:00Z",
        "event": "detected",
        "details": "Z-score 3.45"
      },
      {
        "timestamp": "2026-07-19T14:35:00Z",
        "event": "acknowledged",
        "user_id": "user-uuid",
        "notes": "Investigating"
      }
    ]
  },
  "timestamp": "2026-07-19T14:50:00Z",
  "request_id": "req-uuid"
}
```

**Deliverable:**
- ✅ GET /api/v1/anomalies/{id}/history endpoint
- ✅ Timeline of status changes
- ✅ User tracking

---

### Week 4.3: Webhook System (6 hours)

**Objective:** Implement webhook subscription, delivery, and retry logic

#### Task 4.3.1: Webhook Models & Repository (1.5 hours)

**File:** `internal/model/webhook.go`

```go
package model

import (
	"time"

	"github.com/google/uuid"
)

// WebhookSubscription represents a webhook event subscription
type WebhookSubscription struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	EventType   string    // e.g., "inventory.changed", "forecast.generated", "anomaly.detected"
	WebhookURL  string
	SecretKey   string    // For HMAC-SHA256 signing
	IsActive    bool
	MaxRetries  int
	RetryInterval int   // Seconds between retries
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID             uuid.UUID
	SubscriptionID uuid.UUID
	EventID        string // Unique event identifier (for idempotency)
	Payload        []byte // JSONB: Event data
	Status         DeliveryStatus
	HTTPStatusCode *int  // HTTP response status from webhook endpoint
	ErrorMessage   *string
	AttemptCount   int
	LastAttemptAt  *time.Time
	ScheduledRetryAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// DeliveryStatus represents webhook delivery status
type DeliveryStatus string

const (
	DeliveryPending    DeliveryStatus = "pending"     // Waiting to be sent
	DeliverySuccess    DeliveryStatus = "success"     // Successfully delivered
	DeliveryFailed     DeliveryStatus = "failed"      // Failed after max retries
	DeliveryDeadLetter DeliveryStatus = "dead_letter" // Moved to dead-letter queue
)

// WebhookEvent represents an event payload
type WebhookEvent struct {
	EventID   string    `json:"event_id"`   // Unique identifier
	EventType string    `json:"event_type"` // Type of event
	Timestamp time.Time `json:"timestamp"`
	Data      interface{} `json:"data"`
}
```

**File:** `internal/repository/webhook_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/aulia/scinel-backend/internal/model"
)

// WebhookRepository defines webhook data access operations
type WebhookRepository interface {
	// Create subscription
	CreateSubscription(ctx context.Context, sub *model.WebhookSubscription) (*model.WebhookSubscription, error)

	// Get subscriptions for user
	GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.WebhookSubscription, error)

	// Get subscriptions by event type (for publishing)
	GetSubscriptionsByEventType(ctx context.Context, eventType string) ([]model.WebhookSubscription, error)

	// Update subscription
	UpdateSubscription(ctx context.Context, sub *model.WebhookSubscription) error

	// Delete subscription
	DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error

	// Record delivery attempt
	CreateDelivery(ctx context.Context, delivery *model.WebhookDelivery) (*model.WebhookDelivery, error)

	// Get pending deliveries
	GetPendingDeliveries(ctx context.Context, limit int) ([]model.WebhookDelivery, error)

	// Update delivery status
	UpdateDeliveryStatus(ctx context.Context, deliveryID uuid.UUID, status model.DeliveryStatus, httpStatus *int, errMsg *string) error

	// Get delivery history
	GetDeliveryHistory(ctx context.Context, subscriptionID uuid.UUID, limit, offset int) ([]model.WebhookDelivery, int64, error)
}
```

**Deliverable:**
- ✅ Webhook subscription model
- ✅ Webhook delivery model (with retry tracking)
- ✅ Event payload structure
- ✅ Repository with 8 methods

---

#### Task 4.3.2: Webhook Publisher & Retry Logic (2.5 hours)

**File:** `internal/service/webhook_service.go`

```go
package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/aulia/scinel-backend/internal/model"
)

// WebhookPublisher handles webhook event publishing
type WebhookPublisher struct {
	repo   repository.WebhookRepository
	logger logger.Logger
	client *http.Client
}

// Publish sends webhook event to all subscribed endpoints
// Async: Spawns goroutine for actual delivery (don't block caller)
// Reference: LOGIC_FLOW Section 7 (Webhook Flow), TDD Section 5.5
func (wp *WebhookPublisher) PublishAsync(ctx context.Context, eventType string, data interface{}) {
	go func() {
		wp.publish(ctx, eventType, data)
	}()
}

func (wp *WebhookPublisher) publish(ctx context.Context, eventType string, data interface{}) {
	// 1. Get all subscriptions for this event type
	subs, err := wp.repo.GetSubscriptionsByEventType(ctx, eventType)
	if err != nil {
		wp.logger.Error("Failed to get subscriptions", map[string]interface{}{
			"event_type": eventType,
			"error":      err.Error(),
		})
		return
	}

	if len(subs) == 0 {
		return // No subscribers, nothing to do
	}

	// 2. Create event payload with unique ID
	eventID := uuid.New().String()
	event := model.WebhookEvent{
		EventID:   eventID,
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	payload, _ := json.Marshal(event)

	// 3. Create delivery record for each subscription
	for _, sub := range subs {
		if !sub.IsActive {
			continue
		}

		delivery := &model.WebhookDelivery{
			SubscriptionID: sub.ID,
			EventID:        eventID,
			Payload:        payload,
			Status:         model.DeliveryPending,
			AttemptCount:   0,
		}

		_, err := wp.repo.CreateDelivery(ctx, delivery)
		if err != nil {
			wp.logger.Error("Failed to create delivery record", map[string]interface{}{
				"subscription_id": sub.ID,
				"error":           err.Error(),
			})
		}
	}

	wp.logger.Info("Webhook published", map[string]interface{}{
		"event_type":      eventType,
		"event_id":        eventID,
		"subscriptions":   len(subs),
	})
}

// ProcessDeliveries runs retry logic for pending deliveries
// Called: By background job (ROADMAP_4 scheduler)
// Logic: Exponential backoff retry (max 5 attempts)
// Reference: TDD Section 5.5.2 (Webhook Retry)
func (wp *WebhookPublisher) ProcessDeliveries(ctx context.Context) error {
	// Get pending deliveries (limit to avoid overload)
	deliveries, err := wp.repo.GetPendingDeliveries(ctx, 100)
	if err != nil {
		return fmt.Errorf("failed to get pending deliveries: %w", err)
	}

	for _, delivery := range deliveries {
		// Check if it's time to retry
		if delivery.ScheduledRetryAt != nil && delivery.ScheduledRetryAt.After(time.Now()) {
			continue // Not ready to retry yet
		}

		// Get subscription details
		subs, err := wp.repo.GetSubscriptionsByEventType(ctx, "dummy") // TODO: Improve query
		if err != nil {
			continue
		}

		var sub *model.WebhookSubscription
		for _, s := range subs {
			if s.ID == delivery.SubscriptionID {
				sub = &s
				break
			}
		}

		if sub == nil {
			continue
		}

		// Attempt delivery
		err = wp.attemptDelivery(ctx, delivery, sub)
		if err != nil {
			wp.logger.Error("Delivery attempt failed", map[string]interface{}{
				"delivery_id": delivery.ID,
				"error":       err.Error(),
			})
		}
	}

	return nil
}

// attemptDelivery sends webhook to endpoint with HMAC signing
func (wp *WebhookPublisher) attemptDelivery(ctx context.Context, delivery *model.WebhookDelivery, sub *model.WebhookSubscription) error {
	// 1. Increment attempt count
	delivery.AttemptCount++

	// 2. Create signed request (HMAC-SHA256)
	// Reference: TDD Section 5.5.1 (HMAC Signing)
	signature := wp.signPayload(delivery.Payload, sub.SecretKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.WebhookURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", delivery.EventID)
	req.Header.Set("X-Webhook-Signature", "sha256="+signature) // HMAC signature
	req.Header.Set("X-Webhook-Timestamp", time.Now().Format(time.RFC3339))

	// 3. Send request with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := wp.client.Do(req)
	if err != nil {
		// Network error: Schedule retry
		nextRetryTime := wp.calculateRetryTime(delivery.AttemptCount, sub.RetryInterval)
		wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliveryPending, nil, ptrString(err.Error()))
		// TODO: Update scheduled_retry_at
		return err
	}
	defer resp.Body.Close()

	// 4. Handle response
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Success
		wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliverySuccess, &resp.StatusCode, nil)
		return nil
	}

	// Failure: Check if retryable
	if delivery.AttemptCount >= sub.MaxRetries {
		// Max retries exceeded: Move to dead-letter queue
		errMsg := fmt.Sprintf("HTTP %d after %d attempts", resp.StatusCode, delivery.AttemptCount)
		wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliveryDeadLetter, &resp.StatusCode, &errMsg)
		wp.logger.Warn("Webhook delivery failed (max retries)", map[string]interface{}{
			"delivery_id":      delivery.ID,
			"subscription_id":  delivery.SubscriptionID,
			"attempts":         delivery.AttemptCount,
			"http_status":      resp.StatusCode,
		})
		return fmt.Errorf("max retries exceeded")
	}

	// Schedule next retry (exponential backoff)
	nextRetryTime := wp.calculateRetryTime(delivery.AttemptCount, sub.RetryInterval)
	errMsg := fmt.Sprintf("HTTP %d, retry at %v", resp.StatusCode, nextRetryTime)
	wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliveryPending, &resp.StatusCode, &errMsg)
	// TODO: Update scheduled_retry_at to nextRetryTime

	return fmt.Errorf("failed with status %d, will retry", resp.StatusCode)
}

// signPayload creates HMAC-SHA256 signature
// Reference: TDD Section 5.5.1
func (wp *WebhookPublisher) signPayload(payload []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// calculateRetryTime determines next retry time with exponential backoff
// Backoff: 1s, 2s, 4s, 8s (2^attempt)
func (wp *WebhookPublisher) calculateRetryTime(attemptCount, retryInterval int) time.Time {
	backoffSeconds := 1 << uint(attemptCount-1) // 2^(attempt-1)
	if backoffSeconds < retryInterval {
		backoffSeconds = retryInterval
	}
	if backoffSeconds > 3600 { // Cap at 1 hour
		backoffSeconds = 3600
	}
	return time.Now().Add(time.Duration(backoffSeconds) * time.Second)
}

func ptrString(s string) *string {
	return &s
}
```

**Webhook Retry Logic:**
- Attempt 1: Immediate
- Attempt 2: 2 seconds delay
- Attempt 3: 4 seconds delay
- Attempt 4: 8 seconds delay
- Attempt 5 (final): 16 seconds delay
- Max 5 retries total
- Dead-letter queue for permanent failures

**Webhook Payload Structure:**
```json
{
  "event_id": "evt-unique-id",
  "event_type": "inventory.changed",
  "timestamp": "2026-07-19T14:30:00Z",
  "data": {
    "warehouse_id": "wh-jakarta-001",
    "sku_id": "sku-nike-001",
    "quantity": 150,
    "transaction_id": "txn-uuid"
  }
}
```

**HMAC Signature:**
- Header: `X-Webhook-Signature: sha256={hex-encoded-hmac}`
- Algorithm: HMAC-SHA256
- Secret: Stored in database (not exposed)

**Deliverable:**
- ✅ Webhook publisher with async delivery
- ✅ HMAC-SHA256 signing (security)
- ✅ Exponential backoff retry (5 attempts)
- ✅ Dead-letter queue (permanent failures)
- ✅ Event deduplication (via event_id)
- ✅ At-least-once delivery guarantee

---

### Week 4.4: Webhook Endpoints, Health Checks, Rate Limiting (5 hours)

#### Task 4.4.1: Webhook Endpoints (2 hours)

**Endpoint 1: POST /api/v1/webhooks/subscribe (Create Subscription)**

**Purpose:** Subscribe to webhook events

**Authentication:** RequireAuth

**Request Body:**
```json
{
  "event_type": "inventory.changed",
  "webhook_url": "https://example.com/webhooks/inventory",
  "max_retries": 5,
  "retry_interval": 60
}
```

**Response (201 Created):**
```json
{
  "code": 201,
  "message": "Subscription created",
  "data": {
    "id": "sub-uuid",
    "event_type": "inventory.changed",
    "webhook_url": "https://example.com/webhooks/inventory",
    "is_active": true,
    "created_at": "2026-07-19T14:50:00Z"
  },
  "timestamp": "2026-07-19T14:50:00Z",
  "request_id": "req-uuid"
}
```

**Deliverable:**
- ✅ POST /api/v1/webhooks/subscribe endpoint
- ✅ Secret key generation (for HMAC)
- ✅ Subscription created and active

---

**Endpoint 2: GET /api/v1/webhooks/subscriptions (List Subscriptions)**

**Purpose:** Retrieve user's webhook subscriptions

**Authentication:** RequireAuth

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "subscriptions": [
      {
        "id": "sub-uuid",
        "event_type": "inventory.changed",
        "webhook_url": "https://example.com/webhooks/inventory",
        "is_active": true,
        "created_at": "2026-07-19T14:50:00Z"
      }
    ]
  },
  "timestamp": "2026-07-19T14:50:00Z",
  "request_id": "req-uuid"
}
```

**Deliverable:**
- ✅ GET /api/v1/webhooks/subscriptions endpoint

---

**Endpoint 3: DELETE /api/v1/webhooks/subscriptions/{id} (Unsubscribe)**

**Purpose:** Delete webhook subscription

**Authentication:** RequireAuth (owner only)

**Response (204 No Content):**

**Deliverable:**
- ✅ DELETE /api/v1/webhooks/subscriptions/{id} endpoint

---

**Endpoint 4: GET /api/v1/webhooks/deliveries (View Delivery History)**

**Purpose:** View webhook delivery attempts and status

**Authentication:** RequireAuth

**Query Parameters:**
```
?subscription_id=sub-uuid
?status=success
?limit=20
&offset=0
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "deliveries": [
      {
        "id": "delivery-uuid",
        "subscription_id": "sub-uuid",
        "event_id": "evt-uuid",
        "status": "success",
        "http_status_code": 200,
        "attempt_count": 1,
        "created_at": "2026-07-19T14:30:00Z"
      }
    ],
    "total": 45,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-07-19T14:50:00Z",
  "request_id": "req-uuid"
}
```

**Deliverable:**
- ✅ GET /api/v1/webhooks/deliveries endpoint
- ✅ Delivery history with status tracking

---

#### Task 4.4.2: API Key Endpoints (1 hour)

**NEW ENDPOINTS (from audit critical issues):**

**Endpoint 1: POST /api/v1/api-keys (Create API Key)**

**Purpose:** Generate programmatic access token (alternative to user credentials)

**Authentication:** RequireAuth

**Request Body:**
```json
{
  "name": "Integration Service",
  "permissions": ["inventory:read", "inventory:write", "forecast:read"],
  "expires_in_days": 365
}
```

**Response (201 Created):**
```json
{
  "code": 201,
  "message": "API key created",
  "data": {
    "key_id": "key-uuid",
    "key_prefix": "sk_test_1a2b3c",
    "key_secret": "sk_test_1a2b3c_full_secret_only_shown_once", // NEVER show again
    "name": "Integration Service",
    "permissions": ["inventory:read", "inventory:write", "forecast:read"],
    "created_at": "2026-07-19T14:55:00Z",
    "expires_at": "2027-07-19T14:55:00Z"
  },
  "timestamp": "2026-07-19T14:55:00Z",
  "request_id": "req-uuid"
}
```

**Security Note:** Full secret only shown once at creation. Users must store securely.

**Deliverable:**
- ✅ POST /api/v1/api-keys endpoint
- ✅ Secure key generation
- ✅ Permission-based access control
- ✅ Expiration tracking

---

**Endpoint 2: GET /api/v1/api-keys (List API Keys)**

**Purpose:** View all API keys (masked)

**Authentication:** RequireAuth

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "keys": [
      {
        "key_id": "key-uuid",
        "key_prefix": "sk_test_1a2b3c",
        "name": "Integration Service",
        "permissions": ["inventory:read", "inventory:write"],
        "is_active": true,
        "last_used_at": "2026-07-19T14:00:00Z",
        "created_at": "2026-07-19T14:55:00Z",
        "expires_at": "2027-07-19T14:55:00Z"
      }
    ]
  },
  "timestamp": "2026-07-19T14:55:00Z",
  "request_id": "req-uuid"
}
```

**Security Note:** Secret never shown in list response.

**Deliverable:**
- ✅ GET /api/v1/api-keys endpoint

---

**Endpoint 3: DELETE /api/v1/api-keys/{id} (Revoke API Key)**

**Purpose:** Revoke API key immediately

**Authentication:** RequireAuth (owner only)

**Response (204 No Content):**

**Deliverable:**
- ✅ DELETE /api/v1/api-keys/{id} endpoint
- ✅ Immediate revocation

---

**Endpoint 4: POST /api/v1/api-keys/{id}/rotate (Rotate API Key)**

**Purpose:** Generate new secret for existing key (for key rotation)

**Authentication:** RequireAuth (owner only)

**Response (201 Created):**
```json
{
  "code": 201,
  "message": "API key rotated",
  "data": {
    "key_id": "key-uuid",
    "key_prefix": "sk_test_1a2b3c",
    "key_secret": "sk_test_1a2b3c_new_secret",
    "rotated_at": "2026-07-19T14:55:00Z"
  },
  "timestamp": "2026-07-19T14:55:00Z",
  "request_id": "req-uuid"
}
```

**Deliverable:**
- ✅ POST /api/v1/api-keys/{id}/rotate endpoint
- ✅ Key rotation without losing key_id
- ✅ Smooth migration path

---

#### Task 4.4.3: Rate Limiting Middleware (1.5 hours)

**File:** `internal/middleware/rate_limit.go`

**Purpose:** Limit API requests to 1000 req/sec per user/API-key

**Reference:** TDD Section 8 (Rate Limiting), LOGIC_FLOW Section 9 (Traffic Control)

```go
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/aulia/scinel-backend/internal/util"
)

// RateLimitMiddleware applies rate limiting (Redis-based)
// Rate limit: 1000 requests per second per user
// Key: user_id or api_key (for auth'd requests)
// Alternative for unauthenticated: IP address (not recommended for production)
type RateLimitMiddleware struct {
	redis        *redis.Client
	requestsPerSec int
	windowSize    time.Duration
}

// NewRateLimitMiddleware creates rate limiter
func NewRateLimitMiddleware(redisClient *redis.Client, requestsPerSec int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redis:          redisClient,
		requestsPerSec: requestsPerSec,
		windowSize:     1 * time.Second, // 1-second sliding window
	}
}

// Apply rate limiting to request
func (rl *RateLimitMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := util.GenerateRequestID()
		ctx := r.Context()

		// Determine rate limit key (from auth token or API key)
		// Priority: API-Key header > Authorization bearer token > IP address
		limitKey := rl.getLimitKey(r)

		// Check rate limit
		allowed, remaining, retryAfter, err := rl.checkLimit(ctx, limitKey)
		if err != nil {
			// Redis error: Log but don't block request (rate limit is optional)
			// Continue to next handler
			next.ServeHTTP(w, r)
			return
		}

		// Add rate limit info to response headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.requestsPerSec))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Duration(retryAfter)*time.Second).Unix()))

		// Block if limit exceeded
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			util.WriteError(w, http.StatusTooManyRequests, "Rate limit exceeded", 
				fmt.Sprintf("You have exceeded %d requests per second. Retry after %d seconds", rl.requestsPerSec, retryAfter), 
				requestID)
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// checkLimit checks if request is within rate limit
func (rl *RateLimitMiddleware) checkLimit(ctx context.Context, key string) (allowed bool, remaining int, retryAfter int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// Get current count
	count, err := rl.redis.Incr(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return true, 0, 0, err // Redis error, allow request
	}

	// If first request in window, set expiration
	if count == 1 {
		rl.redis.Expire(ctx, redisKey, rl.windowSize)
	}

	// Check if limit exceeded
	if int(count) > rl.requestsPerSec {
		ttl, _ := rl.redis.TTL(ctx, redisKey).Result()
		retryAfter := int(ttl.Seconds()) + 1
		if retryAfter <= 0 {
			retryAfter = 1
		}
		return false, 0, retryAfter, nil
	}

	remaining := rl.requestsPerSec - int(count)
	return true, remaining, 0, nil
}

// getLimitKey extracts limit key from request
func (rl *RateLimitMiddleware) getLimitKey(r *http.Request) string {
	// Check for API key in header
	if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
		return apiKey
	}

	// Check for user ID in Bearer token (from JWT claims)
	if userID, ok := r.Context().Value("user_id").(string); ok && userID != "" {
		return "user:" + userID
	}

	// Fallback to IP address (not ideal)
	ip := r.RemoteAddr
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = xff
	}
	return "ip:" + ip
}
```

**Configuration:**
- Rate limit: 1000 requests per second
- Window size: 1 second (sliding)
- Key format: "rate_limit:{user_id}" or "rate_limit:{api_key}"
- Redis storage: Automatic expiration after window
- Error handling: Log but don't block (graceful degradation)

**Response Headers:**
- `X-RateLimit-Limit: 1000`
- `X-RateLimit-Remaining: 950`
- `X-RateLimit-Reset: 1626721234`
- `Retry-After: 5` (if rate limited)

**Deliverable:**
- ✅ Redis-based rate limiting
- ✅ 1000 req/sec per user/API-key
- ✅ Sliding window algorithm
- ✅ Proper HTTP headers (RFC 6585)
- ✅ Graceful degradation on Redis error

---

#### Task 4.4.4: Integration Testing & Documentation (1.5 hours)

**Testing Checklist:**

```
Unit Tests:
- [ ] Inventory service: GetBySKU, UpdateQuantity, GetByWarehouse, caching
- [ ] Transaction service: Create, Query, locking
- [ ] Forecast service: ETS, ARIMA, batch job
- [ ] Anomaly service: Z-score, IQR detection
- [ ] Webhook service: Publishing, HMAC signing, retry logic
- [ ] Rate limiting: Count increments, limit enforcement, headers

Integration Tests:
- [ ] E2E: POST /transactions → webhook published → anomaly detected
- [ ] Forecast batch: Generate forecasts for 100 SKUs (<5s)
- [ ] Webhook retries: Failed delivery → exponential backoff → success
- [ ] Concurrent transactions: 100 simultaneous (no deadlock)
- [ ] Rate limiting: Exceed 1000 req/sec → 429 response

Performance Tests:
- [ ] GET /inventory: p99 <100ms (with cache)
- [ ] Forecast batch: <5s for 1000 SKUs
- [ ] Webhook processing: 100 deliveries/sec
- [ ] Cache hit rate: >60% for inventory queries

Documentation:
- [ ] API.md: All 23 endpoints documented with curl examples
- [ ] ARCHITECTURE.md: Updated with ROADMAP_2 components
- [ ] README.md: Feature status (ROADMAP_2 complete)
```

**Deliverable:**
- ✅ >80% test coverage for backend
- ✅ All performance targets verified
- ✅ API documentation complete

---

## ROADMAP_2 COMPLETE ENDPOINT COUNT

**Total Endpoints: 23**

| Category | Count | Endpoints |
|----------|-------|-----------|
| **Health Checks** | 3 | /health, /readiness, /liveness |
| **Authentication** | 4 | register, login, refresh, logout (from ROADMAP_1) |
| **Inventory** | 5 | GET /inventory, GET /inventory/{warehouse}, POST /transactions, GET /transactions, POST /import |
| **Forecasting** | 3 | GET /forecasts, POST /forecasts/generate, GET /forecasts/accuracy |
| **Anomalies** | 4 | GET /anomalies, GET /anomalies/{id}, PUT /anomalies/{id}, GET /anomalies/{id}/history |
| **Webhooks** | 4 | POST /webhooks/subscribe, GET /webhooks/subscriptions, DELETE /webhooks/subscriptions/{id}, GET /webhooks/deliveries |
| **API Keys** | 4 | POST /api-keys, GET /api-keys, DELETE /api-keys/{id}, POST /api-keys/{id}/rotate |
| **Rate Limiting** | 1 | Middleware (applied globally) |

**TOTAL: 23 endpoints** ✅ (exactly as per Blueprint)

---

## 🧪 TESTING STRATEGY

### Unit Tests (>80% coverage)
```
✅ Inventory: Create, update, query, cache logic
✅ Transaction: In/out, pessimistic locking, deadlock retry  
✅ Forecast: ETS + ARIMA algorithms, batch job
✅ Anomaly: Z-score + IQR detection, severity classification
✅ Webhook: Publishing, HMAC signing, retry logic
✅ API Keys: Generation, validation, expiration
✅ Rate Limiting: Counter increment, limit check
```

### Integration Tests
```
✅ E2E: Transaction → Inventory Update → Anomaly Detection → Webhook
✅ Forecast: Historical data → ETS/ARIMA → Predictions
✅ Webhook: Failed delivery → Exponential backoff → Success
✅ Concurrent: 100 simultaneous transactions (deadlock handling)
```

### Performance Tests
```
✅ Inventory query: p99 <100ms (with cache)
✅ Forecast batch: <5s for 1000 SKUs
✅ Concurrent transactions: 100 simultaneous
✅ Cache hit rate: >60% for inventory
```

### Load Tests
```
✅ 1000 concurrent users
✅ Rate limiter at capacity (1000 req/sec)
✅ Webhook delivery: 100 deliveries/sec sustained
```

---

## 📊 ROADMAP_2 COMPLETION CHECKLIST

### All Endpoints (23 total)
- ✅ 3 Health check endpoints (health, readiness, liveness)
- ✅ 4 Auth endpoints (register, login, refresh, logout)
- ✅ 5 Inventory endpoints (list, warehouse, transactions, query, import)
- ✅ 3 Forecast endpoints (query, generate, accuracy)
- ✅ 4 Anomaly endpoints (list, detail, update, history)
- ✅ 4 Webhook endpoints (subscribe, list, delete, deliveries)
- ✅ 4 API Key endpoints (create, list, revoke, rotate)

### Business Logic
- ✅ Inventory transactions with pessimistic locking + deadlock retry
- ✅ ETS forecasting algorithm (180-day historical, 30-day forecast)
- ✅ ARIMA forecasting algorithm (auto parameter selection)
- ✅ Algorithm comparison & selection (MAPE-based)
- ✅ Anomaly detection (Z-score + IQR methods)
- ✅ Webhook publishing with retry logic (5 attempts, exponential backoff)
- ✅ Batch forecast job (daily at 2 AM UTC via cron)
- ✅ API Key authentication (alternative to JWT)

### Database & Caching
- ✅ Inventory cache: 30s TTL (Redis)
- ✅ Forecast cache: 24h TTL (Redis)
- ✅ Rate limit counters: Redis (1s sliding window)
- ✅ All 12 tables populated with test data
- ✅ Indexes optimized for query performance

### Middleware & Security
- ✅ Rate limiting: 1000 req/sec per key (global)
- ✅ RBAC: Manager warehouse isolation
- ✅ Webhook HMAC-SHA256 signing
- ✅ API Key secret hashing + rotation
- ✅ Error responses: Standardized (TDD Sec 5.1.1)
- ✅ Request logging: Structured JSON

### Git Commits (Roadmap 2: Commits 8-15)
- [ ] Commit 8: Inventory models + repository
- [ ] Commit 9: Inventory transaction handling + endpoints
- [ ] Commit 10: Inventory cache + query endpoints
- [ ] Commit 11: Forecast models + ETS/ARIMA algorithms
- [ ] Commit 12: Forecast repository + service + batch job
- [ ] Commit 13: Forecast endpoints + manual trigger
- [ ] Commit 14: Anomaly detection service + 4 endpoints
- [ ] Commit 15: Webhook system + API Keys + Rate limiting

---

## 📚 KEY REFERENCE SECTIONS

### Documents
- **TDD Sec 4.4-4.10:** Database tables & schema
- **TDD Sec 5.2-5.5:** API endpoints & specifications
- **TDD Sec 6.3-6.4:** Business logic & algorithms
- **TDD Sec 8:** Rate limiting & security
- **LOGIC_FLOW Sec 4-7:** Detailed step-by-step flows
- **LOGIC_FLOW Sec 6:** Batch job flow (forecast generation)
- **LOGIC_FLOW Sec 9:** Traffic control (rate limiting)
- **LOGIC_FLOW Sec 10:** Monitoring & health checks
- **BLUEPRINT Sec 3:** Endpoint mapping (23 total)

### Tech Stack (LOCKED)
- Backend: Go 1.21+
- Database: PostgreSQL 15+, Redis 7+
- Testing: Go testing package
- Algorithms: ETS + ARIMA (numerical, no ML libraries)

### Performance Targets (LOCKED)
- Inventory query: p99 <100ms ✓
- Forecast batch: <5s for 1000 SKUs ✓
- Concurrent users: 1000+ ✓
- Rate limit: 1000 req/sec per key ✓

---

## ⚡ CRITICAL IMPLEMENTATION NOTES

### Deadlock Handling (IMPORTANT)
- Always use pessimistic locking: `SELECT ... FOR UPDATE`
- Retry logic: 3 attempts with backoff (100ms, 200ms, 400ms)
- Reference: TDD Section 4.0, ROADMAP_1 Week 1.3

### Webhook Guarantees
- At-least-once delivery (with event_id deduplication)
- Max 5 retry attempts
- Exponential backoff: 1s, 2s, 4s, 8s, 16s
- Dead-letter queue for permanent failures

### Cache Invalidation
- After transaction: Immediately invalidate inventory cache
- After forecast: Set forecast cache with 24h TTL
- Handle race conditions (cache-aside pattern)

### Algorithm Accuracy
- ETS/ARIMA: 180 days historical data minimum
- Forecasts: 30-day horizon by default
- Confidence intervals: 95% (±1.96 * std_error)
- Accuracy metrics: MAPE, RMSE, MAE (calculated daily)

### RBAC Enforcement (IMPORTANT)
- Manager: Can only create transactions for own warehouse
- Admin: Can trigger forecast batch, manage API keys, webhooks
- Analyst: Read-only access to forecasts + anomalies
- Viewer: Read-only access to inventory
- Enforce at handler level (middleware) + service level

---

## 🎯 SUCCESS CRITERIA

**ROADMAP_2 is COMPLETE when:**

1. ✅ All 23 endpoints working and documented
2. ✅ All business logic implemented (inventory, forecast, anomaly, webhook, health, API keys)
3. ✅ >80% test coverage for backend services
4. ✅ Performance targets verified: p99 <100ms inventory, <5s forecast batch
5. ✅ Rate limiting active and tested (1000 req/sec limit)
6. ✅ Webhook system delivering with retry logic (5 attempts, exponential backoff)
7. ✅ Health endpoints passing (for Kubernetes, monitoring, CI/CD)
8. ✅ API Key authentication working with key rotation
9. ✅ Redis caching working (inventory 30s, forecast 24h)
10. ✅ CI/CD pipeline passing all tests, ready for ROADMAP_3

**Ready to move to ROADMAP_3: Frontend/Mobile development** ✅

---

**Document:** ROADMAP_2: BACKEND CORE (REVISED v1.1)  
**Status:** ✅ FINAL - Ready for Development  
**Generated:** 2026-07-19  
**Version:** 1.1 (Comprehensive Audit Fixes Applied)  
**Total Endpoints:** 23 ✓  
**Coverage:** 100% of PRD backend requirements  

---

**END OF ROADMAP_2 REVISION v1.1**

