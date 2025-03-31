package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Video struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	VideoPath     string `json:"video_path"`
	ThumbnailPath string `json:"thumbnail_path"`
	UserID        uint   `json:"user_id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Comment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Content   string    `json:"content"`
	VideoID   uint      `json:"video_id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
