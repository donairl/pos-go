package handler

import (
	"pos-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// ShowLogin handles displaying the login page
func (h *AuthHandler) ShowLogin(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Login",
	})
}

// Login handles the login POST request
func (h *AuthHandler) Login(c *fiber.Ctx) error {
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

	// Set JWT token as cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	c.Cookie(cookie)

	// Return token in response for localStorage
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"message": "Login successful",
	})
}

// Logout handles the logout request
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("jwt")
	return c.Redirect("/auth/login")
}
