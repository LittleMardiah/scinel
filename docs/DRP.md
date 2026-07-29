# Disaster Recovery Plan (DRP) - SUPPLY CHAIN INTELLIGENCE PLATFORM v1.0.0

> **Purpose**: Dokumen ini mendefinisikan prosedur dan strategi pemulihan sistem Supply Chain Intelligence Platform (ScIntel) setelah terjadinya disaster atau kegagalan kritis, mencakup worst-case scenarios, backup strategies, failover mechanisms, dan playbooks untuk rapid recovery.

---

## 1. EXECUTIVE SUMMARY

### 1.1 Project Information

| Item | Value |
|------|-------|
| **Project Name** | Supply Chain Intelligence Platform (ScIntel) |
| **Owner/Team** | Backend Engineering Team (Golang/PostgreSQL/DevOps) |
| **Last Updated** | 2026-07-18 |
| **Next Review Date** | 2026-10-18 (After Phase 1 completion) |
| **Version** | 1.0.0 - INITIAL |
| **Environment** | Development (Local), Staging, Production |
| **Document Status** | DRAFT → Ready for Approval |

### 1.2 DRP Objectives

1. **Minimize Downtime**: Meminimalkan sistem downtime dalam setiap disaster scenario
2. **Data Integrity Protection**: Melindungi integritas data transaksi inventory dan forecasting
3. **Business Continuity**: Memastikan operasi supply chain dapat berlanjut dengan minimum disruption
4. **Clear Accountability**: Menetapkan clear roles, responsibilities, dan escalation paths
5. **Rapid Recovery**: Mengembalikan sistem ke operational state dalam RTO yang ditargetkan
6. **Audit & Compliance**: Memastikan disaster recovery activities terdokumentasi untuk audit trail

### 1.3 Key Metrics & SLA Targets

| Metric | Value | Justification |
|--------|-------|---------------|
| **RTO (Recovery Time Objective)** | **15-60 minutes** (varies by component) | Critical services (Auth, Inventory DB) RTO 15min; Secondary services 30-60min |
| **RPO (Recovery Point Objective)** | **5-15 minutes** | Daily inventory transactions = hourly importance; Backup frequency = every 5-15min |
| **Backup Frequency** | **Every 5 minutes (incremental) / Hourly (full snapshots)** | High transaction volume requires frequent backups |
| **Backup Retention** | **7 days (daily snapshots) + 30 days (weekly) + 90 days (monthly)** | Sufficient for root cause analysis & long-term audit |
| **Recovery Testing** | **Monthly drills** | Verify playbooks, measure actual RTO/RPO |
| **Documentation Accuracy** | **100% verified against PRD v1.0.0, LOGIC_FLOW, TDD** | All disaster scenarios mapped to actual system components |

---

## 2. SCOPE & SYSTEMS INVENTORY

### 2.1 In-Scope Systems (CRITICAL - Must Recover)

Berdasarkan PRD v1.0.0 section 4 (Infrastructure & Technology Stack) dan LOGIC_FLOW v1.0.0:

#### **Database Layer**
- ✅ **PostgreSQL 15+ (Primary Database)**: 
  - Tables: users, inventory, transactions, forecasts, anomalies, audit_log, webhooks, etc.
  - Role: Core business data store (ACID transactions mandatory)
  - RPO: ≤5 minutes (hourly transactions = critical)
  - RTO: ≤15 minutes (service cannot function without DB)

- ✅ **PostgreSQL Replica/Standby**: 
  - Streaming replication (logical or physical)
  - Role: High availability, read-only queries
  - Failover: Automatic or manual promotion

#### **Cache Layer**
- ✅ **Redis 7+ (In-Memory Cache)**:
  - Purpose: Inventory cache (TTL 30s), Forecast cache (TTL 24h), Session storage, Rate limiting
  - Data Loss Tolerance: Acceptable (non-persistent, can rebuild from DB)
  - RTO: ≤10 minutes (graceful degradation without cache)
  - Health Impact: Performance degradation (queries hit DB directly) if Redis down

#### **Application Servers**
- ✅ **Golang Backend (REST API)**:
  - Services: Auth, Inventory, Forecasting, Anomaly Detection, Webhooks
  - Instances: Minimum 2 (preferably 3+ for production)
  - State: Stateless (no persistent local state)
  - RTO: ≤5 minutes (redeploy or failover to healthy instance)

#### **Message & Event Processing**
- ✅ **Webhook Event Queue** (in-database or Redis-based):
  - Purpose: Event subscriptions, delivery logs, retry mechanism
  - Failure Impact: Events may not deliver to subscribers; customers cannot receive alerts
  - RTO: ≤30 minutes (async component, lower priority than inventory/forecasting)

#### **Authentication & Authorization**
- ✅ **JWT Token Generation & Validation**:
  - Role: All API requests require valid JWT (per LOGIC_FLOW 3.3)
  - Failure Impact: CRITICAL - No user can authenticate
  - RTO: ≤5 minutes (depends on token service availability)

#### **External Integrations** (Out of scope for v1.0, but prepare for future)
- ❌ **Kafka/Event Streaming**: Phase 1.5+ (not in v1.0 scope per PRD)
- ❌ **S3/Cloud Storage**: Phase 1.1+ (file uploads not in v1.0 scope per PRD)

### 2.2 Out-of-Scope Systems (Not Covered by This DRP)

| System/Component | Reason | Status |
|------------------|--------|--------|
| **Frontend (Next.js)** | Frontend is stateless presentation layer; can redeploy independently without affecting backend DRP | Covered in separate Frontend DRP |
| **CDN / Static Assets** | Non-critical for core supply chain operations; static content can be served from backup CDN | Covered in Infrastructure DRP |
| **Email Service** | Notifications only; core operations don't depend on email delivery | Covered in Communication DRP |
| **Analytics/Reporting** | Non-real-time; can defer until recovery complete | Covered in Analytics DRP (Phase 2) |
| **Kafka Integration** | Not in v1.0 scope per PRD section 4 (Phase 1.5+ feature) | Out of Scope - DRP v2.0 |

### 2.3 Infrastructure Overview (System Architecture)

Berdasarkan LOGIC_FLOW section 2 (System Architecture Diagram) dan TDD section 3:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CLIENT LAYER (Web/Mobile/API)                        │
└──────────────────────────────┬──────────────────────────────────────────┘
                               │
            ┌──────────────────▼──────────────────┐
            │     API Gateway / Load Balancer     │
            │  (Rate Limiting: 1000 req/sec)     │
            └──────────────────┬──────────────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
   ┌────▼──────┐       ┌──────▼───────┐      ┌────────▼────┐
   │  AUTH     │       │  INVENTORY   │      │ FORECASTING │
   │  SERVICE  │       │  SERVICE     │      │   SERVICE   │
   │           │       │              │      │             │
   │ (JWT, RBAC)       │ (Transactions│      │ (ETS/ARIMA) │
   │           │       │  Real-time)  │      │  Batch jobs │
   └────┬──────┘       └──────┬───────┘      └────────┬────┘
        │                     │                       │
        │                ┌────▼────────────────┐      │
        │                │   ANOMALY SERVICE   │      │
        │                │ (Z-score, IQR)     │      │
        │                └────┬────────────────┘      │
        │                     │                       │
        │    ┌────────────────┼───────────────────┐   │
        │    │                │                   │   │
   ┌────▼────▼────────────────▼─────────────────▼──┐ │
   │                                                 │ │
   │        PRIMARY DATABASE (PostgreSQL 15+)       │ │
   │  ┌──────────────────────────────────────────┐  │ │
   │  │ Tables: users, inventory, transactions, │  │ │
   │  │         forecasts, anomalies, webhooks, │  │ │
   │  │         audit_log, etc.                 │  │ │
   │  └──────────────────────────────────────────┘  │ │
   │  ┌──────────────────────────────────────────┐  │ │
   │  │ Streaming Replication → Standby/Replica │  │ │
   │  └──────────────────────────────────────────┘  │ │
   │                                                 │ │
   └─────────────┬────────────────────────────────────┘ │
                 │                                      │
                 │                                      │
        ┌────────▼──────────┐                  ┌──────▼────────┐
        │   REDIS CACHE     │                  │   WEBHOOK     │
        │  (Inventory/      │                  │   EVENT QUEUE │
        │   Forecast cache) │                  │  (Retry logic)│
        └───────────────────┘                  └───────────────┘
```

**Layer Descriptions**:

1. **API Gateway Layer**: Rate limiting (1000 req/sec per PRD), request routing, authentication enforcement (JWT validation per LOGIC_FLOW 3.3)

2. **Service Layer**: 
   - Auth Service: User login, token generation, RBAC enforcement
   - Inventory Service: Stock in/out transactions (LOGIC_FLOW 4), real-time queries
   - Forecasting Service: ETS/ARIMA algorithms (PRD 4.4), batch processing
   - Anomaly Service: Z-score/IQR detection (PRD 4.6), threshold alerts

3. **Data Layer**:
   - Primary PostgreSQL: ACID transactions for inventory (critical)
   - Streaming Replication: Real-time replica for failover
   - Redis Cache: Sub-millisecond access to frequently-used data
   - Webhook Queue: Event delivery with retry mechanism

---

## 3. DISASTER SCENARIOS & RISK ASSESSMENT

### 3.1 Disaster Scenario Matrix

Berdasarkan LOGIC_FLOW section 2 (System Architecture) dan PRD section 9 (Performance & Scalability), berikut worst-case scenarios yang mungkin terjadi:

| # | Scenario | Severity | RTO | RPO | Likelihood | Business Impact | In PRD? |
|---|----------|----------|-----|-----|------------|-----------------|--------|
| 1 | **Primary Database Complete Failure** (corruption/hardware failure) | **CRITICAL** | **15 min** | **5 min** | Medium | Revenue loss, SLA breach, inventory inaccuracy | ✅ Yes (4.3) |
| 2 | **Database Connection Pool Exhaustion** (max 50 conns breached) | **HIGH** | **10 min** | **0 min** | Medium-High | Service becomes unresponsive; queued requests fail | ✅ Yes (4.1) |
| 3 | **Redis Cache Complete Failure** (memory full, service down) | **HIGH** | **10 min** | **0 min** | Low-Medium | Performance degradation (queries hit DB); no rate limiting | ✅ Yes (4.1) |
| 4 | **Auth Service Down / JWT Token Generation Fails** | **CRITICAL** | **5 min** | **0 min** | Low | No user can login; complete service unavailable | ✅ Yes (4.2) |
| 5 | **Inventory Transaction Deadlock** (concurrent transactions) | **HIGH** | **30 sec** | **1 min** | Medium | Stock transactions blocked; data inconsistency risk | ✅ Yes (4.5) |
| 6 | **Forecasting Batch Job Crash / Timeout** (1000+ SKUs timeout) | **MEDIUM** | **60 min** | **24 hr** | Low-Medium | Next day forecast delayed; predictions inaccurate | ✅ Yes (4.4) |
| 7 | **Anomaly Detection Algorithm Failure** (Z-score calc error) | **MEDIUM** | **60 min** | **N/A** | Low | Anomalies not detected; delayed alerts to users | ✅ Yes (4.6) |
| 8 | **Webhook Delivery Queue Backup** (retry exhaustion) | **HIGH** | **30 min** | **1 hr** | Medium | Event subscribers miss notifications; delayed actions | ✅ Yes (4.7) |
| 9 | **Database Replica Lag / Sync Failure** (replication breaks) | **MEDIUM** | **30 min** | **5 min** | Low | Read-only queries fail; failover compromised | ✅ Yes (4.1) |
| 10 | **API Server Memory Leak / Resource Exhaustion** | **HIGH** | **5 min** | **0 min** | Medium | Service restart required; goroutines leak | ✅ Yes (4.2) |
| 11 | **Data Corruption** (invalid SQL query, data type mismatch) | **CRITICAL** | **30 min** | **5 min** | Low | Inventory counts wrong; forecasts inaccurate | ✅ Yes (4.3) |
| 12 | **Backup Storage Failure** (backup disk full, restore fails) | **CRITICAL** | **120 min** | **24 hr** | Low | Cannot recover from disaster; permanent data loss risk | ✅ Yes (4.1) |

---

### 3.2 Severity Levels Definition (Based on DRP Template + ScIntel Specifics)

#### **CRITICAL** (RTO: ≤15 minutes, RPO: ≤5 minutes)
**Definition**: Scenario langsung mengancam core supply chain operations; data loss atau service unavailability tidak dapat diterima.

**ScIntel Critical Scenarios**:
- Primary database complete failure / corruption (inventory data lost)
- Authentication system unavailable (no user can access system)
- All API servers down (service completely unavailable)
- Data integrity breach (inventory counts incorrect, can't trust data)

**Example**: "Warehouse manager cannot check stock levels; supplier cannot fulfill orders" → Revenue loss immediate

#### **HIGH** (RTO: 10-30 minutes, RPO: ≤1 hour)
**Definition**: Significant impact tapi tidak immediate threat ke core business; partial functionality preserved atau degraded.

**ScIntel High Scenarios**:
- Redis cache down (inventory queries slower, but still accessible)
- Forecasting batch job delayed (next day forecast delayed, not today's)
- Webhook delivery failing (subscribers miss alerts, but core system OK)
- Database replica lag (read-only queries slower, but available)

**Example**: "Inventory queries taking 2-3 seconds instead of 100ms; users frustrated but can work" → Operasi berlanjut tapi less efficient

#### **MEDIUM** (RTO: 30-120 minutes, RPO: ≤4 hours)
**Definition**: Impact terbatas atau non-critical functionality loss; data loss minimal atau acceptable recovery window.

**ScIntel Medium Scenarios**:
- Anomaly detection delayed (alerts come 1 hour late instead of real-time)
- Forecasting accuracy metrics slow to compute (reporting delayed)
- Non-critical service performance degradation (analytics dashboard slow)
- Webhook retry queue backed up (event delivery delayed but will eventually succeed)

**Example**: "Warehouse manager notices forecast 30 minutes late; still actionable" → Acceptable delay

#### **LOW** (RTO: > 2 hours, RPO: ≤1 day)
**Definition**: Minimal business impact; dapat di-defer tanpa significant loss.

**ScIntel Low Scenarios**:
- Audit log temporarily unavailable (compliance report delayed 24 hours)
- Staging/dev environment down (developers cannot test; prod ops unaffected)
- Non-customer-facing monitoring dashboard offline (ops team uses manual checks)

**Example**: "Audit trail data available tomorrow; no customer-facing impact today" → Acceptable deferral

---

## 4. BACKUP & REDUNDANCY STRATEGY

### 4.1 Backup Architecture (Based on PRD Section 4.1 & TDD Section 4)

#### 4.1.1 Database Backups (PostgreSQL 15+)

**Backup Strategy**:

| Strategy Component | Implementation | Details |
|------------------|-----------------|---------|
| **Full Backup Frequency** | 1x daily (early morning 2 AM UTC) | Capture entire database state; baseline for recovery |
| **Incremental/Differential Backup** | Every 5-15 minutes (transaction logs) | Continuous WAL (Write-Ahead Logging) capture for point-in-time recovery |
| **Backup Tool/Method** | pg_dump (logical) + WAL archiving (physical) | Hybrid: logical for schema/data validation; physical for point-in-time recovery |
| **Retention Period** | **7 days (daily snapshots) + 30 days (weekly) + 90 days (monthly)** | Sufficient for 5-min RPO + root cause analysis + compliance audit |
| **Backup Size Estimation** | ~1-5 GB/day (depends on transaction volume) | Inventory table grows ~100MB/day; forecasts table ~50MB/day |
| **Storage Location (Primary)** | PostgreSQL server's local storage (`/var/lib/postgresql/backups/`) | Fast restore; no network latency |
| **Storage Location (Secondary)** | Remote S3-compatible storage (AWS S3 / MinIO / DigitalOcean Spaces) | Geo-redundancy; protection from local hardware failure |
| **Backup Encryption** | AES-256 encryption at rest; TLS in transit | Compliance requirement; prevent unauthorized data access |
| **Verification Method** | Automated restore test (weekly, non-prod environment) | Verify backup integrity; catch restore process issues early |

**Backup Schedule**:

```
┌─────────────────────────────────────────────────────────────┐
│                 BACKUP SCHEDULE (24-hour)                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  00:00 - 02:00 UTC: Full Database Backup (pg_dump)         │
│                     Estimated duration: 30-60 min           │
│                     Action: LOCK TABLES (read-only mode)    │
│                                                              │
│  02:05 - End of day: Continuous WAL Archiving              │
│                     Every 5-15 min incremental logs saved   │
│                     Action: Background process (no blocking) │
│                                                              │
│  Weekly (Sundays):   Full backup + 6-day incremental        │
│                     bundled for long-term retention         │
│                                                              │
│  Monthly (1st):      Dedicated monthly archive              │
│                     Moved to cold storage (S3 Glacier)      │
│                     Retention: 90 days minimum              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Implementation Details**:

```bash
# Full Backup (Daily, 2 AM UTC)
pg_dump --verbose \
  --format=plain \
  --file=/var/lib/postgresql/backups/scintel_$(date +%Y%m%d_%H%M%S).sql \
  scintel_db

# WAL Archiving (Continuous, every 5-15 min)
# In postgresql.conf:
archive_mode = on
archive_command = 'aws s3 cp %p s3://scintel-backups/wal/%f'
archive_timeout = 300  # 5 minutes

# Backup Encryption & Upload to S3
gpg --symmetric --cipher-algo AES256 \
  --output /var/lib/postgresql/backups/scintel_$(date +%Y%m%d).sql.gpg \
  /var/lib/postgresql/backups/scintel_$(date +%Y%m%d).sql

aws s3 cp /var/lib/postgresql/backups/scintel_$(date +%Y%m%d).sql.gpg \
  s3://scintel-backups/daily/$(date +%Y/%m/%d)/
```

