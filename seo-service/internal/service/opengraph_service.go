package service

import (
	"context"
	"fmt"

	"seo-service/internal/models"
	"seo-service/pkg/logger"

	"go.uber.org/zap"
)

// OpenGraphService сервис для генерации Open Graph мета-тегов
// Важно для VK, Telegram, Одноклассники
type OpenGraphService interface {
	GenerateOGTags(ctx context.Context, meta *models.SEOMeta) map[string]string
	ValidateOGImage(imageURL string) error
}

type openGraphService struct {
	logger         *logger.Logger
	defaultOGImage string
}

// NewOpenGraphService создает новый OpenGraph сервис
func NewOpenGraphService(logger *logger.Logger, defaultOGImage string) OpenGraphService {
	return &openGraphService{
		logger:         logger,
		defaultOGImage: defaultOGImage,
	}
}

// GenerateOGTags генерирует Open Graph мета-теги для соцсетей
func (s *openGraphService) GenerateOGTags(ctx context.Context, meta *models.SEOMeta) map[string]string {
	tags := make(map[string]string)

	// Основные OG теги
	tags["og:title"] = meta.OGTitle
	tags["og:description"] = meta.OGDescription
	tags["og:type"] = meta.OGType
	tags["og:locale"] = meta.OGLocale

	// URL
	if meta.CanonicalURL != "" {
		tags["og:url"] = meta.CanonicalURL
	}

	// Site name
	if meta.OGSiteName != "" {
		tags["og:site_name"] = meta.OGSiteName
	}

	// Image
	if meta.OGImage != "" {
		tags["og:image"] = meta.OGImage
		tags["og:image:alt"] = meta.OGTitle
		// Рекомендуемый размер для VK, Telegram, OK: 1200x630
		tags["og:image:width"] = "1200"
		tags["og:image:height"] = "630"
	} else if s.defaultOGImage != "" {
		tags["og:image"] = s.defaultOGImage
	}

	// Для статей добавляем дополнительные теги
	if meta.OGType == "article" {
		tags["article:published_time"] = meta.CreatedAt.Format("2006-01-02T15:04:05-07:00")
		tags["article:modified_time"] = meta.UpdatedAt.Format("2006-01-02T15:04:05-07:00")
	}

	s.logger.Debug("Generated OG tags", zap.Int("tags_count", len(tags)), zap.String("slug", meta.Slug))

	return tags
}

// ValidateOGImage проверяет валидность URL изображения
func (s *openGraphService) ValidateOGImage(imageURL string) error {
	if imageURL == "" {
		return nil // Пустой URL допустим
	}

	// Проверяем, что URL начинается с http:// или https://
	if len(imageURL) < 7 || (imageURL[:7] != "http://" && imageURL[:8] != "https://") {
		return fmt.Errorf("invalid image URL format: must start with http:// or https://")
	}

	// Проверяем расширение файла
	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
	hasValidExtension := false

	for _, ext := range validExtensions {
		if len(imageURL) >= len(ext) && imageURL[len(imageURL)-len(ext):] == ext {
			hasValidExtension = true
			break
		}
	}

	if !hasValidExtension {
		s.logger.Warn("Image URL has unsupported extension", zap.String("url", imageURL))
		// Не возвращаем ошибку, просто предупреждение
	}

	return nil
}

// GenerateVKTags генерирует специфичные теги для ВКонтакте
func (s *openGraphService) GenerateVKTags(meta *models.SEOMeta) map[string]string {
	tags := s.GenerateOGTags(context.Background(), meta)

	// VK использует стандартные OG теги, но можно добавить специфичные
	// vk:image для лучшей оптимизации
	if meta.OGImage != "" {
		tags["vk:image"] = meta.OGImage
	}

	return tags
}

// GenerateTelegramTags генерирует теги для Telegram Instant View
func (s *openGraphService) GenerateTelegramTags(meta *models.SEOMeta) map[string]string {
	tags := s.GenerateOGTags(context.Background(), meta)

	// Telegram использует OG теги для превью
	// Можно добавить специфичные meta для Instant View
	tags["telegram:channel"] = "@yournewschannel" // Если есть канал

	return tags
}
