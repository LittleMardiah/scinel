# ROADMAP_3: FRONTEND/MOBILE (REVISED - COMPLETE v1.1)
**Supply Chain Intelligence Platform (ScIntel) - Weeks 5-8**

**Phase:** Frontend/Mobile - Production-Grade UI/UX  
**Duration:** 4 weeks (60 hours)  
**Status:** Ready to Start (after ROADMAP_1 & ROADMAP_2 complete)  
**Target Completion:** End of Week 8  
**Team:** Solo/Small (1-2 developers)  
**Revision:** v1.1 (Comprehensive Audit Fixes + HALAMAN.md v2.1 Alignment 2026-07-19)  
**Blocks:** ROADMAP_4 (Testing & Deployment) — cannot start without ROADMAP_3 complete

---

## 📋 ROADMAP_3 OVERVIEW

### Goal
Build production-grade frontend dashboard & mobile app for supply chain management. **Full user interface** with real-time inventory tracking, demand forecasting visualization, anomaly alerts, transaction management, and comprehensive settings. Frontend communicates with all 23 backend endpoints from ROADMAP_2 (auth 4 + inventory 5 + forecast 3 + anomaly 4 + webhooks 4 + API keys 4 + health 3 + rate limiting 1).

**Key Differentiator:** Frontend is 100% mapped to HALAMAN.md v2.1 specifications - 16 pages, 40+ components, all interactions precisely defined per Bahasa Indonesia UX flow.

### Effort Breakdown
- Week 5: 18 hours (Authentication pages + Dashboard setup + component foundation)
- Week 6: 16 hours (Inventory & Transaction pages + real-time updates)
- Week 7: 16 hours (Forecast & Anomaly pages + interactive charts)
- Week 8: 10 hours (Real-time features polish, mobile optimization, testing & performance)
- **Total: 60 hours**

### Key Deliverables
✅ **16 pages** (4 auth + 12 protected) per HALAMAN.md v2.1 mapping  
✅ **40+ React components** (reusable, tested, TypeScript-typed)  
✅ **Real-time inventory updates** (polling 30s default, WebSocket optional Week 8)  
✅ **Interactive charts & visualizations** (Recharts: demand forecast, accuracy metrics, anomaly trends)  
✅ **Mobile responsive design** (tested: 375px mobile, 768px tablet, 1920px desktop)  
✅ **Dark mode toggle** (theme persistence with localStorage, smooth transitions)  
✅ **Comprehensive error handling** (error boundaries, toast notifications, retry logic, 401/429/5xx)  
✅ **Loading states** (skeleton screens, spinners, suspense boundaries)  
✅ **>50% component test coverage** (unit + integration tests with Vitest + React Testing Library)  
✅ **5+ end-to-end user flows** passing (login → dashboard → transaction → anomaly acknowledgement)  
✅ **Lighthouse score >90** (performance, accessibility, SEO, best practices)  
✅ **API integration complete** (23 endpoints fully integrated, error handling, retry logic)  
✅ **RBAC UI enforcement** (manager sees only own warehouse data, admin sees all, role-based menu)  
✅ **Form validation** (frontend HTML5 + custom validators + backend error display)  
✅ **Rate limiting UI** (429 Too Many Requests → "Terlalu banyak permintaan" message with retry)  

### Dependencies
- ✅ ROADMAP_1 (Foundation) — must be 100% complete
  - Database: PostgreSQL with 12 tables
  - Auth: JWT + RBAC implemented
  - Graceful shutdown: 15-second timeout
  - Database idle timeout: 30 seconds
  - Password reset token TTL: 15 minutes
- ✅ ROADMAP_2 (Backend Core) — all 23 APIs must be working
  - 3 Health check endpoints (GET /health, /readiness, /liveness)
  - 4 Auth endpoints (POST register, login, refresh, logout)
  - 5 Inventory endpoints (GET all/warehouse, POST transactions, GET transactions, POST import)
  - 3 Forecast endpoints (GET forecasts, POST generate, GET accuracy)
  - 4 Anomaly endpoints (GET list, GET detail, PUT update, GET history)
  - 4 Webhook endpoints (POST subscribe, GET subscriptions, DELETE, GET deliveries)
  - 4 API Key endpoints (POST create, GET list, DELETE revoke, POST rotate)
  - Rate limiting: 1000 req/sec per user/API-key
  - Database: PostgreSQL with caching (inventory 30s TTL, forecast 24h TTL)

### Tech Stack (LOCKED - No Changes)
- **Framework:** Next.js 14+ with App Router (React Server Components)
- **Language:** TypeScript 5+ (strict mode: true)
- **UI Library:** React 18+ with hooks (useContext, useReducer, useCallback, useMemo)
- **Styling:** Tailwind CSS v3+ (utility-first, dark mode support)
- **State Management:** Context API + useReducer (no Redux/Zustan/Recoil for MVP)
- **Charts:** Recharts (lightweight SVG charts, no D3.js, no ApexCharts)
- **Icons:** Lucide React (SVG icons, tree-shakeable, 400+ icons)
- **HTTP Client:** Axios with interceptors (JWT refresh, error handling, timeout 30s)
- **Form Validation:** HTML5 built-in + custom regex validators
- **Cookie Management:** js-cookie (secure, httpOnly support)
- **Package Manager:** pnpm (faster than npm, disk-efficient)
- **Testing:** Vitest (faster than Jest) + React Testing Library (@testing-library/react)
- **Build Tool:** Next.js built-in Webpack (automatic code splitting, SSR, ISR)
- **CSS-in-JS:** Tailwind (zero runtime CSS, static output, dark mode via class)
- **Date Formatting:** date-fns (tree-shakeable, alternative to moment.js)
- **Utility Functions:** clsx + tailwind-merge (className composition)

### Success Criteria
1. ✅ All 16 pages built and functional (per HALAMAN.md v2.1)
2. ✅ All 40+ components created, tested, and styled
3. ✅ Real-time inventory updates working (polling 30s + webhook integration)
4. ✅ Dark mode toggle functional (theme persistence with localStorage)
5. ✅ All interactive charts rendering correctly (Recharts line/bar/area charts)
6. ✅ Mobile responsive (tested: 375px mobile, 768px tablet, 1920px desktop)
7. ✅ >50% component test coverage (unit + integration tests)
8. ✅ 5+ E2E user flows passing (documented, repeatable)
9. ✅ Lighthouse score >90 (performance, accessibility, SEO, best practices)
10. ✅ Error handling + loading states complete (401, 429, 5xx, network errors)
11. ✅ RBAC UI enforcement working (manager warehouse isolation, role-based menu)
12. ✅ Form validation (frontend + backend error display)
13. ✅ Rate limiting UI (429 → "Terlalu banyak permintaan" with retry countdown)
14. ✅ All pages styled with Tailwind + dark mode support
15. ✅ Ready for ROADMAP_4 (Comprehensive testing, documentation, deployment)

---

## 📄 PAGE MAPPING (16 pages per HALAMAN.md v2.1)

### Authentication Pages (4 pages - public, no auth required)

**Page 1.1: Login Page**
- Route: `/login` or `/auth/login`
- Endpoint: POST /api/v1/auth/login
- Components: LoginForm (email + password), error display, forgot password link, register link
- Flow per HALAMAN.md 1.1: Email validation → backend verification → JWT token storage → redirect /dashboard
- Error cases: 401 (invalid credentials), 403 (unverified/suspended), 429 (rate limited), 500 (server error)
- Special: Failed attempt tracking (5 attempts in 5 min → lock for 5 min)

**Page 1.2: Register Page**
- Route: `/register` or `/auth/register`
- Endpoint: POST /api/v1/auth/register
- Components: RegisterForm (email, password, confirm password, full name)
- Validation: Email format, password strength (min 8 chars), password match
- Flow: Form validation → backend registration → auto-login → redirect /dashboard
- Error cases: 409 (duplicate email), 400 (validation), 500

**Page 1.3: Forgot Password Page**
- Route: `/forgot-password`
- Endpoint: POST /api/v1/auth/password-reset (future endpoint)
- Components: Email input form, success message
- Flow: Enter email → backend sends reset link → success toast

**Page 1.4: Reset Password Page**
- Route: `/reset-password?token=TOKEN_HERE`
- Endpoint: PUT /api/v1/auth/password-reset/{token} (future endpoint)
- Components: Password input form, token validation
- Flow: Validate token → set new password → redirect /login

### Protected Pages (12 pages - requires auth)

**Page 2.1: Dashboard**
- Route: `/dashboard`
- Endpoints: GET /health, GET /readiness, GET /inventory (count), GET /anomalies (open count), GET /forecasts (accuracy)
- Components: StatsCard (4 KPIs), RecentTransactions (table), ForecastChart (line chart), SystemStatus (health check), QuickActionsCard
- RBAC: All roles can view (data filtered by role on backend)
- Real-time: Polling every 30 seconds for data refresh
- Bahasa: "Dashboard", "Ringkasan supply chain Anda", "Total SKU", "Stok Rendah", "Anomali Aktif", "Akurasi Forecast"

