package config

import (
	"log"
	"os"
	"strconv"

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

	AuthServiceURL  string
	MediaServiceURL string

	ConsulHost string
	ConsulPort string

	Environment string
	LogLevel    string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "news-service"),
		GRPCPort:    getEnv("GRPC_PORT", "8082"),
		HTTPPort:    getEnv("HTTP_PORT", "8092"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "news_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		AuthServiceURL:  getEnv("AUTH_SERVICE_URL", "localhost:8081"),
		MediaServiceURL: getEnv("MEDIA_SERVICE_URL", "localhost:8084"),

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

func (c *Config) Validate() error {
	if c.DBName == "" {
		log.Println("WARNING: DB_NAME is not set")
	}
	return nil
}
