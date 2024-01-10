package redis

import (
	"github.com/redis/go-redis/v9"
	"tiktok/src/constant/config"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.EnvCfg.RedisAddr,
		Password: config.EnvCfg.RedisPassword,
		DB:       config.EnvCfg.RedisDB,
	})
}