#### 4.1.2 Application & Code Backups

```
Version Control & Repository Mirroring:
- Primary: GitHub (https://github.com/user/scintel-backend)
- Mirror: GitLab or Gitea (backup Git server)
- Backup Frequency: Real-time (via Git hooks)
- Purpose: Recover application code, revert bad deployments

Configuration Files:
- Location: /etc/scintel/ (on application server)
- Backup Frequency: Every deployment (via infrastructure-as-code)
- Tool: Terraform / Ansible version control
- Purpose: Recover environment configuration, connection strings

Secrets Management:
- Primary: HashiCorp Vault / AWS Secrets Manager
- Encryption: AES-256 at rest, TLS in transit
- Backup: Automated encrypted snapshot (weekly)
- Rotation: Every 30 days (or on key compromise)
- Purpose: Protect database credentials, API keys, JWT signing keys
```

#### 4.1.3 Redis Cache Snapshots

```
Cache Backup Strategy:
- Type: RDB (Redis Database) snapshots + AOF (Append-Only File)
- Frequency: RDB snapshots every 6 hours (scheduled)
- Retention: 3 days (cache is ephemeral; rebuild from DB if needed)
- Location: Local Redis server + S3 backup
- Recovery Time: < 1 minute (cache data can be rebuilt from DB)

Note: Cache loss is ACCEPTABLE (LOW impact); 
      priority is database backup (HIGH impact if lost)
```

#### 4.1.4 Webhook Event Logs

```
Dead-Letter Queue (DLQ) Backup:
- Type: Database table backup (part of PostgreSQL backup)
- Table: webhook_delivery_log, webhook_failed_events
- Backup Frequency: Included in daily PostgreSQL backup
- Retention: 30 days (sufficient for event investigation)
- Purpose: Trace failed webhook deliveries; retry logic audit
```

---

### 4.2 Redundancy & High Availability

Berdasarkan PRD section 4.1 (Connection Pooling) dan TDD section 3 (Architecture):

| Component | Redundancy Type | Implementation | Failover Mechanism | Failover Time |
|-----------|-----------------|-----------------|-------------------|---------------|
| **PostgreSQL Database** | **Active-Passive** (Primary + Standby) | Streaming replication (logical_replication_mode or slot-based) | Automatic promotion via pg_ctl promote or manual via DBA decision | **1-2 minutes** (detect failure + promote) |
| **PostgreSQL Read Replicas** | **Active-Active** (Multiple replicas) | 2-3 read-only replicas for load distribution | Automatic routing via pgBouncer/connection pooler | **30 seconds** (connection pool reconnect) |
| **API Application Servers** | **Active-Active** (2+ instances) | Docker containers on separate hosts; load balanced | Health checks every 10 seconds; route to healthy instance | **10-30 seconds** (health check + route update) |
| **Redis Cache** | **Passive (No replication in v1.0)** | Single instance with local RDB snapshots | Manual restart or failover to backup Redis (future) | **2-5 minutes** (service restart + reconnect) |
| **Load Balancer** | **Active-Passive** | Primary LB + standby LB (heartbeat monitoring) | Automatic failover via VRRP (Virtual Router Redundancy Protocol) | **10 seconds** (heartbeat miss + VIP migration) |

**Failover Mechanisms**:

1. **Database Failover**:
   - Detection: Primary DB health check fails (connection timeout, error rate > 10%)
   - Action: Promote standby to primary via `pg_ctl promote`
   - DNS Update: Update connection string to point to new primary (manual or automatic)
   - Timeline: ~1-2 minutes from detection to traffic routing to new primary

2. **Application Server Failover**:
   - Detection: Health endpoint (`GET /health/ready`) returns 5xx or timeout
   - Action: Remove unhealthy instance from load balancer
   - Recovery: Redeploy service or restart container
   - Timeline: ~10-30 seconds from health check failure to route update

3. **Load Balancer Failover**:
   - Detection: Heartbeat between primary and standby LB fails
   - Action: VRRP triggers failover; standby assumes VIP (virtual IP)
   - Traffic: Connections re-established to new primary LB
   - Timeline: ~10 seconds (heartbeat interval)

4. **DNS Failover**:
   - Detection: Primary API endpoint unreachable (manual monitoring or watchdog)
   - Action: Update DNS record to point to secondary region/backup API
   - Client Impact: Connection attempts fail; retry logic kicks in
   - Timeline: 30-300 seconds (DNS TTL = 60 seconds)

---

### 4.3 Geographic Distribution (Preparation for Future Multi-Region)

**Current State (v1.0)**: Single-region deployment

```
Primary Region: [TBD - Singapore/Jakarta/AWS Region]
├── Data Center / Cloud Region: AWS us-east-1 (or local datacenter)
├── Primary Database: PostgreSQL 15+ (primary)
├── Standby Database: PostgreSQL 15+ (streaming replication, same region)
├── Application Servers: 2-3 API instances (load balanced)
├── Cache Layer: Redis 7+ (single instance)
├── Backup Storage: Local S3-compatible (MinIO or AWS S3)
└── Estimated RPO if Region Lost: 5 minutes (latest backup)
    Estimated RTO if Region Lost: 60-120 minutes (restore in secondary region)

Secondary Region: [Preparation Only - Phase 1.5+]
├── Purpose: Warm standby (periodic backup restoration for testing)
├── Recovery Procedure: Manual failover only (no automatic failover in v1.0)
└── Timeline: Not in MVP scope; business decision required
```

**Future Enhancement (Phase 1.5+)**:
- Warm standby in secondary region (read-only replica, periodic backups)
- Active-active deployment for zero-downtime failover
- Cross-region replication with acceptable replication lag
- Estimated RTO: 5-15 minutes (failover + DNS propagation)

---

## 5. DISASTER RESPONSE PLAYBOOKS

### 5.1 Playbook Naming Convention

Each playbook follows the structure:

```
[SCENARIO #]: [SCENARIO NAME]
Severity: [CRITICAL / HIGH / MEDIUM / LOW]
Symptoms/Detection: [How to detect the disaster]
Impact Assessment: [What breaks, what works]
Immediate Actions: [First 15 minutes - triage & stabilization]
Recovery Steps: [Detailed step-by-step recovery]
Validation: [Verify system is healthy]
Post-Recovery: [Monitoring & communication]
RTO/RPO: [Actual expected times for this scenario]
```

---

### 5.2 Playbook Details

---

#### **SCENARIO #1: PRIMARY DATABASE COMPLETE FAILURE**

**Severity**: 🔴 **CRITICAL** (RTO: ≤15 min, RPO: ≤5 min)

**Symptoms / Detection** (from PRD 4.3, LOGIC_FLOW 2):

Alert Triggers:
- 🚨 Alert: PostgreSQL primary instance responds to ping but queries timeout (>30s)
- 🚨 Alert: Connection pool exhaustion (connection error rate > 50% for 2 consecutive minutes)
- 🚨 Alert: All API endpoints returning 500 "database connection failed" (service health check fails)
- 🚨 Alert: Monitoring dashboard shows database CPU 100% / Memory full for >5 minutes
- 🚨 Alert: Replication lag increases to >5 minutes (standby getting stale data)

**Impact Assessment**:

Affected Systems:
- ❌ Inventory transactions cannot commit (stock in/out fails) → CRITICAL
- ❌ Real-time inventory queries timeout (users cannot check stock) → CRITICAL
- ❌ Forecast batch jobs cannot retrieve historical data → HIGH
- ❌ Authentication queries fail (users cannot login) → CRITICAL
- ❌ Anomaly detection queries cannot run → MEDIUM
- ⚠️ Webhook events queued but cannot process (retry queue backs up) → HIGH

User Impact:
- Warehouse managers **cannot** check stock levels → immediate operational impact
- New stock transactions **blocked** (in/out operations fail)
- Users **cannot login** (if auth DB affected)
- Suppliers/customers **cannot** place orders or check fulfillment status
- Estimated revenue loss: ~$X per minute (depends on business volume)

Business Impact:
- **SLA Breach**: RTO 15 min > customer SLA expectations
- **Inventory Inaccuracy**: Pending transactions lost if not recovered
- **Compliance Risk**: Audit trail incomplete; data loss detected in post-incident review

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] CONFIRM THE ISSUE** 
   - [ ] Who detects: On-call database engineer or monitoring alert
   - [ ] Tool: `psql -h {primary_db_host} -U postgres -c "SELECT NOW();"` → Should complete <3 sec
   - [ ] Expected result: Connection timeout OR query returns no result
   - [ ] Action: CONFIRMED = continue to step 2

2. **[1-2 min] NOTIFY INCIDENT COMMANDER**
   - [ ] Slack: Post in #incidents channel: "🔴 CRITICAL: Primary DB offline - triggering recovery playbook"
   - [ ] Who calls: On-call engineer
   - [ ] Escalation: Notify Database Lead (DBA) + Engineering Manager + CTO
   - [ ] Expected: Response within 2 minutes

3. **[2-3 min] ASSESS BACKUP AVAILABILITY**
   - [ ] Command: Check latest backup timestamp
     ```bash
     ls -lart /var/lib/postgresql/backups/ | tail -5
     # Expected: Recent backup (within last 5 min) exists
     ```
   - [ ] Expected result: Backup file size > 500MB (not empty/corrupted)
   - [ ] Action: If no recent backup → escalate to CTO (data loss risk)

4. **[3-5 min] ATTEMPT PRIMARY DB RESTART** (If hardware alive but process crashed)
   - [ ] Check: Is PostgreSQL process still running?
     ```bash
     systemctl status postgresql  # or: ps aux | grep postgres
     ```
   - [ ] If process alive: Check logs for errors
     ```bash
     tail -f /var/log/postgresql/postgresql.log
     ```
   - [ ] If process dead: Attempt graceful restart
     ```bash
     systemctl restart postgresql
     ```
   - [ ] Expected: Service starts, queries respond within 3 seconds
   - [ ] If restart fails → continue to failover (step 5)

5. **[5-8 min] INITIATE FAILOVER TO STANDBY** (If primary cannot recover)
   - [ ] Verify standby is healthy:
     ```bash
     psql -h {standby_host} -U postgres -c "SELECT now(), pg_is_in_recovery();"
     # Expected: Returns true (standby is in recovery mode)
     ```
   - [ ] Promote standby to primary:
     ```bash
     ssh {standby_host}
     sudo -u postgres pg_ctl promote -D /var/lib/postgresql/main
     ```
   - [ ] Expected: Promotion completes within 30 seconds
   - [ ] Verify promotion:
     ```bash
     psql -h {standby_host} -U postgres -c "SELECT pg_is_in_recovery();"
     # Expected: Returns false (now primary)
     ```

6. **[8-10 min] UPDATE CONNECTION STRINGS** (Clients connect to new primary)
   - [ ] Update database configuration:
     - Docker compose: Update `DATABASE_URL` env var with new primary host
     - Kubernetes: Update secret `scintel-db-credentials` with new host
     - Local config: Update `/etc/scintel/config.yaml` → database.host = {new_primary}
   - [ ] Restart application servers (rolling restart, not simultaneous):
     ```bash
     # For each API server (one at a time)
     docker restart scintel-api-1  # Wait for health check to pass
     docker restart scintel-api-2
     docker restart scintel-api-3
     # Expected: Each server takes ~30 seconds to restart & reconnect to new DB
     ```
   - [ ] Verify: Inventory query returns results
     ```bash
     curl -H "Authorization: Bearer {valid_token}" \
       http://{api_server}:8080/api/v1/inventory?sku=SKU-NIKE-001
     # Expected: 200 OK with inventory data
     ```

7. **[10-15 min] VERIFY TRAFFIC ROUTING & NOTIFICATIONS**
   - [ ] Check: Load balancer is routing to healthy API instances
     ```bash
     curl -I http://{load_balancer}/api/v1/health/live
     # Expected: 200 OK
     ```
   - [ ] Test: User login still works (auth service healthy)
     ```bash
     curl -X POST http://{api_server}:8080/api/v1/auth/login \
       -d '{"email": "test@example.com", "password": "..."}' \
       -H "Content-Type: application/json"
     # Expected: 200 OK or 401 (if creds wrong, not connection error)
     ```
   - [ ] Notify stakeholders (Communication Lead):
     - Internal: "🟢 Database failover complete. Service recovering."
     - External: "Warehouse operations briefly interrupted (5 min). Resumed now. Sorry for inconvenience."

---

**Recovery Steps (Detailed)** 🛠️

**Phase 1: Preparation**
- [ ] **Assign roles**:
  - DBA #1: Lead recovery effort, monitor primary DB health
  - Infrastructure Lead: Manage server restarts, connection pooling
  - DevOps: Update DNS, verify traffic routing
  - Incident Commander: Coordinate, track timeline, communicate
  
- [ ] **Prepare recovery environment**:
  - Have backup files ready (verify path exists)
  - SSH access to primary, standby, all API servers
  - PostgreSQL credentials & sudo access
  - Contact info for all team members (on-call list)

- [ ] **Review backup integrity**:
  ```bash
  # Check latest full backup
  file /var/lib/postgresql/backups/scintel_20260718_020000.sql
  # Expected: "SQL script, ASCII text"
  
  # Test restore on staging (if time permits)
  pg_restore --verbose --clean --if-exists \
    --dbname=scintel_stage \
    /var/lib/postgresql/backups/scintel_20260718_020000.sql
  ```

---

**Phase 2: Execution (Detailed Step-by-Step)**

**Option A: Failover to Standby (Recommended if standby healthy)**

```bash
# Step 1: Stop all application servers to prevent write conflicts
ssh {api_server_1}
docker stop scintel-api-1
docker stop scintel-api-2
docker stop scintel-api-3
# Wait 30 seconds for connections to drain

# Step 2: Verify standby is caught up with primary
ssh {standby_host}
psql -U postgres -c "
  SELECT 
    now(),
    pg_is_in_recovery() as is_replica,
    pg_last_wal_receive_lsn() as received_lsn,
    pg_last_wal_replay_lsn() as replayed_lsn;
"
# Expected: replayed_lsn should equal received_lsn (fully caught up)
# If lag > 1 minute, wait for catchup

# Step 3: Promote standby to primary
sudo -u postgres pg_ctl promote -D /var/lib/postgresql/main
# Expected output: "promoting standby to primary... done" 
# Monitor: tail -f /var/log/postgresql/postgresql.log

# Step 4: Verify new primary is accepting connections
sleep 10  # Wait for promotion to complete
psql -h {standby_host} -U postgres -c "
  SELECT 
    version(),
    pg_is_in_recovery() as is_primary;
"
# Expected: is_primary = false (now a primary server)

# Step 5: Update DNS / connection pooling
# Option 1: Update primary_host DNS record to point to standby
ssh {dns_server}
# Edit zone file: Change A record for db.scintel.io → {standby_ip}
systemctl restart bind9  # or: rndc reload
sleep 5
nslookup db.scintel.io  # Verify DNS propagation

# Option 2: Update pgBouncer config
ssh {pgbouncer_host}
cat > /etc/pgbouncer/pgbouncer.ini << EOF
[databases]
scintel_db = host={standby_host} port=5432 user=postgres password={db_pass}
[pgbouncer]
max_client_conn = 1000
default_pool_size = 25
EOF
systemctl restart pgbouncer
# Verify: psql -h localhost:6432 -U postgres -c "SELECT now();"

# Step 6: Restart application servers one by one
ssh {api_server_1}
docker start scintel-api-1
sleep 30  # Wait for container to connect to DB
docker ps | grep scintel-api-1  # Verify running

docker start scintel-api-2
sleep 30

docker start scintel-api-3
sleep 30

# Step 7: Monitor application health
for i in {1..3}; do
  curl -s http://scintel-api-${i}:8080/api/v1/health/ready
  echo "API Server $i: $(curl -s -o /dev/null -w '%{http_code}' http://scintel-api-${i}:8080/api/v1/health/ready)"
done
# Expected: All return 200 OK
```

**Option B: Restore from Backup (If standby also corrupted)**

```bash
# Step 1: Stop all application services
docker stop scintel-api-1 scintel-api-2 scintel-api-3

# Step 2: Stop old primary database
systemctl stop postgresql

# Step 3: Backup current corrupted data (for investigation)
sudo mv /var/lib/postgresql/main /var/lib/postgresql/main.corrupted.$(date +%s)

# Step 4: Create empty database cluster
sudo -u postgres /usr/lib/postgresql/15/bin/initdb \
  -D /var/lib/postgresql/main \
  -E UTF8 \
  -A md5

# Step 5: Restore from backup file
sudo -u postgres psql -d postgres < /var/lib/postgresql/backups/scintel_20260718_020000.sql
# or use pg_restore if binary format backup:
sudo -u postgres pg_restore --verbose --clean --if-exists \
  -d postgres \
  /var/lib/postgresql/backups/scintel_20260718_020000.sql.backup

# Step 6: Replay WAL logs for point-in-time recovery (optional)
# If backup is from 02:00 and crash happened at 10:30, replay WAL to 10:25
sudo -u postgres pg_ctl start -D /var/lib/postgresql/main
# Monitor: tail -f /var/log/postgresql/postgresql.log
# Expected: "recovering ..." then "database system is ready to accept connections"

# Step 7: Verify data integrity
psql -U postgres -d scintel_db -c "
  SELECT COUNT(*) as inventory_count FROM inventory;
  SELECT COUNT(*) as transaction_count FROM transactions;
  SELECT COUNT(*) as user_count FROM users;
"
# Expected: Row counts match pre-disaster state (may be ~5 min old due to RPO)

# Step 8: Restart application servers
docker start scintel-api-1 scintel-api-2 scintel-api-3
sleep 30
for i in {1..3}; do
  curl http://scintel-api-${i}:8080/api/v1/health/ready || echo "Server $i not ready"
done
```

