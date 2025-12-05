package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 实际配置应从配置文件读取
	})
	// Ping test
	if err := Client.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect redis: " + err.Error())
	}
}
