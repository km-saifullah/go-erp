package main

import (
	"log"

	"github.com/km-saifullah/go-erp/internal/cache"
	"github.com/km-saifullah/go-erp/internal/config"
	"github.com/km-saifullah/go-erp/internal/database"
	"github.com/km-saifullah/go-erp/internal/server"
)

func main() {
	cfg := config.Load()

	// MySQL
	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	defer db.Close()

	log.Println("MySQL connected successfully")

	// Redis
	redis, err := cache.NewRedis(cfg)
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	defer redis.Close()

	log.Println("Redis connected successfully")

	// HTTP server
	router := server.New(cfg)

	log.Printf(
		"ERP API running on port %s [%s]",
		cfg.AppPort,
		cfg.AppEnv,
	)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
