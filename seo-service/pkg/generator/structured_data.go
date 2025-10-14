package generator

import (
	"encoding/json"
	"fmt"
	"time"

	"seo-service/internal/models"

	"github.com/google/uuid"
)

// StructuredDataGenerator генерирует Schema.org JSON-LD
type StructuredDataGenerator struct {
	baseURL      string
	siteName     string
	publisherURL string
	logoURL      string
}

// NewStructuredDataGenerator создает новый генератор структурированных данных
func NewStructuredDataGenerator(baseURL, siteName, publisherURL, logoURL string) *StructuredDataGenerator {
	return &StructuredDataGenerator{
		baseURL:      baseURL,
		siteName:     siteName,
		publisherURL: publisherURL,
		logoURL:      logoURL,
	}
}

// NewsArticle схема для новостной статьи (Schema.org NewsArticle)
type NewsArticle struct {
	Context          string      `json:"@context"`
	Type             string      `json:"@type"`
	Headline         string      `json:"headline"`
	Description      string      `json:"description,omitempty"`
	Image            []string    `json:"image,omitempty"`
	DatePublished    string      `json:"datePublished"`
	DateModified     string      `json:"dateModified"`
	Author           *Author     `json:"author,omitempty"`
	Publisher        *Publisher  `json:"publisher"`
	MainEntityOfPage *EntityPage `json:"mainEntityOfPage,omitempty"`
}

// Author схема автора
type Author struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Publisher схема издателя
type Publisher struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	Logo *Logo  `json:"logo,omitempty"`
}

// Logo схема логотипа
type Logo struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
}

// EntityPage схема главной страницы сущности
type EntityPage struct {
	Type string `json:"@type"`
	ID   string `json:"@id"`
}

// GenerateNewsArticle генерирует Schema.org NewsArticle для новостной статьи
func (g *StructuredDataGenerator) GenerateNewsArticle(meta *models.SEOMeta, authorName string, publishedAt, updatedAt time.Time) (map[string]interface{}, error) {
	article := &NewsArticle{
		Context:       "https://schema.org",
		Type:          "NewsArticle",
		Headline:      meta.Title,
		Description:   meta.Description,
		DatePublished: publishedAt.Format(time.RFC3339),
		DateModified:  updatedAt.Format(time.RFC3339),
		Publisher: &Publisher{
			Type: "Organization",
			Name: g.siteName,
			Logo: &Logo{
				Type: "ImageObject",
				URL:  g.logoURL,
			},
		},
		MainEntityOfPage: &EntityPage{
			Type: "WebPage",
			ID:   fmt.Sprintf("%s/news/%s", g.baseURL, meta.Slug),
		},
	}

	// Добавляем изображения
	if meta.OGImage != "" {
		article.Image = []string{meta.OGImage}
	}

	// Добавляем автора
	if authorName != "" {
		article.Author = &Author{
			Type: "Person",
			Name: authorName,
		}
	}

	// Конвертируем в map для сохранения в JSONB
	var result map[string]interface{}
	data, err := json.Marshal(article)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal article: %w", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
	}

	return result, nil
}

// GenerateBreadcrumbList генерирует Schema.org BreadcrumbList
func (g *StructuredDataGenerator) GenerateBreadcrumbList(items []BreadcrumbItem) (map[string]interface{}, error) {
	breadcrumb := map[string]interface{}{
		"@context":        "https://schema.org",
		"@type":           "BreadcrumbList",
		"itemListElement": make([]map[string]interface{}, 0, len(items)),
	}

	for i, item := range items {
		element := map[string]interface{}{
			"@type":    "ListItem",
			"position": i + 1,
			"name":     item.Name,
			"item":     item.URL,
		}
		breadcrumb["itemListElement"] = append(breadcrumb["itemListElement"].([]map[string]interface{}), element)
	}

	return breadcrumb, nil
}

// BreadcrumbItem элемент хлебных крошек
type BreadcrumbItem struct {
	Name string
	URL  string
}

// GenerateWebSite генерирует Schema.org WebSite (для главной страницы)
func (g *StructuredDataGenerator) GenerateWebSite() (map[string]interface{}, error) {
	website := map[string]interface{}{
		"@context": "https://schema.org",
		"@type":    "WebSite",
		"name":     g.siteName,
		"url":      g.baseURL,
		"potentialAction": map[string]interface{}{
			"@type":       "SearchAction",
			"target":      fmt.Sprintf("%s/search?q={search_term_string}", g.baseURL),
			"query-input": "required name=search_term_string",
		},
	}

	return website, nil
}

// GenerateOrganization генерирует Schema.org Organization
func (g *StructuredDataGenerator) GenerateOrganization() (map[string]interface{}, error) {
	organization := map[string]interface{}{
		"@context": "https://schema.org",
		"@type":    "Organization",
		"name":     g.siteName,
		"url":      g.baseURL,
		"logo":     g.logoURL,
	}

	return organization, nil
}

// ToJSONLD конвертирует map в JSON-LD строку
func (g *StructuredDataGenerator) ToJSONLD(data map[string]interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON-LD: %w", err)
	}
	return string(bytes), nil
}

// GenerateFromMeta генерирует полный набор Schema.org данных из SEOMeta
func (g *StructuredDataGenerator) GenerateFromMeta(meta *models.SEOMeta, newsID uuid.UUID, authorName string, publishedAt time.Time) (map[string]interface{}, error) {
	// Если схема уже существует в meta.SchemaData, возвращаем её
	if meta.SchemaData != nil && len(meta.SchemaData) > 0 {
		return meta.SchemaData, nil
	}

	// Иначе генерируем новую
	return g.GenerateNewsArticle(meta, authorName, publishedAt, meta.UpdatedAt)
}
