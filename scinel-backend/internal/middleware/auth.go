package middleware

import (
"context"
"net/http"
"strings"

"github.com/LittleMardiah/scinel/internal/util"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
// UserContextKey is the key for user info in context
UserContextKey ContextKey = "user"
)

// UserInfo holds user data extracted from JWT
type UserInfo struct {
UserID string
Email  string
Role   string
}

// AuthMiddleware validates JWT tokens for protected routes
func AuthMiddleware(jwtManager *util.JWTManager) func(http.Handler) http.Handler {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Public paths that don't require authentication
publicPaths := []string{
    "/api/v1/auth/register",
    "/api/v1/auth/login",
    "/api/v1/auth/refresh",
    "/api/v1/health",
    "/api/v1/readiness",
    "/api/v1/liveness",
    "/health",
    "/ready",
}

// Check if path is public
for _, path := range publicPaths {
if r.URL.Path == path {
next.ServeHTTP(w, r)
return
}
}

// Get token from Authorization header
authHeader := r.Header.Get("Authorization")
if authHeader == "" {
util.WriteError(w, http.StatusUnauthorized, "Missing authorization header", "", util.GenerateRequestID())
return
}

// Extract token (Bearer <token>)
parts := strings.Split(authHeader, " ")
if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
util.WriteError(w, http.StatusUnauthorized, "Invalid authorization format", "Expected: Bearer <token>", util.GenerateRequestID())
return
}
tokenString := parts[1]

// Verify token
claims, err := jwtManager.VerifyAccessToken(tokenString)
if err != nil {
util.WriteError(w, http.StatusUnauthorized, "Invalid or expired token", err.Error(), util.GenerateRequestID())
return
}

// Store user info in context
userInfo := &UserInfo{
UserID: claims.UserID.String(),
Email:  claims.Email,
Role:   claims.Role,
}
ctx := context.WithValue(r.Context(), UserContextKey, userInfo)

// Call next handler with updated context
next.ServeHTTP(w, r.WithContext(ctx))
})
}
}