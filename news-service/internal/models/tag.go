package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag модель тега
type Tag struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug      string         `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook для генерации UUID
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TableName задает имя таблицы
func (Tag) TableName() string {
	return "tags"
}

// TagRequest DTO для создания/обновления тега
type TagRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Slug string `json:"slug" validate:"omitempty,min=2,max=100"`
}

// TagResponse DTO для ответа
type TagResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
