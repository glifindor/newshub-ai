package handler

import (
	"net/http"

	"seo-service/internal/models"
	"seo-service/internal/service"
	"seo-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// HTTPHandler обработчик HTTP запросов
type HTTPHandler struct {
	seoService     service.SEOService
	ogService      service.OpenGraphService
	sitemapService service.SitemapService
	robotsService  service.RobotsService
	logger         *logger.Logger
}

// NewHTTPHandler создает новый HTTP handler
func NewHTTPHandler(
	seoService service.SEOService,
	ogService service.OpenGraphService,
	sitemapService service.SitemapService,
	robotsService service.RobotsService,
	logger *logger.Logger,
) *HTTPHandler {
	return &HTTPHandler{
		seoService:     seoService,
		ogService:      ogService,
		sitemapService: sitemapService,
		robotsService:  robotsService,
		logger:         logger,
	}
}

// RegisterRoutes регистрирует все роуты
func (h *HTTPHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// SEO метаданные
		api.GET("/seo/:slug", h.GetSEOBySlug)
		api.POST("/seo/create", h.CreateSEO)
		api.PUT("/seo/update", h.UpdateSEO)
		api.DELETE("/seo/:newsId", h.DeleteSEO)

		// Webhook от news-service
		api.POST("/webhook/news-published", h.NewsPublishedWebhook)

		// Open Graph теги
		api.GET("/seo/:slug/og-tags", h.GetOGTags)
	}

	// Публичные endpoints (без /api/v1)
	r.GET("/sitemap.xml", h.GetSitemap)
	r.GET("/robots.txt", h.GetRobotsTxt)

	// Health check
	r.GET("/health", h.HealthCheck)
}

// GetSEOBySlug получает SEO метаданные по slug
// @Summary Получить SEO метаданные
// @Tags SEO
// @Param slug path string true "Slug новости"
// @Success 200 {object} models.SEOResponse
// @Router /api/v1/seo/{slug} [get]
func (h *HTTPHandler) GetSEOBySlug(c *gin.Context) {
	slug := c.Param("slug")

	seoMeta, err := h.seoService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		h.logger.Error("Failed to get SEO by slug", zap.Error(err), zap.String("slug", slug))
		c.JSON(http.StatusNotFound, gin.H{"error": "SEO metadata not found"})
		return
	}

	c.JSON(http.StatusOK, seoMeta)
}

// CreateSEO создает SEO метаданные
// @Summary Создать SEO метаданные
// @Tags SEO
// @Accept json
// @Produce json
// @Param request body models.CreateSEORequest true "SEO данные"
// @Success 201 {object} models.SEOResponse
// @Router /api/v1/seo/create [post]
func (h *HTTPHandler) CreateSEO(c *gin.Context) {
	var req models.CreateSEORequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	seoMeta, err := h.seoService.CreateOrUpdate(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create SEO", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SEO metadata"})
		return
	}

	c.JSON(http.StatusCreated, seoMeta)
}

