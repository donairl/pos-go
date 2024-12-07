package repository

import (
	"pos-go/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	GetUsers(page, limit int) ([]domain.User, int64, error)
	Create(user *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUsers(page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&domain.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}
