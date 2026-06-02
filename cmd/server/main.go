package main

import (
	"goshield/internal/gbl"
	log "goshield/internal/log"
	"time"

	"goshield/internal/config"
	"goshield/internal/conn"
	"goshield/internal/middleware"
	"goshield/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	go middleware.CleanupClients()
}

func main() {
	log.SetupLogger()
	log.Logger.Infof("Application starting")

	app := fiber.New(fiber.Config{})

	app.Use(middleware.JWTAuth())

	app.Use(middleware.RateLimiter(
		100,
		1*time.Minute),
	)

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${protocol} ${host} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: gbl.DateTimeDispFmt,
		Output:     log.GetAccessLogWriter(),
	}))

	config.LoadConfig("configs/env_config.conf")

	err := conn.InitDbConn()
	if err != nil {
		log.Logger.Errorf("Shutting down application as there is error while connection to database")
		return
	}

	routes.NewRoutes().ProxyRoutes(app)

	if err := app.Listen(config.GetProxyAddr()); err != nil {
		log.Logger.Errorf("Error while starting server")
	}
}
