package models

import (
	"time"
)

// Article model
type Article struct {
	ID        uint      `gorm:"primaryKey" json:"id" db:"id"`
	Title     string    `gorm:"not null" json:"title" db:"title"`
	Content   string    `gorm:"not null" json:"content" db:"content"`
	Author    string    `gorm:"not null" json:"author" db:"author"`
	AuthorID  uint      `gorm:"not null" json:"authorid" db:"authorid"`
	Image     string    `json:"image" db:"image"`
	Likes     uint      `gorm:"default:0" json:"likes" db:"likes"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" db:"updated_at"`
}
