================================================================
PRD (PRODUCT REQUIREMENT DOCUMENT) - SUPPLY CHAIN INTELLIGENCE PLATFORM v1.0.0
AI-Powered Demand Forecasting & Warehouse Operations Management System
================================================================

## 1. HEADER & IDENTITAS

| Item | Value |
|------|-------|
| **Nama Proyek** | Supply Chain Intelligence Platform (ScIntel) |
| **Versi Dokumen** | 1.0.0 (MVP - Portfolio Edition) |
| **Tanggal** | 2026-07-18 |
| **Penulis** | M. Arif Aulia (Backend Engineer / Product Manager) |
| **Status** | ✅ Ready for Development |
| **Target Rating** | 10/10 untuk Portfolio Recruiter Backend Golang |
| **Backend Platform** | Golang 1.21+ |
| **Database** | PostgreSQL 15+ |
| **Cache Layer** | Redis 7+ |
| **Message Queue** | Kafka (Phase 1.5+, optional) |
| **Timeline** | October 2026 (12 weeks / 3 bulan development) |

---

## 2. RINGKASAN EKSEKUTIF

### LATAR BELAKANG MASALAH

Warehouse & supply chain managers menghadapi **6 major pain points** dengan existing solutions:

1. **Inaccurate Demand Forecasting** → Stockouts atau overstock menghabiskan budget
2. **Reactive Inventory Management** → Keputusan manual = delays dan costly errors
3. **Fragmented Systems** → Manual data aggregation dari 3+ warehouse = time-wasting
4. **No Real-Time Visibility** → Can't respond to market changes quickly
5. **Anomaly Blind Spots** → Unusual patterns missed (fraud, system errors, demand spikes)
6. **Scalability Issues** → Existing solutions break at 5K+ SKUs atau 100+ warehouses

**Existing Solutions** (SAP, Oracle SCM, Odoo, Shopify):
- ❌ Expensive ($10K-100K+ per tahun)
- ❌ Overkill & kompleks untuk mid-market (slow implementation)
- ❌ Poor Indonesia localization
- ❌ Limited anomaly detection (mahal add-on)
- ❌ Rigid architecture (sulit customize & integrate)

### SOLUSI PRODUK

**Supply Chain Intelligence Platform (ScIntel)** = enterprise-grade backend system yang solve semua 6 gaps:

- **Backend**: Golang (high-performance, concurrent, enterprise-ready)
- **Database**: PostgreSQL 15+ (ACID transactions, full-text search, window functions)
- **Cache**: Redis 7+ (sub-millisecond latency)
- **Design**: Microservices-ready, modular, horizontal scaling
- **DevOps**: Docker, Kubernetes-ready, CI/CD automated

### VALUE PROPOSITION

ScIntel deliver 3 core capabilities yang solve supply chain problems:

✅ **Accurate Forecasting** — Time-series algorithms (Exponential Smoothing, ARIMA) predict demand dengan 90%+ accuracy

✅ **Real-Time Inventory** — Sub-second visibility across multi-warehouse: click 1x, total stock dari semua warehouse instant tau

✅ **Anomaly Detection** — Statistical alerts untuk unusual patterns (demand spike, system errors, suspicious activities)

**Unique Position**: Portfolio-grade backend yang showcase master-level Golang + PostgreSQL + system design thinking.

---

## 3. TUJUAN PRODUK

### PRIMARY GOAL (Portfolio)

Build **production-grade backend system** yang demonstrate enterprise-level engineering:

- Showcase **advanced Golang patterns** (concurrency, error handling, interface design)
- Demonstrate **PostgreSQL mastery** (query optimization, indexing, transactions)
- Implement **system design at scale** (microservices, caching, async processing)
- Achieve **portfolio rating 10/10** dari backend-focused recruiters & interviewers
- Create **reusable architecture patterns** buat interview talking points

### SECONDARY GOALS

- Deep Golang knowledge (goroutines, channels, context)
- PostgreSQL expertise (indexes, query planning, CTEs, window functions)
- System reliability (circuit breakers, retry logic, graceful degradation)
- Professional documentation (API specs, deployment guides, architecture decisions)
- Production-ready code (Docker, CI/CD, Kubernetes-compatible)
- Clean code practices (SOLID principles, design patterns)

---

## 4. RUANG LINGKUP (SCOPE V1.0 - BALANCED APPROACH)

### 🎯 DECISION: 80% Balanced (dari conversation AI sebelah)

**Focus**: Forecasting + Inventory (DEEP), Anomaly Detection (LITE)
**Trade-off**: Quality over quantity. Master 2 domains perfectly, 1 domain solid.

---

### IN SCOPE (MVP - WAJIB IMPLEMENTASIKAN)

#### 4.1 INFRASTRUCTURE & ARCHITECTURE

| Item | Detail | Priority |
|------|--------|----------|
| Golang Backend | 1.21+, net/http, standard library focus | P0 |
| PostgreSQL | 15+, ACID transactions, full-text search | P0 |
| Redis | 7+, caching, rate limiting, session mgmt | P0 |
| Docker & Compose | Local dev, production deployment | P0 |
| GitHub Actions | CI/CD: lint, test, build, push | P0 |
| Connection Pooling | Efficient resource management | P0 |
| Graceful Shutdown | Clean resource cleanup | P0 |
| Structured Logging | JSON logs, appropriate levels | P0 |
| Database Migrations | Version control, rollback capability | P0 |

#### 4.2 AUTHENTICATION & AUTHORIZATION

