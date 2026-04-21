package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var ready atomic.Bool

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	setupLogger()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/readyz", func(c *gin.Context) {
		if ready.Load() {
			c.String(http.StatusOK, "ok")
			return
		}
		c.Status(http.StatusServiceUnavailable)
	})

	r.GET("/api/v1/hello", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			name = "world"
		}
		slog.Info("request", "path", c.Request.URL.Path, "name", name)
		c.JSON(http.StatusOK, gin.H{"message": "hello, " + name})
	})

	// 第 6 章：设置 BACKEND_BASE_URL，例如 http://backend.demo.svc.cluster.local:8080
	if base := strings.TrimRight(os.Getenv("BACKEND_BASE_URL"), "/"); base != "" {
		r.GET("/api/v1/chain", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/internal/hello", nil)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			c.JSON(http.StatusOK, gin.H{"via": "api", "backend_status": resp.StatusCode, "backend_body": string(body)})
		})
	}

	// 用于练习 1.5：优雅退出时仍在途的请求
	r.GET("/slow", func(c *gin.Context) {
		slog.Info("slow request start")
		time.Sleep(8 * time.Second)
		c.String(http.StatusOK, "done")
	})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		time.Sleep(3 * time.Second)
		ready.Store(true)
		slog.Info("readiness enabled")
	}()

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown", "err", err)
	}
	slog.Info("server stopped")
}

func setupLogger() {
	level := slog.LevelInfo
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(h))
	// 练习 3.3：从环境变量读取密钥，勿打印到日志
	if k := os.Getenv("API_KEY"); k != "" {
		_ = len(k) // 仅表示已注入；业务里按需使用
		slog.Debug("api key configured")
	}
}