**Page 2.2: Inventory List**
- Route: `/inventory`
- Endpoints: GET /inventory (with filters: warehouse_id, status), POST /transactions, POST /inventory/import
- Components: InventoryTable (sortable, filterable), InventoryFilters (warehouse dropdown, status radio), ImportButton
- Columns: SKU Code, Name, Quantity, Status (normal/low/critical), Warehouse, Last Restock, Actions
- Actions: View detail, Create transaction, Edit
- Real-time: Auto-refresh every 30s, cache invalidation on update

**Page 2.3: Warehouse-Specific Inventory**
- Route: `/inventory/:warehouse`
- Endpoints: GET /inventory/:warehouse (filtered by warehouse)
- Components: InventoryTable (same as 2.2, but pre-filtered)
- RBAC: Manager can only see own warehouse, Admin/Analyst can see any
- Breadcrumb: Dashboard > Inventory > [Warehouse Name]

**Page 2.4: New Transaction**
- Route: `/inventory/transaction/new`
- Endpoints: POST /transactions
- Components: TransactionForm (warehouse select, SKU select, type radio [in/out], quantity input, reason input, details textarea)
- Validation: All fields required, quantity > 0, sufficient stock check for "out" type
- Flow: Form fill → validate → POST /transactions → success toast + redirect /inventory/transactions
- Error: 409 (insufficient stock), 400 (validation), 401 (auth), 429 (rate limit)
- Bahasa: "Buat Transaksi Baru", "Jenis", "Masuk", "Keluar", "Kuantitas", "Alasan", "Rincian Tambahan"

**Page 2.5: Transaction History**
- Route: `/inventory/transactions`
- Endpoints: GET /transactions (with pagination: limit, offset), GET /transactions?warehouse_id=X&type=in/out
- Components: TransactionTable (date, SKU, type, quantity, warehouse, reason), Pagination
- Filters: Date range, warehouse, transaction type (in/out)
- Sort: Date (default DESC), SKU, quantity, warehouse
- Pagination: 20 items per page, prev/next navigation
- Real-time: Auto-refresh every 30s

**Page 2.6: Forecast Display**
- Route: `/forecast`
- Endpoints: GET /forecasts (with filters: warehouse_id, sku_id, days)
- Components: ForecastChart (Recharts LineChart with confidence intervals - lower_bound/upper_bound), ForecastTable (SKU, forecast_date, predicted_demand, lower/upper bound)
- Visualization: Line chart with shaded confidence interval (95%), legend showing predicted vs bounds
- Filters: Warehouse select, SKU select, forecast days (default 30)
- Bahasa: "Forecast Permintaan", "Prediksi Permintaan", "Interval Kepercayaan (95%)", "Algoritma", "Tren"

**Page 2.7: Forecast Accuracy Metrics**
- Route: `/forecast/accuracy`
- Endpoints: GET /forecasts/accuracy (with filters: warehouse_id, algorithm)
- Components: AccuracyTable (algorithm, MAPE %, RMSE, MAE, evaluation_date), AlgorithmComparison chart (ETS vs ARIMA)
- Metrics Display: MAPE (%), RMSE (units), MAE (units), evaluation date
- Comparison Chart: Recharts BarChart comparing MAPE/RMSE for ETS vs ARIMA
- Interpretation: Lower MAPE/RMSE/MAE = better accuracy
- Admin Feature: Button to manually trigger forecast batch job (POST /forecasts/generate)

**Page 2.8: Anomaly List**
- Route: `/anomalies`
- Endpoints: GET /anomalies (with filters: status [open/acknowledged/resolved], severity [low/medium/high/critical])
- Components: AnomalyTable (SKU, type, severity badge, description, z_score, detected_at, status, actions), Filters (status radio, severity radio)
- Severity Colors: low=yellow, medium=orange, high=red, critical=dark red
- Actions: View detail, Acknowledge, Resolve
- Real-time: Auto-refresh every 30s
- Bahasa: "Anomali yang Terdeteksi", "Jenis", "Tingkat Keparahan", "Deskripsi", "Z-Score", "Status"

**Page 2.9: Anomaly Detail**
- Route: `/anomalies/:id`
- Endpoints: GET /anomalies/:id (detail + full info), PUT /anomalies/:id (update status), GET /anomalies/:id/history (timeline)
- Components: AnomalyHeader (SKU, type, severity), AnomalyTimeline (detected → acknowledged → resolved), ResolutionModal (update status + notes), RelatedTransactions (recent transactions for this SKU)
- Display: Historical average, detected value, z-score calculation, trend
- Actions: Acknowledge (status update), Resolve with notes
- Timeline: Show all status changes with timestamps and user info

**Page 2.10: Anomaly Analytics**
- Route: `/anomalies/analytics`
- Endpoints: GET /anomalies (all), GET /forecasts/accuracy (trends)
- Components: AnomalyTrendChart (Recharts AreaChart showing anomaly count over time), SeverityBreakdown (pie chart), TopAffectedSKUs (bar chart)
- Trends: Anomalies by severity over last 30 days, detection rate, resolution rate
- Analysis: Which SKUs have most anomalies, which warehouses, which time patterns

**Page 2.11: User Profile**
- Route: `/profile`
- Endpoints: GET /auth/me (future endpoint), PUT /auth/me (update profile)
- Components: ProfileForm (email [read-only], full name [editable], role [read-only], last login, created date)
- Actions: Edit name, change password (separate modal)
- Display: User info, account creation date, last login date/time

**Page 2.12: Settings**
- Route: `/settings`
- Endpoints: Various (theme stored in localStorage, API keys from GET /api-keys, webhooks from GET /webhooks/subscriptions)
- Components: ThemeToggle (light/dark/system), APIKeyManager (list, create, revoke, rotate), WebhookManager (list, subscribe, delete), NotificationPreferences
- Features:
  - Theme: Light/Dark/System preference (radio buttons)
  - API Keys: Table with key_prefix, name, permissions, last_used, expires, actions (revoke, rotate)
  - Webhooks: Table with event_type, webhook_url, is_active toggle, last_delivery_status
  - Notifications: Email notifications preferences (Bahasa: "Pengaturan")

### Special Pages (not in HALAMAN.md but required)

**Page 3.1: 404 Not Found**
- Route: `/404`
- Display: "Halaman tidak ditemukan", back to dashboard button

**Page 3.2: 500 Server Error**
- Route: `/500`
- Display: "Terjadi kesalahan server", request_id, retry button, support link

**Page 3.3: Unauthorized (401)**
- Route: `/401`
- Display: "Anda tidak berhak mengakses halaman ini", role info, back to dashboard

---

## 📅 DETAILED BREAKDOWN (60 hours total)

---

## WEEK 5: AUTHENTICATION PAGES & DASHBOARD (18 hours)

### Week 5.1: Project Setup & Component Foundation (4 hours)

#### Task 5.1.1: Next.js Project Initialization (1.5 hours)

**Commands:**
```bash
# Create Next.js with TypeScript, Tailwind, App Router
npx create-next-app@latest scinel-frontend \
  --typescript \
  --tailwind \
  --app \
  --eslint \
  --no-git \
  --import-alias '@/*'

cd scinel-frontend

# Install dependencies (30+ packages)
pnpm install \
  recharts \
  lucide-react \
  axios \
  js-cookie \
  date-fns \
  clsx \
  tailwind-merge \
  next-themes

# Dev dependencies
pnpm add -D \
  @testing-library/react \
  @testing-library/jest-dom \
  @testing-library/user-event \
  vitest \
  @vitest/ui \
  jsdom \
  @types/js-cookie
```

**Project Structure (40+ directories):**

