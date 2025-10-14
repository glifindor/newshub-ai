package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/database"
	"auth-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация логгера
	if err := logger.InitLogger(cfg.Environment); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Auth Service",
		zap.String("version", "1.0.0"),
		zap.String("environment", cfg.Environment),
	)

	// Подключение к PostgreSQL через GORM
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Автомиграция
	if err := database.AutoMigrate(db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	// Проверка подключения к Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Successfully connected to Redis")

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(redisClient)

	// Инициализация сервисов
	tokenService := service.NewTokenService(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)
	authService := service.NewAuthService(userRepo, sessionRepo, tokenService)
	userService := service.NewUserService(userRepo)

	// HTTP Server (Gin)
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	httpHandler := handler.NewHTTPHandler(authService, userService)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// Public routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", httpHandler.Register)
		v1.POST("/login", httpHandler.Login)
		v1.POST("/refresh-token", httpHandler.RefreshToken)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.POST("/logout", httpHandler.Logout)
		protected.GET("/profile", httpHandler.GetProfile)
		protected.PUT("/profile", httpHandler.UpdateProfile)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "auth-service",
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

	logger.Info("Shutting down Auth Service...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced to shutdown", zap.Error(err))
	}

	logger.Info("Auth Service stopped gracefully")
}
