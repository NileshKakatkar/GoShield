package main

import (
	"log"

	"goshield/internal/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{})
	config.LoadServerConfig("configs/env_config.conf")
	log.Fatal(app.Listen(":8000"))
}
