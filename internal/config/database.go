package config

import (
	"fmt"
	"log"
	"os"
	"pos-go/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("Connected Successfully to Database")

	// Auto Migrate the models
	DB.AutoMigrate(&domain.User{}, &domain.Category{}, &domain.Product{}, &domain.Transaction{}, &domain.TransactionItem{})

	// Seed the database
	seedDatabase()
}

func seedDatabase() {
	// Seed admin user if it doesn't exist
	var adminUser domain.User
	if err := DB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		adminUser = domain.User{
			Username: "admin",
			Password: "$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa", // admin4321
			Role:     "admin",
		}
		DB.Create(&adminUser)
	}

	// Seed default categories if they don't exist
	categories := []domain.Category{
		{Name: "Electronics"},
		{Name: "Food & Beverages"},
		{Name: "Clothing"},
		{Name: "Books"},
		{Name: "Others"},
	}

	for _, category := range categories {
		var existingCategory domain.Category
		if err := DB.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			DB.Create(&category)
		}
	}

	// Seed test products if they don't exist
	var productCount int64
	DB.Model(&domain.Product{}).Count(&productCount)
	if productCount == 0 {
		// Get the "Electronics" category
		var electronicsCategory domain.Category
		DB.Where("name = ?", "Electronics").First(&electronicsCategory)
		var foodCategory domain.Category
		DB.Where("name = ?", "Food & Beverages").First(&foodCategory)

		products := []domain.Product{
			{
				Name:       "Laptop",
				Price:      999.99,
				Stock:      10,
				CategoryID: electronicsCategory.ID,
			},
			{
				Name:       "Smartphone",
				Price:      499.99,
				Stock:      20,
				CategoryID: electronicsCategory.ID,
			},
			{
				Name:       "Tablet",
				Price:      299.99,
				Stock:      15,
				CategoryID: electronicsCategory.ID,
			},
			{
				Name:       "Orange Juice",
				Price:      2.29,
				Stock:      50,
				CategoryID: foodCategory.ID,
			},
			{
				Name:       "Mineral Water ABC",
				Price:      1.99,
				Stock:      100,
				CategoryID: foodCategory.ID,
			},
		}

		for _, product := range products {
			DB.Create(&product)
		}
	}
}
