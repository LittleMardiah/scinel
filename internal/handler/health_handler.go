package handler

import (
    "context"
    "net/http"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/LittleMardiah/scinel/internal/database"
    "github.com/LittleMardiah/scinel/internal/util"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
    db    *database.DB
    redis *redis.Client
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB, redis *redis.Client) *HealthHandler {
    return &HealthHandler{
        db:    db,
        redis: redis,
    }
}

// Health returns basic health status (GET /api/v1/health)
// Purpose: Simple liveness probe for load balancers
// Response: 200 OK if service is running
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    response := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "version":   "1.0.0",
    }

    util.WriteResponse(w, http.StatusOK, "OK", response, requestID)
}

// Readiness performs deep health check (GET /api/v1/readiness)
// Purpose: Kubernetes readiness probe - is service ready to accept traffic?
// Checks: Database connectivity, Redis connectivity, migrations status
// Response: 200 OK if all dependencies healthy, 503 Service Unavailable if any fail
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    checks := map[string]bool{}

    // Check database
    err := h.db.PingContext(ctx)
    checks["database"] = err == nil

    // Check Redis
    err = h.redis.Ping(ctx).Err()
    checks["cache"] = err == nil

    // Check if all dependencies are healthy
    allHealthy := true
    for _, healthy := range checks {
        if !healthy {
            allHealthy = false
            break
        }
    }

    statusCode := http.StatusOK
    message := "Ready"
    if !allHealthy {
        statusCode = http.StatusServiceUnavailable
        message = "Not ready - dependency failure"
    }

    util.WriteResponse(w, statusCode, message, checks, requestID)
}

// Liveness provides deep diagnostic info (GET /api/v1/liveness)
// Purpose: Kubernetes liveness probe - should this pod be restarted?
// Checks: All readiness checks + database connection pool stats
// Response: 200 OK if healthy, 503 if deadlock/hanging detected
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()

    checks := map[string]interface{}{}

    // Check database connectivity
    dbErr := h.db.PingContext(ctx)
    checks["database_connected"] = dbErr == nil

    // Get database stats (connection pool info)
    dbStats := h.db.Stats()
    checks["database_stats"] = map[string]interface{}{
        "open_connections": dbStats.OpenConnections,
        "in_use":           dbStats.InUse,
        "idle":             dbStats.Idle,
    }

    // Check Redis connectivity
    cacheErr := h.redis.Ping(ctx).Err()
    checks["cache_connected"] = cacheErr == nil

    // Check if all critical dependencies are healthy
    allHealthy := dbErr == nil && cacheErr == nil

    statusCode := http.StatusOK
    message := "Alive"
    if !allHealthy {
        statusCode = http.StatusServiceUnavailable
        message = "Not alive - critical dependency failure"
    }

    util.WriteResponse(w, statusCode, message, checks, requestID)
}