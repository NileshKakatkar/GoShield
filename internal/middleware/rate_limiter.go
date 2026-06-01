package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Client struct {
	Count     int
	ResetTime time.Time
}

var clients = make(map[string]*Client)
var mu sync.Mutex

func RateLimiter(maxRequest int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userID").(string)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"error": "Invalid user",
				},
			)
		}

		mu.Lock()
		defer mu.Unlock()

		client, exists := clients[userId]

		if !exists {
			clients[userId] = &Client{
				Count:     1,
				ResetTime: time.Now().Add(window),
			}
			return c.Next()
		}

		if time.Now().After(client.ResetTime) {
			client.Count = 1
			client.ResetTime = time.Now().Add(window)
			return c.Next()
		}

		if client.Count >= maxRequest {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				fiber.Map{
					"error": "rate limit exceeds",
				})
		}
		client.Count++

		return c.Next()
	}
}

func CleanupClients() {
	ticker := time.NewTicker(10 * time.Minute)

	for range ticker.C {
		mu.Lock()

		now := time.Now()

		for userID, client := range clients {
			if now.After(client.ResetTime) {
				delete(clients, userID)
			}
		}

		mu.Unlock()
	}
}
