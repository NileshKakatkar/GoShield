package helthcheck

import (
	"context"
	"goshield/internal/model"
	"time"
)

type HealthChecker interface {
	Start(ctx context.Context)
	CheckAllServers()
	CheckServer(server model.ServerInfo)
}

type healthChecker struct {
	interval time.Duration
}

func NewHealthChecker(interval time.Duration) HealthChecker {
	return &healthChecker{interval: interval}
}

func (h *healthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(h.interval)

	go func() {
		defer ticker.Stop()

		for {
			h.CheckAllServers()
			<-ticker.C
		}
	}()
}
