package main

import (
	"log"
	"os"
	"pos-go/internal/config"
	"pos-go/internal/handler"
	"pos-go/internal/middleware"
	"pos-go/internal/repository"
	"pos-go/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	config.ConnectDB()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(cors.New())
	app.Static("/public", "./public")

	// Initialize repositories
	userRepo := repository.NewUserRepository(config.DB)
	productRepo := repository.NewProductRepository(config.DB)
	transactionRepo := repository.NewTransactionRepository(config.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	setupRoutes(app, authHandler, productHandler, transactionHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App,
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	transactionHandler *handler.TransactionHandler) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("layouts/main", fiber.Map{
			"Title": "POS System",
		})
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Get("/login", authHandler.ShowLogin)
	auth.Post("/login", authHandler.Login)
	auth.Get("/logout", authHandler.Logout)

	// Protected routes
	api := app.Group("/api", middleware.Protected())

	// Product routes
	products := api.Group("/products")
	products.Get("/", productHandler.GetProducts)
	products.Post("/", productHandler.CreateProduct)
	products.Get("/:id", productHandler.GetProduct)
	products.Put("/:id", productHandler.UpdateProduct)
	products.Delete("/:id", productHandler.DeleteProduct)

	// Transaction routes
	transactions := api.Group("/transactions")
	transactions.Get("/", transactionHandler.GetTransactions)
	transactions.Post("/", transactionHandler.CreateTransaction)
	transactions.Get("/:id", transactionHandler.GetTransaction)

	// Report routes
	reports := api.Group("/reports")
	reports.Get("/sales", transactionHandler.GetSalesReport)
	reports.Get("/export", transactionHandler.ExportSalesReport)
}
