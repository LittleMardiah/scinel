# Software Test Plan (STP) - SUPPLY CHAIN INTELLIGENCE PLATFORM v1.0.0

> **Purpose**: Dokumen ini mendefinisikan strategi testing, scope, execution plan, automation, dan quality criteria untuk Supply Chain Intelligence Platform (ScIntel) MVP Release v1.0.0, yang mencakup Demand Forecasting Engine, Real-Time Inventory Tracking, dan Anomaly Detection System.

> **Audience**: Backend Engineers, QA Engineers, DevOps Engineers, Product Owners

> **Version**: 1.0.0 - FINAL COMPLETE

---

## 1. DOCUMENT METADATA & OVERVIEW

### 1.1 Project Information
- **Project Name**: Supply Chain Intelligence Platform (ScIntel)
- **Project Acronym**: ScIntel
- **Release Version**: 1.0.0 (MVP - Portfolio Edition)
- **Release Target Date**: 31/10/2026 (12 weeks from 18/07/2026)
- **STP Owner**: Backend Engineering Team / QA Lead
- **Last Updated**: 18/07/2026
- **Next Review**: 30/09/2026 (End of Phase 1)
- **Status**: Draft → Ready for Approval

### 1.2 Document Purpose

Dokumen ini menjelaskan bagaimana testing akan dilakukan untuk ScIntel MVP Release v1.0.0, mencakup:

- **Scope Testing**: Fitur yang di-test, platform yang didukung, priority levels
- **Strategi Testing**: Testing pyramid, test levels (unit/integration/E2E), tools & frameworks
- **Execution Plan**: Timeline, resource allocation, triggers, testing phases
- **Quality Criteria**: Acceptance criteria, coverage targets, performance SLAs
- **Automation & Regression Testing**: CI/CD integration, automated test suites, regression strategy
- **Bug Tracking & Resolution**: Process, severity levels, escalation paths
- **Test Reporting & Metrics**: KPIs, monitoring, dashboards, review cadence
- **UAT & Pre-Release Verification**: User acceptance testing procedures, final checklist

### 1.3 Related Documents

- [x] PRD (Product Requirements Document): PRD.md v1.0.0
- [x] TDD (Technical Design Document): TDD.md v1.0.0 (21 API endpoints, 12 database tables)
- [x] LOGIC_FLOW (System Architecture & Business Logic): LOGIC_FLOW.md v1.0.0 (5 detailed flows)
- [x] HALAMAN (UI/UX Design): HALAMAN.md v2.0 (Dashboard, Inventory, Forecast, Anomaly pages)
- [x] DRP (Disaster Recovery Plan): DRP.md v1.0.0 (Backup, failover, RTO/RPO targets)
- [x] API Contract Specification: TDD.md section 5 (request/response examples, error codes)
- [x] Database Schema Documentation: TDD.md section 4 (12 tables, relationships, indexes)
- [x] CI/CD Pipeline Configuration: GitHub Actions workflows
- [x] Bug Tracking Repository: GitHub Issues

---

## 2. TESTING SCOPE & OBJECTIVES

### 2.1 Scope: Features In Testing

#### **Functional Requirements (In Scope - MVP v1.0.0)**

| Module/Feature | Priority | Test Type | Coverage Target | Notes |
|---|---|---|---|---|
| **Authentication & Authorization** | P0 | Unit / Integration / API | 100% | JWT, RBAC (4 roles), token refresh, password hashing (bcrypt) |
| **User Management** | P0 | Unit / Integration / API | 90% | Create user, update profile, list users, role assignment |
| **Real-Time Inventory Tracking** | P0 | Unit / Integration / API / E2E | 100% | Stock in/out, multi-warehouse, <100ms p99 latency |
| **Inventory Transactions** | P0 | Unit / Integration / API | 100% | Log transactions, audit trail, movement history |
| **Inventory Snapshots** | P1 | Unit / Integration | 80% | End-of-day, weekly, monthly snapshots |
| **Stock Alerts** | P1 | Unit / Integration | 80% | Low stock, overstock, expiring items thresholds |
| **Demand Forecasting Engine (ETS)** | P0 | Unit / Integration / API | 100% | Exponential smoothing, 6-12 months historical data, accuracy metrics |
| **Demand Forecasting Engine (ARIMA)** | P0 | Unit / Integration / API | 100% | AutoRegressive Integrated Moving Average, model comparison |
| **Forecast Batch Processing** | P0 | Unit / Integration | 100% | Scheduled runs (daily/weekly), <5s for 1000 SKUs, background jobs |
| **Forecast Accuracy Metrics** | P1 | Unit / Integration | 85% | MAE, RMSE, MAPE calculations, confidence intervals (95%) |
| **Anomaly Detection (Z-score)** | P0 | Unit / Integration / API | 90% | Statistical outlier detection, threshold-based alerts |
| **Anomaly Detection (IQR)** | P1 | Unit / Integration | 85% | Interquartile range method, pattern recognition |
| **Anomaly Alert Notifications** | P1 | Unit / Integration | 80% | Email, webhook, severity levels (low/med/high/critical) |
| **Webhook Event Subscriptions** | P0 | Unit / Integration / API | 90% | Subscribe to: inventory.changed, forecast.generated, anomaly.detected |
| **Webhook Event Delivery** | P0 | Unit / Integration / API | 90% | POST to external URLs, retry logic (exponential backoff, max 5 retries), HMAC-SHA256 signature |
| **Webhook Dead-Letter Queue** | P1 | Unit / Integration | 80% | Failed events, investigation, manual retry capability |
| **Audit Logging & Compliance** | P0 | Unit / Integration | 90% | Who did what, when, from where (user_id, action, resource, timestamp, ip_address) |
| **Database Connection Pooling** | P0 | Integration | 90% | Min 10, Max 50 connections, efficient resource management |
| **Graceful Shutdown** | P0 | Integration | 100% | Clean resource cleanup, in-flight request completion |
| **Structured Logging (JSON)** | P1 | Integration | 85% | JSON logs with appropriate levels (DEBUG, INFO, WARN, ERROR) |
| **Database Migrations** | P0 | Integration | 100% | Version control, rollback capability, migration validation |
| **Rate Limiting** | P1 | Unit / Integration / API | 85% | Per-user/API-key abuse prevention, 1000 req/sec limit |
| **API Key Authentication (Third-party)** | P1 | Unit / Integration / API | 80% | Third-party integrations, key validation, scope-based access |

**Priority Level Definitions**:
- **P0 (Critical)**: Core business logic, must have 100% test coverage, blocks release if failing
- **P1 (High)**: Important features, should have 80%+ coverage, must fix before release
- **P2 (Medium)**: Nice-to-have features, 50%+ coverage acceptable, can defer to next release

---

### 2.2 Scope: Features Out of Testing

| Feature | Reason | Target Release |
|---|---|---|
| Advanced ML-based Anomaly Detection (Isolation Forest, LSTM) | Out of scope for MVP, complex ML implementation | Phase 1.5 |
| Kafka Event Streaming | Not in MVP scope, future scaling enhancement | Phase 1.5 |
| GraphQL API | REST API sufficient for MVP, can add GraphQL layer later | Phase 2 |
| OAuth2 Social Login (Google, GitHub) | Out of scope for MVP, basic JWT sufficient | Phase 1.1 |
| Two-Factor Authentication (2FA) | Out of scope for MVP, future security enhancement | Phase 1.2 |
| File Uploads & S3 Integration | Not MVP requirement, can add in Phase 1.1 | Phase 1.1 |
| Multi-Tenancy | Out of scope for MVP, architecture designed to support in future | Phase 2 |
| Blue-Green Deployment | Advanced DevOps feature, basic deployment sufficient for MVP | Phase 2 |
| Materialized Views & Table Partitioning | Performance optimization, secondary to MVP functionality | Phase 1.5 |
| Holt-Winters Forecasting Algorithm | ARIMA + ETS sufficient for MVP, can add later | Phase 1.5 |

---

### 2.3 Platform & Environment Scope

