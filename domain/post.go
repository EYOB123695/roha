package domain

import "time"

type Post struct {
	ID        uint
	UserID    uint
	User      User
	MediaURL  string
	MediaType string
	Caption   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostRepository interface {
	Create(post *Post) error
	GetAll() ([]Post, error)
	GetByID(id uint) (*Post, error)
	Update(post *Post) error
	Delete(id uint) error
}

type PostUseCase interface {
	CreatePost(userID uint, mediaURL, mediaType, caption string) (*Post, error)
	GetPosts() ([]Post, error)
	GetPost(id uint) (*Post, error)
	UpdatePost(userID, postID uint, caption string) (*Post, error)
	DeletePost(userID, postID uint) error
}
