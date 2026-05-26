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

type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
}

type UserUseCase interface {
	Signup(username, email, password, avatarURL string) error
	Login(email, password string) (string, error)
	GetByID(id uint) (*User, error)
}
