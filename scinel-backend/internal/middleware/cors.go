package middleware

import (
"net/http"
"os"
"strings"
)

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
origin := r.Header.Get("Origin")

// Allow if origin is in whitelist
for _, allowed := range allowedOrigins {
if strings.TrimSpace(allowed) == origin {
w.Header().Set("Access-Control-Allow-Origin", origin)
break
}
}

// Set common CORS headers
w.Header().Set("Access-Control-Allow-Methods", os.Getenv("CORS_ALLOWED_METHODS"))
w.Header().Set("Access-Control-Allow-Headers", os.Getenv("CORS_ALLOWED_HEADERS"))
w.Header().Set("Access-Control-Max-Age", "86400")

// Handle preflight (OPTIONS) requests
if r.Method == http.MethodOptions {
w.WriteHeader(http.StatusOK)
return
}

next.ServeHTTP(w, r)
})
}

// SecurityHeaders adds security headers to every response
func SecurityHeaders(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Prevent MIME type sniffing
w.Header().Set("X-Content-Type-Options", "nosniff")
// Prevent clickjacking
w.Header().Set("X-Frame-Options", "DENY")
// Enable XSS protection
w.Header().Set("X-XSS-Protection", "1; mode=block")
// Enforce HTTPS (for production)
w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

next.ServeHTTP(w, r)
})
}