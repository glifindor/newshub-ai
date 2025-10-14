package service

import (
	"errors"
	"time"

	"auth-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewTokenService(secret string, accessExpiry, refreshExpiry time.Duration) *TokenService {
	return &TokenService{
		secret:        secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (s *TokenService) GenerateTokenPair(user *models.User) (*models.TokenPair, error) {
	// Access Token
	accessToken, err := s.generateToken(user, s.accessExpiry, "access")
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshToken, err := s.generateToken(user, s.refreshExpiry, "refresh")
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  s.accessExpiry,
		RefreshExpiresIn: s.refreshExpiry,
	}, nil
}

func (s *TokenService) generateToken(user *models.User, expiry time.Duration, tokenType string) (string, error) {
	now := time.Now()
	claims := &models.TokenClaims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Permissions: s.getRolePermissions(user.Role),
		TokenType:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *TokenService) ValidateAccessToken(tokenString string) (*models.TokenClaims, error) {
	return s.validateToken(tokenString, "access")
}

func (s *TokenService) ValidateRefreshToken(tokenString string) (*models.TokenClaims, error) {
	return s.validateToken(tokenString, "refresh")
}

func (s *TokenService) validateToken(tokenString, expectedType string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func (s *TokenService) getRolePermissions(role string) []string {
	permissions := map[string][]string{
		"admin": {
			"create_news", "edit_news", "delete_news",
			"moderate", "manage_users", "manage_categories",
		},
		"editor": {
			"create_news", "edit_news", "manage_categories",
		},
		"moderator": {
			"moderate", "edit_news",
		},
		"user": {
			"read_news", "comment",
		},
	}

	if perms, ok := permissions[role]; ok {
		return perms
	}
	return permissions["user"]
}
