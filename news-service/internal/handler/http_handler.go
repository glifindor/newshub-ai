package handler

import (
	"net/http"
	"strconv"

	"news-service/internal/models"
	"news-service/internal/repository"
	"news-service/internal/service"
	"news-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type HTTPHandler struct {
	newsService     service.NewsService
	categoryService service.CategoryService
	tagService      service.TagService
}

func NewHTTPHandler(
	newsService service.NewsService,
	categoryService service.CategoryService,
	tagService service.TagService,
) *HTTPHandler {
	return &HTTPHandler{
		newsService:     newsService,
		categoryService: categoryService,
		tagService:      tagService,
	}
}

// News handlers

// CreateNews создает новость
// @Summary Create news
// @Tags news
// @Accept json
// @Produce json
// @Param request body models.NewsRequest true "News data"
// @Success 201 {object} models.NewsResponse
// @Router /api/v1/news [post]
func (h *HTTPHandler) CreateNews(c *gin.Context) {
	var req models.NewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// TODO: Получить author_id из JWT токена
	authorID := uuid.MustParse("00000000-0000-0000-0000-000000000000") // Временно

	news, err := h.newsService.Create(c.Request.Context(), authorID, &req)
	if err != nil {
		logger.Error("Failed to create news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news"})
		return
	}

	c.JSON(http.StatusCreated, news)
}

// GetNewsBySlug получает новость по slug
// @Summary Get news by slug
// @Tags news
// @Produce json
// @Param slug path string true "News slug"
// @Success 200 {object} models.NewsResponse
// @Router /api/v1/news/{slug} [get]
func (h *HTTPHandler) GetNewsBySlug(c *gin.Context) {
	slug := c.Param("slug")

	news, err := h.newsService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		logger.Error("News not found", zap.Error(err), zap.String("slug", slug))
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	// Инкремент просмотров
	go h.newsService.IncrementViewCount(c.Request.Context(), news.ID)

	c.JSON(http.StatusOK, news)
}

// UpdateNews обновляет новость
// @Summary Update news
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Param request body models.NewsRequest true "News data"
// @Success 200 {object} models.NewsResponse
// @Router /api/v1/news/{id} [put]
func (h *HTTPHandler) UpdateNews(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	var req models.NewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	news, err := h.newsService.Update(c.Request.Context(), id, &req)
	if err != nil {
		logger.Error("Failed to update news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news"})
		return
	}

	c.JSON(http.StatusOK, news)
}

