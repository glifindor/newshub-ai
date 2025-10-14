package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken      string        `json:"access_token"`
	RefreshToken     string        `json:"refresh_token"`
	AccessExpiresIn  time.Duration `json:"access_expires_in"`
	RefreshExpiresIn time.Duration `json:"refresh_expires_in"`
}

type TokenClaims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}
