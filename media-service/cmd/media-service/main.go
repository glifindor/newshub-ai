package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"media-service/internal/config"
	"media-service/internal/handler"
	"media-service/internal/repository"
	"media-service/internal/service"
	"media-service/pkg/database"
	"media-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	// Инициализация логгера
	if err := logger.InitLogger(cfg.Environment); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Media Service",
		zap.String("version", "1.0.0"),
		zap.String("environment", cfg.Environment),
	)

	// PostgreSQL через GORM
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Автомиграция
	if err := database.AutoMigrate(db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Repository
	mediaRepo := repository.NewMediaRepository(db)

	// Service с MinIO
	mediaService, err := service.NewMediaService(mediaRepo, cfg)
	if err != nil {
		logger.Fatal("Failed to create media service", zap.Error(err))
	}

	// HTTP Server (Gin)
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Лимит размера запроса
	router.MaxMultipartMemory = cfg.MaxUploadSize

	httpHandler := handler.NewHTTPHandler(mediaService)

	// Public routes
	v1 := router.Group("/api/v1/media")
	{
		v1.GET("/file/:filename", httpHandler.ServeFile)
		v1.GET("/:id", httpHandler.GetByID)
		v1.GET("", httpHandler.List)
	}

	// Protected routes (требуют авторизации)
	protected := router.Group("/api/v1/media")
	// protected.Use(authMiddleware.RequireAuth())
	{
		protected.POST("/upload", httpHandler.Upload)
		protected.PUT("/:id", httpHandler.Update)
		protected.DELETE("/:id", httpHandler.Delete)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "media-service",
		})
	})

	// Запуск HTTP сервера
	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		logger.Info("HTTP server started", zap.String("port", cfg.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Media Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced to shutdown", zap.Error(err))
	}

	logger.Info("Media Service stopped gracefully")
}
