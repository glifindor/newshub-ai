package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"seo-service/internal/models"
	"seo-service/internal/repository"
	"seo-service/pkg/generator"
	"seo-service/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SEOService интерфейс для работы с SEO
type SEOService interface {
	// CRUD операции
	CreateOrUpdate(ctx context.Context, req *models.CreateSEORequest) (*models.SEOResponse, error)
	GetBySlug(ctx context.Context, slug string) (*models.SEOResponse, error)
	GetByNewsID(ctx context.Context, newsID uuid.UUID) (*models.SEOResponse, error)
	Delete(ctx context.Context, newsID uuid.UUID) error

	// Генерация SEO из новости
	GenerateFromNews(ctx context.Context, webhook *models.NewsPublishedWebhook) (*models.SEOResponse, error)

	// Пакетные операции
	GetAllIndexable(ctx context.Context) ([]*models.SEOResponse, error)
}

type seoService struct {
	repo              repository.SEORepository
	structuredDataGen *generator.StructuredDataGenerator
	logger            *logger.Logger
	baseURL           string
}

// NewSEOService создает новый SEO сервис
func NewSEOService(
	repo repository.SEORepository,
	structuredDataGen *generator.StructuredDataGenerator,
	logger *logger.Logger,
	baseURL string,
) SEOService {
	return &seoService{
		repo:              repo,
		structuredDataGen: structuredDataGen,
		logger:            logger,
		baseURL:           baseURL,
	}
}

