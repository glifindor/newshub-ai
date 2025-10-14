package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MediaType тип медиафайла
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypePDF   MediaType = "pdf"
)

// Media модель медиафайла
type Media struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OriginalName string         `gorm:"size:500;not null" json:"original_name"`
	FileName     string         `gorm:"size:500;not null;uniqueIndex" json:"file_name"`
	FilePath     string         `gorm:"size:1000;not null" json:"file_path"`
	FileSize     int64          `gorm:"not null" json:"file_size"`
	MimeType     string         `gorm:"size:100;not null" json:"mime_type"`
	MediaType    MediaType      `gorm:"type:varchar(20);not null;index" json:"media_type"`
	Width        int            `json:"width,omitempty"`
	Height       int            `json:"height,omitempty"`
	Duration     int            `json:"duration,omitempty"` // для видео (секунды)
	ThumbnailURL string         `gorm:"size:1000" json:"thumbnail_url,omitempty"`
	AltText      string         `gorm:"size:500" json:"alt_text"`
	Caption      string         `gorm:"type:text" json:"caption"`
	UploadedBy   uuid.UUID      `gorm:"type:uuid;not null;index" json:"uploaded_by"`
	Folder       string         `gorm:"size:255;index" json:"folder"`
	IsPublic     bool           `gorm:"default:true" json:"is_public"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook для генерации UUID
func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// TableName задает имя таблицы
func (Media) TableName() string {
	return "media"
}

// MediaUploadRequest DTO для загрузки файла
type MediaUploadRequest struct {
	AltText  string `json:"alt_text" validate:"omitempty,max=500"`
	Caption  string `json:"caption" validate:"omitempty,max=5000"`
	Folder   string `json:"folder" validate:"omitempty,max=255"`
	IsPublic bool   `json:"is_public"`
}

// MediaResponse DTO для ответа
type MediaResponse struct {
	ID           uuid.UUID `json:"id"`
	OriginalName string    `json:"original_name"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	MediaType    MediaType `json:"media_type"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	Duration     int       `json:"duration,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	AltText      string    `json:"alt_text"`
	Caption      string    `json:"caption"`
	UploadedBy   uuid.UUID `json:"uploaded_by"`
	Folder       string    `json:"folder"`
	IsPublic     bool      `json:"is_public"`
	URL          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MediaListResponse DTO для списка файлов
type MediaListResponse struct {
	Items      []MediaResponse `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