| Item | Detail | Priority |
|------|--------|----------|
| JWT Auth | Stateless, secure token generation | P0 |
| Password Hashing | bcrypt, proper salt handling | P0 |
| RBAC | 4 roles: admin, manager, analyst, viewer | P0 |
| Token Refresh | Access + refresh token pattern | P0 |
| API Key Auth | Third-party integrations | P1 |
| Audit Logging | Who did what, when, from where | P1 |
| Rate Limiting | Per user/API key abuse prevention | P1 |

#### 4.3 DATABASE SCHEMA & DESIGN

| Item | Detail | Priority |
|------|--------|----------|
| Inventory Table | warehouse_id, sku, quantity, last_updated | P0 |
| Transaction Log | in/out movements, timestamps, user tracking | P0 |
| Forecast Data | sku, predicted_demand, confidence_interval | P0 |
| Anomaly Records | detected anomalies, severity, resolution | P1 |
| User Management | users, roles, permissions | P0 |
| Audit Trail | Immutable event log of all changes | P1 |
| Webhook Subscriptions | Event subscriptions & delivery logs | P1 |
| Proper Indexing | Search performance optimization | P0 |
| Foreign Key Constraints | Data integrity enforcement | P0 |
| Partitioning Strategy | Historical data isolation (P1) | P1 |

#### 4.4 DEMAND FORECASTING ENGINE (DEPTH - PRIORITAS UTAMA)

| Item | Detail | Priority |
|------|--------|----------|
| Historical Data Aggregation | Collect 6-12 months inventory movements | P0 |
| Time-Series Analysis | Trend, seasonality, cyclical patterns | P0 |
| Algorithm: Exponential Smoothing (ETS) | Simple, double, triple exponential | P0 |
| Algorithm: ARIMA | AutoRegressive Integrated Moving Average | P0 |
| Accuracy Metrics | MAE, RMSE, MAPE for evaluation | P0 |
| Confidence Intervals | 95% range for decision-making | P0 |
| Scheduled Forecast Runs | Daily/weekly updates via cron | P0 |
| Batch Processing | Efficient bulk forecast generation | P0 |
| SKU-Level Forecasting | Individual product predictions | P0 |
| Category-Level Aggregation | Roll-up forecasts | P1 |
| Model Comparison | Side-by-side algorithm performance | P1 |
| Seasonal Decomposition | Break down trend + seasonality | P1 |

#### 4.5 REAL-TIME INVENTORY TRACKING (DEPTH - PRIORITAS UTAMA)

| Item | Detail | Priority |
|------|--------|----------|
| Stock In/Out Transactions | Logged with timestamp, user, reason | P0 |
| Multi-Warehouse Support | Track inventory across 100+ locations | P0 |
| Real-Time Stock Level API | <100ms response time target | P0 |
| Inventory Snapshots | End-of-day, weekly, monthly records | P0 |
| Stock Alerts | Low stock, overstock, expiring items | P0 |
| Movement History | Complete audit trail of changes | P0 |
| Physical Count Reconciliation | Actual vs system quantities | P1 |
| SKU Lifecycle Tracking | New, active, discontinued status | P1 |
| Bulk Import/Export | CSV support for batch operations | P1 |
| Stock Level Predictions | Forecast when stock will deplete | P1 |

#### 4.6 ANOMALY DETECTION SYSTEM (LITE - SOLID BUT NOT DEEP)

| Item | Detail | Priority |
|------|--------|----------|
| Statistical Outlier Detection | Z-score, IQR methods | P0 |
| Pattern Recognition | Unusual ordering, return patterns | P1 |
| Threshold-Based Alerts | Customizable per warehouse/SKU | P0 |
| Severity Levels | Low, medium, high, critical | P0 |
| Alert Notifications | Email, webhook, dashboard | P1 |
| Resolution Workflow | Acknowledge, investigate, resolve | P1 |
| Historical Anomaly Analysis | Trend over time | P2 |

#### 4.7 WEBHOOK SYSTEM (EVENT-DRIVEN ARCHITECTURE)

| Item | Detail | Priority |
|------|--------|----------|
| Event Subscriptions | Subscribe to: inventory, forecast, anomaly | P0 |
| Webhook Delivery | POST to external URLs with retry logic | P0 |
| Retry Mechanism | Exponential backoff (3-5 retries max) | P0 |
| Signature Verification | HMAC-SHA256 for security | P0 |
| Dead-Letter Queue | Failed events for investigation | P1 |
| Event History | Audit trail of all webhook calls | P1 |
| Idempotency Keys | Prevent duplicate event processing | P0 |

#### 4.8 TESTING STRATEGY

| Item | Detail | Priority |
|------|--------|----------|
| Unit Tests | >80% code coverage target | P0 |
| Integration Tests | Database, cache, API layers | P0 |
| Load Testing | Benchmark 1000+ concurrent users | P1 |
| Smoke Tests | Quick health checks | P0 |
| E2E Tests | Critical user journeys (sample) | P1 |

#### 4.9 PERFORMANCE & SCALABILITY

| Item | Detail | Target |
|------|--------|--------|
| Real-Time Stock Query Latency | p99 <100ms | P0 |
| Forecast Batch Latency | <5 seconds for 1000 SKUs | P0 |
| Concurrent Users | Support 1000+ without degradation | P0 |
| Request Rate Limit | 1000 req/sec per API key | P0 |
| Database Connection Pool | Min 10, Max 50 connections | P0 |

#### 4.10 DEPLOYMENT & DEVOPS