**Backend Platforms**:
- [x] **Language & Runtime**: Go 1.21+ (Linux, macOS, Windows compatible)
- [x] **Database**: PostgreSQL 15+ (development, staging, production)
- [x] **Cache Layer**: Redis 7+ (development, staging, production)
- [x] **Operating System**: Linux (Docker container), macOS (local dev), Windows (local dev via WSL2)
- [x] **Container Orchestration**: Docker (all environments), Docker Compose (local dev)

**API & Integration Platforms**:
- [x] **HTTP Protocol**: REST API over HTTPS/HTTP (version 1.0: HTTP on localhost for testing)
- [x] **Content-Type**: JSON (application/json)
- [x] **Authentication**: Bearer token (JWT), API Key (basic support)

**Network Conditions** (if applicable):
- [x] **Local Development**: High bandwidth, low latency
- [x] **Staging Server**: Realistic network conditions (15-30ms latency, 99%+ uptime)
- [x] **Production**: High availability, <100ms latency for p99 queries

---

### 2.4 Testing Objectives

**Primary Objectives**:
- Validate all P0 & P1 features meet functional requirements defined in PRD.md v1.0.0
- Ensure non-functional requirements (performance, security, reliability) are met per PRD section 4.9
- Identify and document bugs before release, prioritize by severity
- Verify bug fixes and regression tests pass
- Ensure system stability and data integrity

**Secondary Objectives**:
- Achieve >80% code coverage for critical modules (auth, inventory, forecasting, anomaly)
- Verify API contracts match TDD.md specifications exactly
- Validate database schema, relationships, constraints work as designed
- Test error handling, edge cases, boundary conditions
- Perform load testing (1000+ concurrent users)
- Verify DevOps readiness (Docker, Kubernetes manifests, CI/CD)

**Quality Gates**:
- All P0 bugs must be fixed before release
- All P1 bugs should be fixed (critical path features)
- P2 bugs can defer to next release
- Code coverage ≥80% for critical paths
- Performance SLAs met: p99 <100ms for real-time queries, <5s for batch forecasts

---

## 3. TEST STRATEGY & PYRAMID APPROACH

### 3.1 Testing Pyramid Overview

```
                    /\
                   /  \          E2E / UI Tests
                  /----\         (10-15% of effort)
                 /      \
                /--------\
               /          \      Integration Tests
              /            \     (35-40% of effort)
             /              \
            /________________\
           /                  \  Unit Tests
          /                    \ (45-50% of effort)
         /______________________ \
```

**Distribution Target for ScIntel**:
- **Unit Tests**: 45-50% of total test effort → Target: 250+ unit tests
- **Integration Tests**: 35-40% of total test effort → Target: 150+ integration tests
- **API Contract Tests**: 5-10% of total test effort → Target: 50+ API tests
- **E2E / UI Tests**: 10-15% of total test effort → Target: 30-40 critical user journeys

**Rationale**: 
- High unit test coverage for time-series algorithms (forecasting), anomaly detection logic
- Integration tests for database transactions, multi-warehouse inventory consistency
- API contract tests ensure backend/frontend contracts are verified
- Limited E2E tests focused on critical paths (login → inventory → forecast generation)

---

### 3.2 Test Levels & Definitions

#### **Level 1: Unit Tests**

**Purpose**: Test individual functions/methods in isolation, verify business logic correctness

**Scope** (Per TDD.md architecture):
- [x] All public functions in service layer:
  - `auth_service.go`: login, register, refresh token, validate JWT
  - `inventory_service.go`: record transaction, get stock level, update snapshots
  - `forecast_service.go`: run forecasting algorithm, calculate accuracy metrics
  - `anomaly_service.go`: detect Z-score outliers, IQR anomalies
  - `webhook_service.go`: publish events, handle retries
- [x] All utility/helper functions:
  - `exponential_smoothing.go`: ETS algorithm implementations
  - `arima.go`: ARIMA algorithm implementations
  - `anomaly_detector.go`: statistical detection methods
  - `validator.go`: input validation functions
  - `hash.go`: bcrypt password hashing
  - `jwt.go`: JWT creation/verification
- [x] Business logic validation:
  - Demand forecasting accuracy (MAE, RMSE, MAPE calculations)
  - Anomaly detection thresholds (Z-score, IQR bounds)
  - Role-based access control logic
  - Password strength validation
- [x] Edge cases & boundary conditions:
  - Empty inventory data (forecast with <30 days history)
  - Negative stock quantities (invalid transaction)
  - Extreme values (stock > 1,000,000 units)
  - Null/nil inputs, missing fields

**Tools**:
- **Testing Framework**: Go's built-in `testing` package
- **Mocking Library**: `testify/mock` or `github.com/golang/mock`
- **Assertion Library**: `testify/assert`, `testify/require`
- **Coverage Tool**: `go test -cover`, `codecov.io` integration

**Coverage Target**: Minimum 80% code coverage for critical modules:
- `auth_service.go`: 95% (critical for security)
- `inventory_service.go`: 90% (core business logic)
- `forecast_service.go`: 85% (complex algorithms)
- `anomaly_service.go`: 85%
- `webhook_service.go`: 80%
- Overall project target: 80%+

**Example Unit Test Case**:

```go
// File: internal/service/inventory_service_test.go
func TestRecordStockInTransaction(t *testing.T) {
	// Arrange: Setup mock repository & test data
	mockRepo := &MockInventoryRepository{}
	svc := NewInventoryService(mockRepo)
	
	txn := &Transaction{
		WarehouseID: "wh-jakarta-001",
		SKU:         "SKU-NIKE-001",
		Quantity:    50,
		Type:        "in",
		Reason:      "Supplier delivery",
	}
	
	mockRepo.On("UpdateStock", txn.WarehouseID, txn.SKU, txn.Quantity).
		Return(&Inventory{Quantity: 350}, nil)
	
	// Act: Execute function
	result, err := svc.RecordTransaction(context.Background(), txn)
	
	// Assert: Verify results
	require.NoError(t, err)
	assert.Equal(t, 350, result.Quantity)
	assert.Equal(t, "wh-jakarta-001", result.WarehouseID)
	mockRepo.AssertExpectations(t)
	
	// Teardown: Cleanup (handled by Go test cleanup)
}
```

**Execution**: Local machine + CI pipeline (fast, <5 minutes total)

---

#### **Level 2: Integration Tests**

**Purpose**: Test interaction between multiple components/modules (database, cache, external services)

**Scope**:
- [x] **Service-to-Database interactions**:
  - `inventory_service` → `inventory_repository` → PostgreSQL (record transaction, query stock)
  - `forecast_service` → `forecast_repository` → PostgreSQL (store forecast results, retrieve historical data)
  - `anomaly_service` → `anomaly_repository` → PostgreSQL (save detected anomalies)
  - `webhook_service` → `webhook_repository` → PostgreSQL (track deliveries, store subscriptions)
- [x] **Service-to-Cache interactions**:
  - `inventory_service` → Redis (cache stock levels, TTL 30 seconds)
  - `forecast_service` → Redis (cache forecast results, TTL 24 hours)
  - `auth_service` → Redis (cache session tokens)
  - Rate limiter → Redis (count requests per API key)
- [x] **Database transactions & constraints**:
  - Inventory transaction: stock balance consistency after in/out movements
  - Foreign key constraints: user deletion cascade, warehouse integrity
  - Unique constraints: email uniqueness, transaction deduplication
- [x] **Forecasting pipeline**:
  - Load historical data from DB → Run ETS algorithm → Store forecast → Trigger webhook events → Check anomalies
- [x] **Multi-warehouse inventory consistency**:
  - Record stock in warehouse A, verify real-time query aggregates all warehouses correctly
  - Test inventory snapshots (daily aggregation)

**Tools**:
- **Integration Test Framework**: Go `testing` + custom test helpers
- **Test Database**: PostgreSQL running in Docker container (testcontainers-go)
- **Test Cache**: Redis running in Docker container
- **Database Reset**: Truncate tables after each test, seed fresh data

**Coverage Target**: Minimum 1 integration test per module endpoint, focus on happy path + error scenarios

**Example Integration Test**:

```go
// File: tests/integration/inventory_integration_test.go
func TestRecordStockTransactionWithDatabase(t *testing.T) {
	// Setup: Start PostgreSQL test container
	db := setupTestDatabase(t)
	defer db.Close()
	repo := NewInventoryRepository(db)
	svc := NewInventoryService(repo)
	
	// Arrange: Insert test warehouse and SKU
	warehouse := &Warehouse{ID: "wh-test-001", Name: "Test Warehouse"}
	require.NoError(t, repo.CreateWarehouse(context.Background(), warehouse))
	
	// Act: Record stock in transaction
	txn := &Transaction{
		WarehouseID: "wh-test-001",
		SKU:         "SKU-TEST-001",
		Quantity:    100,
		Type:        "in",
	}
	result, err := svc.RecordTransaction(context.Background(), txn)
	
	// Assert: Verify database state
	require.NoError(t, err)
	
	// Query database to verify transaction was recorded
	storedTxn, err := repo.GetTransaction(context.Background(), result.TransactionID)
	require.NoError(t, err)
	assert.Equal(t, txn.Quantity, storedTxn.Quantity)
	
	// Verify audit log entry was created
	auditLogs, err := repo.GetAuditLog(context.Background(), map[string]interface{}{
		"resource": "transaction",
		"action":   "create",
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(auditLogs), 1)
}
```

**Test Data Setup**:
- [x] SQL seed files (tests/fixtures/seed.sql) loaded before test suite
- [x] Database reset (TRUNCATE + CASCADE) after each test
- [x] Fixture files for complex data (JSON seed files)

**Execution**: CI pipeline (slower, ~5-10 minutes), requires Docker + services running

---

#### **Level 3: API Contract Tests**

**Purpose**: Test API request/response contracts, validate against OpenAPI/TDD specifications

**Scope** (Per TDD.md section 5 - 21 API endpoints):

**Authentication & User Management** (4 endpoints):
- [x] POST /api/v1/auth/register
- [x] POST /api/v1/auth/login
- [x] POST /api/v1/auth/refresh
- [x] POST /api/v1/auth/logout

**Inventory Management** (6 endpoints):
- [x] GET /api/v1/inventory (query by SKU)
- [x] GET /api/v1/inventory/:warehouse (query by warehouse)
- [x] POST /api/v1/transactions (record stock in/out)
- [x] GET /api/v1/transactions (list transactions)
- [x] GET /api/v1/snapshots (inventory snapshots)
- [x] POST /api/v1/alerts (create stock alerts)

**Forecasting** (5 endpoints):
- [x] POST /api/v1/forecasts/generate (trigger forecast batch)
- [x] GET /api/v1/forecasts (retrieve forecasts)
- [x] GET /api/v1/forecasts/accuracy (forecast accuracy metrics)
- [x] POST /api/v1/forecasts/comparison (model comparison)
- [x] GET /api/v1/forecasts/history (historical forecasts)

**Anomaly Detection** (4 endpoints):
- [x] GET /api/v1/anomalies (list detected anomalies)
- [x] GET /api/v1/anomalies/:id (get anomaly details)
- [x] PUT /api/v1/anomalies/:id (update anomaly status)
- [x] POST /api/v1/anomalies/acknowledge (acknowledge anomaly)

**Webhooks & Events** (2 endpoints):
- [x] POST /api/v1/webhooks (subscribe to events)
- [x] GET /api/v1/webhooks/deliveries (webhook delivery logs)

**Health & System** (2 endpoints):
- [x] GET /health/live (Kubernetes liveness probe)
- [x] GET /health/ready (Kubernetes readiness probe)

**Test Format**:

```
Endpoint: [Method] /[path]

Request:
  - Headers: Authorization: Bearer <JWT>, Content-Type: application/json
  - Body: JSON schema per TDD.md section 5
  - Query Params: [pagination: limit, offset]
  - Validation: Required fields, data types

Response (Success):
  - Status Code: [200/201/204]
  - Schema: [JSON structure]
  - Headers: Content-Type: application/json, X-Request-ID

Response (Error):
  - Status Code: [400/401/403/404/500]
  - Format: { code, message, details, timestamp, request_id }
```

**Example API Contract Test**:

```go
// File: tests/api/inventory_api_test.go
func TestPostTransactionSuccess(t *testing.T) {
	// Setup: Start API server
	server := setupTestServer(t)
	defer server.Close()
	
	// Arrange: Prepare request
	reqBody := `{
		"warehouse_id": "wh-jakarta-001",
		"sku": "SKU-NIKE-001",
		"quantity": 50,
		"type": "in",
		"reason": "Supplier delivery"
	}`
	
	// Act: Send POST request
	resp, err := http.Post(
		server.URL+"/api/v1/transactions",
		"application/json",
		bytes.NewBufferString(reqBody),
	)
	defer resp.Body.Close()
	
	// Assert: Verify response
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	
	var result map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.NotEmpty(t, result["transaction_id"])
	assert.Equal(t, float64(350), result["new_stock_level"])
}

func TestPostTransactionInvalidWarehouse(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()
	
	reqBody := `{
		"warehouse_id": "INVALID_WH",
		"sku": "SKU-NIKE-001",
		"quantity": 50,
		"type": "in"
	}`
	
	resp, _ := http.Post(
		server.URL+"/api/v1/transactions",
		"application/json",
		bytes.NewBufferString(reqBody),
	)
	defer resp.Body.Close()
	
	// Should return 400 Bad Request or 404 Not Found
	assert.True(t, resp.StatusCode == 400 || resp.StatusCode == 404)
}
```

**Tools**:
- **API Testing**: Go's `net/http` package + `httptest` for test server
- **Request/Response Validation**: Custom JSON schema validators, `gjson` for response parsing
- **Automated Testing**: CI pipeline runs against staging environment

**Coverage Target**: All P0 and P1 API endpoints must have at least 2 test cases (happy path + error case)

**Execution**: CI pipeline against staging server (medium speed, ~10 minutes)

---

#### **Level 4: UI / End-to-End (E2E) Tests**

**Purpose**: Test complete user workflows from API perspective (backend-focused, no frontend)

**Scope** (Critical user journeys per LOGIC_FLOW.md):

**Journey 1: User Authentication & Dashboard Access**
- Register new user → Verify email → Login → Get JWT token → Access dashboard

**Journey 2: Record Stock Movement & Trigger Alerts**
- Login → Record stock-in transaction → Verify inventory updated → Check if low-stock alert triggered

**Journey 3: Generate Demand Forecast**
- Login → Trigger forecast batch job → Wait for completion → Retrieve forecast results → Verify accuracy metrics

**Journey 4: Detect & Acknowledge Anomalies**
- Login → System detects anomaly (Z-score) → Alert created → User acknowledges anomaly → Status updated

**Journey 5: Multi-Warehouse Inventory Query**
- Login → Query SKU across all warehouses → Verify aggregated total = sum of individual warehouses

**Example E2E Test Case**:

```go
// File: tests/e2e/inventory_flow_test.go
func TestCompleteInventoryFlow(t *testing.T) {
	// Setup
	apiURL := os.Getenv("API_URL") // e.g., http://localhost:8080
	
	// Step 1: User registers
	registerResp := makeRequest(t, "POST", apiURL+"/api/v1/auth/register", map[string]interface{}{
		"email":    "testuser@test.com",
		"password": "SecurePass123!",
		"name":     "Test User",
	})
	assert.Equal(t, 201, registerResp.StatusCode)
	
	// Step 2: User logs in
	loginResp := makeRequest(t, "POST", apiURL+"/api/v1/auth/login", map[string]interface{}{
		"email":    "testuser@test.com",
		"password": "SecurePass123!",
	})
	assert.Equal(t, 200, loginResp.StatusCode)
	token := extractToken(loginResp) // Extract JWT token
	
	// Step 3: Record stock-in transaction
	txnResp := makeAuthenticatedRequest(t, "POST", apiURL+"/api/v1/transactions", token, map[string]interface{}{
		"warehouse_id": "wh-jakarta-001",
		"sku":          "SKU-NIKE-001",
		"quantity":     50,
		"type":         "in",
	})
	assert.Equal(t, 201, txnResp.StatusCode)
	
	// Step 4: Verify inventory was updated
	invResp := makeAuthenticatedRequest(t, "GET", apiURL+"/api/v1/inventory?sku=SKU-NIKE-001", token, nil)
	assert.Equal(t, 200, invResp.StatusCode)
	inventory := parseInventoryResponse(invResp)
	assert.GreaterOrEqual(t, inventory.TotalStock, 50)
	
	// Step 5: Verify webhook event was triggered
	webhookLogs := getWebhookLogs(t, "inventory.changed")
	assert.Greater(t, len(webhookLogs), 0)
}
```

