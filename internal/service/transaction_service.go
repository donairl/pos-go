package service

import (
	"encoding/csv"
	"fmt"
	"pos-go/internal/domain"
	"pos-go/internal/repository"
	"strings"
	"time"
)

type TransactionService interface {
	GetTransactions(page, limit int, userID uint) ([]domain.Transaction, int64, error)
	GetTransactionByID(id uint) (*domain.Transaction, error)
	CreateTransaction(transaction *domain.Transaction) error
	GetDailySalesReport() (interface{}, error)
	GetMonthlySalesReport() (interface{}, error)
	ExportSalesReport(startDate, endDate time.Time, format string) ([]byte, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) GetTransactions(page, limit int, userID uint) ([]domain.Transaction, int64, error) {
	return s.transactionRepo.GetTransactions(page, limit, userID)
}

func (s *transactionService) GetTransactionByID(id uint) (*domain.Transaction, error) {
	return s.transactionRepo.GetTransactionByID(id)
}

func (s *transactionService) CreateTransaction(transaction *domain.Transaction) error {
	return s.transactionRepo.Create(transaction)
}

func (s *transactionService) GetDailySalesReport() (interface{}, error) {
	total, count, err := s.transactionRepo.GetDailySales(time.Now())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"date":  time.Now().Format("2006-01-02"),
		"total": total,
		"count": count,
	}, nil
}

func (s *transactionService) GetMonthlySalesReport() (interface{}, error) {
	now := time.Now()
	total, count, err := s.transactionRepo.GetMonthlySales(now.Year(), now.Month())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"month": now.Format("2006-01"),
		"total": total,
		"count": count,
	}, nil
}

func (s *transactionService) ExportSalesReport(startDate, endDate time.Time, format string) ([]byte, error) {
	transactions, err := s.transactionRepo.GetSalesByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Write header
	writer.Write([]string{"Date", "Transaction ID", "Total", "Payment Method"})

	// Write data
	for _, t := range transactions {
		writer.Write([]string{
			t.TransactionDate.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%d", t.ID),
			fmt.Sprintf("%.2f", t.Total),
			t.PaymentMethod,
		})
	}

	writer.Flush()
	return []byte(builder.String()), nil
}
