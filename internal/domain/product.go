package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"not null"`
	CategoryID uint           `json:"category_id" gorm:"not null"`
	Category   Category       `json:"category" gorm:"foreignKey:CategoryID"`
	Price      float64        `json:"price" gorm:"not null"`
	Stock      int            `json:"stock" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
