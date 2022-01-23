package common

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func Redis() {
	client := redis.NewClient(&redis.Options{
		Addr : viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB: 0,	// 使用默认 DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	RedisClient = client
}
