package config

import (
	"fmt"
	"log"
	"os"
	"pos-go/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schemas
	db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Transaction{},
		&domain.TransactionItem{},
	)

	DB = db
	log.Println("Database connected successfully")

	// Seed default admin user
	seedAdminUser()
	// Seed test products
	seedTestProducts()
}

func seedAdminUser() {
	var user domain.User
	if err := DB.Where("username = ?", "admin").First(&user).Error; err != nil {
		// User doesn't exist, create it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin4321"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		adminUser := domain.User{
			Username: "admin",
			Password: string(hashedPassword),
			Role:     domain.RoleAdmin,
		}

		if err := DB.Create(&adminUser).Error; err != nil {
			log.Fatal("Failed to create admin user:", err)
		}
		log.Println("Default admin user created successfully")
	}
}

func seedTestProducts() {
	var count int64
	DB.Model(&domain.Product{}).Count(&count)
	if count == 0 {
		testProducts := []domain.Product{
			{
				Name:  "Product 1",
				Price: 10000,
				Stock: 100,
			},
			{
				Name:  "Product 2",
				Price: 20000,
				Stock: 50,
			},
			{
				Name:  "Product 3",
				Price: 15000,
				Stock: 75,
			},
		}

		for _, product := range testProducts {
			if err := DB.Create(&product).Error; err != nil {
				log.Printf("Failed to create test product %s: %v", product.Name, err)
			}
		}
		log.Println("Test products created successfully")
	}
}
