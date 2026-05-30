package usecase

import (
	"errors"

	"github.com/EYOB123695/roha/domain"
)

type commentUseCase struct {
	commentRepo domain.CommentRepository
	postRepo    domain.PostRepository
}

// NewCommentUseCase constructs a new CommentUseCase with injected repositories.
func NewCommentUseCase(commentRepo domain.CommentRepository, postRepo domain.PostRepository) domain.CommentUseCase {
	return &commentUseCase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (u *commentUseCase) AddComment(userID uint, postID uint, body string) (*domain.Comment, error) {
	// 1. Core Rule: Verify body text is not empty
	if body == "" {
		return nil, errors.New("comment body cannot be empty")
	}

	// 2. Business Integrity: Verify that the parent post exists
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// 3. Construct domain comment entity
	comment := &domain.Comment{
		UserID: userID,
		PostID: postID,
		Body:   body,
	}

	// 4. Persistence call
	err = u.commentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
