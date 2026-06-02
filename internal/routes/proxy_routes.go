package routes

import (
	"goshield/internal/handler/proxyhandler"

	"github.com/gofiber/fiber/v2"
)

func ProxyRoutes(app *fiber.App) {
	app.Use("/*", func(c *fiber.Ctx) error {
		return proxyhandler.NewProxyHandler().ForwardRequest(c)
	})
}
