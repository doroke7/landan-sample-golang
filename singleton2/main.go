package bootstrap

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

// go 喜歡這種 package-level-global 變量當作唯一連結變量（singleton）
var (
	redisClient *redis.Client
	once        sync.Once
)

func InitRedis() *redis.Client {
	once.Do(func() {
		fmt.Println("[INFO] 初始化 Redis 連線...")
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Username: "backend",
			Password: "redis_pass_cli",
			DB:       0,
		})

		ctx := context.Background()
		if _, err := redisClient.Ping(ctx).Result(); err != nil {
			panic("Redis connection failed: " + err.Error())
		}
	})
	return redisClient
}
