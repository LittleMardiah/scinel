# 🧪 COMPREHENSIVE TEST PLAN - ScIntel Platform
**Supply Chain Intelligence Platform Testing Strategy**

**Status:** Ready for Implementation  
**Version:** 1.0  
**Last Updated:** 2026-07-29  
**Coverage Target:** 
- Backend: >90% unit coverage
- Frontend: >50% component coverage
- Integration: 5+ critical flows
- E2E: 5+ complete user journeys

---

## 📋 TEST OVERVIEW

Based on complete ROADMAP 1-4 analysis, this document outlines all MANDATORY tests required to achieve production-grade quality for ScIntel.

**Total Test Deliverables:**
- 40+ Backend Unit Test Files
- 9 Backend Integration Test Suites
- 50+ Frontend Unit Test Files
- 5+ E2E Test Scenarios
- 6 Load/Performance Test Suites

---

## 🔴 UNIT TESTS (BACKEND - GO)

### Test Coverage Targets
- **Auth Service:** 95% (registration, login, JWT, refresh, logout)
- **Inventory Service:** 92% (CRUD, stock levels, transactions, caching)
- **Forecast Service:** 90% (ETS/ARIMA algorithms, batch jobs, accuracy)
- **Anomaly Service:** 90% (Z-score, IQR detection, severity)
- **Webhook Service:** 88% (publishing, retry logic, dead-letter queue)
- **Middleware:** 95% (JWT, RBAC, rate limiting, CORS)
- **Repositories:** 90% (database queries, error handling)
- **Utilities:** 85% (validators, hash functions, formatters)

### Test Files to Create

#### 1️⃣ Authentication Tests (Roadmap 1.4 - JWT & RBAC)

**File:** `tests/unit/services/auth_service_test.go`

```
MUST TEST:
- ✅ User Registration (valid/invalid email, weak password, duplicate email)
- ✅ User Login (correct credentials, wrong password, non-existent user, locked account)
- ✅ JWT Token Generation (valid token, expiry time, refresh mechanism)
- ✅ Token Refresh (valid refresh token, expired refresh, invalid token)
- ✅ Logout (token invalidation, Redis cleanup)
- ✅ Password Hashing (bcrypt strength, salt variation)
- ✅ RBAC Authorization (4 roles: admin, manager, analyst, viewer)
- ✅ Permission Checking (role-based access control validation)
- ✅ Token Validation (signature verification, expiry check, format validation)
- ✅ Account Locking (after 5 failed attempts in 5 minutes)
```

**File:** `tests/unit/middleware/auth_middleware_test.go`

```
MUST TEST:
- ✅ JWT Extraction from Authorization header
- ✅ Token Signature Validation
- ✅ Token Expiry Check
- ✅ Invalid Token Handling (malformed, missing, expired)
- ✅ User Context Injection
- ✅ Authorization Header Format (Bearer token)
- ✅ Role-Based Access Control Enforcement
- ✅ API Key Authentication (alternative to JWT)
```

**File:** `tests/unit/middleware/rbac_middleware_test.go`

```
MUST TEST:
- ✅ Role Extraction from JWT Token
- ✅ Permission Matrix (admin→all, manager→warehouse, analyst→read, viewer→view)
- ✅ Endpoint Authorization (admin-only, manager-only routes)
- ✅ Warehouse Isolation (managers see only own warehouse)
- ✅ Unauthorized Access Rejection (403 Forbidden)
- ✅ Role Hierarchy Validation
```

---

#### 2️⃣ Inventory Service Tests (Roadmap 2.1-2.2 - Stock Management)

**File:** `tests/unit/services/inventory_service_test.go`

```
MUST TEST:
- ✅ Get Inventory by SKU (single item retrieval)
- ✅ Get Inventory by Warehouse (bulk retrieval)
- ✅ Update Stock Quantity (increment/decrement)
- ✅ Stock Level Warnings (normal/low/critical status)
- ✅ Low Stock Alerts (threshold triggers)
- ✅ Transaction Processing (in/out transactions)
- ✅ Concurrent Transaction Handling (deadlock prevention with SELECT FOR UPDATE)
- ✅ Cache Invalidation (after updates, 30s TTL)
- ✅ Inventory History Tracking
- ✅ Insufficient Stock Validation (reject negative quantities)
- ✅ Warehouse Assignment Validation
- ✅ SKU Code Validation (format, uniqueness)
```

**File:** `tests/unit/services/transaction_service_test.go`

```
MUST TEST:
- ✅ Create Transaction (in/out type)
- ✅ Transaction Type Validation
- ✅ Quantity Validation (positive, non-zero)
- ✅ Warehouse Existence Check
- ✅ SKU Existence Check
- ✅ Sufficient Stock Check (for "out" transactions)
- ✅ Pessimistic Locking (SELECT FOR UPDATE)
- ✅ Deadlock Retry Logic (3 attempts, exponential backoff)
- ✅ Transaction History Storage
- ✅ Reason & Details Logging
- ✅ Timestamp Recording (UTC timezone)
- ✅ Concurrent Transaction Ordering
- ✅ Rollback on Failure (ACID compliance)
```

**File:** `tests/unit/repositories/inventory_repository_test.go`

```
MUST TEST:
- ✅ Create Inventory Record
- ✅ Get by SKU
- ✅ Get by Warehouse
- ✅ Update Quantity
- ✅ Delete Inventory (cascade rules)
- ✅ Query with Filters (warehouse, status, date range)
- ✅ Sorting (by SKU, warehouse, quantity, timestamp)
- ✅ Pagination (limit, offset)
- ✅ Index Usage (fast queries <100ms)
- ✅ Foreign Key Validation
- ✅ Unique Constraint Checking (warehouse + SKU)
```

---

#### 3️⃣ Forecasting Tests (Roadmap 2.3 - ETS & ARIMA Algorithms)

**File:** `tests/unit/services/forecast_service_test.go`

```
MUST TEST:
- ✅ ETS Algorithm (Exponential Smoothing)
  - Correct parameter initialization
  - 180-day historical data requirement
  - 30-day forecast generation
  - Confidence interval calculation (95%, ±1.96 * std_error)
  - Accuracy metrics (MAPE, RMSE, MAE)
  
- ✅ ARIMA Algorithm (AutoRegressive Integrated Moving Average)
  - Auto parameter selection (p, d, q)
  - 180-day historical minimum validation
  - Stationarity checking (ADF test)
  - Forecast generation
  - Residual analysis
  
- ✅ Algorithm Comparison
  - MAPE-based selection (which algorithm is better)
  - Accuracy metric calculation
  - Confidence interval output
  
- ✅ Batch Job Processing
  - Generate forecasts for multiple SKUs
  - Batch timeout handling (<5s for 1000 SKUs)
  - Error recovery (partial failure handling)
  - SKU filtering (only active items)
  
- ✅ Data Validation
  - Minimum data points check (180 days)
  - Missing data handling (interpolation)
  - Outlier detection
  - Timezone consistency
  
- ✅ Cache Management
  - Store forecast in Redis (24h TTL)
  - Cache key generation (SKU-based)
  - Cache hit rate (target >60%)
```

**File:** `tests/unit/services/forecast_accuracy_test.go`

```
MUST TEST:
- ✅ MAPE Calculation (Mean Absolute Percentage Error)
- ✅ RMSE Calculation (Root Mean Square Error)
- ✅ MAE Calculation (Mean Absolute Error)
- ✅ Confidence Interval Width (95% default)
- ✅ Accuracy Trending (daily calculation)
- ✅ Algorithm Ranking (which performs best)
- ✅ Seasonal Decomposition (trend, seasonality, residuals)
- ✅ Forecast Recalibration (daily at 2 AM UTC)
```

---

#### 4️⃣ Anomaly Detection Tests (Roadmap 2.4 - Z-Score & IQR)

**File:** `tests/unit/services/anomaly_service_test.go`