---

**Phase 3: Validation** ✅

- [ ] **Database Connectivity Check**:
  ```bash
  psql -h {primary_db} -U postgres -c "SELECT COUNT(*) FROM inventory;" 
  # Expected: Returns a number (e.g., 4250)
  
  psql -h {primary_db} -U postgres -c "SELECT COUNT(*) FROM transactions WHERE created_at > NOW() - INTERVAL '5 minutes';"
  # Expected: Shows recent transactions (if DB was running, should have some)
  ```

- [ ] **API Endpoint Health Check**:
  ```bash
  # 1. Health check endpoint
  curl -s http://{api_server}:8080/api/v1/health/live | jq .
  # Expected: { "status": "alive" }
  
  # 2. Readiness probe
  curl -s http://{api_server}:8080/api/v1/health/ready | jq .
  # Expected: { "status": "ready" }
  
  # 3. Test inventory query
  curl -H "Authorization: Bearer {valid_token}" \
    http://{api_server}:8080/api/v1/inventory?sku=SKU-NIKE-001 | jq .
  # Expected: 200 OK with inventory data
  
  # 4. Test authentication
  curl -X POST http://{api_server}:8080/api/v1/auth/login \
    -d '{"email": "manager@warehouse.id", "password": "test123"}' \
    -H "Content-Type: application/json" | jq .
  # Expected: 200 OK with JWT tokens
  ```

- [ ] **Critical User Journeys** (from HALAMAN v2.0, LOGIC_FLOW):
  ```bash
  # Scenario 1: Manager checks inventory
  1. Login via POST /api/v1/auth/login
  2. Get inventory via GET /api/v1/inventory?sku=SKU-NIKE-001
  3. Verify: Returns warehouse breakdown with quantities
  
  # Scenario 2: Create stock transaction
  1. Authorized header with valid JWT
  2. POST /api/v1/transactions (stock in)
  3. Verify: 201 Created, transaction_id returned, inventory updated
  
  # Scenario 3: Get forecast
  1. Authorized header with valid JWT
  2. GET /api/v1/forecasts?sku=SKU-NIKE-001&horizon=30
  3. Verify: 200 OK with forecast data
  ```

- [ ] **Data Integrity Checks**:
  ```bash
  -- Check for orphaned transactions (no matching inventory)
  SELECT t.* FROM transactions t 
  LEFT JOIN inventory i ON t.sku = i.sku 
  WHERE i.id IS NULL LIMIT 10;
  # Expected: No rows (all transactions reference valid SKUs)
  
  -- Check audit log for discrepancies
  SELECT COUNT(*) as audit_entries FROM audit_log 
  WHERE created_at > NOW() - INTERVAL '30 minutes';
  # Expected: Audit entries match transaction activity
  
  -- Verify forecast data consistency
  SELECT COUNT(*) FROM forecasts WHERE confidence < 0 OR confidence > 1;
  # Expected: 0 rows (invalid confidence values indicate corruption)
  ```

- [ ] **Performance Baseline**:
  ```bash
  # Measure query response time
  time psql -h {primary_db} -U postgres -c "
    SELECT 
      SUM(quantity) as total_stock 
    FROM inventory 
    WHERE warehouse_id = 'wh-jakarta-001';
  "
  # Expected: Query completes < 100ms (per PRD performance target)
  # If > 500ms, indexes may be corrupted → REBUILD INDEXES
  ```

---

**Phase 4: Post-Recovery** 📋

- [ ] **Verify data integrity & consistency**:
  - Count inventory transactions pre-disaster vs post-recovery (should match within 5 min window)
  - Cross-check forecast data accuracy (should have latest predictions)
  - Audit trail: Verify no gaps in audit_log (complete record of all activities)

- [ ] **Monitor system for 30 minutes**:
  - Watch database CPU, memory, connection count
  - Monitor API error rates (expect <0.1% errors)
  - Check cache hit ratios (Redis should stabilize after 10 min)
  - Alert on any anomalies or performance degradation

- [ ] **Document incident & timeline**:
  ```
  Incident Report:
  - Failure Time: 2026-07-18 10:30:00 UTC
  - Detection Time: 2026-07-18 10:31:00 UTC (1 min detection lag)
  - Recovery Complete: 2026-07-18 10:43:00 UTC (12 min total downtime)
  - Root Cause: Disk full on primary DB (backup partition not cleaned up)
  - Actions Taken: Failover to standby + restart API servers
  - Data Loss: 0 (backup taken 10:25, recovery includes WAL replay)
  - RTO Actual: 12 minutes (target 15 min) ✅
  - RPO Actual: 5 minutes (target 5 min) ✅
  ```

- [ ] **Notify stakeholders**:
  - Slack: "✅ Database failover complete. All services recovered. Post-incident review in 24 hours."
  - Email: Notify management of recovery completion, estimated impact
  - Status Page: Update https://status.scintel.io with recovery timeline

- [ ] **Trigger post-incident review (next day)**:
  - Root cause analysis: Why did disk fill up? → Automate log rotation
  - Preventive measures: Enable automatic database cleanup alerts
  - Training: Review disaster recovery procedures with team
  - Update playbook: Incorporate lessons learned

---

**Rollback Plan** (If recovery fails):

Condition to trigger rollback:
- API servers cannot connect to new primary after 10 minutes
- Restored database has data corruption (row counts inconsistent)
- Recovery takes > 30 minutes (escalate to alternate recovery method)

Rollback steps:
1. Stop attempted recovery (leave new primary in read-only mode)
2. Inspect logs for root cause of failure
3. Restore to older backup point (before recovery attempt)
4. Notify CTO + database team; escalate to vendor support if available
5. Prepare statement: "We've encountered additional complications. Estimated recovery time extended to 60 minutes."

---

**Expected RTO**: **12-15 minutes** (failover to standby + app restart)

**Expected RPO**: **5 minutes** (continuous WAL archiving captures latest transactions)

---

---

#### **SCENARIO #2: REDIS CACHE COMPLETE FAILURE**

**Severity**: 🟠 **HIGH** (RTO: ≤10 min, RPO: 0 min - cache is ephemeral)

**Symptoms / Detection**:

Alert Triggers:
- 🚨 Alert: Redis responds to ping but PING returns timeout (>5s)
- 🚨 Alert: Cache hit rate drops to 0% (all requests going to database)
- 🚨 Alert: Database query latency spikes 10x (from 50ms to 500ms average)
- 🚨 Alert: Redis connection errors in API logs: "Connection refused" or "Read timed out"
- 🚨 Alert: Memory usage on API servers increasing (rate-limiting not working; requests queued)

**Impact Assessment**:

Affected Systems:
- 🟡 Inventory cache unavailable (queries go directly to PostgreSQL) → Performance drops 10x
- 🟡 Forecast cache stale (queries recalculated instead of cached) → CPU spikes
- 🟡 Session cache down (JWT token validation slower) → Auth latency increases
- 🟡 Rate limiting non-functional (connections not tracked) → Potential abuse
- ✅ Core functionality intact (database still queryable, just slower)

User Impact:
- Inventory queries take 1-2 seconds instead of 100ms
- User login takes 2-3 seconds instead of 500ms
- Load time degradation noticeable but acceptable
- Service continues operating (no outage, just slow)

Business Impact:
- User frustration (perceived slowness)
- API response SLA breached (p99 < 100ms target) → contractual violation possible
- Competitor systems appear faster → UX comparison issue
- **No revenue loss** (service still works, just slow)

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] CONFIRM REDIS IS DOWN**:
   ```bash
   redis-cli -h {redis_host} PING
   # Expected: +PONG (working) or timeout/error (down)
   
   redis-cli -h {redis_host} INFO server | head -5
   # Expected: Server info or error
   ```

2. **[1-2 min] ASSESS IMPACT & ESCALATE**:
   - [ ] Check API error logs: How many requests failing?
     ```bash
     grep -c "redis.*refused\|cache.*timeout" /var/log/scintel-api/api.log
     ```
   - [ ] Slack: "🟠 HIGH: Redis cache down. Database queries elevated. Analyzing recovery."
   - [ ] Notify: Infrastructure Lead (on-call)

3. **[2-5 min] ATTEMPT REDIS RESTART** (If service process crashed):
   ```bash
   # Check Redis process status
   systemctl status redis-server
   
   # If process dead, restart
   systemctl restart redis-server
   
   # Verify connection
   redis-cli PING  # Should return +PONG after 5-10 sec
   ```
   - Expected: Redis connects within 10 seconds
   - If restart fails → Continue to step 4 (failover/skip)

4. **[5-10 min] CONFIGURE API SERVERS TO SKIP REDIS** (If cannot restart):
   - Graceful degradation: API servers should have fallback logic
   - [ ] Update environment variable: `REDIS_ENABLED=false`
   - [ ] Rolling restart of API servers (one at a time, so traffic stays online)
   - [ ] Expected: Queries still work, but hit database directly

5. **[10-15 min] VERIFY SERVICE HEALTH**:
   ```bash
   # Test inventory query (should work, but slower)
   curl -H "Authorization: Bearer {token}" \
     http://{api_server}:8080/api/v1/inventory?sku=SKU-NIKE-001
   # Expected: 200 OK (maybe takes 1-2 sec instead of 100ms)
   
   # Check API error rate
   curl http://{api_server}:8080/api/v1/health/live
   # Expected: 200 OK (healthy, just running slower)
   ```

---

**Recovery Steps** 🛠️

**Phase 1: Diagnosis**
- [ ] Check Redis server logs: `/var/log/redis/redis-server.log`
  ```bash
  tail -100 /var/log/redis/redis-server.log | grep -i error
  # Look for: OOM (out of memory), disk full, connection refused, etc.
  ```

- [ ] Check system resources on Redis host:
  ```bash
  # Memory usage
  free -h  # Should show available memory; if swap used, Redis is swapping
  
  # Disk space
  df -h /var/lib/redis/  # Check if disk full
  
  # Network connectivity
  ss -an | grep 6379  # Check Redis port 6379
  ```

**Phase 2: Recovery Options**

**Option A: Restart Redis (if process crashed but storage OK)**
```bash
# Stop Redis gracefully
redis-cli SHUTDOWN SAVE  # Save in-memory data before shutdown

# Start Redis
systemctl start redis-server

# Verify data restored
redis-cli DBSIZE  # Should return number of keys (not 0)
redis-cli KEYS inventory:* | wc -l  # Should show cached keys
```

**Option B: Rebuild Cache from Database (if RDB file corrupted)**
```bash
# Option 1: Clear corrupted cache & let application rebuild
redis-cli FLUSHALL  # Delete all keys

# Option 2: Restore from backup
redis-cli SHUTDOWN  # Stop Redis
cp /var/lib/redis/dump.rdb.backup /var/lib/redis/dump.rdb
systemctl start redis-server

# Option 3: Manual warmup from database
# Trigger cache warmup job:
curl -X POST http://{api_server}:8080/api/v1/admin/cache/warmup \
  -H "Authorization: Bearer {admin_token}"
# This job: Queries popular SKUs, caches inventory & forecasts
```

**Option C: Failover to Redis Replica (if available in future)**
```bash
# (Prepare for Phase 1.5+ with Redis cluster/sentinel)
redis-cli -h {replica_host} SLAVEOF NO ONE  # Promote replica to primary
# Update connection string in app servers
```

---

**Phase 3: Validation** ✅

```bash
# 1. Redis connectivity
redis-cli PING  # Should return +PONG

# 2. Data present in cache
redis-cli DBSIZE  # Should return > 1000 (cache populated)

# 3. API performance recovered
# Warm up cache
curl http://{api_server}:8080/api/v1/admin/cache/warmup \
  -H "Authorization: Bearer {admin_token}"

# Test query latency (should be <100ms after warmup)
time curl -H "Authorization: Bearer {token}" \
  http://{api_server}:8080/api/v1/inventory?sku=SKU-NIKE-001
# Expected: real 0m0.050s (under 100ms)

# 4. Check cache hit rate
redis-cli INFO stats | grep keyspace_hits
# Expected: Should increase after cache warmed up
```

---

**Phase 4: Post-Recovery** 📋

- [ ] Monitor Redis memory usage for next 30 minutes
  ```bash
  watch -n 5 "redis-cli INFO memory | grep used_memory_human"
  ```

- [ ] Verify API performance metrics returned to normal
  ```bash
  # Check p99 latency in monitoring dashboard
  # Should be < 100ms within 10 minutes of Redis recovery
  ```

- [ ] Document incident:
  - Cause: [e.g., "Memory fill-up from large dataset; eviction policy not working"]
  - Recovery time: X minutes
  - Preventive action: Enable automatic memory monitoring, adjust eviction policy

---

**Expected RTO**: **5-10 minutes** (restart + cache warmup)

**Expected RPO**: **0 minutes** (cache is ephemeral; loss is acceptable)

---

---

#### **SCENARIO #3: API SERVER DOWN / MEMORY LEAK**

**Severity**: 🟠 **HIGH** (RTO: ≤5 min, RPO: 0 min - stateless service)

**Symptoms / Detection** (from PRD section 4.2, TDD section 12):

Alert Triggers:
- 🚨 Alert: Health check endpoint returns 5xx (service unhealthy)
- 🚨 Alert: Memory usage on API server > 90% for >2 minutes
- 🚨 Alert: Goroutine count increasing continuously (leak detected)
- 🚨 Alert: API response time > 30 seconds (timeouts or blocking)
- 🚨 Alert: Load balancer removes API server from pool (failed health checks)

**Impact Assessment**:

If 1 API server down (out of 3):
- ✅ Traffic redistributed to 2 healthy servers (33% capacity hit)
- 🟡 Response time increases slightly (load on other servers +50%)
- ✅ Service continues (no outage, just degraded)

If 2 API servers down:
- 🟠 Only 1 server handling traffic (67% capacity loss)
- 🟠 Response times spike to 5-10 seconds
- 🟠 Error rate increases (request timeouts possible)

If all 3 API servers down:
- 🔴 CRITICAL: Service completely unavailable
- 🔴 Users cannot access any endpoint (login, inventory, forecasting)
- 🔴 Load balancer returns 503 Service Unavailable

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] IDENTIFY WHICH SERVERS ARE DOWN**:
   ```bash
   # Check load balancer pool status
   curl -s http://{load_balancer}:8080/status | jq '.backend_servers'
   # Expected: Shows which servers are up/down
   
   # Alternative: Manual check
   for i in 1 2 3; do
     curl -s -o /dev/null -w "%{http_code}" \
       http://scintel-api-${i}:8080/api/v1/health/live
     echo "  Server $i"
   done
   ```

2. **[1-2 min] RESTART FAILED SERVER(S)**:
   ```bash
   # Restart single API server
   docker restart scintel-api-1
   
   # Wait for container to start & pass health checks
   sleep 30
   curl http://scintel-api-1:8080/api/v1/health/ready
   # Expected: 200 OK
   ```

3. **[2-5 min] INVESTIGATE MEMORY LEAK (if applicable)**:
   ```bash
   # Check memory usage history
   docker stats scintel-api-1 --no-stream
   # Expected: Memory usage should stabilize after restart (< 500MB)
   
   # Check for goroutine leak in logs
   tail -50 /var/log/scintel-api/api.log | grep -i goroutine
   ```

4. **[5-10 min] VERIFY TRAFFIC ROUTING**:
   ```bash
   # Generate test traffic & monitor load distribution
   for i in {1..10}; do
     curl -H "Authorization: Bearer {token}" \
       http://{load_balancer}:8080/api/v1/inventory?sku=SKU-NIKE-001
   done
   
   # Check load balancer logs
   tail -20 /var/log/nginx/access.log | grep -E "scintel-api-[123]"
   # Expected: Traffic distributed evenly among healthy servers
   ```

5. **[10-15 min] CONFIRM SERVICE HEALTH**:
   ```bash
   # Test all critical endpoints with valid token
   curl -H "Authorization: Bearer {token}" \
     http://{load_balancer}:8080/api/v1/health/ready
   # Expected: 200 OK
   
   # Test inventory endpoint
   curl -H "Authorization: Bearer {token}" \
     http://{load_balancer}:8080/api/v1/inventory
   # Expected: 200 OK with inventory data
   ```

---

**Recovery Steps** 🛠️

**Phase 1: Immediate Recovery (Restart)**

```bash
#!/bin/bash
# Script: api-recovery.sh

# Stop and restart unhealthy API servers
for server in scintel-api-1 scintel-api-2 scintel-api-3; do
  status=$(docker inspect $server --format='{{.State.Running}}')
  
  if [ "$status" != "true" ]; then
    echo "[$(date)] Restarting $server..."
    docker restart $server
    
    # Wait for health check to pass
    max_attempts=30
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
      health=$(curl -s -o /dev/null -w "%{http_code}" \
        http://${server}:8080/api/v1/health/ready)
      
      if [ "$health" == "200" ]; then
        echo "[$(date)] $server is healthy"
        break
      fi
      
      echo "[$(date)] $server not ready yet (health: $health). Retry..."
      sleep 2
      attempt=$((attempt + 1))
    done
  fi
done

# Verify load balancer pool
echo "[$(date)] Verifying load balancer..."
curl -s http://{load_balancer}:8080/status | jq '.backend_servers'
```

