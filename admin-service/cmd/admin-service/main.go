package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"admin-service/internal/client"
	"admin-service/internal/config"
	"admin-service/internal/handler"
	"admin-service/internal/middleware"
	"admin-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализируем логгер
	log, err := logger.New(cfg.Environment)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Admin Service",
		zap.String("environment", cfg.Environment),
		zap.String("port", cfg.ServerPort),
	)

	// Инициализируем клиенты для других сервисов
	authClient := client.NewAuthClient(cfg.AuthServiceURL, log)
	newsClient := client.NewNewsClient(cfg.NewsServiceURL, log)
	seoClient := client.NewSEOClient(cfg.SEOServiceURL, log)

	log.Info("Service clients initialized",
		zap.String("auth_service", cfg.AuthServiceURL),
		zap.String("news_service", cfg.NewsServiceURL),
		zap.String("seo_service", cfg.SEOServiceURL),
	)

	// Инициализируем middleware
	authMiddleware := middleware.NewAuthMiddleware(authClient, log)

	// Инициализируем handlers
	adminHandler := handler.NewAdminHandler(
		authClient,
		newsClient,
		seoClient,
		log,
		cfg.DefaultPageSize,
		cfg.MaxPageSize,
	)

	// Настраиваем Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORS(cfg.AllowedOrigins))

	// Регистрируем роуты
	adminHandler.RegisterRoutes(router, authMiddleware)

	// Статические файлы фронтенда (будут добавлены позже)
	router.Static("/assets", "./frontend/dist/assets")
	router.StaticFile("/", "./frontend/dist/index.html")
	router.NoRoute(func(c *gin.Context) {
		// SPA fallback - все неизвестные роуты отдаем index.html
		c.File("./frontend/dist/index.html")
	})

	// Запускаем HTTP сервер
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)

	log.Info("Starting HTTP server",
		zap.String("address", serverAddr),
		zap.Strings("allowed_origins", cfg.AllowedOrigins),
	)

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	log.Info("Admin Service started successfully")

	// Ожидаем сигнал завершения
	<-quit
	log.Info("Shutting down server...")

	// Graceful shutdown с таймаутом 5 секунд
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	<-shutdownCtx.Done()
	log.Info("Server stopped")
}
