package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config структура конфигурации
type Config struct {
	Environment string
	ServerPort  string
	JWTSecret   string

	// Service URLs
	AuthServiceURL  string
	NewsServiceURL  string
	SEOServiceURL   string
	MediaServiceURL string

	// gRPC Addresses
	AuthServiceGRPC  string
	NewsServiceGRPC  string
	MediaServiceGRPC string

	// CORS
	AllowedOrigins []string

	// Logging
	LogLevel string

	// Prometheus
	PrometheusURL string

	// File Upload
	MaxUploadSize    int64
	AllowedFileTypes []string

	// Pagination
	DefaultPageSize int
	MaxPageSize     int
}

// Load загружает конфигурацию из .env
func Load() (*Config, error) {
	// Загружаем .env файл (игнорируем ошибку если файл не найден)
	_ = godotenv.Load()

	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	defaultPageSize, _ := strconv.Atoi(getEnv("DEFAULT_PAGE_SIZE", "20"))
	maxPageSize, _ := strconv.Atoi(getEnv("MAX_PAGE_SIZE", "100"))

	return &Config{
		Environment:      getEnv("ENVIRONMENT", "development"),
		ServerPort:       getEnv("SERVER_PORT", "8084"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		AuthServiceURL:   getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		NewsServiceURL:   getEnv("NEWS_SERVICE_URL", "http://localhost:8082"),
		SEOServiceURL:    getEnv("SEO_SERVICE_URL", "http://localhost:8093"),
		MediaServiceURL:  getEnv("MEDIA_SERVICE_URL", "http://localhost:8085"),
		AuthServiceGRPC:  getEnv("AUTH_SERVICE_GRPC", "localhost:8081"),
		NewsServiceGRPC:  getEnv("NEWS_SERVICE_GRPC", "localhost:8082"),
		MediaServiceGRPC: getEnv("MEDIA_SERVICE_GRPC", "localhost:8085"),
		AllowedOrigins:   strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:5173"), ","),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		PrometheusURL:    getEnv("PROMETHEUS_URL", "http://localhost:9090"),
		MaxUploadSize:    maxUploadSize,
		AllowedFileTypes: strings.Split(getEnv("ALLOWED_FILE_TYPES", "image/jpeg,image/png,image/webp"), ","),
		DefaultPageSize:  defaultPageSize,
		MaxPageSize:      maxPageSize,
	}, nil
}

// IsProduction проверяет production окружение
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment проверяет development окружение
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
