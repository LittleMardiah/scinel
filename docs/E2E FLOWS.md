# E2E Testing — Manual Flows

> **Dibuat:** 2026-08-07  
> **Tujuan:** Memverifikasi 5 critical user flows berfungsi dengan benar di environment lokal/staging/production.  
> **Status:** ⏳ Belum dijalankan (template).

---

## Flow 1: Login → Dashboard → Logout

**Preconditions:**
- Backend & frontend running.
- Akun test tersedia (email: `test@example.com`, password: `password123`).

**Steps:**
1. Buka `http://localhost:3000/login`.
2. Masukkan email `test@example.com` dan password `password123`.
3. Klik tombol **Masuk**.
4. Tunggu redirect ke `/dashboard`.
5. Klik menu user (avatar/username) di pojok kanan atas.
6. Pilih **Keluar**.
7. Konfirmasi logout (jika ada modal).

**Expected Result:**
- ✅ Setelah login, user diarahkan ke `/dashboard`.
- ✅ Token JWT tersimpan di cookie (atau localStorage).
- ✅ Sidebar menampilkan menu sesuai role.
- ✅ Setelah logout, redirect ke `/login`.
- ✅ Token dihapus, sidebar tidak muncul.

**Verification Checklist:**
- [ ] URL setelah login adalah `/dashboard`.
- [ ] Ada elemen dashboard (KPI cards, chart, dll).
- [ ] Cookie `access_token` ada.
- [ ] Sidebar muncul dengan menu yang sesuai.
- [ ] Setelah logout, cookie `access_token` hilang.
- [ ] URL setelah logout adalah `/login`.

**Actual Result:** *(diisi setelah testing)*

---

## Flow 2: Create Transaction → See in History

**Preconditions:**
- User sudah login.
- Minimal 1 warehouse dan 1 SKU tersedia di database.

**Steps:**
1. Dari dashboard, klik menu **Inventory** → **Buat Transaksi** (atau langsung ke `/inventory/transaction/new`).
2. Pilih warehouse dari dropdown.
3. Pilih SKU dari dropdown (atau ketik untuk search).
4. Pilih tipe transaksi: **Masuk** atau **Keluar**.
5. Masukkan jumlah (quantity) > 0.
6. Isi alasan (reason) (contoh: "Testing E2E").
7. Klik tombol **Simpan Transaksi**.
8. Tunggu toast success muncul.
9. Navigasi ke **Inventory** → **Riwayat Transaksi** (atau `/inventory/transactions`).
10. Cari transaksi yang baru dibuat di daftar (bisa dengan filter atau scroll).

**Expected Result:**
- ✅ Transaksi berhasil dibuat (toast success).
- ✅ Transaksi muncul di daftar riwayat dengan data yang sesuai (SKU, warehouse, type, quantity, reason).
- ✅ Jumlah stok berubah sesuai transaksi (bisa dicek di inventory).

**Verification Checklist:**
- [ ] Toast success muncul: "Transaksi berhasil dibuat".
- [ ] Redirect ke `/inventory/transactions` setelah sukses (atau user bisa navigasi sendiri).
- [ ] Transaksi baru ada di tabel riwayat dengan data yang tepat.
- [ ] Kolom `previous_stock` dan `new_stock` terlihat (jika ada).
- [ ] Stok di inventory berubah sesuai delta.

**Actual Result:** *(diisi setelah testing)*

---

## Flow 3: View Forecast → Check Accuracy

**Preconditions:**
- User sudah login.
- Ada data historis transaksi (minimal 30 hari) untuk SKU tertentu.
- Forecast sudah pernah di-generate (atau generate manual via admin).

**Steps:**
1. Dari dashboard, klik menu **Forecast** (atau `/forecast`).
2. Pilih warehouse dan SKU yang memiliki forecast.
3. Pastikan chart menampilkan prediksi permintaan 30 hari ke depan dengan confidence interval (area biru transparan).
4. Klik tab atau link **Akurasi** (atau `/forecast/accuracy`).
5. Lihat metrics: MAPE (%), RMSE, MAE.
6. Lihat perbandingan algoritma ETS vs ARIMA (bar chart atau tabel).

