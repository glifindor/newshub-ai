package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"seo-service/internal/config"
	"seo-service/internal/handler"
	"seo-service/internal/repository"
	"seo-service/internal/service"
	"seo-service/pkg/database"
	"seo-service/pkg/generator"
	"seo-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

	log.Info("Starting SEO Service",
		zap.String("environment", cfg.Environment),
		zap.String("port", cfg.ServerPort),
	)

	// Подключаемся к PostgreSQL
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}

	log.Info("Connected to PostgreSQL",
		zap.String("host", cfg.DBHost),
		zap.String("database", cfg.DBName),
	)

	// Подключаемся к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	log.Info("Connected to Redis", zap.String("addr", cfg.RedisAddr))

	// Инициализируем слои приложения
	// Repository
	seoRepo := repository.NewSEORepository(db)

	// Generators
	sitemapGen := generator.NewSitemapGenerator(cfg.BaseURL)
	robotsGen := generator.NewRobotsGenerator(cfg.BaseURL)
	structuredDataGen := generator.NewStructuredDataGenerator(
		cfg.BaseURL,
		cfg.SiteName,
		cfg.BaseURL,
		fmt.Sprintf("%s/static/logo.png", cfg.BaseURL),
	)

	// Services
	seoService := service.NewSEOService(seoRepo, structuredDataGen, log, cfg.BaseURL)
	ogService := service.NewOpenGraphService(log, fmt.Sprintf("%s/static/default-og.jpg", cfg.BaseURL))
	sitemapService := service.NewSitemapService(seoRepo, sitemapGen, redisClient, log, time.Hour)
	robotsService := service.NewRobotsService(robotsGen, redisClient, log, 24*time.Hour)

	// HTTP Handler
	httpHandler := handler.NewHTTPHandler(seoService, ogService, sitemapService, robotsService, log)

	// Настраиваем Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware для CORS (для российских сайтов)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Регистрируем роуты
	httpHandler.RegisterRoutes(router)

	// Запускаем HTTP сервер
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)

	log.Info("Starting HTTP server",
		zap.String("address", serverAddr),
		zap.String("base_url", cfg.BaseURL),
	)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Канал для перехвата сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	log.Info("SEO Service started successfully")

	// Ожидаем сигнал завершения
	<-quit
	log.Info("Shutting down server...")

	// Graceful shutdown с таймаутом 5 секунд
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Закрываем соединения
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Error("Failed to close Redis connection", zap.Error(err))
		}
	}

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Error("Failed to close PostgreSQL connection", zap.Error(err))
		}
	}

	<-shutdownCtx.Done()
	log.Info("Server stopped")
}
