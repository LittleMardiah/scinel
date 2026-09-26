package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/LittleMardiah/scinel/internal/config"
	"github.com/LittleMardiah/scinel/internal/middleware"
	"github.com/LittleMardiah/scinel/internal/model"
	"github.com/LittleMardiah/scinel/internal/repository"
	"github.com/LittleMardiah/scinel/internal/service"
	"github.com/LittleMardiah/scinel/internal/util"
)

// AuthDeps carries the optional collaborators of AuthHandler.
//
// It is passed as a variadic argument so existing call sites
// (handler.NewAuthHandler(userRepo, jwtManager, cfg)) keep compiling: with no
// AuthDeps the handler behaves exactly as before, i.e. no audit writes and no
// token blacklist.
type AuthDeps struct {
	// AuditSvc receives LOGIN / LOGOUT audit entries. nil disables auditing.
	AuditSvc *service.AuditService
	// RedisClient backs the token blacklist. nil disables revocation.
	RedisClient *redis.Client
}

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userRepo   repository.UserRepository
	jwtManager *util.JWTManager
	cfg        *config.Config
	auditSvc   *service.AuditService
	redis      *redis.Client
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo repository.UserRepository, jwtManager *util.JWTManager, cfg *config.Config, deps ...AuthDeps) *AuthHandler {
	h := &AuthHandler{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		cfg:        cfg,
	}
	for _, d := range deps {
		if d.AuditSvc != nil {
			h.auditSvc = d.AuditSvc
		}
		if d.RedisClient != nil {
			h.redis = d.RedisClient
		}
	}
	return h
}

// ====================
// REQUEST/RESPONSE TYPES
// ====================

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type tokenResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token,omitempty"`
	TokenType    string              `json:"token_type"`
	ExpiresIn    int                 `json:"expires_in"`
	User         *model.UserResponse `json:"user"`
}

// ====================
// AUDIT + BLACKLIST HELPERS
// ====================

// blacklistTimeout bounds every Redis blacklist round-trip so a slow or
// hanging Redis cannot stall logout/refresh.
const blacklistTimeout = 200 * time.Millisecond

// logAudit appends one audit entry, attributing it to the authenticated user.
// It never fails the request: AuditService.Log is itself non-fatal, and a nil
// service short-circuits.
func (h *AuthHandler) logAudit(ctx context.Context, r *http.Request, action string, details map[string]interface{}) {
	var userID *uuid.UUID
	if info, ok := ctx.Value(middleware.UserContextKey).(*middleware.UserInfo); ok && info != nil {
		if parsed, err := uuid.Parse(info.UserID); err == nil {
			userID = &parsed
		}
	}
	h.logAuditAs(ctx, r, userID, action, details)
}

// logAuditAs appends one audit entry for an explicitly resolved actor. Used by
// /auth/logout, which is a public route and therefore has no auth context.
func (h *AuthHandler) logAuditAs(ctx context.Context, r *http.Request, userID *uuid.UUID, action string, details map[string]interface{}) {
	if h.auditSvc == nil {
		return
	}

	entry := &model.AuditLog{
		Action:    action,
		UserID:    userID,
		IPAddress: optionalString(service.ExtractIP(r)),
		UserAgent: optionalString(service.ExtractUserAgent(r)),
	}
	if raw, err := json.Marshal(details); err == nil {
		entry.Details = raw
	}

	h.auditSvc.Log(ctx, entry)
}

// bearerToken extracts the raw token from an Authorization header.
func bearerToken(r *http.Request) (string, bool) {
	if r == nil {
		return "", false
	}
	parts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

// blacklistCtx bounds the Redis round-trip of the revocation check.
func (h *AuthHandler) blacklistCtx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), blacklistTimeout)
}

