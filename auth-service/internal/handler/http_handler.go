package handler

import (
	"net/http"

	"auth-service/internal/models"
	"auth-service/internal/service"
	"auth-service/pkg/logger"
	"auth-service/pkg/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewHTTPHandler(authService *service.AuthService, userService *service.UserService) *HTTPHandler {
	return &HTTPHandler{
		authService: authService,
		userService: userService,
	}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/register [post]
func (h *HTTPHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Валидация
	if err := validator.ValidateStruct(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		FullName:     req.FullName,
		Role:         req.Role,
	}

	if err := h.authService.Register(c.Request.Context(), user); err != nil {
		logger.Error("Failed to register user", zap.Error(err), zap.String("email", req.Email))

		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	logger.Info("User registered successfully", zap.String("user_id", user.ID), zap.String("email", user.Email))

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"user": models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Login godoc
// @Summary Вход в систему
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *HTTPHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Валидация
	if err := validator.ValidateStruct(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		logger.Warn("Failed login attempt", zap.Error(err), zap.String("email", req.Email))

		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	logger.Info("User logged in successfully", zap.String("email", req.Email))

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    int(tokens.AccessExpiresIn.Seconds()),
		"token_type":    "Bearer",
	})
}

// RefreshToken godoc
// @Summary Обновление токенов
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func (h *HTTPHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.Warn("Failed to refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	logger.Info("Tokens refreshed successfully")

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    int(tokens.AccessExpiresIn.Seconds()),
		"token_type":    "Bearer",
	})
}

// Logout godoc
// @Summary Выход из системы
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *HTTPHandler) Logout(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		logger.Error("Failed to logout", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	logger.Info("User logged out successfully", zap.String("user_id", userID.(string)))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// GetProfile godoc
// @Summary Получить профиль текущего пользователя
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/profile [get]
func (h *HTTPHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		logger.Error("Failed to get profile", zap.Error(err), zap.String("user_id", userID.(string)))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
}

// UpdateProfile godoc
// @Summary Обновить профиль
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Данные для обновления"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/profile [put]
func (h *HTTPHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	var updates struct {
		FullName string `json:"full_name,omitempty"`
		Email    string `json:"email,omitempty" validate:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if updates.FullName != "" {
		user.FullName = updates.FullName
	}

	if updates.Email != "" {
		user.Email = updates.Email
	}

	if err := h.userService.UpdateProfile(c.Request.Context(), user); err != nil {
		logger.Error("Failed to update profile", zap.Error(err), zap.String("user_id", userID.(string)))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	logger.Info("Profile updated successfully", zap.String("user_id", user.ID))

	c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
}
