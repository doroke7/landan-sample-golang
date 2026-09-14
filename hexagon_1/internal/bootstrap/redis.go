package bootstrap

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	sAddr := fmt.Sprintf("%s:%s", CONFIG.REDIS.HOST, CONFIG.REDIS.PORT)

	oOptions := &redis.Options{
		Addr:     sAddr,
		Username: CONFIG.REDIS.USERNAME,
		Password: CONFIG.REDIS.PASSWORD,
		DB:       CONFIG.REDIS.DB,
	}
	oClient := redis.NewClient(oOptions)

	oContext := context.Background()
	oPingResult := oClient.Ping(oContext)
	oErr := oPingResult.Err()
	if oErr != nil {
		return nil, oErr
	}

	return oClient, nil
}
