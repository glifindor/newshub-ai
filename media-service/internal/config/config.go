package config

import (
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

	// MinIO S3 Storage
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool

	// Upload limits
	MaxUploadSize int64 // bytes
	AllowedTypes  []string

	ConsulHost string
	ConsulPort string

	Environment string
	LogLevel    string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "media-service"),
		GRPCPort:    getEnv("GRPC_PORT", "8084"),
		HTTPPort:    getEnv("HTTP_PORT", "8094"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "media_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "media"),
		MinIOUseSSL:    getEnvAsBool("MINIO_USE_SSL", false),

		MaxUploadSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 10*1024*1024), // 10MB
		AllowedTypes: []string{
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"video/mp4", "video/webm",
			"application/pdf",
		},

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

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (c *Config) Validate() error {
	return nil
}
