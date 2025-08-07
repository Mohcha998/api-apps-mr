package middleware

import (
	"apps/utils/jwt"
	"apps/utils/response"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth() fiber.Handler {
	return JWT()
}

func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Errors:  "No token provided",
			})
		}

		payload, err := jwt.ValidateToken(token)
		if err != nil || payload == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Errors:  "Your token is not valid!",
			})
		}

		c.Locals("customerID", payload["id"])
		c.Locals("customerEmail", payload["email"])
		return c.Next()
	}
}
