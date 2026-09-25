# ScIntel API Documentation

> **Version:** 1.0.0
> **Base URL:** `http://localhost:8080/api/v1`
> **Authentication:** Bearer JWT token (kecuali endpoint publik)
> **Content-Type:** `application/json`

---

## Ringkasan Endpoint

**Authentication (4)**
| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| POST | `/auth/register` | Register user baru | Public |
| POST | `/auth/login` | Login & dapat token | Public |
| POST | `/auth/refresh` | Refresh access token | Public |
| POST | `/auth/logout` | Logout | JWT |

**Health (3)**
| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| GET | `/health` | Liveness check | Public |
| GET | `/readiness` | Deep check (DB + Redis) | Public |
| GET | `/liveness` | Deep check + DB stats | Public |

**Inventory (5)**
| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| GET | `/inventory` | Query inventory dengan filter | JWT |
| POST | `/transactions` | Buat transaksi stok masuk/keluar | JWT |
| GET | `/transactions` | Riwayat transaksi | JWT |
| GET | `/inventory/{warehouse_id}` | Stok per gudang | JWT |
| POST | `/inventory/import` | Bulk import (placeholder) | JWT |

**Forecast (3)**
| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| GET | `/forecasts` | Query forecast dengan filter | JWT |
| POST | `/forecasts/generate` | Trigger generate forecast | Admin |
| GET | `/forecasts/accuracy` | Metrik akurasi | JWT |

**Anomaly (4)**
| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| GET | `/anomalies` | List anomali dengan filter | JWT |
| GET | `/anomalies/{id}` | Detail anomali | JWT |
| PUT | `/anomalies/{id}` | Update status (ack/resolve) | JWT |
| GET | `/anomalies/{id}/history` | Timeline status | JWT |

**API Keys (4)**
| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| POST | `/api-keys` | Create API key | Admin |
| GET | `/api-keys` | List API keys (masked) | Admin |
| DELETE | `/api-keys/{id}` | Revoke API key | Admin |
| POST | `/api-keys/{id}/rotate` | Rotasi secret | Admin |

**Webhooks (4)**
| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| POST | `/webhooks/subscribe` | Subscribe event | JWT |
| GET | `/webhooks/subscriptions` | List subscription | JWT |
| DELETE | `/webhooks/subscriptions/{id}` | Unsubscribe | JWT |
| GET | `/webhooks/deliveries` | Riwayat delivery | JWT |

**Total: 27 routes / 23 endpoint unik**

---

## Format Respons Standar

Semua respons menggunakan wrapper konsisten:

```json
{
  "code": 200,
  "message": "OK",
  "data": { ... },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

## 1. Authentication

### 1.1 Register User

**Endpoint:** `POST /auth/register`
**Auth:** Public

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}
```

