package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/internal/hello", func(c *gin.Context) {
		slog.Info("backend request")
		c.JSON(http.StatusOK, gin.H{"service": "backend", "msg": "ok"})
	})
	slog.Info("backend listening", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("run", "err", err)
		os.Exit(1)
	}
}
