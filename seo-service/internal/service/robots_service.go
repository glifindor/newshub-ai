package service

import (
	"context"
	"time"

	"seo-service/internal/models"
	"seo-service/pkg/generator"
	"seo-service/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RobotsService сервис для генерации и кеширования robots.txt
type RobotsService interface {
	GenerateRobotsTxt(ctx context.Context) (string, error)
	GenerateCustomRobotsTxt(ctx context.Context, config *models.RobotsConfig) (string, error)
	InvalidateCache(ctx context.Context) error
}

type robotsService struct {
	generator *generator.RobotsGenerator
	redis     *redis.Client
	logger    *logger.Logger
	cacheTTL  time.Duration
}

// NewRobotsService создает новый Robots сервис
func NewRobotsService(
	generator *generator.RobotsGenerator,
	redis *redis.Client,
	logger *logger.Logger,
	cacheTTL time.Duration,
) RobotsService {
	return &robotsService{
		generator: generator,
		redis:     redis,
		logger:    logger,
		cacheTTL:  cacheTTL,
	}
}

const robotsCacheKey = "robots:txt"

// GenerateRobotsTxt генерирует robots.txt с настройками по умолчанию
func (s *robotsService) GenerateRobotsTxt(ctx context.Context) (string, error) {
	// Проверяем кеш
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, robotsCacheKey).Result()
		if err == nil && cached != "" {
			s.logger.Debug("Robots.txt loaded from cache")
			return cached, nil
		}
	}

	// Генерируем robots.txt
	robotsTxt := s.generator.GenerateDefault()

	// Сохраняем в кеш
	if s.redis != nil {
		if err := s.redis.Set(ctx, robotsCacheKey, robotsTxt, s.cacheTTL).Err(); err != nil {
			s.logger.Warn("Failed to cache robots.txt", zap.Error(err))
		} else {
			s.logger.Debug("Robots.txt cached", zap.Duration("ttl", s.cacheTTL))
		}
	}

	s.logger.Info("Robots.txt generated (default)")

	return robotsTxt, nil
}

// GenerateCustomRobotsTxt генерирует robots.txt с кастомной конфигурацией
func (s *robotsService) GenerateCustomRobotsTxt(ctx context.Context, config *models.RobotsConfig) (string, error) {
	// Генерируем с кастомной конфигурацией
	robotsTxt := s.generator.Generate(config)

	// Сохраняем в кеш
	if s.redis != nil {
		if err := s.redis.Set(ctx, robotsCacheKey, robotsTxt, s.cacheTTL).Err(); err != nil {
			s.logger.Warn("Failed to cache custom robots.txt", zap.Error(err))
		}
	}

	s.logger.Info("Robots.txt generated (custom)")

	return robotsTxt, nil
}

// InvalidateCache сбрасывает кеш robots.txt
func (s *robotsService) InvalidateCache(ctx context.Context) error {
	if s.redis == nil {
		return nil
	}

	if err := s.redis.Del(ctx, robotsCacheKey).Err(); err != nil {
		s.logger.Warn("Failed to invalidate robots cache", zap.Error(err))
		return err
	}

	s.logger.Info("Robots.txt cache invalidated")
	return nil
}
