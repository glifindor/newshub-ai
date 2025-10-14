package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"media-service/internal/models"
	"media-service/internal/repository"
	"media-service/internal/service"
	"media-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type HTTPHandler struct {
	mediaService service.MediaService
}

func NewHTTPHandler(mediaService service.MediaService) *HTTPHandler {
	return &HTTPHandler{
		mediaService: mediaService,
	}
}

// Upload загружает файл
// @Summary Upload file
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param alt_text formData string false "Alt text"
// @Param caption formData string false "Caption"
// @Param folder formData string false "Folder"
// @Param is_public formData bool false "Is public"
// @Success 201 {object} models.MediaResponse
// @Router /api/v1/media/upload [post]
func (h *HTTPHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("Failed to get file from request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// TODO: Получить user_id из JWT
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	req := &models.MediaUploadRequest{
		AltText:  c.PostForm("alt_text"),
		Caption:  c.PostForm("caption"),
		Folder:   c.PostForm("folder"),
		IsPublic: c.PostForm("is_public") == "true",
	}

	media, err := h.mediaService.Upload(c.Request.Context(), file, userID, req)
	if err != nil {
		logger.Error("Failed to upload file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Генерация публичного URL
	publicURL := fmt.Sprintf("/api/v1/media/file/%s", media.FileName)

	response := models.MediaResponse{
		ID:           media.ID,
		OriginalName: media.OriginalName,
		FileName:     media.FileName,
		FilePath:     media.FilePath,
		FileSize:     media.FileSize,
		MimeType:     media.MimeType,
		MediaType:    media.MediaType,
		Width:        media.Width,
		Height:       media.Height,
		Duration:     media.Duration,
		ThumbnailURL: media.ThumbnailURL,
		AltText:      media.AltText,
		Caption:      media.Caption,
		UploadedBy:   media.UploadedBy,
		Folder:       media.Folder,
		IsPublic:     media.IsPublic,
		URL:          publicURL,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetByID получает медиафайл по ID
// @Summary Get media by ID
// @Tags media
// @Produce json
// @Param id path string true "Media ID"
// @Success 200 {object} models.MediaResponse
// @Router /api/v1/media/{id} [get]
func (h *HTTPHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	media, err := h.mediaService.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.Error("Media not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	publicURL := fmt.Sprintf("/api/v1/media/file/%s", media.FileName)

	response := models.MediaResponse{
		ID:           media.ID,
		OriginalName: media.OriginalName,
		FileName:     media.FileName,
		FilePath:     media.FilePath,
		FileSize:     media.FileSize,
		MimeType:     media.MimeType,
		MediaType:    media.MediaType,
		Width:        media.Width,
		Height:       media.Height,
		Duration:     media.Duration,
		ThumbnailURL: media.ThumbnailURL,
		AltText:      media.AltText,
		Caption:      media.Caption,
		UploadedBy:   media.UploadedBy,
		Folder:       media.Folder,
		IsPublic:     media.IsPublic,
		URL:          publicURL,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// ServeFile отдает файл
// @Summary Serve file
// @Tags media
// @Produce application/octet-stream
// @Param filename path string true "File name"
// @Success 200
// @Router /api/v1/media/file/{filename} [get]
func (h *HTTPHandler) ServeFile(c *gin.Context) {
	fileName := c.Param("filename")

	url, err := h.mediaService.GetPublicURL(c.Request.Context(), fileName)
	if err != nil {
		logger.Error("Failed to get file URL", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Redirect к MinIO URL
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Update обновляет метаданные файла
// @Summary Update media metadata
// @Tags media
// @Accept json
// @Produce json
// @Param id path string true "Media ID"
// @Param request body object true "Metadata"
// @Success 200 {object} models.MediaResponse
// @Router /api/v1/media/{id} [put]
func (h *HTTPHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	var req struct {
		AltText string `json:"alt_text"`
		Caption string `json:"caption"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	media, err := h.mediaService.Update(c.Request.Context(), id, req.AltText, req.Caption)
	if err != nil {
		logger.Error("Failed to update media", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update media"})
		return
	}

	publicURL := fmt.Sprintf("/api/v1/media/file/%s", media.FileName)

	response := models.MediaResponse{
		ID:           media.ID,
		OriginalName: media.OriginalName,
		FileName:     media.FileName,
		FilePath:     media.FilePath,
		FileSize:     media.FileSize,
		MimeType:     media.MimeType,
		MediaType:    media.MediaType,
		Width:        media.Width,
		Height:       media.Height,
		ThumbnailURL: media.ThumbnailURL,
		AltText:      media.AltText,
		Caption:      media.Caption,
		UploadedBy:   media.UploadedBy,
		Folder:       media.Folder,
		IsPublic:     media.IsPublic,
		URL:          publicURL,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Delete удаляет файл
// @Summary Delete media
// @Tags media
// @Param id path string true "Media ID"
// @Success 204
// @Router /api/v1/media/{id} [delete]
func (h *HTTPHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	if err := h.mediaService.Delete(c.Request.Context(), id); err != nil {
		logger.Error("Failed to delete media", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
		return
	}

	c.Status(http.StatusNoContent)
}

// List возвращает список файлов
// @Summary List media
// @Tags media
// @Produce json
// @Param media_type query string false "Media type filter"
// @Param folder query string false "Folder filter"
// @Param search query string false "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} models.MediaListResponse
// @Router /api/v1/media [get]
func (h *HTTPHandler) List(c *gin.Context) {
	filters := repository.MediaFilters{
		Page:     1,
		PageSize: 20,
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		filters.Page = page
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil && pageSize > 0 {
		filters.PageSize = pageSize
	}

	if mediaType := c.Query("media_type"); mediaType != "" {
		mt := models.MediaType(mediaType)
		filters.MediaType = &mt
	}

	if folder := c.Query("folder"); folder != "" {
		filters.Folder = folder
	}

	if search := c.Query("search"); search != "" {
		filters.Search = search
	}

	media, total, err := h.mediaService.List(c.Request.Context(), filters)
	if err != nil {
		logger.Error("Failed to list media", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list media"})
		return
	}

	totalPages := int(total) / filters.PageSize
	if int(total)%filters.PageSize != 0 {
		totalPages++
	}

	response := models.MediaListResponse{
		Items:      make([]models.MediaResponse, len(media)),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	for i, m := range media {
		publicURL := fmt.Sprintf("/api/v1/media/file/%s", m.FileName)
		response.Items[i] = models.MediaResponse{
			ID:           m.ID,
			OriginalName: m.OriginalName,
			FileName:     m.FileName,
			FilePath:     m.FilePath,
			FileSize:     m.FileSize,
			MimeType:     m.MimeType,
			MediaType:    m.MediaType,
			Width:        m.Width,
			Height:       m.Height,
			Duration:     m.Duration,
			ThumbnailURL: m.ThumbnailURL,
			AltText:      m.AltText,
			Caption:      m.Caption,
			UploadedBy:   m.UploadedBy,
			Folder:       m.Folder,
			IsPublic:     m.IsPublic,
			URL:          publicURL,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}
