package generator

import (
	"fmt"
	"strings"

	"seo-service/internal/models"
)

// RobotsGenerator генерирует robots.txt
type RobotsGenerator struct {
	baseURL string
}

// NewRobotsGenerator создает новый генератор robots.txt
func NewRobotsGenerator(baseURL string) *RobotsGenerator {
	return &RobotsGenerator{baseURL: baseURL}
}

// Generate генерирует robots.txt из конфигурации
func (g *RobotsGenerator) Generate(config *models.RobotsConfig) string {
	var builder strings.Builder

	// Основные правила для поисковиков
	builder.WriteString("# robots.txt для новостного портала\n")
	builder.WriteString("# Оптимизировано для Яндекс, Google, Mail.ru\n\n")

	// User-agent: *
	builder.WriteString("User-agent: *\n")

	// Allowed paths
	for _, path := range config.AllowedPaths {
		builder.WriteString(fmt.Sprintf("Allow: %s\n", path))
	}

	// Disallowed paths
	for _, path := range config.DisallowedPaths {
		builder.WriteString(fmt.Sprintf("Disallow: %s\n", path))
	}

	builder.WriteString("\n")

	// Яндекс специфичные директивы
	builder.WriteString("# Яндекс\n")
	builder.WriteString("User-agent: Yandex\n")
	builder.WriteString("Allow: /\n")
	builder.WriteString(fmt.Sprintf("Crawl-delay: %d\n", config.CrawlDelay))
	builder.WriteString("Clean-param: utm_source&utm_medium&utm_campaign\n\n")

	// Google
	builder.WriteString("# Google\n")
	builder.WriteString("User-agent: Googlebot\n")
	builder.WriteString("Allow: /\n")
	builder.WriteString(fmt.Sprintf("Crawl-delay: %d\n\n", config.CrawlDelay))

	// Mail.ru
	builder.WriteString("# Mail.ru\n")
	builder.WriteString("User-agent: Mail.RU_Bot\n")
	builder.WriteString("Allow: /\n")
	builder.WriteString(fmt.Sprintf("Crawl-delay: %d\n\n", config.CrawlDelay))

	// Sitemap
	builder.WriteString(fmt.Sprintf("Sitemap: %s/sitemap.xml\n", g.baseURL))

	// Host (только для Яндекс)
	builder.WriteString(fmt.Sprintf("Host: %s\n", strings.TrimPrefix(g.baseURL, "http://")))

	return builder.String()
}

// GenerateDefault генерирует robots.txt с настройками по умолчанию
func (g *RobotsGenerator) GenerateDefault() string {
	config := &models.RobotsConfig{
		AllowedPaths: []string{
			"/",
			"/news/",
			"/api/v1/seo/",
		},
		DisallowedPaths: []string{
			"/admin/",
			"/api/v1/auth/",
			"/private/",
			"/tmp/",
			"/*.json$",
			"/*?*page=",
			"/*?*sort=",
		},
		CrawlDelay: 1,
	}

	return g.Generate(config)
}
