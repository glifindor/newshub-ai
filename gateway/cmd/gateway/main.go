package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway/internal/clients"
	"gateway/internal/config"
	"gateway/internal/middleware"
	"gateway/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	// Redis для rate limiting
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	// Инициализация gRPC клиентов
	grpcClients, err := clients.NewClients(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer grpcClients.Close()

	// Настройка Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Глобальные middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.RateLimiter(redisClient))

	// Настройка маршрутов
	router.SetupRoutes(r, grpcClients, redisClient)

	// HTTP сервер
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		log.Printf("API Gateway started on port %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("API Gateway stopped")
}
