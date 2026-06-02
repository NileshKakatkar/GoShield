package config

import (
	"encoding/json"
	"fmt"
	log "goshield/internal/log"
	"os"
	"sync"

	"goshield/internal/model"
)

var appConfig model.Config

var (
	current int
	mu      sync.Mutex
)

func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Logger.Infof("Error reading config file: %+v\n", err)
		return
	}

	if err := json.Unmarshal(data, &appConfig); err != nil {
		log.Logger.Infof("Error unmarshalling config: %+v\n", err)
		return
	}
}

func GetNextBackend() string {
	mu.Lock()
	defer mu.Unlock()

	backend := appConfig.Backends[current]

	current = (current + 1) % len(appConfig.Backends)

	return GetBackendAddr(backend)
}

func GetBackendAddr(backend model.BackendConfig) string {
	return fmt.Sprintf(
		"%s://%s:%d",
		backend.Protocol,
		backend.Host,
		backend.Port,
	)
}

func GetProxyAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		appConfig.Proxy.Host,
		appConfig.Proxy.Port,
	)
}

func GetDbConfig() model.DbConfig {
	return appConfig.DbConfig
}