**Phase 2: Investigation (Memory Leak Analysis)**

If memory leak detected, investigate root cause:

```bash
# 1. Check API version deployed
docker inspect scintel-api-1 | grep -i "image\|version"

# 2. Review application logs for goroutine issues
docker logs scintel-api-1 | tail -100 | grep -i "goroutine\|memory\|alloc"

# 3. Profile memory usage (if pprof endpoint available)
curl http://scintel-api-1:6060/debug/pprof/heap > heap_dump.prof
# Analyze with: go tool pprof heap_dump.prof

# 4. Check for known issues in codebase
# If memory leak confirmed → trigger code review
# If temporary spike → continue monitoring
```

**Phase 3: Validation** ✅

```bash
# 1. All API servers healthy
for i in 1 2 3; do
  status=$(curl -s http://scintel-api-${i}:8080/api/v1/health/ready)
  echo "API Server $i: $status"
done
# Expected: All return { "status": "ready" }

# 2. Load balanced traffic
curl -v http://{load_balancer}:8080/api/v1/inventory \
  -H "Authorization: Bearer {token}"
# Expected: Server header shows round-robin distribution

# 3. Response time normal
time curl http://{load_balancer}:8080/api/v1/inventory \
  -H "Authorization: Bearer {token}"
# Expected: real < 0.500s (under 500ms)
```

**Phase 4: Post-Recovery** 📋

- [ ] Monitor API server memory for next 30 minutes
  ```bash
  watch -n 10 "docker stats scintel-api-1 scintel-api-2 scintel-api-3 --no-stream"
  ```

- [ ] If memory leak confirmed:
  - Schedule code review of recent changes
  - Look for: Unfreed buffers, circular references, goroutine that never exits
  - Plan hotfix release for next deployment

- [ ] Document incident:
  - Which server(s) failed: [server names]
  - Root cause: [memory leak / timeout / disk full / etc]
  - Recovery action: Restart
  - Time to recovery: X minutes
  - Preventive: [Enable memory monitoring / set resource limits / etc]

---

**Expected RTO**: **2-5 minutes** (detect failure + restart container + health check pass)

**Expected RPO**: **0 minutes** (stateless service; no data loss)

---

---

#### **SCENARIO #4: AUTHENTICATION SYSTEM FAILURE (JWT Token Generation Fails)**

**Severity**: 🔴 **CRITICAL** (RTO: ≤5 min, RPO: 0 min)

**Symptoms / Detection**:

Alert Triggers:
- 🚨 Alert: POST /api/v1/auth/login returns 500 (token generation error)
- 🚨 Alert: Login failure rate > 50% for >1 minute
- 🚨 Alert: Error in logs: "JWT signing key not found" or "Secrets Manager timeout"
- 🚨 Alert: RBAC validation failures (role claims missing from token)
- 🚨 Alert: Token expiry misconfigured (tokens expire in 1 second instead of 1 hour)

**Impact Assessment**:

Affected Systems:
- 🔴 User login blocked (POST /api/v1/auth/login fails)
- 🔴 No new JWT tokens generated (existing tokens still valid until expiry)
- 🟠 Existing sessions can continue (until token expires; then logout forced)
- ✅ API endpoints with valid token still work (if client has unexpired token)

User Impact:
- **New users cannot login** (immediate stoppage)
- Existing logged-in users can continue working (until token expires)
- If scope: All tokens affected → Complete service unavailability
- Typical session duration: 1-24 hours → Majority users locked out within 1 hour

Business Impact:
- Revenue impact: HIGH (users cannot access platform)
- SLA breach: CRITICAL (100% login failure rate)
- Compliance risk: Potential audit trail gap (audit logs may not record failed login attempts)

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] CONFIRM AUTH FAILURE**:
   ```bash
   # Test login endpoint
   curl -X POST http://{api_server}:8080/api/v1/auth/login \
     -d '{"email": "test@example.com", "password": "test123"}' \
     -H "Content-Type: application/json"
   
   # Expected (working): 200 OK or 401 Unauthorized
   # Actual (failure): 500 Internal Server Error
   ```

2. **[1-2 min] CHECK JWT SIGNING KEY AVAILABILITY**:
   ```bash
   # Key location: Secrets Manager or environment variable
   
   # Check if key exists in secrets
   aws secretsmanager get-secret-value --secret-id scintel-jwt-secret
   # or
   echo $JWT_SECRET  # If using environment variable
   
   # Expected: Non-empty secret key (should be 256-bit or longer)
   ```

3. **[2-3 min] CHECK AUTH SERVICE LOGS**:
   ```bash
   tail -50 /var/log/scintel-api/api.log | grep -i "auth\|jwt\|token\|secret"
   
   # Look for errors like:
   # - "JWT signing key not found"
   # - "Secrets Manager connection timeout"
   # - "Token validation failed"
   # - "HMAC signature verification failed"
   ```

4. **[3-5 min] ATTEMPT AUTH SERVICE RESTART**:
   ```bash
   # Restart API server with auth service
   docker restart scintel-api-1
   
   # Wait for service to start
   sleep 20
   
   # Test login again
   curl -X POST http://scintel-api-1:8080/api/v1/auth/login \
     -d '{"email": "test@example.com", "password": "test123"}' \
     -H "Content-Type: application/json"
   # Expected: 200 OK or 401 (not 500)
   ```

5. **[5-10 min] CHECK SECRETS MANAGER / KEY STORE**:
   ```bash
   # If secrets not accessible, check connectivity
   aws secretsmanager list-secrets  # Should succeed
   
   # If timeout → Check AWS service status
   # If IAM permission denied → Check role permissions
   
   # Fallback: Regenerate JWT key (if backup available)
   # Store in: AWS Secrets Manager or environment variable
   ```

6. **[10-15 min] NOTIFY & ALLOW EXISTING SESSIONS**:
   - [ ] Slack: "🔴 CRITICAL: Authentication system recovering. New logins blocked, existing sessions continue. ETA 5-10 min."
   - [ ] Meanwhile: Users with valid tokens can continue working
   - [ ] Expected: Within 15 min, login should work again

---

**Recovery Steps** 🛠️

**Phase 1: Root Cause Diagnosis**

```bash
# 1. Check JWT secret configuration
docker exec scintel-api-1 bash -c "env | grep JWT"
# Expected: JWT_SECRET=... (non-empty)

# 2. Test token generation function locally
# If access to API source code:
go test -v ./internal/auth -run TestTokenGeneration
# Expected: PASS

# 3. Verify database (for user lookup during login)
psql -U postgres -d scintel_db -c "
  SELECT id, email, password_hash FROM users LIMIT 1;
"
# Expected: Returns user record (database OK)

# 4. Check API service startup logs
docker logs scintel-api-1 | grep -i "starting\|ready\|error"
# Look for errors during initialization
```

**Phase 2: Recovery**

**Option A: Restart Auth Service**
```bash
# Gracefully stop and restart (handles ongoing requests)
docker restart scintel-api-1
sleep 30

# Test
curl -X POST http://scintel-api-1:8080/api/v1/auth/login \
  -d '{"email": "test@warehouse.id", "password": "test123"}' \
  -H "Content-Type: application/json" | jq .
# Expected: { "access_token": "...", "refresh_token": "...", ... }
```

**Option B: Restore JWT Secret (if corrupted)**
```bash
# 1. Get backup JWT secret
aws s3 cp s3://scintel-backups/secrets/jwt-secret-backup /tmp/jwt-secret.bak

# 2. Restore to Secrets Manager
aws secretsmanager update-secret \
  --secret-id scintel-jwt-secret \
  --secret-string $(cat /tmp/jwt-secret.bak)

# 3. Restart API service to pick up updated secret
docker restart scintel-api-1

# 4. Verify
curl -X POST http://scintel-api-1:8080/api/v1/auth/login \
  -d '{"email": "test@warehouse.id", "password": "test123"}' \
  -H "Content-Type: application/json"
# Expected: 200 OK (login works)
```

**Option C: Regenerate JWT Secret (if backup unavailable)**
```bash
# WARNING: This will INVALIDATE all existing tokens!
# Only do if no alternative. Users will need to login again.

# 1. Generate new random secret (256-bit = 32 bytes)
openssl rand -hex 32 > /tmp/new-jwt-secret.txt

# 2. Store in Secrets Manager
aws secretsmanager update-secret \
  --secret-id scintel-jwt-secret \
  --secret-string $(cat /tmp/new-jwt-secret.txt)

# 3. Restart API services
docker restart scintel-api-1 scintel-api-2 scintel-api-3

# 4. Notify users: "Session expired due to security refresh. Please login again."
```

---

**Phase 3: Validation** ✅

```bash
# 1. Test login with valid credentials
curl -X POST http://{api_server}:8080/api/v1/auth/login \
  -d '{"email": "manager@warehouse.id", "password": "securepass123"}' \
  -H "Content-Type: application/json" | jq .

# Expected response:
# {
#   "access_token": "eyJhbGc...",
#   "refresh_token": "eyJhbGc...",
#   "token_type": "Bearer",
#   "expires_in": 3600,
#   "user": { "id": "...", "email": "...", "role": "manager" }
# }

# 2. Test login with invalid credentials (should return 401, not 500)
curl -X POST http://{api_server}:8080/api/v1/auth/login \
  -d '{"email": "manager@warehouse.id", "password": "wrongpass"}' \
  -H "Content-Type: application/json" | jq .

# Expected: 401 Unauthorized (not 500)

# 3. Test token refresh (optional)
curl -X POST http://{api_server}:8080/api/v1/auth/refresh \
  -d '{"refresh_token": "eyJhbGc..."}' \
  -H "Content-Type: application/json"

# Expected: 200 OK with new access_token

# 4. Test protected endpoint with new token
curl -H "Authorization: Bearer {new_token}" \
  http://{api_server}:8080/api/v1/inventory

# Expected: 200 OK with inventory data (RBAC OK)
```

---

**Phase 4: Post-Recovery** 📋

- [ ] Monitor login success rate for next 30 minutes
  ```bash
  # Check error logs for auth failures
  docker logs scintel-api-1 | tail -50 | grep -i "auth\|login"
  # Expected: No error messages, just info logs
  ```

- [ ] Notify users if JWT secret was regenerated:
  - Slack: "✅ Authentication restored. Please login again if session expired."
  - Status page: Update to "Resolved"

- [ ] Document incident:
  - Root cause: [e.g., "JWT signing key rotated; old key deleted prematurely"]
  - Recovery action: [Restart / Restore secret / Regenerate]
  - Time to recovery: X minutes
  - Preventive: [Automate secret rotation; test new secret before activation]

---

**Expected RTO**: **5-10 minutes** (identify issue + restart + test)

**Expected RPO**: **0 minutes** (authentication is stateless)

---

---

#### **SCENARIO #5: FORECASTING BATCH JOB TIMEOUT / CRASH**

**Severity**: 🟡 **MEDIUM** (RTO: ≤60 min, RPO: ≤24 hr)

**Symptoms / Detection**:

Alert Triggers:
- 🚨 Alert: Forecast batch job did not complete within 5-second window (timeout)
- 🚨 Alert: Cron job log: "ETS algorithm crashed: memory allocation error"
- 🚨 Alert: No new forecasts generated in last 24 hours
- 🚨 Alert: Forecast API returns stale data (generated_at = 24 hours ago)
- 🚨 Alert: CPU spike on dedicated forecasting server (1000+ goroutines processing)

**Impact Assessment** (from PRD section 4.4):

Affected Systems:
- 🟡 Daily forecast generation blocked (runs @ 2 AM UTC per PRD)
- 🟡 Managers see 24-hour-old forecasts instead of fresh predictions
- 🟡 Confidence intervals outdated (demand patterns may have changed)
- ✅ Core inventory operations unaffected (stock transactions continue)
- ✅ Anomaly detection works (uses independent algorithm)

User Impact:
- Warehouse managers use yesterday's forecast to make today's decisions
- Prediction accuracy decreases (older data = less relevant)
- Demand forecasts become stale by end of day
- Acceptable for most use cases (tolerance: 24 hours)

Business Impact:
- Operational risk: Managers make decisions with outdated data
- Inventory misalignment risk: Stock levels don't match predicted demand
- No immediate revenue impact (forecast is supporting tool, not critical path)

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] VERIFY BATCH JOB STATUS**:
   ```bash
   # Check cron job log (runs daily @ 2 AM UTC)
   grep "forecast.*batch\|ETS\|ARIMA" /var/log/cron
   
   # Check if process currently running
   ps aux | grep -i "forecast\|batch"
   
   # Check database for latest forecast generation
   psql -U postgres -d scintel_db -c "
     SELECT COUNT(*) as forecast_count, MAX(generated_at) as latest_forecast
     FROM forecasts;
   "
   # Expected: Should show recent generated_at (within last 24 hours)
   ```

2. **[1-2 min] ASSESS FORECAST STALENESS**:
   ```bash
   # Query API to see forecast age
   curl -H "Authorization: Bearer {token}" \
     http://{api_server}:8080/api/v1/forecasts?sku=SKU-NIKE-001 | jq '.generated_at'
   
   # If generated_at < 24 hours → Minor issue (not urgent)
   # If generated_at = null or error → Batch failed, need immediate action
   ```

3. **[2-5 min] REVIEW BATCH JOB LOGS**:
   ```bash
   # For Golang forecast service
   tail -100 /var/log/scintel-forecasting/forecast-batch.log | grep -i "error\|failed\|timeout"
   
   # Look for errors like:
   # - "context deadline exceeded" (timeout)
   # - "out of memory" (insufficient resources)
   # - "SQL error: connection lost" (database disconnect)
   # - "algorithm crash: NaN detected" (algorithm error)
   ```

4. **[5-10 min] MANUAL TRIGGER FORECAST GENERATION** (If automatic failed):
   ```bash
   # Option 1: Call admin API endpoint
   curl -X POST http://{api_server}:8080/api/v1/admin/forecasts/generate \
     -H "Authorization: Bearer {admin_token}" \
     -d '{"warehouse_id": "wh-jakarta-001", "algorithm": "exponential_smoothing"}'
   
   # Option 2: Run batch job manually
   docker exec scintel-forecasting ./cmd/forecast-batch/main \
     --config /etc/scintel/forecast-config.yaml \
     --debug
   
   # Monitor: Should complete within 5 seconds for 1000 SKUs
   ```

5. **[10-15 min] VERIFY FORECAST FRESHNESS**:
   ```bash
   # Query latest forecast timestamp
   curl -H "Authorization: Bearer {token}" \
     http://{api_server}:8080/api/v1/forecasts?sku=SKU-NIKE-001 | jq '.generated_at'
   
   # Expected: Should be current timestamp (or within last 1 hour if delayed)
   ```

---

**Recovery Steps** 🛠️

**Phase 1: Diagnosis**

```bash
# 1. Check forecast service health
docker logs scintel-forecasting | tail -50 | grep -i "error\|warning"

# 2. Verify database connectivity
psql -U postgres -d scintel_db -c "SELECT COUNT(*) FROM transactions LIMIT 1;"
# Expected: Returns number (database OK)

# 3. Check system resources on forecasting server
free -h  # Memory usage
df -h    # Disk space
top -b -n 1 | head -20  # CPU / processes

# 4. Review forecast accuracy metrics (to see if algorithm working)
psql -U postgres -d scintel_db -c "
  SELECT 
    sku, 
    algorithm, 
    mape, 
    accuracy_score,
    generated_at
  FROM forecasts 
  ORDER BY generated_at DESC 
  LIMIT 5;
"
# Expected: Recent forecasts with valid MAPE/accuracy scores
```

**Phase 2: Recovery**

**Option A: Restart Forecasting Service**
```bash
# Stop batch job if running
pkill -f "forecast.*batch"

# Restart service
systemctl restart scintel-forecasting
sleep 10

# Manually trigger forecast generation
curl -X POST http://{api_server}:8080/api/v1/admin/forecasts/generate \
  -H "Authorization: Bearer {admin_token}"

# Wait for completion (monitor logs)
tail -f /var/log/scintel-forecasting/forecast-batch.log
# Expected: "Forecast generation complete. Generated 4250 forecasts in 3.5 seconds"
```

**Option B: Run Single-SKU Forecast (if batch too slow)**
```bash
# If 1000+ SKU batch times out, try:
# 1. Reduce scope (e.g., only 100 most-active SKUs)
# 2. Increase timeout (if infrastructure can handle)
# 3. Parallel processing (run multiple threads)

# Test with single SKU first
curl -X POST http://{api_server}:8080/api/v1/admin/forecasts/generate \
  -H "Authorization: Bearer {admin_token}" \
  -d '{
    "sku": "SKU-NIKE-001",
    "algorithm": "exponential_smoothing",
    "warehouse_id": "wh-jakarta-001"
  }'

# If succeeds → Batch job needs optimization (algorithmic or infrastructure)
# If fails → Issue with algorithm or database query
```

**Option C: Increase Batch Job Resources (if timeout)**
```bash
# If batch times out (>5 seconds for 1000 SKUs), increase resources:

# 1. Increase memory allocation
docker update --memory=4g scintel-forecasting

# 2. Increase timeout in cron job
# In /etc/cron.d/scintel-forecasting:
# Old: 0 2 * * * timeout 5 /app/forecast-batch
# New: 0 2 * * * timeout 30 /app/forecast-batch  # Allow 30 seconds

# 3. Optimize algorithm (Phase 1.5+)
# - Cache historical data to avoid re-query
# - Use GPU acceleration for matrix operations
# - Implement algorithm selection (pick fast algorithm if data large)
```

