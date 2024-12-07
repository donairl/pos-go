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