**Response (201 Created):**
```json
{
  "code": 201,
  "message": "User registered successfully",
  "data": {
    "id": "user-uuid",
    "email": "user@example.com",
    "role": "viewer"
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `400` (email/password tidak valid), `409` (email duplikat)

---

### 1.2 Login

**Endpoint:** `POST /auth/login`
**Auth:** Public

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 3600,
    "user": {
      "id": "user-uuid",
      "email": "user@example.com",
      "full_name": "John Doe",
      "role": "manager"
    }
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `400` (body invalid), `401` (email/password salah), `403` (akun suspended/pending)

---

### 1.3 Refresh Token

**Endpoint:** `POST /auth/refresh`
**Auth:** Public (menggunakan refresh token)

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Token refreshed",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 3600
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `401` (refresh token invalid/expired)

---

### 1.4 Logout

**Endpoint:** `POST /auth/logout`
**Auth:** JWT

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Logout successful",
  "data": null,
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

## 2. Health Check

### 2.1 Health (Liveness)

**Endpoint:** `GET /health`
**Auth:** Public

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "status": "healthy",
    "timestamp": "2026-09-04T14:30:00Z",
    "version": "1.0.0"
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

### 2.2 Readiness

**Endpoint:** `GET /readiness`
**Auth:** Public
**Deskripsi:** Deep check — memverifikasi koneksi DB & Redis.

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Ready",
  "data": {
    "status": "ready",
    "database": "connected",
    "redis": "connected"
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `503` (DB/Redis tidak tersedia)

---

### 2.3 Liveness

**Endpoint:** `GET /liveness`
**Auth:** Public
**Deskripsi:** Deep check + statistik database.

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Alive",
  "data": {
    "status": "alive",
    "db_stats": {
      "max_open_conns": 50,
      "open_conns": 3,
      "in_use": 2,
      "idle": 1
    }
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `503` (DB tidak tersedia)

---

## 3. Inventory

### 3.1 Query Inventory

**Endpoint:** `GET /inventory`
**Auth:** JWT

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `warehouse_id` | UUID | Filter by gudang |
| `sku_id` | UUID | Filter by SKU |
| `status` | string | `normal`, `low_stock`, `critical` |
| `limit` | int | Default 20, max 100 |
| `offset` | int | Default 0 |

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
        "updated_at": "2026-09-04T14:30:00Z"
      }
    ],
    "total": 324,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

> Cache TTL 30s di Redis. Headers respons menyertakan `X-Cache` saat di-serve dari cache.

---

### 3.2 Create Transaction

**Endpoint:** `POST /transactions`
**Auth:** JWT

**Request Body:**
```json
{
  "warehouse_id": "wh-jakarta-001",
  "sku_id": "sku-nike-001",
  "sku_code": "SKU-NIKE-001",
  "type": "out",
  "quantity": 5,
  "reason": "sales_order",
  "details": "Order #12345"
}
```

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
    "type": "out",
    "quantity": 5,
    "reason": "sales_order",
    "created_at": "2026-09-04T14:30:00Z",
    "created_by": "user-uuid"
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:**
- `400` (validasi gagal)
- `404` (inventory/target tidak ditemukan)
- `409` (stok tidak cukup untuk transaksi keluar)

> Menggunakan **pessimistic locking** — baris inventory di-lock saat transaksi agar tidak ada race condition.

---

### 3.3 Query Transactions (History)

**Endpoint:** `GET /transactions`
**Auth:** JWT

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `warehouse_id` | UUID | Filter by gudang |
| `sku_id` | UUID | Filter by SKU |
| `type` | string | `in`, `out`, `adjustment` |
| `reason` | string | Filter by alasan (ILIKE) |
| `from` | date | Tanggal awal (YYYY-MM-DD) |
| `to` | date | Tanggal akhir |
| `limit` | int | Default 20 |
| `offset` | int | Default 0 |

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "transactions": [
      {
        "id": "txn-uuid",
        "warehouse_id": "wh-jakarta-001",
        "sku_id": "sku-nike-001",
        "sku_code": "SKU-NIKE-001",
        "type": "out",
        "quantity": 5,
        "reason": "sales_order",
        "created_at": "2026-09-04T14:30:00Z",
        "created_by": "user-uuid"
      }
    ],
    "total": 145,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

### 3.4 Get Warehouse Inventory

**Endpoint:** `GET /inventory/{warehouse_id}`
**Auth:** JWT

**Path Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `warehouse_id` | UUID | ID gudang |

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "warehouse_id": "wh-jakarta-001",
    "inventory": [
      {
        "id": "inv-uuid",
        "sku_id": "sku-nike-001",
        "sku_code": "SKU-NIKE-001",
        "quantity": 150,
        "reorder_point": 100,
        "status": "normal"
      }
    ],
    "total_in_warehouse": 25
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `404` (warehouse tidak ditemukan)

---

### 3.5 Bulk Import Inventory

**Endpoint:** `POST /inventory/import`
**Auth:** JWT
**Status:** ⚠️ Placeholder (belum full implementasi)

**Request Body:**
```json
{
  "warehouse_id": "wh-jakarta-001",
  "items": [
    { "sku_code": "SKU-NIKE-001", "quantity": 500 },
    { "sku_code": "SKU-ADIDAS-002", "quantity": 300 }
  ]
}
```

**Response (202 Accepted):**
```json
{
  "code": 202,
  "message": "Import accepted",
  "data": {
    "items_received": 2
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

## 4. Forecasting

### 4.1 Query Forecasts

**Endpoint:** `GET /forecasts`
**Auth:** JWT

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `warehouse_id` | UUID | Filter by gudang |
| `sku_id` | UUID | Filter by SKU |
| `days` | int | Horizon forecast (1-90, default 30) |
| `limit` | int | Default 20 |
| `offset` | int | Default 0 |

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
        "forecast_date": "2026-10-04",
        "predicted_demand": 150,
        "lower_bound": 120,
        "upper_bound": 180,
        "confidence": 0.95,
        "algorithm": "ets",
        "generated_at": "2026-09-04T02:00:00Z"
      }
    ],
    "total": 45,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

> Cache Redis TTL 24 jam. Kolom `algorithm` bernilai `ets` atau `arima`.

---

### 4.2 Generate Forecasts (Manual Trigger)

**Endpoint:** `POST /forecasts/generate`
**Auth:** Admin (RBAC)

**Request Body (opsional):**
```json
{
  "algorithm": "ets"
}
```
Algoritma: `ets` (default), `arima`, atau `auto` (pilih MAPE terbaik).

