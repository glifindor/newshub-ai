package handler

import (
	"net/http"
	"strconv"

	"admin-service/internal/client"
	"admin-service/internal/middleware"
	"admin-service/internal/models"
	"admin-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AdminHandler обработчик админ панели
type AdminHandler struct {
	authClient      *client.AuthClient
	newsClient      *client.NewsClient
	seoClient       *client.SEOClient
	logger          *logger.Logger
	defaultPageSize int
	maxPageSize     int
}

// NewAdminHandler создает новый handler
func NewAdminHandler(
	authClient *client.AuthClient,
	newsClient *client.NewsClient,
	seoClient *client.SEOClient,
	logger *logger.Logger,
	defaultPageSize, maxPageSize int,
) *AdminHandler {
	return &AdminHandler{
		authClient:      authClient,
		newsClient:      newsClient,
		seoClient:       seoClient,
		logger:          logger,
		defaultPageSize: defaultPageSize,
		maxPageSize:     maxPageSize,
	}
}

// RegisterRoutes регистрирует роуты
func (h *AdminHandler) RegisterRoutes(r *gin.Engine, authMiddleware *middleware.AuthMiddleware) {
	// Публичные endpoints
	r.GET("/health", h.HealthCheck)
	r.POST("/api/admin/login", h.Login)

	// Защищенные endpoints - только admin и editor
	api := r.Group("/api/admin")
	api.Use(authMiddleware.Authenticate())
	api.Use(authMiddleware.RequireRole("admin", "editor"))
	{
		// Dashboard
		api.GET("/dashboard", h.GetDashboard)

		// News CRUD
		api.GET("/news", h.GetNews)
		api.GET("/news/:id", h.GetNewsByID)
		api.POST("/news", h.CreateNews)
		api.PUT("/news/:id", h.UpdateNews)
		api.DELETE("/news/:id", h.DeleteNews)

		// SEO Management
		api.GET("/news/:id/seo", h.GetNewsSEO)
		api.PUT("/news/:id/seo", h.UpdateNewsSEO)

		// User info
		api.GET("/me", h.GetCurrentUser)
	}

	// Admin only endpoints
	adminOnly := r.Group("/api/admin")
	adminOnly.Use(authMiddleware.Authenticate())
	adminOnly.Use(authMiddleware.RequireRole("admin"))
	{
		adminOnly.GET("/users", h.GetUsers)
		adminOnly.DELETE("/news/:id/force", h.ForceDeleteNews)
	}
}

// HealthCheck проверка работоспособности
func (h *AdminHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "admin-service",
		"version": "1.0.0",
	})
}

// Login авторизация в админ панели
func (h *AdminHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Авторизуемся через auth-service
	loginResp, err := h.authClient.Login(req.Username, req.Password)
	if err != nil {
		h.logger.Error("Login failed", zap.Error(err), zap.String("username", req.Username))
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid credentials",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Проверяем роль (только admin и editor)
	if loginResp.User.Role != "admin" && loginResp.User.Role != "editor" {
		h.logger.Warn("Access denied for non-admin user",
			zap.String("username", req.Username),
			zap.String("role", loginResp.User.Role),
		)
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "Access denied. Admin or editor role required.",
			Code:    http.StatusForbidden,
		})
		return
	}

	h.logger.Info("User logged in", zap.String("username", req.Username), zap.String("role", loginResp.User.Role))

	c.JSON(http.StatusOK, loginResp)
}