| Item | Detail | Priority |
|------|--------|----------|
| Docker Image | Multi-stage build, <100MB size | P0 |
| Docker Compose | Local development environment | P0 |
| Kubernetes Manifests | Deployment, service, ConfigMap, secrets | P1 |
| Health Checks | Readiness & liveness probes | P0 |
| Horizontal Scaling | Stateless design for scale-out | P1 |
| Config Management | Environment variables, config files | P0 |
| Blue-Green Deployment | Zero-downtime deployments | P2 |

---

### OUT OF SCOPE (FOR V1.0 - ROADMAP FASE 1.5+)

| Item | Reason | Target |
|------|--------|--------|
| Advanced Anomaly Detection | ML models (isolation forest, LSTM) | Fase 1.5 |
| Kafka Integration | Event streaming for high-scale | Fase 1.5 |
| GraphQL API | REST is sufficient for MVP | Fase 2 |
| Machine Learning | Prophet-style models, deep learning | Fase 2 |
| Mobile App UI | Backend-focused, no frontend | N/A |
| Payment Gateway | Out of scope for supply chain | N/A |
| International Markets | Indonesia-first, global later | Fase 3 |
| Advanced Analytics Dashboard | Backend API, frontend separate | Fase 2 |
| Real-Time WebSocket | REST + webhooks sufficient | Fase 1.5 |
| Multi-Tenancy | Single tenant MVP, multi-tenant Fase 2 | Fase 2 |

---

## 5. USER PERSONAS

### 5.1 WAREHOUSE MANAGER (Primary)

| Attribute | Value |
|-----------|-------|
| **Name** | Bambang (age 35-45, 10+ years experience) |
| **Role** | Oversee daily warehouse operations, manage staff, handle customer orders |
| **Pain Points** | Never know actual stock levels, manual checking 3 warehouses = time-wasting, frequent stockouts = customer complaints |
| **Goals** | Real-time stock visibility, quick decision-making, reduce manual errors |
| **Usage Pattern** | Daily (morning briefing, throughout day), check inventory status, respond to alerts |

### 5.2 SUPPLY CHAIN ANALYST (Secondary)

| Attribute | Value |
|-----------|-------|
| **Name** | Siti (age 28-35, data-driven mindset) |
| **Role** | Analyze trends, plan orders, optimize inventory levels |
| **Pain Points** | Can't predict demand accurately, no trend analysis tools, guessing game leads to overstock/understock |
| **Goals** | Accurate forecasts, trend insights, data-driven recommendations |
| **Usage Pattern** | Weekly (planning cycles), view forecasts, run simulations, export reports |

### 5.3 OPERATIONS MANAGER (Tertiary)

| Attribute | Value |
|-----------|-------|
| **Name** | Rudi (age 40-50, executive perspective) |
| **Role** | KPIs, performance monitoring, strategic decisions |
| **Pain Points** | No visibility into system health, can't detect anomalies (fraud, errors), reactive instead of proactive |
| **Goals** | Automated alerts, anomaly detection, performance dashboards |
| **Usage Pattern** | Daily (morning briefing), view KPIs, investigate alerts |

---

## 6. USER JOURNEYS & USE CASES

### USE CASE 1: Real-Time Stock Check (Warehouse Manager)

**Flow:**
1. Bambang logs into ScIntel dashboard via JWT
2. Clicks "Inventory Status"
3. API returns total stock across all 3 warehouses in <100ms
4. Sees: Nike = 395 (Jakarta 150 + Surabaya 45 + Bandung 200), Adidas = 250
5. **Result**: Quick decision without manual calls to 3 warehouses

**Alternate Path (Error):**
- If warehouse is down: System still returns cached data (Redis), notifies Bambang of data freshness
- Graceful degradation: Better stale data than no data

---

### USE CASE 2: Forecast Planning (Supply Chain Analyst)

**Flow:**
1. Siti logs in, selects "Demand Forecast"
2. Chooses SKU: "Nike-AIR-MAX", Warehouse: "Jakarta", Horizon: "30 days"
3. System returns:
   - Predicted daily demand for next 30 days
   - 95% confidence interval (upper/lower bounds)
   - Accuracy metrics (MAPE: 8.2%)
4. Siti sees: "Nike-AIR-MAX demand will peak Saturday (150 units), plan stock accordingly"
5. **Result**: Data-driven ordering, prevents stockout

**Technical Detail:**
- System uses 3 algorithms (ETS, ARIMA, Holt-Winters), compares accuracy, recommends best
- Forecast runs daily via cron, pre-computed for speed

---

### USE CASE 3: Anomaly Alert (Operations Manager)

**Flow:**
1. Nike stock in Jakarta suddenly drops 100 units in 1 hour (unusual pattern)
2. System detects anomaly (Z-score > 3 sigma) at 2am
3. Sends webhook event to Rudi's dashboard (real-time)
4. Also sends email alert: "ANOMALY DETECTED: Nike stock dropped 100 units, investigate"
5. Rudi clicks link, sees:
   - Time of anomaly
   - Previous pattern (normal ~5 units/hour)
   - Possible causes (system error? employee theft? demand spike?)
6. **Result**: Quick investigation, prevent losses

---

## 7. FUNCTIONAL REQUIREMENTS (FR)

