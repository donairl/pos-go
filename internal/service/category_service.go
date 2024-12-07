package service

import (
	"pos-go/internal/domain"
	"pos-go/internal/repository"
)

type CategoryService interface {
	GetCategories(page, limit int) ([]domain.Category, int64, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) error
	DeleteCategory(id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetCategories(page, limit int) ([]domain.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.categoryRepo.GetCategories(page, limit)
}

func (s *categoryService) GetCategoryByID(id uint) (*domain.Category, error) {
	return s.categoryRepo.GetCategoryByID(id)
}

func (s *categoryService) CreateCategory(category *domain.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *categoryService) UpdateCategory(category *domain.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *categoryService) DeleteCategory(id uint) error {
	return s.categoryRepo.Delete(id)
}
