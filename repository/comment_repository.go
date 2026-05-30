package repository 

import (

	 "github.com/EYOB123695/roha/domain"
	"gorm.io/gorm"

)

type commentRepository struct { 
	db *gorm.DB
}


func NewCommentRepository (db *gorm.DB) domain.CommentRepository {
	return &commentRepository{db: db}
}


func (r*commentRepository) Create( c* domain.Comment) error {
         gormComment := Comment{
		PostID: c.PostID,
		UserID: c.UserID,
		Body:   c.Body,
	}

	// Insert record into DB using GORM
	result := r.db.Create(&gormComment)
	if result.Error != nil {
		return result.Error
	}


	// Fetch user details to populate the associated User in the domain model
	var u User
	if err := r.db.First(&u, c.UserID).Error; err == nil {
		c.User = domain.User{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			AvatarURL: u.AvatarURL,
		}
	}
  
	// Assign the generated database primary key and timestamps back to the domain model
	c.ID = gormComment.ID
	c.CreatedAt = gormComment.CreatedAt
	c.UpdatedAt = gormComment.UpdatedAt
	return nil







}




