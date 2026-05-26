package usecase

import (
	"errors"

	"github.com/EYOB123695/roha/domain"
)

type postUseCase struct {
	postRepo domain.PostRepository
}

func NewPostUseCase(postRepo domain.PostRepository) domain.PostUseCase {
	return &postUseCase{postRepo: postRepo}
}

func (u *postUseCase) CreatePost(userID uint, mediaURL, mediaType, caption string) (*domain.Post, error) {
	post := &domain.Post{
		UserID:    userID,
		MediaURL:  mediaURL,
		MediaType: mediaType,
		Caption:   caption,
	}

	err := u.postRepo.Create(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (u *postUseCase) GetPosts() ([]domain.Post, error) {
	return u.postRepo.GetAll()
}

func (u *postUseCase) GetPost(id uint) (*domain.Post, error) {
	post, err := u.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (u *postUseCase) UpdatePost(userID, postID uint, caption string) (*domain.Post, error) {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	if post.UserID != userID {
		return nil, errors.New("unauthorized to update this post")
	}

	post.Caption = caption
	err = u.postRepo.Update(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (u *postUseCase) DeletePost(userID, postID uint) error {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post not found")
	}

	if post.UserID != userID {
		return errors.New("unauthorized to delete this post")
	}

	return u.postRepo.Delete(postID)
}
