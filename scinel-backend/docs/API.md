# Supply Chain Intelligence Platform - API Documentation

> **Version:** 1.0.0
> **Base URL:** http://localhost:8080/api/v1
> **Authentication:** Bearer JWT Token (except public endpoints)

---

## 1. Authentication

### 1.1 Register User
**Endpoint:** POST /auth/register
**Public:** ✅ Yes
**Description:** Create a new user account

**Request Body:**
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}

**Response (201 Created):**
{
  "code": 201,
  "message": "User registered successfully",
  "data": {
    "id": "user-uuid",
    "email": "user@example.com",
    "role": "viewer"
  },
  "timestamp": "2026-07-20T14:30:00Z",
  "request_id": "req-uuid"
}

---

## 2. Health Checks

### 2.1 Health (Liveness)
**Endpoint:** GET /health
**Public:** ✅ Yes
**Description:** Basic liveness check (always returns 200)

**Response (200 OK):**
{
  "code": 200,
  "message": "OK",
  "data": {
    "status": "healthy",
    "timestamp": "2026-07-20T14:30:00Z",
    "version": "1.0.0"
  },
  "timestamp": "2026-07-20T14:30:00Z",
  "request_id": "req-uuid"
}

---

## 3. Inventory

### 3.1 Query Inventory
**Endpoint:** GET /inventory
**Public:** ❌ Requires Auth
**Description:** Get inventory levels with filters

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| "warehouse_id" | UUID | Filter by warehouse |
| "sku_id" | UUID | Filter by SKU |
| "status" | string | "normal", "low_stock", "critical" |
| "limit" | int | Default 20, max 100 |
| "offset" | int | Default 0 |

**Response (200 OK):**
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
        "updated_at": "2026-07-20T14:30:00Z"
      }
    ],
    "total": 324,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-07-20T14:30:00Z",
  "request_id": "req-uuid"
}

---

## 4. Forecasting

### 4.1 Query Forecasts
**Endpoint:** GET /forecasts
**Public:** ❌ Requires Auth
**Description:** Get demand forecasts with filters

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| "warehouse_id" | UUID | Filter by warehouse |
| "sku_id" | UUID | Filter by SKU |
| "days" | int | Forecast horizon (1-90, default 30) |
| "limit" | int | Default 20, max 100 |
| "offset" | int | Default 0 |

**Response (200 OK):**
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
        "generated_at": "2026-07-20T02:00:00Z"
      }
    ],
    "total": 45,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-07-20T14:30:00Z",
  "request_id": "req-uuid"
}

---

## 5. Anomaly Detection

### 5.1 List Anomalies
**Endpoint:** GET /anomalies
**Public:** ❌ Requires Auth
**Description:** List detected anomalies with filters

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| "status" | string | "open", "acknowledged", "resolved" |
| "severity" | string | "low", "medium", "high", "critical" |
| "warehouse_id" | UUID | Filter by warehouse |
| "limit" | int | Default 20, max 100 |
| "offset" | int | Default 0 |

**Response (200 OK):**
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
        "detected_at": "2026-07-20T14:30:00Z"
      }
    ],
    "total": 12,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-07-20T14:30:00Z",
  "request_id": "req-uuid"
}

---

## 6. Webhooks

### 6.1 Create Subscription
**Endpoint:** POST /webhooks/subscribe
**Public:** ❌ Requires Auth
**Description:** Subscribe to webhook events

**Request Body:**
{
  "event_type": "inventory.changed",
  "webhook_url": "https://example.com/webhooks/inventory",
  "max_retries": 5,
  "retry_interval": 60
}

**Response (201 Created):**
{
  "code": 201,
  "message": "Subscription created",
  "data": {
    "id": "sub-uuid",
    "event_type": "inventory.changed",
    "webhook_url": "https://example.com/webhooks/inventory",
    "is_active": true,
    "secret_key": "whsec_abcdef123456...",
    "created_at": "2026-07-20T14:30:00Z"
  },
  "timestamp": "2026-07-20T14:30:00Z",
  "request_id": "req-uuid"
}

---

## 7. API Keys

### 7.1 Create API Key
**Endpoint:** POST /api-keys
**Public:** ❌ Requires Auth
**Description:** Generate programmatic access key

**Request Body:**
{
  "name": "Integration Service",
  "permissions": ["inventory:read", "inventory:write", "forecast:read"],
  "expires_in_days": 365
}

**Response (201 Created):**
{
  "code": 201,
  "message": "API key created",
  "data": {
    "key_id": "key-uuid",
    "key_prefix": "sk_test_1a2b3c",
    "key_secret": "sk_test_1a2b3c_full_secret_only_shown_once",
    "name": "Integration Service",
    "permissions": ["inventory:read", "inventory:write", "forecast:read"],
    "is_active": true,
    "created_at": "2026-07-20T14:30:00Z",
    "expires_at": "2027-07-20T14:30:00Z"
  },
  "timestamp": "2026-07-20T14:30:00Z",
  "request_id": "req-uuid"
}

---

## 8. Error Handling

### 8.1 Standard Error Response
{
  "code": 400,
  "message": "Validation error",
  "details": "SKU length must be 5-20 characters",
  "timestamp": "2026-07-20T14:30:00Z",
  "request_id": "req-uuid"
}

### 8.2 Common Error Codes

| Code | Meaning |
|------|---------|
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 429 | Too Many Requests |
| 500 | Internal Server Error |
| 503 | Service Unavailable |
