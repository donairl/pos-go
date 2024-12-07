package handler

import (
	"pos-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ShowLogin handles displaying the login page
func (h *UserHandler) ShowLogin(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Login",
	})
}

// Login handles the login POST request
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token, err := h.userService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

// Logout handles the logout request
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("jwt")
	return c.Redirect("/auth/login")
}

// GetUsers handles getting all users (admin only)
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	users, total, err := h.userService.GetUsers(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": users,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