```
scinel-frontend/
├── app/
│   ├── layout.tsx                    # Root layout + providers
│   ├── page.tsx                      # Home redirect
│   ├── (auth)/                       # Auth layout group
│   │   ├── layout.tsx
│   │   ├── login/page.tsx
│   │   ├── register/page.tsx
│   │   ├── forgot-password/page.tsx
│   │   └── reset-password/page.tsx
│   ├── (protected)/                  # Protected routes layout
│   │   ├── layout.tsx
│   │   ├── dashboard/page.tsx
│   │   ├── inventory/
│   │   │   ├── page.tsx
│   │   │   ├── [warehouse]/page.tsx
│   │   │   ├── transaction/new/page.tsx
│   │   │   └── transactions/page.tsx
│   │   ├── forecast/
│   │   │   ├── page.tsx
│   │   │   └── accuracy/page.tsx
│   │   ├── anomalies/
│   │   │   ├── page.tsx
│   │   │   ├── [id]/page.tsx
│   │   │   └── analytics/page.tsx
│   │   ├── profile/page.tsx
│   │   └── settings/page.tsx
│   ├── api/
│   │   └── revalidate-tag/route.ts   # ISR revalidation
│   └── globals.css                   # Tailwind + global styles
├── components/
│   ├── layout/
│   │   ├── Navbar.tsx
│   │   ├── Sidebar.tsx
│   │   └── Footer.tsx
│   ├── auth/
│   │   ├── LoginForm.tsx
│   │   ├── RegisterForm.tsx
│   │   ├── ForgotPasswordForm.tsx
│   │   └── ResetPasswordForm.tsx
│   ├── dashboard/
│   │   ├── StatsCard.tsx
│   │   ├── KPIGrid.tsx
│   │   ├── RecentTransactions.tsx
│   │   ├── ForecastChart.tsx
│   │   ├── SystemStatus.tsx
│   │   └── QuickActionsCard.tsx
│   ├── inventory/
│   │   ├── InventoryTable.tsx
│   │   ├── InventoryFilters.tsx
│   │   ├── TransactionForm.tsx
│   │   ├── InventoryStatusBadge.tsx
│   │   └── InventoryHistory.tsx
│   ├── forecast/
│   │   ├── ForecastChart.tsx
│   │   ├── ForecastTable.tsx
│   │   ├── AccuracyMetrics.tsx
│   │   ├── AlgorithmComparison.tsx
│   │   └── ConfidenceIntervalBadge.tsx
│   ├── anomaly/
│   │   ├── AnomalyList.tsx
│   │   ├── AnomalyCard.tsx
│   │   ├── SeverityBadge.tsx
│   │   ├── AnomalyTimeline.tsx
│   │   ├── ResolutionModal.tsx
│   │   └── AnomalyChart.tsx
│   ├── common/
│   │   ├── Button.tsx
│   │   ├── Card.tsx
│   │   ├── Modal.tsx
│   │   ├── Toast.tsx
│   │   ├── Loading.tsx
│   │   ├── ErrorBoundary.tsx
│   │   ├── Skeleton.tsx
│   │   ├── Pagination.tsx
│   │   ├── ConfirmDialog.tsx
│   │   └── Breadcrumb.tsx
│   └── ui/
│       ├── Input.tsx
│       ├── Select.tsx
│       ├── Checkbox.tsx
│       ├── Radio.tsx
│       ├── DatePicker.tsx
│       ├── TextArea.tsx
│       └── Badge.tsx
├── hooks/
│   ├── useAuth.ts
│   ├── useApi.ts
│   ├── useInventory.ts
│   ├── useForecast.ts
│   ├── useAnomalies.ts
│   ├── useLocalStorage.ts
│   ├── useForm.ts
│   ├── useDebounce.ts
│   ├── useTheme.ts
│   └── useInfiniteScroll.ts
├── context/
│   ├── AuthContext.tsx
│   └── ThemeContext.tsx
├── lib/
│   ├── api.ts
│   ├── auth.ts
│   ├── validators.ts
│   ├── utils.ts
│   ├── constants.ts
│   └── types.ts
├── types/
│   ├── index.ts
│   ├── api.ts
│   ├── auth.ts
│   ├── inventory.ts
│   ├── forecast.ts
│   ├── anomaly.ts
│   └── webhook.ts
├── tests/
│   ├── components/
│   │   ├── Button.test.tsx
│   │   ├── LoginForm.test.tsx
│   │   ├── StatsCard.test.tsx
│   │   ├── InventoryTable.test.tsx
│   │   ├── ForecastChart.test.tsx
│   │   └── AnomalyList.test.tsx
│   ├── hooks/
│   │   ├── useAuth.test.ts
│   │   ├── useApi.test.ts
│   │   └── useTheme.test.ts
│   └── utils/
│       └── validators.test.ts
├── public/
│   ├── logo.svg
│   ├── favicon.ico
│   └── icons/
├── .env.example
├── next.config.js
├── tailwind.config.ts
├── tsconfig.json
├── vitest.config.ts
├── package.json
├── pnpm-lock.yaml
└── README.md
```

**Configuration Files:**

**File: `.env.example`**
```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_API_TIMEOUT=30000

# Real-time Configuration
NEXT_PUBLIC_POLLING_INTERVAL=30000  # 30 seconds for inventory updates
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:8080/ws  # Optional

# Authentication
NEXT_PUBLIC_AUTH_REDIRECT=/dashboard
NEXT_PUBLIC_LOGIN_REDIRECT=/login
NEXT_PUBLIC_TOKEN_STORAGE=cookie  # or 'localStorage'

# Feature Flags
NEXT_PUBLIC_ENABLE_WEBSOCKET=false  # Enable in Week 8
NEXT_PUBLIC_ENABLE_DARK_MODE=true
NEXT_PUBLIC_ENABLE_ANALYTICS=true
NEXT_PUBLIC_ENABLE_ERROR_TRACKING=false  # Sentry/Rollbar

# Monitoring (optional)
# NEXT_PUBLIC_SENTRY_DSN=https://...@sentry.io/...
```

**File: `tailwind.config.ts`**
```typescript
import type { Config } from 'tailwindcss'
import defaultTheme from 'tailwindcss/defaultTheme'

const config: Config = {
  darkMode: 'class', // Dark mode via .dark class
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
          800: '#075985',
          900: '#0c3d66',
        },
        success: {
          50: '#f0fdf4',
          600: '#16a34a',
          700: '#15803d',
        },
        warning: {
          50: '#fffbeb',
          600: '#f59e0b',
          700: '#d97706',
        },
        error: {
          50: '#fef2f2',
          600: '#dc2626',
          700: '#b91c1c',
        },
        neutral: {
          50: '#fafafa',
          100: '#f5f5f5',
          200: '#e5e5e5',
          300: '#d4d4d4',
          400: '#a3a3a3',
          500: '#737373',
          600: '#525252',
          700: '#404040',
          800: '#262626',
          900: '#171717',
        },
      },
      fontFamily: {
        sans: ['Inter', ...defaultTheme.fontFamily.sans],
      },
      spacing: {
        xs: '0.5rem',   // 8px
        sm: '0.75rem',  // 12px
        md: '1rem',     // 16px
        lg: '1.5rem',   // 24px
        xl: '2rem',     // 32px
        '2xl': '3rem',  // 48px
      },
      fontSize: {
        xs: ['0.75rem', { lineHeight: '1rem' }],
        sm: ['0.875rem', { lineHeight: '1.25rem' }],
        base: ['1rem', { lineHeight: '1.5rem' }],
        lg: ['1.125rem', { lineHeight: '1.75rem' }],
        xl: ['1.25rem', { lineHeight: '1.75rem' }],
        '2xl': ['1.5rem', { lineHeight: '2rem' }],
        '3xl': ['1.875rem', { lineHeight: '2.25rem' }],
      },
      boxShadow: {
        sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        md: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
        lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1)',
      },
    },
  },
  plugins: [],
}
export default config
```

**File: `next.config.js`**
```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  typescript: {
    strictNullChecks: true,
  },
  compress: true,
  poweredByHeader: false,
  reactStrictMode: true,
  swcMinify: true,
  images: {
    formats: ['image/avif', 'image/webp'],
  },
}

module.exports = nextConfig
```

**File: `vitest.config.ts`**
```typescript
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    include: ['**/*.test.ts', '**/*.test.tsx'],
    exclude: ['node_modules', 'dist', '.idea', '.git', '.cache'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      statements: 50,
      branches: 50,
      functions: 50,
      lines: 50,
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './'),
    },
  },
})
```

**Deliverable:**
- ✅ Next.js 14+ project with TypeScript
- ✅ Tailwind CSS v3+ configured with custom colors
- ✅ Dark mode support (class strategy)
- ✅ File structure (40+ directories prepared)
- ✅ Environment configuration template
- ✅ Testing setup (Vitest + React Testing Library)

---

#### Task 5.1.2: Core Utilities, Providers & Hooks (1.5 hours)

**File: `lib/api.ts` (API client with JWT + error handling)**

```typescript
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from 'axios';
import Cookies from 'js-cookie';
import { Toast } from '@/components/common/Toast';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const API_TIMEOUT = parseInt(process.env.NEXT_PUBLIC_API_TIMEOUT || '30000');

export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: API_TIMEOUT,
});

// Request interceptor: Add JWT token to all requests
apiClient.interceptors.request.use((config) => {
  const token = Cookies.get('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor: Handle errors & token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };

    // Handle 401 Unauthorized - try to refresh token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = Cookies.get('refresh_token');
        if (!refreshToken) {
          // No refresh token available
          Cookies.remove('access_token');
          Cookies.remove('refresh_token');
          window.location.href = '/login';
          return Promise.reject(error);
        }

        // Attempt to refresh access token
        const response = await axios.post(
          `${API_BASE}/auth/refresh`,
          { refresh_token: refreshToken },
          {
            headers: { 'Content-Type': 'application/json' },
            timeout: API_TIMEOUT,
          }
        );

        const { access_token, refresh_token } = response.data.data;

        // Store new tokens securely
        Cookies.set('access_token', access_token, {
          secure: true,
          sameSite: 'strict',
          expires: 1 / 24, // 1 hour
        });
        Cookies.set('refresh_token', refresh_token, {
          secure: true,
          sameSite: 'strict',
          expires: 7, // 7 days
        });

        // Retry original request with new token
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${access_token}`;
        }
        return apiClient(originalRequest);
      } catch (refreshError) {
        // Refresh failed, redirect to login
        Cookies.remove('access_token');
        Cookies.remove('refresh_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    // Handle 429 Too Many Requests
    if (error.response?.status === 429) {
      const retryAfter = error.response.headers['retry-after'] || '60';
      const message = `Terlalu banyak permintaan. Coba lagi dalam ${retryAfter} detik.`;
      Toast.error(message);
      return Promise.reject({
        ...error,
        message,
        retryAfter: parseInt(retryAfter),
      });
    }

    // Handle 403 Forbidden (insufficient permissions)
    if (error.response?.status === 403) {
      Toast.error('Anda tidak berhak mengakses resource ini.');
      return Promise.reject(error);
    }

    // Handle 404 Not Found
    if (error.response?.status === 404) {
      Toast.error('Resource tidak ditemukan.');
      return Promise.reject(error);
    }

    // Handle 500+ Server Errors
    if (error.response?.status && error.response.status >= 500) {
      const requestId = error.response.headers['x-request-id'] || 'unknown';
      Toast.error(`Terjadi kesalahan server (ID: ${requestId}). Tim kami sedang menangani.`);
      return Promise.reject({
        ...error,
        requestId,
      });
    }

    // Handle network errors (no response from server)
    if (!error.response) {
      Toast.error('Koneksi gagal. Periksa internet Anda dan coba lagi.');
      return Promise.reject(error);
    }

    // Generic error
    Toast.error(error.message || 'Terjadi kesalahan. Coba lagi nanti.');
    return Promise.reject(error);
  }
);

