package repository

import (
	"pos-go/internal/domain"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProducts(page, limit int, search string) ([]domain.Product, int64, error)
	GetProductByID(id uint) (*domain.Product, error)
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProducts(page, limit int, search string) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.Model(&domain.Product{}).Preload("Category")

	// Apply search filter if provided
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// Apply stock filter if provided
	if stockBelow := r.db.Statement.Context.Value("stock_below"); stockBelow != nil {
		if limit, ok := stockBelow.(int); ok {
			query = query.Where("stock < ?", limit)
		}
	}

	// Get total count before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) GetProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}