```
MUST TEST:
- ✅ Z-Score Detection Method
  - Calculate Z-score for data point
  - Threshold setting (default 3.0)
  - Outlier identification
  - Mean & standard deviation calculation
  
- ✅ IQR (Interquartile Range) Method
  - Calculate Q1, Q2, Q3
  - IQR = Q3 - Q1
  - Lower/upper fence calculation (Q1-1.5×IQR, Q3+1.5×IQR)
  - Outlier identification
  
- ✅ Severity Classification
  - Critical (>3.0 Z-score)
  - High (2.0-3.0 Z-score)
  - Medium (1.5-2.0 Z-score)
  - Low (<1.5 Z-score)
  
- ✅ Anomaly Status Tracking
  - Create anomaly record (open status)
  - Update status (open→acknowledged→resolved)
  - Status timeline preservation
  - Timestamp accuracy
  
- ✅ Anomaly Acknowledgment
  - Mark as acknowledged
  - Add notes/comments
  - Update timestamp
  
- ✅ Historical Data Retrieval
  - Get anomaly history
  - Filter by date range
  - Status filtering
  - Sorting (newest/oldest)
  
- ✅ Concurrent Anomaly Detection
  - Multiple SKU processing
  - Race condition handling
  - Duplicate prevention (same anomaly twice)
```

**File:** `tests/unit/models/anomaly_model_test.go`

```
MUST TEST:
- ✅ Anomaly Record Creation
- ✅ Severity Enum Validation
- ✅ Status Enum Validation (open, acknowledged, resolved)
- ✅ Timestamp Validation (created_at, updated_at, acknowledged_at, resolved_at)
- ✅ Persistence (database storage)
- ✅ Retrieval by ID, SKU, warehouse, date range
```

---

#### 5️⃣ Webhook Tests (Roadmap 2.5 - Event Publishing & Retry)

**File:** `tests/unit/services/webhook_service_test.go`

```
MUST TEST:
- ✅ Webhook Subscription
  - Create subscription (valid URL, valid event types)
  - List subscriptions
  - Delete subscription
  
- ✅ Event Publishing
  - Publish event to HTTP endpoint
  - HMAC-SHA256 signing
  - Request header format (X-Webhook-Signature)
  - JSON payload encoding
  
- ✅ Retry Logic
  - Attempt 1 (immediate)
  - Attempt 2 (1s delay)
  - Attempt 3 (2s delay)
  - Attempt 4 (4s delay)
  - Attempt 5 (8s delay)
  - Dead-letter queue after 5 failures
  
- ✅ Failure Handling
  - HTTP 500 error → retry
  - HTTP 404 error → dead-letter (don't retry)
  - Timeout (5s) → retry
  - Connection refused → retry
  
- ✅ Event Deduplication
  - event_id tracking (at-least-once delivery guarantee)
  - Idempotent processing
  
- ✅ Webhook Delivery Tracking
  - Store delivery attempts
  - Track HTTP status code
  - Record error messages
  - Store request/response details
  
- ✅ HMAC Signing
  - Secret key handling
  - Signature format (sha256=hexdigest)
  - Verification on subscriber side
```

**File:** `tests/unit/services/webhook_retry_test.go`

```
MUST TEST:
- ✅ Exponential Backoff Algorithm
- ✅ Max Retry Attempts (5 total)
- ✅ Dead-Letter Queue Insertion
- ✅ Failure Rate Tracking (monitoring)
- ✅ Retry Job Scheduling (distributed task queue)
```

---

#### 6️⃣ Rate Limiting Tests (Roadmap 2.6 - Traffic Control)

**File:** `tests/unit/middleware/rate_limit_middleware_test.go`

```
MUST TEST:
- ✅ Request Counter Increment
- ✅ Sliding Window Algorithm (1-second window)
- ✅ Limit Enforcement (1000 req/sec per user/key)
- ✅ 429 Too Many Requests Response
- ✅ Rate Limit Headers (X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset)
- ✅ Redis Counter Management
- ✅ Key Extraction (user_id or API key)
- ✅ Limit Reset Timing
- ✅ Edge Cases (concurrent requests, cleanup)
```

---

#### 7️⃣ Repository/Database Tests (Roadmap 1.3 - Data Access)

**File:** `tests/unit/repositories/user_repository_test.go`

```
MUST TEST:
- ✅ Create User (unique email validation)
- ✅ Get User by ID
- ✅ Get User by Email
- ✅ Update User (email, full_name, status)
- ✅ Delete User (soft delete, preserve audit trail)
- ✅ List Users (pagination, filtering by role)
- ✅ Update Password Hash
- ✅ Update Last Login
- ✅ Update Account Status (active/suspended/deleted)
```

**File:** `tests/unit/repositories/transaction_repository_test.go`

```
MUST TEST:
- ✅ Create Transaction
- ✅ Get Transaction by ID
- ✅ List Transactions (filter by warehouse, SKU, date, type)
- ✅ Query Transactions (pagination, sorting)
- ✅ Update Transaction (status, notes)
- ✅ Delete Transaction (soft delete)
```

**File:** `tests/unit/repositories/forecast_repository_test.go`

```
MUST TEST:
- ✅ Create Forecast
- ✅ Get Forecast by SKU
- ✅ Get Forecast by Date Range
- ✅ List Forecasts (filter by warehouse, algorithm)
- ✅ Update Forecast (accuracy, status)
- ✅ Delete Old Forecasts (cleanup, retention policy)
```

**File:** `tests/unit/repositories/anomaly_repository_test.go`

```
MUST TEST:
- ✅ Create Anomaly
- ✅ Get Anomaly by ID
- ✅ List Anomalies (filter by status, severity, warehouse, SKU)
- ✅ Update Anomaly Status
- ✅ Get Anomaly History (timeline)
- ✅ Delete/Archive Anomalies
```

**File:** `tests/unit/repositories/webhook_repository_test.go`

```
MUST TEST:
- ✅ Create Subscription
- ✅ Get Subscription by ID
- ✅ List Subscriptions (filter by user, event_type)
- ✅ Delete Subscription
- ✅ Store Delivery Attempt
- ✅ Get Delivery History
```

**File:** `tests/unit/repositories/api_key_repository_test.go`

```
MUST TEST:
- ✅ Create API Key
- ✅ Get Key by Prefix
- ✅ Validate Key Hash
- ✅ List Keys (filter by user, status)
- ✅ Rotate Key (generate new hash)
- ✅ Revoke Key (soft delete)
- ✅ Check Expiration
```

---

#### 8️⃣ Utility Function Tests (Roadmap 1.5 - Helpers)

**File:** `tests/unit/util/validators_test.go`

```
MUST TEST:
- ✅ Email Validation (format, domain)
- ✅ Password Strength (min 8 chars, uppercase, digit, special char)
- ✅ SKU Code Format (alphanumeric, no special chars)
- ✅ Warehouse Name Validation
- ✅ URL Validation (for webhooks)
- ✅ Date Format Validation (RFC3339)
- ✅ Timezone Validation
```

**File:** `tests/unit/util/crypto_test.go`

```
MUST TEST:
- ✅ Password Hashing (bcrypt, cost 12)
- ✅ Password Verification (bcrypt comparison)
- ✅ JWT Generation (HS256, correct payload)
- ✅ JWT Verification (signature validation, expiry)
- ✅ HMAC-SHA256 Signing
- ✅ Hash Consistency (same input = same output)
- ✅ Salt Generation (random, unique)
```

**File:** `tests/unit/util/response_formatter_test.go`

```
MUST TEST:
- ✅ Success Response Format (JSON structure)
- ✅ Error Response Format (error message, code)
- ✅ Pagination Response (total, limit, offset)
- ✅ Timestamp Formatting (RFC3339)
- ✅ Error Code Mapping (400, 401, 403, 404, 429, 500)
- ✅ Custom Field Inclusion/Exclusion
```

**File:** `tests/unit/util/timezone_test.go`

```
MUST TEST:
- ✅ UTC Conversion
- ✅ Timezone Detection
- ✅ Timestamp Parsing
- ✅ Timestamp Formatting
- ✅ Daylight Saving Time Handling
```

---

#### 9️⃣ Configuration Tests (Roadmap 1.2)

**File:** `tests/unit/config/config_test.go`

```
MUST TEST:
- ✅ Environment Variable Loading
- ✅ Default Value Assignment
- ✅ Database Connection String
- ✅ Redis Connection String
- ✅ JWT Secret Validation
- ✅ Server Port Configuration
- ✅ Log Level Configuration
- ✅ Timeout Settings (server, database, API)
- ✅ Configuration Validation (all required values present)
```

---

#### 🔟 Database Connection Tests (Roadmap 1.3)

**File:** `tests/unit/database/connection_test.go`

