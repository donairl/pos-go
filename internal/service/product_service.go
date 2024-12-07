package service

import (
	"pos-go/internal/domain"
	"pos-go/internal/repository"
)

type ProductService interface {
	GetProducts(page, limit int, search string) ([]domain.Product, int64, error)
	GetProductByID(id uint) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	UpdateProduct(product *domain.Product) error
	DeleteProduct(id uint) error
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) GetProducts(page, limit int, search string) ([]domain.Product, int64, error) {
	return s.productRepo.GetProducts(page, limit, search)
}

func (s *productService) GetProductByID(id uint) (*domain.Product, error) {
	return s.productRepo.GetProductByID(id)
}

func (s *productService) CreateProduct(product *domain.Product) error {
	return s.productRepo.Create(product)
}

func (s *productService) UpdateProduct(product *domain.Product) error {
	return s.productRepo.Update(product)
}

func (s *productService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}
