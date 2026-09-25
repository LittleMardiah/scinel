# Contributing to ScIntel Backend

Terima kasih sudah mempertimbangkan untuk berkontribusi pada ScIntel Backend! 🎉

> **Catatan:** Panduan kontribusi lengkap (meliputi seluruh monorepo) tersedia di [`docs/CONTRIBUTING.md`](../docs/CONTRIBUTING.md). Dokumen ini berisi ringkasan khusus backend Go.

## Code of Conduct

Dengan berpartisipasi, Anda setuju: bersikap respek, profesional, dan konstruktif.

## Alur Kerja

1. Fork repository.
2. Buat branch: `git checkout -b feature/your-feature-name`
3. Buat perubahan.
4. Jalankan test: `make test`
5. Jalankan linter: `make lint`
6. Commit (lihat konvensi di bawah).
7. Push ke fork & buka Pull Request.

## Standar Development

### Go Code Style
- Jalankan `gofmt -s -w .` sebelum commit (termasuk di `make format`).
- Gunakan `golangci-lint` untuk menangkap masalah umum.
- Ikuti konvensi penamaan Go (CamelCase exported, camelCase unexported).
- Ikuti Clean Architecture: handler → service → repository. Jangan menaruh logika bisnis di handler.

### Commit Conventions
Kami mengikuti [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` — fitur baru
- `fix:` — perbaikan bug
- `docs:` — pembaruan dokumentasi
- `test:` — perubahan test
- `chore:` — tugas maintenance

**Contoh:** `feat(auth): implement JWT refresh token endpoint`

### Persyaratan Testing
- **Coverage**: kode baru harus mempertahankan >80% coverage.
- **Unit Test**: harus lolos lokal sebelum push.
- **Integration Test**: gunakan `sqlmock` (bukan testcontainers) agar berjalan di WSL tanpa Docker.
- **Jangan skip test** tanpa alasan; catat ke `docs/TECHNICAL_DEBT.txt` bila terpaksa.

## Pull Request Process

1. Deskripsi PR jelas: apa berubah, mengapa, dan cara verifikasi.
2. Update dokumentasi bila perlu.
3. Tunggu CI hijau. Semua check harus lulus.
4. Minta review (atau self-review untuk proyek solo).

## Melaporkan Bug

Gunakan GitHub Issues dengan label `bug`. Sertakan:
- Langkah reproduksi
- Perilaku yang diharapkan
- Perilaku aktual
- Environment (OS, versi Go, dll.)

Terima kasih sudah berkontribusi! 🚀
