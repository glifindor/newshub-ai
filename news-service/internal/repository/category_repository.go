package repository

import (
	"context"

	"news-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*models.Category, error)
	GetTree(ctx context.Context) ([]*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&category, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&category, "slug = ?", slug).Error

	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id).Error
}

func (r *categoryRepository) List(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) GetTree(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL AND is_active = ?", true).
		Preload("Children", "is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error

	return categories, err
}
