package repository

import (
	"errors"

	"github.com/EYOB123695/roha/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u *domain.User) error {
	gormUser := User{
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		AvatarURL: u.AvatarURL,
	}

	result := r.db.Create(&gormUser)
	if result.Error != nil {
		return result.Error
	}

	u.ID = gormUser.ID
	u.CreatedAt = gormUser.CreatedAt
	u.UpdatedAt = gormUser.UpdatedAt
	return nil
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var gormUser User
	result := r.db.First(&gormUser, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &domain.User{
		ID:        gormUser.ID,
		Username:  gormUser.Username,
		Email:     gormUser.Email,
		Password:  gormUser.Password,
		AvatarURL: gormUser.AvatarURL,
		CreatedAt: gormUser.CreatedAt,
		UpdatedAt: gormUser.UpdatedAt,
	}, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var gormUser User
	result := r.db.Where("email = ?", email).First(&gormUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &domain.User{
		ID:        gormUser.ID,
		Username:  gormUser.Username,
		Email:     gormUser.Email,
		Password:  gormUser.Password,
		AvatarURL: gormUser.AvatarURL,
		CreatedAt: gormUser.CreatedAt,
		UpdatedAt: gormUser.UpdatedAt,
	}, nil
}

func (r *userRepository) GetProfileByID(id uint) (*domain.UserProfileDTO, error) {
	var gormUser User 
	err := r.db.First(&gormUser, id).Error
	if err != nil { 
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// count followers from the join table 
	var followersCount int64 
	r.db.Table("followers").Where("following_id = ?", id).Count(&followersCount)

	// count following
	var followingCount int64
	r.db.Table("followers").Where("follower_id = ?", id).Count(&followingCount)
    
	var gormPosts []Post
	r.db.Where("user_id = ? ", id).Find(&gormPosts)

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
		})
	}

	return &domain.UserProfileDTO{
		ID:             gormUser.ID,
		Username:       gormUser.Username,
		Email:          gormUser.Email,
		AvatarURL:      gormUser.AvatarURL,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		Posts:          posts,
	}, nil
}


func (r *userRepository) Follow(followerID, followingID uint) error {
	return r.db.Table("followers").Create(map[string]interface{}{
		"follower_id":  followerID,
		"following_id": followingID,
	}).Error
}

func (r *userRepository) UnFollow(followerID, followingID uint) error {
	return r.db.Table("followers").
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(nil).Error
}

func (r *userRepository) IsFollowing(followerID, followingID uint) (bool, error) {
	var count int64
	err := r.db.Table("followers").Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) GetFollowers(userID uint) ([]domain.FollowerDTO, error) {
	// Step 1: Get all follower IDs from the join table
	var followerIDs []uint
	err := r.db.Table("followers").
		Where("following_id = ?", userID).
		Pluck("follower_id", &followerIDs).Error
	if err != nil {
		return nil, err
	}

	// Step 2: If no followers, return empty list (not nil)
	if len(followerIDs) == 0 {
		return []domain.FollowerDTO{}, nil
	}

	// Step 3: Load the actual user records for those IDs
	var gormUsers []User
	err = r.db.Where("id IN ?", followerIDs).Find(&gormUsers).Error
	if err != nil {
		return nil, err
	}

	// Step 4: Convert GORM models to domain DTOs
	var followers []domain.FollowerDTO
	for _, gu := range gormUsers {
		followers = append(followers, domain.FollowerDTO{
			ID:        gu.ID,
			Username:  gu.Username,
			AvatarURL: gu.AvatarURL,
		})
	}
	return followers, nil
}