**Option D: Restore from Backup Forecast (if data corrupted)**
```bash
# If forecast data corrupted (NaN values, invalid confidence):
# Restore from database backup and re-run

# 1. Query valid historical forecasts
psql -U postgres -d scintel_db -c "
  SELECT * FROM forecasts 
  WHERE generated_at = (
    SELECT MAX(generated_at) FROM forecasts 
    WHERE generated_at < NOW() - INTERVAL '24 hours'
  )
  LIMIT 10;
"

# 2. If corrupted, restore from backup
# (See Scenario #1: Database Backup Restore)
```

---

**Phase 3: Validation** ✅

```bash
# 1. Verify forecast was generated recently
curl -H "Authorization: Bearer {token}" \
  http://{api_server}:8080/api/v1/forecasts?sku=SKU-NIKE-001 | jq '.generated_at'

# Expected: Timestamp within last 2 hours

# 2. Check forecast data quality
curl -H "Authorization: Bearer {token}" \
  http://{api_server}:8080/api/v1/forecasts?sku=SKU-NIKE-001 | jq '.model_info'

# Expected:
# {
#   "algorithm_used": "exponential_smoothing",
#   "accuracy_metrics": {
#     "mape": 0.082,
#     "rmse": 12.5,
#     "mae": 10.2
#   }
# }

# 3. Verify confidence intervals valid
curl -H "Authorization: Bearer {token}" \
  http://{api_server}:8080/api/v1/forecasts?sku=SKU-NIKE-001 | jq '.forecasts[0]'

# Expected: confidence between 0.90-0.99, lower_bound < upper_bound
# Invalid: NaN, negative values, confidence > 1.0
```

---

**Phase 4: Post-Recovery** 📋

- [ ] Monitor next batch job run (tomorrow @ 2 AM)
  - [ ] Set alarm to check forecast completion
  - [ ] Verify email alert sent if successful

- [ ] Review batch job performance
  - [ ] If timeout: Investigate optimization opportunities
  - [ ] If data issues: Check data quality in transactions table
  - [ ] If algorithm crash: Review error logs for specific algorithm bug

- [ ] Document incident:
  - Root cause: [e.g., "Database query for historical data took 8 seconds; timeout set to 5"]
  - Recovery: [Restart service / Manual trigger / Increase timeout]
  - Preventive: [Monitor batch job duration; alert if > 3 seconds]

---

**Expected RTO**: **10-60 minutes** (depends on recovery method; manual trigger is quickest)

**Expected RPO**: **24 hours** (acceptable; forecasts generated daily)

---

---

#### **SCENARIO #6: WEBHOOK DELIVERY QUEUE BACKUP / RETRY FAILURE**

**Severity**: 🟠 **HIGH** (RTO: ≤30 min, RPO: ≤1 hr)

**Symptoms / Detection**:

Alert Triggers:
- 🚨 Alert: Webhook delivery queue size > 10,000 events (backed up)
- 🚨 Alert: Failed webhooks in database table webhook_failed_events > 1000 (exponential backoff exhausted)
- 🚨 Alert: Dead-letter queue backing up (delivery retries all failed)
- 🚨 Alert: Error in logs: "Webhook delivery failed after 5 retries" repeated >100 times/min
- 🚨 Alert: External subscriber reports not receiving inventory change notifications

**Impact Assessment** (from PRD section 4.7, LOGIC_FLOW section 7):

Affected Systems:
- 🟡 Webhook subscribers don't receive events (inventory changes, anomaly alerts)
- 🟡 External integrations get stale data (not notified of stock changes)
- ✅ Core system operations unaffected (internal inventory updates continue)
- ✅ Database still logging events (can be replayed later)

User Impact:
- Warehouse managers using external dashboard miss inventory updates
- Third-party apps (inventory dashboard, ERP integrations) out of sync with actual stock
- Alerts delivered late or not at all (anomaly notifications delayed)
- Customers placing orders see stale availability (negative experience)

Business Impact:
- Partner integrations fail silently (no error indication)
- Order fulfillment delayed (external system waiting for notification)
- SLA breach if webhook delivery is contractual requirement
- Low immediate revenue impact (inventory operations continue, just notifications delayed)

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] CONFIRM WEBHOOK QUEUE BACKUP**:
   ```bash
   # Check database for pending webhooks
   psql -U postgres -d scintel_db -c "
     SELECT 
       COUNT(*) as pending_count,
       MAX(created_at) as oldest_pending
     FROM webhook_events 
     WHERE status = 'pending';
   "
   # Expected: Shows number of pending events
   
   # Check failed webhooks
   psql -U postgres -d scintel_db -c "
     SELECT COUNT(*) as failed_count FROM webhook_failed_events WHERE retry_count >= 5;
   "
   # Expected: Number of exhausted retries (max 5 per PRD)
   ```

2. **[1-2 min] IDENTIFY ROOT CAUSE OF DELIVERY FAILURES**:
   ```bash
   # Check webhook delivery logs
   tail -50 /var/log/scintel-api/webhook-delivery.log | grep -i "error\|failed\|timeout"
   
   # Look for patterns:
   # - All deliveries to same endpoint failing → Subscriber's server down
   # - Random failures across endpoints → Network issue
   # - Connection refused → Subscriber server not accepting connections
   # - Timeout → Subscriber's server slow
   ```

3. **[2-5 min] CHECK WEBHOOK SUBSCRIBER STATUS** (Are they reachable?):
   ```bash
   # List active webhook subscriptions
   psql -U postgres -d scintel_db -c "
     SELECT 
       id,
       subscriber_url,
       status,
       last_delivery_at,
       last_error
     FROM webhook_subscriptions
     ORDER BY last_delivery_at DESC
     LIMIT 10;
   "
   
   # Test connectivity to top 3 subscribers
   for subscriber_url in $(psql -U postgres -d scintel_db -c \
       "SELECT subscriber_url FROM webhook_subscriptions WHERE status='active' LIMIT 3" -t); do
     
     status=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$subscriber_url" \
       -d '{"test": "connection"}' \
       -H "Content-Type: application/json" \
       --max-time 5)
     
     echo "Subscriber $subscriber_url: HTTP $status"
   done
   
   # Expected: 200-299 (success) or 4xx (subscriber issue), not timeout
   ```

4. **[5-10 min] ASSESS WEBHOOK SERVICE HEALTH**:
   ```bash
   # Check webhook processor goroutines
   curl -s http://{api_server}:6060/debug/pprof/goroutine | grep webhook
   
   # Check error rate
   docker logs scintel-api-1 | tail -30 | grep -i "webhook"
   
   # Check if webhook service running
   ps aux | grep webhook
   systemctl status scintel-webhook-processor
   ```

5. **[10-15 min] NOTIFY SUBSCRIBERS & DECISION POINT**:
   - Slack: "🟠 HIGH: Webhook delivery queue backed up. Events queued for retry. ETA recovery 30 min."
   - If critical subscriber affected: Call to notify them (e.g., major ERP integration)
   - Decision: Immediately retry failed webhooks OR wait for delivery window retry?

---

**Recovery Steps** 🛠️

**Phase 1: Quick Recovery (Retry Immediately)**

```bash
# Option 1: Restart webhook processor service
systemctl restart scintel-webhook-processor

# Monitor delivery queue
psql -U postgres -d scintel_db -c "
  SELECT COUNT(*) FROM webhook_events WHERE status='pending';
"
# Should decrease as events are delivered

# Option 2: Trigger manual retry of failed webhooks
curl -X POST http://{api_server}:8080/api/v1/admin/webhooks/retry-failed \
  -H "Authorization: Bearer {admin_token}" \
  -d '{"retry_max_age_hours": 24}'  # Retry events from last 24 hours

# Expected: Returns number of events queued for retry
# Monitor: {"retried": 500, "skipped": 200}  # 200 already retried >5 times
```

**Phase 2: Investigation & Issue Resolution**

```bash
# 1. If subscriber server is down, notify them
# 2. If network issue, check firewall rules
firewall-cmd --list-all | grep webhook
# 3. If DNS issue, verify subscriber URL resolution
nslookup subscriber.example.com

# 4. For each failed subscriber, decide:
# - Temporarily disable subscription until they recover
# - Extend retry window (instead of giving up after 5 retries)
# - Switch to manual/polling mode (subscriber polls instead of receiving push)
```

**Phase 3: Validate & Resume**

```bash
# 1. Monitor webhook delivery queue decreasing
watch -n 10 "psql -U postgres -d scintel_db -c \
  'SELECT COUNT(*) FROM webhook_events WHERE status=\"pending\";'"

# 2. Verify delivered webhooks reaching subscribers
# Ask partner: "Are you receiving inventory update notifications?"
# Check partner dashboard: Last event timestamp should be fresh

# 3. Clear dead-letter queue if all retries exhausted
psql -U postgres -d scintel_db -c "
  DELETE FROM webhook_failed_events 
  WHERE retry_count >= 5 AND created_at < NOW() - INTERVAL '48 hours';
"
# This frees up space; events older than 48 hours discarded
```

---

**Phase 4: Post-Recovery** 📋

- [ ] Webhook delivery metrics normalized
  ```bash
  # Check delivery success rate (target: >99%)
  psql -U postgres -d scintel_db -c "
    SELECT 
      COUNT(CASE WHEN status='delivered' THEN 1 END)::float / COUNT(*) * 100 as success_rate
    FROM webhook_events
    WHERE created_at > NOW() - INTERVAL '1 hour';
  "
  ```

- [ ] Notify affected subscribers of recovery
  - Email: "Webhook delivery issue resolved. All pending events delivered. Thanks for patience."

- [ ] Document incident:
  - Root cause: [e.g., "Subscriber's server was down for maintenance"]
  - Recovery: [Restart service + retry]
  - Time: X minutes
  - Preventive: [Setup subscriber uptime monitoring; alert if endpoint down]

---

**Expected RTO**: **5-30 minutes** (restart + retry processing)

**Expected RPO**: **1 hour** (events logged in database; can replay)

---

---

#### **SCENARIO #7: ANOMALY DETECTION ALGORITHM FAILURE**

**Severity**: 🟡 **MEDIUM** (RTO: ≤60 min, RPO: N/A)

**Symptoms / Detection**:

Alert Triggers:
- 🚨 Alert: Anomaly detection job did not complete within 2-minute window
- 🚨 Alert: Error logs: "Z-score calculation failed: NaN result" OR "IQR method returned invalid range"
- 🚨 Alert: No new anomalies detected in last 24 hours (anomaly_created_at stale)
- 🚨 Alert: Webhook anomaly events not firing (downstream integrations quiet)
- 🚨 Alert: Algorithm service CPU spike then crash (goroutine panic)

**Impact Assessment**:

Affected Systems:
- 🟡 Anomaly detection paused (unusual patterns not flagged)
- 🟡 Alert subscribers miss notifications (stock spikes, system errors undetected)
- ✅ Inventory operations unaffected (stock transactions continue normally)
- ✅ Forecasting unaffected (independent service)
- ✅ Core business operations continue

User Impact:
- Warehouse managers miss unusual stock movements (e.g., unexpected drop of 1000 units)
- Fraud/system errors potentially slip by undetected
- Acceptable delay: 24 hours (anomalies generated daily)
- **No immediate operational impact** (supporting feature, not core path)

Business Impact:
- Risk: Inventory discrepancies missed
- Low immediate revenue impact (alerts are detection layer, not transactional)
- Compliance: May need to manually audit anomalies during recovery period

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] CONFIRM ANOMALY SERVICE STATUS**:
   ```bash
   # Check if service running
   systemctl status scintel-anomaly-detection
   
   # Check latest anomaly record
   psql -U postgres -d scintel_db -c "
     SELECT id, sku, severity, created_at, z_score, status
     FROM anomalies
     ORDER BY created_at DESC
     LIMIT 1;
   "
   # Expected: Recent timestamp (within last 1-2 hours)
   ```

2. **[1-2 min] REVIEW ERROR LOGS**:
   ```bash
   tail -100 /var/log/scintel-api/anomaly-detection.log | grep -i "error\|failed\|panic\|nan"
   
   # Look for:
   # - "math: result is NaN"
   # - "divide by zero"
   # - "invalid sample size"
   # - "algorithm panic"
   ```

3. **[2-5 min] CHECK DATA QUALITY IN DATABASE**:
   ```bash
   # Verify transaction history exists (algorithm needs historical data)
   psql -U postgres -d scintel_db -c "
     SELECT 
       COUNT(*) as transaction_count,
       DATE(MIN(created_at)) as oldest_date,
       DATE(MAX(created_at)) as newest_date
     FROM transactions;
   "
   # Expected: Multiple records spanning weeks/months
   
   # Check for data anomalies (NULL values, invalid quantities)
   psql -U postgres -d scintel_db -c "
     SELECT COUNT(*) FROM transactions 
     WHERE quantity IS NULL OR quantity < 0;
   "
   # Expected: 0 (no corrupted data)
   ```

4. **[5-10 min] RESTART ANOMALY SERVICE**:
   ```bash
   systemctl restart scintel-anomaly-detection
   sleep 15
   
   # Trigger manual anomaly detection
   curl -X POST http://{api_server}:8080/api/v1/admin/anomalies/detect \
     -H "Authorization: Bearer {admin_token}"
   
   # Monitor logs for success
   tail -f /var/log/scintel-api/anomaly-detection.log
   # Expected: "Anomaly detection complete. Detected 5 anomalies in 1.2 seconds"
   ```

5. **[10-15 min] VERIFY ANOMALIES DETECTED**:
   ```bash
   # Query newly detected anomalies
   psql -U postgres -d scintel_db -c "
     SELECT id, sku, severity, z_score, created_at
     FROM anomalies
     WHERE created_at > NOW() - INTERVAL '10 minutes'
     ORDER BY created_at DESC;
   "
   # Expected: New anomaly records (if any anomalies exist in current data)
   ```

---

**Recovery Steps** 🛠️

**Phase 1: Diagnosis**

```bash
# 1. Verify algorithm implementation
# Check if algorithm selection correct:
psql -U postgres -d scintel_db -c "
  SELECT 
    algorithm_type,
    COUNT(*) as count,
    AVG(z_score) as avg_zscore,
    STDDEV(z_score) as stddev_zscore
  FROM anomalies
  GROUP BY algorithm_type;
"
# Expected: Z-scores distributed around 0 with reasonable stddev (~1-2)
# Invalid: All NaN, all zeros, or stddev >> 100

# 2. Test algorithm with sample data (if algorithm parameterizable)
curl -X POST http://{api_server}:8080/api/v1/admin/anomalies/test \
  -H "Authorization: Bearer {admin_token}" \
  -d '{
    "algorithm": "z_score",
    "test_values": [1, 2, 3, 4, 5, 100],  # 100 should be flagged
    "threshold": 3.0
  }' | jq .

# Expected: 
# {
#   "test_results": [
#     { "value": 1, "z_score": -1.46, "is_anomaly": false },
#     ...
#     { "value": 100, "z_score": 3.92, "is_anomaly": true }
#   ]
# }

# 3. Check if algorithm parameters out of bounds
# Common issues:
# - Sample size < 10 (insufficient data for statistics)
# - Z-score threshold misconfigured (too high/low)
# - IQR method: Q1/Q3 calculation wrong
```

**Phase 2: Recovery**

**Option A: Restart Service with Clean State**
```bash
# Stop service
systemctl stop scintel-anomaly-detection

# Clear any cached algorithm state (if cached)
redis-cli DEL "anomaly:algorithm:*"

# Restart
systemctl start scintel-anomaly-detection

# Trigger reprocessing
curl -X POST http://{api_server}:8080/api/v1/admin/anomalies/reprocess \
  -H "Authorization: Bearer {admin_token}" \
  -d '{"days_back": 7}'  # Reprocess last 7 days

# Monitor
tail -f /var/log/scintel-api/anomaly-detection.log | grep -i "anomaly\|complete\|error"
```

**Option B: Algorithm Parameter Adjustment**
```bash
# If algorithm crashing due to parameter issues:
# Update anomaly detection config

cat > /etc/scintel/anomaly-config.yaml << EOF
anomaly_detection:
  algorithm: "z_score"  # or "iqr"
  z_score_threshold: 3.0  # Default: flag if |z| > 3
  iqr_multiplier: 1.5  # Default: flag if outlier > Q3 + 1.5*IQR
  min_sample_size: 20  # Require >= 20 historical points
  lookback_days: 90  # Analyze last 90 days
  exclude_outliers: false  # Don't remove detected outliers from calculation
EOF

# Restart service
systemctl restart scintel-anomaly-detection

# Re-run detection
curl -X POST http://{api_server}:8080/api/v1/admin/anomalies/detect \
  -H "Authorization: Bearer {admin_token}"
```

**Option C: Skip Algorithm, Use Threshold-Based Detection (Fallback)**
```bash
# If algorithm broken beyond quick fix:
# Fall back to simple threshold-based alerts

curl -X POST http://{api_server}:8080/api/v1/admin/anomalies/detect-threshold \
  -H "Authorization: Bearer {admin_token}" \
  -d '{
    "method": "threshold",
    "stock_drop_threshold": 50,  # Alert if stock drops >50 units/hour
    "unusual_order_qty": 1000   # Alert if single order > 1000 units
  }'

# This is temporary; real algorithm should be fixed
```

