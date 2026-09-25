# ScIntel Troubleshooting

Kumpulan masalah umum yang ditemukan selama development ScIntel beserta solusinya. Diambil dari pengalaman debugging nyata yang tercatat di [LOGBOOK](LOGBOOK.txt) dan [TECHNICAL_DEBT](TECHNICAL_DEBT.txt).

---

## Daftar Isi
1. [Backend tidak bisa konek ke PostgreSQL (28000)](#1-backend-tidak-bisa-konek-ke-postgresql-28000)
2. [Playwright baseURL salah port](#2-playwright-baseurl-salah-port)
3. [Error 42725 — current_date + parameter](#3-error-42725--current_date--parameter)
4. [Error 42P08 — parameter tipe tidak konsisten](#4-error-42p08--parameter-tipe-tidak-konsisten)
5. [Permissions TEXT[] tidak bisa di-scan ke []string](#5-permissions-text-tidak-bisa-di-scan-ke-string)
6. [Webhook payload body nil saat attemptDelivery](#6-webhook-payload-body-nil-saat-attemptdelivery)
7. [Test testcontainers gagal di WSL](#7-test-testcontainers-gagal-di-wsl)
8. [Rate limiter test di-skip (miniredis)](#8-rate-limiter-test-di-skip-miniredis)
9. [BOM pada .env.local](#9-bom-pada-envlocal)
10. [Placeholder ganda di query SQL ($$4)](#10-placeholder-ganda-di-query-sql-4)
11. [Frontend — 401 loop pada refresh token](#11-frontend--401-loop-pada-refresh-token)
12. [Error 409 transisi status anomaly](#12-error-409-transisi-status-anomaly)

---

## 1. Backend tidak bisa konek ke PostgreSQL (28000)

**Gejala:**
```
ERROR: role "scinel_user" does not exist (SQLSTATE 28000)
```

**Penyebab:** Role dan database belum di-provision di instance PostgreSQL lokal. `.env.local` memakai `scinel_user`/`scinel_db` namun role tsb belum dibuat.

**Solusi:**
```sql
-- jalankan sebagai superuser (postgres)
CREATE ROLE scinel_user LOGIN PASSWORD 'scinel_password';
CREATE DATABASE scinel_db OWNER scinel_user;
```

Kemudian pastikan `DATABASE_URL` di `.env.local` sesuai:
```
DATABASE_URL=postgres://scinel_user:scinel_password@localhost:5432/scinel_db?sslmode=disable
```

**Pencegahan:** Dokumentasikan provision role di README/setup guide.

---

## 2. Playwright baseURL salah port

**Gejala:** Test E2E gagal konek karena akses ke port 3001, padahal frontend Next.js berjalan di 3000.

**Penyebab:** `playwright.config.ts` memakai `baseURL: http://localhost:3001`.

**Solusi:**
```ts
// scinel-frontend/playwright.config.ts
use: {
  baseURL: 'http://localhost:3000',
}
```

---

## 3. Error 42725 — current_date + parameter

**Gejala:**
```
ERROR: operator is not unique: date + unknown (SQLSTATE 42725)
```

**Lokasi:** `forecast_repository_impl.go` pada `GetByWarehouseAndSKU`.

**Penyebab:** `CURRENT_DATE + $3` ambigu karena Go mengirim `$3` sebagai `int64` sementara PostgreSQL tidak tahu tipe mana yang dipakai.

**Solusi:** Cast eksplisit pada parameter:
```sql
WHERE forecast_date BETWEEN CURRENT_DATE AND CURRENT_DATE + $3::int
```

---

## 4. Error 42P08 — parameter tipe tidak konsisten

**Gejala:**
```
ERROR: could not determine data type of parameter $1 (SQLSTATE 42P08)
```

**Lokasi:** `anomaly_repository_impl.go` pada `UpdateStatus` saat memakai parameter yang sama 2x dengan tipe berbeda dalam `CASE WHEN`.

**Solusi:** Cast eksplisit tipe:
```sql
CASE WHEN $2 = 'resolved' THEN NOW() END,
status = $1::varchar
```

---

## 5. Permissions TEXT[] tidak bisa di-scan ke []string

**Gejala:**
```
unsupported Scan, storing driver.Value type []string into type *[]string
```

**Lokasi:** `api_key_repository_impl.go` (kolom `permissions TEXT[]`).

**Penyebab:** Driver `lib/pq` memerlukan `pq.Array(...)` untuk meng-scan/insert tipe array PostgreSQL.

**Solusi:** Bungkus semua read & write `permissions`:
```go
// membaca
var perms []string
rows.Scan(&id, &name, pq.Array(&perms), ...)

// menulis
db.Exec(`INSERT ... permissions = $3`, ..., pq.Array(item.Permissions))
```

---

## 6. Webhook payload body nil saat attemptDelivery

**Gejala:** Request webhook terkirim tanpa body payload (`nil`), meski event dipublish.

**Penyebab:** `attemptDelivery` tidak melakukan marshal payload ke JSON sebelum `http.NewRequest`.

**Solusi:** Marshal payload terlebih dahulu:
```go
body, err := json.Marshal(payload)
if err != nil { ... }
req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
```

---

## 7. Test testcontainers gagal di WSL

**Gejala:** Integration test memakai `testcontainers-go` (spawn Docker container) gagal di lingkungan WSL karena ketiadaan/keterbatasan Docker.

**Penyebab:** Dependency Docker eksternal tidak tersedia secara konsisten di WSL.

**Solusi:** Tulis ulang integration test memakai `sqlmock` — repository asli + DB mock tanpa container:
```go
// tests/integration/auth_flow_test.go
db, _ := sqlmock.New()
repo := repository.NewUserRepository(db)
// ... test alur register → login → refresh → logout
```
Hasil: test berjalan tanpa Docker. Dependency `testcontainers-go` masih ada di `go.mod` tapi tidak dipakai untuk auth flow.

---

## 8. Rate limiter test di-skip (miniredis)

**Gejala:** Test rate limiter di-skip (`t.Skip`) karena tidak ada Redis untuk unit test.

**Solusi:** Tambahkan dependency `github.com/alicebob/miniredis/v2` (in-memory Redis) dan tulis ulang test tanpa skip:
```go
mr, _ := miniredis.Run()
client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
```
Hasil: coverage middleware naik ke 99.0%.

---

## 9. BOM pada .env.local

**Gejala:** Ketika `source .env.local` muncul `command not found: ...` di baris pertama.

**Penyebab:** File `.env.local` diawali BOM (byte order mark).

**Solusi:** Bersihkan BOM:
```bash
sed -i '1s/^\xEF\xBB\xBF//' .env.local
```

**Pencegahan:** Simpan file dengan encoding UTF-8 tanpa BOM di editor.

---

## 10. Placeholder ganda di query SQL ($$4)

**Gejala:** Error bind parameter di query, mis. `error: unexpected $$4`.

**Lokasi:** `transaction_repository_impl.go` pada `Query`.

**Penyebab:** Copy-paste placeholder menghasilkan `$$4` (dolar ganda) alih-alih `$4`.

**Solusi:** Ganti `$$4` → `$4` untuk seluruh query.

---

## 11. Frontend — 401 loop pada refresh token

**Gejala:** Aplikasi logout atau infinite loop karena request 401 berulang.

**Penyebab:** Interceptor API belum membedakan `401` dari request `/auth/refresh` itu sendiri, sehingga refresh memicu refresh lagi.

**Solusi:** Di `lib/api.ts`, beri guard — hanya trigger refresh sekali; jika request itu sendiri adalah `/auth/refresh`, jangan refresh ulang, langsung logout:
```ts
if (isRetryRequest || isRefreshRequest) {
  // logout & redirect
} else {
  // call refresh, retry original request
}
```

---

## 12. Error 409 transisi status anomaly

**Gejala:** `PUT /anomalies/{id}` dengan status `resolved` langsung dari `open` gagal dengan 409.

**Penyebab:** Workflow status mengunci urutan: `open → acknowledged → resolved`. Melompati `acknowledged` ditolak.

**Solusi:** Ikuti urutan yang benar, atau gunakan endpoint history untuk update bertahap:
```json
// langkah 1
{ "status": "acknowledged" }
// langkah 2
{ "status": "resolved" }
```

---

## Masalah yang Belum Terselesaikan

Lihat [TECHNICAL_DEBT.txt](TECHNICAL_DEBT.txt) untuk:
- `[D-001]` Hardcoded rate limit (1000 req/s) → sebaiknya dari config.
- `[D-002]` Tidak ada per-route rate limit tiers untuk auth.
- `[D-003]` RBAC hanya mencakup rute API Keys.
- `[D-004]` Role "admin" hardcoded.

## Tips Debugging Cepat

```bash
# lihat log backend (JSON)
make run

# jalankan hanya satu test
go test ./internal/repository/... -run TestName -v

# verifikasi semua endpoint backend hidup
curl -s http://localhost:8080/api/v1/health
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"...","password":"..."}'

# cek koneksi db
docker-compose exec postgres pg_isready

# cek koneksi redis
docker-compose exec redis redis-cli ping
```