```
MUST TEST:
- ✅ Database Connection Pooling (min 10, max 50 connections)
- ✅ Connection Reuse
- ✅ Idle Connection Timeout (30 seconds)
- ✅ Connection Timeout (30 seconds)
- ✅ Connection String Parsing
- ✅ SSL Mode Handling
- ✅ Connection Error Recovery
- ✅ Ping/Health Check
```

**File:** `tests/unit/database/transaction_test.go`

```
MUST TEST:
- ✅ Begin Transaction
- ✅ Commit Transaction
- ✅ Rollback on Error
- ✅ Nested Transaction Handling (savepoints)
- ✅ Transaction Isolation Levels (READ COMMITTED)
- ✅ Deadlock Handling
- ✅ Transaction Timeout
- ✅ ACID Compliance
```

---

### Unit Test Checklist

```
BACKEND UNIT TESTS CHECKLIST:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Auth Service & Middleware:
  ☐ Registration (9 test cases)
  ☐ Login (6 test cases)
  ☐ Token Management (7 test cases)
  ☐ JWT Validation (5 test cases)
  ☐ RBAC Enforcement (6 test cases)
  ☐ Account Locking (3 test cases)
  Subtotal: 36 test cases

Inventory Management:
  ☐ Inventory Service (12 test cases)
  ☐ Transaction Service (13 test cases)
  ☐ Inventory Repository (11 test cases)
  ☐ Transaction Repository (6 test cases)
  Subtotal: 42 test cases

Forecasting:
  ☐ ETS Algorithm (8 test cases)
  ☐ ARIMA Algorithm (7 test cases)
  ☐ Algorithm Comparison (5 test cases)
  ☐ Batch Processing (4 test cases)
  ☐ Accuracy Metrics (7 test cases)
  ☐ Cache Management (5 test cases)
  ☐ Forecast Repository (6 test cases)
  Subtotal: 42 test cases

Anomaly Detection:
  ☐ Z-Score Detection (7 test cases)
  ☐ IQR Detection (7 test cases)
  ☐ Severity Classification (5 test cases)
  ☐ Status Tracking (5 test cases)
  ☐ Anomaly Model (6 test cases)
  ☐ Anomaly Repository (6 test cases)
  Subtotal: 36 test cases

Webhooks:
  ☐ Webhook Subscription (6 test cases)
  ☐ Event Publishing (8 test cases)
  ☐ Retry Logic (7 test cases)
  ☐ HMAC Signing (4 test cases)
  ☐ Webhook Repository (6 test cases)
  Subtotal: 31 test cases

Rate Limiting:
  ☐ Rate Limit Middleware (9 test cases)
  Subtotal: 9 test cases

Database & Config:
  ☐ Database Connection (8 test cases)
  ☐ Database Transactions (8 test cases)
  ☐ Configuration (9 test cases)
  ☐ Validators (7 test cases)
  ☐ Crypto Utilities (7 test cases)
  ☐ Response Formatter (6 test cases)
  ☐ Timezone Utilities (5 test cases)
  Subtotal: 50 test cases

TOTAL UNIT TEST CASES: 246 test cases
TARGET COVERAGE: >90%
```

---

## 🟠 UNIT TESTS (FRONTEND - REACT/TYPESCRIPT)

### Test Coverage Targets
- **Components:** >50% coverage (40+ components)
- **Hooks:** >80% coverage (custom hooks)
- **Context:** >90% coverage (AuthContext, ThemeContext)
- **API Client:** >85% coverage (axios interceptors, error handling)
- **Utilities:** >80% coverage (formatters, validators)

### Test Files to Create

#### 1️⃣ Component Tests (40+ Components)

**Authentication Components:**
- `LoginForm.test.tsx` (form validation, submission, error display, 6 test cases)
- `RegisterForm.test.tsx` (registration flow, password validation, 7 test cases)
- `PasswordResetForm.test.tsx` (reset flow, token validation, 5 test cases)

**Dashboard Components:**
- `StatsCard.test.tsx` (data display, formatting, 4 test cases)
- `RecentTransactions.test.tsx` (table rendering, sorting, pagination, 8 test cases)
- `ForecastChart.test.tsx` (Recharts integration, data binding, 6 test cases)
- `SystemStatus.test.tsx` (health check display, status colors, 4 test cases)
- `KPIGrid.test.tsx` (grid layout, responsive, 5 test cases)

**Inventory Components:**
- `InventoryTable.test.tsx` (sorting, filtering, pagination, 10 test cases)
- `InventoryFilters.test.tsx` (warehouse filter, status filter, 6 test cases)
- `StatusBadge.test.tsx` (status color mapping, text, 4 test cases)
- `TransactionForm.test.tsx` (form validation, submission, error display, 8 test cases)
- `ImportButton.test.tsx` (file upload, parsing, 5 test cases)

**Forecast Components:**
- `ForecastChart.test.tsx` (chart rendering, data points, confidence intervals, 8 test cases)
- `ForecastTable.test.tsx` (forecast display, sorting, filtering, 7 test cases)
- `AccuracyMetrics.test.tsx` (MAPE/RMSE display, formatting, 5 test cases)
- `AlgorithmComparison.test.tsx` (algorithm comparison chart, selection, 6 test cases)

**Anomaly Components:**
- `AnomalyList.test.tsx` (list rendering, filtering by status/severity, 8 test cases)
- `AnomalyCard.test.tsx` (anomaly details display, actions, 6 test cases)
- `SeverityBadge.test.tsx` (severity color mapping, icons, 4 test cases)
- `AnomalyTimeline.test.tsx` (timeline rendering, status progression, 6 test cases)
- `ResolutionModal.test.tsx` (modal display, form submission, 5 test cases)
- `AnomalyChart.test.tsx` (trend visualization, data points, 6 test cases)

**Common Components:**
- `Button.test.tsx` (click handler, disabled state, variants, 5 test cases)
- `Card.test.tsx` (layout, shadow, responsive, 3 test cases)
- `Input.test.tsx` (text input, validation, error display, 6 test cases)
- `Select.test.tsx` (dropdown, option selection, filtering, 7 test cases)
- `Modal.test.tsx` (open/close, overlay, keyboard escape, 5 test cases)
- `Toast.test.tsx` (toast notification display, auto-dismiss, 5 test cases)
- `Loading.test.tsx` (spinner display, loading states, 4 test cases)
- `ErrorBoundary.test.tsx` (error catching, fallback display, 4 test cases)
- `Skeleton.test.tsx` (skeleton animation, structure, 3 test cases)
- `Pagination.test.tsx` (page navigation, active state, disabled, 6 test cases)
- `ConfirmDialog.test.tsx` (confirmation modal, cancel/confirm actions, 5 test cases)
- `Breadcrumb.test.tsx` (navigation trail, links, 4 test cases)
- `Badge.test.tsx` (badge display, variants, 3 test cases)

**Layout Components:**
- `Navbar.test.tsx` (navigation links, user menu, logo, 5 test cases)
- `Sidebar.test.tsx` (menu items, active link highlight, collapse, 6 test cases)
- `Footer.test.tsx` (footer links, copyright, 2 test cases)

**Subtotal: 200+ test cases for components**

---

#### 2️⃣ Hook Tests

**File:** `tests/unit/hooks/useAuth.test.ts`

```
MUST TEST:
- ✅ Login Hook (API call, token storage, user context)
- ✅ Logout Hook (token cleanup, redirect)
- ✅ Register Hook (registration flow, auto-login)
- ✅ Token Refresh Hook (refresh logic, retry)
- ✅ Error Handling (401, 403, 500 errors)
- ✅ Loading State (pending, resolved, rejected)
```

**File:** `tests/unit/hooks/useTheme.test.ts`

```
MUST TEST:
- ✅ Theme Toggle (light ↔ dark)
- ✅ Theme Persistence (localStorage)
- ✅ System Preference Detection (prefers-color-scheme)
- ✅ Theme Application (CSS class application)
- ✅ Provider Context (theme available to components)
```

**File:** `tests/unit/hooks/useInventory.test.ts`

```
MUST TEST:
- ✅ Fetch Inventory (API call, caching)
- ✅ Filter Inventory (by warehouse, status)
- ✅ Sort Inventory (by SKU, quantity)
- ✅ Pagination (page navigation)
- ✅ Real-time Updates (polling, WebSocket)
- ✅ Error Handling (network errors, 500)
- ✅ Loading State
- ✅ Cache Invalidation (on update)
```