| ID | Deskripsi | Priority | Owner |
|----|-----------|-----------| -----|
| **FR-01** | User dapat register dengan email, validasi email (OTP) | P0 | Backend Auth |
| **FR-02** | User dapat login dengan JWT token, refresh token mechanism | P0 | Backend Auth |
| **FR-03** | User dapat logout (token invalidation via blacklist) | P0 | Backend Auth |
| **FR-04** | Admin dapat assign roles ke user (RBAC: admin, manager, analyst, viewer) | P0 | Backend Auth |
| **FR-05** | System enforce permission checks sesuai role | P0 | Backend Auth |
| **FR-06** | User dapat upload inventory via CSV (batch) | P1 | Backend Inventory |
| **FR-07** | System dapat create inventory in/out transaction | P0 | Backend Inventory |
| **FR-08** | Real-time stock level query (<100ms) | P0 | Backend Inventory |
| **FR-09** | System return total inventory across all warehouses | P0 | Backend Inventory |
| **FR-10** | Inventory alerts (low stock, overstock) triggered automatically | P0 | Backend Inventory |
| **FR-11** | Movement history audit trail immutable & queryable | P0 | Backend Inventory |
| **FR-12** | Forecast engine compute demand prediction (30 days) | P0 | Backend Forecasting |
| **FR-13** | Forecast include confidence intervals (95%) | P0 | Backend Forecasting |
| **FR-14** | Forecast accuracy metrics displayed (MAE, RMSE, MAPE) | P0 | Backend Forecasting |
| **FR-15** | Forecast run automatically daily via cron | P0 | Backend Forecasting |
| **FR-16** | Anomaly detection triggered on unusual patterns | P0 | Backend Anomaly |
| **FR-17** | Anomaly severity levels assigned (low/medium/high/critical) | P0 | Backend Anomaly |
| **FR-18** | Webhook events published on: inventory change, forecast update, anomaly | P0 | Backend Webhooks |
| **FR-19** | Webhook delivery with retry logic (exponential backoff) | P0 | Backend Webhooks |
| **FR-20** | Webhook signature verification (HMAC-SHA256) | P0 | Backend Webhooks |
| **FR-21** | API rate limiting per user (default 1000 req/sec) | P1 | Backend API |
| **FR-22** | API versioning (v1, v2, etc) | P1 | Backend API |
| **FR-23** | API error responses standardized (with request_id) | P0 | Backend API |
| **FR-24** | Database migrations versioned & rollback-able | P0 | Backend Infrastructure |
| **FR-25** | System graceful shutdown (cleanup resources in 30 sec) | P0 | Backend Infrastructure |
| **FR-26** | Structured JSON logging (request_id, user_id, latency) | P0 | Backend Infrastructure |
| **FR-27** | Panic recovery & error tracking (Sentry/DataDog) | P1 | Backend Infrastructure |
| **FR-28** | Health check endpoint for Kubernetes | P0 | Backend Infrastructure |

---

## 8. NON-FUNCTIONAL REQUIREMENTS (NFR)

| Category | Requirement | Target |
|----------|-------------|--------|
| **PERFORMANCE** | Real-time stock query latency (p99) | <100ms |
| **PERFORMANCE** | Forecast batch computation (1000 SKUs) | <5 seconds |
| **PERFORMANCE** | Database connection pool efficiency | Min 10, Max 50 |
| **SCALABILITY** | Concurrent users support | 1000+ without degradation |
| **SCALABILITY** | Request throughput | 1000+ req/sec |
| **SCALABILITY** | SKU count supported | 10K+ SKUs |
| **SCALABILITY** | Warehouse count supported | 100+ warehouses |
| **AVAILABILITY** | Uptime target (SLA) | 99% (2 nines) |
| **AVAILABILITY** | Graceful shutdown time | <30 seconds |
| **RELIABILITY** | Database backup frequency | Daily + transaction logs |
| **RELIABILITY** | Webhook retry max attempts | 5 attempts max |
| **RELIABILITY** | Dead-letter queue for failed events | Indefinite retention |
| **SECURITY** | Password hashing | bcrypt (salt rounds: 12) |
| **SECURITY** | Token expiry (access token) | 1 hour |
| **SECURITY** | Token expiry (refresh token) | 7 days |
| **SECURITY** | SQL injection prevention | Parameterized queries only |
| **SECURITY** | Rate limiting per API key | 1000 req/sec |
| **SECURITY** | Audit logging (PII excluded) | All user actions |
| **MAINTAINABILITY** | Code coverage target | >80% |
| **MAINTAINABILITY** | Documentation | README, API specs, ADRs, deployment guides |
| **MAINTAINABILITY** | Code modularity | Service-oriented structure |
| **USABILITY** | API response format consistency | JSON standard |
| **USABILITY** | Error message clarity | Actionable, not generic |

---

## 9. BATASAN & ASUMSI

### CONSTRAINTS (Keterbatasan)

| Constraint | Detail |
|-----------|--------|
| **Timeline** | 12 weeks (3 months) untuk MVP, tidak flexible |
| **Team** | Solo developer (Aulia), no external resources |
| **Browser/Client Support** | API-first, no browser requirement (REST via Postman/curl) |
| **Payment Gateway** | Out of scope, warehouse management only |
| **Real-Time WebSocket** | REST + webhooks sufficient, no WebSocket in V1 |
| **Multi-Tenancy** | Single tenant MVP, multi-tenant Fase 2 |
| **Data Retention** | 2 years historical data minimum, beyond = archive |
| **Regional Compliance** | Indonesia-focused, may not comply with GDPR/CCPA |

### ASSUMPTIONS

| Assumption | Rationale |
|-----------|-----------|
| **Database Availability** | PostgreSQL always running (or quick recovery) |
| **Network Connectivity** | Stable internet (webhooks, external APIs) |
| **User Skills** | API consumers can read documentation & use REST clients |
| **Data Integrity** | Warehouse staff accurate in inventory input (garbage in = garbage out) |
| **Forecasting Data** | 6-12 months historical data minimum untuk accuracy |
| **Seasonality** | Business has clear seasonal patterns (analyzable via ETS/ARIMA) |
| **No Extreme Events** | No COVID-like disruptions affecting demand prediction |