export default apiClient;
```

**File: `context/AuthContext.tsx` (Authentication state)**

```typescript
'use client'; // Client component

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  ReactNode,
} from 'react';
import Cookies from 'js-cookie';
import apiClient from '@/lib/api';

export interface User {
  id: string;
  email: string;
  full_name?: string;
  role: 'admin' | 'manager' | 'analyst' | 'viewer';
  status: 'active' | 'suspended' | 'pending_email_verification';
  created_at: string;
  last_login_at?: string;
}

export interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Decode JWT and extract user info
  const decodeToken = useCallback((token: string): User | null => {
    try {
      const parts = token.split('.');
      if (parts.length !== 3) return null;

      const payload = JSON.parse(atob(parts[1]));
      return {
        id: payload.user_id,
        email: payload.email,
        role: payload.role,
        status: payload.status || 'active',
        created_at: new Date(payload.iat * 1000).toISOString(),
      };
    } catch {
      return null;
    }
  }, []);

  // Initialize auth on mount
  useEffect(() => {
    const token = Cookies.get('access_token');
    if (token) {
      const decodedUser = decodeToken(token);
      if (decodedUser) {
        setUser(decodedUser);
      } else {
        // Invalid token, clear cookies
        Cookies.remove('access_token');
        Cookies.remove('refresh_token');
      }
    }
    setIsLoading(false);
  }, [decodeToken]);

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      const response = await apiClient.post('/auth/login', { email, password });
      const { access_token, refresh_token, user: userData } = response.data.data;

      // Store tokens securely
      Cookies.set('access_token', access_token, {
        secure: true,
        sameSite: 'strict',
        expires: 1 / 24,
      });
      Cookies.set('refresh_token', refresh_token, {
        secure: true,
        sameSite: 'strict',
        expires: 7,
      });

      setUser(userData);
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (email: string, password: string, fullName: string) => {
    setIsLoading(true);
    try {
      const response = await apiClient.post('/auth/register', {
        email,
        password,
        full_name: fullName,
      });
      const { access_token, refresh_token, user: userData } = response.data.data;

      Cookies.set('access_token', access_token, {
        secure: true,
        sameSite: 'strict',
        expires: 1 / 24,
      });
      Cookies.set('refresh_token', refresh_token, {
        secure: true,
        sameSite: 'strict',
        expires: 7,
      });

      setUser(userData);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      await apiClient.post('/auth/logout');
    } finally {
      Cookies.remove('access_token');
      Cookies.remove('refresh_token');
      setUser(null);
    }
  };

  const refreshUser = async () => {
    const token = Cookies.get('access_token');
    if (token) {
      const decodedUser = decodeToken(token);
      setUser(decodedUser);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
```

**File: `context/ThemeContext.tsx` (Dark mode)**

```typescript
'use client';

import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from 'react';

type Theme = 'light' | 'dark' | 'system';

interface ThemeContextType {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  isDark: boolean;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setThemeState] = useState<Theme>('system');
  const [isDark, setIsDark] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    // Read theme from localStorage or use system preference
    const storedTheme = (localStorage.getItem('theme') as Theme | null) || 'system';
    setThemeState(storedTheme);

    // Determine if should be dark
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const shouldBeDark = storedTheme === 'dark' || (storedTheme === 'system' && prefersDark);

    // Apply theme to DOM
    const root = document.documentElement;
    if (shouldBeDark) {
      root.classList.add('dark');
      setIsDark(true);
    } else {
      root.classList.remove('dark');
      setIsDark(false);
    }

    setMounted(true);
  }, []);

  const setTheme = (newTheme: Theme) => {
    setThemeState(newTheme);
    localStorage.setItem('theme', newTheme);

    const root = document.documentElement;
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const shouldBeDark = newTheme === 'dark' || (newTheme === 'system' && prefersDark);

    if (shouldBeDark) {
      root.classList.add('dark');
      setIsDark(true);
    } else {
      root.classList.remove('dark');
      setIsDark(false);
    }
  };

  // Don't render until mounted to avoid hydration mismatch
  if (!mounted) {
    return <>{children}</>;
  }

  return (
    <ThemeContext.Provider value={{ theme, setTheme, isDark }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within ThemeProvider');
  }
  return context;
}
```

**File: `hooks/useApi.ts` (Custom hook for API calls with loading/error state)**

```typescript
import { useState, useCallback, useEffect } from 'react';
import apiClient from '@/lib/api';

interface UseApiOptions {
  autoFetch?: boolean;
  onSuccess?: (data: any) => void;
  onError?: (error: any) => void;
}

export function useApi<T = any>(
  url: string,
  options: UseApiOptions = {}
) {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<any>(null);

  const fetch = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(url);
      setData(response.data.data);
      options.onSuccess?.(response.data.data);
    } catch (err) {
      setError(err);
      options.onError?.(err);
    } finally {
      setIsLoading(false);
    }
  }, [url, options]);

  useEffect(() => {
    if (options.autoFetch !== false) {
      fetch();
    }
  }, [fetch, options.autoFetch]);

  return { data, isLoading, error, refetch: fetch };
}
```

**File: `app/layout.tsx` (Root layout with providers)**

```typescript
import type { Metadata } from 'next';
import { AuthProvider } from '@/context/AuthContext';
import { ThemeProvider } from '@/context/ThemeContext';
import './globals.css';

export const metadata: Metadata = {
  title: 'ScIntel - Supply Chain Intelligence',
  description: 'Platform manajemen supply chain dengan AI forecasting',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id" suppressHydrationWarning>
      <body className="bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-white">
        <ThemeProvider>
          <AuthProvider>
            {children}
          </AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
```

**Deliverable:**
- ✅ API client with JWT refresh + error handling
- ✅ AuthContext with login/register/logout
- ✅ ThemeContext with dark mode persistence
- ✅ useApi hook for data fetching
- ✅ Root layout with providers
- ✅ Secure cookie management (httpOnly, secure, sameSite)

---

#### Task 5.1.3: Reusable Component Library (1 hour)

Due to space constraints, here are 5 core component examples. Implement remaining 35+ components following the same pattern:

**File: `components/common/Button.tsx`**

```typescript
import React from 'react';
import { cn } from '@/lib/utils';

interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    { className, variant = 'primary', size = 'md', isLoading, disabled, children, ...props },
    ref
  ) => {
    const baseClasses =
      'inline-flex items-center justify-center font-medium rounded-md transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';

    const variantClasses = {
      primary: 'bg-blue-600 text-white hover:bg-blue-700 active:bg-blue-800 focus-visible:ring-blue-500',
      secondary: 'bg-gray-200 text-gray-900 hover:bg-gray-300 active:bg-gray-400 focus-visible:ring-gray-500 dark:bg-gray-700 dark:text-white',
      danger: 'bg-red-600 text-white hover:bg-red-700 active:bg-red-800 focus-visible:ring-red-500',
      ghost: 'hover:bg-gray-100 dark:hover:bg-gray-800 focus-visible:ring-gray-500',
    };

    const sizeClasses = {
      sm: 'px-3 py-1.5 text-sm',
      md: 'px-4 py-2 text-base',
      lg: 'px-6 py-3 text-lg',
    };

    return (
      <button
        className={cn(
          baseClasses,
          variantClasses[variant],
          sizeClasses[size],
          className
        )}
        disabled={disabled || isLoading}
        ref={ref}
        {...props}
      >
        {isLoading && <span className="mr-2 inline-block animate-spin">⊙</span>}
        {children}
      </button>
    );
  }
);
Button.displayName = 'Button';
```

**File: `components/common/Card.tsx`**

```typescript
import React from 'react';
import { cn } from '@/lib/utils';

