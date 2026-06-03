package helthcheck

import (
	"goshield/internal/model"
	"time"
)

type HealthChecker interface {
	Start()
	CheckAllServers()
	CheckServer(server model.ServerInfo)
}

type healthChecker struct {
	interval time.Duration
}

func NewHealthChecker(interval time.Duration) HealthChecker {
	return &healthChecker{interval: interval}
}

func (h *healthChecker) Start() {
	ticker := time.NewTicker(h.interval)

	go func() {
		defer ticker.Stop()

		for {
			h.CheckAllServers()
			<-ticker.C
		}
	}()
}