**File:** `tests/unit/hooks/useForecast.test.ts`

```
MUST TEST:
- ✅ Fetch Forecasts (API call)
- ✅ Filter Forecasts (by SKU, warehouse)
- ✅ Algorithm Comparison (ETS vs ARIMA)
- ✅ Accuracy Metrics (MAPE, RMSE)
- ✅ Batch Generation (multiple SKUs)
- ✅ Error Handling
- ✅ Cache Management (24h TTL)
```

**File:** `tests/unit/hooks/useAnomaly.test.ts`

```
MUST TEST:
- ✅ Fetch Anomalies (API call)
- ✅ Filter by Status (open, acknowledged, resolved)
- ✅ Filter by Severity (critical, high, medium, low)
- ✅ Acknowledge Anomaly (status update)
- ✅ Resolve Anomaly (status update)
- ✅ Get History (timeline)
- ✅ Real-time Updates (polling)
```

**Subtotal: 50+ test cases for hooks**

---

#### 3️⃣ Context Tests

**File:** `tests/unit/context/AuthContext.test.tsx`

```
MUST TEST:
- ✅ Provide Auth Context
- ✅ Login Action (update user, token)
- ✅ Logout Action (clear user, token)
- ✅ Token Refresh Action
- ✅ User Role Access
- ✅ Authenticated State
```

**File:** `tests/unit/context/ThemeContext.test.tsx`

```
MUST TEST:
- ✅ Provide Theme Context
- ✅ Toggle Theme Action
- ✅ Read Theme from localStorage
- ✅ Write Theme to localStorage
- ✅ Theme Value Propagation to Children
```

**Subtotal: 15+ test cases for contexts**

---

#### 4️⃣ API Client Tests

**File:** `tests/unit/api/client.test.ts`

```
MUST TEST:
- ✅ Axios Instance Creation
- ✅ Request Interceptor (attach JWT token)
- ✅ Response Interceptor (handle success)
- ✅ Error Interceptor (handle errors)
- ✅ Token Refresh on 401 (automatic refresh)
- ✅ Retry Logic (exponential backoff)
- ✅ Timeout Handling (30s)
- ✅ Request Headers (Content-Type, Authorization)
- ✅ Request/Response Transformation
```

**File:** `tests/unit/api/endpoints.test.ts`

```
MUST TEST:
- ✅ Auth Endpoints (login, register, logout, refresh)
- ✅ Inventory Endpoints (GET all, GET warehouse, POST transaction)
- ✅ Forecast Endpoints (GET list, POST generate, GET accuracy)
- ✅ Anomaly Endpoints (GET list, PUT status, GET history)
- ✅ Webhook Endpoints (POST subscribe, DELETE subscription)
- ✅ API Key Endpoints (POST create, GET list, DELETE revoke)
```

**Subtotal: 25+ test cases for API client**

---

#### 5️⃣ Utility Function Tests

**File:** `tests/unit/utils/formatters.test.ts`

```
MUST TEST:
- ✅ Format Currency (IDR, 2 decimal places)
- ✅ Format Date (DD/MM/YYYY, DD Mon YYYY)
- ✅ Format DateTime (RFC3339)
- ✅ Format Number (with thousand separator)
- ✅ Format Percentage (2 decimal places, % symbol)
- ✅ Format Timestamp (relative: "2 hours ago")
```

**File:** `tests/unit/utils/validators.test.ts`

```
MUST TEST:
- ✅ Email Validation
- ✅ Password Strength (min 8, uppercase, digit, special)
- ✅ URL Validation (for webhooks)
- ✅ Quantity Validation (positive integer)
- ✅ SKU Format Validation
```

**File:** `tests/unit/utils/localStorage.test.ts`

```
MUST TEST:
- ✅ Store Theme Setting
- ✅ Retrieve Theme Setting
- ✅ Store JWT Token
- ✅ Retrieve JWT Token
- ✅ Clear All Data (logout)
- ✅ Error Handling (quota exceeded)
```

**Subtotal: 20+ test cases for utilities**

---

### Frontend Unit Test Checklist

```
FRONTEND UNIT TESTS CHECKLIST:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Components (40+):
  ☐ Auth Components (3 components × 6 tests) = 18
  ☐ Dashboard Components (5 components × 5.5 tests) = 27
  ☐ Inventory Components (5 components × 7.2 tests) = 36
  ☐ Forecast Components (4 components × 6.5 tests) = 26
  ☐ Anomaly Components (6 components × 5.8 tests) = 35
  ☐ Common Components (13 components × 4.3 tests) = 56
  ☐ Layout Components (3 components × 4.3 tests) = 13
  Subtotal: 211 test cases

Hooks (6 custom hooks):
  ☐ useAuth Hook (6 test cases)
  ☐ useTheme Hook (5 test cases)
  ☐ useInventory Hook (8 test cases)
  ☐ useForecast Hook (6 test cases)
  ☐ useAnomaly Hook (7 test cases)
  ☐ Custom hooks (5 more) × 5 tests = 25
  Subtotal: 57 test cases

Contexts (2 contexts):
  ☐ AuthContext (6 test cases)
  ☐ ThemeContext (5 test cases)
  Subtotal: 11 test cases

API Client:
  ☐ Axios Configuration (9 test cases)
  ☐ API Endpoints (6 test cases)
  Subtotal: 15 test cases

Utilities:
  ☐ Formatters (6 test cases)
  ☐ Validators (5 test cases)
  ☐ localStorage (6 test cases)
  Subtotal: 17 test cases

TOTAL FRONTEND UNIT TEST CASES: 311 test cases
TARGET COVERAGE: >50%
```

---

## 🟡 INTEGRATION TESTS

### Backend Integration Tests (5+ Critical Flows)

#### Flow 1️⃣: Authentication Complete Cycle

**File:** `tests/integration/auth_flow_test.go`

**Test Cases:**
1. ✅ Register new user (POST /auth/register)
2. ✅ Verify user created in database
3. ✅ Login with new credentials (POST /auth/login)
4. ✅ Verify JWT token received
5. ✅ Use JWT to access protected endpoint (GET /inventory)
6. ✅ Verify data returned matches user's role
7. ✅ Refresh JWT token (POST /auth/refresh)
8. ✅ Verify new token valid
9. ✅ Logout (POST /auth/logout)
10. ✅ Verify token invalidated in Redis
11. ✅ Attempt access with invalidated token → 401 error

**Integration Points:**
- User table
- Auth service
- JWT generation
- Database transaction
- Redis token storage

---

#### Flow 2️⃣: Transaction & Inventory Update Cycle

**File:** `tests/integration/transaction_flow_test.go`

**Test Cases:**
1. ✅ Authenticate as manager
2. ✅ Check current inventory (GET /inventory)
3. ✅ Create outbound transaction (POST /transactions)
4. ✅ Verify inventory quantity decremented
5. ✅ Verify transaction stored in database
6. ✅ Create inbound transaction (POST /transactions)
7. ✅ Verify inventory quantity incremented
8. ✅ Query transaction history (GET /transactions)
9. ✅ Verify all transactions visible
10. ✅ Verify cache invalidated (fetch returns latest data)
11. ✅ Attempt transaction exceeding stock → error
12. ✅ Verify inventory unchanged on error

**Integration Points:**
- Inventory service
- Transaction service
- Database transactions (deadlock handling)
- Redis cache invalidation
- Warehouse isolation (RBAC)

---

#### Flow 3️⃣: Forecast Generation & Accuracy Cycle

**File:** `tests/integration/forecast_flow_test.go`

**Test Cases:**
1. ✅ Seed 180 days of historical data
2. ✅ Authenticate as admin
3. ✅ Trigger forecast generation (POST /forecasts/generate)
4. ✅ Verify batch job created
5. ✅ Wait for completion (polling)
6. ✅ Verify forecasts stored in database
7. ✅ Verify forecasts cached in Redis (24h TTL)
8. ✅ Query forecasts (GET /forecasts)
9. ✅ Verify ETS algorithm results
10. ✅ Verify ARIMA algorithm results
11. ✅ Get accuracy metrics (GET /forecasts/accuracy)
12. ✅ Verify MAPE, RMSE, MAE calculated
13. ✅ Verify algorithm comparison (which is better)
14. ✅ Update historical data
15. ✅ Trigger recalibration
16. ✅ Verify new forecasts generated

