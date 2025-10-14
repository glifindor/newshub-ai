package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// InitLogger инициализирует глобальный логгер
func InitLogger(environment string) error {
	var err error
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	log, err = config.Build()
	return err
}

// Info логирует информационное сообщение
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Error логирует ошибку
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Debug логирует отладочное сообщение
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Warn логирует предупреждение
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Fatal логирует фатальную ошибку и завершает программу
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Sync сбрасывает буферы логгера
func Sync() error {
	return log.Sync()
}

// GetLogger возвращает экземпляр логгера
func GetLogger() *zap.Logger {
	return log
}
