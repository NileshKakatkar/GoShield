package proxyhandler

import (
	"goshield/internal/log"
	routing "goshield/internal/service/routing"
	"net/http"
	"strings"

	"github.com/gofiber/contrib/websocket"
	gws "github.com/gorilla/websocket"
)

func (p *proxyHandler) ForwardWebsocket(clientConn *websocket.Conn, headers http.Header, originalUrl string) {
	serverAddr, err := routing.NewLoadBalancer().GetBackendServer()
	if err != nil {
		log.Logger.Errorf("Error while getting servers data: %+v\n", err)
	}

	clientConn.Conn.UnderlyingConn()
	wsServerAddr := "ws://" + strings.Split(serverAddr, "://")[1]

	backendConn, _, err := gws.DefaultDialer.Dial(
		wsServerAddr,
		headers,
	)

	if err != nil {
		log.Logger.Errorf("backend websocket dial failed: %v", err)
		clientConn.Close()
		return
	}

	defer backendConn.Close()
	defer clientConn.Close()

	errChan := make(chan error, 2)

	go proxyClientToBackend(clientConn, backendConn, errChan)
	go proxyBackendToClient(clientConn, backendConn, errChan)

	<-errChan

}

func proxyClientToBackend(clientConn *websocket.Conn, backendConn *gws.Conn, errChan chan error) {
	for {
		msgType, msg, err := clientConn.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}

		if err := backendConn.WriteMessage(msgType, msg); err != nil {
			errChan <- err
			return
		}
	}
}

func proxyBackendToClient(clientConn *websocket.Conn, backendConn *gws.Conn, errChan chan error) {
	for {
		msgType, msg, err := backendConn.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}

		if err := clientConn.WriteMessage(msgType, msg); err != nil {
			errChan <- err
			return
		}
	}
}