---

## 10. KRITERIA SUKSES & METRICS

### KEBERHASILAN PRODUK

| Metric | Target | Measurement |
|--------|--------|-------------|
| **System Uptime** | 99% (2 nines) | Monitoring via Kubernetes health checks |
| **Stock Query Latency (p99)** | <100ms | Load test results (1000 concurrent users) |
| **Forecast Accuracy (MAPE)** | <15% | Historical backtest on real data |
| **Code Coverage** | >80% | Jest/Go coverage reports |
| **API Documentation** | 100% endpoints documented | OpenAPI/Swagger spec |
| **Test Pass Rate** | 100% (CI/CD) | GitHub Actions pipeline |

### PORTFOLIO VALUE

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Recruiter Rating** | 9-10/10 | Feedback from 5+ backend engineers |
| **Interview Readiness** | Explain every design decision | Mock interview pass |
| **GitHub Stars** | 50+ (optional) | Community interest |
| **Blog Posts** | 2-3 technical writeups | Document journey & learnings |

---

## 11. DEPENDENSI & INTEGRASI EKSTERNAL

### EXTERNAL SERVICES

| Service | Purpose | Optional? | Reason |
|---------|---------|-----------|--------|
| PostgreSQL | Primary database | No (P0) | Core requirement |
| Redis | Caching layer | No (P0) | Performance critical |
| GitHub Actions | CI/CD pipeline | No (P0) | Automated testing |
| Docker Hub | Image registry | Yes (P1) | Container distribution |
| Email Service | Alert notifications | Yes (P1) | Can use SMTP fallback |
| Kafka | Event streaming | Yes (P2) | Phase 1.5+, for scale |

### EXTERNAL APIS

| API | Purpose | Priority |
|-----|---------|----------|
| Timezone API | Localize timestamps to region | P2 |
| Currency Exchange | Multi-currency support | P2 |

### LIBRARIES/FRAMEWORKS

| Technology | Usage | Rationale |
|-----------|-------|-----------|
| **Golang stdlib** | HTTP server, JSON, database/sql | Standard, no external deps |
| **pgx** | PostgreSQL driver | Best-in-class performance |
| **redis** | Redis client | Reliable, maintained |
| **gorilla/mux** | HTTP routing (optional) | Or use stdlib ServeMux |
| **jwt-go** | JWT token handling | Industry standard |
| **testify** | Unit testing assertions | Better error messages |
| **Docker** | Containerization | DevOps standard |

---

## 12. RISIKO TINGKAT TINGGI

### TECHNICAL RISKS

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|-----------|
| **Database Performance Degradation** | Queries slow down at 10K+ SKUs | Medium | Query optimization, indexing strategy, EXPLAIN ANALYZE |
| **Forecast Inaccuracy** | Predictions wrong, wrong ordering decisions | Medium | Test 3+ algorithms, use accuracy metrics to pick best, backtest on real data |
| **Webhook Delivery Failures** | External systems don't receive events | Low | Retry logic (exponential backoff), dead-letter queue, audit trail |
| **Security Breach (SQL Injection)** | Data leak, system compromise | Very Low | Parameterized queries only, input validation, security review |
| **Concurrent Request Spike** | System overload, crashes | Medium | Load testing, rate limiting, horizontal scaling ready |

### TIMELINE RISKS

| Risk | Impact | Mitigation |
|-----|--------|-----------|
| **Underestimated Complexity** | Miss deadline | Use "80% Balanced" scope (Forecasting deep, Anomaly lite) |
| **Scope Creep** | 12 weeks not enough | Freeze scope mid-development, defer to Fase 1.5 |
| **External Dependencies Delay** | Wait for PostgreSQL/Redis setup | Use Docker Compose for quick local setup |

### BUSINESS RISKS

| Risk | Mitigation |
|-----|-----------|
| **Portfolio Rejection** | Interview-proof design decisions, document tradeoffs clearly |
| **Algorithm Complexity** | Start with simple ETS, add ARIMA if time permits |
| **Forecasting False Positives** | Threshold tuning, manual review of high-severity anomalies |

---

## 13. DECISION LOG

### DECISION 1: Scope Approach (80% Balanced)

**Question**: Should we build forecasting + inventory (deep) + anomaly (lite), OR just forecasting + inventory?

**Options**:
- **100% Aggressive**: Everything (forecasting, inventory, anomaly detection deep). Risk: too ambitious, tight deadline.
- **80% Balanced**: Forecasting + inventory DEEP, anomaly LITE. Best portfolio value.
- **70% Safe**: Just forecasting + inventory, polished to perfection. Conservative.

**Decision**: **80% Balanced** ✅

**Rationale**: 
- Forecasting + Inventory showcase both algorithms AND database mastery
- Anomaly detection lite (basic Z-score) still impressive without over-commitment
- 12 weeks realistic for this scope
- Better to master 2 domains deeply than 3 domains shallowly

**Trade-off**: Advanced ML anomaly detection (LSTM, isolation forest) deferred to Fase 1.5

---

### DECISION 2: Database Complexity Level

**Question**: How deep should database optimization go?

**Options**:
- **Go Deep**: Advanced indexing, partitioning, materialized views, CTEs, window functions. 15% more complex.
- **Practical Level**: Good schema, solid indexes, standard queries. Sufficient for MVP.

**Decision**: **Practical Level + selective advanced features** ✅

