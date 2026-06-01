package config

import (
	"encoding/json"
	"fmt"
	"goshield/internal/model"
	"log"
	"os"
)

var serverAddr string

func LoadServerConfig(path string) {
	var serverConf model.Server
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error while reading config file: %+v\n", err)
		return
	}
	err = json.Unmarshal(data, &serverConf)
	if err != nil {
		log.Printf("Error while deserializing data: %+v\n", err)
		return
	}
	serverAddr = fmt.Sprintf("%s://%s:%d", serverConf.HttpMethod, serverConf.ServerAddress, serverConf.ServerPort)
}

func GetServerAddr() string {
	return serverAddr
}
