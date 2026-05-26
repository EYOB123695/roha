package model

import "gorm.io/gorm"

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
