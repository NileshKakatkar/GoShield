package service

import (
	"fmt"
	"goshield/internal/db"
	"goshield/internal/model"
	"sync"
)

type LoadBalancer interface {
	GetBackendServer() (string, error)
}

type loadBalancer struct{}

func NewLoadBalancer() LoadBalancer {
	return &loadBalancer{}
}

var (
	current int
	mu      sync.Mutex
)

func (l *loadBalancer) GetBackendServer() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	serversList, _ := db.NewServerRepo().GetAllActiveServers()
	if len(serversList) == 0 {
		return "", fmt.Errorf("No backend server present")
	}
	backend := serversList[current]

	current = (current + 1) % len(serversList)
	addr := getAddress(backend)
	return addr, nil
}

func getAddress(backend model.ServerInfo) string {
	return fmt.Sprintf(
		"%s://%s:%d",
		backend.Method,
		backend.HostAddress,
		backend.Port,
	)
}
