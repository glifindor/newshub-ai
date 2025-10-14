package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger обертка над zap.Logger
type Logger struct {
	*zap.Logger
}

// New создает новый логгер
func New(environment string) (*Logger, error) {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger}, nil
}

// Sync синхронизирует буферы
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// GetLogger возвращает базовый zap.Logger
func (l *Logger) GetLogger() *zap.Logger {
	return l.Logger
}
