package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string
	GRPCPort    string
	HTTPPort    string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	ConsulHost string
	ConsulPort string

	Environment string
	LogLevel    string
}

func Load() *Config {
	// Загрузка .env файла (игнорируем ошибку, если его нет)
	_ = godotenv.Load()

	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "auth-service"),
		GRPCPort:    getEnv("GRPC_PORT", "8081"),
		HTTPPort:    getEnv("HTTP_PORT", "8091"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "auth_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTAccessExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 168*time.Hour),

		ConsulHost: getEnv("CONSUL_HOST", "localhost"),
		ConsulPort: getEnv("CONSUL_PORT", "8500"),

		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (c *Config) Validate() error {
	if c.JWTSecret == "your-secret-key-change-in-production" {
		log.Println("WARNING: Using default JWT secret. Change it in production!")
	}
	return nil
}
