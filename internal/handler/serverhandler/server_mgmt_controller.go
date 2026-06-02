package serverhandler

import (
	"goshield/internal/db"
	"goshield/internal/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ServerHandler interface {
	GetServersList(c *fiber.Ctx) error
	AddServerData(c *fiber.Ctx) error
	GetServerDetails(c *fiber.Ctx) error
}

type serverHandler struct{}

func NewServerHandler() ServerHandler {
	return &serverHandler{}
}

func (s *serverHandler) GetServersList(c *fiber.Ctx) error {
	servers, err := db.NewServerRepo().GetAllServers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"response": model.ServerResp{
					Status:  fiber.StatusInternalServerError,
					Message: "request failed",
				},
			})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"response": model.ServerResp{
				ServersList: servers,
				Message:     "request successful",
			},
		})
}

func (s *serverHandler) AddServerData(c *fiber.Ctx) error {
	var server model.ServerInfo

	if err := c.BodyParser(&server); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"response": model.ServerResp{
					Status:  fiber.StatusBadRequest,
					Message: "Invalid request format",
				},
			})
	}

	results := db.NewServerRepo().AddServer(server)

	if results.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"response": model.ServerResp{
					Status:  fiber.StatusInternalServerError,
					Message: "request failed",
				},
			})
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"response": model.ServerResp{
				Message: "Server configuration added",
			},
		})
}

func (s *serverHandler) GetServerDetails(c *fiber.Ctx) error {
	serverId := c.FormValue("serverId")
	uniqId, _ := strconv.Atoi(serverId)
	if uniqId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"response": model.ServerResp{
					Status:  fiber.StatusBadRequest,
					Message: "Invalid request",
				},
			})
	}
	server, _ := db.NewServerRepo().GetServerByUnqiqId(uniqId)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"response": model.ServerResp{
				Server:  *server,
				Message: "Server configuration added",
			},
		})
}
