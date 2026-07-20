package handler

import (
    "crypto/rand"
    "encoding/base64"
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

// APIKeyHandler handles API key endpoints
type APIKeyHandler struct {
    repo repository.APIKeyRepository
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(repo repository.APIKeyRepository) *APIKeyHandler {
    return &APIKeyHandler{repo: repo}
}

// ==================== REQUEST/RESPONSE TYPES ====================

type createAPIKeyRequest struct {
    Name          string   `json:"name"`
    Permissions   []string `json:"permissions"`
    ExpiresInDays int      `json:"expires_in_days,omitempty"`
}

type apiKeyResponse struct {
    KeyID       string    `json:"key_id"`
    KeyPrefix   string    `json:"key_prefix"`
    KeySecret   string    `json:"key_secret,omitempty"` // only on create/rotate
    Name        string    `json:"name"`
    Permissions []string  `json:"permissions"`
    IsActive    bool      `json:"is_active"`
    LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type rotateResponse struct {
    KeyID     string    `json:"key_id"`
    KeyPrefix string    `json:"key_prefix"`
    KeySecret string    `json:"key_secret"`
    RotatedAt time.Time `json:"rotated_at"`
}

// ==================== ENDPOINT 1: POST /api/v1/api-keys ====================

// CreateAPIKey generates a new API key
// POST /api/v1/api-keys
func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    var req createAPIKeyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error(), requestID)
        return
    }

    if req.Name == "" {
        util.WriteError(w, http.StatusBadRequest, "Name is required", "", requestID)
        return
    }
    if len(req.Permissions) == 0 {
        util.WriteError(w, http.StatusBadRequest, "At least one permission is required", "", requestID)
        return
    }

    // Validate permissions (simple check)
    validPermissions := map[string]bool{
        "inventory:read":  true,
        "inventory:write": true,
        "forecast:read":   true,
        "forecast:write":  true,
        "anomaly:read":    true,
        "anomaly:write":   true,
        "webhook:manage":  true,
    }
    for _, p := range req.Permissions {
        if !validPermissions[p] {
            util.WriteError(w, http.StatusBadRequest, "Invalid permission: "+p, "Allowed: inventory:read, inventory:write, forecast:read, forecast:write, anomaly:read, anomaly:write, webhook:manage", requestID)
            return
        }
    }

    userUUID, err := uuid.Parse(userInfo.UserID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Invalid user ID", err.Error(), requestID)
        return
    }

    // Generate key prefix (first 6 chars of random)
    randomBytes := make([]byte, 6)
    if _, err := rand.Read(randomBytes); err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to generate key", err.Error(), requestID)
        return
    }
    prefix := "sk_test_" + base64.URLEncoding.EncodeToString(randomBytes)[:6]

    // Generate full secret (32 bytes)
    secretBytes := make([]byte, 32)
    if _, err := rand.Read(secretBytes); err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to generate secret", err.Error(), requestID)
        return
    }
    fullSecret := prefix + "_" + base64.URLEncoding.EncodeToString(secretBytes)

    // Hash full secret for storage (bcrypt would be better, but we'll use SHA256 for speed)
    // For MVP, we'll store the full secret hashed with SHA256 (in production, use bcrypt)
    keyHash := util.HashSHA256(fullSecret) // kita tambahkan helper

    // Calculate expiry
    var expiresAt *time.Time
    if req.ExpiresInDays > 0 {
        t := time.Now().UTC().AddDate(0, 0, req.ExpiresInDays)
        expiresAt = &t
    } else {
        // Default 365 days
        t := time.Now().UTC().AddDate(1, 0, 0)
        expiresAt = &t
    }

    // Create key record
    key := &model.APIKey{
        UserID:      userUUID,
        KeyPrefix:   prefix,
        KeyHash:     keyHash,
        Name:        req.Name,
        Permissions: req.Permissions,
        IsActive:    true,
        ExpiresAt:   expiresAt,
        CreatedAt:   time.Now().UTC(),
        UpdatedAt:   time.Now().UTC(),
    }

    saved, err := h.repo.Create(r.Context(), key)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to create API key", err.Error(), requestID)
        return
    }

    // Key secret only shown once
    resp := apiKeyResponse{
        KeyID:       saved.ID.String(),
        KeyPrefix:   saved.KeyPrefix,
        KeySecret:   fullSecret, // full secret only here
        Name:        saved.Name,
        Permissions: saved.Permissions,
        IsActive:    saved.IsActive,
        CreatedAt:   saved.CreatedAt,
        ExpiresAt:   saved.ExpiresAt,
    }

    util.WriteResponse(w, http.StatusCreated, "API key created", resp, requestID)
}

