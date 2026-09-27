# Changelog

Semua perubahan penting pada proyek ScIntel akan dicatat di file ini.

Format mengikuti [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
dan versioning mengikuti [Semantic Versioning](https://semver.org/).

## [v1.0.0] — 2026-09-05

Rilis stabil pertama ScIntel — Supply Chain Intelligence Platform.

### Added

#### Backend (Go / Clean Architecture)
- **Auth & RBAC**: register, login, refresh token, logout dengan JWT access/refresh; RBAC lintas role (admin, manager, analyst, viewer); audit trail.
- **Inventory**: manajemen stok multi-gudang, transaksi masuk/keluar dengan pessimistic locking, status stok normal/low/critical.
- **Forecasting**: prediksi permintaan 30 hari dengan ETS & ARIMA, confidence interval 95%, metrik akurasi MAPE/RMSE/MAE.
- **Anomaly Detection**: deteksi stok anomali via Z-score & IQR, klasifikasi severity (low/medium/high/critical), workflow acknowledge → resolve.
- **Webhooks**: notifikasi event real-time (`inventory.changed`, `anomaly.detected`) dengan retry exponential backoff, dead-letter queue, dan HMAC signing.
- **API Keys**: kunci API ter-hash (SHA-256) untuk integrasi programmatic, rotasi & revoke.
- **Rate Limiting**: sliding window berbasis Redis (1000 req/s per key), graceful degradation (fail-open) saat Redis turun.
- **Security Headers**: CORS whitelist, X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, HSTS, Content-Security-Policy, Referrer-Policy.
- **Database**: 12 tabel via golang-migrate + schema_id SQL; index untuk query umum (warehouse_id, created_at DESC, detected_at DESC, email).
- 23 endpoint REST client-callable (dokumentasi lengkap di `docs/API.md`).

#### Frontend (Next.js 14 / TypeScript)
- App Router + Tailwind CSS v3, dark mode via next-themes.
- Halaman: login/dashboard, inventory per-gudang, transaksi (form, history, filter, sorting, export), forecast & accuracy, anomaly management, settings & theme.
- Protected routes, navbar, sidebar; API client Axios dengan JWT interceptor.

#### Testing
- Unit tests backend: model, util, repository (sqlmock), service, handler (testify + sqlmock + miniredis), middleware (99.0% coverage).
- Integration tests (sqlmock, tanpa Docker/testcontainers).
- Frontend: 71 test suites via Jest + React Testing Library; E2E Playwright untuk 5 critical flow.

### Changed
- Refactor struktur repo menjadi monorepo `scinel-backend/` + `scinel-frontend/` + `docs/`.
- Dokumentasi lengkap: README, API docs (23 endpoint), ARCHITECTURE, TROUBLESHOOTING, CONTRIBUTING, LOGBOOK, TECHNICAL_DEBT.

### Fixed
- Integration test di WSL (pindah ke sqlmock, lepas dari Docker).
- BOM di `.env.local` backend (menyebabkan "command not found" saat `source`).
- Playwright baseURL salah port (3001 → 3000).

### Security
- Tidak ada secret yang dikomit; `.env.local` ter-ignore.
- Password di-hash dengan bcrypt, API key dengan SHA-256.
- HMAC signing untuk payload webhook.