package proxyhandler

import (
	"bytes"
	routing "goshield/internal/service/routing"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (p *proxyHandler) ForwardHttp(c *fiber.Ctx, headers http.Header) error {
	serverAddr, err := routing.NewLoadBalancer().GetBackendServer()

	if err != nil {
		return c.Status(fiber.StatusBadGateway).
			SendString("Backend unavailable")
	}

	req, err := http.NewRequest(
		c.Method(),
		serverAddr+c.OriginalURL(),
		bytes.NewReader(c.Body()),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString(err.Error())
	}

	c.Request().Header.VisitAll(func(k, v []byte) {
		headers.Add(string(k), string(v))
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).
			SendString("Backend unavailable")
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Append(key, value)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString(err.Error())
	}

	modifiedBody := body

	return c.Status(resp.StatusCode).Send(modifiedBody)
}
