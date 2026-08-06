package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/km-saifullah/go-erp/internal/config"
)

func New(cfg config.Config) *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "ERP API is healthy",
			"time":    time.Now().UTC(),
		})
	})

	return router
}
