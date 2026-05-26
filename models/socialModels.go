package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	PostID uint   `gorm:"not null"`
	UserID uint   `gorm:"not null"`
	User   User   `gorm:"foreignKey:UserID"`
	Body   string `gorm:"not null"`
}

type Tag struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Posts []Post `gorm:"many2many:post_tags;"`
}

type Like struct {
	UserID    uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"primaryKey"`
	CreatedAt time.Time
}

type Bookmark struct {
	UserID    uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"primaryKey"`
	CreatedAt time.Time
}

type UserActivityLog struct {
	gorm.Model
	UserID        uint   `gorm:"not null"`
	PostID        uint   `gorm:"not null"`
	ActionType    string `gorm:"not null"` // view, click, share
	WatchDuration int    // in seconds
}

type UserInterest struct {
	UserID        uint      `gorm:"primaryKey"`
	TagID         uint      `gorm:"primaryKey"`
	InterestScore float64   `gorm:"type:numeric(5,2);not null"`
	UpdatedAt     time.Time
}
