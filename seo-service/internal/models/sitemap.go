package models

import "encoding/xml"

// Sitemap структура для sitemap.xml
type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapURL отдельный URL в sitemap
type SitemapURL struct {
	Loc        string  `xml:"loc"`                  // URL страницы
	LastMod    string  `xml:"lastmod,omitempty"`    // Дата последнего изменения (ISO 8601)
	ChangeFreq string  `xml:"changefreq,omitempty"` // always, hourly, daily, weekly, monthly, yearly, never
	Priority   float32 `xml:"priority,omitempty"`   // 0.0 - 1.0
}

// NewSitemap создает новый sitemap с правильным namespace
func NewSitemap() *Sitemap {
	return &Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]SitemapURL, 0),
	}
}

// AddURL добавляет URL в sitemap
func (s *Sitemap) AddURL(url SitemapURL) {
	s.URLs = append(s.URLs, url)
}

// SitemapIndex для больших сайтов (несколько sitemap файлов)
type SitemapIndex struct {
	XMLName  xml.Name          `xml:"sitemapindex"`
	XMLNS    string            `xml:"xmlns,attr"`
	Sitemaps []SitemapIndexURL `xml:"sitemap"`
}

// SitemapIndexURL ссылка на отдельный sitemap файл
type SitemapIndexURL struct {
	Loc     string `xml:"loc"`               // URL sitemap файла
	LastMod string `xml:"lastmod,omitempty"` // Дата последнего изменения
}

// NewSitemapIndex создает новый sitemap index
func NewSitemapIndex() *SitemapIndex {
	return &SitemapIndex{
		XMLNS:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: make([]SitemapIndexURL, 0),
	}
}
