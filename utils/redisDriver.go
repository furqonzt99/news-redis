package utils

import (
	config "github.com/furqonzt99/news-redis/configs"
	"github.com/go-redis/redis/v8"
)

func InitRedis(config *config.AppConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       0,
	})

	return rdb
}
