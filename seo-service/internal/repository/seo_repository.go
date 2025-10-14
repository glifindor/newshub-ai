package repository

import (
	"context"
	"fmt"

	"seo-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SEORepository интерфейс для работы с SEO метаданными
type SEORepository interface {
	// CRUD операции
	Create(ctx context.Context, meta *models.SEOMeta) error
	GetBySlug(ctx context.Context, slug string) (*models.SEOMeta, error)
	GetByNewsID(ctx context.Context, newsID uuid.UUID) (*models.SEOMeta, error)
	Update(ctx context.Context, meta *models.SEOMeta) error
	Delete(ctx context.Context, newsID uuid.UUID) error

	// Для sitemap.xml
	GetAllIndexable(ctx context.Context) ([]*models.SEOMeta, error)
	GetRecentlyUpdated(ctx context.Context, limit int) ([]*models.SEOMeta, error)
}

type seoRepository struct {
	db *gorm.DB
}

// NewSEORepository создает новый экземпляр SEO repository
func NewSEORepository(db *gorm.DB) SEORepository {
	return &seoRepository{db: db}
}

// Create создает новую SEO запись
func (r *seoRepository) Create(ctx context.Context, meta *models.SEOMeta) error {
	if err := r.db.WithContext(ctx).Create(meta).Error; err != nil {
		return fmt.Errorf("failed to create SEO meta: %w", err)
	}
	return nil
}

// GetBySlug получает SEO метаданные по slug
func (r *seoRepository) GetBySlug(ctx context.Context, slug string) (*models.SEOMeta, error) {
	var meta models.SEOMeta

	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&meta).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("SEO meta not found for slug: %s", slug)
		}
		return nil, fmt.Errorf("failed to get SEO meta by slug: %w", err)
	}

	return &meta, nil
}

// GetByNewsID получает SEO метаданные по ID новости
func (r *seoRepository) GetByNewsID(ctx context.Context, newsID uuid.UUID) (*models.SEOMeta, error) {
	var meta models.SEOMeta

	err := r.db.WithContext(ctx).
		Where("news_id = ?", newsID).
		First(&meta).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("SEO meta not found for news_id: %s", newsID.String())
		}
		return nil, fmt.Errorf("failed to get SEO meta by news_id: %w", err)
	}

	return &meta, nil
}

// Update обновляет SEO метаданные
func (r *seoRepository) Update(ctx context.Context, meta *models.SEOMeta) error {
	result := r.db.WithContext(ctx).
		Model(&models.SEOMeta{}).
		Where("id = ?", meta.ID).
		Updates(meta)

	if result.Error != nil {
		return fmt.Errorf("failed to update SEO meta: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("SEO meta not found: %s", meta.ID.String())
	}

	return nil
}

// Delete удаляет SEO метаданные по ID новости
func (r *seoRepository) Delete(ctx context.Context, newsID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("news_id = ?", newsID).
		Delete(&models.SEOMeta{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete SEO meta: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("SEO meta not found for news_id: %s", newsID.String())
	}

	return nil
}

// GetAllIndexable получает все SEO метаданные для индексации (для sitemap.xml)
// Возвращает только записи с robots_index = true, отсортированные по приоритету и дате
func (r *seoRepository) GetAllIndexable(ctx context.Context) ([]*models.SEOMeta, error) {
	var metaList []*models.SEOMeta

	err := r.db.WithContext(ctx).
		Where("robots_index = ?", true).
		Order("sitemap_priority DESC, updated_at DESC").
		Find(&metaList).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get indexable SEO meta: %w", err)
	}

	return metaList, nil
}

// GetRecentlyUpdated получает недавно обновленные SEO метаданные
// Полезно для инкрементального обновления sitemap
func (r *seoRepository) GetRecentlyUpdated(ctx context.Context, limit int) ([]*models.SEOMeta, error) {
	var metaList []*models.SEOMeta

	err := r.db.WithContext(ctx).
		Where("robots_index = ?", true).
		Order("updated_at DESC").
		Limit(limit).
		Find(&metaList).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recently updated SEO meta: %w", err)
	}

	return metaList, nil
}
