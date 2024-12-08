package repository

import (
	"fmt"
	"pos-go/internal/domain"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetTransactions(page, limit int, userID uint) ([]domain.Transaction, int64, error)
	GetTransactionByID(id uint) (*domain.Transaction, error)
	Create(transaction *domain.Transaction) error
	GetDailySales(date time.Time) (float64, int64, error)
	GetMonthlySales(year int, month time.Month) (float64, int64, error)
	GetSalesByDateRange(startDate, endDate time.Time) ([]domain.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) GetTransactions(page, limit int, userID uint) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var total int64

	query := r.db.Model(&domain.Transaction{})
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Preload("Items.Product").Offset(offset).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *transactionRepository) GetTransactionByID(id uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.Preload("Items.Product").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) Create(transaction *domain.Transaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the transaction first
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// Update stock for each item
		for _, item := range transaction.Items {
			// Get current product
			var product domain.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return err
			}

			// Check if enough stock
			if product.Stock < item.Quantity {
				return fmt.Errorf("insufficient stock for product %s: have %d, need %d", product.Name, product.Stock, item.Quantity)
			}

			// Update stock
			if err := tx.Model(&product).Update("stock", product.Stock-item.Quantity).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *transactionRepository) GetDailySales(date time.Time) (float64, int64, error) {
	var total float64
	var count int64

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Model(&domain.Transaction{}).
		Where("transaction_date BETWEEN ? AND ?", startOfDay, endOfDay).
		Select("COALESCE(SUM(total), 0) as total, COUNT(*) as count").
		Row().
		Scan(&total, &count)

	if err != nil {
		return 0, 0, err
	}

	return total, count, nil
}

func (r *transactionRepository) GetMonthlySales(year int, month time.Month) (float64, int64, error) {
	var total float64
	var count int64

	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	err := r.db.Model(&domain.Transaction{}).
		Where("transaction_date BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Select("COALESCE(SUM(total), 0) as total, COUNT(*) as count").
		Row().Scan(&total, &count)

	return total, count, err
}

func (r *transactionRepository) GetSalesByDateRange(startDate, endDate time.Time) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := r.db.Preload("Items.Product").
		Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Find(&transactions).Error

	return transactions, err
}
