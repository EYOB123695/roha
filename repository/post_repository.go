package repository

import (
	"errors"

	"github.com/EYOB123695/roha/domain"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) domain.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(p *domain.Post) error {
	gormPost := Post{
		UserID:    p.UserID,
		MediaURL:  p.MediaURL,
		MediaType: p.MediaType,
		Caption:   p.Caption,
	}

	result := r.db.Create(&gormPost)
	if result.Error != nil {
		return result.Error
	}

	p.ID = gormPost.ID
	p.CreatedAt = gormPost.CreatedAt
	p.UpdatedAt = gormPost.UpdatedAt
	return nil
}

func (r *postRepository) GetAll() ([]domain.Post, error) {
	var gormPosts []Post
	result := r.db.Preload("User").Find(&gormPosts)
	if result.Error != nil {
		return nil, result.Error
	}

	var posts []domain.Post
	for _, gp := range gormPosts {
		posts = append(posts, domain.Post{
			ID:        gp.ID,
			UserID:    gp.UserID,
			MediaURL:  gp.MediaURL,
			MediaType: gp.MediaType,
			Caption:   gp.Caption,
			CreatedAt: gp.CreatedAt,
			UpdatedAt: gp.UpdatedAt,
			User: domain.User{
				ID:        gp.User.ID,
				Username:  gp.User.Username,
				Email:     gp.User.Email,
				AvatarURL: gp.User.AvatarURL,
			},
		})
	}
	return posts, nil
}

func (r *postRepository) GetByID(id uint) (*domain.Post, error) {
	var gp Post
	result := r.db.Preload("User").First(&gp, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &domain.Post{
		ID:        gp.ID,
		UserID:    gp.UserID,
		MediaURL:  gp.MediaURL,
		MediaType: gp.MediaType,
		Caption:   gp.Caption,
		CreatedAt: gp.CreatedAt,
		UpdatedAt: gp.UpdatedAt,
		User: domain.User{
			ID:        gp.User.ID,
			Username:  gp.User.Username,
			Email:     gp.User.Email,
			AvatarURL: gp.User.AvatarURL,
		},
	}, nil
}

func (r *postRepository) Update(p *domain.Post) error {
	var gp Post
	result := r.db.First(&gp, p.ID)
	if result.Error != nil {
		return result.Error
	}

	result = r.db.Model(&gp).Updates(Post{
		Caption: p.Caption,
	})
	if result.Error != nil {
		return result.Error
	}

	p.UpdatedAt = gp.UpdatedAt
	return nil
}

func (r *postRepository) Delete(id uint) error {
	result := r.db.Delete(&Post{}, id)
	return result.Error
}


func (r *postRepository) GetFeed(userID uint) ([]domain.Post, error) { 
	var followingIDs [] uint
	err:= r.db.Table("followers").Where("follower_id = ?", userID).Pluck("following_id", &followingIDs).Error
	if err ! nil { 
		return nil,err
	}

	if len(folloeingIDs) == 0 { 
		return []domain.Post{} ,nil
	}
    var gormPosts []Post
    err = r.db.Preload("User").Where("user_id IN ?", followingIDs).Order("created_at desc").
		Find(&gormPosts).Error
	if err != nil {
		return nil, err
	}


	var posts []domain.Post
	for _, gp := range gormPosts {
		posts = append(posts, domain.Post{
			ID: gp.ID,
			UserID: gp.UserID,
			MediaURL: gp.MediaURL,
			MediaType: gp.MediaType,
			Caption: gp.Caption,
			CreatedAt: gp.CreatedAt,
			UpdatedAt: gp.UpdatedAt,
			User: domain.User{
				ID: gp.User.ID,
				Username: gp.User.Username,
				Email: gp.User.Email,
				AvatarURL: gp.User.AvatarURL,
			},
		})
	}
	return posts, nil



}