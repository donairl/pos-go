package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"pos-go/internal/domain"
	"pos-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	productService     service.ProductService
}

func NewTransactionHandler(transactionService service.TransactionService, productService service.ProductService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		productService:     productService,
	}
}

func (h *TransactionHandler) ShowTransactionPage(c *fiber.Ctx) error {
	// Get initial products for the transaction page
	products, _, err := h.productService.GetProducts(1, 100, "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Render("transactions/index", fiber.Map{
		"Title":    "New Transaction",
		"Products": products,
	}, "layouts/main")
}

func (h *TransactionHandler) GetTransactions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	userID := c.Locals("userID").(uint)

	transactions, total, err := h.transactionService.GetTransactions(page, limit, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": transactions,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	transaction, err := h.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": transaction,
	})
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var transaction domain.Transaction
	if err := c.BodyParser(&transaction); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set user ID and transaction date
	transaction.UserID = c.Locals("userID").(uint)
	transaction.TransactionDate = time.Now()

	// Calculate total
	var total float64
	for i, item := range transaction.Items {
		// Get product to verify price and stock
		product, err := h.productService.GetProductByID(item.ProductID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Product not found for item %d", i+1),
			})
		}

		// Set product price at time of transaction
		item.Price = product.Price
		item.Subtotal = product.Price * float64(item.Quantity)
		total += item.Subtotal

		// Update the item in the slice
		transaction.Items[i] = item
	}

	// Set transaction total
	transaction.Total = total

	if err := h.transactionService.CreateTransaction(&transaction); err != nil {
		// Check if it's a stock-related error
		if strings.Contains(err.Error(), "insufficient stock") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": transaction,
	})
}

func (h *TransactionHandler) GetSalesReport(c *fiber.Ctx) error {
	reportType := c.Query("type", "daily")
	var report interface{}
	var err error

	switch reportType {
	case "daily":
		report, err = h.transactionService.GetDailySalesReport()
	case "monthly":
		report, err = h.transactionService.GetMonthlySalesReport()
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid report type",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": report,
	})
}

func (h *TransactionHandler) ExportSalesReport(c *fiber.Ctx) error {
	startDate, err := time.Parse("2006-01-02", c.Query("start_date"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start date",
		})
	}

	endDate, err := time.Parse("2006-01-02", c.Query("end_date"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end date",
		})
	}

	format := c.Query("format", "csv")

	data, err := h.transactionService.ExportSalesReport(startDate, endDate, format)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", "attachment; filename=sales_report."+format)
	return c.Send(data)
}
