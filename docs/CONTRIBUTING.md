# Contributing to ScIntel

Terima kasih sudah ingin berkontribusi pada Supply Chain Intelligence Platform (ScIntel)! 🎉

Dokumen ini berisi standar kode dan proses pull request yang berlaku untuk **seluruh monorepo** (backend Go + frontend Next.js).

---

## Code of Conduct

Dengan berpartisipasi, Anda setuju untuk:
- Bersikap **respek**, **profesional**, dan **konstruktif**.
- Fokus pada diskusi teknis, bukan personal.
- Menerima kritik dengan terbuka dan memberikan umpan balik yang membangun.

---

## Alur Kerja

1. **Fork** repository (atau buat branch untuk member).
2. Buat branch baru:
   ```bash
   git checkout -b feature/nama-fitur
   # atau
   git checkout -b fix/nama-bug
   ```
3. Buat perubahan Anda.
4. Jalankan pengujian & linting (lihat bagian per-subproyek).
5. Commit dengan [Conventional Commits](#commit-conventions).
6. Push ke branch & buka Pull Request.
7. Tunggu CI hijau & review.

---

## Commit Conventions

Kami mengikuti [Conventional Commits](https://www.conventionalcommits.org/):

| Prefix | Deskripsi |
|--------|-----------|
| `feat:` | Fitur baru |
| `fix:` | Perbaikan bug |
| `docs:` | Perubahan dokumentasi |
| `test:` | Perubahan test |
| `chore:` | Tugas maintenance (deps, config) |
| `refactor:` | Refactoring tanpa ubah perilaku |

**Format:**
```
<type>(<scope>): <deskripsi>

contoh:
feat(auth): implement JWT refresh token endpoint
fix(anomaly): resolve 42P08 cast pada UpdateStatus
docs(api): tambahkan endpoint webhook delivery
```

Gunakan imperative mood dalam `type(scope):`. Branch untuk satu concern hanya boleh satu commit scope.

---

## Backend (Go) — Standar

### Kode
- Jalankan `gofmt -s -w .` sebelum commit (termasuk di `make format`).
- Gunakan `golangci-lint` via `make lint`.
- Ikuti konvensi penamaan Go (CamelCase exported, camelCase unexported).
- **Jangan mengubah interface repository/service** bila hanya menambah method — tambah method baru tanpa mengubah yang sudah ada (hindari breaking change).
- Struktur mengikuti Clean Architecture:
  - `handler/` hanya urus HTTP (decode/validate/respond), tidak logika bisnis.
  - `service/` mengorkestrasi business logic.
  - `repository/` berisi SQL, dipanggil via interface.

### Testing
- **Semua kode baru harus punya test**.
- Target coverage package kritis: **>80%** (config, model, util, handler, middleware).
- Test berjalan tanpa Docker di WSL — gunakan `sqlmock` & `miniredis` untuk unit/integration test.
- Jangan *skip* test tanpa alasan; catat ke [TECHNICAL_DEBT](../docs/TECHNICAL_DEBT.txt) bila terpaksa.

```bash
cd scinel-backend
make test                 # unit test semua
make test-coverage        # + coverage
go test ./...             # semua paket
go vet ./...              # static analysis
```

---

## Frontend (Next.js) — Standar

### Kode
- TypeScript, tanpa `any` (kecuali justified).
- Gunakan komponen existing, jangan duplikasi.
- Ikuti pola folder yang ada (`components/<modul>/`, `hooks/`, `lib/`).
- Tailwind utility untuk styling; dukung **dark mode** (via `next-themes`).
- Fokus aksesibilitas: ARIA labels, touch target min 44px, kontras warna.

### Testing
- Setiap komponen/feature baru tambahkan test (Jest + React Testing Library).
- Test harus deterministik (mock global: `next/navigation`, `next-themes`, `recharts`).
- E2E: gunakan Playwright untuk flow kritis (login, transaksi, forecast, anomaly, theme).

```bash
cd scinel-frontend
pnpm test                 # unit test (Jest)
pnpm test:coverage        # + coverage
pnpm test:ci              # CI variant (Vitest)
pnpm lint                 # ESLint
pnpm test:e2e             # Playwright (butuh backend + frontend hidup)
```

---

## Pull Request Process

1. Judul PR jelas & deskriptif, referensi issue bila ada (mis. `Closes #12`).
2. Deskripsi PR: **apa** berubah, **mengapa**, dan **cara verifikasi**.
3. Lampirkan screenshot/gif untuk perubahan UI bila relevan.
4. Update dokumentasi bila diperlukan (README, API.md, ARCHITECTURE.md).
5. Pastikan **semua CI check hijau** (lint, build, test, coverage).
6. Catat perubahan signifikan di `docs/LOGBOOK.txt`.
7. Untuk solo project: lakukan self-review sebelum merge.

### Checklist PR
- [ ] `go build ./...` / `pnpm build` lolos
- [ ] Unit test lolos & coverage tidak turun
- [ ] Lint tidak ada error
- [ ] Dokumentasi diperbarui bila perlu
- [ ] Tidak ada secret dikomit (cek `.env.local`)

---

## Menambahkan Endpoint Baru

1. Buat method di `service/` (bila ada business logic) atau langsung di `handler/`.
2. Tambahkan method interface di `repository/` bila butuh data baru.
3. Implementasi query di `*_repository_impl.go` (pakai `$n` placeholder, hindari `$$n`).
4. Daftarkan route di `cmd/api/main.go`.
5. Tambahkan test handler (sqlmock).
6. Update [API.md](../docs/API.md) — request/response contoh.

---

## Menambahkan Migrasi Database

1. Buat file `migrations/0000NN_name.up.sql` (dan `.down.sql`).
2. Nomor increment dari yang terakhir; jangan mengubah file migration yang sudah di-commit.
3. Jalankan `make db-migrate` & verifikasi `make db-reset` clean.
4. Update skema di [ARCHITECTURE.md](../docs/ARCHITECTURE.md).

---

## Melaporkan Bug

Gunakan GitHub Issues dengan label `bug`. Sertakan:
- Langkah reproduksi.
- Perilaku yang diharapkan vs aktual.
- Environment (OS, Go/Node version, log error).
- Screenshot & request_id (untuk masalah API).

---

## Referensi

- [README](../README.md) — overview & quick start
- [API.md](../docs/API.md) — dokumentasi endpoint
- [ARCHITECTURE.md](../docs/ARCHITECTURE.md) — arsitektur & skema DB
- [TROUBLESHOOTING.md](../docs/TROUBLESHOOTING.md) — masalah umum
- [TECHNICAL_DEBT.txt](../docs/TECHNICAL_DEBT.txt) — utang teknis yang diketahui

Terima kasih sudah berkontribusi! 🚀