// revokeToken records the token digest in the blacklist for the remainder of
// its own lifetime. Best effort: a Redis failure is logged and the request
// still succeeds, because logout must not be able to lock a user out of the
// endpoint that tells them the session ended.
func (h *AuthHandler) revokeToken(ctx context.Context, rawToken string, expiresAt time.Time) {
	if h.redis == nil || rawToken == "" {
		return
	}
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return
	}
	if err := middleware.RevokeToken(ctx, h.redis, rawToken, ttl); err != nil {
		log.Printf("warn: failed to blacklist token on logout: %v", err)
	}
}

// isBlacklisted reports whether the token was revoked on logout. A Redis error
// is deliberately reported as "not blacklisted" (fail-open) so an outage
// cannot take authentication down.
func (h *AuthHandler) isBlacklisted(ctx context.Context, rawToken string) bool {
	if h.redis == nil || rawToken == "" {
		return false
	}
	revoked, err := middleware.IsTokenRevoked(ctx, h.redis, rawToken)
	if err != nil {
		log.Printf("warn: token blacklist lookup failed, allowing request: %v", err)
		return false
	}
	return revoked
}

// optionalString returns nil for an empty string so nullable audit_log columns
// are stored as SQL NULL rather than an empty value.
func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// ====================
// HANDLER METHODS
// ====================

// Register POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error(), requestID)
		return
	}

	// Validasi input
	if req.Email == "" {
		util.WriteError(w, http.StatusBadRequest, "Email is required", "", requestID)
		return
	}
	if req.Password == "" {
		util.WriteError(w, http.StatusBadRequest, "Password is required", "", requestID)
		return
	}
	if len(req.Password) < 8 {
		util.WriteError(w, http.StatusBadRequest, "Password must be at least 8 characters", "", requestID)
		return
	}
	if req.FullName == "" {
		util.WriteError(w, http.StatusBadRequest, "Full name is required", "", requestID)
		return
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password, 0) // 0 = default 12
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to hash password", err.Error(), requestID)
		return
	}

	// Create user (default role: viewer, status: active)
	user, err := h.userRepo.CreateUser(r.Context(), req.Email, hashedPassword, req.FullName, "viewer")
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "already exists") {
			util.WriteError(w, http.StatusConflict, "Email already registered", "", requestID)
			return
		}
		util.WriteError(w, http.StatusInternalServerError, "Failed to create user", err.Error(), requestID)
		return
	}

	util.WriteResponse(w, http.StatusCreated, "User registered successfully", user.ToResponse(), requestID)
}

// Login POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error(), requestID)
		return
	}

	if req.Email == "" || req.Password == "" {
		util.WriteError(w, http.StatusBadRequest, "Email and password are required", "", requestID)
		return
	}

	// Cari user by email
	user, err := h.userRepo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, "Invalid email or password", "", requestID)
		return
	}

	// Verifikasi password
	if err := util.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		util.WriteError(w, http.StatusUnauthorized, "Invalid email or password", "", requestID)
		return
	}

	// Cek status user
	switch user.Status {
	case "suspended":
		util.WriteError(w, http.StatusForbidden, "Account has been suspended", "Please contact administrator", requestID)
		return
	case "pending_email_verification":
		util.WriteError(w, http.StatusForbidden, "Email not verified", "Please verify your email first", requestID)
		return
	case "active":
		// Lanjutkan
	default:
		util.WriteError(w, http.StatusForbidden, "Account status invalid", "", requestID)
		return
	}

	// Update last login
	_ = h.userRepo.UpdateLastLogin(r.Context(), user.ID)

	// Generate tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to generate access token", err.Error(), requestID)
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to generate refresh token", err.Error(), requestID)
		return
	}

	// Audit: LOGIN. Best effort, never changes the response.
	h.logAuditAs(r.Context(), r, &user.ID, model.AuditActionLogin, map[string]interface{}{
		"email": user.Email,
		"role":  user.Role,
	})

	// Response
	resp := tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(h.cfg.JWT.AccessExpiry.Seconds()),
		User:         &model.UserResponse{ID: user.ID, Email: user.Email, FullName: user.FullName, Role: user.Role, Status: user.Status, CreatedAt: user.CreatedAt},
	}

	util.WriteResponse(w, http.StatusOK, "Login successful", resp, requestID)
}