// ==================== ENDPOINT 2: GET /api/v1/api-keys ====================

// ListAPIKeys retrieves all API keys for the authenticated user
// GET /api/v1/api-keys
func (h *APIKeyHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
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

    keys, err := h.repo.GetByUser(r.Context(), userUUID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve API keys", err.Error(), requestID)
        return
    }

    items := make([]apiKeyResponse, len(keys))
    for i, k := range keys {
        items[i] = apiKeyResponse{
            KeyID:       k.ID.String(),
            KeyPrefix:   k.KeyPrefix,
            Name:        k.Name,
            Permissions: k.Permissions,
            IsActive:    k.IsActive,
            LastUsedAt:  k.LastUsedAt,
            CreatedAt:   k.CreatedAt,
            ExpiresAt:   k.ExpiresAt,
        }
    }

    data := map[string]interface{}{"keys": items}
    util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== ENDPOINT 3: DELETE /api/v1/api-keys/{id} ====================

// RevokeAPIKey revokes an API key (sets is_active = false)
// DELETE /api/v1/api-keys/{id}
func (h *APIKeyHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    path := strings.TrimPrefix(r.URL.Path, "/api/v1/api-keys/")
    if path == "" {
        util.WriteError(w, http.StatusBadRequest, "Key ID required", "", requestID)
        return
    }
    keyID, err := uuid.Parse(path)
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid key ID", err.Error(), requestID)
        return
    }

    userUUID, err := uuid.Parse(userInfo.UserID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Invalid user ID", err.Error(), requestID)
        return
    }

    if err := h.repo.Delete(r.Context(), keyID, userUUID); err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to revoke API key", err.Error(), requestID)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// ==================== ENDPOINT 4: POST /api/v1/api-keys/{id}/rotate ====================

// RotateAPIKey generates a new secret for an existing API key
// POST /api/v1/api-keys/{id}/rotate
func (h *APIKeyHandler) RotateAPIKey(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    path := strings.TrimPrefix(r.URL.Path, "/api/v1/api-keys/")
    if path == "" {
        util.WriteError(w, http.StatusBadRequest, "Key ID required", "", requestID)
        return
    }
    parts := strings.Split(path, "/")
    if len(parts) != 2 || parts[1] != "rotate" {
        util.WriteError(w, http.StatusBadRequest, "Invalid path", "Expected /api/v1/api-keys/{id}/rotate", requestID)
        return
    }
    keyID, err := uuid.Parse(parts[0])
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid key ID", err.Error(), requestID)
        return
    }

    userUUID, err := uuid.Parse(userInfo.UserID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Invalid user ID", err.Error(), requestID)
        return
    }

    // Get existing key
    existing, err := h.repo.GetByID(r.Context(), keyID, userUUID)
    if err != nil {
        util.WriteError(w, http.StatusNotFound, "API key not found", err.Error(), requestID)
        return
    }

    if !existing.IsActive {
        util.WriteError(w, http.StatusConflict, "Cannot rotate inactive key", "Key is already revoked", requestID)
        return
    }

    // Generate new secret
    secretBytes := make([]byte, 32)
    if _, err := rand.Read(secretBytes); err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to generate secret", err.Error(), requestID)
        return
    }
    fullSecret := existing.KeyPrefix + "_" + base64.URLEncoding.EncodeToString(secretBytes)
    newHash := util.HashSHA256(fullSecret)

    // Update key hash
    existing.KeyHash = newHash
    existing.UpdatedAt = time.Now().UTC()
    if err := h.repo.Update(r.Context(), existing); err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to rotate API key", err.Error(), requestID)
        return
    }

    resp := rotateResponse{
        KeyID:     existing.ID.String(),
        KeyPrefix: existing.KeyPrefix,
        KeySecret: fullSecret,
        RotatedAt: existing.UpdatedAt,
    }

    util.WriteResponse(w, http.StatusCreated, "API key rotated", resp, requestID)
}