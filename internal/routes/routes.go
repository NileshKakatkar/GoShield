package routes

import (
	"goshield/internal/handler/proxyhandler"

	"github.com/gofiber/fiber/v2"
)

type SetupProxyRoutes interface {
	ProxyRoutes(c *fiber.App)
}
type setup struct{}

func NewRoutes() SetupProxyRoutes {
	return &setup{}
}

func (s *setup) ProxyRoutes(app *fiber.App) {
	app.Use("/*", func(c *fiber.Ctx) error {
		return proxyhandler.NewProxyHandler().ForwardRequest(c)
	})
}