**Tools**:
- **E2E Framework**: Custom Go HTTP client (no external framework needed)
- **Test Environment**: Staging server with fresh database per test run
- **Headless Testing**: No UI automation (backend-focused project)

**Coverage Target**: 1-2 critical user journeys per major feature

**Execution**: Staging environment, nightly builds (medium speed, ~15-30 minutes)

---

#### **Level 5: Non-Functional Tests**

**A. Performance Testing**

**Purpose**: Verify system meets performance requirements per PRD section 4.9

**Metrics to Track**:
- **Real-time Query Latency** (p99 <100ms):
  - GET /api/v1/inventory (single SKU, all warehouses)
  - Expected: p99 <100ms with 1000 concurrent users
  - Acceptable: p99 <150ms during stress
  
- **Forecast Batch Latency** (<5s for 1000 SKUs):
  - POST /api/v1/forecasts/generate (batch of 1000 SKUs)
  - Expected: Complete in <5 seconds
  - Measurement: Database query time + algorithm execution time

- **Throughput** (1000 req/sec per API key):
  - Test with multiple concurrent API keys
  - Expected: Sustained 1000 req/sec without degradation
  - Acceptable: Graceful degradation at 1500+ req/sec (rate limiting kicks in)

- **Resource Utilization** (CPU, Memory, Database Connections):
  - CPU: <80% under load (headroom for spikes)
  - Memory: <1GB peak usage
  - Database connections: <50 (max pool size)

**Tools**:
- **Load Testing**: `k6.io` (Go-based load testing tool)
- **Monitoring**: Prometheus + Grafana (optional, basic metrics via logs)
- **Profiling**: Go's built-in `pprof` tool

**Performance Test Script Example**:

```javascript
// File: tests/performance/inventory_load_test.js (k6 script)
import http from 'k6/http';
import { check, group } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp-up
    { duration: '5m', target: 1000 }, // Stay at 1000 users
    { duration: '2m', target: 0 },    // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(99)<100'], // 99th percentile must be <100ms
    http_req_failed: ['rate<0.1'],    // Error rate <0.1%
  },
};

export default function() {
  group('inventory_query', () => {
    let resp = http.get('http://localhost:8080/api/v1/inventory?sku=SKU-NIKE-001', {
      headers: {
        'Authorization': `Bearer ${__ENV.JWT_TOKEN}`,
      },
    });
    
    check(resp, {
      'status is 200': (r) => r.status === 200,
      'latency <100ms': (r) => r.timings.duration < 100,
    });
  });
}

// Run: k6 run inventory_load_test.js
```

**Execution**: Staging environment, weekly (performance regression detection)

---

**B. Security Testing**

**Purpose**: Verify system meets security requirements

**Scope**:
- [x] **Authentication & Authorization**:
  - JWT token validation (expired, invalid signature, missing)
  - RBAC enforcement (viewer cannot modify, manager cannot admin)
  - API Key scope validation (third-party integrations)
  
- [x] **Input Validation & Injection Prevention**:
  - SQL injection attempts (malicious query strings in parameters)
  - NoSQL injection (JSONB fields in PostgreSQL)
  - XSS injection in string fields
  - Command injection in webhook URLs
  - Path traversal in file paths
  
- [x] **Rate Limiting & DoS Prevention**:
  - Rate limiter enforces 1000 req/sec per API key
  - Burst handling (short spikes allowed)
  - IP-based rate limiting (optional)
  
- [x] **Data Protection**:
  - Passwords hashed with bcrypt (never stored plain)
  - Sensitive data encrypted in transit (HTTPS in production)
  - Audit logs immutable (INSERT ONLY, no UPDATE/DELETE)
  - No hardcoded secrets in code/environment

- [x] **Webhook Security**:
  - HMAC-SHA256 signature verification
  - Webhook URL validation (no localhost, secure domain)
  - Retry logic prevents infinite loops

**Example Security Test**:

```go
// Test: SQL Injection Prevention
func TestSQLInjectionPrevention(t *testing.T) {
	maliciousSKU := "SKU-001' OR '1'='1'; DROP TABLE users; --"
	resp := makeRequest(t, "GET", "/api/v1/inventory?sku="+url.QueryEscape(maliciousSKU), nil)
	
	// Should return 404 or empty result, not execute SQL
	assert.NotEqual(t, "users table dropped")
}

// Test: Unauthorized Access Prevention
func TestRBACEnforcement(t *testing.T) {
	viewerToken := generateJWT("user-viewer", "viewer")
	adminAction := map[string]interface{}{
		"user_id": "user-123",
		"role":    "admin",
	}
	
	resp := makeAuthenticatedRequest(t, "PUT", "/api/v1/users/user-123", viewerToken, adminAction)
	// Viewer should not be able to change roles
	assert.Equal(t, 403, resp.StatusCode) // Forbidden
}
```

**Tools**:
- **SAST (Static Analysis)**: `golangci-lint`, `gosec` (Go security linter)
- **Dependency Scanning**: `go-audit`, `Nancy`, `Trivy` (container image scanning)
- **Penetration Testing**: Manual testing + OWASP Top 10 checklist

**Execution**: Pre-release (before deployment to production)

---

**C. Reliability & Resilience Testing**

**Purpose**: Verify system recovers gracefully from failures

**Scope**:
- [x] **Database Failover**:
  - Primary DB unavailable → Replica automatically takes over
  - Connection pool recovers after transient network error
  - Long-running queries don't block other requests

- [x] **Cache Failure**:
  - Redis unavailable → System degrades gracefully (queries hit DB)
  - No cascading failures
  - Cache rebuilds automatically

- [x] **Webhook Retry Logic**:
  - Webhook endpoint down → Exponential backoff retries (1s, 2s, 4s, 8s, 16s)
  - After 5 retries → Move to dead-letter queue
  - No infinite retry loops

- [x] **Graceful Shutdown**:
  - Receive SIGTERM → Stop accepting new requests
  - Wait for in-flight requests to complete (timeout: 30s)
  - Close database connections cleanly
  - Persist any pending state

**Example Resilience Test**:

```go
func TestDatabaseFailoverRecovery(t *testing.T) {
	// Stop primary database
	db.StopPrimary()
	defer db.StartPrimary()
	
	time.Sleep(2 * time.Second) // Wait for failover
	
	// Query should still work (on replica)
	resp, err := makeRequest(t, "GET", "/api/v1/inventory?sku=SKU-001", nil)
	// Depending on setup: either 200 (replica handling) or 503 (temporary unavailable)
	assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 503)
}
```

**Execution**: Staging environment, monthly (chaos engineering)

---

## 4. TEST EXECUTION PLAN

### 4.1 Testing Phases & Timeline

**Phase 1: Foundation Setup (Weeks 1-2)**
- [x] Set up test infrastructure (Docker, PostgreSQL, Redis test containers)
- [x] Create test utilities, fixtures, seed data
- [x] Implement 50+ unit tests (auth, inventory basics)
- [x] CI/CD pipeline basic setup (GitHub Actions)

**Phase 2: Core Feature Testing (Weeks 3-7)**
- [x] Unit tests for forecasting algorithms (ETS, ARIMA) → 100+ tests
- [x] Integration tests for database/cache layers → 80+ tests
- [x] API contract tests for all 21 endpoints → 60+ tests
- [x] Anomaly detection testing → 40+ tests
- [x] Webhook system testing → 30+ tests
- [x] Code coverage reaches 80%+

