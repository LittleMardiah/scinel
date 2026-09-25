# ScIntel Architecture

Dokumen ini menjelaskan arsitektur teknis Supply Chain Intelligence Platform (ScIntel) — Clean Architecture pada backend Go, struktur frontend Next.js, alur data, dan skema database.

---

## 1. Gambaran Umum

ScIntel adalah platform web untuk manajemen rantai pasok dengan tiga kemampuan inti: **demand forecasting**, **deteksi anomali stok**, dan **distribusi stok multi-gudang**.

```
┌──────────────────────────────┐
│  PostgreSQL 15               │
│  (sumber kebenaran/data)     │
└──────────────┬───────────────┘
               │ lib/pq (connection pool)
┌──────────────▼───────────────┐      ┌─────────────────┐
│  scinel-backend (Go 1.25)    │◄────►│  Redis 7        │
│  Clean Architecture          │ cache│  (forecast TTL, │
│                              │ rate │   rate limiter) │
└──────────────┬───────────────┘ limit└─────────────────┘
               │ REST/JSON (8080)
┌──────────────▼───────────────┐
│  scinel-frontend (Next.js)   │
│  App Router / port 3000      │
└──────────────────────────────┘
```

---

## 2. Clean Architecture (Backend)

Backend mengikuti **Clean Architecture** — pemisahan yang tegas antar layer dengan dependency injection dari luar ke dalam.

```
┌─────────────────────────────────────────────────────────┐
│ LAYER PRESENTATION                                      │
│  internal/handler/  (HTTP request/response, validasi)   │
│  internal/middleware/ (Auth, RBAC, CORS, RateLimit)     │
├─────────────────────────────────────────────────────────┤
│ LAYER APPLICATION / SERVICE                             │
│  internal/service/  (business logic: forecast, anomaly, │
│                      webhook)                           │
├─────────────────────────────────────────────────────────┤
│ LAYER DOMAIN / REPOSITORY                               │
│  internal/repository/ (interface → implementasi SQL)    │
│  internal/model/      (struct domain + enum)            │
├─────────────────────────────────────────────────────────┤
│ LAYER INFRASTRUCTURE                                    │
│  internal/database/   (connection pool, transaction)    │
│  internal/config/     (environment config)              │
│  internal/util/       (JWT, bcrypt, response wrapper)   │
│  pkg/                 (timeseries, logger)              │
└─────────────────────────────────────────────────────────┘
```

### 2.1 Alur Dependency

```go
// cmd/api/main.go — jalinan utama (composition root)
userRepo   := repository.NewUserRepository(db)
forecastSvc := service.NewForecastService(forecastRepo, invRepo, txnRepo, appLogger)
authHandler := handler.NewAuthHandler(userRepo, jwtManager, cfg)
```

- **Handler** menerima interface *repository/service* (bukan implementasi konkret) → mudah di-test dengan mock.
- **Service** menerima beberapa repository → mengorkestrasi business logic lintas domain.
- **main.go** adalah satu-satunya tempat wiring (composition root). Tidak ada `new` di dalam handler.

### 2.2 Middleware Chain

```
Request
  └─► CORSMiddleware            # izinkan origin terkonfigurasi (frontend)
       └─► SecurityHeaders      # X-Frame-Options, HSTS, dll
            └─► RateLimit       # sliding window Redis, 1000 req/s
                 └─► AuthMiddleware   # validasi & parse JWT → UserInfo di context
                      └─► (ops.) RBACMiddleware("admin")  # pembatasan role
                           └─► http.ServeMux
```

**RBAC hanya di rute API Keys** (admin-only). Rute lain hanya butuh autentikasi.

---

## 3. Struktur Folder Backend