**Integration Points:**
- Historical inventory data
- ETS algorithm service
- ARIMA algorithm service
- Batch job scheduler
- Database storage
- Redis caching
- Accuracy calculation

---

#### Flow 4️⃣: Anomaly Detection → Acknowledgment → Resolution

**File:** `tests/integration/anomaly_flow_test.go`

**Test Cases:**
1. ✅ Create unusual inventory transaction (Z-score > 3)
2. ✅ Verify anomaly detected (Z-score method)
3. ✅ Verify anomaly severity calculated
4. ✅ Verify anomaly status = "open"
5. ✅ Get anomalies list (GET /anomalies) → anomaly visible
6. ✅ Acknowledge anomaly (PUT /anomalies/{id})
7. ✅ Verify status → "acknowledged"
8. ✅ Verify acknowledged_at timestamp set
9. ✅ Get anomaly detail (GET /anomalies/{id})
10. ✅ Verify status updated
11. ✅ Get anomaly history (GET /anomalies/{id}/history)
12. ✅ Verify status transition recorded
13. ✅ Resolve anomaly (PUT /anomalies/{id}, status=resolved)
14. ✅ Verify resolved_at timestamp set
15. ✅ Verify timeline shows all transitions

**Integration Points:**
- Transaction service
- Anomaly detection (Z-score/IQR)
- Anomaly service
- Database storage
- Status tracking
- Timeline generation

---

#### Flow 5️⃣: Webhook Event Publishing with Retry

**File:** `tests/integration/webhook_flow_test.go`

**Test Cases:**
1. ✅ Create webhook subscription (POST /webhooks/subscribe)
  - Target URL: mock HTTP server
  - Event types: transaction.created, anomaly.detected
2. ✅ Create transaction → event triggered
3. ✅ Verify webhook event published
4. ✅ Verify mock server received event
5. ✅ Verify HMAC-SHA256 signature valid
6. ✅ Verify event_id deduplication works
7. ✅ Mock subscriber server → respond with 500 error
8. ✅ Verify retry attempt 1 (immediate)
9. ✅ Verify retry attempt 2 (1s delay)
10. ✅ Verify retry attempt 3 (2s delay)
11. ✅ Verify retry attempt 4 (4s delay)
12. ✅ Verify retry attempt 5 (8s delay)
13. ✅ Mock subscriber → respond with 200 OK
14. ✅ Verify delivery marked successful
15. ✅ Query delivery history (GET /webhooks/deliveries)
16. ✅ Verify all attempts recorded
17. ✅ Test permanent failure (404) → no retry
18. ✅ Verify dead-letter queue storage

**Integration Points:**
- Webhook subscription
- Event publishing
- Retry logic
- HTTP client
- HMAC signing
- Dead-letter queue
- Delivery tracking

---

### Frontend Integration Tests (5+ E2E Flows)

#### E2E Flow 1️⃣: Complete User Registration & Login

```typescript
Test Case: Registration → Automatic Login → Dashboard
1. Navigate to /register
2. Enter email: "testuser@example.com"
3. Enter password: "SecurePass123!"
4. Click Register button
5. Verify API call made (POST /auth/register)
6. Verify redirect to /dashboard
7. Verify navbar shows user name
8. Verify dashboard loads (no 401 error)
9. Verify inventory data displayed
```

---

#### E2E Flow 2️⃣: Create Transaction & View in History

```typescript
Test Case: Create Inventory Transaction → Verify in List
1. Login as manager
2. Navigate to /inventory
3. Click "New Transaction" button
4. Select warehouse (from dropdown)
5. Select SKU (from dropdown)
6. Select "Out" type
7. Enter quantity: 10
8. Enter reason: "Customer order"
9. Click Submit
10. Verify success toast
11. Verify inventory quantity decremented
12. Navigate to /inventory/transactions
13. Verify new transaction in list
14. Filter by warehouse → verify filtered
15. Verify real-time update (30s polling)
```

---

#### E2E Flow 3️⃣: Generate Forecast & View Charts

```typescript
Test Case: Trigger Forecast → View Results → Compare Algorithms
1. Login as admin
2. Navigate to /forecast
3. Click "Generate Forecast" button
4. Verify modal appears
5. Select SKUs (all, or subset)
6. Click Generate
7. Verify loading spinner
8. Wait for completion (backend 5s batch job)
9. Verify chart renders (Recharts)
10. Verify forecast points displayed
11. Verify confidence intervals shown
12. Navigate to /forecast/accuracy
13. Verify accuracy metrics displayed (MAPE, RMSE, MAE)
14. Navigate to /forecast/comparison
15. Verify algorithm comparison chart
16. Verify ETS vs ARIMA scores
```

---

#### E2E Flow 4️⃣: Anomaly Acknowledge & Resolve

```typescript
Test Case: View Anomalies → Acknowledge → Resolve → Verify Timeline
1. Login as analyst
2. Navigate to /anomalies
3. Verify list of anomalies displayed
4. Click anomaly card (first anomaly)
5. Verify detail page loads
6. Click "Acknowledge" button
7. Verify status badge changes to "Acknowledged"
8. Verify timestamp updated
9. Click "Resolve" button
10. Verify modal for resolution notes
11. Enter notes: "Fixed by restocking"
12. Click Submit
13. Verify status changes to "Resolved"
14. Navigate to /anomalies/:id/timeline
15. Verify timeline shows:
    - Created: [timestamp]
    - Acknowledged: [timestamp]
    - Resolved: [timestamp]
```

---

#### E2E Flow 5️⃣: Dark Mode Toggle & Persistence

```typescript
Test Case: Toggle Theme → Refresh → Verify Persistence
1. Navigate to /dashboard
2. Verify light mode active (background white)
3. Click theme toggle (moon icon in navbar)
4. Verify page switches to dark mode (background dark)
5. Verify all components styled for dark mode
6. Refresh page (F5)
7. Verify dark mode still active (localStorage read)
8. Click theme toggle again
9. Verify switch to light mode
10. Refresh page
11. Verify light mode still active
12. Verify contrast ratios meet WCAG AA (4.5:1)
```

---

### Integration Test Checklist

```
INTEGRATION TESTS CHECKLIST:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Backend Integration (5 critical flows):
  ☐ Flow 1: Auth Complete Cycle (11 assertions)
  ☐ Flow 2: Transaction & Inventory Cycle (12 assertions)
  ☐ Flow 3: Forecast Generation Cycle (16 assertions)
  ☐ Flow 4: Anomaly Detection Cycle (15 assertions)
  ☐ Flow 5: Webhook Event Publishing (18 assertions)
  Subtotal: 72 integration test cases

Frontend E2E (5+ critical flows):
  ☐ Flow 1: Registration & Login (11 steps)
  ☐ Flow 2: Transaction Creation (15 steps)
  ☐ Flow 3: Forecast Generation (16 steps)
  ☐ Flow 4: Anomaly Acknowledge & Resolve (15 steps)
  ☐ Flow 5: Dark Mode Toggle & Persistence (11 steps)
  Subtotal: 5+ E2E flows (68+ user actions verified)

TOTAL INTEGRATION ASSERTIONS: 72 + 68 = 140+ verifications
```

---

## 🔴 END-TO-END (E2E) TESTS

### E2E Test Strategy (Playwright/Cypress)

### E2E Test Files to Create

#### 1️⃣ Authentication E2E Tests

**File:** `tests/e2e/auth.spec.ts`

