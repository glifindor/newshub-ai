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

// SEOClient клиент для seo-service
type SEOClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// NewSEOClient создает новый клиент seo-service
func NewSEOClient(baseURL string, logger *logger.Logger) *SEOClient {
	return &SEOClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// GetSEOBySlug получает SEO метаданные по slug
func (c *SEOClient) GetSEOBySlug(slug string) (*models.SEOMeta, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v1/seo/%s", c.baseURL, slug))
	if err != nil {
		c.logger.Error("Failed to get SEO", zap.Error(err), zap.String("slug", slug))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // SEO не найдено - это нормально
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get SEO: status %d", resp.StatusCode)
	}

	var seo models.SEOMeta
	if err := json.NewDecoder(resp.Body).Decode(&seo); err != nil {
		return nil, err
	}

	return &seo, nil
}

// UpdateSEO обновляет SEO метаданные
func (c *SEOClient) UpdateSEO(newsID uuid.UUID, req models.UpdateSEORequest) (*models.SEOMeta, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/seo/update", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Failed to update SEO", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update SEO: %s", string(bodyBytes))
	}

	var seo models.SEOMeta
	if err := json.NewDecoder(resp.Body).Decode(&seo); err != nil {
		return nil, err
	}

	return &seo, nil
}
