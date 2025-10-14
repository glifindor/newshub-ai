package models

import (
	"time"

	"github.com/google/uuid"
)

// User модель пользователя (из auth-service)
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // admin, editor, user
	CreatedAt time.Time `json:"created_at"`
}

// News модель новости
type News struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content"`
	AuthorID    uuid.UUID  `json:"author_id"`
	AuthorName  string     `json:"author_name"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Category    string     `json:"category,omitempty"`
	Tags        []string   `json:"tags"`
	ImageURL    string     `json:"image_url"`
	Status      string     `json:"status"` // draft, published, archived
	ViewCount   int        `json:"view_count"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// SEOMeta модель SEO метаданных
type SEOMeta struct {
	ID                uuid.UUID              `json:"id"`
	NewsID            uuid.UUID              `json:"news_id"`
	Slug              string                 `json:"slug"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	Keywords          string                 `json:"keywords"`
	CanonicalURL      string                 `json:"canonical_url"`
	OGTitle           string                 `json:"og_title"`
	OGDescription     string                 `json:"og_description"`
	OGImage           string                 `json:"og_image"`
	RobotsIndex       bool                   `json:"robots_index"`
	RobotsFollow      bool                   `json:"robots_follow"`
	SitemapChangeFreq string                 `json:"sitemap_change_freq"`
	SitemapPriority   float64                `json:"sitemap_priority"`
	SchemaData        map[string]interface{} `json:"schema_data,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// DashboardStats статистика для дашборда
type DashboardStats struct {
	TotalNews     int64             `json:"total_news"`
	PublishedNews int64             `json:"published_news"`
	DraftNews     int64             `json:"draft_news"`
	TotalViews    int64             `json:"total_views"`
	TodayViews    int64             `json:"today_views"`
	PopularNews   []News            `json:"popular_news"`
	RecentNews    []News            `json:"recent_news"`
	CategoryStats map[string]int64  `json:"category_stats"`
	ViewsTrend    []ViewsTrendPoint `json:"views_trend"`
}

// ViewsTrendPoint точка графика просмотров
type ViewsTrendPoint struct {
	Date  string `json:"date"`
	Views int64  `json:"views"`
}

// CreateNewsRequest запрос на создание новости
type CreateNewsRequest struct {
	Title      string     `json:"title" binding:"required,min=3,max=255"`
	Slug       string     `json:"slug" binding:"required,min=3,max=255"`
	Summary    string     `json:"summary" binding:"required,min=10,max=500"`
	Content    string     `json:"content" binding:"required,min=50"`
	CategoryID *uuid.UUID `json:"category_id"`
	Tags       []string   `json:"tags"`
	ImageURL   string     `json:"image_url"`
	Status     string     `json:"status" binding:"required,oneof=draft published"`
}

// UpdateNewsRequest запрос на обновление новости
type UpdateNewsRequest struct {
	Title      string     `json:"title" binding:"min=3,max=255"`
	Summary    string     `json:"summary" binding:"min=10,max=500"`
	Content    string     `json:"content" binding:"min=50"`
	CategoryID *uuid.UUID `json:"category_id"`
	Tags       []string   `json:"tags"`
	ImageURL   string     `json:"image_url"`
	Status     string     `json:"status" binding:"oneof=draft published archived"`
}

// UpdateSEORequest запрос на обновление SEO
type UpdateSEORequest struct {
	Title             string  `json:"title" binding:"max=70"`
	Description       string  `json:"description" binding:"max=160"`
	Keywords          string  `json:"keywords" binding:"max=255"`
	OGTitle           string  `json:"og_title" binding:"max=70"`
	OGDescription     string  `json:"og_description" binding:"max=160"`
	OGImage           string  `json:"og_image"`
	RobotsIndex       bool    `json:"robots_index"`
	RobotsFollow      bool    `json:"robots_follow"`
	SitemapChangeFreq string  `json:"sitemap_change_freq" binding:"oneof=always hourly daily weekly monthly yearly never"`
	SitemapPriority   float64 `json:"sitemap_priority" binding:"min=0,max=1"`
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse ответ при входе
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// PaginationParams параметры пагинации
type PaginationParams struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Search   string `form:"search"`
	Status   string `form:"status" binding:"oneof= draft published archived"`
	Category string `form:"category"`
	Tag      string `form:"tag"`
	SortBy   string `form:"sort_by" binding:"oneof= created_at updated_at views title"`
	SortDir  string `form:"sort_dir" binding:"oneof= asc desc"`
}

// PaginatedResponse ответ с пагинацией
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// ErrorResponse ошибка
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}