```typescript
describe("Authentication E2E Tests", () => {
  test("Complete registration flow", async ({ page }) => {
    // Given: User on register page
    await page.goto("/register");
    
    // When: User fills form
    await page.fill('input[name="email"]', "newuser@example.com");
    await page.fill('input[name="password"]', "SecurePass123!");
    await page.fill('input[name="confirmPassword"]', "SecurePass123!");
    await page.click('button[type="submit"]');
    
    // Then: Redirected to dashboard
    await page.waitForURL("/dashboard");
    expect(page.url()).toContain("/dashboard");
  });
  
  test("Login with correct credentials", async ({ page }) => {
    await page.goto("/login");
    await page.fill('input[name="email"]', "testuser@example.com");
    await page.fill('input[name="password"]', "SecurePass123!");
    await page.click('button[type="submit"]');
    
    await page.waitForURL("/dashboard");
    expect(page.locator("text=Dashboard").isVisible()).toBeTruthy();
  });
  
  test("Login with incorrect credentials", async ({ page }) => {
    await page.goto("/login");
    await page.fill('input[name="email"]', "testuser@example.com");
    await page.fill('input[name="password"]', "WrongPassword");
    await page.click('button[type="submit"]');
    
    expect(page.locator("text=Invalid credentials").isVisible()).toBeTruthy();
    expect(page.url()).toContain("/login");
  });
  
  test("Forgot password flow", async ({ page }) => {
    await page.goto("/login");
    await page.click("text=Forgot Password?");
    
    await page.fill('input[name="email"]', "testuser@example.com");
    await page.click('button[type="submit"]');
    
    expect(page.locator("text=Reset link sent").isVisible()).toBeTruthy();
  });
  
  test("Logout clears session", async ({ page }) => {
    // Login first
    await page.goto("/login");
    await page.fill('input[name="email"]', "testuser@example.com");
    await page.fill('input[name="password"]', "SecurePass123!");
    await page.click('button[type="submit"]');
    await page.waitForURL("/dashboard");
    
    // Logout
    await page.click("button[aria-label='User menu']");
    await page.click("text=Logout");
    
    // Verify redirected to login
    await page.waitForURL("/login");
    expect(page.url()).toContain("/login");
  });
});
```

---

#### 2️⃣ Inventory Management E2E Tests

**File:** `tests/e2e/inventory.spec.ts`

```typescript
describe("Inventory Management E2E Tests", () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto("/login");
    await page.fill('input[name="email"]', "manager@example.com");
    await page.fill('input[name="password"]', "SecurePass123!");
    await page.click('button[type="submit"]');
    await page.waitForURL("/dashboard");
  });
  
  test("View inventory list", async ({ page }) => {
    await page.goto("/inventory");
    
    // Verify table renders
    expect(page.locator("table").isVisible()).toBeTruthy();
    expect(page.locator("thead").isVisible()).toBeTruthy();
    
    // Verify columns
    expect(page.locator("text=SKU Code").isVisible()).toBeTruthy();
    expect(page.locator("text=Quantity").isVisible()).toBeTruthy();
    expect(page.locator("text=Status").isVisible()).toBeTruthy();
  });
  
  test("Filter inventory by warehouse", async ({ page }) => {
    await page.goto("/inventory");
    
    // Open filter
    await page.click('button:has-text("Filter")');
    
    // Select warehouse
    await page.selectOption('select[name="warehouse"]', { value: "wh-001" });
    
    // Verify table updated
    await page.waitForTimeout(500);
    expect(page.locator("tbody tr").count()).toBeGreaterThan(0);
  });
  
  test("Create inbound transaction", async ({ page }) => {
    await page.goto("/inventory/transaction/new");
    
    // Fill form
    await page.selectOption('select[name="warehouse"]', { value: "wh-001" });
    await page.selectOption('select[name="sku"]', { value: "sku-001" });
    await page.click('label:has-text("In")');
    await page.fill('input[name="quantity"]', "100");
    await page.fill('textarea[name="reason"]', "Restock from supplier");
    
    // Submit
    await page.click('button[type="submit"]');
    
    // Verify success
    expect(page.locator("text=Transaction created successfully").isVisible()).toBeTruthy();
  });
  
  test("Create outbound transaction", async ({ page }) => {
    await page.goto("/inventory/transaction/new");
    
    // Fill form
    await page.selectOption('select[name="warehouse"]', { value: "wh-001" });
    await page.selectOption('select[name="sku"]', { value: "sku-001" });
    await page.click('label:has-text("Out")');
    await page.fill('input[name="quantity"]', "50");
    await page.fill('textarea[name="reason"]', "Customer order #12345");
    
    // Submit
    await page.click('button[type="submit"]');
    
    // Verify success
    expect(page.locator("text=Transaction created successfully").isVisible()).toBeTruthy();
  });
  
  test("View transaction history", async ({ page }) => {
    await page.goto("/inventory/transactions");
    
    // Verify table
    expect(page.locator("table").isVisible()).toBeTruthy();
    expect(page.locator("text=Date").isVisible()).toBeTruthy();
    expect(page.locator("text=Type").isVisible()).toBeTruthy();
    
    // Verify transactions displayed
    expect(page.locator("tbody tr").count()).toBeGreaterThan(0);
  });
});
```

---

#### 3️⃣ Forecasting E2E Tests

**File:** `tests/e2e/forecast.spec.ts`

```typescript
describe("Forecasting E2E Tests", () => {
  test.beforeEach(async ({ page }) => {
    // Login as admin
    await page.goto("/login");
    await page.fill('input[name="email"]', "admin@example.com");
    await page.fill('input[name="password"]', "SecurePass123!");
    await page.click('button[type="submit"]');
    await page.waitForURL("/dashboard");
  });
  
  test("View forecast list", async ({ page }) => {
    await page.goto("/forecast");
    
    // Verify table renders
    expect(page.locator("table").isVisible()).toBeTruthy();
    expect(page.locator("text=SKU").isVisible()).toBeTruthy();
    expect(page.locator("text=Forecast").isVisible()).toBeTruthy();
  });
  
  test("Generate forecast batch job", async ({ page }) => {
    await page.goto("/forecast");
    
    // Click generate button
    await page.click('button:has-text("Generate Forecast")');
    
    // Verify modal
    expect(page.locator('text="Generate Forecast"').isVisible()).toBeTruthy();
    
    // Select all SKUs
    await page.click('input[type="checkbox"][value="all"]');
    
    // Submit
    await page.click('button:has-text("Start Generation")');
    
    // Wait for generation
    await page.waitForTimeout(6000); // Backend 5s + overhead
    
    // Verify success message
    expect(page.locator("text=Forecasts generated successfully").isVisible()).toBeTruthy();
  });
  
  test("View forecast chart", async ({ page }) => {
    await page.goto("/forecast");
    
    // Click on forecast detail
    await page.click('tbody tr:first-child');
    
    // Verify chart renders
    expect(page.locator("svg").first()).toBeTruthy();
    
    // Verify confidence intervals visible
    expect(page.locator("text=Confidence Interval").isVisible()).toBeTruthy();
  });
  
  test("View accuracy metrics", async ({ page }) => {
    await page.goto("/forecast/accuracy");
    
    // Verify metrics displayed
    expect(page.locator("text=MAPE").isVisible()).toBeTruthy();
    expect(page.locator("text=RMSE").isVisible()).toBeTruthy();
    expect(page.locator("text=MAE").isVisible()).toBeTruthy();
  });
  
  test("Compare forecast algorithms", async ({ page }) => {
    await page.goto("/forecast/comparison");
    
    // Verify comparison chart
    expect(page.locator("svg").first()).toBeTruthy();
    
    // Verify algorithms compared
    expect(page.locator("text=ETS").isVisible()).toBeTruthy();
    expect(page.locator("text=ARIMA").isVisible()).toBeTruthy();
  });
});
```

---

#### 4️⃣ Anomaly Detection E2E Tests

**File:** `tests/e2e/anomalies.spec.ts`

