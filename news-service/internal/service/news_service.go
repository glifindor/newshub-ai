package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"news-service/internal/models"
	"news-service/internal/repository"
	"news-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type NewsService interface {
	Create(ctx context.Context, authorID uuid.UUID, req *models.NewsRequest) (*models.News, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.News, error)
	GetBySlug(ctx context.Context, slug string) (*models.News, error)
	Update(ctx context.Context, id uuid.UUID, req *models.NewsRequest) (*models.News, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters repository.NewsFilters) ([]*models.News, int64, error)
	Publish(ctx context.Context, id uuid.UUID) error
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	GetFeatured(ctx context.Context, limit int) ([]*models.News, error)
	GetBreaking(ctx context.Context, limit int) ([]*models.News, error)

	// Full Text Search
	Search(ctx context.Context, query string, page, pageSize int) ([]*models.News, int64, error)
}

type newsService struct {
	newsRepo     repository.NewsRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
	redis        *redis.Client
	seoService   SEOService
}

func NewNewsService(
	newsRepo repository.NewsRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
	redis *redis.Client,
	seoService SEOService,
) NewsService {
	return &newsService{
		newsRepo:     newsRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		redis:        redis,
		seoService:   seoService,
	}
}

func (s *newsService) Create(ctx context.Context, authorID uuid.UUID, req *models.NewsRequest) (*models.News, error) {
	// Проверка существования категории
	_, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		logger.Error("Category not found", zap.Error(err), zap.String("category_id", req.CategoryID.String()))
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Генерация slug
	newsSlug := req.Slug
	if newsSlug == "" {
		newsSlug = slug.Make(req.Title)
	}

	news := &models.News{
		Title:           req.Title,
		Slug:            newsSlug,
		Summary:         req.Summary,
		Content:         req.Content,
		FeaturedImage:   req.FeaturedImage,
		AuthorID:        authorID,
		CategoryID:      req.CategoryID,
		Status:          req.Status,
		PublishedAt:     req.PublishedAt,
		IsFeatured:      req.IsFeatured,
		IsBreaking:      req.IsBreaking,
		AllowComments:   req.AllowComments,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		MetaKeywords:    req.MetaKeywords,
	}

	// Если статус не указан, ставим draft
	if news.Status == "" {
		news.Status = models.StatusDraft
	}

	// Автогенерация SEO метатегов если не указаны
	if req.SEOTitle == "" || req.SEODescription == "" || req.SEOKeywords == "" {
		seoTitle, seoDesc, keywords := s.seoService.GenerateMetaTags(ctx, req.Title, req.Content)

		if req.SEOTitle == "" {
			news.SEOTitle = seoTitle
		} else {
			news.SEOTitle = req.SEOTitle
		}

		if req.SEODescription == "" {
			news.SEODescription = seoDesc
		} else {
			news.SEODescription = req.SEODescription
		}

		if req.SEOKeywords == "" {
			news.SEOKeywords = keywords
		} else {
			news.SEOKeywords = req.SEOKeywords
		}

		if req.SEOImage != "" {
			news.SEOImage = req.SEOImage
		} else if req.FeaturedImage != "" {
			news.SEOImage = req.FeaturedImage
		}

		logger.Info("Auto-generated SEO meta tags",
			zap.String("title", news.Title),
			zap.String("seo_title", news.SEOTitle),
		)
	} else {
		news.SEOTitle = req.SEOTitle
		news.SEODescription = req.SEODescription
		news.SEOKeywords = req.SEOKeywords
		news.SEOImage = req.SEOImage
	}

	// Получение тегов
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.GetByIDs(ctx, req.TagIDs)
		if err != nil {
			logger.Error("Failed to get tags", zap.Error(err))
			return nil, err
		}
		news.Tags = make([]models.Tag, len(tags))
		for i, tag := range tags {
			news.Tags[i] = *tag
		}
	}

	if err := s.newsRepo.Create(ctx, news); err != nil {
		logger.Error("Failed to create news", zap.Error(err))
		return nil, err
	}

	logger.Info("News created successfully",
		zap.String("news_id", news.ID.String()),
		zap.String("title", news.Title),
		zap.String("author_id", authorID.String()),
	)

	// Инвалидация кэша
	s.invalidateCache(ctx)

	return news, nil
}

