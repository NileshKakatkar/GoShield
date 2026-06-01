package middleware

import (
	"goshield/internal/gbl"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "b4c1c8f4e9d5a7b2f3c8d1e6a9b7c4d2f8e1a6c3b9d5e7f2a4c8d1b6e9f3a7c5"

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == gbl.EmptyString {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authorization header missing",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(
					fiber.StatusUnauthorized,
					"unexpected signing method",
				)
			}

			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		return c.Next()
	}
}
