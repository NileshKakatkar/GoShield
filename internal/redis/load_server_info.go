package redis

import (
	"encoding/json"
	"goshield/internal/db"
	"goshield/internal/model"
)

type RedisServive interface {
	LoadActiveServers() error
	GetAllActiveServers() ([]model.ServerInfo, error)
}

type redisService struct{}

func NewRedisService() RedisServive {
	return &redisService{}
}

func (r *redisService) LoadActiveServers() error {
	servers, err := db.NewServerRepo().GetAllActiveServers()
	if err != nil {
		return err
	}

	data, _ := json.Marshal(servers)

	return redisClient.Set(
		ctx,
		"active_servers",
		data,
		0,
	).Err()
}

func (r *redisService) GetAllActiveServers() ([]model.ServerInfo, error) {
	val, err := redisClient.Get(
		ctx,
		"active_servers",
	).Result()

	if err != nil {
		return nil, err
	}

	var servers []model.ServerInfo

	json.Unmarshal(
		[]byte(val),
		&servers,
	)

	return servers, nil
}
