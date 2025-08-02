// @title           Beta Payment API
// @version         1.0
// @description     API service to manage payment using Clean Architecture

// @contact.name   ADF Code
// @contact.url    https://github.com/adf-code

// @host      localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Using token header using the Bearer scheme. Example: "Bearer {token}"

package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/adf-code/beta-payment-api/config"
	_ "github.com/adf-code/beta-payment-api/docs"
	deliveryHttp "github.com/adf-code/beta-payment-api/internal/delivery/http"
	pkgDatabase "github.com/adf-code/beta-payment-api/internal/pkg/database"
	pkgLogger "github.com/adf-code/beta-payment-api/internal/pkg/logger"
	"github.com/adf-code/beta-payment-api/internal/repository"
	"github.com/adf-code/beta-payment-api/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_ = godotenv.Load() // Load .env

	// Load env config
	cfg := config.LoadConfig()

	logger := pkgLogger.InitLoggerWithTelemetry(cfg)
	postgresClient := pkgDatabase.NewPostgresClient(cfg, logger)
	db := postgresClient.InitPostgresDB()

	// Repository and HTTP handler
	paymentRepo := repository.NewPaymentRepo(db)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, db, logger)
	handler := deliveryHttp.SetupHandler(paymentUC, logger)

	// HTTP server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: handler,
	}

	// Run server in goroutine
	go func() {
		logger.Info().Msgf("🟢 Server running on http://localhost:%s", cfg.Port)
		logger.Info().Msgf("📚 Swagger running on http://localhost:%s/swagger/index.html", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msgf("❌ Server failed: %v", err)
		}
	}()

	// Setup signal listener
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msgf("🛑 Gracefully shutting down server...")

	// Graceful shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msgf("❌ Server shutdown failed: %v", err)
	}

	// ✅ Close PostgreSQL DB
	closePostgres(db, logger)

	logger.Info().Msgf("✅ Server shutdown completed.")
}

func closePostgres(db *sql.DB, logger zerolog.Logger) {
	if err := db.Close(); err != nil {
		logger.Info().Msgf("⚠️ Failed to close PostgreSQL connection: %v", err)
	} else {
		logger.Info().Msgf("🔒 PostgreSQL connection closed.")
	}
}