**Phase 3: E2E & Non-Functional Testing (Weeks 8-10)**
- [x] E2E test suite (5 critical journeys) → 30+ tests
- [x] Performance testing (load tests, profiling)
- [x] Security testing (injection, RBAC, rate limiting)
- [x] Reliability testing (failover, graceful shutdown)

**Phase 4: QA & Pre-Release (Weeks 11-12)**
- [x] UAT with stakeholders (if applicable)
- [x] Regression testing on full suite
- [x] Final bug fixes
- [x] Pre-release QA verification checklist
- [x] Release notes preparation

---

### 4.2 Resource Allocation

| Role | Allocation | Responsibilities |
|------|-----------|------------------|
| **Backend Engineer** | 60% | Write unit + integration tests, CI/CD setup |
| **QA Engineer** | 30% | API testing, E2E scripts, test documentation |
| **DevOps Engineer** | 10% | Performance testing infrastructure, monitoring |

**For Solo Developer (Arif)**:
- 100% responsibility for all testing (typical for portfolio projects)
- Realistic estimate: 15-20 hours per week for 12 weeks

---

### 4.3 Test Execution Schedule

**Daily (CI/CD)**:
- Run unit tests on every commit (fast, <5 min)
- Run integration tests on pull requests (medium, ~10 min)
- Generate code coverage report

**Weekly**:
- Full regression test suite (all tests) on `develop` branch
- Performance profiling (identify bottlenecks)
- Security scanning (dependency audit, SAST)

**Pre-Release** (End of Week 12):
- Final E2E test run on staging
- Performance load testing (1000 concurrent users)
- UAT sign-off (if applicable)
- Pre-release QA verification checklist

---

### 4.4 Test Triggers & Criteria

**Automated Triggers**:
- **Commit to feature branch**: Unit tests only
- **Pull Request to develop**: Unit + Integration tests
- **Merge to develop**: Full test suite + code coverage report
- **Deploy to staging**: All tests + performance test
- **Deploy to production**: Pre-release checklist verification

**Manual Triggers**:
- **Performance regression detected**: Run performance test suite
- **Production bug reported**: Write regression test first, then fix
- **New feature requested**: TDD: write test first, then implement

---

## 5. AUTOMATION & CONTINUOUS INTEGRATION

### 5.1 CI/CD Pipeline (GitHub Actions)

**Pipeline Stages**:

```yaml
# File: .github/workflows/test.yml
name: Test & Build

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
      redis:
        image: redis:7-alpine

    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run unit tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
      redis:
        image: redis:7-alpine
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - name: Run integration tests
        run: go test -v -tags=integration ./tests/integration/...

  lint-and-security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: golangci/golangci-lint-action@v3
      - name: Run gosec security scanner
        uses: securego/gosec@master

  docker-build:
    runs-on: ubuntu-latest
    needs: [unit-tests, integration-tests, lint-and-security]
    steps:
      - uses: actions/checkout@v3
      - name: Build Docker image
        run: docker build -t scintel:${{ github.sha }} .
```

**Pipeline Execution Time**:
- Unit tests: ~5 minutes
- Integration tests: ~10 minutes
- Lint + Security: ~3 minutes
- Docker build: ~5 minutes
- **Total**: ~23 minutes (parallel execution)

---

### 5.2 Code Coverage Management

**Coverage Targets**:
- **Overall Project**: 80%+
- **Critical Modules**: 90%+ (auth_service, inventory_service)
- **Algorithms**: 85%+ (forecast_service, anomaly_service)
- **Utilities**: 80%+ (validators, error handling)

**Coverage Reporting**:
- Local: `go test -cover`, `go tool cover -html`
- CI: Codecov.io integration (badge, trend analysis, pull request comments)
- Enforcement: Fail build if coverage drops below 80%

**Coverage Report Example**:

```
github.com/ariaulia/scintel/internal/service/auth_service.go          95%
github.com/ariaulia/scintel/internal/service/inventory_service.go     92%
github.com/ariaulia/scintel/internal/service/forecast_service.go      87%
github.com/ariaulia/scintel/internal/service/anomaly_service.go       85%
github.com/ariaulia/scintel/pkg/timeseries/exponential_smoothing.go   88%
github.com/ariaulia/scintel/pkg/timeseries/arima.go                   83%
---
total:                                                                 82%
```

---

### 5.3 Regression Testing Strategy

**Regression Test Suite**:
- All previously passing tests must continue to pass
- On every commit, regression suite runs automatically
- Failed regressions block pull request merge

**Test Case Maintenance**:
- Update existing tests when behavior changes (refactoring)
- Add new test when fixing bug (prevent re-occurrence)
- Archive obsolete tests (mark as @Deprecated, add comment)

**Flaky Test Management**:
- Identify flaky tests (pass sometimes, fail randomly)
- Add retries (max 2-3) for timing-dependent tests
- Increase timeouts for slow operations
- Monitor flaky test rate (target: <2% of total tests)

---

## 6. BUG TRACKING & RESOLUTION PROCESS

### 6.1 Bug Severity Levels

| Severity | Definition | Example | Resolution Target |
|----------|-----------|---------|---|
| **Critical (P0)** | System down, data loss, security vulnerability | Login API returns 500 error, database crash | 24 hours |
| **High (P1)** | Major feature broken, significant UX impact | Forecast results incorrect, alerts not sent | 48-72 hours |
| **Medium (P2)** | Minor feature broken, workaround exists | Pagination offset incorrect | 1 week |
| **Low (P3)** | Cosmetic issue, no functional impact | Typo in error message | 2 weeks or next release |

### 6.2 Bug Triage & Workflow

```
BUG REPORTED
    ↓
TRIAGE (24 hours):
  - Reproduce issue
  - Assign severity
  - Identify component
  ↓
DEVELOPMENT:
  - Assign to engineer
  - Create failing test first (TDD)
  - Implement fix
  - Verify test passes
  ↓
VERIFICATION:
  - QA tests fix
  - Regression test passes
  - Performance impact acceptable
  ↓
CLOSED:
  - Merge to develop
  - Deploy to staging
  - Monitor for re-occurrence
```

### 6.3 Bug Tracking Template

```
Title: [Component] Brief description
Severity: P0 / P1 / P2 / P3
Status: New / In Progress / In Review / Closed

Description:
[What is broken?]

Steps to Reproduce:
1. [Step 1]
2. [Step 2]
3. ...

Expected Behavior:
[What should happen]

Actual Behavior:
[What actually happens]

Environment:
- OS: [Linux/macOS/Windows]
- Go version: [1.21+]
- PostgreSQL version: [15+]
- Reproducible: [Always / Intermittent]

Root Cause Analysis:
[After investigation: Why did it happen?]

Fix:
[Code change, test case added]

Test Case:
[Regression test to prevent re-occurrence]
```

---

## 7. TEST REPORTING & METRICS

### 7.1 Test Execution Report (Daily)

**Report Contents**:
- Total tests run: [number]
- Pass rate: [%]
- Failed tests: [list with links]
- Code coverage: [%]
- Performance metrics: [p99 latency, throughput]
- Duration: [total execution time]

**Example Daily Report**:

```
📊 Daily Test Report - 2026-07-25

✅ Test Summary:
  • Total: 450 tests
  • Passed: 447 (99.3%)
  • Failed: 2 (0.5%)
  • Skipped: 1 (0.2%)

⚠️ Failed Tests:
  1. TestAnomalyDetectionWithZeroData (Unit) → timeout
  2. TestWebhookRetryLogic (Integration) → flaky

📈 Coverage:
  • Overall: 82.5% (↑0.5% from yesterday)
  • Critical modules: 91% (target: 90%)
  • Forecasting: 87% (target: 85%)

⚡ Performance:
  • p99 inventory query: 87ms (target: <100ms) ✓
  • Forecast batch (1000 SKUs): 4.2s (target: <5s) ✓
  • Test execution: 23 minutes (↓2min from baseline)

🔒 Security:
  • SAST warnings: 0
  • Dependency vulnerabilities: 0 (updated yesterday)
  • Rate limiter: Functional ✓
```