// CreateOrUpdate создает или обновляет SEO метаданные
func (s *seoService) CreateOrUpdate(ctx context.Context, req *models.CreateSEORequest) (*models.SEOResponse, error) {
	// Проверяем, существует ли запись
	existing, _ := s.repo.GetByNewsID(ctx, req.NewsID)

	var meta *models.SEOMeta

	if existing != nil {
		// Обновление существующей записи
		meta = existing
		s.updateMetaFromRequest(meta, req)

		if err := s.repo.Update(ctx, meta); err != nil {
			s.logger.Error("Failed to update SEO meta", zap.Error(err), zap.String("news_id", req.NewsID.String()))
			return nil, fmt.Errorf("failed to update SEO meta: %w", err)
		}

		s.logger.Info("SEO meta updated", zap.String("slug", meta.Slug), zap.String("news_id", req.NewsID.String()))
	} else {
		// Создание новой записи
		meta = &models.SEOMeta{
			ID:        uuid.New(),
			NewsID:    req.NewsID,
			Slug:      req.Slug,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		s.updateMetaFromRequest(meta, req)

		if err := s.repo.Create(ctx, meta); err != nil {
			s.logger.Error("Failed to create SEO meta", zap.Error(err), zap.String("news_id", req.NewsID.String()))
			return nil, fmt.Errorf("failed to create SEO meta: %w", err)
		}

		s.logger.Info("SEO meta created", zap.String("slug", meta.Slug), zap.String("news_id", req.NewsID.String()))
	}

	return s.metaToResponse(meta), nil
}

// GetBySlug получает SEO метаданные по slug
func (s *seoService) GetBySlug(ctx context.Context, slug string) (*models.SEOResponse, error) {
	meta, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return s.metaToResponse(meta), nil
}

// GetByNewsID получает SEO метаданные по ID новости
func (s *seoService) GetByNewsID(ctx context.Context, newsID uuid.UUID) (*models.SEOResponse, error) {
	meta, err := s.repo.GetByNewsID(ctx, newsID)
	if err != nil {
		return nil, err
	}

	return s.metaToResponse(meta), nil
}

// Delete удаляет SEO метаданные
func (s *seoService) Delete(ctx context.Context, newsID uuid.UUID) error {
	if err := s.repo.Delete(ctx, newsID); err != nil {
		s.logger.Error("Failed to delete SEO meta", zap.Error(err), zap.String("news_id", newsID.String()))
		return err
	}

	s.logger.Info("SEO meta deleted", zap.String("news_id", newsID.String()))
	return nil
}

// GenerateFromNews автоматически генерирует SEO из новости (webhook)
func (s *seoService) GenerateFromNews(ctx context.Context, webhook *models.NewsPublishedWebhook) (*models.SEOResponse, error) {
	// Генерируем SEO title (оптимизация для поисковиков)
	seoTitle := s.optimizeTitle(webhook.Title)

	// Генерируем SEO description
	seoDescription := s.generateDescription(webhook.Summary, webhook.Content)

	// Извлекаем ключевые слова
	keywords := s.extractKeywords(webhook.Title, webhook.Summary, webhook.Content)

	// Генерируем Schema.org JSON-LD
	schemaData, err := s.structuredDataGen.GenerateNewsArticle(
		&models.SEOMeta{
			Title:       seoTitle,
			Description: seoDescription,
			Slug:        webhook.Slug,
			OGImage:     webhook.ImageURL,
		},
		webhook.AuthorName,
		webhook.PublishedAt,
		time.Now(),
	)
	if err != nil {
		s.logger.Warn("Failed to generate schema data", zap.Error(err))
		schemaData = nil
	}

	// Создаем request
	req := &models.CreateSEORequest{
		NewsID:       webhook.NewsID,
		Slug:         webhook.Slug,
		Title:        seoTitle,
		Description:  seoDescription,
		Keywords:     keywords,
		CanonicalURL: fmt.Sprintf("%s/news/%s", s.baseURL, webhook.Slug),

		// Open Graph для VK, Telegram, OK
		OGTitle:       seoTitle,
		OGDescription: seoDescription,
		OGImage:       webhook.ImageURL,
		OGType:        "article",
		OGLocale:      "ru_RU",
		OGSiteName:    "Новостной портал",

		// Robots
		RobotsIndex:  true,
		RobotsFollow: true,

		// Sitemap
		SitemapChangeFreq: "daily",
		SitemapPriority:   0.8,

		// Schema.org
		SchemaData: schemaData,
	}

	return s.CreateOrUpdate(ctx, req)
}

// GetAllIndexable получает все индексируемые SEO метаданные
func (s *seoService) GetAllIndexable(ctx context.Context) ([]*models.SEOResponse, error) {
	metaList, err := s.repo.GetAllIndexable(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.SEOResponse, 0, len(metaList))
	for _, meta := range metaList {
		responses = append(responses, s.metaToResponse(meta))
	}

	return responses, nil
}

// updateMetaFromRequest обновляет SEOMeta из request
func (s *seoService) updateMetaFromRequest(meta *models.SEOMeta, req *models.CreateSEORequest) {
	meta.Slug = req.Slug
	meta.Title = req.Title
	meta.Description = req.Description
	meta.Keywords = req.Keywords
	meta.CanonicalURL = req.CanonicalURL

	// Open Graph
	meta.OGTitle = req.OGTitle
	meta.OGDescription = req.OGDescription
	meta.OGImage = req.OGImage
	meta.OGType = req.OGType
	meta.OGLocale = req.OGLocale
	meta.OGSiteName = req.OGSiteName

	// Robots
	meta.RobotsIndex = req.RobotsIndex
	meta.RobotsFollow = req.RobotsFollow

	// Sitemap
	meta.SitemapChangeFreq = req.SitemapChangeFreq
	meta.SitemapPriority = req.SitemapPriority

	// Schema.org
	if req.SchemaData != nil {
		meta.SchemaData = req.SchemaData
	}

	meta.UpdatedAt = time.Now()
}

// metaToResponse конвертирует SEOMeta в response
func (s *seoService) metaToResponse(meta *models.SEOMeta) *models.SEOResponse {
	return &models.SEOResponse{
		ID:                meta.ID,
		NewsID:            meta.NewsID,
		Slug:              meta.Slug,
		Title:             meta.Title,
		Description:       meta.Description,
		Keywords:          meta.Keywords,
		CanonicalURL:      meta.CanonicalURL,
		OGTitle:           meta.OGTitle,
		OGDescription:     meta.OGDescription,
		OGImage:           meta.OGImage,
		OGType:            meta.OGType,
		OGLocale:          meta.OGLocale,
		OGSiteName:        meta.OGSiteName,
		RobotsIndex:       meta.RobotsIndex,
		RobotsFollow:      meta.RobotsFollow,
		SitemapChangeFreq: meta.SitemapChangeFreq,
		SitemapPriority:   meta.SitemapPriority,
		SchemaData:        meta.SchemaData,
		CreatedAt:         meta.CreatedAt,
		UpdatedAt:         meta.UpdatedAt,
	}
}

// optimizeTitle оптимизирует заголовок для поисковиков (макс 70 символов)
func (s *seoService) optimizeTitle(title string) string {
	// Удаляем лишние пробелы
	title = strings.TrimSpace(title)

	// Если заголовок слишком длинный, обрезаем
	if len(title) > 70 {
		// Обрезаем до последнего пробела перед 67 символом (оставляем место для "...")
		if idx := strings.LastIndex(title[:67], " "); idx > 0 {
			title = title[:idx] + "..."
		} else {
			title = title[:67] + "..."
		}
	}

	return title
}

// generateDescription генерирует SEO description (макс 160 символов)
func (s *seoService) generateDescription(summary, content string) string {
	var description string

	// Используем summary если есть, иначе content
	if summary != "" {
		description = summary
	} else if content != "" {
		description = content
	} else {
		return ""
	}

	// Удаляем HTML теги и лишние пробелы
	description = s.stripHTML(description)
	description = strings.TrimSpace(description)

	// Обрезаем до 160 символов
	if len(description) > 160 {
		if idx := strings.LastIndex(description[:157], " "); idx > 0 {
			description = description[:idx] + "..."
		} else {
			description = description[:157] + "..."
		}
	}

	return description
}

// extractKeywords извлекает ключевые слова из текста
func (s *seoService) extractKeywords(title, summary, content string) string {
	// Объединяем все тексты
	text := strings.ToLower(title + " " + summary + " " + content)

	// Удаляем HTML и спецсимволы
	text = s.stripHTML(text)

	// Разбиваем на слова
	words := strings.Fields(text)

	// Частотный анализ (упрощенный)
	wordCount := make(map[string]int)
	stopWords := s.getStopWords()

	for _, word := range words {
		// Очищаем от знаков препинания
		word = strings.Trim(word, ".,!?;:\"'()[]{}«»")

		// Пропускаем короткие слова и стоп-слова
		if len(word) < 3 || stopWords[word] {
			continue
		}

		wordCount[word]++
	}

	// Сортируем по частоте и берем топ-10
	type wordFreq struct {
		word  string
		count int
	}

	freqs := make([]wordFreq, 0, len(wordCount))
	for word, count := range wordCount {
		freqs = append(freqs, wordFreq{word, count})
	}

	// Простая сортировка пузырьком (для небольших данных достаточно)
	for i := 0; i < len(freqs); i++ {
		for j := i + 1; j < len(freqs); j++ {
			if freqs[j].count > freqs[i].count {
				freqs[i], freqs[j] = freqs[j], freqs[i]
			}
		}
	}

	// Берем топ-10
	keywords := make([]string, 0, 10)
	for i := 0; i < len(freqs) && i < 10; i++ {
		keywords = append(keywords, freqs[i].word)
	}

	return strings.Join(keywords, ", ")
}

// stripHTML удаляет HTML теги из текста
func (s *seoService) stripHTML(text string) string {
	// Простое удаление тегов (для production лучше использовать bluemonday)
	inTag := false
	var result strings.Builder

	for _, char := range text {
		if char == '<' {
			inTag = true
		} else if char == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// getStopWords возвращает список стоп-слов (русские и английские)
func (s *seoService) getStopWords() map[string]bool {
	return map[string]bool{
		// Русские
		"и": true, "в": true, "во": true, "не": true, "что": true, "он": true, "на": true,
		"я": true, "с": true, "со": true, "как": true, "а": true, "то": true, "все": true,
		"она": true, "так": true, "его": true, "но": true, "да": true, "ты": true, "к": true,
		"у": true, "же": true, "вы": true, "за": true, "бы": true, "по": true, "только": true,
		"ее": true, "мне": true, "было": true, "вот": true, "от": true, "меня": true, "еще": true,
		"нет": true, "о": true, "из": true, "ему": true, "теперь": true, "когда": true, "даже": true,
		"ну": true, "вдруг": true, "ли": true, "если": true, "уже": true, "или": true, "ни": true,
		"быть": true, "был": true, "него": true, "до": true, "вас": true, "нибудь": true, "опять": true,
		"уж": true, "вам": true, "ведь": true, "там": true, "потом": true, "себя": true, "ничего": true,

		// Английские
		"the": true, "and": true, "is": true, "in": true, "to": true, "of": true, "a": true,
		"for": true, "on": true, "with": true, "as": true, "at": true, "by": true, "from": true,
		"or": true, "an": true, "be": true, "this": true, "that": true, "are": true, "was": true,
	}
}
