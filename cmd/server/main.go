package main

import (
	"context"
	"goshield/internal/gbl"
	log "goshield/internal/log"
	"goshield/internal/redis"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goshield/internal/config"
	"goshield/internal/conn"
	"goshield/internal/middleware"
	"goshield/internal/routes"

	healthcheck "goshield/internal/service/healthcheck"

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

	initializeSystem()

	ctx, cancel := context.WithCancel(context.Background())

	healthcheck.NewHealthChecker(30 * time.Second).Start(ctx)

	routes.SetAdminRoutes(app)
	routes.ProxyRoutes(app)

	certs := config.GetCertificates()

	go func() {
		log.Logger.Infof("TLS Server starting")
		app.ListenTLS(config.GetTlsProxyAddr(), certs.PublicKey, certs.PrivateKey)
	}()

	go func() {
		if err := app.Listen(config.GetProxyAddr()); err != nil {
			log.Logger.Errorf("Server stopped: %+v\n", err)
		}
	}()

	gracefulShutdown(app, cancel)
}

func initializeSystem() {

	config.LoadConfig("configs/env_config.conf")

	err := conn.InitDbConn()
	if err != nil {
		log.Logger.Errorf("Shutting down application as there is error while connection to database")
		return
	}

	redisConfig := config.GetRedisConfig()
	err = redis.InitRedisConn(redisConfig)

	if err != nil {
		return
	}

	redis.NewRedisService().LoadActiveServers()
}

func gracefulShutdown(app *fiber.App, healthCtx context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	sig := <-sigChan

	log.Logger.Infof("Received shutdown signal: %v", sig)

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	defer healthCtx()

	done := make(chan struct{})

	go func() {
		defer close(done)

		if err := app.Shutdown(); err != nil {
			log.Logger.Errorf("Error during shutdown: %v", err)
		}

		if db := conn.GetDb(); db != nil {
			if err := db.Close(); err != nil {
				log.Logger.Errorf("Error closing DB: %v", err)
			}
		}
	}()

	select {
	case <-done:
		log.Logger.Info("Graceful shutdown completed")

	case <-shutdownCtx.Done():
		log.Logger.Warn("Graceful shutdown timeout reached")
	}
}
