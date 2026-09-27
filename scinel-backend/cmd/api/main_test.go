package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LittleMardiah/scinel/internal/middleware"
	"github.com/LittleMardiah/scinel/internal/util"
)

const (
	testJWTSecret = "composition-test-secret"
	testUserID    = "11111111-2222-3333-4444-555555555555"
)

// newCompositionChain builds the production middleware stack from main.go:146-153
// (CORS -> SecurityHeaders -> Auth -> RateLimit -> mux) on top of a stub mux, so
// the tests below can observe each layer without a database or real repositories.
//
// The order in the stack is additionally asserted against main.go's own source by
// TestMainComposition_MiddlewareOrderInMainSource, so this helper cannot silently
// drift away from the wiring that main() actually installs.
func newCompositionChain(t *testing.T) (http.Handler, *miniredis.Miniredis, *util.JWTManager, *int) {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = redisClient.Close() })

	jwtManager := util.NewJWTManager(testJWTSecret, time.Hour, 24*time.Hour)

	muxCalls := 0
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/inventory", func(w http.ResponseWriter, r *http.Request) {
		muxCalls++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":200,"message":"OK"}`))
	})
	// /health is in AuthMiddleware's publicPaths list (auth.go:60).
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		muxCalls++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":200,"message":"OK"}`))
	})

	rateLimiter := middleware.NewRateLimitMiddleware(redisClient, 1000)
	handler := middleware.CORSMiddleware(
		middleware.SecurityHeaders(
			middleware.AuthMiddleware(jwtManager, redisClient)(
				rateLimiter.Apply(mux),
			),
		),
	)

	return handler, mr, jwtManager, &muxCalls
}

func rateLimitKeys(mr *miniredis.Miniredis) []string {
	var keys []string
	for _, k := range mr.Keys() {
		if strings.HasPrefix(k, "rate_limit:") {
			keys = append(keys, k)
		}
	}
	return keys
}

// Test 1: a protected route without an Authorization header is rejected by the
// auth layer with 401, and never reaches the mux.
func TestMainComposition_ProtectedRouteWithoutAuthReturns401(t *testing.T) {
	handler, _, _, muxCalls := newCompositionChain(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, 0, *muxCalls, "auth must reject before the mux handler runs")
}

// Test 2: a public route without an Authorization header bypasses auth and
// reaches the handler with 200.
func TestMainComposition_PublicRouteWithoutAuthReturns200(t *testing.T) {
	handler, _, _, muxCalls := newCompositionChain(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 1, *muxCalls, "public route must bypass auth and hit the mux")
}

// Test 3a: an invalid token is rejected by auth, which sits outside the rate
// limiter, so no rate_limit:* key is ever created for the rejected request.
func TestMainComposition_AuthRunsBeforeRateLimit(t *testing.T) {
	handler, mr, _, muxCalls := newCompositionChain(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, 0, *muxCalls)
	assert.Empty(t, rateLimitKeys(mr), "request rejected by auth must not consume a rate limit bucket")
}

// Test 3b: a valid token passes auth and is then rate limited, proving the rate
// limiter is reachable only after auth and keys the bucket by authenticated user.
func TestMainComposition_RateLimitAppliesAfterSuccessfulAuth(t *testing.T) {
	handler, mr, jwtManager, muxCalls := newCompositionChain(t)

	token, err := jwtManager.GenerateAccessToken(uuid.MustParse(testUserID), "chain@test.dev", "admin")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 1, *muxCalls)
	assert.Equal(t, []string{"rate_limit:user:" + testUserID}, rateLimitKeys(mr))
	assert.Equal(t, "1000", rr.Header().Get("X-RateLimit-Limit"))
}

// TestMainComposition_MiddlewareOrderInMainSource asserts the order main() really
// installs, by parsing main.go:146-153. Without this, the helper above would only
// test a copy of the stack, and reordering the wiring in main() would go unnoticed.
func TestMainComposition_MiddlewareOrderInMainSource(t *testing.T) {
	chain, innermost := middlewareChainInMain(t)

	assert.Equal(t, []string{"CORSMiddleware", "SecurityHeaders", "AuthMiddleware", "Apply"}, chain,
		"main.go must install CORS -> SecurityHeaders -> Auth -> RateLimit -> mux")
	assert.Equal(t, "mux", innermost, "the innermost handler of the chain must be the mux")
}

// middlewareChainInMain returns the middleware call names of the handler chain
// built in main(), outermost first, plus the identifier the chain terminates in.
func middlewareChainInMain(t *testing.T) (chain []string, innermost string) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", nil, 0)
	require.NoError(t, err, "main.go must stay parseable")

	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			return true
		}
		outer, ok := assign.Rhs[0].(*ast.CallExpr)
		if !ok || callName(outer) != "CORSMiddleware" {
			return true
		}

		found = true
		chain = nil
		expr := ast.Expr(outer)
		for {
			call, ok := expr.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				break
			}
			chain = append(chain, callName(call))
			expr = call.Args[0]
		}
		if id, ok := expr.(*ast.Ident); ok {
			innermost = id.Name
		}
		return false
	})
	require.True(t, found, "no CORSMiddleware(...) assignment found in main.go")
	return chain, innermost
}

// callName returns the callee name of a call, using the selector for method
// calls. Middleware constructors are applied in curried form
// (AuthMiddleware(jwtManager, redisClient)(handler)), so a call whose Fun is
// itself a call is named after the inner constructor.
func callName(call *ast.CallExpr) string {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name
	case *ast.SelectorExpr:
		return fun.Sel.Name
	case *ast.CallExpr:
		return callName(fun)
	default:
		return ""
	}
}
