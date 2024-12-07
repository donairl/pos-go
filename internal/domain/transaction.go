package domain

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID              uint              `json:"id" gorm:"primaryKey"`
	UserID          uint              `json:"user_id" gorm:"not null"`
	User            User              `json:"user" gorm:"foreignKey:UserID"`
	Total           float64           `json:"total" gorm:"not null"`
	PaymentAmount   float64           `json:"payment_amount" gorm:"not null"`
	PaymentMethod   string            `json:"payment_method" gorm:"type:varchar(20);not null"`
	TransactionDate time.Time         `json:"transaction_date" gorm:"not null"`
	Items           []TransactionItem `json:"items" gorm:"foreignKey:TransactionID"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"deleted_at" gorm:"index"`
}