**Response (202 Accepted):**
```json
{
  "code": 202,
  "message": "Forecast generation started",
  "data": {
    "job_id": "job-uuid",
    "status": "processing"
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `403` (bukan admin)

---

### 4.3 Get Forecast Accuracy

**Endpoint:** `GET /forecasts/accuracy`
**Auth:** JWT

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `warehouse_id` | UUID | Filter by gudang |
| `sku_id` | UUID | Filter by SKU |

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
        "algorithm": "ets",
        "mape": 12.5,
        "rmse": 45.3,
        "mae": 32.1,
        "evaluated_at": "2026-09-04T02:00:00Z"
      },
      {
        "warehouse_id": "wh-jakarta-001",
        "sku_id": "sku-nike-001",
        "algorithm": "arima",
        "mape": 14.2,
        "rmse": 50.1,
        "mae": 35.4,
        "evaluated_at": "2026-09-04T02:00:00Z"
      }
    ]
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

## 5. Anomaly Detection

### 5.1 List Anomalies

**Endpoint:** `GET /anomalies`
**Auth:** JWT

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `status` | string | `open`, `acknowledged`, `resolved` |
| `severity` | string | `low`, `medium`, `high`, `critical` |
| `warehouse_id` | UUID | Filter by gudang |
| `limit` | int | Default 20 |
| `offset` | int | Default 0 |

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
        "detected_at": "2026-09-04T14:30:00Z"
      }
    ],
    "total": 12,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

### 5.2 Get Anomaly Detail

**Endpoint:** `GET /anomalies/{id}`
**Auth:** JWT

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "id": "anomaly-uuid",
    "warehouse_id": "wh-jakarta-001",
    "sku_id": "sku-nike-001",
    "sku_code": "SKU-NIKE-001",
    "type": "unusual_stock_level",
    "severity": "high",
    "description": "Value 250 is 3.45 std deviations from mean 150",
    "z_score": 3.45,
    "iqr_lower": 80,
    "iqr_upper": 220,
    "status": "open",
    "detected_at": "2026-09-04T14:30:00Z",
    "acknowledged_at": null,
    "resolved_at": null
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `404` (anomali tidak ditemukan)

---

### 5.3 Update Anomaly Status

**Endpoint:** `PUT /anomalies/{id}`
**Auth:** JWT

**Request Body:**
```json
{
  "status": "acknowledged",
  "notes": "Sedang diinvestigasi oleh tim gudang"
}
```

Status: `acknowledged` → `resolved` (harus berurutan).

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Anomaly status updated",
  "data": {
    "id": "anomaly-uuid",
    "status": "acknowledged",
    "acknowledged_at": "2026-09-04T15:00:00Z"
  },
  "timestamp": "2026-09-04T15:00:00Z",
  "request_id": "req-uuid"
}
```

**Error:**
- `400` (status tidak valid)
- `403` (tidak berhak)
- `404` (tidak ditemukan)
- `409` (transisi status tidak valid / sudah resolved)

---

### 5.4 Get Anomaly History (Timeline)

