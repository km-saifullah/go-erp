package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/km-saifullah/go-erp/internal/config"
	"github.com/km-saifullah/go-erp/internal/middleware"
)

func New(cfg config.Config) *http.Server {
	router := gin.New()

	// Basic Gin middleware.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	// Temporary health endpoint.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "ERP API is healthy",
		})
	})

	return &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,

		// Prevent clients from keeping connections open forever.
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,

		// Maximum time allowed to read request headers.
		ReadHeaderTimeout: 5 * time.Second,
	}
}