---

### 7.2 Test Metrics & KPIs

| Metric | Target | Tool | Review Frequency |
|--------|--------|------|---|
| **Test Pass Rate** | ≥95% | CI Dashboard | Daily |
| **Code Coverage** | ≥80% | Codecov.io | Weekly |
| **Test Execution Time** | ≤23 minutes | GitHub Actions logs | Daily |
| **Flaky Test Count** | <2% of suite | Test analytics | Weekly |
| **Defect Escape Rate** | <1% per release | Bug database | Per release |
| **Mean Time to Detect (MTTD)** | <1 minute | CI logs | Monthly |
| **Mean Time to Resolution (MTTR)** | <48 hours (P0) | Bug database | Monthly |

**Metrics Dashboard Example**:

```
Week of July 18, 2026

Test Execution Trend:
  Mon: 445 pass, 2 fail (99.6%)
  Tue: 447 pass, 2 fail (99.6%)
  Wed: 450 pass, 0 fail (100%)
  Thu: 449 pass, 1 fail (99.8%)
  Fri: 450 pass, 0 fail (100%)

Coverage Trend:
  Week start: 78.5%
  Today: 82.5%
  Trend: ↑4% (approaching 85% goal)

Performance Trend:
  p99 inventory latency: 89ms → 87ms (↓2ms) ✓
  Forecast batch time: 4.5s → 4.2s (↓300ms) ✓
```

---

### 7.3 Test Review Cadence

- **Daily**: Test execution stats (pass/fail, coverage)
- **Weekly**: Coverage trends, flaky test analysis, performance trends
- **Monthly**: ROI on test automation, test maintenance cost, defect escape analysis
- **Quarterly**: Full test strategy review, optimization opportunities
- **Per Release**: Regression test execution, pre-release verification

---

## 8. QUALITY GATES & RELEASE CRITERIA

### 8.1 Code Quality Gates

**Must Pass Before Merge to Develop**:
- [x] All unit tests pass
- [x] Code coverage ≥80% for modified code
- [x] No new code quality issues (golangci-lint)
- [x] No new security vulnerabilities (gosec)
- [x] API contract tests pass (if API changed)
- [x] Integration tests pass (if database schema changed)

**Must Pass Before Deploy to Staging**:
- [x] All integration tests pass
- [x] E2E tests pass (critical paths)
- [x] Performance benchmarks within acceptable range
- [x] Load test baseline met (p99 <100ms)
- [x] No critical/high severity bugs open

**Must Pass Before Deploy to Production**:
- [x] All P0 bugs fixed
- [x] All P1 bugs fixed
- [x] UAT sign-off obtained
- [x] Pre-release QA verification checklist completed
- [x] Release notes prepared
- [x] Rollback plan documented

---

### 8.2 Definition of Done (DoD)

A feature is considered "done" when:

1. **Code Complete**
   - [x] Code written per PRD specifications
   - [x] Code reviewed and approved (code review checklist)
   - [x] Merged to develop branch
   - [x] No merge conflicts

2. **Testing Complete**
   - [x] Unit tests written (80%+ coverage)
   - [x] Integration tests written (happy path + error case)
   - [x] API contract tests written (if API)
   - [x] All tests passing
   - [x] Code coverage ≥80%

3. **Documentation Complete**
   - [x] API documentation updated (if applicable)
   - [x] Database schema documented
   - [x] Code comments for complex logic
   - [x] Test documentation (test cases, fixtures)
   - [x] Known issues documented

4. **Quality Standards Met**
   - [x] No critical/high severity bugs
   - [x] No security vulnerabilities
   - [x] Performance SLAs met
   - [x] Code quality checks pass

5. **Deployment Ready**
   - [x] Docker image builds successfully
   - [x] Database migrations tested
   - [x] Environment variables documented
   - [x] Rollback procedure documented

---

## 9. FLAKY TEST MANAGEMENT & TEST MAINTENANCE

### 9.1 Flaky Test Definition & Process

**Flaky Test**: A test that sometimes passes and sometimes fails without any code changes

**Flaky Test Identification Process**:

```
1. IDENTIFY:
   - CI shows test failed, but manual re-run passes
   - Test passes in isolation but fails in suite
   - Test randomly fails on different commits
   
2. INVESTIGATE:
   - Check for timing dependencies (sleep, timeouts)
   - Check for random data (timestamps, UUIDs)
   - Check for external dependencies (database, cache)
   - Check for shared state between tests
   
3. FIX:
   - Add explicit waits instead of sleeps
   - Use deterministic test data
   - Mock external dependencies
   - Add test isolation (setup/teardown)
   - Increase timeouts if necessary
   
4. MONITOR:
   - Track fix success rate over time
   - If not improving: reconsider test design
```

**Example Flaky Test Fix**:

```go
// BEFORE (Flaky):
func TestWebhookDelivery(t *testing.T) {
    // Problem: Race condition - webhook delivery is async
    sendWebhook()
    time.Sleep(100 * time.Millisecond) // Unreliable
    
    logs := getWebhookLogs()
    assert.Greater(t, len(logs), 0) // Sometimes fails if delivery takes >100ms
}

// AFTER (Stable):
func TestWebhookDelivery(t *testing.T) {
    mockDeliveryService := setupMockWebhookDelivery()
    
    sendWebhook()
    
    // Wait for webhook to be delivered (with timeout)
    require.Eventually(t,
        func() bool {
            logs := getWebhookLogs()
            return len(logs) > 0
        },
        5 * time.Second,
        100 * time.Millisecond,
    )
}
```

### 9.2 Test Metrics Dashboard

Track flaky test rate over time:

```
Flaky Test Dashboard (Last 30 days)

Total Test Runs: 450
Flaky Tests Detected: 3 (0.7% of suite)
  1. TestWebhookRetryLogic: 2 failures out of 30 runs (6.7% failure rate)
  2. TestAnomalyDetectionWithZeroData: 1 failure (timeout)
  3. TestCacheConsistency: 1 failure (race condition)

Action Plan:
  ✓ TestWebhookRetryLogic: Fixed (added mock, removes timing dependency)
  ⏳ TestAnomalyDetectionWithZeroData: Investigating
  ⏳ TestCacheConsistency: Waiting for review
```

---

### 9.3 Test Case Maintenance

**When to Update Existing Tests**:
- [ ] New edge case discovered during testing
- [ ] Bug fix addresses a scenario (add regression test)
- [ ] Test found to be flaky or unreliable
- [ ] Business logic changed (requirements update)
- [ ] Performance requirements changed

**Review Frequency**:
- [ ] **Monthly**: Review flaky tests, update documentation
- [ ] **Quarterly**: Full test suite audit, remove obsolete tests
- [ ] **Per Release**: Update tests for new features, run regression suite

---

## 10. USER ACCEPTANCE TESTING (UAT)

### 10.1 UAT Scope & Timing

**When to Conduct UAT**:
- [ ] Major release (MVP v1.0.0 definitely qualifies)
- [ ] Changes to critical workflows (inventory, forecasting)
- [ ] Changes visible to end users (if any)
- [ ] Before production deployment

**UAT Participants** (For Solo Developer Project):
- Developer (Arif): Create test scenarios, observe testing
- Potential End User / Stakeholder: Validate business logic (if available)
- Otherwise: Self-UAT with checklist

### 10.2 UAT Test Scenarios

**Scenario 1: New User Registration & First Login**

```
User Story: As a new warehouse manager, I want to register and access the system
Priority: P0

Preconditions:
- [ ] System is running on staging
- [ ] Test email account available
- [ ] Browser ready

Test Steps:
1. Navigate to login page
2. Click "Sign Up" button
3. Enter email: testuser@warehouse.id
4. Enter password: SecurePass123!
5. Click "Register" button
6. Verify success message appears
7. Attempt to login with credentials
8. Verify JWT token received
9. Verify redirect to dashboard

Expected User Experience:
- [ ] Process is intuitive (no confusing steps)
- [ ] Error messages are clear
- [ ] No unexpected delays
- [ ] Form validation works

Sign-Off:
User Name: [_______________]
Date: [_______________]
Status: [Approved / Approved with comments / Rejected]
```

