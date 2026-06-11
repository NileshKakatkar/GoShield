package redis

import (
	"encoding/json"
	"goshield/internal/db"
	"goshield/internal/log"
	"goshield/internal/model"
)

type RedisServive interface {
	LoadActiveServers() error
	GetAllActiveServers() ([]model.ServerInfo, error)
	UpdateActiveServers(servers []model.ServerInfo) error
	GetAllServers() ([]model.ServerInfo, error)
	LoadAllServers() error
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

func (r *redisService) LoadAllServers() error {
	servers, err := db.NewServerRepo().GetAllServers()
	if err != nil {
		return err
	}

	data, _ := json.Marshal(servers)

	return redisClient.Set(
		ctx,
		"all_servers",
		data,
		0,
	).Err()
}

func (r *redisService) GetAllServers() ([]model.ServerInfo, error) {
	var servers []model.ServerInfo

	val, err := redisClient.Get(
		ctx,
		"all_servers",
	).Result()

	if err != nil {
		log.Logger.Errorf("Error while getting data from redis: %+v\n", err)
		return servers, err
	}

	err = json.Unmarshal([]byte(val), &servers)

	if err != nil {
		log.Logger.Errorf("Error while converting data: %+v\n", err)
	}

	return servers, err
}

func (r *redisService) UpdateActiveServers(servers []model.ServerInfo) error {
	data, _ := json.Marshal(servers)

	return redisClient.Set(
		ctx,
		"active_servers",
		data,
		0,
	).Err()
}
