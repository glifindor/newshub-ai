package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"news-service/internal/config"
	"news-service/internal/handler"
	"news-service/internal/repository"
	"news-service/internal/service"
	"news-service/pkg/database"
	"news-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	// Инициализация логгера
	if err := logger.InitLogger(cfg.Environment); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting News Service",
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

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Successfully connected to Redis")

	// Repositories
	newsRepo := repository.NewNewsRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// SEO Service
	seoService := service.NewSEOService(logger.GetLogger())

	// Services
	newsService := service.NewNewsService(newsRepo, categoryRepo, tagRepo, redisClient, seoService)
	categoryService := service.NewCategoryService(categoryRepo, redisClient)
	tagService := service.NewTagService(tagRepo)

	// HTTP Server (Gin)
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	httpHandler := handler.NewHTTPHandler(newsService, categoryService, tagService)

	// Public routes
	v1 := router.Group("/api/v1")
	{
		// News
		v1.GET("/news", httpHandler.ListNews)
		v1.GET("/news/search", httpHandler.SearchNews) // Full Text Search
		v1.GET("/news/:slug", httpHandler.GetNewsBySlug)
		v1.GET("/news/featured", httpHandler.GetFeaturedNews)
		v1.GET("/news/breaking", httpHandler.GetBreakingNews)

		// Categories
		v1.GET("/categories", httpHandler.ListCategories)
		v1.GET("/categories/tree", httpHandler.GetCategoryTree)
		v1.GET("/categories/:slug", httpHandler.GetCategoryBySlug)

		// Tags
		v1.GET("/tags", httpHandler.ListTags)
		v1.GET("/tags/search", httpHandler.SearchTags)
	}

	// Protected routes (требуют авторизации)
	protected := router.Group("/api/v1")
	// protected.Use(authMiddleware.RequireAuth())
	{
		// News management
		protected.POST("/news", httpHandler.CreateNews)
		protected.PUT("/news/:id", httpHandler.UpdateNews)
		protected.DELETE("/news/:id", httpHandler.DeleteNews)
		protected.POST("/news/:id/publish", httpHandler.PublishNews)

		// Category management
		protected.POST("/categories", httpHandler.CreateCategory)
		protected.PUT("/categories/:id", httpHandler.UpdateCategory)
		protected.DELETE("/categories/:id", httpHandler.DeleteCategory)

		// Tag management
		protected.POST("/tags", httpHandler.CreateTag)
		protected.PUT("/tags/:id", httpHandler.UpdateTag)
		protected.DELETE("/tags/:id", httpHandler.DeleteTag)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "news-service",
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

	// gRPC Server (disabled for now)
	/*
		grpcServer := grpc.NewServer()
		newsHandler := handler.NewNewsHandler(newsService, categoryService, tagService)
		pb.RegisterNewsServiceServer(grpcServer, newsHandler)
		reflection.Register(grpcServer)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			logger.Fatal("Failed to listen on gRPC port", zap.String("port", cfg.GRPCPort), zap.Error(err))
		}

		go func() {
			logger.Info("gRPC server started", zap.String("port", cfg.GRPCPort))
			if err := grpcServer.Serve(listener); err != nil {
				logger.Fatal("Failed to serve gRPC", zap.Error(err))
			}
		}()
	*/

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down News Service...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced to shutdown", zap.Error(err))
	}

	// grpcServer.GracefulStop() // Disabled
	logger.Info("News Service stopped gracefully")
}