export const Card = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      'rounded-lg border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-700 dark:bg-gray-800',
      className
    )}
    {...props}
  />
));
Card.displayName = 'Card';

export const CardHeader = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('mb-4 pb-4 border-b border-gray-200 dark:border-gray-700', className)}
    {...props}
  />
));
CardHeader.displayName = 'CardHeader';

export const CardTitle = React.forwardRef<
  HTMLHeadingElement,
  React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
  <h2
    ref={ref}
    className={cn('text-xl font-bold text-gray-900 dark:text-white', className)}
    {...props}
  />
));
CardTitle.displayName = 'CardTitle';

export const CardContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div ref={ref} className={className} {...props} />
));
CardContent.displayName = 'CardContent';

export const CardFooter = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('mt-4 pt-4 border-t border-gray-200 dark:border-gray-700', className)}
    {...props}
  />
));
CardFooter.displayName = 'CardFooter';
```

**Remaining components to create (35+):**

- `components/common/Input.tsx` - Text input with validation state
- `components/common/Select.tsx` - Dropdown select
- `components/common/Modal.tsx` - Dialog with backdrop
- `components/common/Toast.tsx` - Notification system
- `components/common/Loading.tsx` - Loading spinner
- `components/common/ErrorBoundary.tsx` - Error boundary
- `components/common/Skeleton.tsx` - Skeleton screen
- `components/common/Pagination.tsx` - Pagination controls
- `components/common/ConfirmDialog.tsx` - Confirmation modal
- `components/common/Breadcrumb.tsx` - Breadcrumb navigation
- `components/layout/Navbar.tsx` - Top navigation bar
- `components/layout/Sidebar.tsx` - Side navigation menu
- `components/layout/Footer.tsx` - Footer section
- `components/auth/LoginForm.tsx` - Login form
- `components/auth/RegisterForm.tsx` - Register form
- `components/auth/PasswordResetForm.tsx` - Password reset
- `components/dashboard/StatsCard.tsx` - KPI card
- `components/dashboard/KPIGrid.tsx` - 4-column KPI layout
- `components/dashboard/RecentTransactions.tsx` - Transactions table
- `components/dashboard/ForecastChart.tsx` - Line chart
- `components/dashboard/SystemStatus.tsx` - Health check
- `components/dashboard/QuickActionsCard.tsx` - Quick action buttons
- `components/inventory/InventoryTable.tsx` - Data table
- `components/inventory/InventoryFilters.tsx` - Filter controls
- `components/inventory/TransactionForm.tsx` - Transaction form
- `components/inventory/InventoryStatusBadge.tsx` - Status badge
- `components/inventory/InventoryHistory.tsx` - Transaction history
- `components/forecast/ForecastChart.tsx` - Recharts line chart
- `components/forecast/ForecastTable.tsx` - Forecast table
- `components/forecast/AccuracyMetrics.tsx` - Metrics display
- `components/forecast/AlgorithmComparison.tsx` - ETS vs ARIMA
- `components/anomaly/AnomalyList.tsx` - Anomalies table
- `components/anomaly/AnomalyCard.tsx` - Anomaly card
- `components/anomaly/SeverityBadge.tsx` - Severity badge
- `components/anomaly/AnomalyTimeline.tsx` - Status timeline
- `components/anomaly/ResolutionModal.tsx` - Resolution form
- `components/anomaly/AnomalyChart.tsx` - Trend chart
- `components/ui/Badge.tsx` - Badge component
- `components/ui/TextArea.tsx` - Textarea input
- Plus `Input.tsx`, `Select.tsx`, `Checkbox.tsx`, `Radio.tsx`, `DatePicker.tsx`

**Deliverable:**
- ✅ 40+ reusable React components
- ✅ TypeScript types for all components
- ✅ Tailwind CSS styling + dark mode
- ✅ Accessible ARIA labels
- ✅ Composition patterns (Header, Title, Content, Footer)

---

### Week 5.2: Authentication Pages (5 hours)

#### Task 5.2.1: Login Page (2 hours)

**File: `app/(auth)/login/page.tsx`**

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { LoginForm } from '@/components/auth/LoginForm';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Toast } from '@/components/common/Toast';

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async (email: string, password: string) => {
    setIsLoading(true);
    setError(null);
    try {
      await login(email, password);
      Toast.success('Berhasil masuk! Membawa Anda ke dashboard...');
      router.push('/dashboard');
    } catch (err: any) {
      const message =
        err.response?.data?.message ||
        err.message ||
        'Login gagal. Coba lagi.';
      setError(message);
      Toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
      <div className="w-full max-w-md">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-2xl mb-2">Masuk ke ScIntel</CardTitle>
            <p className="text-gray-600 dark:text-gray-400">
              Kelola supply chain Anda dengan cerdas
            </p>
          </CardHeader>

          <CardContent>
            <LoginForm
              onSubmit={handleLogin}
              isLoading={isLoading}
              error={error}
            />

            <div className="mt-6 space-y-3">
              <div className="text-center">
                <Link
                  href="/forgot-password"
                  className="text-blue-600 hover:text-blue-700 text-sm font-medium"
                >
                  Lupa kata sandi?
                </Link>
              </div>
              <div className="text-center">
                <p className="text-gray-600 dark:text-gray-400 text-sm">
                  Belum punya akun?{' '}
                  <Link
                    href="/register"
                    className="text-blue-600 hover:text-blue-700 font-medium"
                  >
                    Daftar sekarang
                  </Link>
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Footer */}
        <div className="text-center mt-6 text-xs text-gray-500">
          <p>ScIntel v1.0.0 • © 2026 Supply Chain Intelligence</p>
        </div>
      </div>
    </div>
  );
}
```

**File: `components/auth/LoginForm.tsx`**

```typescript
'use client';

import { useState, FormEvent } from 'react';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';

interface LoginFormProps {
  onSubmit: (email: string, password: string) => Promise<void>;
  isLoading?: boolean;
  error?: string | null;
}

export function LoginForm({ onSubmit, isLoading, error }: LoginFormProps) {
  const [formData, setFormData] = useState({ email: '', password: '' });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const errors: Record<string, string> = {};

    if (!formData.email) {
      errors.email = 'Email diperlukan';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      errors.email = 'Format email tidak valid';
    }

    if (!formData.password) {
      errors.password = 'Kata sandi diperlukan';
    } else if (formData.password.length < 6) {
      errors.password = 'Kata sandi minimal 6 karakter';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    await onSubmit(formData.email, formData.password);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-200 rounded-md text-sm border border-red-200 dark:border-red-800">
          {error}
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Email
        </label>
        <Input
          type="email"
          placeholder="nama@perusahaan.com"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          error={formErrors.email}
          disabled={isLoading}
          autoComplete="email"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Kata Sandi
        </label>
        <Input
          type="password"
          placeholder="Masukkan kata sandi"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          error={formErrors.password}
          disabled={isLoading}
          autoComplete="current-password"
        />
      </div>

      <Button
        type="submit"
        className="w-full"
        disabled={isLoading || Object.keys(formErrors).length > 0}
        isLoading={isLoading}
      >
        {isLoading ? 'Memproses...' : 'Masuk'}
      </Button>
    </form>
  );
}
```

**Deliverable:**
- ✅ Login page with form validation
- ✅ Error handling + display
- ✅ JWT token storage + redirect
- ✅ Responsive design + dark mode
- ✅ Bahasa Indonesia text

---

#### Task 5.2.2: Register Page (1.5 hours)

**Similar pattern to login with:**
- Email validation
- Password strength check (min 8 chars, uppercase, lowercase, number)
- Password confirmation matching
- Full name input
- Error handling for duplicate email (409 Conflict)
- Auto-login after registration

**Deliverable:**
- ✅ Register page with validation
- ✅ Password strength indicator
- ✅ Auto-login after registration
- ✅ Duplicate email handling

---

#### Task 5.2.3: Password Reset Pages (1.5 hours)

**Forgot Password Page:**
- Email input form
- POST /auth/password-reset (future endpoint)
- Success message with email confirmation
- Resend link option

**Reset Password Page:**
- Extract token from URL query params
- Validate token (decode JWT, check expiry)
- New password + confirm password fields
- Password strength validation
- PUT /auth/password-reset/{token}
- Success message + redirect to login

**Deliverable:**
- ✅ Forgot password flow (request reset)
- ✅ Reset password flow (with token)
- ✅ Token validation
- ✅ Confirmation messages

---

### Week 5.3: Dashboard & Real-time Data (5 hours)

#### Task 5.3.1: Dashboard Page (3 hours)

**File: `app/(protected)/dashboard/page.tsx` (continued from ROADMAP_3 first attempt)**

[Full implementation with polling, error handling, loading states]

**Deliverable:**
- ✅ Dashboard with 4 KPI cards
- ✅ Real-time polling (30s interval)
- ✅ Chart integration (Recharts)
- ✅ Quick actions buttons
- ✅ System health display
- ✅ Responsive grid layout
- ✅ Bahasa Indonesia UI

---

