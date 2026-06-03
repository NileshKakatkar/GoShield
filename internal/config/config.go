package config

import (
	"encoding/json"
	"fmt"
	log "goshield/internal/log"
	"os"
	"path/filepath"
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

func GetProxyAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		appConfig.Proxy.Host,
		appConfig.Proxy.Port,
	)
}

func GetTlsProxyAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		appConfig.Proxy.Host,
		appConfig.Proxy.TlsPort,
	)
}

func GetDbConfig() model.DbConfig {
	return appConfig.DbConfig
}

func GetCertificates() model.Certificates {
	basePath, _ := os.Getwd()

	return model.Certificates{
		PublicKey:  filepath.Join(basePath, appConfig.Certificates.PublicKey),
		PrivateKey: filepath.Join(basePath, appConfig.Certificates.PrivateKey),
	}
}
