package main

import (
	"log"

	"goshield/internal/config"
	"goshield/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{})
	app.Use(middleware.JWTAuth())
	config.LoadServerConfig("configs/env_config.conf")
	log.Fatal(app.Listen(":8000"))
}