func (s *newsService) GetByID(ctx context.Context, id uuid.UUID) (*models.News, error) {
	// Попытка получить из кэша
	cacheKey := fmt.Sprintf("news:id:%s", id.String())
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var news models.News
		if json.Unmarshal([]byte(cached), &news) == nil {
			return &news, nil
		}
	}

	news, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Сохранение в кэш на 5 минут
	data, _ := json.Marshal(news)
	s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return news, nil
}

func (s *newsService) GetBySlug(ctx context.Context, newsSlug string) (*models.News, error) {
	// Попытка получить из кэша
	cacheKey := fmt.Sprintf("news:slug:%s", newsSlug)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var news models.News
		if json.Unmarshal([]byte(cached), &news) == nil {
			return &news, nil
		}
	}

	news, err := s.newsRepo.GetBySlug(ctx, newsSlug)
	if err != nil {
		return nil, err
	}

	// Сохранение в кэш на 5 минут
	data, _ := json.Marshal(news)
	s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return news, nil
}

func (s *newsService) Update(ctx context.Context, id uuid.UUID, req *models.NewsRequest) (*models.News, error) {
	news, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Обновление полей
	titleChanged := news.Title != req.Title
	contentChanged := news.Content != req.Content

	news.Title = req.Title
	if req.Slug != "" {
		news.Slug = req.Slug
	} else {
		news.Slug = slug.Make(req.Title)
	}
	news.Summary = req.Summary
	news.Content = req.Content
	news.FeaturedImage = req.FeaturedImage
	news.CategoryID = req.CategoryID
	news.Status = req.Status
	news.PublishedAt = req.PublishedAt
	news.IsFeatured = req.IsFeatured
	news.IsBreaking = req.IsBreaking
	news.AllowComments = req.AllowComments
	news.MetaTitle = req.MetaTitle
	news.MetaDescription = req.MetaDescription
	news.MetaKeywords = req.MetaKeywords

	// Регенерация SEO метатегов если изменился title или content
	if (titleChanged || contentChanged) && (req.SEOTitle == "" || req.SEODescription == "" || req.SEOKeywords == "") {
		seoTitle, seoDesc, keywords := s.seoService.GenerateMetaTags(ctx, req.Title, req.Content)

		if req.SEOTitle == "" {
			news.SEOTitle = seoTitle
		} else {
			news.SEOTitle = req.SEOTitle
		}

		if req.SEODescription == "" {
			news.SEODescription = seoDesc
		} else {
			news.SEODescription = req.SEODescription
		}

		if req.SEOKeywords == "" {
			news.SEOKeywords = keywords
		} else {
			news.SEOKeywords = req.SEOKeywords
		}

		logger.Info("Regenerated SEO meta tags after content update",
			zap.String("news_id", id.String()),
		)
	} else if req.SEOTitle != "" || req.SEODescription != "" || req.SEOKeywords != "" || req.SEOImage != "" {
		news.SEOTitle = req.SEOTitle
		news.SEODescription = req.SEODescription
		news.SEOKeywords = req.SEOKeywords
		news.SEOImage = req.SEOImage
	}

	// Обновление тегов
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.GetByIDs(ctx, req.TagIDs)
		if err != nil {
			return nil, err
		}
		news.Tags = make([]models.Tag, len(tags))
		for i, tag := range tags {
			news.Tags[i] = *tag
		}
	}

	if err := s.newsRepo.Update(ctx, news); err != nil {
		logger.Error("Failed to update news", zap.Error(err))
		return nil, err
	}

	logger.Info("News updated successfully", zap.String("news_id", news.ID.String()))

	// Инвалидация кэша
	s.invalidateCache(ctx)
	s.redis.Del(ctx, fmt.Sprintf("news:id:%s", id.String()))
	s.redis.Del(ctx, fmt.Sprintf("news:slug:%s", news.Slug))

	return news, nil
}