**Expected Result:**
- ✅ Chart forecast tampil dengan benar (line + shaded area).
- ✅ Tooltip menampilkan nilai prediksi, lower, upper bound.
- ✅ Accuracy metrics ditampilkan (MAPE, RMSE, MAE).
- ✅ Ada perbandingan algoritma (ETS vs ARIMA) dengan penanda algoritma terbaik.

**Verification Checklist:**
- [ ] Chart muncul tanpa error (Recharts render).
- [ ] Confidence interval (area) terlihat.
- [ ] MAPE < 15% (idealnya) atau sesuai data.
- [ ] Ada label "Best Algorithm" atau perbandingan yang jelas.
- [ ] Filter warehouse dan SKU berfungsi.

**Actual Result:** *(diisi setelah testing)*

---

## Flow 4: Acknowledge Anomaly → Update Status

**Preconditions:**
- User sudah login dengan role yang memiliki akses ke anomalies (admin/manager).
- Ada minimal 1 anomaly dengan status `open` di database (bisa dibuat manual via seeding atau dari transaksi yang memicu anomaly).

**Steps:**
1. Dari dashboard, klik menu **Anomalies** (atau `/anomalies`).
2. Lihat daftar anomaly, pilih salah satu yang statusnya `open`.
3. Klik baris anomaly untuk melihat detail (atau klik tombol **Detail**).
4. Di halaman detail, periksa informasi: SKU, severity, z-score, description, timeline.
5. Klik tombol **Akui** atau **Acknowledge**.
6. (Opsional) Tambahkan catatan di modal.
7. Konfirmasi acknowledge.
8. Lihat status berubah menjadi `acknowledged`.
9. Klik tombol **Selesai** atau **Resolve**.
10. Tambahkan catatan resolusi (contoh: "False alarm, checked manually").
11. Konfirmasi resolve.
12. Status berubah menjadi `resolved`.

**Expected Result:**
- ✅ Status anomaly berubah dari `open` → `acknowledged` → `resolved`.
- ✅ Timeline menampilkan perubahan status dengan timestamp dan user.
- ✅ Catatan resolusi tersimpan dan terlihat di detail.
- ✅ Tidak ada error saat update status.

**Verification Checklist:**
- [ ] Detail anomaly menampilkan semua informasi (SKU, severity, z-score, description).
- [ ] Tombol Acknowledge dan Resolve tersedia sesuai status.
- [ ] Setelah acknowledge, status badge berubah.
- [ ] Setelah resolve, status badge berubah dan resolved_at terisi.
- [ ] Timeline menampilkan semua event dengan urutan yang benar.

**Actual Result:** *(diisi setelah testing)*

---

## Flow 5: Toggle Dark Mode → Persists

**Preconditions:**
- User sudah login (atau bahkan di halaman login pun bisa).

**Steps:**
1. Buka aplikasi (misal di `/dashboard`).
2. Perhatikan tema saat ini (light/dark).
3. Klik tombol toggle tema di navbar (ikon matahari/bulan).
4. Amati perubahan tema (background, teks, card, button).
5. Refresh halaman.
6. Periksa apakah tema tetap seperti setelah toggle.

**Expected Result:**
- ✅ Tema berubah secara instan (light ↔ dark) saat toggle.
- ✅ Semua elemen UI mengikuti tema (background, teks, border, shadow).
- ✅ Setelah refresh, tema tetap sama (persistence via localStorage).
- ✅ Tema toggle bekerja di semua halaman.

**Verification Checklist:**
- [ ] Toggle mengubah class `dark` di elemen `<html>`.
- [ ] localStorage key `theme` tersimpan dengan nilai `light` atau `dark`.
- [ ] Warna background berubah (misal dari white ke gray-900).
- [ ] Teks menjadi putih di dark mode.
- [ ] Card, button, input menyesuaikan warna.
- [ ] Setelah refresh, tema tetap sesuai pilihan.

**Actual Result:** *(diisi setelah testing)*

---

## 📋 EKSEKUSI TESTING

**Cara menjalankan:**

1. Pastikan backend & frontend running:
   ```bash
   # Terminal 1: Backend
   cd backend && go run ./cmd/server

   # Terminal 2: Frontend
   cd frontend && pnpm dev