package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category модель категории новостей
type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Slug        string         `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	ParentID    *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	MetaTitle   string         `gorm:"size:255" json:"meta_title"`
	MetaDesc    string         `gorm:"size:500" json:"meta_description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook для генерации UUID
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName задает имя таблицы
func (Category) TableName() string {
	return "categories"
}

// CategoryRequest DTO для создания/обновления категории
type CategoryRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=255"`
	Slug        string     `json:"slug" validate:"omitempty,min=2,max=255"`
	Description string     `json:"description" validate:"omitempty,max=5000"`
	ParentID    *uuid.UUID `json:"parent_id" validate:"omitempty,uuid"`
	SortOrder   int        `json:"sort_order" validate:"omitempty,min=0"`
	IsActive    bool       `json:"is_active"`
	MetaTitle   string     `json:"meta_title" validate:"omitempty,max=255"`
	MetaDesc    string     `json:"meta_description" validate:"omitempty,max=500"`
}

// CategoryResponse DTO для ответа
type CategoryResponse struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description,omitempty"`
	ParentID    *uuid.UUID          `json:"parent_id,omitempty"`
	Parent      *CategoryResponse   `json:"parent,omitempty"`
	Children    []*CategoryResponse `json:"children,omitempty"`
	SortOrder   int                 `json:"sort_order"`
	IsActive    bool                `json:"is_active"`
	MetaTitle   string              `json:"meta_title,omitempty"`
	MetaDesc    string              `json:"meta_description,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
