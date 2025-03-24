package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First try to get token from Authorization header
		token := c.Get("Authorization")
		if token != "" {
			// Remove 'Bearer ' prefix if present
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			// Fallback to cookie if header not present
			token = c.Cookies("jwt")
		}

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Verify the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !parsedToken.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims := parsedToken.Claims.(jwt.MapClaims)
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])
		c.Locals("userID", uint(claims["user_id"].(float64)))

		return c.Next()
	}
}