---

**Phase 3: Validation** ✅

```bash
# 1. Verify anomalies detected correctly
curl -H "Authorization: Bearer {token}" \
  http://{api_server}:8080/api/v1/anomalies \
  -d '{"severity": "high", "limit": 10}' | jq '.anomalies[0]'

# Expected: 
# {
#   "id": "anom-xyz",
#   "sku": "SKU-NIKE-001",
#   "warehouse_id": "wh-jakarta-001",
#   "type": "unusual_stock_drop",
#   "severity": "high",
#   "z_score": 3.5,
#   "detected_at": "2026-07-18T10:30:00Z"
# }

# 2. Verify webhook events fired
psql -U postgres -d scintel_db -c "
  SELECT 
    event_type,
    COUNT(*) as count,
    MAX(created_at) as latest
  FROM webhook_events
  WHERE event_type LIKE 'anomaly%' AND created_at > NOW() - INTERVAL '10 minutes'
  GROUP BY event_type;
"

# 3. Check subscribers received notifications
# Ask warehouse manager: "Did you receive anomaly alerts?"
```

---

**Expected RTO**: **10-30 minutes** (restart + re-detection)

**Expected RPO**: **N/A** (anomalies are notifications, not transactional data)

---

---

#### **SCENARIO #9: DATABASE REPLICA LAG / SYNC FAILURE**

**Severity**: 🟡 **MEDIUM** (RTO: ≤30 min, RPO: ≤5 min)

**Symptoms / Detection** (from PRD 4.1, TDD section 3):

Alert Triggers:
- 🚨 Alert: Replication lag > 5 minutes (standby behind primary)
- 🚨 Alert: Read replica connection fails (read-only queries fail)
- 🚨 Alert: WAL receiver status shows "streaming failed" or "connection lost"
- 🚨 Alert: Primary database replication slot reports "inactive" (no one subscribing)
- 🚨 Alert: Standby lag_bytes increasing continuously (not catching up)

**Impact Assessment**:

Affected Systems:
- 🟡 Read-only queries (analytics, reporting) may see stale data (up to 5+ min old)
- 🟡 Failover capability compromised (standby not caught up → potential data loss)
- 🟡 High-availability queries degraded (can't distribute read load)
- ✅ Write operations unaffected (all to primary)
- ✅ Real-time inventory queries OK if routing to primary

User Impact:
- Analytics dashboard shows data from 5-10 minutes ago (acceptable delay)
- Reports run slower (queries hitting primary instead of read replica)
- Warehouse operations continue (reads hitting primary, just slower)

Business Impact:
- If **not fixed**: Failover will lose recent data (RPO breach)
- If **left unfixed**: Replica becomes worthless for disaster recovery
- **Immediate impact**: Medium (reporting delayed, but ops continue)

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] VERIFY REPLICATION STATUS**:
   ```bash
   # On PRIMARY database
   psql -U postgres -c "
     SELECT 
       slot_name,
       slot_type,
       active,
       restart_lsn,
       confirmed_flush_lsn
     FROM pg_replication_slots;
   "
   # Expected: slot_type='physical', active=true
   
   # On STANDBY database
   psql -U postgres -c "
     SELECT 
       now() - pg_last_xact_replay_time() as replication_lag,
       pg_is_in_recovery(),
       pg_last_wal_receive_lsn(),
       pg_last_wal_replay_lsn();
   "
   # Expected: replication_lag close to 0 seconds
   ```

2. **[1-2 min] CHECK NETWORK CONNECTIVITY**:
   ```bash
   # Verify network between primary and standby
   ping -c 5 {standby_host}
   
   # Check if port 5432 (PostgreSQL) open
   telnet {standby_host} 5432
   
   # Check network bandwidth
   iftop -i eth0 | head -20  # Show network traffic
   ```

3. **[2-5 min] REVIEW POSTGRESQL LOGS**:
   ```bash
   # On PRIMARY:
   tail -50 /var/log/postgresql/postgresql.log | grep -i "wal\|replication\|error"
   
   # Look for:
   # - "connection lost with the standby"
   # - "replication slot became inactive"
   # - "WAL sender error"
   
   # On STANDBY:
   tail -50 /var/log/postgresql/postgresql.log | grep -i "wal receiver\|error\|connection"
   ```

4. **[5-10 min] CHECK DISK SPACE ON PRIMARY**:
   ```bash
   # WAL files taking up space?
   du -sh /var/lib/postgresql/pg_wal/
   
   # If WAL directory > 10GB, replication backlog
   ls -lh /var/lib/postgresql/pg_wal/ | tail -20
   ```

5. **[10-15 min] DECISION POINT**:
   - If lag < 1 min: Continue monitoring (acceptable)
   - If lag 1-5 min: Investigate cause, but not urgent
   - If lag > 10 min: Needs immediate action (failover compromised)

---

**Recovery Steps** 🛠️

**Phase 1: Identify Root Cause**

```bash
# 1. Is standby consuming WAL?
psql -U postgres -c "
  SELECT 
    slot_name,
    active,
    restart_lsn,
    confirmed_flush_lsn,
    slot_type
  FROM pg_replication_slots;
"

# If active=false OR confirmed_flush_lsn = restart_lsn:
# → Standby not reading WAL; connection broken

# 2. Is standby catching up?
psql -U postgres -c "
  SELECT 
    now(),
    pg_last_xact_replay_time(),
    NOW() - pg_last_xact_replay_time() as lag
  FROM pg_stat_replication;
"

# 3. Check write load on primary
psql -U postgres -c "
  SELECT 
    datname,
    tup_inserted,
    tup_updated,
    tup_deleted
  FROM pg_stat_user_tables
  ORDER BY (tup_inserted + tup_updated + tup_deleted) DESC
  LIMIT 10;
"
# If huge insert/update rate: Primary outpacing standby's ability to catch up
```

**Phase 2: Recovery**

**Option A: Restart Standby WAL Receiver**
```bash
# If WAL receiver process crashed:

# Stop standby
systemctl stop postgresql

# Restart
systemctl start postgresql

# Verify replication resuming
sleep 10
psql -U postgres -c "SELECT now() - pg_last_xact_replay_time() as lag;"
# Expected: Lag should decrease over time
```

**Option B: Rebuild Replication Slot**
```bash
# If replication slot corrupted:

# On PRIMARY, drop the slot
psql -U postgres -c "
  SELECT pg_drop_replication_slot('standby_slot_1');
"

# Recreate slot
psql -U postgres -c "
  SELECT * FROM pg_create_physical_replication_slot('standby_slot_1');
"

# On STANDBY, restart to reconnect
systemctl restart postgresql

# Verify reconnection
psql -U postgres -c "
  SELECT now() - pg_last_xact_replay_time() as lag;
"
# Expected: Standby catches up within few seconds
```

**Option C: Full Resync (If Slot Beyond Recovery)**
```bash
# If standby too far behind, need fresh copy:

# 1. Stop standby
ssh {standby_host}
systemctl stop postgresql

# 2. Remove standby data directory
rm -rf /var/lib/postgresql/main

# 3. Create new empty cluster
sudo -u postgres /usr/lib/postgresql/15/bin/initdb -D /var/lib/postgresql/main

# 4. From PRIMARY, take backup
pg_basebackup -h {primary_host} \
  -D /var/lib/postgresql/main \
  -U postgres \
  -v -P \
  -W  # Prompt for password

# 5. Start standby
systemctl start postgresql

# 6. Verify replication
psql -U postgres -c "SELECT now() - pg_last_xact_replay_time() as lag;"
# Expected: Should be caught up within seconds
```

---

**Phase 3: Validation** ✅

```bash
# 1. Replication lag minimal
psql -U postgres -c "
  SELECT 
    now() - pg_last_xact_replay_time() as replication_lag_seconds,
    CASE 
      WHEN now() - pg_last_xact_replay_time() < '1 second'::interval THEN '✅ EXCELLENT'
      WHEN now() - pg_last_xact_replay_time() < '10 second'::interval THEN '✅ GOOD'
      WHEN now() - pg_last_xact_replay_time() < '60 second'::interval THEN '⚠️ ACCEPTABLE'
      ELSE '❌ TOO HIGH'
    END as status;
"

# 2. Slot actively consuming WAL
psql -U postgres -c "
  SELECT 
    slot_name,
    active,
    restart_lsn,
    confirmed_flush_lsn,
    (confirmed_flush_lsn IS NOT NULL) as slot_progressing
  FROM pg_replication_slots;
"
# Expected: active=true, confirmed_flush_lsn progressing

# 3. Read replica queries working
psql -h {standby_host} -U postgres -c "SELECT COUNT(*) FROM inventory;"
# Expected: Returns count (read-only query works)
```

---

**Expected RTO**: **5-30 minutes** (restart standby + resync if needed)

**Expected RPO**: **5 minutes** (lag acceptable if < 5 min)

---

---

#### **SCENARIO #11: DATA CORRUPTION**

**Severity**: 🔴 **CRITICAL** (RTO: ≤30 min, RPO: ≤5 min)

**Symptoms / Detection**:

Alert Triggers:
- 🚨 Alert: Database consistency check fails (`REINDEX`, `ANALYZE` errors)
- 🚨 Alert: Query returns impossible values (e.g., negative inventory, future dates)
- 🚨 Alert: Inventory totals don't match transaction sums (accounting error)
- 🚨 Alert: Foreign key constraint violation (broken referential integrity)
- 🚨 Alert: Duplicate primary keys detected (should be unique)
- 🚨 Alert: Data type violation (string in numeric column, etc.)

**Impact Assessment**:

Affected Systems:
- 🔴 Inventory data integrity compromised (cannot trust stock numbers)
- 🔴 Forecasting based on corrupted data → inaccurate predictions
- 🔴 Financial reconciliation impossible (audit log unreliable)
- 🔴 Compliance risk (corrupted audit trail)

User Impact:
- **Critical**: Cannot rely on inventory data for business decisions
- Stock levels may be wrong by 10%, 50%, or completely false
- Forecasts garbage in, garbage out
- Must audit manually or restore backup

Business Impact:
- **Severe**: Loss of data trust (core asset)
- May need to manually verify inventory (warehouse count)
- Regulatory compliance issue (audit trail corrupted)
- Revenue impact: Cannot fulfill orders confidently

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] CONFIRM CORRUPTION**:
   ```bash
   # Run database consistency check
   psql -U postgres -d scintel_db -c "
     -- Check for constraint violations
     SELECT COUNT(*) FROM inventory 
     WHERE quantity < 0;  -- Negative stock impossible
   "
   
   psql -U postgres -d scintel_db -c "
     -- Check inventory vs transaction sum
     WITH inv_sum AS (
       SELECT 
         sku,
         warehouse_id,
         SUM(quantity_delta) as calculated_quantity
       FROM transactions
       WHERE created_at > NOW() - INTERVAL '30 days'
       GROUP BY sku, warehouse_id
     )
     SELECT 
       i.sku,
       i.warehouse_id,
       i.quantity as actual_inventory,
       COALESCE(s.calculated_quantity, 0) as calculated_inventory,
       (i.quantity - COALESCE(s.calculated_quantity, 0)) as discrepancy
     FROM inventory i
     LEFT JOIN inv_sum s ON i.sku = s.sku AND i.warehouse_id = s.warehouse_id
     WHERE ABS(i.quantity - COALESCE(s.calculated_quantity, 0)) > 10
     LIMIT 20;
   "
   
   # If results show significant discrepancies → CORRUPTION CONFIRMED
   ```

2. **[1-2 min] ASSESS SCOPE OF CORRUPTION**:
   ```bash
   # Which tables affected?
   SELECT 
     'inventory'::text as table_name,
     COUNT(*) as problem_count
   FROM inventory WHERE quantity < 0
   UNION ALL
   SELECT 
     'transactions'::text,
     COUNT(*) FROM transactions WHERE quantity_delta IS NULL
   UNION ALL
   SELECT 
     'forecasts'::text,
     COUNT(*) FROM forecasts 
     WHERE confidence < 0 OR confidence > 1
   ORDER BY problem_count DESC;
   ```

3. **[2-5 min] DECIDE: RESTORE VS. FIX**:
   ```bash
   # If < 100 corrupted rows → Try to fix specific records
   # If > 100 corrupted rows → Restore from backup
   # If entire table corrupted → Restore from backup
   
   # Quick assessment:
   psql -U postgres -d scintel_db -c "
     SELECT 
       CASE 
         WHEN (SELECT COUNT(*) FROM inventory WHERE quantity < 0) > 0 
           THEN '❌ RESTORE NEEDED (inventory corrupted)'
         WHEN (SELECT COUNT(*) FROM transactions WHERE quantity_delta IS NULL) > 0 
           THEN '❌ RESTORE NEEDED (transactions corrupted)'
         WHEN (SELECT COUNT(*) FROM forecasts WHERE confidence < 0 OR confidence > 1) > 0 
           THEN '⚠️ FIX POSSIBLE (minor corruption)'
         ELSE '✅ NO CORRUPTION DETECTED'
       END as verdict;
   "
   ```

4. **[5-10 min] TAKE IMMEDIATE BACKUP** (Before any changes):
   ```bash
   # Backup corrupted state (for investigation)
   pg_dump --format=custom \
     --file=/var/lib/postgresql/backups/scintel_corrupted_$(date +%s).backup \
     scintel_db
   
   # Mark timestamp so team knows corruption was discovered at this time
   echo "Corruption discovered at $(date)" > /var/lib/postgresql/backups/CORRUPTION_MARK
   ```

5. **[10-15 min] INITIATE RECOVERY**:
   ```bash
   # Stop application servers (prevent further writes)
   docker stop scintel-api-1 scintel-api-2 scintel-api-3
   
   # Restore from backup
   # (See Scenario #1 for detailed restore steps)
   ```

---

**Recovery Steps** 🛠️

**Phase 1: Investigate Root Cause (After initial restore)**

```bash
# 1. How did corruption happen?
# - Bad SQL query (e.g., UPDATE without WHERE clause)?
# - Hardware failure (bit flip)?
# - Concurrent transaction race condition?
# - Backup itself was corrupted?

# Check transaction logs around corruption time:
tail -1000 /var/log/postgresql/postgresql.log | grep -i "error\|corrupt\|constraint"

# Check application logs:
tail -1000 /var/log/scintel-api/api.log | grep -i "update\|delete\|error"

# 2. Examine corrupted values
# If inventory negative: trace transactions leading to negative balance
psql -U postgres -d scintel_db -c "
  SELECT 
    t.transaction_id,
    t.sku,
    t.warehouse_id,
    t.quantity_delta,
    t.created_at,
    t.created_by
  FROM transactions t
  WHERE t.sku = 'SKU-CORRUPTED' AND t.warehouse_id = 'wh-corrupted'
  ORDER BY t.created_at DESC
  LIMIT 20;
"

# 3. If caused by buggy application update:
# Review recent code changes that touched inventory update logic
```

**Phase 2: Recovery**

**Option A: Restore from Clean Backup (Recommended)**
```bash
# See Scenario #1 (Database Failure) for full restore procedure
# Key steps:
# 1. Stop API servers
# 2. Restore database from backup
# 3. Replay transaction logs to recover missing transactions
# 4. Restart API servers

# Quick restore:
systemctl stop postgresql

# Restore from backup
pg_restore --clean --if-exists \
  --dbname=scintel_db \
  /var/lib/postgresql/backups/scintel_20260718_020000.backup

# Restart
systemctl start postgresql
```

**Option B: Fix Specific Corrupted Records (If isolated corruption)**
```bash
# ONLY if < 100 records corrupted AND you know the fix

# Example: Fix negative inventory
UPDATE inventory
SET quantity = 0
WHERE quantity < 0
AND warehouse_id = 'wh-jakarta-001';

# Example: Re-calculate inventory from transactions
UPDATE inventory i
SET quantity = (
  SELECT 
    COALESCE(SUM(quantity_delta), 0)
  FROM transactions t
  WHERE t.sku = i.sku 
    AND t.warehouse_id = i.warehouse_id
    AND t.created_at <= NOW()
)
WHERE i.sku = 'SKU-NIKE-001';

# ALWAYS backup before running these!
# ALWAYS run in transaction and verify:
BEGIN;
-- Run fix
-- ROLLBACK if wrong; COMMIT if correct
COMMIT;
```

**Option C: Rebuild Table from Backup Slice + Live Transactions**
```bash
# If only one table corrupted:

# 1. Restore that table from backup
pg_restore --clean --if-exists \
  --dbname=scintel_db \
  --table=inventory \
  /var/lib/postgresql/backups/scintel_20260718_020000.backup

# 2. Replay transactions created after backup
INSERT INTO inventory (sku, warehouse_id, quantity, last_updated)
SELECT 
  sku,
  warehouse_id,
  SUM(quantity_delta),
  MAX(created_at)
FROM transactions
WHERE created_at > '2026-07-18 02:00:00'  -- Backup time
GROUP BY sku, warehouse_id
ON CONFLICT (sku, warehouse_id) DO UPDATE
SET quantity = EXCLUDED.quantity,
    last_updated = EXCLUDED.last_updated;
```

---

**Phase 3: Validation** ✅

