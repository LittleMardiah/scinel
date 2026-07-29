================================================================
LOGIC FLOW DOCUMENTATION - SUPPLY CHAIN INTELLIGENCE PLATFORM v1.0.0
Alur Logika Sistem Dari Awal Hingga Akhir
================================================================

## 1. OVERVIEW & SYSTEM ARCHITECTURE

ScIntel adalah **event-driven, asynchronous backend system** dengan 3 core flows:

1. **AUTHENTICATION FLOW** → User login, JWT token generation
2. **INVENTORY FLOW** → Stock in/out, real-time tracking, alerts
3. **FORECASTING FLOW** → Historical data analysis, demand prediction
4. **ANOMALY DETECTION FLOW** → Pattern recognition, alerts
5. **WEBHOOK FLOW** → Event publishing, retry logic, notification

---

## 2. SYSTEM ARCHITECTURE DIAGRAM

```
┌─────────────────────────────────────────────────────────────────────┐
│                           API Client (Postman/Frontend/Mobile)      │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                 ┌─────────▼──────────┐
                 │   API Gateway      │
                 │  (Rate Limiting)   │
                 └─────────┬──────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼──────┐   ┌──────▼───────┐   ┌──────▼───────┐
   │   AUTH    │   │  INVENTORY   │   │ FORECASTING  │
   │  Service  │   │   Service    │   │    Service   │
   │           │   │              │   │              │
   │ • Login   │   │ • Check stock │   │ • ETS Algo   │
   │ • Register│   │ • In/Out txn │   │ • ARIMA Algo │
   │ • Refresh │   │ • Alerts     │   │ • Batch jobs │
   │ • Revoke  │   │ • Snapshots  │   │ • Metrics    │
   └────┬──────┘   └──────┬───────┘   └──────┬───────┘
        │                 │                  │
        │          ┌──────▼───────┐          │
        │          │   ANOMALY    │          │
        │          │   Service    │          │
        │          │              │          │
        │          │ • Z-score    │          │
        │          │ • IQR method │          │
        │          │ • Threshold  │          │
        │          └──────┬───────┘          │
        │                 │                  │
        │     ┌───────────┼──────────────┐   │
        │     │           │              │   │
   ┌────▼─────▼────┐  ┌───▼──────┐  ┌──▼──────┐
   │  Audit Log DB │  │  Cache   │  │ Message │
   │  (PostgreSQL) │  │ (Redis)  │  │ Queue   │
   └───────────────┘  └──────────┘  └──────────┘
        ▲
        │
   ┌────┴─────────────────────┐
   │  WEBHOOK SERVICE         │
   │ (Event Publishing)       │
   │ • Subscribe endpoints    │
   │ • Retry with backoff     │
   │ • HMAC signing           │
   │ • Dead-letter queue      │
   └──────────────────────────┘
```

---

## 3. DETAILED FLOW 1: AUTHENTICATION & AUTHORIZATION

### 3.1 REGISTRATION FLOW

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. USER SUBMITS REGISTRATION REQUEST                            │
│    POST /api/v1/auth/register                                   │
│    { "email": "bambang@warehouse.id", "password": "...", ... }  │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 2. VALIDATION LAYER                                             │
│    • Check email format valid?                                  │
│    • Check password meets requirements (8+ chars, mixed)?       │
│    • Check email already exists in DB?                          │
│    IF validation fails → Return 400 Bad Request + error message │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 3. PASSWORD HASHING (bcrypt)                                    │
│    • Generate salt (rounds: 12)                                 │
│    • Hash password + salt                                       │
│    • Store hashed password (never store plain password!)        │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 4. CREATE USER RECORD IN DATABASE                               │
│    INSERT INTO users (                                          │
│      id, email, password_hash, role, status, created_at         │
│    ) VALUES (...)                                               │
│    Default role: "viewer" (least privilege)                     │
│    Status: "active" or "pending_email_verification"             │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 5. AUDIT LOG ENTRY                                              │
│    INSERT INTO audit_log (                                      │
│      user_id, action, resource, timestamp, ip_address           │
│    ) = ("user-new", "create", "user:email", now(), "...")       │
│    → Traceable for compliance & security review                 │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 6. RETURN SUCCESS RESPONSE                                      │
│    201 Created                                                  │
│    { "user_id": "user-bambang-001",                             │
│      "email": "bambang@warehouse.id",                           │
│      "role": "viewer",                                          │
│      "message": "User created. Please verify email." }          │
└─────────────────────────────────────────────────────────────────┘
```

---

### 3.2 LOGIN & TOKEN GENERATION FLOW

```
┌──────────────────────────────────────────────────────────────┐
│ 1. USER SUBMITS LOGIN REQUEST                                │
│    POST /api/v1/auth/login                                   │
│    { "email": "bambang@warehouse.id", "password": "..." }    │
└──────────┬───────────────────────────────────────────────────┘
           │
┌──────────▼───────────────────────────────────────────────────┐
│ 2. FIND USER BY EMAIL                                        │
│    SELECT * FROM users WHERE email = ?                       │
│    IF NOT FOUND → Return 401 Unauthorized                    │
└──────────┬───────────────────────────────────────────────────┘
           │
┌──────────▼───────────────────────────────────────────────────┐
│ 3. VERIFY PASSWORD (bcrypt.CompareHashAndPassword)           │
│    Compare provided password vs stored hash                  │
│    IF mismatch → Return 401 Unauthorized                     │
└──────────┬───────────────────────────────────────────────────┘
           │
┌──────────▼───────────────────────────────────────────────────┐
│ 4. CHECK ACCOUNT STATUS                                      │
│    IF status = "suspended" → Return 403 Forbidden            │
│    IF status = "pending_email_verification"                  │
│       → Return 403 + message: "Please verify email"          │
└──────────┬───────────────────────────────────────────────────┘
           │
