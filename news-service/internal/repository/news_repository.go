package repository

import (
	"context"
	"fmt"

	"news-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewsRepository interface {
	Create(ctx context.Context, news *models.News) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.News, error)
	GetBySlug(ctx context.Context, slug string) (*models.News, error)
	Update(ctx context.Context, news *models.News) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters NewsFilters) ([]*models.News, int64, error)
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	GetFeatured(ctx context.Context, limit int) ([]*models.News, error)
	GetBreaking(ctx context.Context, limit int) ([]*models.News, error)

	// Full Text Search
	FullTextSearch(ctx context.Context, query string, page, pageSize int) ([]*models.News, int64, error)
}

type NewsFilters struct {
	Status     *models.NewsStatus
	CategoryID *uuid.UUID
	AuthorID   *uuid.UUID
	IsFeatured *bool
	IsBreaking *bool
	Search     string
	Page       int
	PageSize   int
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) Create(ctx context.Context, news *models.News) error {
	return r.db.WithContext(ctx).Create(news).Error
}

func (r *newsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.News, error) {
	var news models.News
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Tags").
		First(&news, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsRepository) GetBySlug(ctx context.Context, slug string) (*models.News, error) {
	var news models.News
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Tags").
		First(&news, "slug = ?", slug).Error

	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsRepository) Update(ctx context.Context, news *models.News) error {
	return r.db.WithContext(ctx).Save(news).Error
}

func (r *newsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.News{}, "id = ?", id).Error
}

func (r *newsRepository) List(ctx context.Context, filters NewsFilters) ([]*models.News, int64, error) {
	var news []*models.News
	var total int64

	query := r.db.WithContext(ctx).Model(&models.News{})

	// Применяем фильтры
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.AuthorID != nil {
		query = query.Where("author_id = ?", *filters.AuthorID)
	}
	if filters.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filters.IsFeatured)
	}
	if filters.IsBreaking != nil {
		query = query.Where("is_breaking = ?", *filters.IsBreaking)
	}
	if filters.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", filters.Search)
		query = query.Where("title ILIKE ? OR summary ILIKE ? OR content ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Подсчет общего количества
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Пагинация
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}
	offset := (filters.Page - 1) * filters.PageSize

	err := query.
		Preload("Category").
		Preload("Tags").
		Order("created_at DESC").
		Limit(filters.PageSize).
		Offset(offset).
		Find(&news).Error

	return news, total, err
}

func (r *newsRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.News{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).
		Error
}

func (r *newsRepository) GetFeatured(ctx context.Context, limit int) ([]*models.News, error) {
	var news []*models.News
	err := r.db.WithContext(ctx).
		Where("is_featured = ? AND status = ?", true, models.StatusPublished).
		Preload("Category").
		Preload("Tags").
		Order("published_at DESC").
		Limit(limit).
		Find(&news).Error

	return news, err
}

func (r *newsRepository) GetBreaking(ctx context.Context, limit int) ([]*models.News, error) {
	var news []*models.News
	err := r.db.WithContext(ctx).
		Where("is_breaking = ? AND status = ?", true, models.StatusPublished).
		Preload("Category").
		Preload("Tags").
		Order("published_at DESC").
		Limit(limit).
		Find(&news).Error

	return news, err
}

// FullTextSearch выполняет полнотекстовый поиск с использованием PostgreSQL FTS
func (r *newsRepository) FullTextSearch(ctx context.Context, query string, page, pageSize int) ([]*models.News, int64, error) {
	var news []*models.News
	var total int64

	// Подготовка поискового запроса с ранжированием
	// Используем search_vector колонку с весами (A=title, B=summary, C=content, D=keywords)
	searchDB := r.db.WithContext(ctx).
		Table("news").
		Select("news.*, ts_rank(search_vector, plainto_tsquery('english', ?)) as rank", query).
		Where("search_vector @@ plainto_tsquery('english', ?)", query).
		Where("status = ?", models.StatusPublished).
		Order("rank DESC, published_at DESC")

	// Подсчет общего количества результатов
	if err := searchDB.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Пагинация
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// Получение результатов с пагинацией
	err := searchDB.
		Preload("Category").
		Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Find(&news).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to search news: %w", err)
	}

	return news, total, nil
}