```bash
# 1. Run data integrity checks
psql -U postgres -d scintel_db -c "
  -- Check inventory >= 0
  SELECT COUNT(*) as invalid_count FROM inventory WHERE quantity < 0;
  -- Expected: 0
  
  -- Check forecasts valid
  SELECT COUNT(*) FROM forecasts 
  WHERE confidence < 0 OR confidence > 1 OR confidence IS NULL;
  -- Expected: 0
  
  -- Check foreign keys
  SELECT COUNT(*) FROM inventory i
  LEFT JOIN warehouses w ON i.warehouse_id = w.id
  WHERE w.id IS NULL;
  -- Expected: 0 (all inventory linked to valid warehouses)
"

# 2. Reconcile inventory sums
psql -U postgres -d scintel_db -c "
  WITH trans_sum AS (
    SELECT 
      sku,
      warehouse_id,
      SUM(quantity_delta) as calculated_qty
    FROM transactions
    GROUP BY sku, warehouse_id
  )
  SELECT 
    COUNT(CASE WHEN ABS(i.quantity - ts.calculated_qty) <= 1 THEN 1 END)::float 
    / COUNT(*) * 100 as percent_reconciled
  FROM inventory i
  LEFT JOIN trans_sum ts ON i.sku = ts.sku AND i.warehouse_id = ts.warehouse_id;
"
# Expected: > 99% reconciled

# 3. Manual spot checks
# Ask warehouse: "Are these stock levels correct?"
SELECT sku, warehouse_id, quantity, last_updated
FROM inventory
WHERE warehouse_id IN ('wh-jakarta-001', 'wh-surabaya-001')
LIMIT 20;
```

---

**Post-Recovery Actions** 📋

- [ ] **RCA (Root Cause Analysis)**: How did corruption happen?
  - Bug in application code? (fix code + deploy patch)
  - Hardware issue? (run disk check, replace if faulty)
  - Operator error? (audit permissions + add safeguards)
  
- [ ] **Add Safeguards**:
  - Enable table replication to standby (for backup)
  - Add triggers to validate data on insert/update
  - Implement checksums on critical tables
  - Schedule regular REINDEX jobs to detect corruption early

- [ ] **Team Audit**: How close were we to undetected corruption?
  - Did monitoring catch this, or manual discovery?
  - How long before corruption detected? (impact window)

---

**Expected RTO**: **20-30 minutes** (restore from backup + validation)

**Expected RPO**: **5 minutes** (backup frequency)

---

---

#### **SCENARIO #12: BACKUP STORAGE FAILURE**

**Severity**: 🔴 **CRITICAL** (RTO: ≤120 min, RPO: ≤24 hr)

**Symptoms / Detection**:

Alert Triggers:
- 🚨 Alert: Backup job fails (cron error: "Permission denied" OR "Disk full")
- 🚨 Alert: S3/backup storage unreachable (connection timeout, 403 Forbidden)
- 🚨 Alert: Backup integrity check fails (corrupted backup files detected)
- 🚨 Alert: Backup disk is full (`/var/lib/postgresql/backups/` 100%)
- 🚨 Alert: Encryption key for backups unavailable (cannot decrypt restore)
- 🚨 Alert: No valid backup exists for last 48 hours

**Impact Assessment**:

Affected Systems:
- 🔴 **NO BACKUPS AVAILABLE** = Cannot recover from disaster
- 🔴 If primary database crashes → **PERMANENT DATA LOSS** (no recovery path)
- ✅ Current system operations unaffected (backup failure doesn't impact running system)
- ✅ Only matters IF disaster occurs (like having fire insurance, useless until fire)

User Impact:
- **None immediately** (backups are preventive, not operational)
- **Catastrophic if disaster strikes**: All data lost, business cannot recover
- Compliance breach (required backups unavailable)

Business Impact:
- **WORST CASE**: Database crashes + no backup → **Total data loss + business stoppage**
- Cannot meet RTO/RPO targets
- Regulatory fines (GDPR/compliance: required backups unavailable)
- Customer trust destroyed

---

**Immediate Actions (First 15 minutes)** ⏱️

1. **[0-1 min] VERIFY BACKUP FAILURE**:
   ```bash
   # Check last backup timestamp
   ls -lh /var/lib/postgresql/backups/ | tail -5
   
   # Expected: Recent backup (created today or yesterday)
   # Actual: Backup from 48+ hours ago → PROBLEM
   
   # Check backup logs
   tail -50 /var/log/cron | grep backup
   tail -50 /var/log/postgresql/backup.log | grep -i "error\|failed"
   ```

2. **[1-2 min] ASSESS BACKUP STORAGE**:
   ```bash
   # Local backup storage
   df -h /var/lib/postgresql/backups/
   # Expected: > 50% free space
   # Problem: 90%+ full (disk full)
   
   # S3/Remote storage
   aws s3 ls s3://scintel-backups/ --recursive | tail -20
   aws s3 ls s3://scintel-backups/daily/ --human-readable --summarize | tail -5
   ```

3. **[2-5 min] IDENTIFY ROOT CAUSE**:
   ```bash
   # Possible causes:
   # 1. Disk full (local backup directory)
   du -sh /var/lib/postgresql/backups/
   
   # 2. S3 permission denied
   aws s3 ls s3://scintel-backups/ 2>&1 | grep -i "access\|denied\|forbidden"
   
   # 3. Encryption key unavailable
   aws secretsmanager get-secret-value --secret-id scintel-backup-key 2>&1
   
   # 4. Network connectivity to S3
   curl -I https://s3.amazonaws.com
   ```

4. **[5-10 min] IMMEDIATE ACTIONS**:
   - [ ] **If disk full**: Clean up old backups
     ```bash
     # Remove backups older than 30 days (keep monthly)
     find /var/lib/postgresql/backups/ -name "*.sql*" -mtime +30 -delete
     ```
   
   - [ ] **If S3 unreachable**: Check AWS service status, IAM permissions
   
   - [ ] **If encryption key lost**: Retrieve from Secrets Manager or backup copy

5. **[10-15 min] MANUAL BACKUP TRIGGER**:
   ```bash
   # Create backup immediately (don't wait for cron)
   pg_dump --verbose \
     --format=plain \
     --file=/var/lib/postgresql/backups/scintel_$(date +%Y%m%d_%H%M%S).sql \
     scintel_db
   
   # Test restore on staging (verify backup is valid)
   # This ensures we have at least ONE known-good backup
   ```

---

**Recovery Steps** 🛠️

**Phase 1: Fix Immediate Issue**

**If Disk Full**:
```bash
# Identify large files
du -sh /var/lib/postgresql/backups/* | sort -h | tail -20

# Remove old/test backups
rm /var/lib/postgresql/backups/test_*.sql*
rm /var/lib/postgresql/backups/scintel_2026*old*.sql*

# Compress old backups to save space
gzip /var/lib/postgresql/backups/scintel_2026060*.sql  # Compress June backups

# Verify space freed
df -h /var/lib/postgresql/backups/
# Expected: Now has > 50% free space
```

**If S3 Unreachable**:
```bash
# Check S3 bucket permissions
aws s3api head-bucket --bucket scintel-backups

# Verify IAM policy allows uploads
aws iam get-user-policy --user-name scintel-backup-service --policy-name BackupS3Policy

# If permission denied, update policy:
aws iam put-user-policy --user-name scintel-backup-service \
  --policy-name BackupS3Policy \
  --policy-document '{
    "Version": "2012-10-17",
    "Statement": [{
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:GetObject", "s3:ListBucket"],
      "Resource": ["arn:aws:s3:::scintel-backups", "arn:aws:s3:::scintel-backups/*"]
    }]
  }'
```

**If Encryption Key Lost**:
```bash
# Retrieve from AWS Secrets Manager
aws secretsmanager get-secret-value \
  --secret-id scintel-backup-encryption-key \
  --query SecretString \
  --output text > /tmp/backup-key.txt

# If secret not found, generate new key
openssl rand -hex 32 > /tmp/backup-key-new.txt

# Store in Secrets Manager
aws secretsmanager create-secret \
  --name scintel-backup-encryption-key-new \
  --secret-string $(cat /tmp/backup-key-new.txt)

# Update backup config to use new key
```

**Phase 2: Verify Backup Integrity**

```bash
# Test restore (non-prod environment)
psql -U postgres -d scintel_test << EOF
DROP DATABASE IF EXISTS scintel_restore_test;
CREATE DATABASE scintel_restore_test;
EOF

pg_restore --clean --if-exists \
  --dbname=scintel_restore_test \
  /var/lib/postgresql/backups/scintel_latest.backup

# Verify test restore
psql -U postgres -d scintel_restore_test -c "
  SELECT 
    COUNT(*) as table_count,
    SUM(COALESCE(
      (SELECT COUNT(*) FROM inventory),
      (SELECT COUNT(*) FROM transactions),
      0
    )) as row_sample
  FROM information_schema.tables
  WHERE table_schema = 'public';
"

# Expected: Tables exist, rows present = backup OK
```

**Phase 3: Implement Backup Safeguards**

```bash
# 1. Enable automated backup monitoring
cat > /usr/local/bin/backup-health-check.sh << 'EOF'
#!/bin/bash

# Check last backup age
last_backup_time=$(ls -t /var/lib/postgresql/backups/scintel_*.sql 2>/dev/null | head -1 | xargs stat -c %Y)
now=$(date +%s)
backup_age=$(( ($now - $last_backup_time) / 60 ))  # Minutes

if [ $backup_age -gt 1440 ]; then  # 24 hours
  echo "❌ ALERT: Backup older than 24 hours ($backup_age minutes ago)"
  # Send alert to Slack / email
  curl -X POST $SLACK_WEBHOOK -d "{\"text\": \"Backup health check failed: $backup_age min old\"}"
  exit 1
fi

echo "✅ Backup health OK (Age: $backup_age minutes)"
exit 0
EOF

chmod +x /usr/local/bin/backup-health-check.sh

# Add to cron to run every 6 hours
echo "0 */6 * * * /usr/local/bin/backup-health-check.sh" | crontab -

# 2. Set up backup rotation policy
cat > /usr/local/bin/backup-rotation.sh << 'EOF'
#!/bin/bash

BACKUP_DIR=/var/lib/postgresql/backups
DAILY_RETENTION=7  # Keep 7 days of daily backups
WEEKLY_RETENTION=30  # Keep 30 days of weekly backups

# Archive old daily backups to S3 Glacier
find $BACKUP_DIR -name "scintel_*.sql" -mtime +$DAILY_RETENTION | while read file; do
  aws s3 mv "$file" "s3://scintel-backups/archive/$(basename $file)" \
    --storage-class GLACIER
done

# Clean up local backups older than 30 days
find $BACKUP_DIR -name "scintel_*.sql" -mtime +$WEEKLY_RETENTION -delete

echo "✅ Backup rotation complete"
EOF

# 3. Enable S3 versioning (protect against accidental deletion)
aws s3api put-bucket-versioning \
  --bucket scintel-backups \
  --versioning-configuration Status=Enabled
```

---

**Phase 4: Validation** ✅

```bash
# 1. Backup exists and recent
ls -lh /var/lib/postgresql/backups/scintel_*.sql | tail -3
# Expected: Recent backup (today or yesterday)

# 2. Backup can be restored
pg_restore --list /var/lib/postgresql/backups/scintel_latest.backup | head -20
# Expected: No errors, lists table/sequence definitions

# 3. S3 backup accessible
aws s3 ls s3://scintel-backups/daily/
# Expected: Shows uploaded backup files

# 4. Encryption key available
aws secretsmanager get-secret-value --secret-id scintel-backup-key > /dev/null 2>&1
echo "Backup encryption key: $?"
# Expected: Exit code 0 (key accessible)
```

---

**Expected RTO**: **60-120 minutes** (fix storage + restore + validate)

**Expected RPO**: **24 hours** (acceptable for backups; daily frequency)

---

---

## 6. ROLES, RESPONSIBILITIES & ESCALATION

### 6.1 Incident Response Team (Based on PRD Section 9 - Scope & Organization)

| Role | Name/Title | Responsibilities | Contact | On-Call |
|------|-----------|-----------------|---------|---------|
| **Incident Commander** | [TBD - Assigned at deployment] | Coordinate response, decision-making, communication with stakeholders. Authorize emergency actions (e.g., database failover). | [Phone/Slack] | Primary |
| **Database Administrator (DBA)** | [TBD - Assigned at deployment] | Database recovery, backup management, failover coordination, replication monitoring. | [Phone/Slack] | Primary |
| **Infrastructure/DevOps Lead** | [TBD - Assigned at deployment] | Server, network, load balancer failover. Docker/Kubernetes operations. | [Phone/Slack] | Primary |
| **Golang Backend Lead** | [TBD - Assigned at deployment] | Application-specific issues, code review for bugs, API restart decisions. | [Phone/Slack] | Primary |
| **Communications Lead** | [TBD - Assigned at deployment] | Internal & external communications, status page updates, stakeholder notifications. | [Slack/Email] | Secondary |
| **CTO / VP Engineering** | [TBD] | Final escalation, strategic decisions, public apology if needed. | [Phone] | Escalation |

---

### 6.2 Escalation Path (Decision Tree)

```
Initial Alert (0-5 min)
    ↓
On-Call Engineer investigates
    ↓
    ├─ Issue resolved < 5 min? → Document incident, no escalation needed
    │
    └─ Issue NOT resolved after 5 min? → ESCALATE
         ↓
    Incident Commander engaged
    + Tech Lead (DBA / Infra) conference call
    + Communications Lead alerted
         ↓
         ├─ RTO target will be breached? → ESCALATE
         │     ↓
         │   Engineering Manager + CTO alerted
         │   + Prepare customer communication
         │     ↓
         │     └─ RTO exceeded by 30 min? → ESCALATE TO CEO/BOARD
         │          (Financial impact requires C-level awareness)
         │
         └─ Data loss confirmed? → ESCALATE TO CTO IMMEDIATELY
              (Compliance/legal implications)
```

**Escalation Triggers**:
- 🚨 **RTO exceeded by 30 minutes** → Notify Engineering Manager + CTO
- 🚨 **Data loss confirmed** → Notify CTO + Legal immediately
- 🚨 **Multiple systems affected** (e.g., DB + API + Cache) → Notify CTO
- 🚨 **Customer-facing impact** lasting >15 minutes → Notify customer success team
- 🚨 **Potential compliance/data breach** → Notify CISO + Legal

---

### 6.3 Communication Plan (During Incident)

**Internal Communications** (Real-time):

| Channel | Frequency | Audience | Message Template |
|---------|-----------|----------|------------------|
| **#incidents Slack channel** | Every 5-10 min (during incident) | Engineering team | "[STATUS] [COMPONENT] [ACTION] [ETA]\n🔴 CRITICAL: Database failover in progress. ETA 10 min. Will update in 5 min." |
| **War Room Conference Call** | Continuous (if >30 min incident) | Incident Commander + Tech leads | Discussion, real-time troubleshooting, decision-making |
| **Incident Tracking System** | Every 15-30 min | Full team | Detailed timeline, action items, root cause hypothesis |

**External Communications** (Customer-facing):

| Scenario | Timing | Channel | Message |
|----------|--------|---------|---------|
| **Service Degradation** (slow, not down) | Immediate | Status page + email | "We're experiencing degraded performance. Our team is investigating. Sorry for inconvenience." |
| **Service Outage** (complete unavailability) | Immediate | Status page + SMS + email | "Supply Chain Intelligence Platform is currently unavailable. We're working to restore service. Updates every 15 min." |
| **Partial Outage** (some users affected) | Immediate | Status page | "Warehouse inventory queries experiencing delays. Other features unaffected. Being resolved now." |
| **Recovery Update** | Every 15-30 min | Status page + email | "Service recovering. Database failover complete. Expected full recovery in 5 minutes." |
| **All-Clear** | Within 30 min of recovery | Status page + email | "✅ Service fully restored. Details: [root cause summary] + apology + link to post-mortem." |

**Communication Template**:
```
Incident: [NAME]
Status: [INVESTIGATING / IN PROGRESS / RESOLVED]
Severity: [CRITICAL / HIGH / MEDIUM]
Affected Systems: [List which services/features are affected]
Impact: [# users/warehouses affected; estimated revenue loss if applicable]
Last Update: [ISO 8601 timestamp]
Next Update: [In X minutes]
ETA Recovery: [If known]

For updates, visit: https://status.scintel.io/
```

---

## 7. TESTING, VALIDATION & DRILLS

### 7.1 Testing Strategy (Based on PRD Section 8 - Testing & TDD Section 10)

| Test Type | Frequency | Scope | Success Criteria | Owner |
|-----------|-----------|-------|-----------------|-------|
| **Backup Restoration Test** | **Monthly** | Full database restore to staging environment | Restore time < 30 min; data integrity check passes; all tables accessible | DBA |
| **Failover Test** | **Quarterly** (or after major infrastructure changes) | Primary DB → Standby promotion; verify no data loss | Failover time < 2 min; zero data loss; read/write operations immediately resume | DBA + Infra |
| **API Failover Test** | **Monthly** | Graceful shutdown of one API server; verify traffic reroutes | Failover < 30 sec; zero request loss; error rate < 0.1% during failover | Infra Lead |
| **Webhook Delivery Test** | **Monthly** | Simulate webhook delivery queue backup; test retry logic | All events delivered within 5 retries; manual retry endpoint works | Backend Lead |
| **Disaster Recovery Drill** | **Monthly** (tabletop) + **Quarterly** (full simulation) | Simulate specific scenario (choose one per month) | Playbook executed correctly; actual RTO/RPO measured; team confirms procedures | Incident Commander |
| **Application Smoke Tests** | **After every deployment** | Critical user journeys (login, inventory query, forecast fetch) | All endpoints return 2xx; response time < 1s | QA / CI/CD |

---

### 7.2 Disaster Recovery Drills (Monthly / Quarterly)

**Monthly Drill Format** (30-45 minutes, no-impact):

Scenario: Choose one month, one scenario
- Month 1 (July): Scenario #1 - Database failure
- Month 2 (August): Scenario #2 - Redis cache failure
- Month 3 (September): Scenario #3 - API server down
- Etc.

Participants: On-call engineer + tech leads + incident commander (optional)

Execution (no-impact):
1. **Announce drill** (5 min before): "🎯 DRILL: Simulating Database Failover. This is a TEST. No customer impact."
2. **Read scenario** (5 min): Hand out playbook; let team review
3. **Execute playbook** (30 min): Team performs recovery steps against staging/non-prod database
4. **Measure metrics** (5 min): Calculate actual RTO/RPO
5. **Debrief** (10 min): What went well? What was confusing? Update playbook.

**Quarterly Drill Format** (2-4 hours, production-like simulation):

Scenario: Full system recovery test
- Production-like data volume (full database dump in staging)
- All team members participate
- Realistic constraints (e.g., "Primary DB unavailable for 30 minutes")
- Measure: Actual recovery time, team coordination, communication effectiveness

Execution:
1. **Pre-drill setup** (30 min): Backup production state; prepare staging infrastructure
2. **Inject failure** (Simulate disaster - DBA simulates DB crash)
3. **Team responds** (1-2 hours): Execute recovery playbook; provide real-time updates
4. **Validate recovery** (30 min): Verify data integrity, run smoke tests
5. **Post-drill debrief** (30 min): Document lessons learned; identify gaps

---

### 7.3 Drill Execution Checklist

**Pre-Drill** (1 day before):

- [ ] Notify all participants (Slack + email)
- [ ] Schedule meeting/war room (Zoom/conference call)
- [ ] Prepare staging environment (backup production state)
- [ ] Backup current state (if applicable)
- [ ] Set drill duration (e.g., 45 minutes for monthly; 3 hours for quarterly)
- [ ] Define success criteria (specific RTO/RPO targets)
- [ ] Assign roles: Incident Commander, DBA, Infra, QA, Communications
- [ ] Prepare props/simulation (e.g., "kill primary DB process" script)

**During Drill**:

- [ ] Start timer (clock starts when failure injected)
- [ ] Execute playbook steps exactly as written
- [ ] Log all actions & timings in shared document
- [ ] Record decisions made (e.g., "Chose failover vs restore")
- [ ] Note issues/blockers encountered (e.g., "Playbook missing step X")
- [ ] Measure actual RTO/RPO (time from detection to recovery)
- [ ] Capture screenshots/logs for post-mortem review

**After Drill** (same day):

- [ ] Stop timer; calculate actual RTO/RPO
- [ ] Verify system recovery (smoke tests pass)
- [ ] Collect feedback from participants (Slack poll or quick survey)
- [ ] Create after-action report (what went well; what to improve)
- [ ] Update playbook based on findings (add missing steps, clarify procedures)
- [ ] Assign owners for follow-up items (e.g., "DBA to optimize backup restore time")

---

### 7.4 Metrics & KPIs (Tracked from each drill/incident)

Track these metrics over time to measure DRP maturity:

| Metric | Target | Current | Notes |
|--------|--------|---------|-------|
| **Actual RTO (Database Failure)** | ≤15 min | TBD | Measured in drill |
| **Actual RPO (Database Failure)** | ≤5 min | TBD | Data loss acceptable range |
| **Playbook Accuracy** | ≥95% | TBD | % steps that matched actual execution |
| **Team Response Time** | ≤5 min | TBD | Time from alert to incident commander engaged |
| **Communication Adherence** | 100% | TBD | Did team follow communication plan? |
| **Recovery Validation Time** | ≤10 min | TBD | Time to run smoke tests & declare "all clear" |
| **Knowledge Gaps Identified** | ≤2 per drill | TBD | Missing procedures, unclear steps, team confusion |
| **Backup Restoration Success Rate** | 100% | TBD | % backup restoration tests that succeeded |

---

## 8. MAINTENANCE & CONTINUOUS IMPROVEMENT

### 8.1 Regular Maintenance Tasks (Ongoing, every week/month)

| Task | Frequency | Owner | Purpose |
|------|-----------|-------|---------|
| **Backup Verification** | Weekly | DBA | Verify backups completing, no errors; random restore test |
| **Review Playbooks** | Quarterly | Incident Commander | Update for new systems/features; mark outdated sections |
| **Update Contact List** | Monthly | Communications Lead | Verify phone/email/Slack current; test page/call tree |
| **Review RTO/RPO Targets** | Semi-annually | Incident Commander + CTO | Align with business priorities; adjust if infrastructure changes |
| **Audit Backup Storage** | Monthly | DBA | Verify backup size reasonable; encryption working; retention policy honored |
| **Test Failover Procedure** | Quarterly | DBA + Infra | Practice promoting standby; measure actual failover time |
| **Certificate/Secret Rotation** | As needed (monthly check) | Security Lead | JWT secret, database password, TLS certificates |
| **Clean Up Logs / Archives** | Monthly | Infrastructure | Prevent disk fill-up; compress old logs |

---

### 8.2 Post-Incident Review Process (AFTER REAL INCIDENT)

Incident management process to prevent recurrence:

**Within 24 hours**:
```
1. Document what happened:
   - Timeline: When did incident start? When detected? When resolved?
   - Actions taken: What was the recovery procedure? Did playbook match?
   - Data loss: Any transactions lost? Audit trail gaps?
   - Root cause hypothesis: Initial assessment (detailed RCA later)

2. Identify immediate impacts:
   - How many users affected? For how long?
   - Estimated revenue loss?
   - Data integrity issues?
   - Compliance/audit implications?

3. Measure RTO/RPO:
   - Actual RTO: Time from alert to full recovery
   - Actual RPO: Amount of data lost (if any)
   - Did we meet targets?

4. Team feedback:
   - Slack survey: What went well? What was confusing? (3-5 quick questions)
```

**Within 1 week**:
```
1. Complete formal incident report:
   - Title: [Incident Name]
   - Severity: [CRITICAL/HIGH/MEDIUM]
   - Timeline (with actual timestamps)
   - Root cause (detailed technical analysis)
   - Contributing factors
   - Business impact (calculated loss if applicable)
   
2. Identify preventive measures:
   - What could have prevented this? (e.g., monitoring alert, backup test)
   - What could have made recovery faster? (e.g., better runbook, automation)
   - Long-term fixes vs. quick wins
   - Owner + target date for each fix
   
3. Share lessons learned:
   - Post-incident review meeting (1 hour with team)
   - Demo lesson learned to broader engineering
   - Update knowledge base/wiki

4. Update playbooks:
   - Add missing steps
   - Clarify ambiguous sections
   - Incorporate team feedback
   - Version increment (e.g., v1.0 → v1.1)
```

**Within 30 days**:
```
1. Implement preventive measures:
   - Prioritize quick wins (< 1 day fix)
   - Schedule larger projects (2+ weeks)
   - Track progress in incident tracking system

2. Run drill based on learnings:
   - Test new playbook procedures
   - Verify preventive measures effective
   - Team practice on updated procedures

3. Close incident ticket:
   - Mark as "Resolved" once all preventive actions complete
   - Link to post-incident review and action items
```

---

### 8.3 Version Control & Change Log

Track DRP document updates to ensure team always using latest version:

| Version | Date | Changes | Author | Status |
|---------|------|---------|--------|--------|
| 1.0 | 2026-07-18 | Initial DRP creation; 6 critical scenarios; backup strategy; playbooks | M. Arif Aulia (Backend Engineer) | **DRAFT** |
| 1.1 | [Next review] | [Updates based on Phase 1 completion] | [Author] | Pending |
| 2.0 | [Phase 1.5] | [Add Kafka failover; Multi-region DRP; Advanced monitoring] | [Author] | Roadmap |

**Last Updated**: 2026-07-18  
**Next Review**: 2026-10-18 (After Phase 1 completion and first real incident / monthly drill)  
**Status**: **DRAFT → Awaiting Approval**

---

## 9. APPENDIX

### 9.1 Glossary (Disaster Recovery Terms)

| Term | Definition |
|------|-----------|
| **RTO (Recovery Time Objective)** | Maksimum waktu sistem diizinkan offline setelah disaster. Target: ≤15 min untuk critical systems. |
| **RPO (Recovery Point Objective)** | Maksimum data yang diizinkan hilang. Target: ≤5 min untuk critical systems (hourly transactions). |
| **Backup** | Copy database/config files untuk recovery purposes. Disimpan di lokasi terpisah dari primary system. |
| **Failover** | Automatic atau manual switching ke backup systems saat primary fails. Ideally < 1 minute. |
| **Disaster** | Event yang cause system/data unavailable: hardware failure, cyberattack, human error, natural disaster, bugs. |
| **Playbook** | Step-by-step procedure untuk merespons specific disaster scenario. Tested regularly via drills. |
| **War Room** | Real-time incident response center (Zoom call atau physical room) dengan all stakeholders. |
| **Post-Mortem** | Meeting setelah incident untuk analyze root cause dan identify preventive measures. |
| **Replication** | Real-time copying database changes dari primary ke standby/replica. Enables fast failover. |
| **Graceful Degradation** | Service continues operating (slowly) if non-critical components fail. Better than complete outage. |

---

### 9.2 Related Documents

| Document | Purpose | Status |
|----------|---------|--------|
| **PRD.md v1.0.0** | Product requirements; defines systems, features, scope | ✅ APPROVED |
| **LOGIC_FLOW.md v1.0.0** | Detailed system architecture & workflows; maps data flow | ✅ APPROVED |
| **TDD.md v1.0.0** | Technical design; API contracts, database schema, business logic | ✅ APPROVED |
| **HALAMAN.md v2.0** | UI/page specifications; user workflows | ✅ APPROVED |
| **Infrastructure Documentation** | Network diagram, server specs, deployment guide | [Link TBD] |
| **API Contract / Swagger/OpenAPI** | Detailed API endpoint specifications | [Link TBD] |
| **Security Policy** | Authentication, authorization, data protection, audit logging | [Link TBD] |
| **Backup Strategy Document** | Detailed backup procedures, retention, verification | [Link TBD] |
| **Incident Response Policy** | How to escalate, communicate, track incidents | [Link TBD] |
| **Monitoring & Alerting Setup** | Prometheus, Grafana, alert thresholds, notif channels | [Link TBD] |

---

### 9.3 Tools & Resources (To Be Implemented)

| Tool/Resource | Purpose | Status | Access |
|---------------|---------|--------|--------|
| **PostgreSQL** | Primary database; backup/restore via pg_dump, WAL archiving | TBD | [Host/Connection string] |
| **Redis** | Cache layer; RDB snapshots for backup | TBD | [Host:Port/Connection string] |
| **AWS S3 / MinIO** | Backup storage (encrypted, geo-redundant) | TBD | [S3 bucket ARN] |
| **HashiCorp Vault / AWS Secrets Manager** | JWT secret, database credentials, API keys | TBD | [URL/Access method] |
| **Prometheus + Grafana** | Monitoring & alerting; real-time dashboards | TBD | [URL] |
| **Slack / PagerDuty** | Incident notifications, on-call escalation | TBD | [Channel/Integration] |
| **GitHub Actions** | CI/CD pipeline for automated deployments | TBD | [Repository] |
| **Docker & Docker Compose** | Application containerization, local dev environment | TBD | [Hub/Registry] |
| **Kubernetes** | Production orchestration (Phase 1.5+) | TBD | [Cluster info] |
| **Incident Tracking System** | Jira / Linear / GitHub Issues for incident tracking | TBD | [URL/Project] |

---

### 9.4 Recovery Commands Reference (Quick Lookup)

**Database Backup & Restore** (PostgreSQL):

```bash
# Full Backup (Create)
pg_dump --verbose --format=plain scintel_db > /var/lib/postgresql/backups/scintel_$(date +%Y%m%d).sql

# Full Restore
psql -U postgres scintel_db < /var/lib/postgresql/backups/scintel_20260718.sql

# Backup with compression
pg_dump --compress=9 scintel_db > scintel_$(date +%Y%m%d).sql.gz

# Restore from compressed
gunzip -c scintel_20260718.sql.gz | psql -U postgres scintel_db

# Point-in-time recovery (restore to specific timestamp)
pg_ctl restore -X from-wal -D /var/lib/postgresql/main \
  --recovery-target-timeline latest \
  --recovery-target-time '2026-07-18 10:00:00'

# Check database size
psql -U postgres -c "SELECT pg_size_pretty(pg_database_size('scintel_db'));"

# Verify backup integrity
pg_dump --format=tar scintel_db | tar -t > /dev/null
# Expected exit code: 0 (backup readable)
```

**API Server Management** (Docker):

```bash
# Check container status
docker ps | grep scintel-api
docker logs scintel-api-1

# Restart single API server
docker restart scintel-api-1

# Restart all API servers
docker restart scintel-api-1 scintel-api-2 scintel-api-3

# Check health
curl http://scintel-api-1:8080/api/v1/health/ready

# Graceful shutdown (allows in-flight requests to complete)
docker stop -t 30 scintel-api-1  # 30 second shutdown grace period
```

**Cache Layer** (Redis):

```bash
# Connect to Redis CLI
redis-cli -h localhost -p 6379

# Inside redis-cli:
PING  # Test connection
DBSIZE  # Show number of keys
FLUSHALL  # Clear all cache (use with care!)
KEYS "inventory:*"  # List cache keys matching pattern
INFO memory  # Check memory usage
SAVE  # Force RDB backup
SHUTDOWN SAVE  # Graceful shutdown

# Clear specific cache keys
redis-cli DEL "inventory:*"
redis-cli DEL "forecast:*"
```

**Load Balancer / Network Troubleshooting**:

```bash
# Check if port is listening
netstat -tlnp | grep 8080
ss -tlnp | grep 8080

# Test API endpoint via load balancer
curl -v http://{load_balancer}:8080/api/v1/health/live

# Check DNS resolution
nslookup db.scintel.io
dig db.scintel.io

# Ping database server
ping {db_host}

# Test database port connectivity
telnet {db_host} 5432
```

**System Resources & Debugging**:

```bash
# Memory usage
free -h
ps aux | grep scintel | sort -k4 -nr | head -5

# Disk usage
df -h
du -sh /var/lib/postgresql/

# CPU utilization
top -b -n 1 | head -20

# Network connections
ss -an | grep ESTABLISHED
ss -an | grep 5432  # PostgreSQL connections

# View process limits
ulimit -a
# Check if file descriptor limit exceeded:
cat /proc/{pid}/limits
```

---

### 9.5 Contact & On-Call Information (To Be Filled)

**On-Call Rotation** (24/7 coverage):

| Role | Week 1 | Week 2 | Week 3 | Week 4 |
|------|--------|--------|--------|--------|
| **Incident Commander** | [Name] | [Name] | [Name] | [Name] |
| **Database Administrator** | [Name] | [Name] | [Name] | [Name] |
| **Infrastructure Lead** | [Name] | [Name] | [Name] | [Name] |
| **Backend Lead** | [Name] | [Name] | [Name] | [Name] |

**Escalation Numbers**:

```
Level 1 (On-Call Engineer)
  → Slack: @on-call-engineer
  → Phone: [Number]

Level 2 (Tech Lead)
  → Slack: @tech-lead
  → Phone: [Number]

Level 3 (Engineering Manager)
  → Slack: @engineering-manager
  → Phone: [Number]

Level 4 (CTO / VP)
  → Slack: @cto
  → Phone: [Number]
  → Personal cell: [Number]
```

**Status Page & Communication**:

- Status Page: https://status.scintel.io/
- Incident Tracking: [Jira/Linear project URL]
- War Room: [Zoom meeting link - always-on]
- Runbooks: [Confluence/Wiki URL]
- Monitoring Dashboard: [Grafana dashboard URL]

---

## Document Sign-Off

| Role | Name | Date | Status |
|------|------|------|--------|
| **DRP Owner** | [To Be Assigned] | [Sign-off date] | Pending |
| **Engineering Lead** | [To Be Assigned] | [Sign-off date] | Pending |
| **Infrastructure Lead** | [To Be Assigned] | [Sign-off date] | Pending |
| **CTO / Management** | [To Be Assigned] | [Sign-off date] | Pending |

---

**Document Status**: 📝 **DRAFT** (Ready for review & team feedback before approval)

**Last Updated**: 2026-07-18  
**Next Review**: 2026-10-18 (After Phase 1 completion & first production drill)  
**Version**: 1.0.0 - INITIAL

---

## APPENDIX: Additional Scenarios (Phase 1.5+)

The following scenarios are out-of-scope for v1.0 MVP but should be addressed in future DRP updates:

1. **Multi-Region Failover** (Phase 2): Cross-region replication, DNS failover to secondary region
2. **Kafka Integration Failure** (Phase 1.5): Event streaming pipeline failure, message queue backup
3. **Compliance Audit Trail Loss** (Phase 2): Audit log corruption, forensics recovery
4. **Certificate Expiry / TLS Failure** (Phase 1.1): Expired SSL certificates, key rotation
5. **Third-Party API Integration Failures** (Phase 2): Payment gateway down, SMS service down, external service integration

These will be detailed in DRP v2.0 as the platform evolves.

---

**END OF DOCUMENT**

---

**Approved Version**: 1.0.0 (DRAFT)  
**Document Classification**: Internal - Engineering Only  
**Access Control**: GitHub private repo + Confluence wiki (restricted to engineering team)