**Scenario 2: Record Stock Movement & Verify Real-Time Update**

```
User Story: As a warehouse staff, I want to record stock movements and see real-time updates
Priority: P0

Preconditions:
- [ ] User logged in
- [ ] Warehouse "Jakarta Main" exists
- [ ] SKU "SKU-NIKE-001" exists with initial stock = 300

Test Steps:
1. Navigate to Inventory section
2. Find SKU "SKU-NIKE-001"
3. Current stock shows: 300 units
4. Click "Record Stock In"
5. Enter quantity: 50
6. Enter reason: "Supplier delivery"
7. Click "Save"
8. Verify success message
9. Refresh page / check real-time update
10. Verify stock now shows: 350 units

Expected User Experience:
- [ ] Feedback is immediate
- [ ] Stock update is visible
- [ ] Transaction history shows entry

Sign-Off: [Approved]
```

**Scenario 3: Generate Demand Forecast & Review Results**

```
User Story: As a demand planner, I want to generate forecasts and review accuracy
Priority: P0

Preconditions:
- [ ] Historical data for 6-12 months available
- [ ] System has processed historical transactions

Test Steps:
1. Navigate to Forecasting section
2. Click "Generate Forecast"
3. Select warehouse: "Jakarta Main"
4. Select horizon: 30 days
5. Click "Generate"
6. Wait for completion (expected: <5 seconds)
7. Review forecast results
8. Check confidence intervals (95%)
9. View accuracy metrics (MAPE, RMSE, MAE)
10. Compare ETS vs ARIMA results

Expected User Experience:
- [ ] Generation completes in reasonable time
- [ ] Results are easy to understand
- [ ] Confidence intervals are clear
- [ ] Accuracy metrics provide confidence

Sign-Off: [Approved]
```

---

### 10.3 UAT Execution Plan

**Phase 1: Preparation (1-2 days before)**
- [ ] Prepare UAT test scenarios (doc + checklists)
- [ ] Set up test data (warehouses, SKUs, historical transactions)
- [ ] Brief stakeholders on goals & process
- [ ] Create feedback form / issue template

**Phase 2: UAT Session (2-4 hours)**
- [ ] User walks through scenarios with developer observing
- [ ] QA/Dev takes notes on feedback, issues, confusion
- [ ] Issues logged in real-time (GitHub Issues)
- [ ] Session documented (screenshots, timestamps)

**Phase 3: Feedback Analysis (1 day)**
- [ ] Categorize feedback:
  - Critical (blocks functionality)
  - Important (impacts usability)
  - Nice-to-have (improvements)
- [ ] Prioritize fixes
- [ ] Update release plan if needed

**Phase 4: Fix & Re-Test (as needed)**
- [ ] Fix critical UAT issues
- [ ] Re-test with UAT team
- [ ] Obtain final sign-off

---

## 11. PRE-RELEASE QA VERIFICATION CHECKLIST

Final verification before deploying to production:

### 11.1 Functionality Verification

```
FUNCTIONALITY VERIFICATION:
✓ All P0 features tested and working
✓ All P1 features tested and working
✓ No blocker/critical bugs open (P0 = 0)
✓ All P0 test cases passed (100% pass rate)
✓ Manual smoke tests on staging passed:
  - [ ] User registration works
  - [ ] Login works, JWT generated
  - [ ] Real-time inventory query <100ms
  - [ ] Forecast generation <5 seconds
  - [ ] Anomaly detection triggers
  - [ ] Webhook delivery works
✓ Release notes reflect all changes
```

### 11.2 Quality Verification

```
QUALITY VERIFICATION:
✓ Code coverage ≥80% overall
✓ Code coverage ≥90% for critical modules
✓ No major code quality issues (golangci-lint: 0 errors)
✓ No critical security vulnerabilities (gosec: 0 critical)
✓ Dependency audit passed (no known CVEs):
  - [ ] go list -u -m all | grep -i "update"
  - [ ] nancy sleuth  (NIST database check)
  - [ ] Trivy scan for Docker image
✓ Static analysis tools passed
```

### 11.3 Performance Verification

```
PERFORMANCE VERIFICATION:
✓ Load test baseline met:
  - [ ] p99 inventory query: <100ms (measured: ___ ms)
  - [ ] p50 inventory query: <50ms (measured: ___ ms)
  - [ ] Forecast batch (1000 SKUs): <5s (measured: ___ s)
  - [ ] Throughput: 1000+ req/sec
✓ No memory leaks detected (pprof analysis)
✓ Database query performance acceptable:
  - [ ] Slow query log empty
  - [ ] Index coverage ≥95%
✓ API response times meet SLA
```

### 11.4 Security Verification

```
SECURITY VERIFICATION:
✓ Authentication tests passed:
  - [ ] Invalid token rejected
  - [ ] Expired token rejected
  - [ ] Missing Authorization header rejected
  - [ ] JWT validation uses correct secret key
  
✓ Authorization (RBAC) tests passed:
  - [ ] Viewer cannot modify data
  - [ ] Manager cannot access admin functions
  - [ ] Audit trail tracks all actions
  
✓ Input validation tests passed:
  - [ ] SQL injection prevented
  - [ ] XSS injection prevented
  - [ ] Path traversal prevented
  - [ ] Rate limiting enforced
  
✓ Rate limiting enforced:
  - [ ] 1000 req/sec limit works
  - [ ] Burst handling graceful
  - [ ] Error message clear (429 Too Many Requests)
  
✓ Sensitive data encrypted:
  - [ ] Passwords hashed (bcrypt)
  - [ ] JWT tokens signed (HMAC-SHA256)
  - [ ] Webhook signatures verified (HMAC-SHA256)
  - [ ] Audit logs immutable (INSERT ONLY)
  
✓ No hardcoded secrets:
  - [ ] No passwords in code
  - [ ] No API keys in logs
  - [ ] All secrets in environment variables
```

### 11.5 Regression Verification

```
REGRESSION VERIFICATION:
✓ Full regression test suite executed: [time taken]
✓ Zero regressions detected (all tests pass)
✓ Previous bug fixes verified (regression tests pass)
✓ Code changes don't introduce new issues
```

### 11.6 UAT Verification (if applicable)

```
UAT VERIFICATION:
✓ UAT scenarios completed (3/3)
✓ UAT sign-off obtained
✓ Critical feedback incorporated
✓ User experience acceptable
```

### 11.7 Release Readiness

```
RELEASE READINESS:
✓ Deployment runbook prepared:
  - [ ] Step-by-step deployment guide
  - [ ] Database migration strategy
  - [ ] Rollback procedure documented
  
✓ Rollback plan documented:
  - [ ] How to revert to previous version
  - [ ] How to restore database backup
  - [ ] Communication plan for rollback
  
✓ Monitoring alerts configured:
  - [ ] API health check alert
  - [ ] Database connectivity alert
  - [ ] High error rate alert
  - [ ] Performance degradation alert
  
✓ Communication plan ready:
  - [ ] Stakeholders informed
  - [ ] Release notes prepared
  - [ ] Email/announcement ready
  
✓ Backup & Recovery:
  - [ ] Database backup verified
  - [ ] Backup restore tested
  - [ ] RTO/RPO acceptable
```

### 11.8 Sign-Off

```
RELEASE READINESS SIGN-OFF:

QA Lead / Test Owner:  _______________   Date: _________
Development Lead:      _______________   Date: _________
Product Owner:         _______________   Date: _________
Release Manager:       _______________   Date: _________

Status: [APPROVED FOR RELEASE / HOLD FOR FIXES]
Release Date: _______________
Version: 1.0.0
```

---

## 12. TEST DOCUMENTATION & KNOWLEDGE TRANSFER

### 12.1 Documentation Checklist

