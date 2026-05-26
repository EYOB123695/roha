package model

import "gorm.io/gorm"

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