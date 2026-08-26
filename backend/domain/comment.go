package domain 

import "time"

type Comment struct { 
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentRepository interface {
	Create(comment *Comment) error 
	GetByPostID(postID uint) ([]*Comment, error)
	GetByID(id uint) (*Comment, error)
	Delete(id uint) error
}

type CommentUseCase interface { 
	AddComment(userID uint, postID uint, body string) (*Comment, error) 
	GetCommentsByPostID(postID uint) ([]*Comment, error)
	DeleteComment(userID uint, commentID uint) error
} 