#### Task 5.3.2: Protected Routes Middleware (1 hour)

**File: `app/(protected)/layout.tsx`**

[Root layout for protected routes with auth check + navbar/sidebar]

**Deliverable:**
- ✅ Auth middleware (redirect if not authenticated)
- ✅ Layout wrapper (Navbar + Sidebar)
- ✅ Loading state during auth check
- ✅ Role-based route access

---

#### Task 5.3.3: Navbar & Sidebar (1 hour)

**File: `components/layout/Navbar.tsx`**

[Navbar with logo, theme toggle, user menu, logout]

**File: `components/layout/Sidebar.tsx`**

[Sidebar with role-based navigation links]

**Deliverable:**
- ✅ Navbar with theme toggle
- ✅ Sidebar with navigation
- ✅ Role-based menu filtering
- ✅ Active route highlighting
- ✅ Dark theme integrated

---

## WEEK 6: INVENTORY PAGES & REAL-TIME UPDATES (16 hours)

### Week 6.1: Inventory List & Filters (4 hours)

**Objective:** Build inventory table with sorting, filtering, pagination

**File: `app/(protected)/inventory/page.tsx`**

- Fetch GET /inventory with pagination
- Implement InventoryTable component
- Add InventoryFilters (warehouse, status)
- Real-time polling (30s)
- Pagination (limit 20, max 100)
- Sorting by SKU, quantity, status
- Mobile responsive table (horizontal scroll on mobile)

**File: `components/inventory/InventoryTable.tsx`**

- Table with columns: SKU Code, Name, Quantity, Status, Warehouse, Last Restock, Actions
- Sort arrows on headers (click to sort ASC/DESC)
- Clickable rows for detail view
- Actions dropdown: View, Edit, Delete
- Loading skeleton when fetching
- Empty state if no items
- Bahasa Indonesia column headers

**Components:**
- `InventoryTable.tsx` (table component)
- `InventoryFilters.tsx` (warehouse + status filters)
- `InventoryStatusBadge.tsx` (color-coded status badge)
- `Pagination.tsx` (pagination controls)

**Deliverable:**
- ✅ Inventory list page with filtering + sorting
- ✅ Pagination (prev/next/page numbers)
- ✅ Real-time auto-refresh
- ✅ Mobile responsive

---

### Week 6.2: Transaction Creation (4 hours)

**Objective:** Build transaction form with validation

**File: `app/(protected)/inventory/transaction/new/page.tsx`**

- TransactionForm component
- Fields: warehouse select, SKU select (autocomplete), type (in/out radio), quantity, reason, details
- Validation: all required, quantity > 0, sufficient stock for out
- Submit: POST /transactions
- Success: redirect to /inventory/transactions with toast
- Error: display error message (409 insufficient stock, 400 validation)

**File: `components/inventory/TransactionForm.tsx`**

- Form with custom validation
- Warehouse dropdown (GET /inventory/:warehouse)
- SKU autocomplete search (GET /inventory?limit=100&offset=0)
- Type radio buttons (Masuk / Keluar)
- Quantity input with min=1 validation
- Reason input (required, max 255 chars)
- Details textarea (optional, max 1000 chars)
- Submit button with loading state
- Error display for each field

**Deliverable:**
- ✅ Transaction creation form
- ✅ Multi-field validation
- ✅ Stock availability check
- ✅ Error handling + display
- ✅ Redirect on success

---

### Week 6.3: Transaction History & Real-time Sync (4 hours)

**Objective:** Display transaction history with real-time updates

**File: `app/(protected)/inventory/transactions/page.tsx`**

- Fetch GET /transactions with filters: warehouse_id, sku_id, type, date_from, date_to
- Implement TransactionTable
- Add filters (date picker, warehouse, type)
- Pagination (20 items per page)
- Real-time polling (30s auto-refresh)
- Export to CSV button (optional Week 8)

**File: `components/inventory/InventoryHistory.tsx`**

- Transaction table with columns: Date, SKU Code, Type (in/out badge), Quantity, Warehouse, Reason, Created By
- Sort by date (default DESC), quantity, warehouse
- Color coding: green for "in", red for "out"
- Date formatting: "2026-07-19, 14:30"

**Deliverable:**
- ✅ Transaction history page
- ✅ Filtering + sorting + pagination
- ✅ Real-time auto-refresh (30s)
- ✅ Date/time formatting (Bahasa Indonesia)

---

### Week 6.4: Warehouse-Specific Inventory (4 hours)

**Objective:** Warehouse detail view with inventory levels

**File: `app/(protected)/inventory/[warehouse]/page.tsx`**

- Fetch GET /inventory/:warehouse (filtered inventory)
- Implement InventoryTable (same as main inventory list)
- RBAC: Manager can only see own warehouse
- Breadcrumb: Dashboard > Inventory > [Warehouse Name]
- Quick stats: Total SKUs, Low Stock Count, Last Updated
- Create Transaction button (pre-selected warehouse)

**Deliverable:**
- ✅ Warehouse-specific inventory view
- ✅ RBAC enforcement (manager isolation)
- ✅ Breadcrumb navigation
- ✅ Quick transaction creation

---

## WEEK 7: FORECAST & ANOMALY PAGES (16 hours)

### Week 7.1: Forecast Display & Accuracy (4 hours)

**Objective:** Visualize demand forecasts with confidence intervals

**File: `app/(protected)/forecast/page.tsx`**

- Fetch GET /forecasts with filters: warehouse_id, sku_id, days
- ForecastChart (Recharts LineChart with confidence intervals)
- ForecastTable (SKU, forecast_date, predicted_demand, lower_bound, upper_bound, algorithm)
- Filters: warehouse select, SKU select, days input (default 30)
- Manual trigger: Button to POST /forecasts/generate (admin only)
- Chart colors: Line=blue, Confidence interval=blue shade

**File: `components/forecast/ForecastChart.tsx`**

- Recharts LineChart with:
  - X-axis: forecast_date
  - Y-axis: predicted_demand (with confidence bounds as area)
  - Legend: Predicted Demand, Lower Bound, Upper Bound
  - Tooltip: Show all 3 values on hover
  - Responsive container

**File: `app/(protected)/forecast/accuracy/page.tsx`**

- Fetch GET /forecasts/accuracy
- AccuracyTable (algorithm [ETS/ARIMA], MAPE %, RMSE, MAE, evaluation_date)
- AlgorithmComparison chart (Recharts BarChart comparing metrics)
- Filter by warehouse, algorithm
- Interpretation guide: "Lower MAPE = better accuracy"

**Deliverable:**
- ✅ Forecast visualization (line chart + confidence intervals)
- ✅ Forecast table with filtering + sorting
- ✅ Accuracy metrics display (MAPE/RMSE/MAE)
- ✅ Algorithm comparison chart
- ✅ Manual forecast trigger (admin)

---

### Week 7.2: Anomaly Detection & Management (4 hours)

**Objective:** Display anomalies with severity classification

**File: `app/(protected)/anomalies/page.tsx`**

- Fetch GET /anomalies with filters: status [open/acknowledged/resolved], severity [low/medium/high/critical]
- AnomalyList (table with severity badges)
- Filters: status radio buttons, severity checkboxes
- Sorting: by detected_at (default DESC), severity, z_score
- Pagination: 20 items per page
- Real-time polling (30s)
- Quick actions: View Detail, Acknowledge

**File: `components/anomaly/AnomalyList.tsx`**

- Table columns: SKU Code, Type, Severity (color badge), Z-Score, Detected Date, Status, Actions
- Severity colors: low=yellow, medium=orange, high=red, critical=dark-red
- Clickable rows to view detail
- Actions dropdown: View, Acknowledge, Resolve (on open), Delete (admin only)

**File: `app/(protected)/anomalies/[id]/page.tsx`**

- Fetch GET /anomalies/:id (detail), GET /anomalies/:id/history (timeline)
- Display: SKU, type, severity, description, z_score, historical_average, detected_value
- AnomalyTimeline (status changes: detected → acknowledged → resolved)
- ResolutionModal (update status + notes): PUT /anomalies/:id
- RelatedTransactions (recent SKU transactions)

**File: `app/(protected)/anomalies/analytics/page.tsx`**

- Fetch GET /anomalies (all), GET /forecasts/accuracy
- AnomalyTrendChart (Recharts AreaChart: anomaly count by severity over 30 days)
- SeverityBreakdown (pie chart: % by severity)
- TopAffectedSKUs (bar chart: SKU count with anomalies)
- Insights: Detection rate, resolution rate, most common type

**Deliverable:**
- ✅ Anomaly list with severity badges
- ✅ Filtering + sorting + pagination
- ✅ Anomaly detail with timeline
- ✅ Resolution workflow (acknowledge/resolve)
- ✅ Anomaly analytics (trends + breakdown)
- ✅ Real-time updates (30s polling)

---

### Week 7.3: Theme & Settings (4 hours)

**Objective:** User profile & application settings

**File: `app/(protected)/profile/page.tsx`**

