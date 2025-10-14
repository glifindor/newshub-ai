package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NewsStatus статус новости
type NewsStatus string

const (
	StatusDraft     NewsStatus = "draft"
	StatusPublished NewsStatus = "published"
	StatusArchived  NewsStatus = "archived"
)

// News модель новости
type News struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Title           string     `gorm:"size:500;not null;index" json:"title"`
	Slug            string     `gorm:"size:500;not null;uniqueIndex" json:"slug"`
	Summary         string     `gorm:"type:text" json:"summary"`
	Content         string     `gorm:"type:text;not null" json:"content"`
	FeaturedImage   string     `gorm:"size:500" json:"featured_image"`
	AuthorID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"author_id"`
	CategoryID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"category_id"`
	Category        *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags            []Tag      `gorm:"many2many:news_tags;" json:"tags,omitempty"`
	Status          NewsStatus `gorm:"type:varchar(20);default:'draft';index" json:"status"`
	PublishedAt     *time.Time `gorm:"index" json:"published_at,omitempty"`
	ViewCount       int64      `gorm:"default:0" json:"view_count"`
	IsFeatured      bool       `gorm:"default:false;index" json:"is_featured"`
	IsBreaking      bool       `gorm:"default:false;index" json:"is_breaking"`
	AllowComments   bool       `gorm:"default:true" json:"allow_comments"`
	MetaTitle       string     `gorm:"size:255" json:"meta_title"`
	MetaDescription string     `gorm:"size:500" json:"meta_description"`
	MetaKeywords    string     `gorm:"size:500" json:"meta_keywords"`

	// SEO Optimization Fields
	SEOTitle       string `gorm:"size:70" json:"seo_title"`
	SEODescription string `gorm:"size:160" json:"seo_description"`
	SEOKeywords    string `gorm:"size:255" json:"seo_keywords"`
	SEOImage       string `gorm:"size:500" json:"seo_image"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook для генерации UUID
func (n *News) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// TableName задает имя таблицы
func (News) TableName() string {
	return "news"
}

// NewsRequest DTO для создания/обновления новости
type NewsRequest struct {
	Title           string      `json:"title" validate:"required,min=5,max=500"`
	Slug            string      `json:"slug" validate:"omitempty,min=5,max=500"`
	Summary         string      `json:"summary" validate:"omitempty,max=5000"`
	Content         string      `json:"content" validate:"required,min=50"`
	FeaturedImage   string      `json:"featured_image" validate:"omitempty,url"`
	CategoryID      uuid.UUID   `json:"category_id" validate:"required,uuid"`
	TagIDs          []uuid.UUID `json:"tag_ids" validate:"omitempty"`
	Status          NewsStatus  `json:"status" validate:"omitempty,oneof=draft published archived"`
	PublishedAt     *time.Time  `json:"published_at"`
	IsFeatured      bool        `json:"is_featured"`
	IsBreaking      bool        `json:"is_breaking"`
	AllowComments   bool        `json:"allow_comments"`
	MetaTitle       string      `json:"meta_title" validate:"omitempty,max=255"`
	MetaDescription string      `json:"meta_description" validate:"omitempty,max=500"`
	MetaKeywords    string      `json:"meta_keywords" validate:"omitempty,max=500"`

	// SEO Fields
	SEOTitle       string `json:"seo_title" validate:"omitempty,max=70"`
	SEODescription string `json:"seo_description" validate:"omitempty,max=160"`
	SEOKeywords    string `json:"seo_keywords" validate:"omitempty,max=255"`
	SEOImage       string `json:"seo_image" validate:"omitempty,url"`
}

// NewsResponse DTO для ответа
type NewsResponse struct {
	ID              uuid.UUID         `json:"id"`
	Title           string            `json:"title"`
	Slug            string            `json:"slug"`
	Summary         string            `json:"summary,omitempty"`
	Content         string            `json:"content"`
	FeaturedImage   string            `json:"featured_image,omitempty"`
	AuthorID        uuid.UUID         `json:"author_id"`
	CategoryID      uuid.UUID         `json:"category_id"`
	Category        *CategoryResponse `json:"category,omitempty"`
	Tags            []TagResponse     `json:"tags,omitempty"`
	Status          NewsStatus        `json:"status"`
	PublishedAt     *time.Time        `json:"published_at,omitempty"`
	ViewCount       int64             `json:"view_count"`
	IsFeatured      bool              `json:"is_featured"`
	IsBreaking      bool              `json:"is_breaking"`
	AllowComments   bool              `json:"allow_comments"`
	MetaTitle       string            `json:"meta_title,omitempty"`
	MetaDescription string            `json:"meta_description,omitempty"`
	MetaKeywords    string            `json:"meta_keywords,omitempty"`

	// SEO Fields
	SEOTitle       string `json:"seo_title,omitempty"`
	SEODescription string `json:"seo_description,omitempty"`
	SEOKeywords    string `json:"seo_keywords,omitempty"`
	SEOImage       string `json:"seo_image,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewsListResponse DTO для списка новостей с пагинацией
type NewsListResponse struct {
	Items      []NewsResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}
