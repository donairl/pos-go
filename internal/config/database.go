package config

import (
	"fmt"
	"log"
	"os"
	"pos-go/internal/domain"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("Connected Successfully to Database")

	// Auto Migrate the models
	DB.AutoMigrate(&domain.User{}, &domain.Category{}, &domain.Product{}, &domain.Transaction{}, &domain.TransactionItem{})

	// Seed the database
	seedDatabase()
}

// GeneratePassword hashes a password using bcrypt
func GeneratePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func seedDatabase() {
	// Seed admin user if it doesn't exist
	var adminUser domain.User
	if err := DB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		adminPassword, err := GeneratePassword("admin4321")
		if err != nil {
			log.Fatal("Failed to generate admin password: ", err)
		}

		adminUser = domain.User{
			Username: "admin",
			Password: adminPassword,
			Role:     "admin",
		}
		DB.Create(&adminUser)
	}

	// Create another user
	userPassword, err := GeneratePassword("newuser4321")
	if err != nil {
		log.Fatal("Failed to generate user password: ", err)
	}

	newUser := domain.User{
		Username: "donny",
		Password: userPassword,
		Role:     "customer",
	}
	DB.Create(&newUser)

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
				Name:       "Apple MacBook Pro M3 - 14 inch",
				Price:      999.99,
				Stock:      10,
				CategoryID: electronicsCategory.ID,
			},
			{
				Name:       "Samsung Brand Smartphone M-50",
				Price:      499.99,
				Stock:      20,
				CategoryID: electronicsCategory.ID,
			},
			{
				Name:       "Amazon Brand Tablet",
				Price:      299.99,
				Stock:      15,
				CategoryID: electronicsCategory.ID,
			},
			{
				Name:       "Generic Orange Juice",
				Price:      2.29,
				Stock:      50,
				CategoryID: foodCategory.ID,
			},
			{
				Name:       "Mineral Water ABC - 500ml",
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
