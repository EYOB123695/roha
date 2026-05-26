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
