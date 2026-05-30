package repository 

import (
	"errors"

	"github.com/EYOB123695/roha/domain"
	"gorm.io/gorm"
)

type commentRepository struct { 
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) domain.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(c *domain.Comment) error {
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

func (r *commentRepository) GetByPostID(postID uint) ([]*domain.Comment, error) {
	var gormComments []Comment
	result := r.db.Preload("User").Where("post_id = ?", postID).Order("created_at desc").Find(&gormComments)
	if result.Error != nil { 
		return nil, result.Error
	} 

	var comments []*domain.Comment
	for _, gc := range gormComments {
		comments = append(comments, &domain.Comment{
			ID:        gc.ID,
			PostID:    gc.PostID,
			UserID:    gc.UserID,
			Body:      gc.Body,
			CreatedAt: gc.CreatedAt,
			UpdatedAt: gc.UpdatedAt,
			User: domain.User{
				ID:        gc.User.ID,
				Username:  gc.User.Username,
				Email:     gc.User.Email,
				AvatarURL: gc.User.AvatarURL,
			},
		})
	}
	return comments, nil
}

func (r *commentRepository) GetByID(id uint) (*domain.Comment, error) {
	var gc Comment
	result := r.db.First(&gc, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &domain.Comment{
		ID:        gc.ID,
		PostID:    gc.PostID,
		UserID:    gc.UserID,
		Body:      gc.Body,
		CreatedAt: gc.CreatedAt,
		UpdatedAt: gc.UpdatedAt,
	}, nil
}

func (r *commentRepository) Delete(id uint) error {
	result := r.db.Delete(&Comment{}, id)
	return result.Error
}
