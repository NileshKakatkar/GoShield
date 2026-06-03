package helthcheck

import (
	"fmt"
	"goshield/internal/db"
	"goshield/internal/model"
	"net/http"
	"time"
)

func (h *healthChecker) CheckAllServers() {
	servers, _ := db.NewServerRepo().GetAllServers()

	for _, server := range servers {
		go h.CheckServer(server)
	}
}

func (h *healthChecker) CheckServer(server model.ServerInfo) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	targetUrl := fmt.Sprintf("%s:%d", server.HostAddress, server.Port)
	resp, err := client.Get(targetUrl + "/health")

	if err != nil {
		db.NewServerRepo().UpdateServerStatus(server.UniqId, "inactive")
		return
	}

	defer resp.Body.Close()

	healthy := resp.StatusCode == http.StatusOK

	status := "inactive"
	if healthy {
		status = "active"
	}

	db.NewServerRepo().UpdateServerStatus(server.UniqId, status)
}
