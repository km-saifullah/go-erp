package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/km-saifullah/go-erp/internal/cache"
	"github.com/km-saifullah/go-erp/internal/config"
	"github.com/km-saifullah/go-erp/internal/database"
	"github.com/km-saifullah/go-erp/internal/server"
)

func main() {
	// --------------------------------------------------
	// Configuration
	// --------------------------------------------------

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	log.Printf(
		"configuration loaded successfully: env=%s port=%s",
		cfg.AppEnv,
		cfg.AppPort,
	)

	// --------------------------------------------------
	// MySQL
	// --------------------------------------------------

	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	log.Println("mysql connected successfully")

	// --------------------------------------------------
	// Redis
	// --------------------------------------------------

	redis, err := cache.NewRedis(cfg)
	if err != nil {
		db.Close()

		log.Fatalf("redis connection failed: %v", err)
	}

	log.Println("redis connected successfully")

	// --------------------------------------------------
	// HTTP Server
	// --------------------------------------------------

	httpServer := server.New(cfg)

	// --------------------------------------------------
	// Start HTTP Server
	// --------------------------------------------------

	go func() {
		log.Printf("ERP API listening on :%s", cfg.AppPort)

		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			log.Printf("HTTP server error: %v", err)
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

	log.Println("shutdown signal received")

	// --------------------------------------------------
	// Graceful Shutdown
	// --------------------------------------------------

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	log.Println("stopping HTTP server...")

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("HTTP server stopped")

	// --------------------------------------------------
	// Close Redis
	// --------------------------------------------------

	log.Println("closing redis connection...")

	if err := redis.Close(); err != nil {
		log.Printf("redis shutdown error: %v", err)
	} else {
		log.Println("redis connection closed")
	}

	// --------------------------------------------------
	// Close MySQL
	// --------------------------------------------------

	log.Println("closing mysql connection...")

	if err := db.Close(); err != nil {
		log.Printf("mysql shutdown error: %v", err)
	} else {
		log.Println("mysql connection closed")
	}

	log.Println("ERP API shutdown completed")
}
