# ScIntel — Supply Chain Intelligence Platform

Platform manajemen supply chain berbasis web yang menyediakan **prediksi permintaan (demand forecasting)**, **deteksi anomali stok**, dan **distribusi stok multi-gudang** secara real-time — dilengkapi API yang aman untuk integrasi sistem eksternal.

![Status](https://img.shields.io/badge/status-active-brightgreen)
![Go](https://img.shields.io/badge/Go-1.25-blue)
![Next.js](https://img.shields.io/badge/Next.js-14.2.5-black)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue)
![Redis](https://img.shields.io/badge/Redis-7-red)

---

## Fitur Utama

| Modul | Deskripsi |
|-------|-----------|
| **Auth & RBAC** | Register, login, JWT access/refresh token, role-based access control (admin, manager, analyst, viewer) |
| **Inventory** | Manajemen stok multi-gudang, status stok normal/low/critical, transaksi stok masuk/keluar dengan pessimistic locking |
| **Forecasting** | Prediksi permintaan 30 hari ke depan menggunakan algoritma ETS & ARIMA, confidence interval 95%, metrik akurasi MAPE/RMSE/MAE |
| **Anomaly Detection** | Deteksi anomali stok via Z-score & IQR, klasifikasi severity (low/medium/high/critical), workflow acknowledge → resolve |
| **Webhooks** | Langganan event real-time (inventory.changed, anomaly.detected), delivery dengan retry exponential backoff & HMAC signing |
| **API Keys** | Kunci API ter-hash untuk integrasi programmatic, rotasi & revoke |
| **Rate Limiting** | Sliding window berbasis Redis (1000 req/s), graceful degradation jika Redis turun |

## Arsitektur

- **Backend** : Go 1.25 — Clean Architecture (handler → service → repository)
- **Frontend** : Next.js 14 (App Router) + TypeScript + Tailwind CSS
- **Database** : PostgreSQL 15 (12 tabel, golang-migrate)
- **Cache/Optimization** : Redis 7 (caching forecast, rate limiter)

```
┌──────────────────────┐      ┌──────────────────────────┐
│  scinel-frontend     │ ───► │  scinel-backend (Go)     │
│  Next.js 14 / TS     │ HTTP │  Clean Architecture      │
└──────────────────────┘      └────────────┬─────────────┘
                                           │
                          ┌────────────────┴─────────────┐
                          ▼                               ▼
                  ┌──────────────┐               ┌──────────────┐
                  │  PostgreSQL  │               │    Redis     │
                  │  15 (12 tbl) │               │ 7 (cache/RL) │
                  └──────────────┘               └──────────────┘
```

## Quick Start

### Prasyarat
- Go 1.25+
- Node.js 20+ & pnpm
- Docker & Docker Compose
- golang-migrate

### 1. Jalankan database & cache

```bash
cd scinel-backend
docker-compose up -d
```

### 2. Setup backend

```bash
cd scinel-backend
cp .env.example .env.local      # lalu sesuaikan value

# jalankan migration
make db-migrate

# seed data dummy (3 gudang, 5 SKU, 15 inventory)
go run scripts/seed-db.go

# jalankan server (port 8080)
make run
```

### 3. Setup frontend

```bash
cd scinel-frontend
cp .env.example .env.local      # NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
pnpm install
pnpm dev                        # http://localhost:3000
```

### 4. Verifikasi

```bash
# backend health check
curl http://localhost:8080/api/v1/health

# frontend test
cd scinel-frontend && pnpm test
```

## Performance Metrics

| Metrik | Target | Verifikasi |
|--------|--------|------------|
| Test coverage (unit) | >80% | velocity, model, util, handler, middleware |
| Handler coverage | 83.7% | `go test ./internal/handler/... -cover` |
| Middleware coverage | 99.0% | rate limiter, auth, RBAC |
| Test suites (frontend) | 69+ PASS / 71 | `pnpm test` |
| Failure latency | fail-open | Redis down → request tetap lanjut |
| Rate limit | 1000 req/s per key | sliding window Redis |

## Dokumentasi

| Dokumen | Deskripsi |
|---------|-----------|
| [API Docs](docs/API.md) | Dokumentasi lengkap 27 endpoint (request/reponse) |
| [Architecture](docs/ARCHITECTURE.md) | Clean architecture, data flow, skema database |
| [Troubleshooting](docs/TROUBLESHOOTING.md) | Masalah umum & solusi dari pengalaman debugging |
| [Contributing](docs/CONTRIBUTING.md) | Code standards & PR process |
| [Logbook](docs/LOGBOOK.txt) | Catatan perkembangan development |
| [Technical Debt](docs/TECHNICAL_DEBT.txt) | Register utang teknis (open/resolved) |

## Struktur Repository

```
.
├── scinel-backend/          # Go REST API (Clean Architecture)
│   ├── cmd/api/             # Entry point & route registration
│   ├── internal/            # handler, service, repository, middleware, model
│   ├── pkg/                 # timeseries (ETS/ARIMA/anomaly), logger
│   ├── migrations/          # 12 skema tabel (up/down)
│   ├── scripts/             # seed data
│   └── tests/integration/   # integration test (sqlmock)
├── scinel-frontend/         # Next.js web app
│   ├── app/                 # App Router pages
│   ├── components/          # 40+ reusable components
│   ├── context/             # Auth, Theme
│   ├── hooks/               # useApi, useDashboardData
│   └── lib/                 # api client, utils
└── docs/                    # Dokumentasi proyek
```

## Tech Stack

**Backend**
- Go 1.25, net/http (standard library router)
- PostgreSQL 15 + lib/pq + golang-migrate
- Redis 7 (go-redis/v8)
- JWT (golang-jwt/v5), bcrypt
- slog (structured JSON logging)
- Testify, sqlmock, miniredis

**Frontend**
- Next.js 14.2.5 (App Router), TypeScript
- Tailwind CSS v3 (dark mode via next-themes)
- Recharts (visualisasi forecast/anomaly)
- Axios (API client + JWT interceptor)
- Jest + React Testing Library, Playwright (E2E)

## Lisensi

MIT — lihat [LICENSE](LICENSE) untuk detail.
