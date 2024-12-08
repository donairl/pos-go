package handler

import (
	"strconv"

	"pos-go/internal/domain"
	"pos-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService  service.ProductService
	categoryService service.CategoryService
}

func NewProductHandler(productService service.ProductService, categoryService service.CategoryService) *ProductHandler {
	return &ProductHandler{
		productService:  productService,
		categoryService: categoryService,
	}
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	stockBelow, _ := strconv.Atoi(c.Query("stock_below", "0"))

	// Create context with stock filter if provided
	ctx := c.Context()
	if stockBelow > 0 {
		ctx.SetUserValue("stock_below", stockBelow)
	}

	products, total, err := h.productService.GetProducts(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// If it's an API request, return JSON
	if c.Get("Accept") == "application/json" || c.Get("HX-Request") == "true" {
		return c.JSON(fiber.Map{
			"data": products,
			"meta": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	}

	// Get all categories for the dropdown
	categories, _, err := h.categoryService.GetCategories(1, 100)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Calculate total pages
	totalPages := (int(total) + limit - 1) / limit

	return c.Render("products/index", fiber.Map{
		"Title":      "Products",
		"Products":   products,
		"Categories": categories,
		"Meta": fiber.Map{
			"Page":       page,
			"Limit":      limit,
			"Total":      int(total),
			"TotalPages": totalPages,
		},
	}, "layouts/main")
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": product,
	})
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product domain.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.productService.CreateProduct(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": product,
	})
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	product.ID = uint(id)
	if err := h.productService.UpdateProduct(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": product,
	})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
