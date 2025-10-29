package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"size:100;not null" json:"-" validate:"required,min=6"`
	FullName  string         `gorm:"size:100" json:"full_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Articles  []Article      `gorm:"foreignKey:AuthorID" json:"articles,omitempty"`
}
