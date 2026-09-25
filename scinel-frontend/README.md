# ScIntel Frontend

Web application untuk Supply Chain Intelligence Platform (ScIntel) — dibangun dengan **Next.js 14** (App Router), **TypeScript**, dan **Tailwind CSS**. Menyediakan dashboard real-time untuk manajemen inventori multi-gudang, demand forecasting, deteksi anomali, dan pengaturan.

![Next.js](https://img.shields.io/badge/Next.js-14.2.5-black)
![TypeScript](https://img.shields.io/badge/TypeScript-5-blue)
![Tailwind](https://img.shields.io/badge/Tailwind-3.4.1-38bdf8)

---

## Fitur

- **Auth Flow**: Login, register, forgot/reset password (Bahasa Indonesia UI)
- **Dashboard**: 4 KPI cards, chart forecast (Recharts), polling real-time 30s
- **Inventory**: Tabel stok multi-gudang, sorting/filter/pagination, transaksi baru, riwayat transaksi + export CSV
- **Forecast**: Visualisasi demand + confidence interval, perbandingan akurasi ETS vs ARIMA
- **Anomaly**: List/badge severity, detail + timeline, workflow ack/resolve, analitik (trend, distribusi, top SKU)
- **Settings**: Profil, change password, API keys manager, webhook manager, tema
- **Dark Mode**: Light/dark/system (next-themes) — persists di localStorage
- **Responsive**: Touch-friendly (min 44px), mobile navigation, bottom sheet, table scroll

## Tech Stack

| Layer | Teknologi |
|-------|-----------|
| Framework | Next.js 14.2.5 (App Router) |
| Bahasa | TypeScript 5 |
| Styling | Tailwind CSS v3, clsx, tailwind-merge |
| Data Fetching | Axios (JWT interceptor + refresh) |
| Chart | Recharts 3 |
| State | React Context (Auth, Theme) |
| Theme | next-themes |
| Testing | Jest + React Testing Library, Playwright (E2E) |
| Package Manager | pnpm |

## Struktur Proyek

```
scinel-frontend/
├── app/
│   ├── (auth)/               # Route group: login, register, forgot/reset-password
│   │   ├── login/
│   │   ├── register/
│   │   ├── forgot-password/
│   │   └── reset-password/
│   ├── (protected)/          # Route group (auth check): dashboard + modul
│   │   ├── dashboard/
│   │   ├── inventory/        # list, [warehouse], transaction/new, transactions
│   │   ├── forecast/         # page, accuracy
│   │   ├── anomalies/        # page, [id], analytics
│   │   ├── profile/
│   │   └── settings/
│   ├── layout.tsx            # Root layout (ThemeProvider, AuthProvider)
│   └── globals.css
├── components/
│   ├── anomaly/              # AnomalyList, SeverityBadge, Timeline, dll
│   ├── auth/                 # LoginForm, RegisterForm, PasswordResetForm
│   ├── common/               # Button, Card, Input, Select, Modal, dll
│   ├── dashboard/            # StatsCard, KPIGrid, Charts
│   ├── forecast/             # ForecastChart, AccuracyMetrics
│   ├── inventory/            # InventoryTable, TransactionForm, History, dll
│   ├── layout/               # Navbar, Sidebar, Footer, MobileNav
│   ├── settings/             # ApiKeyManager, WebhookManager, ThemeSelector
│   └── ui/                   # Badge, Checkbox, Radio, TextArea, DatePicker
├── context/                  # AuthContext, ThemeContext
├── hooks/                    # useApi, useDashboardData, useBreakpoint
├── lib/                      # api.ts (client), utils.ts
└── tests/                    # unit (jest) + e2e/ (playwright)
```

## Quick Start

### Prasyarat
- Node.js 20+
- pnpm
- Backend berjalan di `http://localhost:8080` (lihat [scinel-backend/README.md](../scinel-backend/README.md))

### 1. Install dependencies

```bash
pnpm install
```

### 2. Setup environment

```bash
cp .env.example .env.local
# NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### 3. Jalankan development server

```bash
pnpm dev       # http://localhost:3000
```

### 4. Login

Buka `http://localhost:3000/login`, gunakan akun yang didaftarkan via backend (atau seed user).

## Testing

```bash
pnpm test                   # semua unit test (Jest)
pnpm test:coverage          # test + coverage
pnpm test:e2e               # Playwright E2E (butuh backend & frontend jalan)
```

**Hasil**: 69+ test suites PASS dari 71, coverage >55% statements.

## Scripts

| Script | Deskripsi |
|--------|-----------|
| `pnpm dev` | Development server (port 3000) |
| `pnpm build` | Production build |
| `pnpm start` | Production server |
| `pnpm lint` | ESLint |
| `pnpm test` | Jest unit tests |
| `pnpm test:coverage` | Jest + coverage |
| `pnpm test:e2e` | Playwright E2E |

## Detail Teknis

- **API Client** (`lib/api.ts`): interceptor menambahkan JWT, auto-refresh saat 401, penanganan error 400/401/403/404/409/429/500.
- **AuthContext**: decode JWT, kelola access/refresh token via cookie, expose login/register/logout/refreshUser.
- **Polling real-time**: dashboard, inventory, anomaly refresh setiap 30 detik.
- **Export CSV**: dengan BOM untuk dukungan character UTF-8 (mis. bahasa Indonesia).
- **Autocomplete SKU**: debounce 300ms pada form transaksi.

## Browser Support

Chromium, Firefox, WebKit (untuk Playwright E2E). Browser dapat diinstal via:

```bash
npx playwright install
```

## Lisensi

MIT
