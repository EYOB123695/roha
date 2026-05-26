package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/EYOB123695/roha/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) Signup(username, email, password, avatarURL string) error {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Encrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return errors.New("failed to encrypt password")
	}

	user := &domain.User{
		Username:  username,
		Email:     email,
		Password:  string(hash),
		AvatarURL: avatarURL,
	}

	return u.userRepo.Create(user)
}

func (u *userUseCase) Login(email, password string) (string, error) {
	// Lookup user
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func (u *userUseCase) GetByID(id uint) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}
