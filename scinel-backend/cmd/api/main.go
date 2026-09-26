package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"

	"github.com/LittleMardiah/scinel/internal/config"
	"github.com/LittleMardiah/scinel/internal/database"
	"github.com/LittleMardiah/scinel/internal/handler"
	"github.com/LittleMardiah/scinel/internal/middleware"
	"github.com/LittleMardiah/scinel/internal/repository"
	"github.com/LittleMardiah/scinel/internal/service"
	"github.com/LittleMardiah/scinel/internal/util"
	"github.com/LittleMardiah/scinel/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.NewLogger(cfg.Logging.Level)
	appLogger.Info("Starting ScIntel application",
		"environment", cfg.Environment,
		"port", cfg.Server.Port,
	)

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()
	appLogger.Info("Database connection established")

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6380"
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		appLogger.Warn("Failed to parse REDIS_URL, using default", "error", err.Error())
		opt = &redis.Options{
			Addr: "localhost:6380",
			DB:   0,
		}
	}
	redisClient := redis.NewClient(opt)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		appLogger.Warn("Failed to connect to Redis", "error", err.Error())
	} else {
		appLogger.Info("Redis connection established")
	}
	defer redisClient.Close()

	jwtManager := util.NewJWTManager(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

    userRepo := repository.NewUserRepository(db)
    forecastRepo := repository.NewForecastRepository(db)
    invRepo := repository.NewInventoryRepository(db)
    anomalyRepo := repository.NewAnomalyRepository(db)
    apiKeyRepo := repository.NewAPIKeyRepository(db)
    webhookRepo := repository.NewWebhookRepository(db)
    txnRepo := repository.NewTransactionRepository(db)
    auditRepo := repository.NewAuditRepository(db)

    // AuditService is wired before the services that hook into it, because
    // AuditService.Log is non-fatal and its own failures only produce a warn.
    auditSvc := service.NewAuditService(auditRepo, appLogger)

    forecastSvc := service.NewForecastService(forecastRepo, invRepo, txnRepo, appLogger)
    webhookPublisher := service.NewWebhookPublisher(webhookRepo, appLogger, auditSvc)
    anomalySvc := service.NewAnomalyService(anomalyRepo, invRepo, txnRepo, appLogger, webhookPublisher, auditSvc)

    authHandler := handler.NewAuthHandler(userRepo, jwtManager, cfg, handler.AuthDeps{
        AuditSvc:    auditSvc,
        RedisClient: redisClient,
    })
    healthHandler := handler.NewHealthHandler(db, redisClient)
    forecastHandler := handler.NewForecastHandler(forecastRepo, forecastSvc, redisClient)
    inventoryHandler := handler.NewInventoryHandler(invRepo, txnRepo, db, anomalySvc, appLogger, handler.InventoryDeps{
        AuditSvc: auditSvc,
    })
    anomalyHandler := handler.NewAnomalyHandler(anomalyRepo)
    apiKeyHandler := handler.NewAPIKeyHandler(apiKeyRepo)
    webhookHandler := handler.NewWebhookHandler(webhookRepo)

	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)

	// Health
	mux.HandleFunc("GET /api/v1/health", healthHandler.Health)
	mux.HandleFunc("GET /api/v1/readiness", healthHandler.Readiness)
	mux.HandleFunc("GET /api/v1/liveness", healthHandler.Liveness)

	// Inventory
	mux.HandleFunc("GET /api/v1/inventory", inventoryHandler.QueryInventory)
	mux.HandleFunc("POST /api/v1/transactions", inventoryHandler.CreateTransaction)
	mux.HandleFunc("GET /api/v1/transactions", inventoryHandler.QueryTransactions)
	mux.HandleFunc("GET /api/v1/inventory/{warehouse_id}", inventoryHandler.GetWarehouseInventory)
	mux.HandleFunc("POST /api/v1/inventory/import", inventoryHandler.ImportInventory)

	// Forecast
	mux.HandleFunc("GET /api/v1/forecasts", forecastHandler.QueryForecasts)
	mux.HandleFunc("POST /api/v1/forecasts/generate", forecastHandler.GenerateForecasts)
	mux.HandleFunc("GET /api/v1/forecasts/accuracy", forecastHandler.GetAccuracy)

	// Anomaly
	mux.HandleFunc("GET /api/v1/anomalies", anomalyHandler.ListAnomalies)
	mux.HandleFunc("GET /api/v1/anomalies/{id}", anomalyHandler.GetAnomaly)
	mux.HandleFunc("PUT /api/v1/anomalies/{id}", anomalyHandler.UpdateAnomalyStatus)
	mux.HandleFunc("GET /api/v1/anomalies/{id}/history", anomalyHandler.GetAnomalyHistory)

	// API Keys (admin-only)
	adminOnly := middleware.RBACMiddleware("admin")
	mux.Handle("POST /api/v1/api-keys", adminOnly(http.HandlerFunc(apiKeyHandler.CreateAPIKey)))
	mux.Handle("GET /api/v1/api-keys", adminOnly(http.HandlerFunc(apiKeyHandler.ListAPIKeys)))
	mux.Handle("DELETE /api/v1/api-keys/{id}", adminOnly(http.HandlerFunc(apiKeyHandler.RevokeAPIKey)))
	mux.Handle("POST /api/v1/api-keys/{id}/rotate", adminOnly(http.HandlerFunc(apiKeyHandler.RotateAPIKey)))

	// Webhooks
	mux.HandleFunc("POST /api/v1/webhooks/subscribe", webhookHandler.CreateSubscription)
	mux.HandleFunc("GET /api/v1/webhooks/subscriptions", webhookHandler.ListSubscriptions)
	mux.HandleFunc("DELETE /api/v1/webhooks/subscriptions/{id}", webhookHandler.DeleteSubscription)
	mux.HandleFunc("GET /api/v1/webhooks/deliveries", webhookHandler.ListDeliveries)

    rateLimiter := middleware.NewRateLimitMiddleware(redisClient, 1000)
    handlerWithMiddleware := middleware.CORSMiddleware(
        middleware.SecurityHeaders(
            rateLimiter.Apply(
                middleware.AuthMiddleware(jwtManager, redisClient)(mux),
            ),
        ),
    )

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      handlerWithMiddleware,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		appLogger.Info("HTTP server listening", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server error", "error", err.Error())
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutdown signal received, shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err.Error())
	}
	appLogger.Info("Server shutdown complete")
}