**Endpoint:** `GET /anomalies/{id}/history`
**Auth:** JWT

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "history": [
      {
        "status": "open",
        "timestamp": "2026-09-04T14:30:00Z",
        "notes": "Anomaly detected automatically"
      },
      {
        "status": "acknowledged",
        "timestamp": "2026-09-04T15:00:00Z",
        "notes": "Sedang diinvestigasi oleh tim gudang"
      }
    ]
  },
  "timestamp": "2026-09-04T15:00:00Z",
  "request_id": "req-uuid"
}
```

---

## 6. API Keys (Admin Only)

### 6.1 Create API Key

**Endpoint:** `POST /api-keys`
**Auth:** Admin (RBAC)

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
    "key_secret": "sk_test_1a2b3c_full_secret_only_shown_once",
    "name": "Integration Service",
    "permissions": ["inventory:read", "inventory:write", "forecast:read"],
    "is_active": true,
    "created_at": "2026-09-04T14:30:00Z",
    "expires_at": "2027-09-04T14:30:00Z"
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

> **PENTING:** `key_secret` hanya ditampilkan sekali. Secret di-hash SHA-256 di database.

**Error:** `403` (bukan admin)

---

### 6.2 List API Keys

**Endpoint:** `GET /api-keys`
**Auth:** Admin (RBAC)

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "api_keys": [
      {
        "key_id": "key-uuid",
        "key_prefix": "sk_test_1a2b3c",
        "name": "Integration Service",
        "permissions": ["inventory:read"],
        "is_active": true,
        "created_at": "2026-09-04T14:30:00Z",
        "expires_at": "2027-09-04T14:30:00Z",
        "last_used_at": null
      }
    ],
    "total": 3
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

> Sandi penuh (`key_secret`) tidak pernah dikembalikan.

---

### 6.3 Revoke API Key

**Endpoint:** `DELETE /api-keys/{id}`
**Auth:** Admin (RBAC)

**Response (204 No Content):**
_(body kosong)_

**Error:** `400` (id invalid), `403` (bukan admin), `404` (tidak ditemukan)

---

### 6.4 Rotate API Key

**Endpoint:** `POST /api-keys/{id}/rotate`
**Auth:** Admin (RBAC)

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "API key rotated",
  "data": {
    "key_id": "key-uuid",
    "key_prefix": "sk_test_9z8y7x",
    "key_secret": "sk_test_9z8y7x_new_secret_only_shown_once"
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

**Error:** `403` (bukan admin), `404` (tidak ditemukan), `409` (key non-aktif)

---

## 7. Webhooks

### 7.1 Create Subscription

**Endpoint:** `POST /webhooks/subscribe`
**Auth:** JWT

**Request Body:**
```json
{
  "event_type": "inventory.changed",
  "webhook_url": "https://example.com/webhooks/inventory",
  "max_retries": 5,
  "retry_interval": 60
}
```

Event types: `inventory.changed`, `anomaly.detected`, `forecast.generated`.

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
    "secret_key": "whsec_abcdef123456...",
    "created_at": "2026-09-04T14:30:00Z"
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

> `secret_key` digunakan untuk memverifikasi signature HMAC-SHA256 di payload.

---

### 7.2 List Subscriptions

**Endpoint:** `GET /webhooks/subscriptions`
**Auth:** JWT

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
        "created_at": "2026-09-04T14:30:00Z"
      }
    ],
    "total": 2
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

### 7.3 Delete Subscription (Unsubscribe)

**Endpoint:** `DELETE /webhooks/subscriptions/{id}`
**Auth:** JWT

**Response (204 No Content):**
_(body kosong)_

**Error:** `403` (bukan pemilik), `404` (tidak ditemukan)

---

### 7.4 List Deliveries

**Endpoint:** `GET /webhooks/deliveries`
**Auth:** JWT

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `subscription_id` | UUID | Filter by subscription |
| `status` | string | `pending`, `success`, `failed`, `dead_letter` |
| `limit` | int | Default 20 |
| `offset` | int | Default 0 |

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
        "event_type": "inventory.changed",
        "status": "success",
        "http_status_code": 200,
        "attempt_count": 1,
        "last_attempt_at": "2026-09-04T14:31:00Z",
        "next_retry_at": null
      }
    ],
    "total": 15,
    "limit": 20,
    "offset": 0
  },
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

---

## 8. Webhook Payload Format

Payload dikirim ke `webhook_url` dengan **POST**:

**Headers:**
```
X-ScIntel-Event: inventory.changed
X-ScIntel-Signature: sha256=<hmac_hex>
X-ScIntel-Delivery: delivery-uuid
Content-Type: application/json
```

**Body:**
```json
{
  "id": "delivery-uuid",
  "event": "inventory.changed",
  "timestamp": "2026-09-04T14:30:00Z",
  "data": {
    "warehouse_id": "wh-jakarta-001",
    "sku_id": "sku-nike-001",
    "quantity": 150,
    "previous_quantity": 155
  }
}
```

Compute signature: `HMAC-SHA256(secret_key, "<timestamp>.<payload_body>")`.

---

## 9. Error Handling

### 9.1 Format Error Standar

```json
{
  "code": 400,
  "message": "Validation error",
  "details": "SKU length must be 5-20 characters",
  "timestamp": "2026-09-04T14:30:00Z",
  "request_id": "req-uuid"
}
```

### 9.2 Kode Error Umum

| Kode | Arti | Kapan Terjadi |
|------|------|---------------|
| 400 | Bad Request | Validasi body/query gagal |
| 401 | Unauthorized | Token JWT invalid/expired |
| 403 | Forbidden | Role tidak cukup (RBAC) / akun suspended |
| 404 | Not Found | Resource tidak ditemukan |
| 409 | Conflict | Email duplikat, stok tidak cukup, transisi status invalid |
| 429 | Too Many Requests | Rate limit terlampaui (1000 req/s) |
| 500 | Internal Server Error | Error server/DB |
| 503 | Service Unavailable | DB/Redis down |

### 9.3 Rate Limit Response (429)

Headers saat rate limit dilanggar:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1
Retry-After: 1
```

---

## 10. Autentikasi via API

Dua cara autentikasi:

1. **Bearer JWT** (untuk pengguna web):
   ```
   Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
   ```

2. **X-API-Key** (untuk integrasi programmatic):
   ```
   X-API-Key: sk_test_1a2b3c...
   ```

JWT dihasilkan dari `/auth/login` dan berlaku 1 jam. Gunakan `/auth/refresh` untuk memperbarui.
