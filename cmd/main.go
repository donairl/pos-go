package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup template engine
	engine := html.New("./views", ".html")

	// Create fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Serve static files
	app.Static("/public", "./public")

	// Setup routes
	setupRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App) {
	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("layouts/main", fiber.Map{
			"Title": "POS System",
		})
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Get("/login", handleLogin)
	auth.Post("/login", handleLoginPost)
	auth.Get("/logout", handleLogout)

	// Protected routes
	api := app.Group("/api", middleware.Protected())

	// Product routes
	products := api.Group("/products")
	products.Get("/", handleGetProducts)
	products.Post("/", handleCreateProduct)
	products.Put("/:id", handleUpdateProduct)
	products.Delete("/:id", handleDeleteProduct)

	// Transaction routes
	transactions := api.Group("/transactions")
	transactions.Get("/", handleGetTransactions)
	transactions.Post("/", handleCreateTransaction)
	transactions.Get("/:id", handleGetTransaction)

	// Report routes
	reports := api.Group("/reports")
	reports.Get("/sales", handleSalesReport)
	reports.Get("/export", handleExportReport)
}
