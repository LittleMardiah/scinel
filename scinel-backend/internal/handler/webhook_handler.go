package handler

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "net/http"
    "net/url"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/middleware"
    "github.com/LittleMardiah/scinel/internal/model"
    "github.com/LittleMardiah/scinel/internal/repository"
    "github.com/LittleMardiah/scinel/internal/util"
)

// WebhookHandler handles webhook-related endpoints
type WebhookHandler struct {
    repo repository.WebhookRepository
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(repo repository.WebhookRepository) *WebhookHandler {
    return &WebhookHandler{repo: repo}
}

// ==================== REQUEST/RESPONSE TYPES ====================

type createSubscriptionRequest struct {
    EventType     string `json:"event_type"`
    WebhookURL    string `json:"webhook_url"`
    MaxRetries    int    `json:"max_retries,omitempty"`
    RetryInterval int    `json:"retry_interval,omitempty"`
}

type subscriptionResponse struct {
    ID         string    `json:"id"`
    EventType  string    `json:"event_type"`
    WebhookURL string    `json:"webhook_url"`
    IsActive   bool      `json:"is_active"`
    SecretKey  string    `json:"secret_key,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type deliveryResponse struct {
    ID             string    `json:"id"`
    SubscriptionID string    `json:"subscription_id"`
    EventID        string    `json:"event_id"`
    Status         string    `json:"status"`
    HTTPStatusCode *int      `json:"http_status_code,omitempty"`
    ErrorMessage   *string   `json:"error_message,omitempty"`
    AttemptCount   int       `json:"attempt_count"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

// ==================== ENDPOINT 1: POST /api/v1/webhooks/subscribe ====================

// CreateSubscription creates a new webhook subscription
// POST /api/v1/webhooks/subscribe
func (h *WebhookHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    var req createSubscriptionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error(), requestID)
        return
    }

    if req.EventType == "" {
        util.WriteError(w, http.StatusBadRequest, "Event type is required", "", requestID)
        return
    }
    if req.WebhookURL == "" {
        util.WriteError(w, http.StatusBadRequest, "Webhook URL is required", "", requestID)
        return
    }

    validEventTypes := map[string]bool{
        "inventory.changed":  true,
        "forecast.generated": true,
        "anomaly.detected":   true,
    }
    if !validEventTypes[req.EventType] {
        util.WriteError(w, http.StatusBadRequest, "Invalid event type", "Must be: inventory.changed, forecast.generated, anomaly.detected", requestID)
        return
    }

    parsedURL, err := url.Parse(req.WebhookURL)
    if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
        util.WriteError(w, http.StatusBadRequest, "Invalid webhook URL", "URL must be valid and include scheme (http/https)", requestID)
        return
    }
    if parsedURL.Scheme != "https" && !strings.Contains(parsedURL.Host, "localhost") && !strings.Contains(parsedURL.Host, "127.0.0.1") {
        util.WriteError(w, http.StatusBadRequest, "Webhook URL must use HTTPS", "For security, webhooks must use HTTPS in production", requestID)
        return
    }

    maxRetries := req.MaxRetries
    if maxRetries <= 0 {
        maxRetries = 5
    }
    retryInterval := req.RetryInterval
    if retryInterval <= 0 {
        retryInterval = 60
    }

    secretBytes := make([]byte, 32)
    if _, err := rand.Read(secretBytes); err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to generate secret key", err.Error(), requestID)
        return
    }
    secretKey := "whsec_" + base64.URLEncoding.EncodeToString(secretBytes)

    userUUID, err := uuid.Parse(userInfo.UserID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Invalid user ID", err.Error(), requestID)
        return
    }

    sub := &model.WebhookSubscription{
        UserID:         userUUID,
        EventType:      req.EventType,
        WebhookURL:     req.WebhookURL,
        SecretKey:      secretKey,
        IsActive:       true,
        MaxRetries:     maxRetries,
        RetryInterval:  retryInterval,
        CreatedAt:      time.Now().UTC(),
        UpdatedAt:      time.Now().UTC(),
    }

    saved, err := h.repo.CreateSubscription(r.Context(), sub)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to create subscription", err.Error(), requestID)
        return
    }

    resp := subscriptionResponse{
        ID:         saved.ID.String(),
        EventType:  saved.EventType,
        WebhookURL: saved.WebhookURL,
        IsActive:   saved.IsActive,
        SecretKey:  saved.SecretKey,
        CreatedAt:  saved.CreatedAt,
        UpdatedAt:  saved.UpdatedAt,
    }

    util.WriteResponse(w, http.StatusCreated, "Subscription created", resp, requestID)
}

// ==================== ENDPOINT 2: GET /api/v1/webhooks/subscriptions ====================

