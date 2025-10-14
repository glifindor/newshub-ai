package repository

import (
	"context"

	"news-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository interface {
	Create(ctx context.Context, tag *models.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*models.Tag, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Tag, error)
	Update(ctx context.Context, tag *models.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*models.Tag, error)
	Search(ctx context.Context, query string, limit int) ([]*models.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).First(&tag, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetBySlug(ctx context.Context, slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).First(&tag, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.WithContext(ctx).Find(&tags, "id IN ?", ids).Error
	return tags, err
}

func (r *tagRepository) Update(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Tag{}, "id = ?", id).Error
}

func (r *tagRepository) List(ctx context.Context) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.WithContext(ctx).Order("name ASC").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) Search(ctx context.Context, query string, limit int) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+query+"%").
		Limit(limit).
		Find(&tags).Error

	return tags, err
}
