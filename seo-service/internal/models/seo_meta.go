package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SEOMeta модель для хранения SEO метаданных новостей
type SEOMeta struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	NewsID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"news_id"`

	// Basic SEO Fields
	Title        string `gorm:"size:70;not null" json:"title"`             // SEO title (макс 60-70 символов)
	Description  string `gorm:"size:160;not null" json:"description"`      // Meta description (макс 160 символов)
	Keywords     string `gorm:"size:255" json:"keywords,omitempty"`        // Keywords через запятую
	Slug         string `gorm:"size:500;not null;uniqueIndex" json:"slug"` // URL slug
	CanonicalURL string `gorm:"size:500" json:"canonical_url,omitempty"`   // Canonical URL

	// Open Graph (Facebook, LinkedIn, etc.)
	OGTitle       string `gorm:"size:100" json:"og_title,omitempty"`       // OG title
	OGDescription string `gorm:"size:200" json:"og_description,omitempty"` // OG description
	OGImage       string `gorm:"size:500" json:"og_image,omitempty"`       // OG image URL
	OGType        string `gorm:"size:50;default:'article'" json:"og_type"` // article, website, video
	OGLocale      string `gorm:"size:10;default:'en_US'" json:"og_locale"` // en_US, ru_RU

	// Twitter Cards
	TwitterCard        string `gorm:"size:50;default:'summary_large_image'" json:"twitter_card"` // summary, summary_large_image
	TwitterTitle       string `gorm:"size:100" json:"twitter_title,omitempty"`
	TwitterDescription string `gorm:"size:200" json:"twitter_description,omitempty"`
	TwitterImage       string `gorm:"size:500" json:"twitter_image,omitempty"`
	TwitterCreator     string `gorm:"size:50" json:"twitter_creator,omitempty"` // @username

	// Schema.org (Structured Data)
	SchemaType string         `gorm:"size:50;default:'NewsArticle'" json:"schema_type"` // NewsArticle, BlogPosting, Article
	SchemaData datatypes.JSON `gorm:"type:jsonb" json:"schema_data,omitempty"`          // Полные данные JSON-LD

	// Indexing & Sitemap
	RobotsIndex       bool    `gorm:"default:true" json:"robots_index"`                      // Разрешить индексацию
	RobotsFollow      bool    `gorm:"default:true" json:"robots_follow"`                     // Разрешить переход по ссылкам
	SitemapPriority   float32 `gorm:"type:decimal(2,1);default:0.5" json:"sitemap_priority"` // 0.0-1.0
	SitemapChangefreq string  `gorm:"size:20;default:'weekly'" json:"sitemap_changefreq"`    // always, hourly, daily, weekly, monthly, yearly, never

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook для генерации UUID
func (s *SEOMeta) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName задает имя таблицы
func (SEOMeta) TableName() string {
	return "seo_meta"
}

// ValidateTitle проверяет длину title
func (s *SEOMeta) ValidateTitle() error {
	if len(s.Title) > 70 {
		s.Title = s.Title[:67] + "..."
	}
	return nil
}

// ValidateDescription проверяет длину description
func (s *SEOMeta) ValidateDescription() error {
	if len(s.Description) > 160 {
		s.Description = s.Description[:157] + "..."
	}
	return nil
}

// AutoFillOG автозаполнение Open Graph из базовых полей
func (s *SEOMeta) AutoFillOG() {
	if s.OGTitle == "" {
		s.OGTitle = s.Title
	}
	if s.OGDescription == "" {
		s.OGDescription = s.Description
	}
}

// AutoFillTwitter автозаполнение Twitter из базовых полей
func (s *SEOMeta) AutoFillTwitter() {
	if s.TwitterTitle == "" {
		s.TwitterTitle = s.Title
	}
	if s.TwitterDescription == "" {
		s.TwitterDescription = s.Description
	}
	if s.TwitterImage == "" {
		s.TwitterImage = s.OGImage
	}
}
