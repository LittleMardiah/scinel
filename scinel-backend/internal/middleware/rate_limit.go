package middleware

import (
    "context"
    "fmt"
    "net"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/LittleMardiah/scinel/internal/util"
)

// RateLimitMiddleware applies rate limiting (Redis-based)
// Rate limit: 1000 requests per second per user
// Key: user_id or api_key (for auth'd requests)
// Alternative for unauthenticated: IP address (not recommended for production)
type RateLimitMiddleware struct {
    redis          *redis.Client
    requestsPerSec int
    windowSize     time.Duration
}

// NewRateLimitMiddleware creates rate limiter
func NewRateLimitMiddleware(redisClient *redis.Client, requestsPerSec int) *RateLimitMiddleware {
    return &RateLimitMiddleware{
        redis:          redisClient,
        requestsPerSec: requestsPerSec,
        windowSize:     1 * time.Second, // 1-second sliding window
    }
}

// Apply rate limiting to request
func (rl *RateLimitMiddleware) Apply(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := util.GenerateRequestID()
        ctx := r.Context()

        // Determine rate limit key (from auth token or API key)
        // Priority: API-Key header > Authorization bearer token > IP address
        limitKey := rl.getLimitKey(r)

        // Check rate limit
        allowed, remaining, retryAfter, err := rl.checkLimit(ctx, limitKey)
        if err != nil {
            // Redis error: Log but don't block request (rate limit is optional)
            // Continue to next handler
            next.ServeHTTP(w, r)
            return
        }

        // Add rate limit info to response headers
        w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.requestsPerSec))
        w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
        w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Duration(retryAfter)*time.Second).Unix()))

        // Block if limit exceeded
        if !allowed {
            w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
            util.WriteError(w, http.StatusTooManyRequests, "Rate limit exceeded",
                fmt.Sprintf("You have exceeded %d requests per second. Retry after %d seconds", rl.requestsPerSec, retryAfter),
                requestID)
            return
        }

        // Continue to next handler
        next.ServeHTTP(w, r)
    })
}

// checkLimit checks if request is within rate limit
func (rl *RateLimitMiddleware) checkLimit(ctx context.Context, key string) (allowed bool, remaining int, retryAfter int, err error) {
    ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
    defer cancel()

    redisKey := fmt.Sprintf("rate_limit:%s", key)

    // Get current count
    count, err := rl.redis.Incr(ctx, redisKey).Result()
    if err != nil && err != redis.Nil {
        return true, 0, 0, err // Redis error, allow request
    }

    // If first request in window, set expiration
    if count == 1 {
        rl.redis.Expire(ctx, redisKey, rl.windowSize)
    }

    // Check if limit exceeded
    if int(count) > rl.requestsPerSec {
        ttl, _ := rl.redis.TTL(ctx, redisKey).Result()
        retryAfter = int(ttl.Seconds()) + 1 // ✅ FIX: pake = bukan :=
        if retryAfter <= 0 {
            retryAfter = 1
        }
        return false, 0, retryAfter, nil
    }

    remaining = rl.requestsPerSec - int(count)
    return true, remaining, 0, nil
}

// getLimitKey extracts limit key from request
func (rl *RateLimitMiddleware) getLimitKey(r *http.Request) string {
    // Check for API key in header
    if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
        return apiKey
    }

    // Check for user ID in Bearer token (from JWT claims)
    // User ID is stored in context by AuthMiddleware
    if userInfo, ok := r.Context().Value(UserContextKey).(*UserInfo); ok && userInfo != nil {
        return "user:" + userInfo.UserID
    }

    // Fallback to IP address (not ideal)
    return "ip:" + clientIP(r)
}

// clientIP returns the address used to bucket unauthenticated requests.
//
// X-Forwarded-For is an append-only chain ("client, proxy1, proxy2"), so only
// the first hop is the originating client. Keying on the whole chain would
// merge every client behind the proxy into a single rate limit bucket. When
// the header is absent, the TCP peer address is used with its port stripped.
func clientIP(r *http.Request) string {
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        if first := strings.TrimSpace(strings.Split(xff, ",")[0]); first != "" {
            return first
        }
    }
    if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
        return host
    }
    return r.RemoteAddr
}