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

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// NewsClient клиент для news-service
type NewsClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// NewNewsClient создает новый клиент news-service
func NewNewsClient(baseURL string, logger *logger.Logger) *NewsClient {
	return &NewsClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// GetNews получает список новостей
func (c *NewsClient) GetNews(params models.PaginationParams) (*models.PaginatedResponse, error) {
	url := fmt.Sprintf("%s/api/v1/news?page=%d&page_size=%d", c.baseURL, params.Page, params.PageSize)

	if params.Search != "" {
		url += fmt.Sprintf("&search=%s", params.Search)
	}
	if params.Status != "" {
		url += fmt.Sprintf("&status=%s", params.Status)
	}
	if params.Category != "" {
		url += fmt.Sprintf("&category=%s", params.Category)
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		c.logger.Error("Failed to get news", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get news: status %d", resp.StatusCode)
	}

	var result models.PaginatedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetNewsByID получает новость по ID
func (c *NewsClient) GetNewsByID(id uuid.UUID) (*models.News, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v1/news/%s", c.baseURL, id.String()))
	if err != nil {
		c.logger.Error("Failed to get news by ID", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("news not found: status %d", resp.StatusCode)
	}

	var news models.News
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, err
	}

	return &news, nil
}

// CreateNews создает новость
func (c *NewsClient) CreateNews(req models.CreateNewsRequest, token string) (*models.News, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/news", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Failed to create news", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create news: %s", string(bodyBytes))
	}

	var news models.News
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, err
	}

	return &news, nil
}

// UpdateNews обновляет новость
func (c *NewsClient) UpdateNews(id uuid.UUID, req models.UpdateNewsRequest, token string) (*models.News, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/news/%s", c.baseURL, id.String()), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Failed to update news", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update news: %s", string(bodyBytes))
	}

	var news models.News
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, err
	}

	return &news, nil
}

// DeleteNews удаляет новость
func (c *NewsClient) DeleteNews(id uuid.UUID, token string) error {
	httpReq, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/news/%s", c.baseURL, id.String()), nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Failed to delete news", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete news: status %d", resp.StatusCode)
	}

	return nil
}

// GetDashboardStats получает статистику для дашборда
func (c *NewsClient) GetDashboardStats() (*models.DashboardStats, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v1/stats/dashboard", c.baseURL))
	if err != nil {
		c.logger.Error("Failed to get dashboard stats", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get stats: status %d", resp.StatusCode)
	}

	var stats models.DashboardStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