┌──────────▼───────────────────────────────────────────────────┐
│ 5. GENERATE ACCESS TOKEN (JWT)                               │
│    Header:    { "alg": "HS256", "typ": "JWT" }               │
│    Payload:   {                                              │
│      "user_id": "user-bambang-001",                          │
│      "email": "bambang@warehouse.id",                        │
│      "role": "manager",                                      │
│      "iat": 1721328600 (issued at),                          │
│      "exp": 1721332200 (expiry = 1 hour later)               │
│    }                                                         │
│    Signature: HMACSHA256(header + payload, secret_key)       │
│    → Example token: eyJhbGc...xyz.abcd...1234.sign...        │
└──────────┬───────────────────────────────────────────────────┘
           │
┌──────────▼───────────────────────────────────────────────────┐
│ 6. GENERATE REFRESH TOKEN (JWT - Long-lived)                 │
│    Same structure as access token, but:                      │
│    • exp: 7 days (instead of 1 hour)                         │
│    • Type claim: "refresh" (to distinguish)                  │
│    → Store in secure, httpOnly cookie (if web app)           │
└──────────┬───────────────────────────────────────────────────┘
           │
┌──────────▼───────────────────────────────────────────────────┐
│ 7. CACHE SESSION IN REDIS (optional, for fast lookups)       │
│    SET cache:session:{user_id}:{access_token_jti}            │
│    EXPIRY: 1 hour (match access token)                       │
│    → Allows quick token blacklist on logout                  │
└──────────┬───────────────────────────────────────────────────┘
           │
┌──────────▼───────────────────────────────────────────────────┐
│ 8. LOG LOGIN EVENT                                           │
│    INSERT INTO audit_log (                                   │
│      user_id, action, resource, timestamp, ip_address        │
│    ) = ("user-bambang-001", "login", "auth", now(), "192...")│
└──────────┬───────────────────────────────────────────────────┘
           │
