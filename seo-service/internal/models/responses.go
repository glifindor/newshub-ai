package models

import (
	"time"

	"github.com/google/uuid"
)

// SEOResponse полный ответ с SEO данными
type SEOResponse struct {
	ID           uuid.UUID        `json:"id"`
	NewsID       uuid.UUID        `json:"news_id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Keywords     string           `json:"keywords,omitempty"`
	Slug         string           `json:"slug"`
	CanonicalURL string           `json:"canonical_url,omitempty"`
	OpenGraph    *OpenGraphData   `json:"open_graph"`
	Twitter      *TwitterCardData `json:"twitter"`
	Schema       interface{}      `json:"schema,omitempty"` // JSON-LD structured data
	Robots       *RobotsData      `json:"robots"`
	Sitemap      *SitemapData     `json:"sitemap"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// OpenGraphData данные Open Graph
type OpenGraphData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image,omitempty"`
	Type        string `json:"type"`
	Locale      string `json:"locale"`
	URL         string `json:"url,omitempty"`
}

// TwitterCardData данные Twitter Card
type TwitterCardData struct {
	Card        string `json:"card"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image,omitempty"`
	Creator     string `json:"creator,omitempty"`
}

// RobotsData данные для robots meta
type RobotsData struct {
	Index  bool `json:"index"`
	Follow bool `json:"follow"`
}

// SitemapData данные для sitemap
type SitemapData struct {
	Priority   float32 `json:"priority"`
	Changefreq string  `json:"changefreq"`
}

// MetaTagsResponse HTML meta теги
type MetaTagsResponse struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Keywords    string            `json:"keywords,omitempty"`
	Canonical   string            `json:"canonical,omitempty"`
	MetaTags    map[string]string `json:"meta_tags"` // Все meta теги в виде ключ-значение
}

// GenerateMetaTags преобразует SEOMeta в набор meta тегов
func (s *SEOResponse) GenerateMetaTags() map[string]string {
	tags := make(map[string]string)

	// Basic meta tags
	tags["title"] = s.Title
	tags["description"] = s.Description
	if s.Keywords != "" {
		tags["keywords"] = s.Keywords
	}
	if s.CanonicalURL != "" {
		tags["canonical"] = s.CanonicalURL
	}

	// Robots
	robotsContent := ""
	if s.Robots.Index {
		robotsContent += "index"
	} else {
		robotsContent += "noindex"
	}
	if s.Robots.Follow {
		robotsContent += ",follow"
	} else {
		robotsContent += ",nofollow"
	}
	tags["robots"] = robotsContent

	// Open Graph
	if s.OpenGraph != nil {
		tags["og:title"] = s.OpenGraph.Title
		tags["og:description"] = s.OpenGraph.Description
		tags["og:type"] = s.OpenGraph.Type
		tags["og:locale"] = s.OpenGraph.Locale
		if s.OpenGraph.Image != "" {
			tags["og:image"] = s.OpenGraph.Image
		}
		if s.OpenGraph.URL != "" {
			tags["og:url"] = s.OpenGraph.URL
		}
	}

	// Twitter
	if s.Twitter != nil {
		tags["twitter:card"] = s.Twitter.Card
		tags["twitter:title"] = s.Twitter.Title
		tags["twitter:description"] = s.Twitter.Description
		if s.Twitter.Image != "" {
			tags["twitter:image"] = s.Twitter.Image
		}
		if s.Twitter.Creator != "" {
			tags["twitter:creator"] = s.Twitter.Creator
		}
	}

	return tags
}
