package routes

import (
	"goshield/internal/handler/serverhandler"

	"github.com/gofiber/fiber/v2"
)

func SetAdminRoutes(app *fiber.App) {
	admin := app.Group("/admin")

	admin.Get("/server/list", func(c *fiber.Ctx) error {
		return serverhandler.NewServerHandler().GetServersList(c)
	})

	admin.Post("/register/server", func(c *fiber.Ctx) error {
		return serverhandler.NewServerHandler().AddServerData(c)
	})

	admin.Get("/server/details", func(c *fiber.Ctx) error {
		return serverhandler.NewServerHandler().GetServerDetails(c)
	})

	admin.Post("/delete/server", func(c *fiber.Ctx) error {
		return serverhandler.NewServerHandler().RemoveServer(c)
	})

	admin.Post("/login", func(c *fiber.Ctx) error {
		return serverhandler.NewAdminMgmt().Authenticate(c)
	})
}
