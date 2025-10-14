package models

import "github.com/google/uuid"

// CreateSEORequest DTO для создания SEO метаданных
type CreateSEORequest struct {
	NewsID      uuid.UUID `json:"news_id" binding:"required"`
	Title       string    `json:"title" binding:"required,max=70"`
	Description string    `json:"description" binding:"required,max=160"`
	Keywords    string    `json:"keywords" binding:"omitempty,max=255"`
	Slug        string    `json:"slug" binding:"required,max=500"`

	// Optional fields
	CanonicalURL      string  `json:"canonical_url" binding:"omitempty,url"`
	OGTitle           string  `json:"og_title" binding:"omitempty,max=100"`
	OGImage           string  `json:"og_image" binding:"omitempty,url"`
	TwitterCard       string  `json:"twitter_card" binding:"omitempty,oneof=summary summary_large_image"`
	TwitterCreator    string  `json:"twitter_creator" binding:"omitempty,max=50"`
	RobotsIndex       *bool   `json:"robots_index"`
	RobotsFollow      *bool   `json:"robots_follow"`
	SitemapPriority   float32 `json:"sitemap_priority" binding:"omitempty,gte=0,lte=1"`
	SitemapChangefreq string  `json:"sitemap_changefreq" binding:"omitempty,oneof=always hourly daily weekly monthly yearly never"`
}

// UpdateSEORequest DTO для обновления SEO метаданных
type UpdateSEORequest struct {
	Title             *string  `json:"title" binding:"omitempty,max=70"`
	Description       *string  `json:"description" binding:"omitempty,max=160"`
	Keywords          *string  `json:"keywords" binding:"omitempty,max=255"`
	OGTitle           *string  `json:"og_title" binding:"omitempty,max=100"`
	OGDescription     *string  `json:"og_description" binding:"omitempty,max=200"`
	OGImage           *string  `json:"og_image" binding:"omitempty,url"`
	TwitterCard       *string  `json:"twitter_card" binding:"omitempty,oneof=summary summary_large_image"`
	TwitterCreator    *string  `json:"twitter_creator" binding:"omitempty,max=50"`
	RobotsIndex       *bool    `json:"robots_index"`
	RobotsFollow      *bool    `json:"robots_follow"`
	SitemapPriority   *float32 `json:"sitemap_priority" binding:"omitempty,gte=0,lte=1"`
	SitemapChangefreq *string  `json:"sitemap_changefreq" binding:"omitempty,oneof=always hourly daily weekly monthly yearly never"`
}

// NewsPublishedWebhook DTO для webhook от news-service
type NewsPublishedWebhook struct {
	NewsID        uuid.UUID `json:"news_id" binding:"required"`
	Title         string    `json:"title" binding:"required"`
	Slug          string    `json:"slug" binding:"required"`
	Summary       string    `json:"summary"`
	Content       string    `json:"content"`
	FeaturedImage string    `json:"featured_image"`
	CategoryName  string    `json:"category_name"`
	Tags          []string  `json:"tags"`
	AuthorName    string    `json:"author_name"`
	PublishedAt   string    `json:"published_at"`
}