**What this means**:
- ✅ Proper indexing strategy (B-tree on warehouse_id, sku, timestamp)
- ✅ Window functions for ranking (e.g., "top 10 moving SKUs")
- ✅ CTEs for complex aggregations (e.g., multi-warehouse summaries)
- ❌ Partitioning: Defer to Fase 1.5 (when 2+ years of historical data)
- ❌ Materialized views: Defer to Fase 1.5 (when analytic queries get complex)

**Rationale**: Showcase PostgreSQL strength without over-engineering for MVP

---

### DECISION 3: Forecasting Algorithm Depth

**Question**: Which algorithms to implement?

**Options**:
- **All 3 (Aggressive)**: ETS + ARIMA + Holt-Winters. ~20 hours math learning. Risky.
- **ETS + Simple Exponential Smoothing (Balanced)**: Sufficient, proven.
- **Use External Library (Safe)**: Less impressive, but fast.

**Decision**: **ETS + ARIMA (Balanced)** ✅

**What this means**:
- ✅ Exponential Smoothing (ETS): Simple, fast, works for most cases
- ✅ ARIMA: More sophisticated, better for complex patterns
- ❌ Holt-Winters: Defer to Fase 1.5 (low ROI, similar to ETS for most use cases)
- ✅ Compare both algorithms, pick best based on accuracy (MAPE), recommend to user

**Rationale**:
- 2 algorithms > 1 (shows comparative thinking)
- Both algorithms well-documented (easier to learn)
- Real code implementation (not library wrapper) = better portfolio value

---

### DECISION 4: Frontend

**Question**: Should we build a dashboard UI?

**Answer**: **NO** ✅

**Rationale**:
- MVP is backend-focused API
- Recruiters care about backend mastery, not fancy UI
- Postman collection + curl examples sufficient for testing
- Fase 2: Separate frontend project (Next.js/React)

---

### DECISION 5: Load Testing & Performance Benchmarks

**Question**: Should we include load testing proof?

**Answer**: **YES** ✅

**What this means**:
- Load test with 1000 concurrent users using Locust/k6
- Capture metrics: p50, p95, p99 latency, throughput
- Document results in README + GitHub
- Impressive benchmark = clear portfolio differentiator

**Timeline**: Week 11-12 (before launch)

---

## 14. TECHNICAL RECOMMENDATIONS (QUICK REF)

### ARCHITECTURE PATTERN

```
┌─────────────────────────────────────────────────┐
│                   API Gateway / Load Balancer    │
└────────────────┬──────────────────────────────────┘
                 │
         ┌───────┴───────┐
         │               │
    ┌────▼────┐     ┌────▼────┐
    │ Backend  │     │ Backend  │  (Horizontal Scale)
    │ Service  │     │ Service  │
    │ (Golang) │     │ (Golang) │
    └────┬────┘     └────┬────┘
         │               │
    ┌────┴───────────────┴────┐
    │                         │
 ┌──▼──┐              ┌──────▼──┐
 │ PostgreSQL         │  Redis  │
 │ (Primary DB)       │ (Cache) │
 └──────┘             └─────────┘
```

**Key Design Principles**:
- **Stateless** services → easy horizontal scaling
- **Database** is single source of truth
- **Redis** for performance (inventory cache, forecast cache)
- **Webhooks** for event-driven notification
- **Cron jobs** for scheduled forecasts

---

### DATABASE SCHEMA (SIMPLIFIED)

```sql
-- Users & Auth
CREATE TABLE users (id UUID, email, password_hash, role, created_at);
CREATE TABLE audit_log (id UUID, user_id, action, resource, timestamp);

-- Inventory
CREATE TABLE warehouses (id UUID, name, location);
CREATE TABLE inventory (warehouse_id UUID, sku TEXT, quantity INT, last_updated);
CREATE TABLE transactions (id UUID, warehouse_id, sku, type, quantity, reason, timestamp, user_id);

-- Forecasting
CREATE TABLE forecasts (id UUID, sku, warehouse_id, date, predicted_demand, confidence_interval, algorithm, generated_at);
CREATE TABLE forecast_accuracy (id UUID, sku, mape, rmse, mae, evaluation_date);

-- Anomalies
CREATE TABLE anomalies (id UUID, warehouse_id, sku, type, severity, detected_at, resolution_status);

-- Webhooks
CREATE TABLE webhook_subscriptions (id UUID, user_id, event_type, url, created_at);
CREATE TABLE webhook_deliveries (id UUID, subscription_id, event_id, status, retries, last_retry);

-- Indexes (CRITICAL)
CREATE INDEX idx_inventory_warehouse_sku ON inventory(warehouse_id, sku);
CREATE INDEX idx_transactions_warehouse_sku ON transactions(warehouse_id, sku, timestamp);
CREATE INDEX idx_forecasts_sku_date ON forecasts(sku, date);
CREATE INDEX idx_audit_log_user_timestamp ON audit_log(user_id, timestamp);
```

---

### API ENDPOINTS (QUICK REF - See TDD for detailed contracts)

