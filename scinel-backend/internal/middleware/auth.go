package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/LittleMardiah/scinel/internal/util"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserContextKey is the key for user info in context
	UserContextKey ContextKey = "user"
)

// blacklistKeyPrefix namespaces revoked-token digests in Redis.
const blacklistKeyPrefix = "blacklist:"

// blacklistTimeout bounds the Redis round-trip of the revocation check so a
// slow (or hanging) Redis cannot stall every authenticated request.
const blacklistTimeout = 200 * time.Millisecond

// UserInfo holds user data extracted from JWT
type UserInfo struct {
	UserID string
	Email  string
	Role   string
}

// AuthMiddleware validates JWT tokens for protected routes.
//
// redisClient is variadic so existing call sites (and unit tests) keep
// compiling. When it is supplied, the middleware also rejects tokens that were
// revoked through POST /api/v1/auth/logout.
func AuthMiddleware(jwtManager *util.JWTManager, redisClient ...*redis.Client) func(http.Handler) http.Handler {
	var revoked *redis.Client
	if len(redisClient) > 0 {
		revoked = redisClient[0]
	}

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

			// Reject tokens revoked via logout.
			// fail-open: if no Redis client is wired, or Redis is unreachable,
			// we log a warning and let the request through instead of 401.
			// Availability wins over revocation enforcement here — a Redis
			// outage must not log every user out of the whole API.
			if revoked != nil {
				ctx, cancel := context.WithTimeout(r.Context(), blacklistTimeout)
				revokedFlag, err := IsTokenRevoked(ctx, revoked, tokenString)
				cancel()
				switch {
				case err != nil:
					// fail-open (see comment above): log warn, continue.
					log.Printf("warn: token blacklist unavailable, allowing request: %v", err)
				case revokedFlag:
					util.WriteError(w, http.StatusUnauthorized, "Token has been revoked", "Session was logged out", util.GenerateRequestID())
					return
				}
			}

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

// IsTokenRevoked reports whether the given token has a digest stored under
// blacklist:{sha256(token)}.
//
// The raw token is never used as a Redis key: only its SHA-256 digest is
// stored, so a Redis dump cannot be replayed against the API.
func IsTokenRevoked(ctx context.Context, client *redis.Client, token string) (bool, error) {
	if client == nil {
		return false, nil
	}
	res, err := client.Get(ctx, BlacklistKey(token)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return res != "", nil
}

// BlacklistKey builds the Redis key holding the revocation marker of a token.
func BlacklistKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return blacklistKeyPrefix + hex.EncodeToString(sum[:])
}

// RevokeToken stores the token's digest in the blacklist for the remainder of
// its lifetime, so the marker disappears exactly when the token expires and
// the keyspace cannot grow without bound.
//
// A nil client or a non-positive ttl is a no-op: the caller keeps its previous
// behaviour (token stays valid) instead of failing the logout request.
func RevokeToken(ctx context.Context, client *redis.Client, token string, ttl time.Duration) error {
	if client == nil || token == "" || ttl <= 0 {
		return nil
	}
	if err := client.Set(ctx, BlacklistKey(token), "1", ttl).Err(); err != nil {
		return err
	}
	return nil
}