┌──────────▼───────────────────────────────────────────────────┐
│ 9. RETURN TOKENS TO CLIENT                                   │
│    200 OK                                                    │
│    {                                                         │
│      "access_token": "eyJhbGc...",                           │
│      "refresh_token": "eyJhbGc...",                          │
│      "token_type": "Bearer",                                 │
│      "expires_in": 3600,                                     │
│      "user": {                                               │
│        "id": "user-bambang-001",                             │
│        "email": "bambang@warehouse.id",                      │
│        "role": "manager"                                     │
│      }                                                       │
│    }                                                         │
└──────────────────────────────────────────────────────────────┘
```

---

### 3.3 API REQUEST WITH JWT AUTHENTICATION

```
┌──────────────────────────────────────────────────────────────┐
│ 1. CLIENT SENDS AUTHENTICATED REQUEST                        │
│    GET /api/v1/inventory                                     │
│    Authorization: Bearer eyJhbGc...                          │
│    Content-Type: application/json                           │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 2. MIDDLEWARE: EXTRACT TOKEN                                │
│    Parse Authorization header                               │
│    Extract token from "Bearer {token}"                      │
│    IF header missing/malformed → Return 401 Unauthorized    │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 3. MIDDLEWARE: VALIDATE TOKEN SIGNATURE                      │
│    Verify JWT signature using secret key                    │
│    IF signature invalid → Return 401 Unauthorized           │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 4. MIDDLEWARE: CHECK TOKEN EXPIRY                            │
│    IF exp < current_time → Return 401 Token Expired         │
│    (Suggest client to use refresh token to get new token)   │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 5. MIDDLEWARE: CHECK TOKEN BLACKLIST (Redis)                │
│    GET cache:blacklist:{token_jti}                          │
│    IF exists → Return 401 Token Revoked                     │
│    (User logged out / token manually revoked)               │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 6. EXTRACT CLAIMS & CONTEXT                                 │
│    user_id = token.claims["user_id"]                        │
│    role = token.claims["role"]                              │
│    Store in request context (Go's context.Context)          │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 7. MIDDLEWARE: RBAC CHECK                                   │
│    Based on endpoint + HTTP method + user role:             │
│    • GET /api/v1/inventory → allowed for ["manager",        │
│                               "analyst", "admin"]            │
│    IF role not in allowed list → Return 403 Forbidden       │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 8. PROCEED TO HANDLER                                       │
│    Request is authenticated & authorized                    │
│    Handler can safely read user_id from context             │
│    Continue to next flow (e.g., INVENTORY FLOW)             │
└──────────────────────────────────────────────────────────────┘
```

---

## 4. DETAILED FLOW 2: INVENTORY IN/OUT TRANSACTION

### 4.1 STOCK IN TRANSACTION FLOW

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. WAREHOUSE MANAGER SUBMITS STOCK IN REQUEST                   │
│    POST /api/v1/transactions                                    │
│    {                                                            │
│      "warehouse_id": "wh-jakarta-001",                          │
│      "sku": "SKU-NIKE-001",                                     │
│      "quantity": 50,                                            │
│      "type": "in",                                              │
│      "reason": "Supplier delivery #PO-2026-001"                 │
│    }                                                            │
│    Authorization: Bearer {token}                                │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 2. VALIDATION LAYER                                             │
│    • Check warehouse exists? → Query DB                        │
│    • Check SKU exists? → Query DB                              │
│    • Check quantity > 0?                                       │
│    • Check user has permission for this warehouse?             │
│    IF validation fails → Return 400 Bad Request                │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 3. START DATABASE TRANSACTION (ACID)                            │
│    BEGIN TRANSACTION;                                          │
│    → Ensure atomicity: either all succeed or all rollback      │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 4. LOCK INVENTORY ROW (Pessimistic)                             │
│    SELECT * FROM inventory                                     │
│    WHERE warehouse_id = ? AND sku = ?                          │
│    FOR UPDATE;  ← Locks row until transaction commits         │
│    → Prevents race conditions (concurrent updates)             │
│    Retrieve: previous_stock_level = 300                        │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 5. UPDATE INVENTORY QUANTITY                                    │
│    UPDATE inventory                                            │
│    SET quantity = quantity + 50,                               │
│        last_updated = NOW()                                    │
│    WHERE warehouse_id = ? AND sku = ?;                         │
│    new_stock_level = 350                                       │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 6. LOG TRANSACTION IN AUDIT TABLE                               │
│    INSERT INTO transactions (                                  │
│      id, warehouse_id, sku, type, quantity_delta,              │
│      previous_stock, new_stock, user_id, reason, timestamp     │
│    ) VALUES (                                                  │
│      'txn-abc123', 'wh-jakarta-001', 'SKU-NIKE-001',           │
│      'in', 50, 300, 350, 'user-bambang-001',                  │
│      'Supplier delivery #PO-2026-001', '2026-07-18 14:30:00'   │
│    );                                                          │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 7. INVALIDATE CACHE (Redis)                                     │
│    DEL cache:inventory:{warehouse_id}:{sku}                    │
│    → Force re-query DB on next GET (stale data prevention)     │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 8. TRIGGER WEBHOOK EVENTS (Async)                              │
│    Publish events to queue:                                    │
│    • Event: "inventory.changed"                                │
│    • Event: "forecast.trigger.scheduled"  (optional)           │
│    → Add to message queue (in-memory channel or Redis queue)   │
│    → Webhooks will be processed asynchronously                 │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 9. TRIGGER ANOMALY CHECK (Async Background Job)                │
│    Queue background task: CheckAnomalies(warehouse_id, sku)    │
│    → Will run asynchronously in worker goroutine               │
│    → See ANOMALY DETECTION FLOW for details                    │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 10. COMMIT DATABASE TRANSACTION                                 │
│     COMMIT;                                                    │
│     → All changes are now permanent & visible to other clients │
└─────────────┬───────────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────────┐
│ 11. RETURN SUCCESS RESPONSE (201 Created)                       │
│     {                                                          │
│       "transaction_id": "txn-abc123",                          │
│       "warehouse_id": "wh-jakarta-001",                        │
│       "sku": "SKU-NIKE-001",                                   │
│       "quantity_delta": 50,                                    │
│       "previous_stock_level": 300,                             │
│       "new_stock_level": 350,                                  │
│       "timestamp": "2026-07-18T14:30:00Z",                     │
│       "created_by": "user-bambang-001",                        │
│       "webhook_events_triggered": [                            │
│         "inventory.changed",                                   │
│         "forecast.trigger.scheduled"                           │
│       ]                                                        │
│     }                                                          │
└──────────────────────────────────────────────────────────────────┘
```

---

### 4.2 REAL-TIME STOCK LEVEL QUERY FLOW

```
┌──────────────────────────────────────────────────────────────┐
│ 1. CLIENT REQUESTS CURRENT STOCK LEVEL                       │
│    GET /api/v1/inventory?sku=SKU-NIKE-001                   │
│    Authorization: Bearer {token}                            │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 2. CHECK REDIS CACHE FIRST (Fast path - <1ms)              │
│    GET cache:inventory:SKU-NIKE-001                         │
│    IF cache HIT → Return cached data immediately            │
│       (Cache TTL: 30 seconds, hit rate: ~80-90%)            │
│    IF cache MISS → Continue to step 3                       │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 3. QUERY DATABASE (Slow path - ~10-50ms)                    │
│    SELECT warehouse_id, quantity                            │
│    FROM inventory                                           │
│    WHERE sku = 'SKU-NIKE-001'                               │
│    ORDER BY warehouse_id;                                   │
│                                                              │
│    Result:                                                  │
│    ┌──────────────────┬──────────┐                          │
│    │ warehouse_id     │ quantity │                          │
│    ├──────────────────┼──────────┤                          │
│    │ wh-jakarta-001   │   150    │                          │
│    │ wh-surabaya-001  │    45    │                          │
│    │ wh-bandung-001   │   200    │                          │
│    └──────────────────┴──────────┘                          │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 4. AGGREGATE RESULTS IN MEMORY (Golang)                     │
│    total_stock = 150 + 45 + 200 = 395                       │
│    warehouse_breakdown = [                                  │
│      { warehouse_id: "wh-jakarta-001", quantity: 150 },     │
│      { warehouse_id: "wh-surabaya-001", quantity: 45 },     │
│      { warehouse_id: "wh-bandung-001", quantity: 200 }      │
│    ]                                                        │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 5. CACHE RESULT IN REDIS (Fast future lookups)              │
│    SET cache:inventory:SKU-NIKE-001                         │
│    {                                                        │
│      "sku": "SKU-NIKE-001",                                 │
│      "total_stock": 395,                                    │
│      "warehouse_breakdown": [...],                          │
│      "last_updated": "2026-07-18T14:35:00Z"                 │
│    }                                                        │
│    EXPIRY: 30 seconds (or until inventory changes)          │
│    → Next query within 30 sec hits cache (instant)          │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 6. RETURN RESPONSE (200 OK)                                 │
│    {                                                        │
│      "sku": "SKU-NIKE-001",                                 │
│      "total_stock": 395,                                    │
│      "warehouse_breakdown": [                               │
│        {                                                    │
│          "warehouse_id": "wh-jakarta-001",                  │
│          "warehouse_name": "Jakarta Main",                  │
│          "quantity": 150                                    │
│        },                                                   │
│        {                                                    │
│          "warehouse_id": "wh-surabaya-001",                 │
│          "warehouse_name": "Surabaya Branch",               │
│          "quantity": 45                                     │
│        },                                                   │
│        {                                                    │
│          "warehouse_id": "wh-bandung-001",                  │
│          "warehouse_name": "Bandung Branch",                │
│          "quantity": 200                                    │
│        }                                                    │
│      ],                                                     │
│      "last_updated": "2026-07-18T14:35:00Z",                │
│      "timestamp": "2026-07-18T14:35:02Z"                    │
│    }                                                        │
│    Response time: p99 <100ms target (usually ~5-20ms)       │
└──────────────────────────────────────────────────────────────┘
```

---

## 5. DETAILED FLOW 3: DEMAND FORECASTING

### 5.1 FORECAST GENERATION PIPELINE

```
┌──────────────────────────────────────────────────────────────┐
│ 1. SCHEDULED TRIGGER (Daily @ 2:00 AM)                      │
│    Cron Job: "0 2 * * *" (every day at 2 AM)               │
│    Goroutine: ForecastScheduler runs in background           │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 2. FETCH ALL ACTIVE SKUS & WAREHOUSES                        │
│    SELECT DISTINCT sku, warehouse_id                         │
│    FROM inventory                                            │
│    WHERE status = 'active'                                   │
│    LIMIT 10000; ← Forecast up to 10K SKU-warehouse combos    │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 3. FOR EACH SKU-WAREHOUSE COMBINATION (Parallel Processing)  │
│    Spawn goroutines to process multiple SKUs concurrently    │
│    Max workers: 100 goroutines (control resource usage)      │
└──────┬───────────────────────────────────────────────────────┘
       │
    ┌──┴───────────────────────────────────────────────────────┐
    │                                                           │
    └─────┬────────────────────────┬─────────────────┬────────┘
          │                        │                 │
    (Worker 1)            (Worker 2)           (Worker 100)
          │                        │                 │
    ┌─────▼────────────────────────▼────────┬────────▼────────┐
    │ 4. FETCH HISTORICAL DATA (6-12 months) │                │
    │    SELECT timestamp, quantity          │                │
    │    FROM transactions                   │                │
    │    WHERE warehouse_id = ? AND sku = ?  │                │
    │    AND timestamp > NOW() - INTERVAL 12 MONTH            │
    │    ORDER BY timestamp ASC;             │                │
    │    → Result: 180-360 daily data points │                │
    └─────┬─────────────────────────────────────┬──────────────┘
          │                                     │
    ┌─────▼─────────────────────────────────────▼──────────────┐
    │ 5. DATA PREPROCESSING (Golang)                           │
    │    • Remove outliers (rolling z-score)                   │
    │    • Handle missing data (interpolation)                 │
    │    • Aggregate to daily level (if hourly data)           │
    │    • Normalize: scale to 0-1 range                       │
    │    → Clean time series: [100, 105, 110, 108, ...]        │
    └─────┬──────────────────────────────────────────────────────┘
          │
    ┌─────▼──────────────────────────────────────────────────────┐
    │ 6. FEATURE EXTRACTION (Time-Series Analysis)               │
    │    • Trend: Upward/downward/stable                        │
    │    • Seasonality: Weekly/monthly patterns                 │
    │    • Cyclicity: Recurring patterns (e.g., holiday sales)  │
    │    • ACF/PACF: Auto-correlation analysis                  │
    │    → Identify which algo will work best                   │
    └─────┬──────────────────────────────────────────────────────┘
          │
    ┌─────▼──────────────────────────────────────────────────────┐
    │ 7. ALGORITHM 1: EXPONENTIAL SMOOTHING (ETS)                │
    │    Simple Exponential Smoothing:                          │
    │      Forecast = α × last_value + (1 - α) × prev_forecast  │
    │      α (alpha) = 0.3 (smoothing parameter)                │
    │                                                           │
    │    Double Exponential Smoothing (Holt's):                 │
    │      Level = α × value + (1 - α) × (level + trend)        │
    │      Trend = β × (level - prev_level) + (1 - β) × trend   │
    │      β (beta) = 0.2 (trend parameter)                     │
    │                                                           │
    │    Result: forecast_ets = [150, 155, 160, 158, ...]       │
    │    Computation time: ~50ms per SKU                        │
    └─────┬──────────────────────────────────────────────────────┘
          │
    ┌─────▼──────────────────────────────────────────────────────┐
    │ 8. ALGORITHM 2: ARIMA (Advanced)                           │
    │    AutoRegressive: Use past values                        │
    │      AR(p): value_t depends on value_t-1, value_t-2, ...  │
    │      p = 3 (use past 3 days)                              │
    │                                                           │
    │    Integrated: Differencing for stationarity               │
    │      d = 1 (1st-order differencing)                       │
    │                                                           │
    │    Moving Average: Use past errors                        │
    │      MA(q): value_t depends on error_t-1, error_t-2, ...  │
    │      q = 2 (use past 2 errors)                            │
    │                                                           │
    │    ARIMA(3,1,2) formula applied...                        │
    │    Result: forecast_arima = [148, 153, 158, 156, ...]     │
    │    Computation time: ~200ms per SKU (more complex)        │
    └─────┬──────────────────────────────────────────────────────┘
          │
    ┌─────▼──────────────────────────────────────────────────────┐
    │ 9. CONFIDENCE INTERVALS (95%)                             │
    │    For ETS:                                               │
    │      lower_bound = forecast - 1.96 × std_error            │
    │      upper_bound = forecast + 1.96 × std_error            │
    │      std_error = std_dev(last 30 residuals)               │
    │                                                           │
    │    Result:                                                │
    │    forecast: 150, lower: 120, upper: 180                 │
    │    → 95% confidence demand between 120-180 units          │
    └─────┬──────────────────────────────────────────────────────┘
          │
    ┌─────▼──────────────────────────────────────────────────────┐
    │ 10. ACCURACY METRICS (Backtest on last 30 days)           │
    │     Calculate accuracy if we ran forecast 30 days ago     │
    │                                                           │
    │     MAE (Mean Absolute Error):                            │
    │       avg(|predicted - actual|) = 10.2 units             │
    │                                                           │
    │     RMSE (Root Mean Squared Error):                       │
    │       sqrt(avg((predicted - actual)^2)) = 12.5 units     │
    │                                                           │
    │     MAPE (Mean Absolute Percentage Error):                │
    │       avg(|predicted - actual| / actual) = 8.2%           │
    │       ✅ < 15% = Good forecast                            │
    │                                                           │
    │     Compare both algos → Pick best (lower MAPE)           │
    └─────┬──────────────────────────────────────────────────────┘
          │
    ┌─────▼──────────────────────────────────────────────────────┐
    │ 11. STORE FORECAST IN DATABASE                            │
    │     INSERT INTO forecasts (                               │
    │       id, sku, warehouse_id, date, predicted_demand,      │
    │       lower_bound, upper_bound, confidence, algorithm,    │
    │       generated_at                                        │
    │     ) VALUES (...)  ← For each day in 30-day horizon     │
    │                                                           │
    │     30 rows inserted per SKU-warehouse combination        │
    │                                                           │
    │     INSERT INTO forecast_accuracy (                       │
    │       id, sku, warehouse_id, mape, rmse, mae,            │
    │       algorithm, evaluation_date                          │
    │     ) = ('fa-001', 'SKU-NIKE-001', 'wh-jakarta-001',     │
    │          0.082, 12.5, 10.2, 'exponential_smoothing',     │
    │          '2026-07-18')                                    │
    └─────┬──────────────────────────────────────────────────────┘
          │
    ┌─────▼──────────────────────────────────────────────────────┐
    │ 12. CACHE FORECAST IN REDIS                               │
    │     SET cache:forecast:SKU-NIKE-001:wh-jakarta-001        │
    │     {forecast data}                                       │
    │     EXPIRY: 24 hours (until next forecast run)            │
    └─────┬──────────────────────────────────────────────────────┘
          │
    ┌─────▼──────────────────────────────────────────────────────┐
    │ 13. PUBLISH WEBHOOK EVENT (Async)                         │
    │     Event: "forecast.generated"                           │
    │     Subscribers notified of new forecast                  │
    └─────┬──────────────────────────────────────────────────────┘
          │
    └──────┬────────────────────────────────────────────────────┘
           │
      (Wait for all workers to complete)
           │
    ┌──────▼────────────────────────────────────────────────────┐
    │ 14. BATCH JOB COMPLETION                                  │
    │     Total forecasts generated: 10,000 SKU combinations    │
    │     Time taken: ~5 seconds (parallel processing)          │
    │     Status: ✅ Success                                    │
    │     Log: "Forecast batch completed at 2026-07-18 02:05"   │
    └────────────────────────────────────────────────────────────┘
```

---

### 5.2 GET FORECAST API FLOW

```
┌──────────────────────────────────────────────────────────────────┐
│ 1. CLIENT REQUESTS FORECAST                                      │
│    GET /api/v1/forecasts?sku=SKU-NIKE-001&warehouse_id=wh-j...  │
│                         &horizon=30&confidence=0.95              │
│    Authorization: Bearer {token}                                 │
└──────┬───────────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────────┐
│ 2. PARAMETER VALIDATION                                          │
│    • sku provided? → Required                                    │
│    • warehouse_id provided? → Required                           │
│    • horizon valid? → 1-90 days (default: 30)                   │
│    • confidence valid? → 0.68-0.99 (default: 0.95)              │
└──────┬───────────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────────┐
│ 3. CHECK CACHE (Redis) - Fast path                               │
│    GET cache:forecast:SKU-NIKE-001:wh-jakarta-001                │
│    IF cache HIT (and TTL valid)                                  │
│      → Return cached forecast immediately (<1ms)                 │
│    IF cache MISS → Continue to step 4                            │
└──────┬───────────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────────┐
│ 4. QUERY DATABASE (Slow path)                                    │
│    SELECT date, predicted_demand, lower_bound, upper_bound,     │
│           confidence, algorithm, trend, seasonality_factor       │
│    FROM forecasts                                                │
│    WHERE sku = ? AND warehouse_id = ?                            │
│    AND date BETWEEN TODAY() AND TODAY() + INTERVAL ? DAY         │
│    ORDER BY date ASC                                             │
│    LIMIT 90;  ← Max 90 days                                      │
│                                                                  │
│    Result: 30 rows (one per day)                                 │
└──────┬───────────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────────┐
│ 5. FETCH ACCURACY METRICS                                        │
│    SELECT mape, rmse, mae, algorithm                             │
│    FROM forecast_accuracy                                        │
│    WHERE sku = ? AND warehouse_id = ?                            │
│    ORDER BY evaluation_date DESC                                 │
│    LIMIT 1;                                                      │
│    Result: MAPE = 0.082, RMSE = 12.5, MAE = 10.2                │
└──────┬───────────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────────┐
│ 6. AGGREGATE & FORMAT RESPONSE                                   │
│    total_predicted_demand = sum(all predicted_demand for days)   │
│                           = 150 + 155 + ... + 148 = 4500         │
│    Build forecast array with all metadata                        │
└──────┬───────────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────────┐
│ 7. CACHE RESULT (30 seconds)                                     │
│    SET cache:forecast:SKU-NIKE-001:wh-jakarta-001                │
│    {forecast_response}                                           │
│    EXPIRY: 30 seconds (or until new forecast generated)          │
└──────┬───────────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────────┐
│ 8. RETURN FORECAST RESPONSE (200 OK)                             │
│    {                                                             │
│      "sku": "SKU-NIKE-001",                                      │
│      "warehouse_id": "wh-jakarta-001",                           │
│      "generated_at": "2026-07-18T02:00:00Z",                     │
│      "forecast_horizon": 30,                                     │
│      "total_predicted_demand": 4500,                             │
│      "forecasts": [                                              │
│        {                                                         │
│          "date": "2026-07-19",                                   │
│          "predicted_demand": 150,                                │
│          "lower_bound": 120,                                     │
│          "upper_bound": 180,                                     │
│          "confidence": 0.95,                                     │
│          "algorithm": "exponential_smoothing",                   │
│          "trend": "stable",                                      │
│          "seasonality_factor": 1.05                              │
│        },                                                        │
│        ... (28 more days)                                        │
│      ],                                                          │
│      "model_info": {                                             │
│        "algorithm_used": "exponential_smoothing",                │
│        "accuracy_metrics": {                                     │
│          "mape": 0.082,                                          │
│          "rmse": 12.5,                                           │
│          "mae": 10.2                                             │
│        },                                                        │
│        "historical_data_points": 180,                            │
│        "last_model_update": "2026-07-18T02:00:00Z"               │
│      }                                                           │
│    }                                                             │
│    Response time: p99 <100ms target (usually ~10-30ms)           │
└──────────────────────────────────────────────────────────────────┘
```

---

## 6. DETAILED FLOW 4: ANOMALY DETECTION

### 6.1 ANOMALY CHECK FLOW (Triggered by inventory changes)

```
┌───────────────────────────────────────────────────────────────┐
│ 1. TRIGGER: Inventory changed (stock in/out)                  │
│    ↓ Previous flow: Step 9 of Inventory Transaction           │
│    Queue background task: CheckAnomalies(warehouse_id, sku)   │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 2. BACKGROUND WORKER PROCESSES TASK                           │
│    Goroutine picks task from queue                            │
│    Executes asynchronously (non-blocking)                     │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 3. FETCH RECENT TRANSACTION HISTORY (24 hours)                │
│    SELECT timestamp, quantity_delta, type                     │
│    FROM transactions                                          │
│    WHERE warehouse_id = ? AND sku = ?                         │
│    AND timestamp > NOW() - INTERVAL 24 HOUR                   │
│    ORDER BY timestamp DESC                                    │
│    LIMIT 100;  ← Last 100 transactions (1-2 days)             │
│                                                               │
│    Result: [                                                  │
│      {timestamp: '14:30', delta: +50},  ← Stock in            │
│      {timestamp: '14:00', delta: -5},   ← Stock out           │
│      {timestamp: '13:30', delta: -3},                         │
│      {timestamp: '13:00', delta: -100}, ← ANOMALY!            │
│      ...                                                      │
│    ]                                                          │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 4. CALCULATE STATISTICS (Latest 24 hours)                     │
│    hourly_deltas = [-3, -5, -100, -2, -4, -3, -6, ...]        │
│    mean = -10 (average stock movement per hour)               │
│    std_dev = 25 (standard deviation)                          │
│    min = -100, max = +50                                      │
│    → Understand "normal" behavior                             │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 5. LATEST TRANSACTION (Just occurred)                         │
│    delta = -100 units                                         │
│    → Compare to baseline                                      │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 6. ANOMALY DETECTION: Z-SCORE METHOD                          │
│    Z-score = (value - mean) / std_dev                         │
│            = (-100 - (-10)) / 25                              │
│            = -90 / 25                                         │
│            = -3.6                                             │
│                                                               │
│    Threshold: |Z-score| > 3 (≈99.7% confidence)              │
│    |−3.6| > 3? → YES ✅ ANOMALY DETECTED!                     │
│                                                               │
│    Interpretation:                                           │
│    -100 is 3.6 standard deviations below mean.               │
│    This is unusual (99.7% confidence)                        │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 7. CALCULATE SEVERITY LEVEL                                   │
│    IF |Z-score| > 4 → CRITICAL 🔴                             │
│    IF |Z-score| > 3 → HIGH 🟠                                 │
│    IF |Z-score| > 2 → MEDIUM 🟡                               │
│    IF |Z-score| > 1.5 → LOW 🟢                                │
│                                                               │
│    Result: severity = "HIGH" (Z-score = 3.6)                  │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 8. STORE ANOMALY RECORD IN DATABASE                           │
│    INSERT INTO anomalies (                                    │
│      id, warehouse_id, sku, type, severity,                   │
│      detected_value, baseline_value, z_score,                 │
│      description, detected_at, resolution_status              │
│    ) VALUES (                                                 │
│      'anom-001',                                              │
│      'wh-jakarta-001',                                        │
│      'SKU-NIKE-001',                                          │
│      'unusual_stock_drop',                                    │
│      'high',                                                  │
│      -100,                                                    │
│      -10,                                                     │
│      -3.6,                                                    │
│      'Stock dropped 100 units (3.6σ from mean)',              │
│      NOW(),                                                   │
│      'pending'                                                │
│    );                                                         │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 9. PUBLISH WEBHOOK EVENT (Async)                              │
│    Event: "anomaly.detected"                                  │
│    Severity: "high"                                           │
│    → Subscribers notified immediately                         │
│    → Email alert sent to operations manager                   │
└───┬───────────────────────────────────────────────────────────┘
    │
┌───▼───────────────────────────────────────────────────────────┐
│ 10. ANOMALY ACKNOWLEDGEMENT (Optional, Manual)                │
│     Operations manager investigates:                          │
│     "Nike stock drop 100 units in 1 hour. Why?"               │
│     Possible causes:                                          │
│     • System error → Rollback transaction                     │
│     • Employee theft → File report                            │
│     • Data input error → Correct manually                     │
│     • Legitimate demand spike → No action                     │
│                                                               │
│     PUT /api/v1/anomalies/anom-001                            │
│     {                                                         │
│       "resolution_status": "resolved",                        │
│       "resolution_notes": "Data entry error, quantity −100    │
│                            corrected to −10",                 │
│       "root_cause": "manual_error"                            │
│     }                                                         │
└───┬───────────────────────────────────────────────────────────┘
    │
    └─► Anomaly record updated, case closed
```

---

## 7. DETAILED FLOW 5: WEBHOOK EVENT DELIVERY

### 7.1 WEBHOOK PUBLISHING & DELIVERY FLOW

```
┌──────────────────────────────────────────────────────────────┐
│ 1. EVENT TRIGGERED (From one of previous flows)              │
│    Events: "inventory.changed", "forecast.generated",        │
│             "anomaly.detected"                               │
│                                                              │
│    Example: Inventory transaction completed                  │
│    → Publishes event to message queue                        │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 2. FETCH SUBSCRIPTIONS (Query DB)                            │
│    SELECT id, url, event_type, is_active                     │
│    FROM webhook_subscriptions                                │
│    WHERE event_type = 'inventory.changed'                    │
│    AND is_active = true;                                     │
│                                                              │
│    Result: 5 active subscriptions                            │
│    [                                                         │
│      {id: 'ws-001', url: 'https://erp.company.com/webhooks'}, │
│      {id: 'ws-002', url: 'https://dashboard.company.com/...'},│
│      {id: 'ws-003', url: 'https://mobile-app.company.com/..'},│
│      {id: 'ws-004', url: 'https://email-service.company.com'},│
│      {id: 'ws-005', url: 'https://analytics.company.com/...'}│
│    ]                                                         │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 3. BUILD WEBHOOK PAYLOAD                                     │
│    {                                                         │
│      "event_id": "evt-xyz789",                               │
│      "event_type": "inventory.changed",                      │
│      "timestamp": "2026-07-18T14:30:00Z",                    │
│      "data": {                                               │
│        "transaction_id": "txn-abc123",                       │
│        "warehouse_id": "wh-jakarta-001",                     │
│        "sku": "SKU-NIKE-001",                                │
│        "quantity_delta": 50,                                 │
│        "new_stock_level": 350                                │
│      }                                                       │
│    }                                                         │
│    payload_json = JSON.stringify(payload)                   │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 4. GENERATE SIGNATURE (HMAC-SHA256)                          │
│    signature = HMAC-SHA256(payload_json, webhook_secret_key) │
│    signature = "sha256=abcd1234efgh5678..."                  │
│    → Allows subscriber to verify payload authenticity        │
└──────┬───────────────────────────────────────────────────────┘
       │
┌──────▼───────────────────────────────────────────────────────┐
│ 5. FOR EACH SUBSCRIBER ENDPOINT (Parallel)                   │
│    Spawn goroutine for each webhook delivery                 │
│    → Non-blocking concurrent delivery                        │
└──────┬───────────────────────────────────────────────────────┘
       │
    ┌──┴───────────────────────────────────────────────────────┐
    │                                                           │
    └─────┬────────────────────────────┬──────────────────┬────┘
          │                            │                  │
    (Subscriber 1)         (Subscriber 2)        (Subscriber 5)
          │                            │                  │
    ┌─────▼────────────────────────────▼──────────────┬────▼───┐
    │ 6. ATTEMPT HTTP POST (1st try)                 │        │
    │    POST {subscriber_url}                       │        │
    │    Content-Type: application/json              │        │
    │    X-Webhook-Signature: sha256=abcd...         │        │
    │    {payload_json}                              │        │
    │                                                 │        │
    │    Timeout: 5 seconds                           │        │
    │    Max retries: 5                               │        │
    │    Result: HTTP 200 OK ✅                        │        │
    │            OR timeout/error ❌                  │        │
    └─────┬─────────────────────────────────────────┬──┴────────┘
          │                                         │
          │         (Subscriber responded 200 OK)  │
          │                                        │ (Subscriber didn't respond)
    ┌─────▼────────────────────────────┐    ┌──────▼──────────────────────┐
    │ 7. DELIVERY SUCCESS 🟢             │    │ 7. DELIVERY FAILED ❌       │
    │    Record in DB:                 │    │    Record attempt 1 failed  │
    │    INSERT INTO                   │    │    Wait 1 second...         │
    │      webhook_deliveries (        │    │    Schedule retry attempt 2 │
    │      subscription_id, event_id,  │    │    (exponential backoff)    │
    │      status, timestamp           │    │                            │
    │    ) = ('ws-001', 'evt-xyz789',  │    └──────┬──────────────────────┘
    │         'success', now())         │           │
    │                                   │      (Wait 1 sec)
    │    Mark as delivered              │           │
    │    No further retries             │    ┌──────▼──────────────────────┐
    │                                   │    │ 8. RETRY ATTEMPT 2          │
    │                                   │    │    POST {subscriber_url}    │
    │                                   │    │    Timeout: 5 seconds       │
    │                                   │    │    Response: Still failing  │
    │                                   │    │    Wait 2 seconds...        │
    │                                   │    │    Schedule retry attempt 3 │
    │                                   │    └──────┬──────────────────────┘
    │                                   │           │
    │                                   │      (Wait 2 sec)
    │                                   │           │
    │                                   │    ┌──────▼──────────────────────┐
    │                                   │    │ 9. RETRY ATTEMPT 3-5        │
    │                                   │    │    Same pattern...          │
    │                                   │    │    5 total attempts         │
    │                                   │    │                            │
    │                                   │    │ IF still failing:           │
    │                                   │    │   Add to Dead-Letter Queue  │
    │                                   │    │   Store in DB for manual    │
    │                                   │    │   investigation             │
    │                                   │    └─────────────────────────────┘
    │                                   │
    └─────────────────────────────────────────────────────────┘
```

---

## 8. OVERALL SYSTEM FLOW SEQUENCE DIAGRAM

```
                                                TIMELINE: 1 DAY
┌────────────────────────────────────────────────────────────────────┐

│ 00:00 - FORECAST BATCH JOB RUNS                                    │
│  ├─ Fetch all active SKUs                                          │
│  ├─ Parallel process forecasts (ETS + ARIMA)                       │
│  ├─ Calculate accuracy metrics (MAPE, RMSE, MAE)                   │
│  ├─ Store forecasts in DB (10K SKU combos × 30 days)               │
│  ├─ Cache forecasts in Redis                                       │
│  └─ Publish "forecast.generated" webhooks → Subscribers notified   │
│                                                                     │
│ 08:00 - WAREHOUSE MANAGER STARTS WORK                              │
│  ├─ Logs in (JWT auth) → Access token generated                    │
│  ├─ Checks stock level (real-time API) → Cache hit, instant        │
│  ├─ Records stock in from supplier                                 │
│  │  ├─ Transaction logged in DB                                    │
│  │  ├─ Inventory updated (pessimistic lock)                        │
│  │  ├─ Cache invalidated                                           │
│  │  ├─ Anomaly check triggered (background)                        │
│  │  └─ "inventory.changed" webhook published                       │
│  │                                                                  │
│  ├─ Webhooks delivered to 5 subscribers                            │
│  │  ├─ ERP system updated (success)                                │
│  │  ├─ Mobile app notified (success)                               │
│  │  ├─ Analytics service (retry 3x, then dead-letter queue)        │
│  │  └─ Email service queues alert email                           │
│                                                                     │
│ 09:00 - SUPPLY CHAIN ANALYST CHECKS FORECAST                       │
│  ├─ Logs in (JWT auth)                                             │
│  ├─ Queries forecast API (SKU: Nike, horizon: 30 days)             │
│  │  ├─ Checks Redis cache → Cache miss                            │
│  │  ├─ Queries DB → Retrieves forecast + accuracy                  │
│  │  ├─ Caches result in Redis (30 sec TTL)                         │
│  │  └─ Returns: 4500 units predicted, MAPE 8.2%                    │
│  └─ Views trends: "Nike demand stable, buy 150/day"                │
│                                                                     │
│ 14:30 - UNUSUAL STOCK DROP DETECTED                                │
│  ├─ Manual stock out entry: -100 units Nike (unusual!)             │
│  ├─ Transaction logged                                             │
│  ├─ Inventory updated                                              │
│  ├─ Anomaly detector runs (Z-score: -3.6)                          │
│  │  ├─ Calculates baseline (historical mean: -10 units/hr)         │
│  │  ├─ Detects anomaly (Z-score > 3)                               │
│  │  ├─ Records anomaly: "unusual_stock_drop", severity: HIGH       │
│  │  └─ Publishes "anomaly.detected" webhook                        │
│  │                                                                  │
│  └─ Operations manager receives webhook                            │
│     ├─ Email alert: "ANOMALY: Nike stock dropped 100 units"        │
│     ├─ Reviews investigation link                                  │
│     ├─ Checks warehouse: "Ah, data entry error. Should be -10"     │
│     ├─ Updates anomaly: "resolved", notes: "corrected to -10"      │
│     └─ Case closed ✅                                              │
│                                                                     │
│ 23:00 - DAILY AUDIT & BACKUP                                       │
│  ├─ Generate inventory snapshot (EOD summary)                       │
│  ├─ Backup all transactions to archive table                       │
│  ├─ Verify data integrity (foreign keys, constraints)              │
│  └─ Ready for next day's forecasting cycle                         │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘
```

---

## 9. ERROR HANDLING & RESILIENCE

### 9.1 RETRY LOGIC (Webhook Delivery Example)

```
Attempt 1: T+0s    → Failed (timeout)
           Wait 1 second...

Attempt 2: T+1s    → Failed (500 error)
           Wait 2 seconds...

Attempt 3: T+3s    → Failed (timeout)
           Wait 4 seconds...

Attempt 4: T+7s    → Failed (502 gateway error)
           Wait 8 seconds...

Attempt 5: T+15s   → Failed (timeout)
           → Move to Dead-Letter Queue for manual investigation

DLQ holds event indefinitely:
- Operator reviews: "Why did analytics service fail 5x?"
- Restarts analytics service
- Manually reprocess webhook from DLQ
- Event delivered on 6th attempt
```

### 9.2 CIRCUIT BREAKER PATTERN (Protecting downstream services)

```
NORMAL STATE: Circuit CLOSED
├─ Requests flow normally to external service
├─ Track success/failure ratio
└─ All good ✅

DEGRADED STATE: Circuit OPEN (if failure rate > 50% in last 10 requests)
├─ Stop sending requests to service (fail-fast)
├─ Return cached response or degraded response to client
├─ Log alert: "Circuit opened for {service}"
└─ Protection mode ⚠️

RECOVERY STATE: Circuit HALF_OPEN (after 30 seconds)
├─ Allow 1 test request through
├─ If successful → Reset to CLOSED, resume normal operation
├─ If fails → Go back to OPEN, wait another 30 seconds
└─ Attempting recovery 🔄
```

---

## 10. PERFORMANCE CHARACTERISTICS

### 10.1 LATENCY BY ENDPOINT

| Endpoint | Cache Miss | Cache Hit | Note |
|----------|-----------|-----------|------|
| POST /transactions | 50-100ms | N/A | DB write + invalidation |
| GET /inventory | 20-50ms | <1ms | Mostly DB, cache TTL 30s |
| GET /forecasts | 30-80ms | <1ms | DB query + aggregation, cache TTL 24h |
| GET /anomalies | 10-30ms | <1ms | Simple DB select |
| POST /auth/login | 100-200ms | N/A | Hash verification + token gen |

**Target**: p99 <100ms across all endpoints ✅

### 10.2 DATABASE QUERY OPTIMIZATION

```
Slow Query (Without indexes):
SELECT * FROM transactions 
WHERE warehouse_id = 'wh-001' AND sku = 'SKU-001'
→ Full table scan (1000ms) ❌

Fast Query (With indexes):
CREATE INDEX idx_txn_warehouse_sku 
  ON transactions(warehouse_id, sku);
→ Index seek (5ms) ✅

Result: 200x faster!
```

---

## 11. DATA CONSISTENCY GUARANTEES

### 11.1 ACID TRANSACTIONS

```
Scenario: Concurrent stock in/out

Transaction A: Add 50 units (T=0)
  BEGIN
  SELECT quantity FOR UPDATE (lock row)
  quantity = 300 + 50 = 350
  COMMIT

Transaction B: Remove 30 units (T=5ms)
  BEGIN
  SELECT quantity FOR UPDATE (waits for Transaction A)
  After A commits → quantity = 350
  quantity = 350 - 30 = 320
  COMMIT

Result: Final stock = 320 (correct)
        Both operations atomic ✅
```

---

## 12. KEY TAKEAWAYS (ALUR LOGIKA KESELURUHAN)

```
1. USER AUTHENTICATES
   └─ JWT token issued + cached in Redis

2. USER MAKES REQUEST
   └─ Middleware validates token → Extract claims → Check RBAC

3. INVENTORY OPERATION (Stock in/out)
   └─ DB transaction (pessimistic lock) → Update quantity
     → Log transaction → Invalidate cache
     → Trigger anomaly check (async)
     → Publish webhook events

4. REAL-TIME STOCK QUERY
   └─ Check Redis cache (fast) 
     └─ If miss: Query DB → Aggregate → Cache result

5. FORECASTING (Daily at 2 AM)
   └─ Fetch 6-12 months historical data
     → Preprocess (clean, normalize)
     → Run 2 algorithms: ETS + ARIMA (parallel)
     → Calculate confidence intervals
     → Backtest accuracy (MAPE, RMSE, MAE)
     → Store forecasts + cache

6. GET FORECAST QUERY
   └─ Check Redis cache (fast)
     └─ If miss: Query DB → Format response → Cache

7. ANOMALY DETECTION
   └─ Calculate baseline (historical stats)
     → Z-score analysis
     → If anomaly detected → Store record + Publish webhook

8. WEBHOOK DELIVERY
   └─ Fetch subscriptions → Build payload → Sign (HMAC-SHA256)
     → Deliver to each subscriber (parallel)
     → Retry with exponential backoff
     → Dead-letter queue on final failure

ALL FLOWS ARE:
✅ Stateless & scalable (horizontal)
✅ Async where possible (non-blocking)
✅ ACID guaranteed (database transactions)
✅ Cached for performance (Redis)
✅ Error resilient (retry logic, circuit breaker)
✅ Auditable (detailed logging)
```

---

**Document Status**: ✅ COMPLETE

**Version**: 1.0.0

**Date**: 2026-07-18

================================================================
