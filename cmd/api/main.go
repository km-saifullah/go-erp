package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/km-saifullah/go-erp/internal/cache"
	"github.com/km-saifullah/go-erp/internal/config"
	"github.com/km-saifullah/go-erp/internal/database"
	"github.com/km-saifullah/go-erp/internal/logger"
	"github.com/km-saifullah/go-erp/internal/server"
)

func main() {
	// --------------------------------------------------
	// Configuration
	// --------------------------------------------------

	cfg, err := config.Load()
	if err != nil {
		slog.Error(
			"configuration error",
			"error", err,
		)

		os.Exit(1)
	}

	log := logger.New(cfg.AppEnv)

	log.Info(
		"configuration loaded successfully",
		"environment", cfg.AppEnv,
		"port", cfg.AppPort,
	)

	// --------------------------------------------------
	// MySQL
	// --------------------------------------------------

	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Error(
			"database connection failed",
			"error", err,
		)

		os.Exit(1)
	}

	log.Info("mysql connected successfully")

	// --------------------------------------------------
	// Redis
	// --------------------------------------------------

	redis, err := cache.NewRedis(cfg)
	if err != nil {
		_ = db.Close()

		log.Error(
			"redis connection failed",
			"error", err,
		)

		os.Exit(1)
	}

	log.Info("redis connected successfully")

	// --------------------------------------------------
	// HTTP Server
	// --------------------------------------------------

	httpServer := server.New(cfg)

	// --------------------------------------------------
	// Start HTTP Server
	// --------------------------------------------------

	go func() {
		log.Info(
			"ERP API server started",
			"port", cfg.AppPort,
		)

		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			log.Error(
				"HTTP server error",
				"error", err,
			)
		}
	}()

	// --------------------------------------------------
	// Wait for Shutdown Signal
	// --------------------------------------------------

	shutdownSignal := make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignal,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-shutdownSignal

	log.Info("shutdown signal received")

	// --------------------------------------------------
	// Graceful Shutdown
	// --------------------------------------------------

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	log.Info("stopping HTTP server")

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error(
			"HTTP server shutdown error",
			"error", err,
		)
	}

	log.Info("HTTP server stopped")

	// --------------------------------------------------
	// Close Redis
	// --------------------------------------------------

	log.Info("closing redis connection")

	if err := redis.Close(); err != nil {
		log.Error(
			"redis shutdown error",
			"error", err,
		)
	} else {
		log.Info("redis connection closed")
	}

	// --------------------------------------------------
	// Close MySQL
	// --------------------------------------------------

	log.Info("closing mysql connection")

	if err := db.Close(); err != nil {
		log.Error(
			"mysql shutdown error",
			"error", err,
		)
	} else {
		log.Info("mysql connection closed")
	}

	log.Info("ERP API shutdown completed")
}
