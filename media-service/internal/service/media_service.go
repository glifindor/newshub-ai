package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"media-service/internal/config"
	"media-service/internal/models"
	"media-service/internal/repository"
	"media-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type MediaService interface {
	Upload(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID, req *models.MediaUploadRequest) (*models.Media, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Media, error)
	GetPublicURL(ctx context.Context, fileName string) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters repository.MediaFilters) ([]*models.Media, int64, error)
	Update(ctx context.Context, id uuid.UUID, altText, caption string) (*models.Media, error)
}

type mediaService struct {
	repo        repository.MediaRepository
	minioClient *minio.Client
	cfg         *config.Config
}

func NewMediaService(repo repository.MediaRepository, cfg *config.Config) (MediaService, error) {
	minioClient, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Создание bucket если не существует
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.MinIOBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.MinIOBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		logger.Info("Created MinIO bucket", zap.String("bucket", cfg.MinIOBucket))
	}

	return &mediaService{
		repo:        repo,
		minioClient: minioClient,
		cfg:         cfg,
	}, nil
}

func (s *mediaService) Upload(ctx context.Context, fileHeader *multipart.FileHeader, userID uuid.UUID, req *models.MediaUploadRequest) (*models.Media, error) {
	// Проверка размера
	if fileHeader.Size > s.cfg.MaxUploadSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.cfg.MaxUploadSize)
	}

	// Проверка типа файла
	mimeType := fileHeader.Header.Get("Content-Type")
	if !s.isAllowedType(mimeType) {
		return nil, fmt.Errorf("file type %s is not allowed", mimeType)
	}

	// Открытие файла
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Генерация уникального имени файла
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	folder := req.Folder
	if folder == "" {
		folder = "uploads"
	}

	objectName := fmt.Sprintf("%s/%s", folder, fileName)

	// Загрузка в MinIO
	_, err = s.minioClient.PutObject(ctx, s.cfg.MinIOBucket, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		logger.Error("Failed to upload to MinIO", zap.Error(err))
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Определение типа медиа
	mediaType := s.getMediaType(mimeType)

	// Создание записи в БД
	media := &models.Media{
		OriginalName: fileHeader.Filename,
		FileName:     fileName,
		FilePath:     objectName,
		FileSize:     fileHeader.Size,
		MimeType:     mimeType,
		MediaType:    mediaType,
		AltText:      req.AltText,
		Caption:      req.Caption,
		UploadedBy:   userID,
		Folder:       folder,
		IsPublic:     req.IsPublic,
	}

	if err := s.repo.Create(ctx, media); err != nil {
		// Откат: удаление файла из MinIO
		s.minioClient.RemoveObject(ctx, s.cfg.MinIOBucket, objectName, minio.RemoveObjectOptions{})
		logger.Error("Failed to create media record", zap.Error(err))
		return nil, err
	}

	logger.Info("File uploaded successfully",
		zap.String("file_id", media.ID.String()),
		zap.String("file_name", fileName),
		zap.String("user_id", userID.String()),
	)

	return media, nil
}

func (s *mediaService) GetByID(ctx context.Context, id uuid.UUID) (*models.Media, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *mediaService) GetPublicURL(ctx context.Context, fileName string) (string, error) {
	media, err := s.repo.GetByFileName(ctx, fileName)
	if err != nil {
		return "", err
	}

	if !media.IsPublic {
		return "", fmt.Errorf("file is not public")
	}

	// Генерация presigned URL (действителен 7 дней)
	url, err := s.minioClient.PresignedGetObject(ctx, s.cfg.MinIOBucket, media.FilePath, 7*24*time.Hour, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (s *mediaService) Delete(ctx context.Context, id uuid.UUID) error {
	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Удаление из MinIO
	err = s.minioClient.RemoveObject(ctx, s.cfg.MinIOBucket, media.FilePath, minio.RemoveObjectOptions{})
	if err != nil {
		logger.Warn("Failed to delete from MinIO", zap.Error(err))
	}

	// Удаление из БД
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete media record", zap.Error(err))
		return err
	}

	logger.Info("Media deleted successfully", zap.String("media_id", id.String()))
	return nil
}

func (s *mediaService) List(ctx context.Context, filters repository.MediaFilters) ([]*models.Media, int64, error) {
	return s.repo.List(ctx, filters)
}

func (s *mediaService) Update(ctx context.Context, id uuid.UUID, altText, caption string) (*models.Media, error) {
	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	media.AltText = altText
	media.Caption = caption

	if err := s.repo.Update(ctx, media); err != nil {
		logger.Error("Failed to update media", zap.Error(err))
		return nil, err
	}

	logger.Info("Media updated successfully", zap.String("media_id", id.String()))
	return media, nil
}

func (s *mediaService) isAllowedType(mimeType string) bool {
	for _, allowed := range s.cfg.AllowedTypes {
		if allowed == mimeType {
			return true
		}
	}
	return false
}

func (s *mediaService) getMediaType(mimeType string) models.MediaType {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return models.MediaTypeImage
	case strings.HasPrefix(mimeType, "video/"):
		return models.MediaTypeVideo
	case mimeType == "application/pdf":
		return models.MediaTypePDF
	default:
		return models.MediaTypeImage
	}
}

func (s *mediaService) DownloadFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	media, err := s.repo.GetByFileName(ctx, fileName)
	if err != nil {
		return nil, err
	}

	object, err := s.minioClient.GetObject(ctx, s.cfg.MinIOBucket, media.FilePath, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("Failed to download from MinIO", zap.Error(err))
		return nil, err
	}

	return object, nil
}
