package proxyhandler

import (
	"goshield/internal/log"
	"net/http"

	"github.com/gofiber/contrib/websocket"
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

	headers := http.Header{}

	c.Request().Header.VisitAll(func(key, value []byte) {
		headers.Add(string(key), string(value))
	})

	addProxyHeaders(&headers, c)

	if websocket.IsWebSocketUpgrade(c) {
		log.Logger.Infof("Forwarding websocket reqeust")

		originalUrl := c.OriginalURL()

		websocket.New(func(ws *websocket.Conn) {
			p.ForwardWebsocket(ws, headers, originalUrl)
		})
		return nil
	}

	log.Logger.Infof("Forwarding http request")

	return p.ForwardHttp(c, headers)
}

func addProxyHeaders(headers *http.Header, c *fiber.Ctx) {
	headers.Add("X-Forwarded-For", c.IP())
	headers.Add("X-Forwarded-Proto", c.Protocol())
	headers.Add("X-Forwarded-Host", c.Hostname())
}
