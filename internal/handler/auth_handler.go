package handler

import (
"encoding/json"
"net/http"
"strings"
"github.com/LittleMardiah/scinel/internal/config"
"github.com/LittleMardiah/scinel/internal/model"
"github.com/LittleMardiah/scinel/internal/repository"
"github.com/LittleMardiah/scinel/internal/util"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
userRepo   repository.UserRepository
jwtManager *util.JWTManager
cfg        *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo repository.UserRepository, jwtManager *util.JWTManager, cfg *config.Config) *AuthHandler {
return &AuthHandler{
userRepo:   userRepo,
jwtManager: jwtManager,
cfg:        cfg,
}
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
AccessToken  string      `json:"access_token"`
RefreshToken string      `json:"refresh_token,omitempty"`
TokenType    string      `json:"token_type"`
ExpiresIn    int         `json:"expires_in"`
User         *model.UserResponse `json:"user"`
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

// Verify refresh token (just to make sure it's valid)
_, err := h.jwtManager.VerifyRefreshToken(req.RefreshToken)
if err != nil {
// Still return 200 OK to prevent user enumeration
util.WriteResponse(w, http.StatusOK, "Logged out successfully", nil, requestID)
return
}

// TODO: Add token to blacklist (Redis) for future requests
// For now, client-side should discard the token

util.WriteResponse(w, http.StatusOK, "Logged out successfully", nil, requestID)
}