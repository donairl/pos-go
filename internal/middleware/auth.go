package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("jwt")
		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // TODO: Use env variable
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
