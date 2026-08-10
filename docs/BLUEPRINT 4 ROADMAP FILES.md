# MASTER BLUEPRINT: 4 ROADMAP FILES FOR SCINTEL
**Single Source of Truth - LOCKED REFERENCE**

**Project:** Supply Chain Intelligence Platform (ScIntel)  
**Purpose:** Blueprint untuk 4 ROADMAP files  
**Created:** 2026-07-19  
**Status:** FINAL - Ready untuk ROADMAP implementation  
**Token-Safe:** Designed untuk avoid hallucination across chat sessions  

---

## 📋 TABLE OF CONTENTS

1. [Project Meta](#project-meta)
2. [4 Roadmap Files Overview](#4-roadmap-files-overview)
3. [ROADMAP_1: Foundation (Weeks 1-2)](#roadmap_1-foundation-weeks-1-2)
4. [ROADMAP_2: Backend Core (Weeks 3-4)](#roadmap_2-backend-core-weeks-3-4)
5. [ROADMAP_3: Frontend/Mobile (Weeks 5-8)](#roadmap_3-frontendmobile-weeks-5-8)
6. [ROADMAP_4: Mastery/Deployment (Week 9)](#roadmap_4-masterydeployment-week-9)
7. [Cross-Phase Dependencies](#cross-phase-dependencies)
8. [Verification Checklist](#verification-checklist)

---

## PROJECT META

### Core Information (from PRD.md)

**Project Name:** Supply Chain Intelligence Platform (ScIntel)  
**Timeline:** 9 weeks total (October 2026)  
**Backend:** Golang 1.21+  
**Database:** PostgreSQL 15+  
**Cache:** Redis 7+  
**Target:** Production-grade portfolio project (10/10 rating)

### Target Users (3 personas)

1. **Warehouse Managers** — Real-time inventory visibility
2. **Supply Chain Analysts** — Demand forecasting & anomaly alerts
3. **System Admins** — Audit logs, user management, webhooks

### Core Features (from PRD.md - IN SCOPE)

| Feature | Priority | Status |
|---------|----------|--------|
| User Authentication (JWT + RBAC) | P0 | Weeks 1-2 |
| Real-time Inventory Tracking | P0 | Weeks 2-4 |
| Demand Forecasting (ETS + ARIMA) | P0 | Weeks 3-5 |
| Anomaly Detection (Z-score) | P1 | Weeks 4-6 |
| Webhook System (events + retry) | P1 | Weeks 3-4 |
| Audit Trail & Compliance Logging | P2 | Weeks 7-8 |
| Rate Limiting & Security | P0 | Weeks 1-2 |
| API Documentation & Deployment | P2 | Week 9 |

### Tech Stack (from TDD.md)

**Backend:**
- Language: Go 1.21+
- Framework: Standard library + minimal packages
- HTTP Router: Custom mux

**Database:**
- Primary: PostgreSQL 15+
- Connection Pool: 10-50 connections
- Isolation Level: READ COMMITTED

**Cache:**
- Redis 7+ (inventory, forecasts, rate limiting)
- TTL: 30s for inventory, 24h for forecasts

**Testing:**
- Unit tests: >80% coverage
- Integration tests: core flows
- Framework: testing + table-driven tests

**DevOps:**
- Container: Docker (multi-stage, <150MB)
- Local: Docker Compose
- Staging/Prod: VPS + systemd or Render.com
- CI/CD: GitHub Actions

---

## 4 ROADMAP FILES OVERVIEW

### Summary Table

| File | Duration | Phase | Main Deliverables | Dependencies |
|------|----------|-------|-------------------|--------------|
| ROADMAP_1 | Weeks 1-2 | Foundation | DB setup, auth, config, project structure | None (start here) |
| ROADMAP_2 | Weeks 3-4 | Backend Core | API endpoints, business logic, webhooks | Depends on ROADMAP_1 |
| ROADMAP_3 | Weeks 5-8 | Frontend/Mobile | Dashboard, charts, real-time features | Depends on ROADMAP_1 + 2 |
| ROADMAP_4 | Week 9 | Mastery/Deployment | Testing, docs, deployment, optimization | Depends on all 1-3 |

### Total Effort Breakdown

```
ROADMAP_1 (Weeks 1-2): 40 hours
  - DB design & migrations: 15h
  - Authentication system: 15h
  - Project setup & testing: 10h

ROADMAP_2 (Weeks 3-4): 40 hours
  - Inventory API + logic: 18h
  - Forecasting API + algorithms: 15h
  - Webhook system: 7h

ROADMAP_3 (Weeks 5-8): 60 hours
  - Dashboard & UI: 25h
  - Charts & visualizations: 15h
  - Real-time features: 12h
  - Mobile responsive: 8h

ROADMAP_4 (Week 9): 20 hours
  - Unit tests: 8h
  - Integration tests: 6h
  - Documentation: 4h
  - Deployment: 2h

TOTAL: ~160 hours = 9 weeks @ 18h/week
```

---

## ROADMAP_1: FOUNDATION (WEEKS 1-2)

### Phase Goals

Build solid foundation with authentication, database, and project structure. **No feature delivery yet** - just infrastructure.

### Week 1: Project Setup & Database

#### Week 1.1: Environment & Project Structure (3 hours)

**Tasks:**
1. Create GitHub repo (public + portfolio-ready README)
2. Setup Go project structure:
   - cmd/api/main.go
   - internal/{handler, service, repository, model, middleware, util, database, config}
   - pkg/{logger, timeseries}
   - migrations/, tests/, docs/
3. Create docker-compose.yml (PostgreSQL, Redis, local dev)
4. Create .env.example with all config variables
5. Setup GitHub Actions CI/CD pipeline (lint, test, build)
6. Create Makefile with common commands (make dev, make test, make build, etc.)

**Deliverables:**
- ✅ Git repo initialized & first commit
- ✅ Folder structure complete
- ✅ docker-compose.yml ready for `docker-compose up -d`
- ✅ GitHub Actions workflow file (.github/workflows/ci.yml)

**Reference:**
- TDD Section 3: Architecture & Project Structure
- PRD Section 6: Technology Stack

---

#### Week 1.2: Database Design & Migrations (6 hours)

**Tasks:**
1. Install & configure PostgreSQL 15+ locally (via docker-compose)
2. Install golang-migrate CLI
3. Create migration files (golang-migrate format):
   - `000001_create_users_table.up.sql` - users table (Section 4.1)
   - `000002_create_warehouses_table.up.sql` - warehouses (Section 4.2)
   - `000003_create_skus_table.up.sql` - SKUs (Section 4.3)
   - `000004_create_inventory_table.up.sql` - inventory (Section 4.4)
   - `000005_create_transactions_table.up.sql` - transactions (Section 4.5)
   - `000006_create_forecast_tables.up.sql` - forecasts + forecast_accuracy (Sections 4.6-4.7)
   - `000007_create_anomalies_table.up.sql` - anomalies (Section 4.8)
   - `000008_create_webhook_tables.up.sql` - webhook_subscriptions + webhook_deliveries (Sections 4.9-4.10)
   - `000009_create_audit_log_table.up.sql` - audit_log (Section 4.11)
   - `000010_create_password_reset_tokens.up.sql` - password_reset_tokens (Section 4.12)
4. Run migrations (`make db-migrate`)
5. Seed test data (warehouses, users, SKUs)
6. Test all table schemas match TDD Section 4

**Deliverables:**
- ✅ All 10 migration files created & tested
- ✅ PostgreSQL running locally with all tables
- ✅ Test data seeded (3 warehouses, 5 users, 100 SKUs, 500 inventory records)
- ✅ No migration errors or warnings

**Reference:**
- TDD Section 4.0-4.12: All database tables
- TDD Section 4.0: Connection pool config (min 10, max 50)
- TDD Section 4.0: Transaction isolation READ COMMITTED

---

#### Week 1.3: Database Connection & Repository Layer (4 hours)

**Tasks:**
1. Create `internal/database/connection.go`:
   - Initialize PostgreSQL connection pool (min 10, max 50)
   - Connection timeout: 5 seconds
   - Idle timeout: 30 seconds
   - Error handling for connection failures
2. Create `internal/database/transaction.go`:
   - Transaction helpers (Begin, Commit, Rollback)
   - Deadlock retry logic (retry 3x, backoff 100/200/400ms)
   - SELECT ... FOR UPDATE for pessimistic locking
3. Create base repository:
   - `internal/repository/base.go` - database interface
   - Dependency injection for database connection
4. Create `internal/util/response.go`:
   - Standard HTTP response wrapper
   - Error response formatter (code, message, details, timestamp, request_id)

**Deliverables:**
- ✅ DB connection pool working
- ✅ Deadlock retry tested
- ✅ Response wrapper tested with error codes (TDD Section 5.1)
- ✅ All repository interfaces defined

**Reference:**
- TDD Section 4.0: Database Configuration
- TDD Section 5.1.1-5.1.2: Error Response Standards

---

### Week 2: Authentication & Authorization

#### Week 2.1: User Model & Repository (3 hours)

**Tasks:**
1. Create `internal/model/user.go`:
   - User struct matching TDD Section 4.1
   - Role enum (admin, manager, analyst, viewer)
   - Status enum (active, suspended, pending_email_verification)
2. Create `internal/repository/user_repository.go`:
   - CreateUser(email, password_hash, full_name, role)
   - GetUserByEmail(email)
   - GetUserByID(id)
   - UpdateUserRole(id, role)
   - UpdateUserStatus(id, status)
   - ListUsers(limit, offset)
3. Create `internal/util/hash.go`:
   - Bcrypt hashing (rounds: 12)
   - VerifyPassword function
4. Test all repository methods with database

**Deliverables:**
- ✅ User repository 100% tested
- ✅ Bcrypt integration working
- ✅ All CRUD operations verified

**Reference:**
- TDD Section 4.1: USERS TABLE
- TDD Section 8: Security (password hashing)

---

#### Week 2.2: JWT & Token Management (4 hours)

**Tasks:**
1. Create `internal/util/jwt.go`:
   - JWT creation (HS256, HMAC-SHA256)
   - Access token payload (user_id, email, role, exp)
   - Refresh token payload (user_id, jti, exp)
   - Token verification with error handling
   - Refresh token rotation
2. Create `internal/config/config.go`:
   - Load environment variables (.env)
   - JWT_SECRET_KEY, JWT_ACCESS_EXPIRY (1 hour), JWT_REFRESH_EXPIRY (7 days)
   - Database URL, Redis URL, etc.
3. Create token expiry handling (refresh before expiry)
4. Test JWT generation & verification

**Deliverables:**
- ✅ JWT tokens generated & validated
- ✅ Access token TTL: 1 hour, Refresh: 7 days
- ✅ Environment variables loaded from .env
- ✅ Token verification tested with expired/invalid tokens

**Reference:**
- TDD Section 8: Security & Authentication
- LOGIC_FLOW Section 2: Authentication Flow

---

#### Week 2.3: Authentication Endpoints (4 hours)

**Tasks:**
1. Create `internal/handler/auth_handler.go`:
   - POST /api/v1/auth/register
     - Request: {email, password, full_name}
     - Response: {user_id, email, role, message}
     - Validate: email format, password strength, unique email
     - Default role: viewer
   - POST /api/v1/auth/login
     - Request: {email, password}
     - Response: {access_token, refresh_token, expires_in, user}
     - Rate limiting: 5 failed attempts in 5 min → lock 5 min
   - POST /api/v1/auth/refresh
     - Request: {refresh_token}
     - Response: {access_token, expires_in, token_type}
   - POST /api/v1/auth/logout
     - Request: {refresh_token}
     - Response: {message: "Successfully logged out"}
2. Create `internal/middleware/auth_middleware.go`:
   - JWT verification
   - Extract user_id, role from token
   - Add to request context
3. Test all endpoints with valid/invalid inputs

**Deliverables:**
- ✅ All 4 auth endpoints working
- ✅ Error responses match TDD Section 5.1 format
- ✅ Auth middleware protecting routes
- ✅ Rate limiting on login (5 attempts = lock)

**Reference:**
- TDD Section 5.2: Authentication Endpoints
- LOGIC_FLOW Section 2: Authentication Flow

---

#### Week 2.4: RBAC & Logging (3 hours)

**Tasks:**
1. Create `internal/middleware/rbac_middleware.go`:
   - Role-based access control
   - Check user role against endpoint requirements
   - Return 403 Forbidden if insufficient role
   - Define role hierarchy: admin > manager > analyst > viewer
2. Create `internal/middleware/logging_middleware.go`:
   - Log all HTTP requests (method, path, status, duration)
   - Structured JSON logging (using pkg/logger)
   - Include request_id for tracing
3. Create `pkg/logger/logger.go`:
   - Structured logging with zap or similar
   - JSON output format
   - Log levels: debug, info, warn, error
4. Create `internal/repository/audit_repository.go`:
   - LogAction(user_id, action, resource, details)
   - For audit trail (Section 4.11)

**Deliverables:**
- ✅ RBAC middleware protecting sensitive endpoints
- ✅ Structured logging in JSON format
- ✅ Request tracing with request_id
- ✅ Audit log recording user actions

**Reference:**
- PRD Section 7: RBAC (4 roles)
- TDD Section 8: Authentication & RBAC

---

### Week 1-2 Deliverables Summary

**Code Commits:**
- [ ] Commit 1: Initial project structure + docker-compose
- [ ] Commit 2: All database migrations + schema
- [ ] Commit 3: Database connection + repository base
- [ ] Commit 4: User model + repository + bcrypt
- [ ] Commit 5: JWT + token management
- [ ] Commit 6: Authentication endpoints
- [ ] Commit 7: RBAC + logging + audit trail

**Testing:**
- [ ] >80% unit test coverage for auth & database layers
- [ ] Docker containers all healthy (`docker ps`)
- [ ] Migrations run clean with no errors
- [ ] All 4 auth endpoints tested (manual Postman or script)

**Documentation:**
- [ ] README.md with setup instructions
- [ ] .env.example with all variables
- [ ] docs/DATABASE.md with schema overview
- [ ] docs/AUTH.md with JWT flow diagrams

**Ready for ROADMAP_2?** YES - All foundation complete, ready for API implementation

---

## ROADMAP_2: BACKEND CORE (WEEKS 3-4)

### Phase Goals

Implement core business logic APIs: inventory management, forecasting, anomalies, webhooks. All require database & auth from ROADMAP_1.

### Week 3: Inventory & Anomaly Core

#### Week 3.1: Inventory Service & Repository (6 hours)

**Tasks:**
1. Create `internal/model/inventory.go`:
   - Inventory struct (warehouse_id, sku_id, quantity, last_updated)
   - InventorySnapshot (historical snapshots)
   - InventoryTransaction (all in/out movements)
2. Create `internal/repository/inventory_repository.go`:
   - GetInventoryByWarehouse(warehouse_id) → all SKUs in warehouse
   - GetInventoryBySKU(sku_id) → all warehouses with this SKU
   - GetTotalStock(sku_id) → sum across all warehouses
   - UpdateInventory(warehouse_id, sku_id, quantity_delta) → with FOR UPDATE lock
   - CreateSnapshot(warehouse_id, timestamp) → daily snapshots
3. Create `internal/service/inventory_service.go`:
   - Business logic layer
   - RecordTransaction(warehouse_id, sku_id, quantity, type, reason)
     - Type: "in" or "out"
     - Validate sufficient stock for "out"
     - Lock inventory row (SELECT ... FOR UPDATE)
     - Update quantity
     - Insert transaction log
     - Check anomalies (async)
     - Invalidate Redis cache
     - Publish webhook "inventory.changed"
   - GetInventoryWithCache(warehouse_id, sku_id)
     - Check Redis (TTL 30s)
     - If miss, query DB & update cache

**Deliverables:**
- ✅ Inventory service fully tested
- ✅ Pessimistic locking working (SELECT ... FOR UPDATE)
- ✅ Cache invalidation tested
- ✅ Transaction log recording all movements

**Reference:**
- TDD Section 4.4-4.5: Inventory & Transactions tables
- LOGIC_FLOW Section 4: Inventory Flow (steps 1-7)

---

#### Week 3.2: Inventory Endpoints (5 hours)

**Tasks:**
1. Create `internal/handler/inventory_handler.go`:
   - GET /api/v1/inventory/{warehouse_id}
     - Get all inventory in warehouse (with SKU details)
   - GET /api/v1/inventory
     - Get all warehouses' inventory (paginated)
   - GET /api/v1/inventory/{warehouse_id}/{sku_id}
     - Get specific SKU in warehouse
   - POST /api/v1/transactions
     - Record inventory movement
     - Request: {warehouse_id, sku_id, quantity, type, reason}
     - Validate role: admin, manager only
     - Manager can only update own warehouse
   - GET /api/v1/transactions
     - Query transaction history (filtered, paginated)
     - Filters: warehouse_id, sku_id, type, date range
2. Create Redis integration:
   - Setup `github.com/redis/go-redis` client
   - Cache key: `inventory:{warehouse_id}:{sku_id}`
   - TTL: 30 seconds
3. Test all endpoints with valid/invalid data

**Deliverables:**
- ✅ 5 inventory endpoints working
- ✅ Pagination working (limit, offset)
- ✅ Redis caching integrated
- ✅ Role-based access control working
- ✅ Error responses match TDD Section 5.1

**Reference:**
- TDD Section 5.2: Inventory Endpoints
- TDD Section 7.1: Rate Limiting (add middleware)

---

#### Week 3.3: Anomaly Detection Setup (4 hours)

**Tasks:**
1. Create `internal/model/anomaly.go`:
   - Anomaly struct (warehouse_id, sku_id, type, severity, z_score, detected_at, status)
   - Type: "unusual_stock_drop", "unusual_spike", "unusual_frequency"
   - Severity: "low", "medium", "high", "critical"
   - Status: "open", "acknowledged", "resolved"
2. Create `internal/repository/anomaly_repository.go`:
   - CreateAnomaly(details)
   - GetAnomalyByID(id)
   - ListAnomalies(filters) → paginated
   - UpdateAnomalyStatus(id, status, notes)
3. Create `pkg/timeseries/anomaly_detector.go`:
   - CalculateZScore(latest_value, values_90days)
     - Formula: Z = (value - mean) / std_dev
   - GetSeverityLevel(z_score)
     - Z > 4: CRITICAL, Z > 3: HIGH, Z > 2: MEDIUM, Z > 1.5: LOW
   - RefreshBaseline(warehouse_id, sku_id) → weekly job
4. No async job yet - just setup infrastructure

**Deliverables:**
- ✅ Anomaly model & repository ready
- ✅ Z-score calculation tested
- ✅ Severity thresholds verified (TDD Section 6.4.2)
- ✅ Baseline calculation working

**Reference:**
- TDD Section 4.8: Anomalies Table
- TDD Section 6.4.1-6.4.2: Anomaly Detection & Z-Score

---

### Week 4: Forecasting & Webhooks

#### Week 4.1: Forecast Model & Repository (5 hours)

**Tasks:**
1. Create `internal/model/forecast.go`:
   - Forecast struct (warehouse_id, sku_id, forecast_date, predicted_demand, lower_bound, upper_bound, confidence, algorithm, model_version, accuracy_metrics JSONB)
   - ForecastAccuracy struct (warehouse_id, sku_id, algorithm, mape, rmse, mae, evaluation_date)
2. Create `internal/repository/forecast_repository.go`:
   - CreateForecast(forecast data)
   - GetForecast(warehouse_id, sku_id, forecast_date)
   - ListForecasts(warehouse_id, sku_id, horizon_days) → 30-day ahead
   - CreateAccuracyMetrics(warehouse_id, sku_id, algorithm, mape, rmse, mae)
   - GetBestAlgorithm(warehouse_id, sku_id) → lowest MAPE
3. Create Redis caching for forecasts:
   - Cache key: `forecast:{warehouse_id}:{sku_id}`
   - TTL: 24 hours (refresh daily @ 2 AM UTC)

**Deliverables:**
- ✅ Forecast model & repository complete
- ✅ Accuracy metrics stored in JSONB
- ✅ Model versioning supported (ets_v1, arima_v1)
- ✅ Cache TTL: 24 hours

**Reference:**
- TDD Section 4.6-4.7: Forecasts & Accuracy tables
- TDD Section 6.3.1: Forecast Schedule (2 AM UTC)

---

#### Week 4.2: Forecasting Algorithms (8 hours)

**Tasks:**
1. Create `pkg/timeseries/exponential_smoothing.go`:
   - Simple exponential smoothing formula
   - Alpha = 0.3 (smoothing parameter)
   - Return: predicted_demand, lower_bound, upper_bound (95% confidence)
   - Test with sample data
2. Create `pkg/timeseries/arima.go`:
   - ARIMA(3,1,2) implementation
   - Or use external package: `github.com/go-echarts/go-echarts` or simple AR implementation
   - Return: predicted_demand, lower_bound, upper_bound
   - Test with sample data
3. Create `pkg/timeseries/accuracy_metrics.go`:
   - CalculateMAE(predicted, actual) → Mean Absolute Error
   - CalculateRMSE(predicted, actual) → Root Mean Squared Error
   - CalculateMAPE(predicted, actual) → Mean Absolute Percentage Error
   - Test accuracy calculations
4. Create `internal/service/forecast_service.go`:
   - GenerateForecast(warehouse_id, sku_id)
     - Fetch 6-12 months historical data
     - Preprocess: remove outliers, interpolate, normalize
     - Run ETS algorithm
     - Run ARIMA algorithm
     - Calculate confidence intervals (95%)
     - Evaluate accuracy on last 30 days
     - Store forecasts & metrics in DB
     - Cache in Redis
     - Publish webhook "forecast.generated"
   - GetForecast(warehouse_id, sku_id, horizon)
     - Check cache first (24h TTL)
     - Select best algorithm (lowest MAPE)
     - Return forecast with confidence intervals

**Deliverables:**
- ✅ Both ETS & ARIMA algorithms implemented & tested
- ✅ Accuracy metrics calculated correctly
- ✅ Forecast service fully functional
- ✅ Algorithm comparison working (best = lowest MAPE)

**Reference:**
- TDD Section 6.3: Forecasting Flow (full description)
- LOGIC_FLOW Section 5: Forecasting Flow

---

#### Week 4.3: Webhook System (6 hours)

**Tasks:**
1. Create `internal/model/webhook.go`:
   - WebhookSubscription (subscription_id, event_type, webhook_url, secret_key, is_active)
   - WebhookEvent (event_id, event_type, timestamp, data)
   - WebhookDelivery (delivery_id, subscription_id, event_id, attempt, status, http_status, error_message)
2. Create `internal/repository/webhook_repository.go`:
   - CreateSubscription(event_type, webhook_url, secret_key)
     - Generate secret_key (32 bytes, return once only)
   - GetSubscriptions(event_type) → active only
   - UpdateSubscriptionStatus(subscription_id, is_active)
   - DeleteSubscription(subscription_id)
   - LogDelivery(subscription_id, event_id, attempt, status, http_status)
3. Create `internal/service/webhook_service.go`:
   - PublishEvent(event_type, data)
     - Fetch all active subscriptions for event_type
     - Build payload JSON
     - Generate HMAC-SHA256 signature (using secret_key)
     - For each subscriber: async delivery (goroutine)
   - DeliverEvent(subscription, payload, signature)
     - Retry logic: 1s, 2s, 4s, 8s (5 attempts total)
     - HTTP timeout: 10 seconds per attempt
     - Log delivery status
     - On success: mark as "success"
     - On failure after 5 attempts: move to dead-letter queue
4. Create `internal/handler/webhook_handler.go`:
   - POST /api/v1/webhooks/subscribe
     - Request: {event_type, webhook_url}
     - Response: {subscription_id, event_type, webhook_url, secret_key}
   - GET /api/v1/webhooks/subscriptions
     - List active subscriptions
   - DELETE /api/v1/webhooks/{subscription_id}
     - Unsubscribe
   - GET /api/v1/webhooks/deliveries
     - View delivery history (admin only)

**Deliverables:**
- ✅ Webhook endpoints working
- ✅ Secret key generation & security verified
- ✅ HMAC-SHA256 signature generation tested
- ✅ Retry logic tested (1-2-4-8s backoff)
- ✅ Async delivery working (goroutines)

**Reference:**
- TDD Section 5.5.3: Webhook Delivery Configuration
- LOGIC_FLOW Section 7: Webhook Flow

---

#### Week 4.4: Rate Limiting & API Middleware (4 hours)

**Tasks:**
1. Create `internal/middleware/rate_limit_middleware.go`:
   - Sliding window algorithm (1-second rolling window)
   - Limits: API key 1000/s, user 100/s, IP 10/s
   - Use Redis sorted sets (ZRANGE + ZADD)
   - Algorithm: Count requests in last 1 second
   - Return 429 Too Many Requests if exceeded
   - Add response headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
2. Create `internal/middleware/error_recovery_middleware.go`:
   - Panic recovery
   - Convert panic to 500 error response
3. Create `internal/middleware/cors_middleware.go`:
   - Allow cross-origin requests
   - Configurable origin whitelist
4. Create `internal/handler/health_handler.go`:
   - GET /health → {status: "live"} (Kubernetes liveness probe)
   - GET /ready → check DB + Redis (readiness probe)

**Deliverables:**
- ✅ Rate limiting working (all 3 dimensions: API key, user, IP)
- ✅ Response headers correct (TDD Section 7.1.2)
- ✅ Panic recovery tested
- ✅ CORS configured
- ✅ Health checks working

**Reference:**
- TDD Section 7.1: Rate Limiting Strategy

---

### Week 3-4 Deliverables Summary

**Code Commits:**
- [ ] Commit 8: Inventory service + repository + endpoints
- [ ] Commit 9: Anomaly detection model + repository + Z-score
- [ ] Commit 10: Forecast model + repository
- [ ] Commit 11: ETS + ARIMA algorithms + accuracy metrics
- [ ] Commit 12: Forecast service + endpoints
- [ ] Commit 13: Webhook system (subscribe, deliver, retry)
- [ ] Commit 14: Rate limiting + middleware + health checks

**Testing:**
- [ ] >80% unit test coverage for service layer
- [ ] Integration tests for inventory transactions (deadlock handling)
- [ ] Integration tests for forecasting algorithm accuracy
- [ ] Integration tests for webhook delivery & retry
- [ ] Load tests for rate limiting (spike 2000 req/s → should throttle)

**API Verification:**
- [ ] 15 endpoints total working (5 inventory, 4 forecast, 4 webhook, 2 health)
- [ ] All error responses match TDD Section 5.1 format
- [ ] All responses include timestamp in UTC (TDD Section 5 - Timezone)
- [ ] All protected endpoints require JWT

**Performance Baseline:**
- [ ] Inventory query: p99 <100ms (with cache hit)
- [ ] Forecast batch: <5 seconds for 1000 SKUs (not implemented yet)
- [ ] Concurrent users: 1000+ (connection pool tested)
- [ ] Rate limiting: handles 1000/s on API key limit

**Ready for ROADMAP_3?** YES - All backend APIs complete, ready for frontend

---

## ROADMAP_3: FRONTEND/MOBILE (WEEKS 5-8)

### Phase Goals

Build user-facing dashboard with real-time inventory, forecasts, anomalies. Requires backend from ROADMAP_1 + 2.

### Stack (from TDD.md)

- **Framework:** Next.js 14+ (App Router, TypeScript)
- **UI Library:** React 18+
- **Styling:** Tailwind CSS v3+
- **State Management:** Context API + hooks (useReducer)
- **API Client:** Axios
- **Package Manager:** pnpm

### Week 5-6: Dashboard & Inventory UI

#### Week 5.1: Frontend Project Setup (4 hours)

**Tasks:**
1. Create Next.js app with TypeScript:
   - `npx create-next-app@14 scinel-frontend --typescript --tailwind`
2. Setup folder structure:
   - app/ (Next.js pages)
   - components/ (React components)
   - lib/ (utilities, API client, auth)
   - types/ (TypeScript interfaces)
   - styles/ (Tailwind config)
   - public/ (static assets)
3. Create environment setup:
   - .env.local with API_URL, etc.
   - .env.example
4. Setup authentication context:
   - `lib/auth.ts` → AuthContext, useAuth hook
   - Store access_token, refresh_token, user data
   - Token refresh on expiry
5. Create API client:
   - `lib/api.ts` → Axios instance with JWT interceptor
   - Auto-add Authorization header
   - Handle 401 (token refresh)

**Deliverables:**
- ✅ Next.js app running locally
- ✅ Auth context & hooks working
- ✅ API client with JWT interceptor
- ✅ TypeScript strict mode enabled

**Reference:**
- TDD Section 2: Frontend Technology Stack
- HALAMAN Section: All pages requiring authentication

---

#### Week 5.2: Authentication Pages (6 hours)

**Tasks:**
1. Create auth pages:
   - `app/(auth)/login/page.tsx` → Login form
   - `app/(auth)/register/page.tsx` → Registration form
   - `app/(auth)/forgot-password/page.tsx` → Send reset email
   - `app/(auth)/reset-password/page.tsx` → Set new password
2. Create `components/auth/LoginForm.tsx`:
   - Form: email, password
   - Validation: email format, password strength
   - Submit: POST /api/v1/auth/login
   - Response: store access_token + refresh_token
   - Redirect: /dashboard
   - Error: show error message (from TDD Section 5.1 error format)
3. Create `components/auth/RegisterForm.tsx`:
   - Form: email, password, full_name
   - Validation: same as login
   - Submit: POST /api/v1/auth/register
   - Response: auto-login or redirect to login
4. Create password reset pages (bonus week 5.3)
5. Styling: Tailwind CSS (dark mode ready)

**Deliverables:**
- ✅ All 4 auth pages functional
- ✅ Form validation working
- ✅ Error messages displayed
- ✅ Tokens stored in localStorage (or httpOnly cookie)
- ✅ Redirect after login working

**Reference:**
- HALAMAN Section: Auth pages (login, register, etc.)

---

#### Week 5.3: Inventory Dashboard (8 hours)

**Tasks:**
1. Create `app/dashboard/layout.tsx`:
   - Navbar with logout
   - Sidebar with navigation
   - User profile dropdown
   - Protected route (redirect to login if no auth)
2. Create `app/dashboard/page.tsx` (main dashboard):
   - Summary cards: total warehouses, total SKUs, total stock value
   - Key metrics: real-time stock, forecast accuracy, anomalies today
   - Charts (using recharts):
     - Inventory trend (last 30 days)
     - Top SKUs by value
     - Anomalies this month
3. Create `app/dashboard/inventory/page.tsx`:
   - Table: warehouse inventory
   - Columns: Warehouse, SKU, Quantity, Status, Last Updated
   - Filters: warehouse, SKU search, quantity range
   - Pagination: 50 items/page
   - Actions: View details, record transaction
4. Create `app/dashboard/inventory/[id]/page.tsx`:
   - Inventory detail page
   - SKU information
   - Stock levels across all warehouses
   - Transaction history (last 100)
   - Charts: stock trend (90 days)
5. Create transaction form (`components/inventory/TransactionForm.tsx`):
   - Type: in/out
   - Quantity
   - Reason
   - Submit: POST /api/v1/transactions
   - Success: refresh inventory data

**Deliverables:**
- ✅ Main dashboard with metrics
- ✅ Inventory table with filters & pagination
- ✅ Inventory detail page with history
- ✅ Transaction form working
- ✅ Real-time data updates (30s polling or WebSocket ready)
- ✅ All timestamps displayed in user's local timezone (Intl.DateTimeFormat)

**Reference:**
- HALAMAN Section 5-8: Inventory pages
- TDD Section 5: Timezone Handling (client-side display)

---

### Week 6: Forecasts & Anomalies UI

#### Week 6.1: Forecast Visualization (8 hours)

**Tasks:**
1. Create `app/dashboard/forecast/page.tsx`:
   - Select warehouse & SKU
   - Display 30-day forecast chart (recharts LineChart)
   - Show: predicted_demand, lower_bound, upper_bound (95% confidence)
   - Confidence band (shaded area between lower/upper)
   - Accuracy metrics: MAPE, RMSE, MAE
   - Algorithm comparison: show both ETS & ARIMA results
2. Create forecast detail component (`components/forecast/ForecastDetail.tsx`):
   - Show forecast table (date, prediction, confidence)
   - Algorithm used (ets_v1 or arima_v1)
   - Model version selector (if multiple versions available)
   - Historical accuracy trend
3. Create forecast accuracy component (`components/forecast/AccuracyMetrics.tsx`):
   - Display MAPE, RMSE, MAE
   - Target: MAPE <15% (green=good, red=poor)
   - Comparison vs. previous forecast
4. Styling: Clean, professional charts

**Deliverables:**
- ✅ Forecast page fully functional
- ✅ Charts rendering correctly with recharts
- ✅ Algorithm comparison visible
- ✅ Accuracy metrics displayed
- ✅ Timezone handling for forecast dates

**Reference:**
- HALAMAN Section 9-10: Forecast pages
- TDD Section 6.3-6.4: Forecast algorithms

---

#### Week 6.2: Anomalies & Alerts UI (8 hours)

**Tasks:**
1. Create `app/dashboard/anomalies/page.tsx`:
   - Table: detected anomalies
   - Columns: Warehouse, SKU, Type, Severity, Detected At, Status
   - Filters: severity, status, date range
   - Pagination: 50 items/page
   - Color coding: CRITICAL=red, HIGH=orange, MEDIUM=yellow, LOW=gray
2. Create anomaly detail page (`app/dashboard/anomalies/[id]/page.tsx`):
   - Full anomaly information
   - Z-score calculation & explanation
   - Historical context (last 30 days data)
   - Resolution form (status, notes)
   - Charts: transaction history with anomaly highlighted
3. Create alert system (`components/anomaly/AnomalyAlert.tsx`):
   - Real-time alerts (polling or WebSocket)
   - Toast notifications for new anomalies (severity: CRITICAL, HIGH)
   - Dismiss/snooze options
4. Create resolution form (`components/anomaly/ResolutionForm.tsx`):
   - Status dropdown: open → acknowledged → resolved
   - Notes textarea
   - Submit: PATCH /api/v1/anomalies/{id}

**Deliverables:**
- ✅ Anomalies page with severity color coding
- ✅ Detail page with Z-score explanation
- ✅ Alert system working
- ✅ Resolution workflow functional
- ✅ Real-time updates (at least 5s polling)

**Reference:**
- HALAMAN Section 11-12: Anomalies pages
- TDD Section 6.4: Anomaly Detection with Z-score

---

### Week 7: Real-time Features & Mobile

#### Week 7.1: Real-Time Updates (4 hours)

**Tasks:**
1. Implement polling-based updates (if WebSocket not ready):
   - useEffect hook: fetch inventory data every 30 seconds
   - useEffect hook: fetch forecasts every 24 hours (or on demand)
   - useEffect hook: fetch anomalies every 5 seconds
2. Setup Webhook integration (manual for now):
   - Show webhook subscriptions page (read-only)
   - Display delivery history
   - Manual refresh button
3. Optional: WebSocket setup (Week 7.2 if time):
   - Socket.io or native WebSocket
   - Real-time inventory updates
   - Real-time anomaly alerts

**Deliverables:**
- ✅ Polling working (30s inventory, 5s anomalies)
- ✅ Data freshness acceptable
- ✅ No excessive API calls
- ✅ Graceful error handling (network failures)

**Reference:**
- TDD Section 5.5: Webhook Endpoints (display only)

---

#### Week 7.2: Mobile Responsive (4 hours)

**Tasks:**
1. Responsive design:
   - Tailwind breakpoints: sm, md, lg, xl
   - Mobile-first approach
   - Navigation: hamburger menu on mobile
   - Tables: collapse to cards on mobile
   - Charts: responsive sizing
2. Test on devices:
   - Desktop (1920x1080)
   - Tablet (768x1024)
   - Mobile (375x812)
3. Performance optimization:
   - Image optimization (next/image)
   - Code splitting (dynamic imports)
   - Lazy loading

**Deliverables:**
- ✅ All pages responsive
- ✅ Mobile navigation working
- ✅ Charts readable on mobile
- ✅ Performance: Lighthouse score >90

**Reference:**
- TDD Section 2: Frontend Stack (Tailwind responsive)

---

### Week 8: Polish & Optimization

#### Week 8.1: UI Polish (4 hours)

**Tasks:**
1. Dark mode toggle
2. Consistent color scheme (brand colors)
3. Icon library (lucide-react or similar)
4. Loading states & spinners
5. Empty states (no data → helpful message)
6. Error boundaries

**Deliverables:**
- ✅ Dark mode working
- ✅ All UI components polished
- ✅ Loading indicators visible
- ✅ Error messages helpful

---

#### Week 8.2: Performance & Testing (4 hours)

**Tasks:**
1. Performance optimization:
   - Lighthouse audit
   - Remove unused dependencies
   - Optimize bundle size
2. Unit tests:
   - API client tests (mocking)
   - Component tests (React Testing Library)
   - Utility function tests
3. E2E tests (optional):
   - Login flow
   - View inventory
   - Create transaction
   - View forecasts

**Deliverables:**
- ✅ Lighthouse score >90
- ✅ >50% test coverage
- ✅ Bundle size <200KB (gzipped)
- ✅ All major flows E2E tested

---

### Week 5-8 Deliverables Summary

**Frontend Commits:**
- [ ] Commit 15: Next.js setup + auth context
- [ ] Commit 16: Auth pages (login, register, reset)
- [ ] Commit 17: Dashboard + inventory table + detail
- [ ] Commit 18: Transaction form + create transaction flow
- [ ] Commit 19: Forecast visualization + accuracy metrics
- [ ] Commit 20: Anomalies page + alerts + resolution
- [ ] Commit 21: Real-time updates + mobile responsive
- [ ] Commit 22: Dark mode + polish + performance

**Pages Created:** 15+ pages total
- Auth: login, register, forgot-password, reset-password
- Dashboard: main, inventory, inventory/:id, transactions
- Forecasts: list, detail
- Anomalies: list, detail
- Settings: profile, webhooks, API keys
- Admin: user management (if admin role)

**Components Created:** 30+ components
- Forms: LoginForm, RegisterForm, TransactionForm, ResolutionForm
- Tables: InventoryTable, TransactionTable, AnomalyTable
- Charts: InventoryTrend, ForecastChart, AnomalyChart
- Common: Navbar, Sidebar, Header, Footer
- UI: Button, Input, Modal, Dropdown, Toast, Spinner

**Testing:**
- [ ] >50% unit test coverage
- [ ] Key user flows E2E tested
- [ ] Mobile responsive verified
- [ ] Dark mode tested

**Performance:**
- [ ] Lighthouse score >90
- [ ] Page load <3s
- [ ] API response time <100ms (average)

**Ready for ROADMAP_4?** YES - All features complete, ready for final polish

---

## ROADMAP_4: MASTERY/DEPLOYMENT (WEEK 9)

### Phase Goals

Final testing, documentation, deployment setup, and project completion. Bridge from dev → production.

### Week 9: Testing, Docs, Deployment

#### Week 9.1: Comprehensive Testing (6 hours)

**Backend Testing:**

1. Unit Tests (increase coverage to >90%):
   - [ ] All service layer methods
   - [ ] All repository methods
   - [ ] Error handling for all error codes (TDD 5.1)
   - [ ] Deadlock retry logic
   - [ ] Rate limiting algorithm
   - [ ] Z-score calculation
   - [ ] Forecast algorithms (ETS, ARIMA)
   - [ ] Webhook retry logic

2. Integration Tests:
   - [ ] Auth flow: register → login → logout
   - [ ] Inventory flow: create transaction → update stock → check cache
   - [ ] Forecast flow: generate forecast → verify accuracy metrics
   - [ ] Anomaly flow: create transaction → detect anomaly → send alert
   - [ ] Webhook flow: publish event → deliver to subscriber → retry on failure
   - [ ] Rate limiting: exceed limit → get 429 → headers correct

3. Load Tests:
   - [ ] Inventory queries: 1000 concurrent users
   - [ ] Forecast batch: process 1000 SKUs in <5 seconds
   - [ ] Rate limiting: handle 2000 req/s (should throttle at limits)
   - [ ] Database connection pool: verify no exhaustion

**Frontend Testing:**

1. Component Tests:
   - [ ] All form validations
   - [ ] All chart renders
   - [ ] Error message display
   - [ ] Loading states

2. E2E Tests (Cypress or Playwright):
   - [ ] Login flow
   - [ ] Create transaction
   - [ ] View forecast chart
   - [ ] Acknowledge anomaly
   - [ ] Refresh token on expiry

3. Performance Tests:
   - [ ] Lighthouse score >90
   - [ ] Page load time <3s
   - [ ] Bundle size gzip <200KB

4. Responsive Design:
   - [ ] Desktop (1920x1080)
   - [ ] Tablet (768x1024)
   - [ ] Mobile (375x812)

**Deliverables:**
- ✅ Backend: >90% unit test coverage
- ✅ Backend: all integration tests passing
- ✅ Backend: load tests passing (capacity verified)
- ✅ Frontend: >50% component test coverage
- ✅ Frontend: 5+ E2E test flows passing
- ✅ All tests passing in CI/CD (GitHub Actions)

**Reference:**
- TDD Section 10: Testing Strategy

---

#### Week 9.2: Documentation (4 hours)

**Deliverables:**

1. README.md (project root):
   - [ ] Project description (1 paragraph)
   - [ ] Features overview
   - [ ] Tech stack
   - [ ] Setup instructions (Docker Compose)
   - [ ] How to run tests
   - [ ] How to deploy

2. docs/API.md:
   - [ ] All 20+ API endpoints documented
   - [ ] Request/response examples
   - [ ] Error codes & meanings (TDD 5.1)
   - [ ] Authentication flow
   - [ ] Rate limiting info
   - [ ] Webhook guide

3. docs/ARCHITECTURE.md:
   - [ ] Architecture pattern (Clean Architecture)
   - [ ] Layer descriptions (handler, service, repository)
   - [ ] Database schema overview
   - [ ] Concurrency model (goroutines, channels)
   - [ ] Caching strategy (Redis)
   - [ ] Logging strategy

4. docs/DEPLOYMENT.md:
   - [ ] Docker build & push
   - [ ] Environment variables (production)
   - [ ] Database migration steps
   - [ ] Health checks (liveness, readiness)
   - [ ] Monitoring & alerting setup
   - [ ] Scaling considerations
   - [ ] Backup strategy

5. docs/FORECASTING_ALGORITHMS.md:
   - [ ] ETS algorithm explanation
   - [ ] ARIMA algorithm explanation
   - [ ] Accuracy metrics (MAPE, RMSE, MAE)
   - [ ] Model versioning strategy
   - [ ] Baseline calculation (90-day window)
   - [ ] Refresh schedule (daily @ 2 AM UTC)

6. Code comments:
   - [ ] Complex algorithms commented
   - [ ] Public functions documented (godoc)
   - [ ] Non-obvious business logic explained

**Deliverables:**
- ✅ All 6 documentation files complete
- ✅ Examples provided for all endpoints
- ✅ Diagrams for architecture
- ✅ Step-by-step deployment guide

**Reference:**
- TDD: All sections (reference for API & design)

---

#### Week 9.3: Deployment Setup (4 hours)

**Backend Deployment:**

1. Docker:
   - [ ] Dockerfile multi-stage build
   - [ ] Final image size <150MB
   - [ ] Healthcheck endpoint (/health)
   - [ ] Graceful shutdown (30s timeout)

2. Environment Configuration:
   - [ ] .env template with all variables
   - [ ] Secrets management (use environment variables, not .env in prod)
   - [ ] Database URL, Redis URL, JWT secret

3. Database Migrations:
   - [ ] Auto-run migrations on startup (or manual via script)
   - [ ] Rollback strategy documented
   - [ ] Backup before deploy

4. Deployment Targets (choose one for MVP):
   - [ ] Option A: Render.com (free tier, easy setup)
   - [ ] Option B: Railway (free tier, Git auto-deploy)
   - [ ] Option C: AWS EC2 + systemd (manual, full control)
   - [ ] Option D: Docker Swarm or Kubernetes (advanced)

5. CI/CD (GitHub Actions):
   - [ ] Lint on push
   - [ ] Tests on push
   - [ ] Build Docker image on merge to main
   - [ ] Push to registry (ghcr.io)
   - [ ] Deploy to staging on successful build
   - [ ] Manual approval for production

**Frontend Deployment:**

1. Build & Optimization:
   - [ ] Next.js build: `npm run build`
   - [ ] Output: out/ (static export) or node server
   - [ ] Bundle size: <200KB gzipped

2. Deployment Targets:
   - [ ] Option A: Vercel (easiest, free tier)
   - [ ] Option B: Netlify (free tier)
   - [ ] Option C: AWS S3 + CloudFront (CDN)
   - [ ] Option D: Self-hosted (on same server as backend)

3. Environment Variables:
   - [ ] NEXT_PUBLIC_API_URL (backend URL)
   - [ ] Different URLs for dev/staging/prod

4. CI/CD:
   - [ ] Build on push
   - [ ] Tests on push
   - [ ] Deploy on merge to main

**Deliverables:**
- ✅ Backend: Dockerfile ready for production
- ✅ Backend: Deployed to staging environment
- ✅ Frontend: Built & deployed to production (or staging)
- ✅ Both: Health checks working
- ✅ Both: CI/CD pipeline working end-to-end

**Reference:**
- TDD Section 11: Deployment & Environment

---

#### Week 9.4: Final Polish & Launch (6 hours)

**Code Quality:**
- [ ] Zero console errors/warnings
- [ ] Linter passing (golangci-lint, eslint)
- [ ] No dead code
- [ ] No TODOs without tickets
- [ ] All secrets removed from repo

**Performance Verification:**
- [ ] Backend: p99 latency <100ms (inventory query)
- [ ] Backend: forecast batch <5s for 1000 SKUs
- [ ] Frontend: Lighthouse score >90
- [ ] Frontend: page load <3s

**Security:**
- [ ] HTTPS enabled (prod)
- [ ] JWT secrets strong (32+ chars)
- [ ] No SQL injection (use parameterized queries)
- [ ] CORS configured (prod URLs only)
- [ ] Rate limiting active
- [ ] Input validation on all endpoints
- [ ] No hardcoded secrets in code

**Monitoring:**
- [ ] Logging: structured JSON logs
- [ ] Error tracking: Sentry or similar (optional)
- [ ] Metrics: Prometheus or DataDog (optional)
- [ ] Alerts: Email on errors >1% rate

**Portfolio Materials:**
- [ ] README.md with badge links (GitHub, demo)
- [ ] Deployment screenshot
- [ ] Architecture diagram
- [ ] API documentation screenshot
- [ ] Performance metrics screenshot (Lighthouse)

**Deliverables:**
- ✅ Production environment live & stable
- ✅ All tests passing (CI/CD green)
- ✅ Documentation complete & polished
- ✅ No console errors/warnings
- ✅ Security checklist passed
- ✅ Performance targets met
- ✅ Portfolio-ready materials prepared

---

### Week 9 Deliverables Summary

**Final Commits:**
- [ ] Commit 23: Comprehensive test suite (unit + integration + E2E)
- [ ] Commit 24: Complete documentation (6 docs)
- [ ] Commit 25: Docker + CI/CD + deployment config
- [ ] Commit 26: Security hardening + final polish
- [ ] Tag: v1.0.0 (release tag)

**Testing Results:**
- [ ] Backend unit tests: >90% coverage
- [ ] Frontend unit tests: >50% coverage
- [ ] Integration tests: all flows passing
- [ ] E2E tests: 5+ critical user flows
- [ ] Load tests: capacity verified
- [ ] All tests passing in CI/CD

**Documentation:**
- [ ] README.md (setup & features)
- [ ] API.md (20+ endpoints documented)
- [ ] ARCHITECTURE.md (system design)
- [ ] DEPLOYMENT.md (step-by-step)
- [ ] FORECASTING_ALGORITHMS.md (algorithm details)
- [ ] Code comments (complex logic)

**Deployment:**
- [ ] Docker image built & tested
- [ ] GitHub Actions CI/CD working
- [ ] Staging environment live
- [ ] Production environment live (or ready)
- [ ] Database migrated & healthy
- [ ] Health checks passing

**Performance Verified:**
- [ ] Backend p99: <100ms (inventory query)
- [ ] Backend: <5s for 1000 SKU forecast batch
- [ ] Frontend: Lighthouse >90
- [ ] Frontend: page load <3s

**Portfolio Quality:**
- [ ] Code quality: high (linter passing, no dead code)
- [ ] Documentation: professional-grade
- [ ] README: compelling for recruiters
- [ ] GitHub repo: portfolio-ready
- [ ] Live demo: working & impressive

**Ready for Production?** YES - All systems go!

---

## CROSS-PHASE DEPENDENCIES

### Dependency Graph

```
ROADMAP_1 (Foundation)
    ↓ (DB, Auth ready)
ROADMAP_2 (Backend Core)
    ↓ (All APIs ready)
ROADMAP_3 (Frontend)
    ↓ (UI complete)
ROADMAP_4 (Testing & Deploy)
    → (Production live)
```

### Phase Blockers

- **ROADMAP_2 blocked by ROADMAP_1:** Cannot implement inventory/forecast APIs without auth & DB
- **ROADMAP_3 blocked by ROADMAP_2:** Cannot build UI without working backend APIs
- **ROADMAP_4 blocked by ROADMAP_3:** Cannot deploy without both backend & frontend complete

### Parallel Work Possible

- **Within each phase:** Multiple engineers can work on different services (e.g., inventory vs. forecasting in ROADMAP_2)
- **ROADMAP_3 frontend:** Can start week 5 while ROADMAP_2 backend finishes week 4
- **ROADMAP_4 testing:** Unit tests can start week 3, integration tests after week 4

### Critical Path

**Minimal path to MVP (without optional features):**
1. ROADMAP_1: Weeks 1-2 (critical)
2. ROADMAP_2: Weeks 3-4 (critical)
3. ROADMAP_3 (Week 5-6): Dashboard + Inventory (critical)
4. ROADMAP_4 (Week 7-9): Testing + Deployment (critical)

**Optional enhancements (if time):**
- Week 7: WebSocket real-time (instead of polling)
- Week 7-8: Advanced forecasting UI (multiple models, detailed accuracy)
- Week 8: Admin dashboard (user management)
- Week 8: API key management (instead of user-based auth)

---

## VERIFICATION CHECKLIST

### Pre-ROADMAP_1 Start

- [ ] All team members read this blueprint
- [ ] GitHub repo created & access granted
- [ ] Docker Desktop installed locally
- [ ] Go 1.21+ installed & tested
- [ ] PostgreSQL 15+ available (via Docker Compose)
- [ ] Redis 7+ available (via Docker Compose)
- [ ] Team familiar with TDD.md, PRD.md, LOGIC_FLOW.md

### ROADMAP_1 Completion Checklist

- [ ] Project structure matches TDD Section 3
- [ ] All 10 migrations run successfully
- [ ] 4 auth endpoints tested & working
- [ ] RBAC middleware protecting sensitive routes
- [ ] Structured JSON logging working
- [ ] Docker Compose stack healthy
- [ ] GitHub Actions CI passing
- [ ] >80% test coverage for auth & DB layers
- [ ] README.md with setup instructions

### ROADMAP_2 Completion Checklist

- [ ] 5 inventory endpoints working
- [ ] Forecast batch job schedulable (not yet running)
- [ ] Anomaly detection model complete
- [ ] Webhook subscribe/deliver working
- [ ] Rate limiting middleware active
- [ ] All 4 database tables (inventory, forecast, anomaly, webhook) populated
- [ ] >80% test coverage for service layer
- [ ] Redis caching verified (30s inventory, 24h forecast)
- [ ] All 15 backend endpoints documented & tested

### ROADMAP_3 Completion Checklist

- [ ] 15+ pages created & functional
- [ ] 30+ components built & tested
- [ ] Dashboard showing real-time inventory
- [ ] Forecast charts rendering correctly
- [ ] Anomaly alerts working
- [ ] Mobile responsive (tested on 3 devices)
- [ ] Dark mode toggle working
- [ ] >50% component test coverage
- [ ] 5+ E2E user flows passing
- [ ] Lighthouse score >90

### ROADMAP_4 Completion Checklist

- [ ] Backend: >90% unit test coverage
- [ ] Backend: all integration tests passing
- [ ] Frontend: >50% component test coverage
- [ ] Frontend: E2E tests passing
- [ ] Load tests: verified capacity (1000 concurrent users)
- [ ] 6 documentation files complete
- [ ] Docker image built (<150MB)
- [ ] CI/CD pipeline working (lint → test → build → deploy)
- [ ] Staged deployment working
- [ ] Production deployment ready
- [ ] All performance targets met
- [ ] Security checklist passed
- [ ] Portfolio materials prepared

---

## FILE REFERENCE MAP

**For each ROADMAP file, reference these source documents:**

| ROADMAP | Key References | Section Numbers |
|---------|----------------|-----------------|
| ROADMAP_1 | TDD, PRD, LOGIC_FLOW | TDD Sec 3, 4.0-4.12, 5.1; PRD Sec 5-6; LOGIC_FLOW Sec 2 |
| ROADMAP_2 | TDD, PRD, LOGIC_FLOW, DRP | TDD Sec 4.4-4.10, 5.2-5.5, 6.3-6.4, 7.1; LOGIC_FLOW Sec 4-7; DRP Sec 2-3 |
| ROADMAP_3 | HALAMAN, TDD Sec 2 | HALAMAN Sec 5-12; TDD Sec 2 (Frontend Stack); TDD Sec 5 (Timezone) |
| ROADMAP_4 | TDD, STP, DRP | TDD Sec 10-11; STP Sec all; DRP Sec 1, 4 |

---

## NOTES FOR CLAUDE WHEN CREATING ROADMAP FILES

**Important constraints to avoid hallucination:**

1. **Tech Stack (LOCKED):**
   - Backend: Go 1.21+, no heavy frameworks
   - Frontend: Next.js 14+, React 18+, Tailwind CSS
   - Database: PostgreSQL 15+, Redis 7+
   - Do NOT invent new tech (no Node.js backend, no Vue.js frontend, no MongoDB)

2. **Timeframe (LOCKED):**
   - Total: 9 weeks, ~160 hours
   - ROADMAP_1: Weeks 1-2 (40h)
   - ROADMAP_2: Weeks 3-4 (40h)
   - ROADMAP_3: Weeks 5-8 (60h)
   - ROADMAP_4: Week 9 (20h)
   - Do NOT add extra phases or extend timeline

3. **Scope (LOCKED from PRD):**
   - IN SCOPE: Auth, Inventory, Forecasting (ETS + ARIMA), Anomaly Detection (Z-score), Webhooks
   - OUT OF SCOPE: OAuth, 2FA, ML models (LSTM, Isolation Forest), Kafka, GraphQL
   - Do NOT add scope creep

4. **API Endpoint Count (LOCKED at 20+):**
   - Auth: 4 endpoints
   - Inventory: 5 endpoints
   - Forecasting: 4 endpoints
   - Anomalies: 4 endpoints
   - Webhooks: 4 endpoints
   - Health: 2 endpoints
   - Total: 23 endpoints
   - Do NOT create extra endpoints not in TDD

5. **Database Tables (LOCKED at 12):**
   - users, warehouses, skus, inventory, transactions, forecasts, forecast_accuracy, anomalies, webhook_subscriptions, webhook_deliveries, audit_log, password_reset_tokens
   - Do NOT add extra tables

6. **Performance Targets (LOCKED):**
   - Inventory query: p99 <100ms
   - Forecast batch: <5 seconds for 1000 SKUs
   - Concurrent users: 1000+
   - API rate limit: 1000 req/sec per key
   - Do NOT claim different targets

7. **Test Coverage (LOCKED):**
   - Backend: >80% (ROADMAP_1-2), >90% (ROADMAP_4)
   - Frontend: >50% (ROADMAP_4)
   - Do NOT claim lower coverage as acceptable

8. **Deployment Targets (FLEXIBLE - pick one):**
   - Render.com (free, easiest)
   - Railway (free, Git auto-deploy)
   - AWS EC2 + systemd (manual)
   - Docker Swarm/Kubernetes (advanced)
   - But be consistent within same ROADMAP file

---

**END OF MASTER BLUEPRINT**

This blueprint is your **single source of truth** for all 4 ROADMAP files. 
Use it to generate consistent, hallucination-free ROADMAP files across separate chat sessions.

Generated: 2026-07-19  
Version: 1.0 (FINAL)  
Token Cost to Generate: ~120k tokens (worth it for accuracy!)

