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

type CategoryService interface {
	Create(ctx context.Context, req *models.CategoryRequest) (*models.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
	Update(ctx context.Context, id uuid.UUID, req *models.CategoryRequest) (*models.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*models.Category, error)
	GetTree(ctx context.Context) ([]*models.Category, error)
}

type categoryService struct {
	repo  repository.CategoryRepository
	redis *redis.Client
}

func NewCategoryService(repo repository.CategoryRepository, redis *redis.Client) CategoryService {
	return &categoryService{
		repo:  repo,
		redis: redis,
	}
}

func (s *categoryService) Create(ctx context.Context, req *models.CategoryRequest) (*models.Category, error) {
	categorySlug := req.Slug
	if categorySlug == "" {
		categorySlug = slug.Make(req.Name)
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        categorySlug,
		Description: req.Description,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
		MetaTitle:   req.MetaTitle,
		MetaDesc:    req.MetaDesc,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		logger.Error("Failed to create category", zap.Error(err))
		return nil, err
	}

	logger.Info("Category created successfully",
		zap.String("category_id", category.ID.String()),
		zap.String("name", category.Name),
	)

	s.invalidateCache(ctx)
	return category, nil
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	cacheKey := fmt.Sprintf("category:id:%s", id.String())
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var category models.Category
		if json.Unmarshal([]byte(cached), &category) == nil {
			return &category, nil
		}
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(category)
	s.redis.Set(ctx, cacheKey, data, 15*time.Minute)

	return category, nil
}

func (s *categoryService) GetBySlug(ctx context.Context, categorySlug string) (*models.Category, error) {
	cacheKey := fmt.Sprintf("category:slug:%s", categorySlug)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var category models.Category
		if json.Unmarshal([]byte(cached), &category) == nil {
			return &category, nil
		}
	}

	category, err := s.repo.GetBySlug(ctx, categorySlug)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(category)
	s.redis.Set(ctx, cacheKey, data, 15*time.Minute)

	return category, nil
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, req *models.CategoryRequest) (*models.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	if req.Slug != "" {
		category.Slug = req.Slug
	} else {
		category.Slug = slug.Make(req.Name)
	}
	category.Description = req.Description
	category.ParentID = req.ParentID
	category.SortOrder = req.SortOrder
	category.IsActive = req.IsActive
	category.MetaTitle = req.MetaTitle
	category.MetaDesc = req.MetaDesc

	if err := s.repo.Update(ctx, category); err != nil {
		logger.Error("Failed to update category", zap.Error(err))
		return nil, err
	}

	logger.Info("Category updated successfully", zap.String("category_id", id.String()))

	s.invalidateCache(ctx)
	s.redis.Del(ctx, fmt.Sprintf("category:id:%s", id.String()))
	s.redis.Del(ctx, fmt.Sprintf("category:slug:%s", category.Slug))

	return category, nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete category", zap.Error(err))
		return err
	}

	logger.Info("Category deleted successfully", zap.String("category_id", id.String()))

	s.invalidateCache(ctx)
	s.redis.Del(ctx, fmt.Sprintf("category:id:%s", id.String()))
	s.redis.Del(ctx, fmt.Sprintf("category:slug:%s", category.Slug))

	return nil
}

func (s *categoryService) List(ctx context.Context) ([]*models.Category, error) {
	cacheKey := "category:list"
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var categories []*models.Category
		if json.Unmarshal([]byte(cached), &categories) == nil {
			return categories, nil
		}
	}

	categories, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(categories)
	s.redis.Set(ctx, cacheKey, data, 30*time.Minute)

	return categories, nil
}

func (s *categoryService) GetTree(ctx context.Context) ([]*models.Category, error) {
	cacheKey := "category:tree"
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var categories []*models.Category
		if json.Unmarshal([]byte(cached), &categories) == nil {
			return categories, nil
		}
	}

	categories, err := s.repo.GetTree(ctx)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(categories)
	s.redis.Set(ctx, cacheKey, data, 30*time.Minute)

	return categories, nil
}

func (s *categoryService) invalidateCache(ctx context.Context) {
	pattern := "category:*"
	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		s.redis.Del(ctx, iter.Val())
	}
}
