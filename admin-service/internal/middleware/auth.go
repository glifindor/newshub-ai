package middleware

import (
	"net/http"
	"strings"

	"admin-service/internal/client"
	"admin-service/internal/models"
	"admin-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware middleware для проверки авторизации
type AuthMiddleware struct {
	authClient *client.AuthClient
	logger     *logger.Logger
}

// NewAuthMiddleware создает новый middleware
func NewAuthMiddleware(authClient *client.AuthClient, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
		logger:     logger,
	}
}

// Authenticate проверяет JWT токен
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization header required",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Извлекаем токен
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid authorization header format",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Проверяем токен через auth-service
		user, err := m.authClient.VerifyToken(token)
		if err != nil {
			m.logger.Warn("Token verification failed", zap.Error(err))
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or expired token",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Сохраняем пользователя в контексте
		c.Set("user", user)
		c.Set("token", token)
		c.Next()
	}
}

// RequireRole проверяет роль пользователя
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "User not authenticated",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		user := userVal.(*models.User)

		// Проверяем, есть ли роль пользователя в списке разрешенных
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			m.logger.Warn("Access denied",
				zap.String("user_id", user.ID.String()),
				zap.String("user_role", user.Role),
				zap.Strings("required_roles", roles),
			)
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "forbidden",
				Message: "Insufficient permissions",
				Code:    http.StatusForbidden,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORS middleware
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Проверяем, разрешен ли origin
		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// GetCurrentUser извлекает пользователя из контекста
func GetCurrentUser(c *gin.Context) (*models.User, error) {
	userVal, exists := c.Get("user")
	if !exists {
		return nil, http.ErrNoCookie
	}

	user, ok := userVal.(*models.User)
	if !ok {
		return nil, http.ErrNoCookie
	}

	return user, nil
}

// GetToken извлекает токен из контекста
func GetToken(c *gin.Context) (string, error) {
	tokenVal, exists := c.Get("token")
	if !exists {
		return "", http.ErrNoCookie
	}

	token, ok := tokenVal.(string)
	if !ok {
		return "", http.ErrNoCookie
	}

	return token, nil
}
