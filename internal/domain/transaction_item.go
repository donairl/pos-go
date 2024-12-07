package domain

import (
	"time"

	"gorm.io/gorm"
)

type TransactionItem struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	TransactionID uint           `json:"transaction_id" gorm:"not null"`
	ProductID     uint           `json:"product_id" gorm:"not null"`
	Product       Product        `json:"product" gorm:"foreignKey:ProductID"`
	Quantity      int            `json:"quantity" gorm:"not null"`
	Price         float64        `json:"price" gorm:"not null"` // Price at the time of transaction
	Subtotal      float64        `json:"subtotal" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
