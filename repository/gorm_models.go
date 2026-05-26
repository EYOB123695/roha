package repository

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	AvatarURL string

	// Self-referencing many-to-many for followers
	Followers []*User `gorm:"many2many:followers;joinForeignKey:following_id;joinReferences:follower_id"`
	Following []*User `gorm:"many2many:followers;joinForeignKey:follower_id;joinReferences:following_id"`
}

type Post struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	MediaURL  string    `gorm:"not null"`
	MediaType string    `gorm:"not null"` // e.g. "image" or "video"
	Caption   string
	Comments  []Comment `gorm:"constraint:OnDelete:CASCADE;"`
	Tags      []Tag     `gorm:"many2many:post_tags;"`
}

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
