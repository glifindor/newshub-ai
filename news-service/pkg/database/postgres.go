package database

import (
	"fmt"
	"time"

	"news-service/internal/config"
	"news-service/internal/models"
	"news-service/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewPostgresDB создает новое подключение к PostgreSQL через GORM
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	var gormLogLevel gormlogger.LogLevel
	if cfg.Environment == "production" {
		gormLogLevel = gormlogger.Error
	} else {
		gormLogLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database")
	return db, nil
}

// AutoMigrate выполняет автомиграцию моделей
func AutoMigrate(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	err := db.AutoMigrate(
		&models.Category{},
		&models.Tag{},
		&models.News{},
	)

	if err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