// ListSubscriptions retrieves all webhook subscriptions for the authenticated user
// GET /api/v1/webhooks/subscriptions
func (h *WebhookHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    userUUID, err := uuid.Parse(userInfo.UserID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Invalid user ID", err.Error(), requestID)
        return
    }

    subs, err := h.repo.GetSubscriptions(r.Context(), userUUID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve subscriptions", err.Error(), requestID)
        return
    }

    items := make([]subscriptionResponse, len(subs))
    for i, sub := range subs {
        items[i] = subscriptionResponse{
            ID:         sub.ID.String(),
            EventType:  sub.EventType,
            WebhookURL: sub.WebhookURL,
            IsActive:   sub.IsActive,
            CreatedAt:  sub.CreatedAt,
            UpdatedAt:  sub.UpdatedAt,
        }
    }

    data := map[string]interface{}{
        "subscriptions": items,
    }

    util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== ENDPOINT 3: DELETE /api/v1/webhooks/subscriptions/{id} ====================

// DeleteSubscription deletes a webhook subscription
// DELETE /api/v1/webhooks/subscriptions/{id}
func (h *WebhookHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    path := strings.TrimPrefix(r.URL.Path, "/api/v1/webhooks/subscriptions/")
    if path == "" {
        util.WriteError(w, http.StatusBadRequest, "Subscription ID required", "", requestID)
        return
    }
    subID, err := uuid.Parse(path)
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid subscription ID", err.Error(), requestID)
        return
    }

    userUUID, err := uuid.Parse(userInfo.UserID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Invalid user ID", err.Error(), requestID)
        return
    }

    subs, err := h.repo.GetSubscriptions(r.Context(), userUUID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to verify ownership", err.Error(), requestID)
        return
    }

    ownsSubscription := false
    for _, sub := range subs {
        if sub.ID == subID {
            ownsSubscription = true
            break
        }
    }

    if !ownsSubscription {
        util.WriteError(w, http.StatusForbidden, "Forbidden", "You do not own this subscription", requestID)
        return
    }

    if err := h.repo.DeleteSubscription(r.Context(), subID); err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to delete subscription", err.Error(), requestID)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// ==================== ENDPOINT 4: GET /api/v1/webhooks/deliveries ====================

// ListDeliveries retrieves webhook delivery history with filters
// GET /api/v1/webhooks/deliveries
func (h *WebhookHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    subscriptionIDStr := r.URL.Query().Get("subscription_id")
    statusFilter := r.URL.Query().Get("status")
    limit, offset := util.ParsePagination(r)

    if limit > 100 {
        limit = 100
    }

    var deliveries []model.WebhookDelivery
    var total int64

    if subscriptionIDStr != "" {
        subID, err := uuid.Parse(subscriptionIDStr)
        if err != nil {
            util.WriteError(w, http.StatusBadRequest, "Invalid subscription ID", err.Error(), requestID)
            return
        }

        // Verify ownership
        userUUID, err := uuid.Parse(userInfo.UserID)
        if err != nil {
            util.WriteError(w, http.StatusInternalServerError, "Invalid user ID", err.Error(), requestID)
            return
        }

        subs, err := h.repo.GetSubscriptions(r.Context(), userUUID)
        if err != nil {
            util.WriteError(w, http.StatusInternalServerError, "Failed to verify ownership", err.Error(), requestID)
            return
        }

        owns := false
        for _, sub := range subs {
            if sub.ID == subID {
                owns = true
                break
            }
        }

        if !owns {
            util.WriteError(w, http.StatusForbidden, "Forbidden", "You do not own this subscription", requestID)
            return
        }

        deliveries, total, err = h.repo.GetDeliveryHistory(r.Context(), subID, limit, offset)
        if err != nil {
            util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve delivery history", err.Error(), requestID)
            return
        }
    } else {
        // If no subscription_id, return empty for MVP
        deliveries = []model.WebhookDelivery{}
        total = 0
    }

    // Filter by status (in-memory)
    filtered := []model.WebhookDelivery{}
    for _, d := range deliveries {
        if statusFilter != "" && string(d.Status) != statusFilter {
            continue
        }
        filtered = append(filtered, d)
    }
    total = int64(len(filtered))

    items := make([]deliveryResponse, len(filtered))
    for i, d := range filtered {
        items[i] = deliveryResponse{
            ID:             d.ID.String(),
            SubscriptionID: d.SubscriptionID.String(),
            EventID:        d.EventID,
            Status:         string(d.Status),
            HTTPStatusCode: d.HTTPStatusCode,
            ErrorMessage:   d.ErrorMessage,
            AttemptCount:   d.AttemptCount,
            CreatedAt:      d.CreatedAt,
            UpdatedAt:      d.UpdatedAt,
        }
    }

    data := map[string]interface{}{
        "deliveries": items,
        "total":      total,
        "limit":      limit,
        "offset":     offset,
    }

    util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}