```
AUTH:
  POST   /api/v1/auth/register          → Create user account
  POST   /api/v1/auth/login             → Get JWT token
  POST   /api/v1/auth/refresh           → Refresh access token
  POST   /api/v1/auth/logout            → Invalidate token

INVENTORY:
  GET    /api/v1/inventory              → Total stock across all warehouses
  GET    /api/v1/inventory/:warehouse   → Stock in specific warehouse
  POST   /api/v1/transactions           → Record stock in/out
  GET    /api/v1/transactions           → Query transaction history
  POST   /api/v1/inventory/import       → Bulk import CSV

FORECASTING:
  GET    /api/v1/forecasts              → Get forecast for SKU (param: sku, horizon)
  POST   /api/v1/forecasts/generate     → Manually trigger forecast
  GET    /api/v1/forecasts/accuracy     → Model accuracy metrics

ANOMALIES:
  GET    /api/v1/anomalies              → List detected anomalies
  GET    /api/v1/anomalies/:id          → Anomaly details
  PUT    /api/v1/anomalies/:id          → Mark as resolved

WEBHOOKS:
  POST   /api/v1/webhooks/subscribe     → Subscribe to events
  GET    /api/v1/webhooks/subscriptions → List subscriptions
  DELETE /api/v1/webhooks/:id           → Unsubscribe

HEALTH:
  GET    /health                        → Liveness probe
  GET    /ready                         → Readiness probe
```

---

## 15. GLOSARIUM

| Term | Definisi | Context |
|------|----------|---------|
| **ACID** | Atomicity, Consistency, Isolation, Durability (database transaction properties) | PostgreSQL ensures ACID for data integrity |
| **ANOMALY** | Unusual/unexpected pattern in data (e.g., demand spike, stock drop) | Detection system flags anomalies for investigation |
| **ARIMA** | AutoRegressive Integrated Moving Average (time-series forecasting algorithm) | Used alongside ETS for demand prediction |
| **AUDIT TRAIL** | Immutable log of all system changes (who, what, when, where) | Compliance & debugging |
| **CACHE** | Fast, temporary storage (Redis) of frequently accessed data | Reduces database load, improves latency |
| **CONCURRENCY** | Goroutines handling multiple tasks simultaneously (Golang) | Efficient resource usage |
| **CONTEXT** | Golang mechanism for cancellation, timeouts, value passing | Clean request lifecycle management |
| **CONNECTION POOL** | Reusing database connections instead of creating new ones | Performance optimization |
| **CTE** | Common Table Expressions (SQL WITH clauses) | Complex query building blocks |
| **EXPONENTIAL SMOOTHING** | Time-series method giving more weight to recent data (ETS) | Forecasting algorithm |
| **GOROUTINE** | Lightweight thread in Golang for concurrent execution | Async background tasks |
| **GRACEFUL SHUTDOWN** | Properly cleaning up resources before application stops | Prevent data loss, connection leaks |
| **HMAC** | Hash-based Message Authentication Code (webhook signature verification) | Security |
| **JWT** | JSON Web Token (stateless authentication) | Secure API token |
| **KAFKA** | Message queue for event streaming (optional, Fase 1.5+) | Async event processing at scale |
| **LATENCY** | Time taken to complete request (milliseconds) | Performance metric |
| **MATERIALIZED VIEW** | Pre-computed database view for performance | Analytics queries (defer to Fase 1.5) |
| **MICROSERVICES** | Small, independent services communicating via APIs | Architecture style |
| **MAPE** | Mean Absolute Percentage Error (forecast accuracy metric) | <15% = good forecast |
| **OAUTH2** | Industry-standard authentication framework | Third-party integrations (Fase 1.5) |
| **PANIC** | Golang runtime error (must be recovered gracefully) | Error handling |
| **P99 LATENCY** | 99th percentile response time (max latency most users see) | <100ms target |
| **PARTITIONING** | Dividing large tables into smaller manageable parts | Database optimization (Fase 1.5) |
| **RATE LIMITING** | Controlling request frequency to prevent abuse | 1000 req/sec per API key |
| **RBAC** | Role-Based Access Control (permission management) | 4 roles: admin, manager, analyst, viewer |
| **SKU** | Stock Keeping Unit (unique product identifier) | E.g., "Nike-AIR-MAX-001" |
| **WEBHOOK** | Event-driven callback (external service notified of changes) | Inventory change → trigger forecast |

---

## 16. TIMELINE & PHASES

### PHASE 1: MVP (12 WEEKS)

| Week | Focus | Deliverables |
|------|-------|--------------|
| 1-2 | Foundation, API design, auth | Project setup, database schema, JWT implementation |
| 3-4 | Inventory core | CRUD endpoints, transaction logging, caching |
| 5-6 | Real-time inventory | Stock query optimization, multi-warehouse queries |
| 7-8 | Forecasting engine | ETS + ARIMA algorithms, scheduled runs, accuracy metrics |
| 9 | Webhooks & anomaly basics | Event subscriptions, basic anomaly detection |
| 10 | Testing & optimization | Unit/integration tests (80%+), load testing |
| 11 | Documentation & DevOps | API specs, README, Docker, CI/CD |
| 12 | Polish, launch | Bug fixes, final optimization, go live |

### PHASE 1.5: POST-MVP ENHANCEMENTS (Optional, Future)

- Holt-Winters algorithm
- Advanced anomaly detection (ML models)
- Kafka integration
- Materialized views
- Table partitioning
- Multi-tenancy

### PHASE 2: ANALYTICS DASHBOARD (Future)

- Next.js frontend
- Real-time charts (forecast, inventory trends)
- Anomaly investigation dashboard
- Report generation

---

## 17. APPENDIX: EXAMPLE API RESPONSES

### 17.1 CREATE TRANSACTION

**Request:**
```bash
POST /api/v1/transactions
Authorization: Bearer eyJhbGc...
Content-Type: application/json

{
  "warehouse_id": "wh-jakarta-001",
  "sku": "SKU-NIKE-001",
  "quantity": 50,
  "type": "in",
  "reason": "Supplier delivery #PO-2026-001"
}
```

