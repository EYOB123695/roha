package domain

import "time"

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}
// UserProfileDTO defines the structured shape of data returned to the client.
type UserProfileDTO struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatar_url"`
	FollowersCount int64  `json:"followers_count"`
	FollowingCount int64  `json:"following_count"`
	Posts          []Post `json:"posts"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetProfileByID(id uint) (*UserProfileDTO, error) 
}

type UserUseCase interface {
	Signup(username, email, password, avatarURL string) error
	Login(email, password string) (string, error)
	GetByID(id uint) (*User, error)
	GetUserProfile(id uint) (*UserProfileDTO, error) 
}

