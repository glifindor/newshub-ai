package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Service
	Environment string
	HTTPPort    string
	GRPCPort    string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// SEO Settings
	BaseURL         string
	SitemapCacheTTL int
	RobotsAllow     string
	RobotsDisallow  string

	// News Service Integration
	NewsServiceGRPC string
	NewsServiceHTTP string

	// Security
	WebhookSecret string
}

func Load() *Config {
	// Загрузка .env файла (игнорируем ошибку в production)
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	sitemapTTL, _ := strconv.Atoi(getEnv("SITEMAP_CACHE_TTL", "3600"))

	return &Config{
		// Service
		Environment: getEnv("ENVIRONMENT", "development"),
		HTTPPort:    getEnv("HTTP_PORT", "8093"),
		GRPCPort:    getEnv("GRPC_PORT", "50053"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "newsportal"),
		DBPassword: getEnv("DB_PASSWORD", "SimplePass123"),
		DBName:     getEnv("DB_NAME", "newsportal_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Redis
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		// SEO Settings
		BaseURL:         getEnv("BASE_URL", "https://newsportal.com"),
		SitemapCacheTTL: sitemapTTL,
		RobotsAllow:     getEnv("ROBOTS_ALLOW", "/"),
		RobotsDisallow:  getEnv("ROBOTS_DISALLOW", "/admin/,/api/"),

		// News Service
		NewsServiceGRPC: getEnv("NEWS_SERVICE_GRPC", "localhost:50052"),
		NewsServiceHTTP: getEnv("NEWS_SERVICE_HTTP", "http://localhost:8092"),

		// Security
		WebhookSecret: getEnv("WEBHOOK_SECRET", "change_me_in_production"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Config) GetDSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
