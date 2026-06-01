package middleware

import (
	"goshield/internal/gbl"
	log "goshield/internal/log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "b4c1c8f4e9d5a7b2f3c8d1e6a9b7c4d2f8e1a6c3b9d5e7f2a4c8d1b6e9f3a7c5"

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		switch path {
		case "/login", "/register":
			log.Logger.Infof("Skipping validations for the login request")
			return c.Next()
		}

		authHeader := c.Get("Authorization")

		if authHeader == gbl.EmptyString {
			log.Logger.Errorf("Authentication header misssing. JWT authentication failed")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authorization header missing",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			log.Logger.Errorf("Invalid authorization format")
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

		claims := token.Claims.(jwt.MapClaims)

		userID := claims["sub"]

		if err != nil || !token.Valid {
			log.Logger.Errorf("Authentication failed for: %s\n", userID)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		c.Locals("userID", userID)
		log.Logger.Infof("JWT authentication successful for: %s", userID)
		return c.Next()
	}
}