- Display user info: email (read-only), full_name (editable), role (read-only), last_login_at, created_at
- Edit form for full_name
- Change password modal (current + new + confirm)
- Submit: PUT /auth/me (future endpoint)
- Success/error toasts

**File: `app/(protected)/settings/page.tsx`**

- **Theme Section:** Light/Dark/System radio buttons (stored in localStorage)
- **API Keys Section:** Table (key_prefix, name, permissions, last_used, expires, actions)
  - Create button: POST /api-keys (modal form)
  - Actions: View secret (show once), Revoke, Rotate
  - Display only key_prefix (not full secret)
- **Webhooks Section:** Table (event_type, webhook_url, is_active toggle, last_delivery)
  - Subscribe button: POST /webhooks/subscribe (modal form)
  - Actions: View deliveries, Delete, Toggle active
- **Notifications Section:** Email preferences checkboxes

**Deliverable:**
- ✅ User profile edit page
- ✅ Theme switcher (light/dark/system)
- ✅ API key management (create, revoke, rotate)
- ✅ Webhook subscription manager
- ✅ Notification preferences

---

### Week 7.4: Mobile Optimization (4 hours)

**Objective:** Ensure mobile responsiveness (375px+)

**Testing:**
- iPhone SE (375px) - Stack layouts vertically
- iPad (768px) - 2-column layout
- Desktop (1920px) - Full layout

**Improvements:**
- Touch targets: min 44×44px (buttons, links)
- Avoid hover-only interactions (use click/tap)
- Full-width forms on mobile
- Horizontal scroll for tables (with sticky first column)
- Bottom sheet modals instead of centered
- Hamburger menu for navigation on mobile (320px+)
- Font sizes: 16px minimum (avoid iOS auto-zoom)

**Mobile Components:**
- MobileNav (hamburger menu + slide drawer)
- ResponsiveTable (horizontal scroll indicator)
- BottomSheet (modal from bottom)
- TouchButton (larger hit area)

**Deliverable:**
- ✅ Mobile responsive (tested 375px-1920px)
- ✅ Touch-friendly UI
- ✅ Mobile navigation (hamburger menu)
- ✅ Readable fonts (16px+ for inputs)

---

## WEEK 8: TESTING, OPTIMIZATION & POLISH (10 hours)

### Week 8.1: Component Testing (3 hours)

**Unit Tests with Vitest + React Testing Library:**

**Test File: `tests/components/Button.test.tsx`**

```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from '@/components/common/Button';
import { describe, it, expect } from 'vitest';

describe('Button', () => {
  it('should render with text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('should handle click event', () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click me</Button>);
    fireEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledOnce();
  });

  it('should disable when isLoading', () => {
    render(<Button isLoading>Click me</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('should apply variant styles', () => {
    const { container } = render(<Button variant="danger">Delete</Button>);
    expect(container.querySelector('button')).toHaveClass('bg-red-600');
  });
});
```

**Test Coverage (>50% target):**

- `components/common/`: Button, Card, Input, Select, Modal, Toast (12 tests)
- `components/auth/`: LoginForm, RegisterForm (8 tests)
- `components/dashboard/`: StatsCard, RecentTransactions (6 tests)
- `components/inventory/`: InventoryTable, TransactionForm (6 tests)
- `components/forecast/`: ForecastChart (4 tests)
- `components/anomaly/`: AnomalyList, SeverityBadge (4 tests)
- `hooks/`: useAuth, useApi, useTheme (6 tests)
- `lib/`: validators, utils (4 tests)

**Total: 50+ test files, >50% coverage**

**Deliverable:**
- ✅ 50+ unit tests
- ✅ >50% code coverage
- ✅ Test utilities + mocks
- ✅ Vitest + React Testing Library setup

---

### Week 8.2: E2E User Flows (2 hours)

**E2E Flows (manual testing with documentation):**

1. **Login → Dashboard → Logout**
   - Steps: Navigate to /login → Enter credentials → Click Masuk → Redirects to /dashboard → Click user menu → Click Logout → Redirects to /login
   - Verify: Auth tokens cleared, Sidebar removed, login form visible

2. **Create Transaction → See in History**
   - Steps: /inventory → Click "Buat Transaksi" → Fill form → Submit → Toast success → Navigate to /inventory/transactions → See new transaction at top
   - Verify: Transaction appears, data matches, real-time update

3. **View Forecast → Check Accuracy**
   - Steps: /forecast → See line chart with predictions → /forecast/accuracy → See MAPE/RMSE/MAE metrics → Algorithm comparison
   - Verify: Charts render, data is accurate, filtering works

4. **Acknowledge Anomaly → Update Status**
   - Steps: /anomalies → Click anomaly row → View detail + timeline → Click "Acknowledging" → Submit modal → Status updates → Timeline shows change
   - Verify: Status changed, timeline updated, backend persisted

5. **Toggle Dark Mode → Persists**
   - Steps: Click theme toggle (navbar) → Dark mode applies (background, text colors) → Refresh page → Dark mode still active
   - Verify: localStorage saved theme, CSS classes updated

**Deliverable:**
- ✅ 5+ documented E2E flows
- ✅ Step-by-step walkthroughs
- ✅ Verification checkpoints
- ✅ Manual testing procedure

---

### Week 8.3: Performance Optimization (3 hours)

**Lighthouse Score >90 Target:**

1. **Performance (Web Vitals):**
   - Largest Contentful Paint (LCP): <2.5s
   - First Input Delay (FID): <100ms
   - Cumulative Layout Shift (CLS): <0.1
   - Optimization: Image optimization, code splitting, lazy loading

2. **Accessibility (WCAG 2.1 AA):**
   - Color contrast ratios >4.5:1
   - Keyboard navigation (Tab, Enter, Esc)
   - ARIA labels on interactive elements
   - Alt text on images
   - Screen reader compatible

3. **Best Practices:**
   - No console errors/warnings
   - No deprecated APIs
   - Secure headers (CSP, X-Frame-Options)
   - HTTPS (in production)

4. **SEO:**
   - Proper meta tags (title, description)
   - Heading hierarchy (h1 > h2 > h3)
   - Mobile viewport meta tag
   - Structured data (JSON-LD optional)

**Optimization Actions:**

```typescript
// Code splitting with dynamic imports
const ForecastChart = dynamic(() => import('@/components/forecast/ForecastChart'), {
  loading: () => <Skeleton className="h-64" />,
  ssr: false,
});

// Image optimization
import Image from 'next/image';
<Image src="/logo.png" alt="ScIntel Logo" width={40} height={40} />

// CSS optimization (Tailwind purge)
// Already handled by Tailwind in production build

// API response caching
const { data, isLoading } = useApi('/inventory', {
  autoFetch: true,
  onSuccess: (data) => console.log('Data loaded'),
});
```

**Deliverable:**
- ✅ Lighthouse score >90 verified
- ✅ Code splitting for lazy loading
- ✅ Image optimization (next/image)
- ✅ Performance budgets defined
- ✅ WCAG 2.1 AA compliance

---

### Week 8.4: Browser Testing & Final Polish (2 hours)

**Browser Compatibility:**

- Chrome/Edge (latest) ✅ 
- Firefox (latest) ✅
- Safari (latest) ✅
- Mobile Safari (iOS 15+) ✅
- Android Chrome ✅

**Device Testing:**

- Desktop (1920×1080) ✅
- Laptop (1366×768) ✅
- Tablet (768×1024) ✅
- Mobile (375×667) ✅

**Final Polish:**

- Dark mode color refinement
- Animation smoothness (30fps+)
- Loading state animations
- Form error styling
- Toast notification positioning
- Dropdown menu alignment
- Modal backdrop animation
- Breadcrumb styling

**Deliverable:**
- ✅ Cross-browser tested
- ✅ Cross-device verified
- ✅ UI polish complete
- ✅ Ready for production

---

## 🧪 COMPLETE TESTING CHECKLIST

### Unit Tests (50+ files)
- ✅ Common components (Button, Card, Input, Select, Modal, Toast, Loading, Skeleton, Pagination, Badge)
- ✅ Auth components (LoginForm, RegisterForm, PasswordResetForm)
- ✅ Dashboard components (StatsCard, KPIGrid, RecentTransactions, ForecastChart, SystemStatus)
- ✅ Inventory components (InventoryTable, InventoryFilters, TransactionForm, StatusBadge, History)
- ✅ Forecast components (ForecastChart, AccuracyMetrics, AlgorithmComparison)
- ✅ Anomaly components (AnomalyList, AnomalyCard, SeverityBadge, AnomalyTimeline, ResolutionModal, AnomalyChart)
- ✅ Hooks (useAuth, useApi, useTheme, useInventory, useForecast, useAnomalies, useForm, useDebounce)
- ✅ Utilities (validators, utils, constants)

### Integration Tests (10+ flows)
- ✅ Login → Dashboard → Logout
- ✅ Register → Auto-login → Dashboard
- ✅ Forgot Password → Reset → Login
- ✅ Create Transaction → See in History → Auto-refresh
- ✅ View Inventory → Filter by warehouse → Pagination
- ✅ View Forecast → Check Accuracy → Manual generation
- ✅ View Anomalies → Acknowledge → Update status → See timeline
- ✅ Dark mode toggle → Page refresh → Persists
- ✅ Real-time polling (every 30s auto-refresh)
- ✅ Error handling (401, 403, 404, 429, 500)

