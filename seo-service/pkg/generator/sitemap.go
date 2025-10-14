package generator

import (
	"encoding/xml"
	"fmt"
	"time"

	"seo-service/internal/models"
)

// SitemapGenerator генерирует sitemap.xml
type SitemapGenerator struct {
	baseURL string
}

// NewSitemapGenerator создает новый генератор sitemap
func NewSitemapGenerator(baseURL string) *SitemapGenerator {
	return &SitemapGenerator{baseURL: baseURL}
}

// GenerateURLSet генерирует основной sitemap.xml из SEO метаданных
func (g *SitemapGenerator) GenerateURLSet(metaList []*models.SEOMeta) (string, error) {
	urlSet := &models.Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]models.SitemapURL, 0, len(metaList)),
	}

	for _, meta := range metaList {
		// Пропускаем если не нужно индексировать
		if !meta.RobotsIndex {
			continue
		}

		url := models.SitemapURL{
			Loc:        fmt.Sprintf("%s/news/%s", g.baseURL, meta.Slug),
			LastMod:    meta.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: meta.SitemapChangeFreq,
			Priority:   meta.SitemapPriority,
		}

		urlSet.URLs = append(urlSet.URLs, url)
	}

	// Маршалинг в XML
	output, err := xml.MarshalIndent(urlSet, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal sitemap: %w", err)
	}

	return xml.Header + string(output), nil
}

// GenerateIndex генерирует sitemap index для больших сайтов
func (g *SitemapGenerator) GenerateIndex(sitemaps []string) (string, error) {
	index := &models.SitemapIndex{
		XMLNS:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: make([]models.SitemapFile, 0, len(sitemaps)),
	}

	now := time.Now().Format("2006-01-02")

	for _, sitemapPath := range sitemaps {
		file := models.SitemapFile{
			Loc:     fmt.Sprintf("%s/%s", g.baseURL, sitemapPath),
			LastMod: now,
		}
		index.Sitemaps = append(index.Sitemaps, file)
	}

	output, err := xml.MarshalIndent(index, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal sitemap index: %w", err)
	}

	return xml.Header + string(output), nil
}

// GenerateStaticPages добавляет статические страницы в sitemap
func (g *SitemapGenerator) GenerateStaticPages() []models.SitemapURL {
	return []models.SitemapURL{
		{
			Loc:        fmt.Sprintf("%s/", g.baseURL),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "daily",
			Priority:   1.0,
		},
		{
			Loc:        fmt.Sprintf("%s/news", g.baseURL),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "hourly",
			Priority:   0.9,
		},
		{
			Loc:        fmt.Sprintf("%s/about", g.baseURL),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   0.5,
		},
		{
			Loc:        fmt.Sprintf("%s/contacts", g.baseURL),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   0.4,
		},
	}
}
