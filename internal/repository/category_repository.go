package repository

import (
	"pos-go/internal/domain"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(page, limit int) ([]domain.Category, int64, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetCategories(page, limit int) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64

	query := r.db.Model(&domain.Category{})

	// Get total count before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) GetCategoryByID(id uint) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}