```
scinel-backend/
├── cmd/api/main.go          # Entry point, wiring, route registration
├── internal/
│   ├── config/config.go     # LoadConfig dari env (type-safe, default values)
│   ├── database/
│   │   ├── connection.go    # Connection pool (min 10, max 50)
│   │   └── transaction.go   # WithTx, deadlock retry (3x, backoff 100/200/400ms)
│   ├── handler/             # auth, health, inventory, forecast, anomaly, api_key, webhook
│   ├── middleware/          # auth.go, rbac.go, cors.go, rate_limit.go, security.go
│   ├── model/               # user, inventory, transaction, forecast, anomaly, webhook, api_key
│   ├── repository/          # *_repository.go (interface) + *_repository_impl.go (SQL)
│   ├── service/             # forecast_service, anomaly_service, webhook_service
│   └── util/                # jwt.go, hash.go, crypto.go, response.go
├── pkg/
│   ├── logger/logger.go     # slog JSON wrapper
│   └── timeseries/          # ets.go, arima.go, anomaly.go
├── migrations/              # 000001-000012 (up/down)
├── scripts/seed-db.go       # 3 gudang, 5 SKU, 15 inventory
└── tests/integration/       # auth_flow_test.go (sqlmock)
```

---

## 4. Alur Data — Request Lifecycle

### Contoh: Create Transaction

```
1. POST /api/v1/transactions
   HTTP request dari Frontend / eksternal

2. Middleware Chain
   CORS → SecurityHeaders → RateLimit → Auth
   (validasi token JWT, set UserInfo ke context)

3. inventoryHandler.CreateTransaction
   - Decode & validasi body (warehouse_id, sku_id/type/quantity/reason)
   - Parsing UUID

4. db.WithTx (transaction helper)
   - BEGIN
   - SELECT ... FOR UPDATE  ← pessimistic lock pada baris inventory
   - cek ketersediaan stok (utk type "out")
   - UPDATE inventory SET quantity = quantity ± n
   - INSERT transactions
   - invalidation cache Redis (inventory TTL 30s)
   - COMMIT (atau ROLLBACK on error)
   - retry otomatis jika deadlock (max 3x)

5. response.WriteResponse (format standar)
```

### Contoh: Generate Forecast (async)

```
1. POST /forecasts/generate (admin)
2. inventoryHandler → forecastSvc.GenerateForecastBatch
3. Ambil semua warehouse×sku via invRepo.GetAll
4. Untuk tiap kombinasi:
   - Ambil 180 hari history "out" via txnRepo.Query
   - aggregateDailyDemand → deret harian
   - Jalankan ETS (eksponensial) & ARIMA (auto p,d,q)
   - Auto-pilih berdasarkan MAPE terendah
   - Simpan 30 hari forecast + 1 record forecast_accuracy
5. Cache Redis di-invalidasi (TTL 24 jam)
```

---

## 5. Skema Database

12 tabel dibuat via golang-migrate (000001-000012).

### Relasi Utama

```
users ──────────────┬── webhook_subscriptions ── webhook_deliveries
                    └── api_keys
                    └── password_reset_tokens

warehouses ──< inventory >── skus
                └────────── transactions

skus/warehouses ──< forecasts ──< forecast_accuracy
                  └─ anomalies

inventory ──< inventory_snapshots
users ──────< audit_log
```

### Tabel & Kolom Kunci

| Tabel | Kolom penting | Constraint |
|-------|---------------|------------|
| `users` | email (uniq), password_hash, full_name, **role** | role: admin/manager/analyst/viewer; status: active/suspended/pending_email_verification |
| `warehouses` | name, code, location | — |
| `skus` | sku_code (uniq), name, category, unit | — |
| `inventory` | warehouse_id, sku_id, quantity, reorder_point | unique(warehouse_id, sku_id); partial index utk stok rendah |
| `transactions` | warehouse_id, sku_id, type, quantity, reason | type: in/out/adjustment/transfer |
| `forecasts` | warehouse_id, sku_id, forecast_date, predicted_demand, lower/upper_bound, **algorithm** | algorithm: exponential_smoothing/arima |
| `forecast_accuracy` | mape, rmse, mae, algorithm | — |
| `anomalies` | type, severity, z_score, status, **resolved_at** | severity: low/medium/high/critical; status: open/acknowledged/resolved |
| `webhook_subscriptions` | user_id, event_type, webhook_url, secret_key | event: inventory.changed/forecast.generated/anomaly.detected |
| `webhook_deliveries` | subscription_id, payload (jsonb), status, attempt_count | status: pending/success/failed/dead_letter |
| `audit_log` | user_id, action, entity_type, entity_id, changes | — |
| `inventory_snapshots` | inventory_id, quantity | snapshot untuk historis |
| `api_keys` | user_id, key_hash, name, permissions (text[]), expires_at | key di-hash SHA-256 |
| `password_reset_tokens` | user_id, token_hash, expires_at | — |

