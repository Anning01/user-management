package domain

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:200;not null" json:"title" validate:"required,min=3,max=200"`
	Content   string         `gorm:"type:text;not null" json:"content" validate:"required,min=10"`
	AuthorID  uint           `gorm:"not null" json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
