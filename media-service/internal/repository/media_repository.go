package repository

import (
	"context"
	"fmt"

	"media-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(ctx context.Context, media *models.Media) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Media, error)
	GetByFileName(ctx context.Context, fileName string) (*models.Media, error)
	Update(ctx context.Context, media *models.Media) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters MediaFilters) ([]*models.Media, int64, error)
}

type MediaFilters struct {
	MediaType  *models.MediaType
	UploadedBy *uuid.UUID
	Folder     string
	IsPublic   *bool
	Search     string
	Page       int
	PageSize   int
}

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *mediaRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Media, error) {
	var media models.Media
	err := r.db.WithContext(ctx).First(&media, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *mediaRepository) GetByFileName(ctx context.Context, fileName string) (*models.Media, error) {
	var media models.Media
	err := r.db.WithContext(ctx).First(&media, "file_name = ?", fileName).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *mediaRepository) Update(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Save(media).Error
}

func (r *mediaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Media{}, "id = ?", id).Error
}

func (r *mediaRepository) List(ctx context.Context, filters MediaFilters) ([]*models.Media, int64, error) {
	var media []*models.Media
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Media{})

	if filters.MediaType != nil {
		query = query.Where("media_type = ?", *filters.MediaType)
	}
	if filters.UploadedBy != nil {
		query = query.Where("uploaded_by = ?", *filters.UploadedBy)
	}
	if filters.Folder != "" {
		query = query.Where("folder = ?", filters.Folder)
	}
	if filters.IsPublic != nil {
		query = query.Where("is_public = ?", *filters.IsPublic)
	}
	if filters.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", filters.Search)
		query = query.Where("original_name ILIKE ? OR alt_text ILIKE ? OR caption ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}
	offset := (filters.Page - 1) * filters.PageSize

	err := query.
		Order("created_at DESC").
		Limit(filters.PageSize).
		Offset(offset).
		Find(&media).Error

	return media, total, err
}
