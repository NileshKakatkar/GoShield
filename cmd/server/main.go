package main

import (
	"log"
	"time"

	"goshield/internal/config"
	"goshield/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func init() {
	go middleware.CleanupClients()
}

func main() {
	app := fiber.New(fiber.Config{})
	app.Use(middleware.RateLimiter(100, 1*time.Minute))
	app.Use(middleware.JWTAuth())
	config.LoadServerConfig("configs/env_config.conf")
	log.Fatal(app.Listen(":8000"))
}
