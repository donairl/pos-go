package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID               uint              `json:"id" gorm:"primaryKey"`
	Name             string            `json:"name" gorm:"not null"`
	Category         string            `json:"category"`
	Price            float64           `json:"price" gorm:"not null"`
	Stock            int               `json:"stock" gorm:"not null"`
	IsDeleted        bool              `json:"is_deleted" gorm:"default:false"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `json:"deleted_at" gorm:"index"`
	TransactionItems []TransactionItem `json:"transaction_items" gorm:"foreignKey:ProductID"`
}
