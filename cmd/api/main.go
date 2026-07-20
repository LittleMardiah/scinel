package main

import (
"context"
"fmt"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"github.com/LittleMardiah/scinel/internal/config"
"github.com/LittleMardiah/scinel/internal/database"
"github.com/LittleMardiah/scinel/internal/handler"
"github.com/LittleMardiah/scinel/internal/middleware"
"github.com/LittleMardiah/scinel/internal/repository"
"github.com/LittleMardiah/scinel/internal/util"
"github.com/LittleMardiah/scinel/pkg/logger"
)

func main() {
// 1. Load Configuration
cfg, err := config.LoadConfig()
if err != nil {
log.Fatalf("Failed to load config: %v", err)
}

// 2. Initialize Logger
appLogger := logger.NewLogger(cfg.Logging.Level)
appLogger.Info("Starting ScIntel application",
"environment", cfg.Environment,
"port", cfg.Server.Port,
)

// 3. Initialize Database Connection
db, err := database.NewConnection(cfg.Database)
if err != nil {
appLogger.Error("Failed to connect to database", "error", err.Error())
os.Exit(1)
}
defer db.Close()
appLogger.Info("Database connection established")

// 4. Initialize JWT Manager
jwtManager := util.NewJWTManager(
cfg.JWT.SecretKey,
cfg.JWT.AccessExpiry,
cfg.JWT.RefreshExpiry,
)

// 5. Initialize Repository
userRepo := repository.NewUserRepository(db)

// 6. Initialize Handler
authHandler := handler.NewAuthHandler(userRepo, jwtManager, cfg)

// 7. Setup Router
mux := http.NewServeMux()

// Public routes (no auth required)
mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.Refresh)
mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)

// Health check endpoints
mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"status":"ok"}`))
})
mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
// Check database connectivity
if err := db.PingContext(r.Context()); err != nil {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusServiceUnavailable)
w.Write([]byte(`{"status":"not_ready","reason":"database_down"}`))
return
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"status":"ready"}`))
})

// Protected route example (requires auth)
mux.HandleFunc("GET /api/v1/protected", func(w http.ResponseWriter, r *http.Request) {
userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
if !ok {
w.WriteHeader(http.StatusInternalServerError)
return
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(fmt.Sprintf(`{"message":"Hello %s, you have role: %s"}`, userInfo.Email, userInfo.Role)))
})

// 8. Apply Global Middleware
// Order: CORS -> Security Headers -> Auth -> Logging (if we had it)
// We apply auth globally, but it skips public paths internally.
handlerWithMiddleware := middleware.CORSMiddleware(
middleware.SecurityHeaders(
middleware.AuthMiddleware(jwtManager)(
mux,
),
),
)

// 9. Create HTTP Server
server := &http.Server{
Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
Handler:      handlerWithMiddleware,
ReadTimeout:  cfg.Server.ReadTimeout,
WriteTimeout: cfg.Server.WriteTimeout,
}

// 10. Start Server in Goroutine
go func() {
appLogger.Info("HTTP server listening", "address", server.Addr)
if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
appLogger.Error("Server error", "error", err.Error())
}
}()

// 11. Graceful Shutdown Handling
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// Block until signal received
<-sigChan
appLogger.Info("Shutdown signal received, shutting down gracefully...")

ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
appLogger.Error("Server forced to shutdown", "error", err.Error())
}

appLogger.Info("Server shutdown complete")
}