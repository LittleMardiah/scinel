package handler

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/middleware"
    "github.com/LittleMardiah/scinel/internal/model"
    "github.com/LittleMardiah/scinel/internal/repository"
    "github.com/LittleMardiah/scinel/internal/util"
)

// AnomalyHandler handles anomaly-related endpoints
type AnomalyHandler struct {
    repo repository.AnomalyRepository
}

// NewAnomalyHandler creates a new anomaly handler
func NewAnomalyHandler(repo repository.AnomalyRepository) *AnomalyHandler {
    return &AnomalyHandler{repo: repo}
}

// ==================== REQUEST/RESPONSE TYPES ====================

type anomalyListResponse struct {
    ID          string    `json:"id"`
    WarehouseID string    `json:"warehouse_id"`
    SKUID       string    `json:"sku_id"`
    Type        string    `json:"type"`
    Severity    string    `json:"severity"`
    Description string    `json:"description"`
    ZScore      float64   `json:"z_score"`
    Status      string    `json:"status"`
    DetectedAt  time.Time `json:"detected_at"`
}

type anomalyDetailResponse struct {
    ID                string     `json:"id"`
    WarehouseID       string     `json:"warehouse_id"`
    SKUID             string     `json:"sku_id"`
    SKUCode           string     `json:"sku_code"`
    Type              string     `json:"type"`
    Severity          string     `json:"severity"`
    Description       string     `json:"description"`
    HistoricalAverage float64    `json:"historical_average"`
    DetectedValue     float64    `json:"detected_value"`
    ZScore            float64    `json:"z_score"`
    Status            string     `json:"status"`
    ResolutionNotes   *string    `json:"resolution_notes,omitempty"`
    ResolvedBy        *string    `json:"resolved_by,omitempty"`
    ResolvedAt        *time.Time `json:"resolved_at,omitempty"`
    DetectedAt        time.Time  `json:"detected_at"`
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`
}

type updateStatusRequest struct {
    Status         string  `json:"status"`
    ResolutionNotes *string `json:"resolution_notes,omitempty"`
}

type timelineEvent struct {
    Timestamp time.Time `json:"timestamp"`
    Event     string    `json:"event"` // "detected", "acknowledged", "resolved"
    UserID    *string   `json:"user_id,omitempty"`
    Notes     *string   `json:"notes,omitempty"`
}

type timelineResponse struct {
    AnomalyID string          `json:"anomaly_id"`
    Timeline  []timelineEvent `json:"timeline"`
}

// ==================== ENDPOINT 1: GET /api/v1/anomalies ====================

// ListAnomalies retrieves all anomalies with filters
// GET /api/v1/anomalies
func (h *AnomalyHandler) ListAnomalies(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    // Parse filters
    statusFilter := r.URL.Query().Get("status")
    severityFilter := r.URL.Query().Get("severity")
    warehouseFilter := r.URL.Query().Get("warehouse_id")
    limit, offset := util.ParsePagination(r)

    if limit > 100 {
        limit = 100
    }

    // Get all anomalies (TODO: implement filter in repository for production)
    anomalies, total, err := h.repo.GetAll(r.Context(), limit, offset)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve anomalies", err.Error(), requestID)
        return
    }

    // Filter in-memory (MVP approach - will be optimized later)
    filtered := []model.Anomaly{}
    for _, a := range anomalies {
        if statusFilter != "" && string(a.Status) != statusFilter {
            continue
        }
        if severityFilter != "" && string(a.Severity) != severityFilter {
            continue
        }
        if warehouseFilter != "" && a.WarehouseID.String() != warehouseFilter {
            continue
        }
        filtered = append(filtered, a)
    }
    total = int64(len(filtered))

    // Build response
    items := make([]anomalyListResponse, len(filtered))
    for i, a := range filtered {
        items[i] = anomalyListResponse{
            ID:          a.ID.String(),
            WarehouseID: a.WarehouseID.String(),
            SKUID:       a.SKUID.String(),
            Type:        a.Type,
            Severity:    string(a.Severity),
            Description: a.Description,
            ZScore:      a.ZScore,
            Status:      string(a.Status),
            DetectedAt:  a.DetectedAt,
        }
    }

    data := map[string]interface{}{
        "anomalies": items,
        "total":     total,
        "limit":     limit,
        "offset":    offset,
    }

    util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== ENDPOINT 2: GET /api/v1/anomalies/{id} ====================

// GetAnomaly retrieves detailed information for a specific anomaly
// GET /api/v1/anomalies/{id}
func (h *AnomalyHandler) GetAnomaly(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    // Extract ID from path
    path := strings.TrimPrefix(r.URL.Path, "/api/v1/anomalies/")
    if path == "" || strings.Contains(path, "/") {
        util.WriteError(w, http.StatusBadRequest, "Invalid anomaly ID", "", requestID)
        return
    }

    anomalyID, err := uuid.Parse(path)
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid anomaly ID format", err.Error(), requestID)
        return
    }

    anomaly, err := h.repo.GetByID(r.Context(), anomalyID)
    if err != nil {
        util.WriteError(w, http.StatusNotFound, "Anomaly not found", err.Error(), requestID)
        return
    }

    resp := anomalyDetailResponse{
        ID:                anomaly.ID.String(),
        WarehouseID:       anomaly.WarehouseID.String(),
        SKUID:             anomaly.SKUID.String(),
        SKUCode:           anomaly.SKUCode,
        Type:              anomaly.Type,
        Severity:          string(anomaly.Severity),
        Description:       anomaly.Description,
        HistoricalAverage: anomaly.HistoricalAverage,
        DetectedValue:     anomaly.DetectedValue,
        ZScore:            anomaly.ZScore,
        Status:            string(anomaly.Status),
        ResolutionNotes:   anomaly.ResolutionNotes,
        DetectedAt:        anomaly.DetectedAt,
        CreatedAt:         anomaly.CreatedAt,
        UpdatedAt:         anomaly.UpdatedAt,
    }
    if anomaly.ResolvedBy != nil {
        resolvedBy := anomaly.ResolvedBy.String()
        resp.ResolvedBy = &resolvedBy
    }
    resp.ResolvedAt = anomaly.ResolvedAt

    util.WriteResponse(w, http.StatusOK, "OK", resp, requestID)
}

// ==================== ENDPOINT 3: PUT /api/v1/anomalies/{id} ====================

// UpdateAnomalyStatus updates anomaly status (acknowledge or resolve)
// PUT /api/v1/anomalies/{id}
func (h *AnomalyHandler) UpdateAnomalyStatus(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    // Extract ID from path
    path := strings.TrimPrefix(r.URL.Path, "/api/v1/anomalies/")
    if path == "" || strings.Contains(path, "/") {
        util.WriteError(w, http.StatusBadRequest, "Invalid anomaly ID", "", requestID)
        return
    }

    anomalyID, err := uuid.Parse(path)
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid anomaly ID format", err.Error(), requestID)
        return
    }

    // Get user info from context
    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    // Parse request body
    var req updateStatusRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error(), requestID)
        return
    }

    // Validate status
    if req.Status != "acknowledged" && req.Status != "resolved" {
        util.WriteError(w, http.StatusBadRequest, "Invalid status", "Status must be 'acknowledged' or 'resolved'", requestID)
        return
    }

    // Get current anomaly to check current status
    current, err := h.repo.GetByID(r.Context(), anomalyID)
    if err != nil {
        util.WriteError(w, http.StatusNotFound, "Anomaly not found", err.Error(), requestID)
        return
    }

    // Prevent invalid transitions
    if current.Status == model.StatusResolved {
        util.WriteError(w, http.StatusConflict, "Anomaly already resolved", "Cannot update resolved anomaly", requestID)
        return
    }
    if current.Status == model.StatusOpen && req.Status == "resolved" {
        util.WriteError(w, http.StatusConflict, "Invalid status transition", "Anomaly must be acknowledged before resolving", requestID)
        return
    }

    resolvedBy, _ := uuid.Parse(userInfo.UserID)
    notes := ""
    if req.ResolutionNotes != nil {
        notes = *req.ResolutionNotes
    }

    err = h.repo.UpdateStatus(r.Context(), anomalyID, model.AnomalyStatus(req.Status), notes, resolvedBy)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to update anomaly status", err.Error(), requestID)
        return
    }

    // Fetch updated anomaly
    updated, err := h.repo.GetByID(r.Context(), anomalyID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve updated anomaly", err.Error(), requestID)
        return
    }

    data := map[string]interface{}{
        "id":         updated.ID.String(),
        "status":     string(updated.Status),
        "updated_at": updated.UpdatedAt,
    }

    util.WriteResponse(w, http.StatusOK, "Anomaly updated", data, requestID)
}

// ==================== ENDPOINT 4: GET /api/v1/anomalies/{id}/history ====================

// GetAnomalyHistory retrieves anomaly resolution timeline
// GET /api/v1/anomalies/{id}/history
func (h *AnomalyHandler) GetAnomalyHistory(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    // Extract ID from path
    path := strings.TrimPrefix(r.URL.Path, "/api/v1/anomalies/")
    if path == "" {
        util.WriteError(w, http.StatusBadRequest, "Invalid anomaly ID", "", requestID)
        return
    }
    parts := strings.Split(path, "/")
    if len(parts) != 2 || parts[1] != "history" {
        util.WriteError(w, http.StatusBadRequest, "Invalid path", "Expected /api/v1/anomalies/{id}/history", requestID)
        return
    }

    anomalyID, err := uuid.Parse(parts[0])
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid anomaly ID format", err.Error(), requestID)
        return
    }

    // Get anomaly
    anomaly, err := h.repo.GetByID(r.Context(), anomalyID)
    if err != nil {
        util.WriteError(w, http.StatusNotFound, "Anomaly not found", err.Error(), requestID)
        return
    }

    // Build timeline from anomaly data
    timeline := []timelineEvent{}

    // Event 1: Detected
    timeline = append(timeline, timelineEvent{
        Timestamp: anomaly.DetectedAt,
        Event:     "detected",
        Notes:     &anomaly.Description,
    })

    // Event 2: Acknowledged (if status is acknowledged or resolved)
    if anomaly.Status == model.StatusAcknowledged || anomaly.Status == model.StatusResolved {
        acknowledgedAt := anomaly.UpdatedAt
        // If resolved_at exists, use it, but for acknowledged we use updated_at
        timeline = append(timeline, timelineEvent{
            Timestamp: acknowledgedAt,
            Event:     "acknowledged",
            Notes:     anomaly.ResolutionNotes,
        })
    }

    // Event 3: Resolved
    if anomaly.Status == model.StatusResolved && anomaly.ResolvedAt != nil {
        timeline = append(timeline, timelineEvent{
            Timestamp: *anomaly.ResolvedAt,
            Event:     "resolved",
            Notes:     anomaly.ResolutionNotes,
        })
    }

    resp := timelineResponse{
        AnomalyID: anomaly.ID.String(),
        Timeline:  timeline,
    }

    util.WriteResponse(w, http.StatusOK, "OK", resp, requestID)
}