// GetDashboard получает статистику для дашборда
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	stats, err := h.newsClient.GetDashboardStats()
	if err != nil {
		h.logger.Error("Failed to get dashboard stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to load dashboard",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetNews получает список новостей с пагинацией и фильтрами
func (h *AdminHandler) GetNews(c *gin.Context) {
	params := h.parsePaginationParams(c)

	news, err := h.newsClient.GetNews(params)
	if err != nil {
		h.logger.Error("Failed to get news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to load news",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, news)
}

// GetNewsByID получает новость по ID
func (h *AdminHandler) GetNewsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid news ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	news, err := h.newsClient.GetNewsByID(id)
	if err != nil {
		h.logger.Error("Failed to get news by ID", zap.Error(err), zap.String("id", id.String()))
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "News not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, news)
}

// CreateNews создает новость
func (h *AdminHandler) CreateNews(c *gin.Context) {
	var req models.CreateNewsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid create news request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	token, _ := middleware.GetToken(c)

	news, err := h.newsClient.CreateNews(req, token)
	if err != nil {
		h.logger.Error("Failed to create news", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create news",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	h.logger.Info("News created", zap.String("id", news.ID.String()), zap.String("title", news.Title))

	c.JSON(http.StatusCreated, news)
}

// UpdateNews обновляет новость
func (h *AdminHandler) UpdateNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid news ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateNewsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid update news request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	token, _ := middleware.GetToken(c)

	news, err := h.newsClient.UpdateNews(id, req, token)
	if err != nil {
		h.logger.Error("Failed to update news", zap.Error(err), zap.String("id", id.String()))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update news",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	h.logger.Info("News updated", zap.String("id", id.String()))

	c.JSON(http.StatusOK, news)
}

// DeleteNews удаляет новость
func (h *AdminHandler) DeleteNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid news ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	token, _ := middleware.GetToken(c)

	if err := h.newsClient.DeleteNews(id, token); err != nil {
		h.logger.Error("Failed to delete news", zap.Error(err), zap.String("id", id.String()))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete news",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	h.logger.Info("News deleted", zap.String("id", id.String()))

	c.Status(http.StatusNoContent)
}

// GetNewsSEO получает SEO метаданные новости
func (h *AdminHandler) GetNewsSEO(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid news ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Сначала получаем новость, чтобы узнать slug
	news, err := h.newsClient.GetNewsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "News not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Получаем SEO по slug
	seo, err := h.seoClient.GetSEOBySlug(news.Slug)
	if err != nil {
		h.logger.Error("Failed to get SEO", zap.Error(err), zap.String("slug", news.Slug))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to load SEO",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if seo == nil {
		// SEO еще не создано
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "SEO metadata not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, seo)
}

// UpdateNewsSEO обновляет SEO метаданные новости
func (h *AdminHandler) UpdateNewsSEO(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid news ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateSEORequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid update SEO request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	seo, err := h.seoClient.UpdateSEO(id, req)
	if err != nil {
		h.logger.Error("Failed to update SEO", zap.Error(err), zap.String("news_id", id.String()))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update SEO",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	h.logger.Info("SEO updated", zap.String("news_id", id.String()))

	c.JSON(http.StatusOK, seo)
}

// GetCurrentUser получает текущего пользователя
func (h *AdminHandler) GetCurrentUser(c *gin.Context) {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers получает список пользователей (только admin)
func (h *AdminHandler) GetUsers(c *gin.Context) {
	// TODO: Implement users list
	c.JSON(http.StatusNotImplemented, models.ErrorResponse{
		Error:   "not_implemented",
		Message: "Users list not implemented yet",
		Code:    http.StatusNotImplemented,
	})
}

// ForceDeleteNews принудительно удаляет новость (только admin)
func (h *AdminHandler) ForceDeleteNews(c *gin.Context) {
	// Используем обычное удаление
	h.DeleteNews(c)
}

// parsePaginationParams парсит параметры пагинации
func (h *AdminHandler) parsePaginationParams(c *gin.Context) models.PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.defaultPageSize)))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = h.defaultPageSize
	}

	if pageSize > h.maxPageSize {
		pageSize = h.maxPageSize
	}

	return models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Search:   c.Query("search"),
		Status:   c.Query("status"),
		Category: c.Query("category"),
		Tag:      c.Query("tag"),
		SortBy:   c.DefaultQuery("sort_by", "created_at"),
		SortDir:  c.DefaultQuery("sort_dir", "desc"),
	}
}
