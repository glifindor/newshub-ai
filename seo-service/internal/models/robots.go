package models

// RobotsConfig конфигурация для robots.txt
type RobotsConfig struct {
	UserAgent   string   // User-agent (например, "*", "Googlebot", "Yandex")
	Allow       []string // Разрешенные пути
	Disallow    []string // Запрещенные пути
	CrawlDelay  int      // Задержка между запросами (в секундах)
	SitemapURLs []string // URL sitemap файлов
}

// NewDefaultRobotsConfig создает конфигурацию по умолчанию
func NewDefaultRobotsConfig(sitemapURL string) *RobotsConfig {
	return &RobotsConfig{
		UserAgent: "*",
		Allow: []string{
			"/",
		},
		Disallow: []string{
			"/admin/",
			"/api/",
			"/private/",
		},
		CrawlDelay:  0,
		SitemapURLs: []string{sitemapURL},
	}
}