// UpdateSEO обновляет SEO метаданные
// @Summary Обновить SEO метаданные
// @Tags SEO
// @Accept json
// @Produce json
// @Param request body models.UpdateSEORequest true "Обновленные SEO данные"
// @Success 200 {object} models.SEOResponse
// @Router /api/v1/seo/update [put]
func (h *HTTPHandler) UpdateSEO(c *gin.Context) {
	var req models.UpdateSEORequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Получаем существующую запись
	existing, err := h.seoService.GetByNewsID(c.Request.Context(), req.NewsID)
	if err != nil {
		h.logger.Error("SEO metadata not found", zap.Error(err), zap.String("news_id", req.NewsID.String()))
		c.JSON(http.StatusNotFound, gin.H{"error": "SEO metadata not found"})
		return
	}

	// Создаем CreateSEORequest из UpdateSEORequest для обновления
	createReq := &models.CreateSEORequest{
		NewsID:            req.NewsID,
		Slug:              existing.Slug, // Slug не меняется
		Title:             req.Title,
		Description:       req.Description,
		Keywords:          req.Keywords,
		OGTitle:           req.OGTitle,
		OGDescription:     req.OGDescription,
		OGImage:           req.OGImage,
		RobotsIndex:       req.RobotsIndex,
		RobotsFollow:      req.RobotsFollow,
		SitemapChangeFreq: req.SitemapChangeFreq,
		SitemapPriority:   req.SitemapPriority,
	}

	seoMeta, err := h.seoService.CreateOrUpdate(c.Request.Context(), createReq)
	if err != nil {
		h.logger.Error("Failed to update SEO", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SEO metadata"})
		return
	}

	// Инвалидируем кеш sitemap
	_ = h.sitemapService.InvalidateCache(c.Request.Context())

	c.JSON(http.StatusOK, seoMeta)
}

// DeleteSEO удаляет SEO метаданные
// @Summary Удалить SEO метаданные
// @Tags SEO
// @Param newsId path string true "ID новости"
// @Success 204
// @Router /api/v1/seo/{newsId} [delete]
func (h *HTTPHandler) DeleteSEO(c *gin.Context) {
	newsIDStr := c.Param("newsId")

	newsID, err := uuid.Parse(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	if err := h.seoService.Delete(c.Request.Context(), newsID); err != nil {
		h.logger.Error("Failed to delete SEO", zap.Error(err), zap.String("news_id", newsID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SEO metadata"})
		return
	}

	// Инвалидируем кеш sitemap
	_ = h.sitemapService.InvalidateCache(c.Request.Context())

	c.Status(http.StatusNoContent)
}

// NewsPublishedWebhook обрабатывает webhook о публикации новости
// @Summary Webhook при публикации новости
// @Tags Webhook
// @Accept json
// @Produce json
// @Param webhook body models.NewsPublishedWebhook true "Данные новости"
// @Success 201 {object} models.SEOResponse
// @Router /api/v1/webhook/news-published [post]
func (h *HTTPHandler) NewsPublishedWebhook(c *gin.Context) {
	var webhook models.NewsPublishedWebhook

	if err := c.ShouldBindJSON(&webhook); err != nil {
		h.logger.Warn("Invalid webhook body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook body"})
		return
	}

	h.logger.Info("News published webhook received",
		zap.String("news_id", webhook.NewsID.String()),
		zap.String("slug", webhook.Slug),
	)

	// Автоматически генерируем SEO из новости
	seoMeta, err := h.seoService.GenerateFromNews(c.Request.Context(), &webhook)
	if err != nil {
		h.logger.Error("Failed to generate SEO from news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate SEO metadata"})
		return
	}

	// Инвалидируем кеш sitemap
	_ = h.sitemapService.InvalidateCache(c.Request.Context())

	c.JSON(http.StatusCreated, seoMeta)
}

// GetOGTags получает Open Graph теги для соцсетей (VK, Telegram, OK)
// @Summary Получить Open Graph теги
// @Tags Open Graph
// @Param slug path string true "Slug новости"
// @Success 200 {object} map[string]string
// @Router /api/v1/seo/{slug}/og-tags [get]
func (h *HTTPHandler) GetOGTags(c *gin.Context) {
	slug := c.Param("slug")

	// Получаем SEO метаданные
	seoMeta, err := h.seoService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		h.logger.Error("Failed to get SEO by slug", zap.Error(err), zap.String("slug", slug))
		c.JSON(http.StatusNotFound, gin.H{"error": "SEO metadata not found"})
		return
	}

	// Конвертируем в модель для генератора
	meta := &models.SEOMeta{
		Slug:          seoMeta.Slug,
		OGTitle:       seoMeta.OGTitle,
		OGDescription: seoMeta.OGDescription,
		OGImage:       seoMeta.OGImage,
		OGType:        seoMeta.OGType,
		OGLocale:      seoMeta.OGLocale,
		OGSiteName:    seoMeta.OGSiteName,
		CanonicalURL:  seoMeta.CanonicalURL,
		CreatedAt:     seoMeta.CreatedAt,
		UpdatedAt:     seoMeta.UpdatedAt,
	}

	// Генерируем OG теги
	ogTags := h.ogService.GenerateOGTags(c.Request.Context(), meta)

	c.JSON(http.StatusOK, gin.H{
		"slug":    slug,
		"og_tags": ogTags,
	})
}

// GetSitemap возвращает sitemap.xml
// @Summary Получить sitemap.xml
// @Tags Sitemap
// @Produce xml
// @Success 200 {string} string "XML sitemap"
// @Router /sitemap.xml [get]
func (h *HTTPHandler) GetSitemap(c *gin.Context) {
	sitemap, err := h.sitemapService.GenerateSitemap(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to generate sitemap", zap.Error(err))
		c.String(http.StatusInternalServerError, "Failed to generate sitemap")
		return
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, sitemap)
}

// GetRobotsTxt возвращает robots.txt
// @Summary Получить robots.txt
// @Tags Robots
// @Produce plain
// @Success 200 {string} string "robots.txt content"
// @Router /robots.txt [get]
func (h *HTTPHandler) GetRobotsTxt(c *gin.Context) {
	robotsTxt, err := h.robotsService.GenerateRobotsTxt(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to generate robots.txt", zap.Error(err))
		c.String(http.StatusInternalServerError, "Failed to generate robots.txt")
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, robotsTxt)
}

// HealthCheck проверка работоспособности сервиса
// @Summary Health check
// @Tags Health
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HTTPHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "seo-service",
		"version": "1.0.0",
	})
}
