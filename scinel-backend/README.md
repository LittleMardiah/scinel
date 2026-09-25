# ScIntel Backend

REST API untuk Supply Chain Intelligence Platform (ScIntel) — dibangun dengan **Go** dan **Clean Architecture**. Menyediakan autentikasi, manajemen inventori multi-gudang, demand forecasting (ETS/ARIMA), deteksi anomali, webhooks, dan API keys.

![Go](https://img.shields.io/badge/Go-1.25-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue)
![Redis](https://img.shields.io/badge/Redis-7-red)

---

## Fitur

- **Auth**: Register, login, refresh & logout dengan JWT (access 1 jam, refresh 7 hari)
- **RBAC**: Role hierarchy `admin > manager > analyst > viewer`
- **Inventory**: Query stok per gudang, transaksi masuk/keluar dengan pessimistic locking
- **Forecasting**: Prediksi 30 hari (ETS & ARIMA), confidence interval 95%, akurasi MAPE/RMSE/MAE
- **Anomaly Detection**: Z-score & IQR, severity classification, workflow ack/resolve
- **Webhooks**: Langganan event `inventory.changed`, `anomaly.detected` dengan retry & HMAC signing
- **API Keys**: Kunci ter-hash SHA-256, dukungan rotasi & revoke
- **Rate Limiting**: Sliding window Redis (1000 req/s), graceful degradation
- **Security**: Security headers (X-Frame-Options, HSTS), CORS, graceful shutdown, logging JSON (slog)

## Tech Stack

| Layer | Teknologi |
|-------|-----------|
| Bahasa | Go 1.25 (net/http standard router) |
| Database | PostgreSQL 15 (lib/pq), golang-migrate |
| Cache | Redis 7 (go-redis/v8) |
| Auth | JWT (golang-jwt/v5), bcrypt |
| Logging | slog (structured JSON) |
| Testing | Testify, go-sqlmock, miniredis |

## Struktur Proyek

```
scinel-backend/
├── cmd/api/               # Entry point & route registration (main.go)
├── internal/
│   ├── config/            # Load config dari environment
│   ├── database/          # Connection pool & transaction helpers
│   ├── handler/           # HTTP handlers (auth, inventory, forecast, dll)
│   ├── middleware/        # Auth, RBAC, CORS, SecurityHeaders, RateLimit
│   ├── model/             # Struct domain (User, Inventory, Forecast, dll)
│   ├── repository/        # Interface + implementasi PostgreSQL
│   ├── service/           # Business logic (Forecast, Anomaly, Webhook)
│   └── util/              # JWT, bcrypt hash, response wrapper
├── pkg/
│   ├── logger/            # Structured JSON logger (slog)
│   └── timeseries/        # ETS, ARIMA, anomaly detection
├── migrations/            # 12 skema tabel (up/down)
├── scripts/               # seed-db.go
├── tests/integration/     # Integration tests (sqlmock)
├── docker-compose.yml     # PostgreSQL 15 + Redis 7
├── Makefile               # Dev command terpusat
└── .env.example           # Template environment
```

## Quick Start

### Prasyarat
- Go 1.25+
- Docker & Docker Compose
- golang-migrate

### 1. Jalankan database & cache

```bash
docker-compose up -d
```

### 2. Setup environment

```bash
cp .env.example .env.local
# edit .env.local sesuai kebutuhan (DATABASE_URL, REDIS_URL, JWT_SECRET_KEY)
```

### 3. Jalankan migration & seed

```bash
make db-migrate     # terapkan 12 tabel
go run scripts/seed-db.go   # 3 gudang, 5 SKU, 15 inventory
```

### 4. Jalankan server

```bash
make run            # http://localhost:8080
```

### Verifikasi

```bash
curl http://localhost:8080/api/v1/health
```

## Command Makefile

```bash
make help           # daftar semua command
make run            # jalankan server
make test           # unit test
make test-coverage  # test + coverage
make lint           # golangci-lint
make build          # build binary
make db-up          # start PostgreSQL + Redis
make db-down        # stop
make db-reset       # reset & migrate ulang
```

## Testing

```bash
make test           # semua unit test
make test-coverage  # lihat persentase coverage
go test ./tests/integration/... -v   # integration test (sqlmock)
```

**Coverage hasil**: handler 83.7%, middleware 99.0%, config/model/util 90-100%.

## Endpoints

Base URL: `http://localhost:8080/api/v1`

| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| POST | `/auth/register` | Register user | Public |
| POST | `/auth/login` | Login & dapat token | Public |
| POST | `/auth/refresh` | Refresh access token | Public |
| POST | `/auth/logout` | Logout | JWT |
| GET | `/health` | Liveness check | Public |
| GET | `/readiness` | Deep check (DB + Redis) | Public |
| GET | `/liveness` | Deep check + DB stats | Public |
| GET | `/inventory` | Query inventory | JWT |
| POST | `/transactions` | Buat transaksi stok | JWT |
| GET | `/transactions` | Riwayat transaksi | JWT |
| GET | `/inventory/{warehouse_id}` | Stok per gudang | JWT |
| POST | `/inventory/import` | Bulk import | JWT |
| GET | `/forecasts` | Query forecast | JWT |
| POST | `/forecasts/generate` | Trigger forecast | Admin |
| GET | `/forecasts/accuracy` | Metrics akurasi | JWT |
| GET | `/anomalies` | List anomali | JWT |
| GET | `/anomalies/{id}` | Detail anomali | JWT |
| PUT | `/anomalies/{id}` | Update status (ack/resolve) | JWT |
| GET | `/anomalies/{id}/history` | Timeline status | JWT |
| POST | `/api-keys` | Create API key | Admin |
| GET | `/api-keys` | List API keys | Admin |
| DELETE | `/api-keys/{id}` | Revoke API key | Admin |
| POST | `/api-keys/{id}/rotate` | Rotasi secret | Admin |
| POST | `/webhooks/subscribe` | Subscribe event | JWT |
| GET | `/webhooks/subscriptions` | List subscription | JWT |
| DELETE | `/webhooks/subscriptions/{id}` | Unsubscribe | JWT |
| GET | `/webhooks/deliveries` | Riwayat delivery | JWT |

Dokumentasi detail: [docs/API.md](../docs/API.md)

## Middleware Chain

```
Request
  └─► CORS
       └─► SecurityHeaders
            └─► RateLimit (Redis, 1000 req/s)
                 └─► Auth (JWT)
                      └─► (opsional) RBAC "admin"
                           └─► Mux / Handler
```

## Teknikal Utama

- **Pessimistic locking** pada transaksi stok untuk mencegah race condition
- **Cache Redis TTL 30s** untuk query inventory, **24 jam** untuk forecast
- **Deadlock retry** (max 3 attempts, exponential backoff 100/200/400ms)
- **Dead-letter queue** pada webhook (max 5 attempts)
- **Graceful shutdown** 15 detik (SIGINT/SIGTERM)

## Lisensi

MIT