```typescript
describe("Anomaly Detection E2E Tests", () => {
  test.beforeEach(async ({ page }) => {
    // Login as analyst
    await page.goto("/login");
    await page.fill('input[name="email"]', "analyst@example.com");
    await page.fill('input[name="password"]', "SecurePass123!");
    await page.click('button[type="submit"]');
    await page.waitForURL("/dashboard");
  });
  
  test("View anomaly list", async ({ page }) => {
    await page.goto("/anomalies");
    
    // Verify table
    expect(page.locator("table").isVisible()).toBeTruthy();
    expect(page.locator("text=SKU").isVisible()).toBeTruthy();
    expect(page.locator("text=Status").isVisible()).toBeTruthy();
  });
  
  test("Filter anomalies by status", async ({ page }) => {
    await page.goto("/anomalies");
    
    // Filter by status
    await page.click('label:has-text("Open")');
    
    // Verify filtered results
    await page.waitForTimeout(500);
    expect(page.locator("text=Open").first()).toBeTruthy();
  });
  
  test("Filter anomalies by severity", async ({ page }) => {
    await page.goto("/anomalies");
    
    // Filter by severity
    await page.click('label:has-text("Critical")');
    
    // Verify filtered results
    await page.waitForTimeout(500);
    expect(page.locator("text=Critical").first()).toBeTruthy();
  });
  
  test("View anomaly detail and acknowledge", async ({ page }) => {
    await page.goto("/anomalies");
    
    // Click first anomaly
    await page.click("tbody tr:first-child");
    
    // Verify detail page
    expect(page.locator("text=Anomaly Detail").isVisible()).toBeTruthy();
    
    // Click acknowledge
    await page.click('button:has-text("Acknowledge")');
    
    // Verify status changed
    expect(page.locator("text=Acknowledged").isVisible()).toBeTruthy();
  });
  
  test("View anomaly timeline", async ({ page }) => {
    await page.goto("/anomalies");
    
    // Click first anomaly
    await page.click("tbody tr:first-child");
    
    // Click timeline tab
    await page.click('button:has-text("Timeline")');
    
    // Verify timeline
    expect(page.locator("text=Created").isVisible()).toBeTruthy();
    expect(page.locator("text=Acknowledged").isVisible()).toBeTruthy();
  });
  
  test("Resolve anomaly", async ({ page }) => {
    await page.goto("/anomalies");
    
    // Click first anomaly
    await page.click("tbody tr:first-child");
    
    // Click resolve
    await page.click('button:has-text("Resolve")');
    
    // Verify modal
    expect(page.locator('text="Resolve Anomaly"').isVisible()).toBeTruthy();
    
    // Enter resolution notes
    await page.fill('textarea[name="notes"]', "Fixed by restocking warehouse");
    
    // Submit
    await page.click('button:has-text("Resolve")');
    
    // Verify status
    expect(page.locator("text=Resolved").isVisible()).toBeTruthy();
  });
});
```

---

#### 5️⃣ User Settings & Theme E2E Tests

**File:** `tests/e2e/settings.spec.ts`

```typescript
describe("Settings & Theme E2E Tests", () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto("/login");
    await page.fill('input[name="email"]', "testuser@example.com");
    await page.fill('input[name="password"]', "SecurePass123!");
    await page.click('button[type="submit"]');
    await page.waitForURL("/dashboard");
  });
  
  test("Toggle dark mode", async ({ page }) => {
    // Verify light mode
    const html = page.locator("html");
    expect(html).not.toHaveAttribute("class", /dark/);
    
    // Click theme toggle
    await page.click('button[aria-label="Toggle theme"]');
    
    // Verify dark mode
    expect(html).toHaveAttribute("class", /dark/);
  });
  
  test("Dark mode persists on refresh", async ({ page }) => {
    // Enable dark mode
    await page.click('button[aria-label="Toggle theme"]');
    
    // Refresh
    await page.reload();
    
    // Verify dark mode still active
    const html = page.locator("html");
    expect(html).toHaveAttribute("class", /dark/);
  });
  
  test("Navigate to settings page", async ({ page }) => {
    await page.goto("/settings");
    
    // Verify settings page
    expect(page.locator("text=Settings").isVisible()).toBeTruthy();
    expect(page.locator("text=Theme").isVisible()).toBeTruthy();
  });
  
  test("View API keys in settings", async ({ page }) => {
    await page.goto("/settings/api-keys");
    
    // Verify table
    expect(page.locator("table").isVisible()).toBeTruthy();
    expect(page.locator("text=Key").isVisible()).toBeTruthy();
  });
  
  test("Create new API key", async ({ page }) => {
    await page.goto("/settings/api-keys");
    
    // Click create button
    await page.click('button:has-text("Create Key")');
    
    // Verify modal
    expect(page.locator('text="Create API Key"').isVisible()).toBeTruthy();
    
    // Enter key name
    await page.fill('input[name="name"]', "Test Key");
    
    // Submit
    await page.click('button:has-text("Create")');
    
    // Verify key created
    expect(page.locator("text=API key created successfully").isVisible()).toBeTruthy();
  });
  
  test("Delete API key", async ({ page }) => {
    await page.goto("/settings/api-keys");
    
    // Click delete button
    await page.click('button[aria-label="Delete"]');
    
    // Confirm deletion
    await page.click('button:has-text("Confirm")');
    
    // Verify deleted
    expect(page.locator("text=Key deleted successfully").isVisible()).toBeTruthy();
  });
});
```

---

### E2E Test Checklist

```
E2E TEST CHECKLIST:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Authentication:
  ☐ Registration complete flow
  ☐ Login with correct credentials
  ☐ Login with incorrect credentials
  ☐ Forgot password flow
  ☐ Logout clears session
  Subtotal: 5 test cases

Inventory Management:
  ☐ View inventory list
  ☐ Filter inventory by warehouse
  ☐ Create inbound transaction
  ☐ Create outbound transaction
  ☐ View transaction history
  Subtotal: 5 test cases

Forecasting:
  ☐ View forecast list
  ☐ Generate forecast batch job
  ☐ View forecast chart
  ☐ View accuracy metrics
  ☐ Compare forecast algorithms
  Subtotal: 5 test cases

Anomaly Detection:
  ☐ View anomaly list
  ☐ Filter anomalies by status
  ☐ Filter anomalies by severity
  ☐ View anomaly detail and acknowledge
  ☐ View anomaly timeline
  ☐ Resolve anomaly
  Subtotal: 6 test cases

Settings & Theme:
  ☐ Toggle dark mode
  ☐ Dark mode persists on refresh
  ☐ Navigate to settings page
  ☐ View API keys in settings
  ☐ Create new API key
  ☐ Delete API key
  Subtotal: 6 test cases

TOTAL E2E TEST CASES: 27 test cases
CRITICAL PATHS TESTED: 5+ complete user journeys
```

---

## 🟢 PERFORMANCE & LOAD TESTS

### Backend Load Testing (k6/Artillery)

**File:** `tests/load/inventory_load_test.js`

```javascript
import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
  stages: [
    { duration: "30s", target: 100 },  // Ramp-up
    { duration: "1m30s", target: 1000 }, // Peak load
    { duration: "30s", target: 0 },      // Ramp-down
  ],
  thresholds: {
    http_req_duration: ["p(99)<100"],  // p99 <100ms
    http_req_failed: ["rate<0.1"],     // <0.1% failure rate
  },
};

export default function () {
  // Setup: Login
  let res = http.post(`${__ENV.API_URL}/auth/login`, {
    email: `user${__VU}@example.com`,
    password: "SecurePass123!",
  });
  
  let token = res.json("token");
  
  // Get inventory (most critical endpoint)
  let inventoryRes = http.get(`${__ENV.API_URL}/inventory`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  
  check(inventoryRes, {
    "status 200": (r) => r.status === 200,
    "p99 <100ms": (r) => r.timings.duration < 100,
  });
  
  sleep(1);
}
```

**Test Metrics:**
- ✅ 1000 concurrent users
- ✅ GET /inventory p99 <100ms
- ✅ Throughput: 1000+ req/sec
- ✅ Error rate: <0.1%

---

### Frontend Performance Testing (Lighthouse)

**Target Scores:**
- ✅ Performance: >90
- ✅ Accessibility: >90
- ✅ Best Practices: >90
- ✅ SEO: >90

**Core Web Vitals:**
- ✅ LCP (Largest Contentful Paint): <2.5s
- ✅ FID (First Input Delay): <100ms
- ✅ CLS (Cumulative Layout Shift): <0.1

---

### Load Test Checklist

```
PERFORMANCE & LOAD TESTS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Backend Load Tests:
  ☐ Inventory query (1000 users, p99 <100ms)
  ☐ Forecast batch (1000 SKUs, <5s)
  ☐ Concurrent transactions (100 users, no deadlock)
  ☐ Webhook delivery (100 deliveries/sec)
  ☐ Rate limiting (exceed 1000 req/sec)
  
Frontend Performance:
  ☐ Lighthouse score >90 (all metrics)
  ☐ Core Web Vitals (LCP <2.5s, FID <100ms, CLS <0.1)
  ☐ Page load time <3s (all pages)
  ☐ Time to Interactive <5s
  ☐ Bundle size <500KB (gzipped)

Security Load Tests:
  ☐ Brute force (1000 login attempts, account locked)
  ☐ CORS validation (cross-origin requests blocked)
  ☐ Rate limiting enforcement (429 response)
```

---

