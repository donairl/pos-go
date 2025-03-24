package service

import (
	"errors"
	"os"
	"pos-go/internal/domain"
	"pos-go/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(username, password string) (string, error)
	GetUsers(page, limit int) ([]domain.User, int64, error)
	CreateUser(user *domain.User) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func (s *userService) GetUsers(page, limit int) ([]domain.User, int64, error) {
	return s.userRepo.GetUsers(page, limit)
}

func (s *userService) CreateUser(user *domain.User) error {
	// Check if username already exists
	existingUser, err := s.userRepo.FindByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Set the hashed password
	user.Password = string(hashedPassword)

	// If role is not set, set default role as "customer"
	if user.Role == "" {
		user.Role = "customer"
	}

	// Create the user
	return s.userRepo.Create(user)
}
