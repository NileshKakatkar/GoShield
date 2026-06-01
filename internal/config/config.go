package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"goshield/internal/model"
)

var appConfig model.Config

func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading config file: %+v\n", err)
		return
	}

	if err := json.Unmarshal(data, &appConfig); err != nil {
		log.Printf("Error unmarshalling config: %+v\n", err)
		return
	}
}
func GetBackendAddr() string {
	return fmt.Sprintf(
		"%s://%s:%d",
		appConfig.Backend.Protocol,
		appConfig.Backend.Host,
		appConfig.Backend.Port,
	)
}
func GetProxyAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		appConfig.Proxy.Host,
		appConfig.Proxy.Port,
	)
}
