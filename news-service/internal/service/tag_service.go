package service

import (
	"context"

	"news-service/internal/models"
	"news-service/internal/repository"
	"news-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"go.uber.org/zap"
)

type TagService interface {
	Create(ctx context.Context, req *models.TagRequest) (*models.Tag, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*models.Tag, error)
	Update(ctx context.Context, id uuid.UUID, req *models.TagRequest) (*models.Tag, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*models.Tag, error)
	Search(ctx context.Context, query string, limit int) ([]*models.Tag, error)
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) Create(ctx context.Context, req *models.TagRequest) (*models.Tag, error) {
	tagSlug := req.Slug
	if tagSlug == "" {
		tagSlug = slug.Make(req.Name)
	}

	tag := &models.Tag{
		Name: req.Name,
		Slug: tagSlug,
	}

	if err := s.repo.Create(ctx, tag); err != nil {
		logger.Error("Failed to create tag", zap.Error(err))
		return nil, err
	}

	logger.Info("Tag created successfully",
		zap.String("tag_id", tag.ID.String()),
		zap.String("name", tag.Name),
	)

	return tag, nil
}

func (s *tagService) GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *tagService) GetBySlug(ctx context.Context, tagSlug string) (*models.Tag, error) {
	return s.repo.GetBySlug(ctx, tagSlug)
}

func (s *tagService) Update(ctx context.Context, id uuid.UUID, req *models.TagRequest) (*models.Tag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tag.Name = req.Name
	if req.Slug != "" {
		tag.Slug = req.Slug
	} else {
		tag.Slug = slug.Make(req.Name)
	}

	if err := s.repo.Update(ctx, tag); err != nil {
		logger.Error("Failed to update tag", zap.Error(err))
		return nil, err
	}

	logger.Info("Tag updated successfully", zap.String("tag_id", id.String()))

	return tag, nil
}

func (s *tagService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete tag", zap.Error(err))
		return err
	}

	logger.Info("Tag deleted successfully", zap.String("tag_id", id.String()))
	return nil
}

func (s *tagService) List(ctx context.Context) ([]*models.Tag, error) {
	return s.repo.List(ctx)
}

func (s *tagService) Search(ctx context.Context, query string, limit int) ([]*models.Tag, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.Search(ctx, query, limit)
}