## 📋 COMPREHENSIVE TEST EXECUTION CHECKLIST

### Phase 1: Backend Unit Tests

```
BACKEND UNIT TESTS EXECUTION:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Week 9.1 Task (2 hours):
  ☐ Auth service tests (36 cases) → PASS
  ☐ Inventory tests (42 cases) → PASS
  ☐ Forecast tests (42 cases) → PASS
  ☐ Anomaly tests (36 cases) → PASS
  ☐ Webhook tests (31 cases) → PASS
  ☐ Rate limiting tests (9 cases) → PASS
  ☐ Database tests (50 cases) → PASS
  
TOTAL: 246 unit test cases
TARGET: >90% coverage
COMMAND: go test -v -coverprofile=coverage.out ./...
EXPECTED RESULT: All tests pass, coverage >90%
```

---

### Phase 2: Frontend Unit Tests

```
FRONTEND UNIT TESTS EXECUTION:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Week 9.2 Task (1 hour):
  ☐ Component tests (211 cases) → PASS
  ☐ Hook tests (57 cases) → PASS
  ☐ Context tests (11 cases) → PASS
  ☐ API client tests (15 cases) → PASS
  ☐ Utility tests (17 cases) → PASS

TOTAL: 311 unit test cases
TARGET: >50% coverage
COMMAND: npm run test -- --coverage
EXPECTED RESULT: All tests pass, coverage >50%
```

---

### Phase 3: Integration Tests

```
INTEGRATION TESTS EXECUTION:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Week 9.1 Task (1 hour):
Backend Integration:
  ☐ Flow 1: Auth Complete Cycle → PASS
  ☐ Flow 2: Transaction & Inventory → PASS
  ☐ Flow 3: Forecast Generation → PASS
  ☐ Flow 4: Anomaly Detection → PASS
  ☐ Flow 5: Webhook Event Publishing → PASS

Frontend E2E:
  ☐ Flow 1: Registration & Login → PASS
  ☐ Flow 2: Transaction Creation → PASS
  ☐ Flow 3: Forecast Generation → PASS
  ☐ Flow 4: Anomaly Acknowledge & Resolve → PASS
  ☐ Flow 5: Dark Mode Persistence → PASS

COMMAND: go test -v ./tests/integration/...
         npm run test:e2e
EXPECTED RESULT: All flows pass, 100% success rate
```

---

### Phase 4: E2E Tests

```
E2E TESTS EXECUTION:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Week 9.2 Task (30 min):
  ☐ Auth tests (5 cases) → PASS
  ☐ Inventory tests (5 cases) → PASS
  ☐ Forecast tests (5 cases) → PASS
  ☐ Anomaly tests (6 cases) → PASS
  ☐ Settings tests (6 cases) → PASS

TOTAL: 27 E2E test cases (documented, repeatable)
COMMAND: npx playwright test
EXPECTED RESULT: All scenarios pass on Chrome, Safari, Firefox
```

---

### Phase 5: Performance Tests

```
PERFORMANCE TESTS EXECUTION:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Week 9.6 Task (1 hour):
Backend Load:
  ☐ k6: Inventory query (1000 users) → p99 <100ms ✓
  ☐ k6: Forecast batch (1000 SKUs) → <5s ✓
  ☐ k6: Concurrent transactions → 100 users, no deadlock ✓
  ☐ k6: Rate limiting → 429 when exceeded ✓

Frontend Performance:
  ☐ Lighthouse: Performance >90 ✓
  ☐ Lighthouse: Accessibility >90 ✓
  ☐ Lighthouse: Best Practices >90 ✓
  ☐ Lighthouse: SEO >90 ✓
  ☐ Core Web Vitals: LCP <2.5s ✓
  ☐ Core Web Vitals: FID <100ms ✓
  ☐ Core Web Vitals: CLS <0.1 ✓

COMMAND: npm run lighthouse
         k6 run tests/load/inventory_load_test.js
EXPECTED RESULT: All targets met
```

---

## 📊 FINAL TEST SUMMARY

```
COMPREHENSIVE TEST COVERAGE SUMMARY:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

UNIT TESTS:
  Backend:    246 test cases  | Target: >90% coverage ✓
  Frontend:   311 test cases  | Target: >50% coverage ✓
  SUBTOTAL:   557 test cases

INTEGRATION TESTS:
  Backend:    72 assertions   | 5 critical flows
  Frontend:   68+ user actions| 5+ E2E scenarios
  SUBTOTAL:   140+ verifications

E2E TESTS:
  27 test cases              | Browser compatibility verified
  5+ critical user journeys  | Repeatable, documented

PERFORMANCE TESTS:
  Load: 1000+ concurrent users
  Latency: p99 <100ms inventory queries
  Throughput: 1000+ req/sec
  Lighthouse: >90 all metrics

SECURITY TESTS:
  OWASP Top 10 coverage
  JWT validation
  RBAC enforcement
  Rate limiting
  SQL injection prevention
  CORS validation

TOTAL TEST COVERAGE:
  Backend: >90% unit + 5 integration flows
  Frontend: >50% unit + 5+ E2E flows
  Performance: Load tested to 1000 users
  Security: Full OWASP coverage

SUCCESS CRITERIA:
  ✅ All 557 unit tests passing
  ✅ All 5+ integration flows passing
  ✅ All 5+ E2E scenarios passing
  ✅ All performance targets met
  ✅ Zero critical security issues
  ✅ >99% uptime in production
```

---

## 🎯 TEST EXECUTION TIMELINE (Week 9)

```
WEEK 9 TESTING SCHEDULE:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

MONDAY (9.1 - 4 hours):
  08:00-09:00 | Unit test execution (backend)
  09:00-10:00 | Integration test execution
  10:00-11:00 | Coverage report analysis
  11:00-12:00 | Bug fixes (if needed)

TUESDAY (9.2 - 2 hours):
  08:00-09:00 | Frontend unit tests
  09:00-10:00 | E2E test execution

WEDNESDAY (9.3 - 2 hours):
  08:00-09:00 | Documentation review
  09:00-10:00 | Test report generation

THURSDAY (9.4 - 5 hours):
  08:00-10:00 | Docker build & test
  10:00-12:00 | Staging deployment
  12:00-13:00 | Production deployment

FRIDAY (9.5-9.7 - 7 hours):
  08:00-09:00 | Security testing
  09:00-10:00 | Monitoring setup
  10:00-11:00 | Performance optimization
  11:00-12:00 | Portfolio materials
  13:00-15:00 | Final QA & release
  15:00-16:00 | v1.0.0 tag & release notes

TOTAL: 20 hours (1 week)
```

---

## ✅ FINAL VERIFICATION CHECKLIST

```
PRE-RELEASE VERIFICATION:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

TESTING COMPLETE:
  ☐ 246 backend unit tests passing (>90% coverage)
  ☐ 311 frontend unit tests passing (>50% coverage)
  ☐ 5 backend integration flows passing
  ☐ 27 E2E test cases passing
  ☐ Load tests: 1000 concurrent users, p99 <100ms
  ☐ Performance tests: Lighthouse >90
  ☐ Security: OWASP Top 10 verified

DOCUMENTATION COMPLETE:
  ☐ README.md
  ☐ API.md (23 endpoints)
  ☐ Deployment Guide
  ☐ Architecture Documentation
  ☐ Troubleshooting Guide
  ☐ Contributing Guidelines

DEPLOYMENT READY:
  ☐ Docker images built (<150MB backend, <100MB frontend)
  ☐ CI/CD pipeline passing
  ☐ Staging environment live & tested
  ☐ Production environment configured
  ☐ Health checks operational
  ☐ Monitoring & logging active

PORTFOLIO READY:
  ☐ GitHub repository complete
  ☐ Live demo accessible
  ☐ Deployment documented
  ☐ Resume updated
  ☐ v1.0.0 release tag created

PROJECT COMPLETE:
  ✅ All requirements met
  ✅ Production-grade quality
  ✅ Ready for deployment
  ✅ Ready for interviews/portfolio
```

---

**Document:** COMPREHENSIVE TEST PLAN - ScIntel Platform  
**Status:** 🟢 READY FOR EXECUTION  
**Total Test Cases:** 900+  
**Coverage Target:** Backend >90%, Frontend >50%  
**Expected Duration:** Week 9 (20 hours)  
**Created:** 2026-07-29
