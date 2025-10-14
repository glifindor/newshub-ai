package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"admin-service/internal/models"
	"admin-service/pkg/logger"

	"go.uber.org/zap"
)

// AuthClient клиент для auth-service
type AuthClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// NewAuthClient создает новый клиент auth-service
func NewAuthClient(baseURL string, logger *logger.Logger) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// VerifyToken проверяет токен и возвращает пользователя
func (c *AuthClient) VerifyToken(token string) (*models.User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/auth/verify", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to verify token", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token: status %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Login авторизация пользователя
func (c *AuthClient) Login(username, password string) (*models.LoginResponse, error) {
	reqBody := models.LoginRequest{
		Username: username,
		Password: password,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/auth/login", c.baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		c.logger.Error("Failed to login", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", string(bodyBytes))
	}

	var loginResp models.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, err
	}

	return &loginResp, nil
}
