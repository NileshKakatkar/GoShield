package proxyhandler

import (
	"bytes"
	routing "goshield/internal/service/routing"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ProxyHandler interface {
	ForwardRequest(c *fiber.Ctx) error
}

type proxyHandler struct{}

func NewProxyHandler() ProxyHandler {
	return &proxyHandler{}
}

func (p *proxyHandler) ForwardRequest(c *fiber.Ctx) error {
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
		req.Header.Set(string(k), string(v))
	})

	setProxyHeaders(req, c)

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

func setProxyHeaders(req *http.Request, c *fiber.Ctx) {
	req.Header.Set("X-Forwarded-For", c.IP())
	req.Header.Set("X-Forwarded-Proto", c.Protocol())
	req.Header.Set("X-Forwarded-Host", c.Hostname())
}
