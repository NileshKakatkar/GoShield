package redis

import (
	"context"
	"fmt"
	"goshield/internal/log"
	"goshield/internal/model"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var client *redis.Client

func InitRedisConn(config model.RedisConfig) error {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client := redis.NewClient(
		&redis.Options{
			Addr: addr,
		})

	err := client.Ping(ctx)

	if err != nil {
		log.Logger.Errorf("Not able to connect to redis server")
		return redis.ErrClosed
	}

	log.Logger.Infof("Connected to redis server")
	return nil
}
