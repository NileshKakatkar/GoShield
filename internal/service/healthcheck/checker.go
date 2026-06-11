package helthcheck

import (
	"fmt"
	"goshield/internal/db"
	"goshield/internal/model"
	"goshield/internal/redis"
	"net/http"
	"time"
)

func (h *healthChecker) CheckAllServers() {
	servers, _ := redis.NewRedisService().GetAllServers()

	var activeServes []model.ServerInfo

	for _, server := range servers {
		go func() {
			server = h.CheckServer(server)

			if server.Status == "active" {
				activeServes = append(servers, server)
			}
		}()
	}

	redis.NewRedisService().UpdateActiveServers(activeServes)
}

func (h *healthChecker) CheckServer(server model.ServerInfo) model.ServerInfo {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	targetUrl := fmt.Sprintf("%s:%d", server.HostAddress, server.Port)
	resp, err := client.Get(targetUrl + "/health")

	if err != nil {
		db.NewServerRepo().UpdateServerStatus(server.UniqId, "inactive")
		server.Status = "inactive"
		return server
	}

	defer resp.Body.Close()

	healthy := resp.StatusCode == http.StatusOK

	status := "inactive"
	if healthy {
		status = "active"
	}

	db.NewServerRepo().UpdateServerStatus(server.UniqId, status)

	server.Status = status
	return server
}