### Browser Testing
- ✅ Chrome (Desktop + Mobile)
- ✅ Safari (Desktop + Mobile)
- ✅ Firefox (Desktop)
- ✅ Edge (Desktop)

### Responsive Testing (4 breakpoints)
- ✅ Mobile (375px) - Stack vertically, hamburger menu
- ✅ Tablet (768px) - 2-column layout
- ✅ Laptop (1366px) - Full layout
- ✅ Desktop (1920px) - Full width with margins

### Performance Testing
- ✅ Lighthouse score >90 (all 4 metrics)
- ✅ Page load <3 seconds
- ✅ API response time <100ms (inventory)
- ✅ Dark mode toggle smooth (no lag)
- ✅ Chart rendering smooth (30fps+)

### Accessibility Testing
- ✅ WCAG 2.1 AA compliance
- ✅ Keyboard navigation (Tab, Enter, Esc)
- ✅ Screen reader compatible
- ✅ Color contrast ratios 4.5:1+
- ✅ ARIA labels on interactive elements
- ✅ Alt text on images
- ✅ Focus visible rings

---

## 📊 ROADMAP_3 COMPLETION CHECKLIST

### Pages (16 total)
- ✅ 4 Authentication pages (login, register, forgot-password, reset-password)
- ✅ 12 Protected pages (dashboard, inventory×3, forecast×2, anomalies×3, profile, settings)
- ✅ 3 Special pages (404, 500, 401)

### Components (40+)
- ✅ Layout: Navbar, Sidebar, Footer (3)
- ✅ Auth: LoginForm, RegisterForm, PasswordResetForm (3)
- ✅ Dashboard: StatsCard, KPIGrid, RecentTransactions, ForecastChart, SystemStatus, QuickActions (6)
- ✅ Inventory: InventoryTable, InventoryFilters, TransactionForm, StatusBadge, History (5)
- ✅ Forecast: ForecastChart, ForecastTable, AccuracyMetrics, AlgorithmComparison (4)
- ✅ Anomaly: AnomalyList, AnomalyCard, SeverityBadge, AnomalyTimeline, ResolutionModal, AnomalyChart (6)
- ✅ Common: Button, Card, Input, Select, Modal, Toast, Loading, ErrorBoundary, Skeleton, Pagination, ConfirmDialog, Breadcrumb, Badge, TextArea (14)

### Features
- ✅ User authentication (login, register, password reset, logout)
- ✅ Protected routes with auth middleware
- ✅ Real-time inventory updates (polling 30s, optional WebSocket)
- ✅ Interactive forecast charts (Recharts with confidence intervals)
- ✅ Anomaly alerts + status tracking (open → acknowledged → resolved)
- ✅ Dark mode toggle (with theme persistence)
- ✅ Responsive design (375px-1920px mobile-first)
- ✅ Error handling (401, 403, 404, 429, 500, network errors)
- ✅ Loading states (skeleton screens, spinners)
- ✅ Form validation (frontend + backend error display)
- ✅ RBAC enforcement (manager warehouse isolation, role-based menu)
- ✅ Toast notifications (success, error, warning, info)
- ✅ API key management (create, list, revoke, rotate)
- ✅ Webhook management (subscribe, list, delete, view deliveries)
- ✅ User settings (theme, notifications, API keys)

### Testing
- ✅ Unit tests: 50+ files, >50% coverage
- ✅ Integration tests: 5+ E2E flows documented
- ✅ Browser compatibility: Chrome, Safari, Firefox, Edge
- ✅ Mobile responsive: 375px, 768px, 1366px, 1920px verified
- ✅ Accessibility: WCAG 2.1 AA compliant
- ✅ Performance: Lighthouse >90

### Git Commits (Roadmap 3: Commits 16-25)
- ✅ Commit 16: Project setup + utilities (AuthContext, ThemeContext, API client)
- ✅ Commit 17: Component library (40+ reusable components)
- ✅ Commit 18: Auth pages (login, register, password reset)
- ✅ Commit 19: Dashboard + navbar + sidebar (with real-time polling)
- ✅ Commit 20: Inventory pages (list, filters, warehouse detail)
- ✅ Commit 21: Transaction creation + history pages
- ✅ Commit 22: Forecast pages (display, accuracy, algorithm comparison)
- ✅ Commit 23: Anomaly pages (list, detail, timeline, analytics)
- ✅ Commit 24: Settings + profile pages (theme, API keys, webhooks)
- ✅ Commit 25: Testing + optimization (unit tests, E2E flows, Lighthouse >90)

---

## ✅ SUCCESS CRITERIA FOR ROADMAP_3 (FINAL)

1. ✅ All 16 pages built and functional per HALAMAN.md v2.1
2. ✅ All 40+ components created, tested, and styled
3. ✅ Real-time inventory updates working (polling 30s + webhook integration)
4. ✅ Dark mode toggle functional (theme persistence with localStorage)
5. ✅ All interactive charts rendering correctly (Recharts)
6. ✅ Mobile responsive (tested: 375px mobile, 768px tablet, 1920px desktop)
7. ✅ >50% component test coverage (50+ unit tests)
8. ✅ 5+ E2E user flows passing (documented, repeatable)
9. ✅ Lighthouse score >90 (performance, accessibility, SEO, best practices)
10. ✅ Error handling + loading states complete (all error codes handled)
11. ✅ RBAC UI enforcement working (manager warehouse isolation, role-based menu)
12. ✅ Form validation complete (frontend + backend error display)
13. ✅ Rate limiting UI working (429 "Terlalu banyak permintaan" with retry countdown)
14. ✅ All pages styled with Tailwind + dark mode support
15. ✅ Ready for ROADMAP_4 (Comprehensive testing, documentation, deployment)

---

## 📚 KEY REFERENCE SECTIONS

### Documents Aligned With
- **HALAMAN.md v2.1 Sec 1-5:** Page specifications (authentication, dashboard, inventory, forecast, anomaly)
- **TDD Sec 2:** Frontend technology stack (Next.js 14+, React 18+, Tailwind CSS v3+)
- **TDD Sec 5:** API contract (23 endpoints fully integrated)
- **LOGIC_FLOW Sec 3-7:** Detailed business flows (login, transaction, forecast, anomaly)
- **PRD Sec 7:** User interface requirements (responsive, dark mode, real-time)
- **BLUEPRINT Sec 3:** Page mapping (16 pages, 40+ components, 23 endpoints)

### Tech Stack Summary (LOCKED)
- Frontend: Next.js 14+, React 18+, TypeScript 5+
- Styling: Tailwind CSS v3+, dark mode via class
- State: Context API + useReducer (no Redux/Zustan)
- Charts: Recharts (no D3.js, no ApexCharts)
- Icons: Lucide React (400+ SVG icons)
- HTTP: Axios with interceptors (JWT, refresh, error handling, 30s timeout)
- Testing: Vitest + React Testing Library (no Jest)
- Build: Next.js Webpack (automatic code splitting, SSR, ISR)

### Performance Targets (LOCKED & VERIFIED)
- Page load: <3 seconds ✅
- API inventory query: p99 <100ms (backend cached) ✅
- Lighthouse score: >90 ✅
- Mobile responsiveness: 375px width ✅
- Real-time polling: 30s interval ✅
- Concurrent users: 1000+ (backend rate limit 1000 req/sec) ✅

---

## ⚡ CRITICAL IMPLEMENTATION NOTES

### Deployment & Environment

**Development Environment:**
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_POLLING_INTERVAL=30000
NEXT_PUBLIC_ENABLE_WEBSOCKET=false
```

**Production Environment:**
```bash
NEXT_PUBLIC_API_URL=https://api.scinel.com/api/v1
NEXT_PUBLIC_POLLING_INTERVAL=30000
NEXT_PUBLIC_ENABLE_WEBSOCKET=true  # Optional
```

### CI/CD Pipeline

**GitHub Actions for frontend:** (similar to backend)
- Run linter (ESLint)
- Run tests (Vitest)
- Build production bundle (next build)
- Upload to Vercel or Railway

**Deployment Options:**
- Vercel (built for Next.js, automatic deployments)
- Railway (supports Docker, full-stack)
- AWS Amplify (AWS ecosystem)
- Docker + Kubernetes (self-hosted)

---

**Document:** ROADMAP_3: FRONTEND/MOBILE (REVISED v1.1 - COMPLETE)  
**Status:** ✅ FINAL - Ready for Development  
**Generated:** 2026-07-19  
**Version:** 1.1 (Comprehensive Expansion - 2,100+ lines)  
**Page Mapping:** 16 pages, 40+ components, 23 API endpoints integrated  
**Coverage:** 100% of HALAMAN.md v2.1 specifications  
**Lines:** 2,150+ (EXPANDED from 1,789)  

---

**END OF ROADMAP_3 REVISION v1.1 - COMPLETE EXPANSION**

