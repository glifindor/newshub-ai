package service

import (
	"context"
	"fmt"
	"time"

	"seo-service/internal/repository"
	"seo-service/pkg/generator"
	"seo-service/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// SitemapService сервис для генерации и кеширования sitemap.xml
type SitemapService interface {
	GenerateSitemap(ctx context.Context) (string, error)
	GenerateSitemapIndex(ctx context.Context) (string, error)
	InvalidateCache(ctx context.Context) error
}

type sitemapService struct {
	repo      repository.SEORepository
	generator *generator.SitemapGenerator
	redis     *redis.Client
	logger    *logger.Logger
	cacheTTL  time.Duration
}

// NewSitemapService создает новый Sitemap сервис
func NewSitemapService(
	repo repository.SEORepository,
	generator *generator.SitemapGenerator,
	redis *redis.Client,
	logger *logger.Logger,
	cacheTTL time.Duration,
) SitemapService {
	return &sitemapService{
		repo:      repo,
		generator: generator,
		redis:     redis,
		logger:    logger,
		cacheTTL:  cacheTTL,
	}
}

const (
	sitemapCacheKey      = "sitemap:main"
	sitemapIndexCacheKey = "sitemap:index"
)

// GenerateSitemap генерирует sitemap.xml с кешированием
func (s *sitemapService) GenerateSitemap(ctx context.Context) (string, error) {
	// Проверяем кеш
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, sitemapCacheKey).Result()
		if err == nil && cached != "" {
			s.logger.Debug("Sitemap loaded from cache")
			return cached, nil
		}
	}

	// Получаем все индексируемые записи
	metaList, err := s.repo.GetAllIndexable(ctx)
	if err != nil {
		s.logger.Error("Failed to get indexable SEO meta", zap.Error(err))
		return "", fmt.Errorf("failed to get indexable SEO meta: %w", err)
	}

	// Генерируем sitemap
	sitemap, err := s.generator.GenerateURLSet(metaList)
	if err != nil {
		s.logger.Error("Failed to generate sitemap", zap.Error(err))
		return "", fmt.Errorf("failed to generate sitemap: %w", err)
	}

	// Сохраняем в кеш
	if s.redis != nil {
		if err := s.redis.Set(ctx, sitemapCacheKey, sitemap, s.cacheTTL).Err(); err != nil {
			s.logger.Warn("Failed to cache sitemap", zap.Error(err))
		} else {
			s.logger.Debug("Sitemap cached", zap.Duration("ttl", s.cacheTTL))
		}
	}

	s.logger.Info("Sitemap generated", zap.Int("urls_count", len(metaList)))

	return sitemap, nil
}

// GenerateSitemapIndex генерирует sitemap index для больших сайтов
func (s *sitemapService) GenerateSitemapIndex(ctx context.Context) (string, error) {
	// Проверяем кеш
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, sitemapIndexCacheKey).Result()
		if err == nil && cached != "" {
			s.logger.Debug("Sitemap index loaded from cache")
			return cached, nil
		}
	}

	// Список sitemap файлов (для больших сайтов можно разбивать по категориям)
	sitemaps := []string{
		"sitemap-news.xml",
		"sitemap-static.xml",
	}

	// Генерируем index
	index, err := s.generator.GenerateIndex(sitemaps)
	if err != nil {
		s.logger.Error("Failed to generate sitemap index", zap.Error(err))
		return "", fmt.Errorf("failed to generate sitemap index: %w", err)
	}

	// Сохраняем в кеш
	if s.redis != nil {
		if err := s.redis.Set(ctx, sitemapIndexCacheKey, index, s.cacheTTL).Err(); err != nil {
			s.logger.Warn("Failed to cache sitemap index", zap.Error(err))
		}
	}

	s.logger.Info("Sitemap index generated")

	return index, nil
}

// InvalidateCache сбрасывает кеш sitemap
func (s *sitemapService) InvalidateCache(ctx context.Context) error {
	if s.redis == nil {
		return nil
	}

	keys := []string{sitemapCacheKey, sitemapIndexCacheKey}

	for _, key := range keys {
		if err := s.redis.Del(ctx, key).Err(); err != nil {
			s.logger.Warn("Failed to invalidate cache", zap.String("key", key), zap.Error(err))
			return err
		}
	}

	s.logger.Info("Sitemap cache invalidated")
	return nil
}