func (s *newsService) Delete(ctx context.Context, id uuid.UUID) error {
	news, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.newsRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete news", zap.Error(err))
		return err
	}

	logger.Info("News deleted successfully", zap.String("news_id", id.String()))

	// Инвалидация кэша
	s.invalidateCache(ctx)
	s.redis.Del(ctx, fmt.Sprintf("news:id:%s", id.String()))
	s.redis.Del(ctx, fmt.Sprintf("news:slug:%s", news.Slug))

	return nil
}

func (s *newsService) List(ctx context.Context, filters repository.NewsFilters) ([]*models.News, int64, error) {
	return s.newsRepo.List(ctx, filters)
}

func (s *newsService) Publish(ctx context.Context, id uuid.UUID) error {
	news, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	news.Status = models.StatusPublished
	now := time.Now()
	news.PublishedAt = &now

	if err := s.newsRepo.Update(ctx, news); err != nil {
		logger.Error("Failed to publish news", zap.Error(err))
		return err
	}

	logger.Info("News published successfully", zap.String("news_id", id.String()))

	// Инвалидация кэша
	s.invalidateCache(ctx)

	return nil
}

func (s *newsService) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return s.newsRepo.IncrementViewCount(ctx, id)
}

func (s *newsService) GetFeatured(ctx context.Context, limit int) ([]*models.News, error) {
	cacheKey := fmt.Sprintf("news:featured:%d", limit)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var news []*models.News
		if json.Unmarshal([]byte(cached), &news) == nil {
			return news, nil
		}
	}

	news, err := s.newsRepo.GetFeatured(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Кэш на 10 минут
	data, _ := json.Marshal(news)
	s.redis.Set(ctx, cacheKey, data, 10*time.Minute)

	return news, nil
}

func (s *newsService) GetBreaking(ctx context.Context, limit int) ([]*models.News, error) {
	cacheKey := fmt.Sprintf("news:breaking:%d", limit)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var news []*models.News
		if json.Unmarshal([]byte(cached), &news) == nil {
			return news, nil
		}
	}

	news, err := s.newsRepo.GetBreaking(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Кэш на 5 минут
	data, _ := json.Marshal(news)
	s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return news, nil
}

func (s *newsService) invalidateCache(ctx context.Context) {
	// Удаление кэша списков
	pattern := "news:*"
	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		s.redis.Del(ctx, iter.Val())
	}
}

// Search выполняет полнотекстовый поиск новостей
func (s *newsService) Search(ctx context.Context, query string, page, pageSize int) ([]*models.News, int64, error) {
	if query == "" {
		return nil, 0, fmt.Errorf("search query cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Проверка кэша
	cacheKey := fmt.Sprintf("news:search:%s:%d:%d", query, page, pageSize)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedData struct {
			News  []*models.News `json:"news"`
			Total int64          `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &cachedData) == nil {
			logger.Info("Search results from cache", zap.String("query", query))
			return cachedData.News, cachedData.Total, nil
		}
	}

	// Поиск в базе данных
	news, total, err := s.newsRepo.FullTextSearch(ctx, query, page, pageSize)
	if err != nil {
		logger.Error("Failed to search news", zap.Error(err), zap.String("query", query))
		return nil, 0, err
	}

	// Кэширование результатов на 5 минут
	cachedData := struct {
		News  []*models.News `json:"news"`
		Total int64          `json:"total"`
	}{
		News:  news,
		Total: total,
	}
	data, _ := json.Marshal(cachedData)
	s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	logger.Info("Search completed",
		zap.String("query", query),
		zap.Int64("total", total),
		zap.Int("results", len(news)),
	)

	return news, total, nil
}