// Refresh POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error(), requestID)
		return
	}

	if req.RefreshToken == "" {
		util.WriteError(w, http.StatusBadRequest, "Refresh token is required", "", requestID)
		return
	}

	// Blacklist guard: a refresh token revoked by Logout must not be able to
	// mint a fresh access token. Checked BEFORE any new token is issued.
	if h.redis != nil {
		ctx, cancel := h.blacklistCtx(r)
		blacklisted := h.isBlacklisted(ctx, req.RefreshToken)
		cancel()
		if blacklisted {
			util.WriteError(w, http.StatusUnauthorized, "Token has been revoked", "Session was logged out", requestID)
			return
		}
	}

	// Verify refresh token
	claims, err := h.jwtManager.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, "Invalid or expired refresh token", err.Error(), requestID)
		return
	}

	// Get user by ID
	user, err := h.userRepo.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, "User not found", err.Error(), requestID)
		return
	}

	if user.Status != "active" {
		util.WriteError(w, http.StatusForbidden, "Account is not active", "", requestID)
		return
	}

	// Generate new access token
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to generate access token", err.Error(), requestID)
		return
	}

	resp := map[string]interface{}{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   int(h.cfg.JWT.AccessExpiry.Seconds()),
	}

	util.WriteResponse(w, http.StatusOK, "Token refreshed successfully", resp, requestID)
}

// Logout POST /api/v1/auth/logout
//
// Both the presented access token and the refresh token from the body are
// blacklisted until their own natural expiry, so neither can be replayed.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	requestID := util.GenerateRequestID()

	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error(), requestID)
		return
	}

	if req.RefreshToken == "" {
		util.WriteError(w, http.StatusBadRequest, "Refresh token is required", "", requestID)
		return
	}

	// 1. Resolve each token's natural expiry (the blacklist TTL must match it,
	//    otherwise revoked keys outlive the token or expire too early).
	var refreshExpiry time.Time
	var refreshUserID *uuid.UUID
	if claims, err := h.jwtManager.VerifyRefreshToken(req.RefreshToken); err == nil {
		if claims.ExpiresAt != nil {
			refreshExpiry = claims.ExpiresAt.Time
		}
		refreshUserID = &claims.UserID
	}

	var accessExpiry time.Time
	if rawToken, ok := bearerToken(r); ok {
		if claims, err := h.jwtManager.VerifyAccessToken(rawToken); err == nil && claims.ExpiresAt != nil {
			accessExpiry = claims.ExpiresAt.Time
		}
	}

	// 2. Blacklist both tokens.
	revoked := 0
	if h.redis != nil {
		ctx, cancel := h.blacklistCtx(r)
		if !refreshExpiry.IsZero() {
			h.revokeToken(ctx, req.RefreshToken, refreshExpiry)
			revoked++
		}
		if rawToken, ok := bearerToken(r); ok && !accessExpiry.IsZero() {
			h.revokeToken(ctx, rawToken, accessExpiry)
			revoked++
		}
		cancel()
	}

	// 3. Audit: LOGOUT, attributed to the refresh token claims because this
	//    route is public and therefore carries no auth context.
	details := map[string]interface{}{"tokens_revoked": revoked}
	if refreshUserID != nil {
		details["user_id"] = refreshUserID.String()
	}
	h.logAuditAs(r.Context(), r, refreshUserID, model.AuditActionLogout, details)

	// Always 200: an unverifiable token is still a logout from the client's
	// point of view, and 4xx here would leak token validity.
	util.WriteResponse(w, http.StatusOK, "Logged out successfully", nil, requestID)
}
