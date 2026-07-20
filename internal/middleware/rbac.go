package middleware

import (
"net/http"

"github.com/LittleMardiah/scinel/internal/util"
)

// RBACMiddleware checks if the authenticated user has the required role
func RBACMiddleware(requiredRoles ...string) func(http.Handler) http.Handler {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Get user info from context (set by AuthMiddleware)
userInfo, ok := r.Context().Value(UserContextKey).(*UserInfo)
if !ok {
util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", util.GenerateRequestID())
return
}

// Check if user's role is in the required roles list
hasRole := false
for _, role := range requiredRoles {
if userInfo.Role == role {
hasRole = true
break
}
}

if !hasRole {
util.WriteError(w, http.StatusForbidden, "Forbidden", "Insufficient permissions", util.GenerateRequestID())
return
}

next.ServeHTTP(w, r)
})
}
}