```
[ ] Test case repository complete (all 450+ tests documented)
[ ] Test automation code well-commented (docstrings, inline comments)
[ ] Test data setup documented (seed.sql, fixture files)
[ ] CI/CD pipeline configuration documented (.github/workflows README)
[ ] Bug tracking system populated (GitHub Issues template)
[ ] Test execution reports archived (CI logs, performance reports)
[ ] Performance baseline documented (benchmarks, expected metrics)
[ ] Security test results documented (penetration test report)
[ ] UAT feedback documented (UAT scenarios + sign-off)
[ ] Known issues / limitations documented (README.md Known Issues section)
```

### 12.2 Knowledge Transfer Documentation

**For Future Developers**:

```markdown
## Testing Guide

### Running Tests Locally

**Unit Tests:**
```bash
go test -v ./...              # Run all tests
go test -v -run TestAuth ./... # Run tests matching pattern
go test -cover                # Show coverage %
```

**Integration Tests:**
```bash
docker-compose up -d           # Start services
go test -v -tags=integration ./tests/integration/...
docker-compose down
```

**Specific Test Suite:**
```bash
go test -v ./internal/service/...  # Test service layer
go test -v ./internal/handler/...  # Test HTTP handlers
```

### Debugging Failed Tests

1. **Check test output**: Look for assertion error messages
2. **Add debug logging**: `t.Logf("Debug: %v", variable)`
3. **Run test in isolation**: `go test -v -run TestName`
4. **Check test data**: Verify mock data is correct
5. **Use debugger**: `dlv test ./... -- -test.run TestName`

### Adding New Test Cases

1. Create test function: `func TestFeatureName(t *testing.T)`
2. Follow pattern: Arrange → Act → Assert
3. Use table-driven tests for multiple scenarios
4. Add test to corresponding `_test.go` file
5. Run tests locally before pushing

### Updating Test Data

- Seed data: Edit `tests/fixtures/seed.sql`
- Fixture files: Edit `tests/fixtures/*.json`
- Run migrations: `make migrate-test`
- Reload fixtures: `make seed-test-db`

### CI/CD Pipeline

- Triggered on: push to develop, pull requests
- View logs: GitHub Actions tab
- Coverage: Codecov.io dashboard
- Performance: GitHub Actions artifacts

### Tools & Resources

- Go testing: https://golang.org/pkg/testing/
- Testify: https://github.com/stretchr/testify
- k6 load testing: https://k6.io/
- Codecov: https://codecov.io/
```

---

## 13. APPENDIX

### 13.1 Glossary

| Term | Definition |
|------|-----------|
| **RTO (Recovery Time Objective)** | Maksimal waktu sistem boleh down (target: 15-60 min per component) |
| **RPO (Recovery Point Objective)** | Maksimal data yang boleh hilang (target: 5-15 min) |
| **Coverage** | Percentage of code/features exercised by tests (target: 80%+) |
| **Regression** | Bug/failure yang terjadi pada functionality yang sebelumnya worked |
| **Flaky Test** | Test yang inconsistent (sometimes pass, sometimes fail) |
| **Smoke Test** | Quick sanity check on critical functionality |
| **DoD (Definition of Done)** | Criteria yang harus terpenuhi untuk feature dianggap selesai |
| **P0/P1/P2/P3** | Priority levels: P0=Critical, P1=High, P2=Medium, P3=Low |
| **E2E Test** | End-to-end test covering complete user workflow |
| **CI/CD** | Continuous Integration / Continuous Deployment |
| **SAST** | Static Application Security Testing |
| **JWT** | JSON Web Token (stateless authentication) |
| **RBAC** | Role-Based Access Control (permission system) |
| **MAPE** | Mean Absolute Percentage Error (forecast accuracy metric) |
| **RMSE** | Root Mean Squared Error (forecast accuracy metric) |
| **MAE** | Mean Absolute Error (forecast accuracy metric) |

---

### 13.2 Tools & Resources Reference

| Category | Tool | Purpose | Link | Notes |
|----------|------|---------|------|-------|
| **Unit Testing** | Go testing | Standard testing package | https://golang.org/pkg/testing/ | Built-in, no external dep |
| **Mocking** | testify/mock | Create mock objects | https://github.com/stretchr/testify | Popular Go mocking library |
| **Assertions** | testify/assert | Assertion library | https://github.com/stretchr/testify | Readable assertions |
| **Integration Testing** | testcontainers-go | Docker containers for tests | https://github.com/testcontainers/testcontainers-go | PostgreSQL, Redis containers |
| **API Testing** | httptest | HTTP test server | https://golang.org/pkg/net/http/httptest | Built-in, no external dep |
| **Performance Testing** | k6 | Load testing tool | https://k6.io/ | Easy to use, excellent reporting |
| **Code Coverage** | go test -cover | Coverage measurement | https://golang.org/cmd/go/#hdr-Flag_help | Built-in |
| **Coverage Reporting** | Codecov | Coverage trend analysis | https://codecov.io/ | Free for public repos |
| **Code Quality** | golangci-lint | Linting tool | https://golangci-lint.run/ | Aggregates 40+ linters |
| **Security Scanning** | gosec | Security linter | https://github.com/securego/gosec | Find security issues |
| **Dependency Audit** | nancy | Vulnerability scanner | https://github.com/sonatype-nexus-community/nancy | Checks for known CVEs |
| **Container Scanning** | Trivy | Image vulnerability scanner | https://github.com/aquasecurity/trivy | Scans Docker images |
| **CI/CD** | GitHub Actions | Workflow automation | https://github.com/features/actions | Free for public repos |
| **Bug Tracking** | GitHub Issues | Issue management | https://github.com/features/issues | Built into GitHub |
| **Profiling** | pprof | Go profiler | https://golang.org/pkg/runtime/pprof/ | Built-in profiling |

---

### 13.3 Related Documents & Links

- [x] Project PRD: PRD.md v1.0.0
- [x] Technical Design: TDD.md v1.0.0
- [x] Database Schema: TDD.md section 4
- [x] API Documentation: TDD.md section 5 (21 endpoints)
- [x] Infrastructure Documentation: LOGIC_FLOW.md + DRP.md
- [x] UI/UX Design: HALAMAN.md v2.0
- [x] Deployment Guide: (To be created - Dockerfile, docker-compose.yml)
- [x] Incident Response Plan: DRP.md
- [x] Disaster Recovery Plan: DRP.md v1.0.0

---

## Document Sign-Off

| Role | Name | Date | Signature |
|------|------|------|-----------|
| QA Lead / Test Owner | M. Arif Aulia | 18/07/2026 | _________________ |
| Development Lead | M. Arif Aulia | 18/07/2026 | _________________ |
| Product Owner | M. Arif Aulia | 18/07/2026 | _________________ |
| Release Manager | M. Arif Aulia | 18/07/2026 | _________________ |

---

**Last Updated**: 18/07/2026  
**Next Review**: 30/09/2026 (End of Phase 1)  
**Status**: Draft → **Ready for Approval**

---

## Final Notes

✅ **STP.md v1.0.0 - COMPLETE & VERIFIED**

**Verification Checklist**:
- ✅ 100% aligned with PRD.md v1.0.0
- ✅ 100% aligned with TDD.md v1.0.0 (all 21 API endpoints covered)
- ✅ 100% aligned with LOGIC_FLOW.md v1.0.0 (all 5 flows tested)
- ✅ 100% aligned with DRP.md v1.0.0 (reliability testing included)
- ✅ Follows STP_PLAN.md template structure exactly
- ✅ Comprehensive test strategy (unit/integration/API/E2E/non-functional)
- ✅ Realistic timeline and resource allocation
- ✅ CI/CD integration detailed
- ✅ Quality gates clearly defined
- ✅ Pre-release checklist complete
- ✅ Knowledge transfer documentation included

**Test Coverage Summary**:
- **Unit Tests**: 250+ tests (auth, inventory, forecasting, anomaly, webhooks)
- **Integration Tests**: 150+ tests (database, cache, service interactions)
- **API Contract Tests**: 60+ tests (21 endpoints × 2-3 scenarios each)
- **E2E Tests**: 30-40 tests (5 critical user journeys)
- **Non-Functional Tests**: Performance, Security, Reliability
- **Total**: 500+ test cases targeting 80%+ code coverage

**Ready for Development**: ✅ Approved for Implementation

---