---

## 6. Strategi Caching (Redis)

| Data | TTL | Kunci |
|------|-----|-------|
| Query inventory | 30 detik | `cache:inventory:{warehouse}:{sku}:{status}` |
| Forecast | 24 jam | `cache:forecast:{warehouse}:{sku}:{days}` |
| Rate limit | 1 detik (sliding window) | `rl:{bucket}:{key}` |

- **Invalidation**: saat transaksi inventory dibuat, cache inventory untuk kombinasi terkait di-invalidate.
- **Graceful degradation**: jika Redis down, request tetap dilanjutkan (fail-open) tanpa cache — hanya rate limit & cache yang kehilangan benefit.

---

## 7. Struktur Frontend (Next.js)

```
app/
├── (auth)/                  # Route group publik
│   ├── login/  register/  forgot-password/  reset-password/
├── (protected)/             # Route group dengan auth check
│   ├── dashboard/           # KPI cards, chart, polling 30s
│   ├── inventory/           # list, [warehouse], transaction/new, transactions
│   ├── forecast/            # page, accuracy
│   ├── anomalies/           # page, [id], analytics
│   ├── profile/  settings/
├── layout.tsx               # Root (ThemeProvider, AuthProvider)
└── globals.css
```

**Data fetching**: `hooks/useApi.ts` membungkus `lib/api.ts` (Axios interceptor):
1. Tambahkan `Authorization: Bearer <token>` ke request.
2. Bila respons `401` → otomatis panggil `/auth/refresh`, lalu retry request.
3. Bila refresh gagal → logout & redirect ke `/login`.

**Polling real-time**: dashboard, inventory, dan anomaly menggunakan `setInterval` 30 detik via `useDashboardData` dan hook serupa.

---

## 8. Keamanan

- **Password**: bcrypt cost 12 (~100ms per hash).
- **JWT**: HMAC-SHA256, access 1 jam, refresh 7 hari.
- **API Keys**: secret di-generate, prefix + full secret ditampilkan sekali; tersimpan sebagai hash SHA-256.
- **Webhooks**: payload ditandatangani HMAC-SHA256 dengan secret per subscription.
- **Security headers**: X-Frame-Options, X-Content-Type-Options, HSTS, Referrer-Policy.
- **Rate limiting**: 1000 req/s per key (API key > JWT user > IP).

---

## 9. Algoritma Time Series

### ETS (Exponential Smoothing)
- Triplet exponential smoothing (level `alpha`, trend `beta`, seasonality `gamma`).
- Digunakan untuk prediksi 30 hari dari 180 hari historical.
- Output: `predicted_demand`, `lower_bound`, `upper_bound` (interval 95%).

### ARIMA
- `AutoRegressive Integrated Moving Average` dengan **grid search** untuk memilih `(p, d, q)` terbaik.
- Memilih parameter via AIC/kriteria.

### Anomaly Detection
- **Z-score**: nilai di luar `mean ± 2.0σ` → anomali.
- **IQR**: nilai di luar `[Q1 - 1.5·IQR, Q3 + 1.5·IQR]` → anomali.
- **Severity**: dipetakan dari besarnya deviasi (low/medium/high/critical).

---

## 10. Deployment & Operasional

- **docker-compose.yml** menyediakan PostgreSQL 15 + Redis 7 untuk development.
- Health probes: `/health` (liveness), `/readiness` (readiness — cek DB & Redis), `/liveness` (deep + DB stats) — siap untuk Kubernetes.
- `Makefile` mengedepankan command: `make run`, `make test`, `make db-migrate`, `make db-reset`, `make docker-up`, dll.
- Graceful shutdown 15 detik (SIGINT/SIGTERM) agar request in-flight selesai.

Lihat [ROADMAP 4. DEPLOYMENT.md](ROADMAP%204.%20DEPLOYMENT.md) untuk detail deployment.