// DeleteNews удаляет новость
// @Summary Delete news
// @Tags news
// @Param id path string true "News ID"
// @Success 204
// @Router /api/v1/news/{id} [delete]
func (h *HTTPHandler) DeleteNews(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	if err := h.newsService.Delete(c.Request.Context(), id); err != nil {
		logger.Error("Failed to delete news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListNews возвращает список новостей с фильтрами
// @Summary List news
// @Tags news
// @Produce json
// @Param status query string false "Status filter"
// @Param category_id query string false "Category ID filter"
// @Param search query string false "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.NewsListResponse
// @Router /api/v1/news [get]
func (h *HTTPHandler) ListNews(c *gin.Context) {
	filters := repository.NewsFilters{
		Page:     1,
		PageSize: 10,
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		filters.Page = page
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10")); err == nil && pageSize > 0 {
		filters.PageSize = pageSize
	}

	if status := c.Query("status"); status != "" {
		s := models.NewsStatus(status)
		filters.Status = &s
	}

	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			filters.CategoryID = &id
		}
	}

	if search := c.Query("search"); search != "" {
		filters.Search = search
	}

	news, total, err := h.newsService.List(c.Request.Context(), filters)
	if err != nil {
		logger.Error("Failed to list news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list news"})
		return
	}

	totalPages := int(total) / filters.PageSize
	if int(total)%filters.PageSize != 0 {
		totalPages++
	}

	response := models.NewsListResponse{
		Items:      make([]models.NewsResponse, len(news)),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	for i, n := range news {
		response.Items[i] = models.NewsResponse{
			ID:              n.ID,
			Title:           n.Title,
			Slug:            n.Slug,
			Summary:         n.Summary,
			Content:         n.Content,
			FeaturedImage:   n.FeaturedImage,
			AuthorID:        n.AuthorID,
			CategoryID:      n.CategoryID,
			Status:          n.Status,
			PublishedAt:     n.PublishedAt,
			ViewCount:       n.ViewCount,
			IsFeatured:      n.IsFeatured,
			IsBreaking:      n.IsBreaking,
			AllowComments:   n.AllowComments,
			MetaTitle:       n.MetaTitle,
			MetaDescription: n.MetaDescription,
			MetaKeywords:    n.MetaKeywords,
			CreatedAt:       n.CreatedAt,
			UpdatedAt:       n.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// PublishNews публикует новость
// @Summary Publish news
// @Tags news
// @Param id path string true "News ID"
// @Success 200 {object} gin.H
// @Router /api/v1/news/{id}/publish [post]
func (h *HTTPHandler) PublishNews(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	if err := h.newsService.Publish(c.Request.Context(), id); err != nil {
		logger.Error("Failed to publish news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish news"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News published successfully"})
}

// GetFeaturedNews возвращает избранные новости
// @Summary Get featured news
// @Tags news
// @Produce json
// @Param limit query int false "Limit" default(5)
// @Success 200 {array} models.NewsResponse
// @Router /api/v1/news/featured [get]
func (h *HTTPHandler) GetFeaturedNews(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	news, err := h.newsService.GetFeatured(c.Request.Context(), limit)
	if err != nil {
		logger.Error("Failed to get featured news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get featured news"})
		return
	}

	c.JSON(http.StatusOK, news)
}

// GetBreakingNews возвращает срочные новости
// @Summary Get breaking news
// @Tags news
// @Produce json
// @Param limit query int false "Limit" default(3)
// @Success 200 {array} models.NewsResponse
// @Router /api/v1/news/breaking [get]
func (h *HTTPHandler) GetBreakingNews(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))

	news, err := h.newsService.GetBreaking(c.Request.Context(), limit)
	if err != nil {
		logger.Error("Failed to get breaking news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get breaking news"})
		return
	}

	c.JSON(http.StatusOK, news)
}

// Category handlers

// CreateCategory создает категорию
func (h *HTTPHandler) CreateCategory(c *gin.Context) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to create category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategoryBySlug получает категорию по slug
func (h *HTTPHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.categoryService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory обновляет категорию
func (h *HTTPHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), id, &req)
	if err != nil {
		logger.Error("Failed to update category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory удаляет категорию
func (h *HTTPHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		logger.Error("Failed to delete category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListCategories возвращает список категорий
func (h *HTTPHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.List(c.Request.Context())
	if err != nil {
		logger.Error("Failed to list categories", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryTree возвращает дерево категорий
func (h *HTTPHandler) GetCategoryTree(c *gin.Context) {
	tree, err := h.categoryService.GetTree(c.Request.Context())
	if err != nil {
		logger.Error("Failed to get category tree", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get category tree"})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// Tag handlers

// CreateTag создает тег
func (h *HTTPHandler) CreateTag(c *gin.Context) {
	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tag, err := h.tagService.Create(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to create tag", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// UpdateTag обновляет тег
func (h *HTTPHandler) UpdateTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tag, err := h.tagService.Update(c.Request.Context(), id, &req)
	if err != nil {
		logger.Error("Failed to update tag", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tag"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag удаляет тег
func (h *HTTPHandler) DeleteTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.tagService.Delete(c.Request.Context(), id); err != nil {
		logger.Error("Failed to delete tag", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTags возвращает список тегов
func (h *HTTPHandler) ListTags(c *gin.Context) {
	tags, err := h.tagService.List(c.Request.Context())
	if err != nil {
		logger.Error("Failed to list tags", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// SearchTags поиск тегов
func (h *HTTPHandler) SearchTags(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tags, err := h.tagService.Search(c.Request.Context(), query, limit)
	if err != nil {
		logger.Error("Failed to search tags", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// SearchNews выполняет полнотекстовый поиск новостей
// @Summary Full Text Search for news
// @Tags news
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/news/search [get]
func (h *HTTPHandler) SearchNews(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query parameter 'q' is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	news, total, err := h.newsService.Search(c.Request.Context(), query, page, limit)
	if err != nil {
		logger.Error("Failed to search news", zap.Error(err), zap.String("query", query))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search news",
			"details": err.Error(),
		})
		return
	}

	// Вычисляем количество страниц
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data": news,
		"pagination": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
		"query": query,
	})
}
