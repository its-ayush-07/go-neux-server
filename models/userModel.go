package models

import (
	"time"
)

// UserAccount model
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id" db:"id"`
	Email     string    `gorm:"not null;unique" json:"email" db:"email"`
	Password  string    `gorm:"not null" json:"password" db:"password"`
	UserName  string    `gorm:"not null;unique" json:"username" db:"username"`
	Token     string    `json:"token" db:"token"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" db:"updated_at"`
}