**Response (201 Created):**
```json
{
  "transaction_id": "txn-abc123",
  "warehouse_id": "wh-jakarta-001",
  "sku": "SKU-NIKE-001",
  "quantity_delta": 50,
  "previous_stock_level": 300,
  "new_stock_level": 350,
  "timestamp": "2026-07-18T14:30:00Z",
  "created_by": "user-bambang-001",
  "webhook_events_triggered": [
    "inventory.changed",
    "forecast.trigger.scheduled"
  ]
}
```

---

### 17.2 GET INVENTORY TOTAL

**Request:**
```bash
GET /api/v1/inventory?sku=SKU-NIKE-001
Authorization: Bearer eyJhbGc...
```

**Response (200 OK):**
```json
{
  "sku": "SKU-NIKE-001",
  "total_stock": 395,
  "warehouse_breakdown": [
    {
      "warehouse_id": "wh-jakarta-001",
      "warehouse_name": "Jakarta Main",
      "quantity": 150
    },
    {
      "warehouse_id": "wh-surabaya-001",
      "warehouse_name": "Surabaya Branch",
      "quantity": 45
    },
    {
      "warehouse_id": "wh-bandung-001",
      "warehouse_name": "Bandung Branch",
      "quantity": 200
    }
  ],
  "last_updated": "2026-07-18T14:35:00Z",
  "timestamp": "2026-07-18T14:35:02Z"
}
```

---

### 17.3 GET FORECAST

**Request:**
```bash
GET /api/v1/forecasts?sku=SKU-NIKE-001&warehouse_id=wh-jakarta-001&horizon=30
Authorization: Bearer eyJhbGc...
```

**Response (200 OK):**
```json
{
  "sku": "SKU-NIKE-001",
  "warehouse_id": "wh-jakarta-001",
  "generated_at": "2026-07-18T02:00:00Z",
  "forecast_horizon": 30,
  "total_predicted_demand": 4500,
  "forecasts": [
    {
      "date": "2026-07-19",
      "predicted_demand": 150,
      "lower_bound": 120,
      "upper_bound": 180,
      "confidence": 0.95,
      "algorithm": "exponential_smoothing",
      "trend": "stable",
      "seasonality_factor": 1.05
    },
    {
      "date": "2026-07-20",
      "predicted_demand": 155,
      "lower_bound": 122,
      "upper_bound": 188,
      "confidence": 0.95,
      "algorithm": "exponential_smoothing",
      "trend": "stable",
      "seasonality_factor": 1.05
    }
  ],
  "model_info": {
    "algorithm_used": "exponential_smoothing",
    "accuracy_metrics": {
      "mape": 0.082,
      "rmse": 12.5,
      "mae": 10.2
    },
    "historical_data_points": 180,
    "last_model_update": "2026-07-18T02:00:00Z"
  }
}
```

---

### 17.4 ANOMALY ALERT (WEBHOOK PAYLOAD)

**Webhook Event (POST to subscriber URL):**
```json
{
  "event_id": "evt-xyz789",
  "event_type": "anomaly.detected",
  "timestamp": "2026-07-18T02:30:00Z",
  "data": {
    "anomaly_id": "anom-001",
    "warehouse_id": "wh-jakarta-001",
    "sku": "SKU-NIKE-001",
    "type": "unusual_stock_drop",
    "severity": "critical",
    "description": "Stock dropped 100 units in 1 hour (unusual pattern detected)",
    "historical_average": "5 units/hour",
    "detected_pattern": "100 units in 60 minutes",
    "z_score": 4.2,
    "investigation_link": "https://scintel.app/anomalies/anom-001"
  },
  "signature": "sha256=abcd1234..."
}
```

---

### 17.5 ERROR RESPONSE

**Request:**
```bash
GET /api/v1/inventory/SKU-INVALID
Authorization: Bearer eyJhbGc...
```

**Response (404 Not Found):**
```json
{
  "code": 404,
  "message": "SKU not found",
  "details": "SKU 'SKU-INVALID' does not exist in any warehouse",
  "timestamp": "2026-07-18T14:30:00Z",
  "request_id": "req-xyz789"
}
```

---

## 18. KESIMPULAN & NEXT STEPS

### RINGKASAN PRD v1.0.0

✅ **Project Name**: Supply Chain Intelligence Platform (ScIntel)

✅ **Scope Decision**: 80% Balanced (Forecasting deep, Inventory deep, Anomaly lite)

✅ **Timeline**: 12 weeks (October 2026)

✅ **Portfolio Target**: 10/10 backend engineering

✅ **Key Differentiators**:
- Advanced Golang patterns (concurrency, error handling)
- PostgreSQL mastery (indexing, query optimization, transactions)
- System design thinking (microservices, caching, webhooks)
- Production-ready code (80%+ test coverage, comprehensive docs)
- DevOps skills (Docker, Kubernetes-ready, CI/CD)

### NEXT STEPS

1. **APPROVE PRD** ✅ (Ready for signature)
2. **Create LOGIC_FLOW.txt** → Flowchart user journeys, API flows
3. **Create TDD (Test-Driven Design)** → Define API contracts & test cases
4. **Create DRP (Disaster Recovery Plan)** → Backup, failover, monitoring strategies
5. **Create STP (Setup & Technical Plan)** → Dev environment, tools, dependencies
6. **Create ROADMAP files** → Week-by-week implementation plan
7. **START CODING** → Begin Phase 1 (Foundation)

---

**Document Status**: ✅ APPROVED FOR DEVELOPMENT

**Approved By**: M. Arif Aulia

**Date**: 2026-07-18

**Version**: 1.0.0 (FINAL)

================================================================
