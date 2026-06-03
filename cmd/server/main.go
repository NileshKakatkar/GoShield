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
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func init() {
	go middleware.CleanupClients()
}

func main() {
	log.SetupLogger()
	log.Logger.Infof("Application starting")

	app := fiber.New(fiber.Config{})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

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

	routes.SetAdminRoutes(app)
	routes.ProxyRoutes(app)

	certs := config.GetCertificates()

	go func() {
		log.Logger.Infof("TLS Server starting")
		app.ListenTLS(config.GetTlsProxyAddr(), certs.PublicKey, certs.PrivateKey)
	}()

	if err := app.Listen(config.GetProxyAddr()); err != nil {
		log.Logger.Errorf("Error while starting server")
	}
}
