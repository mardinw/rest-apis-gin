package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-envconfig"
	"payuoge.com/configs"
)

func Init() (*redis.Client, error) {
	var ctx = context.Background()
	var config configs.AppConfiguration

	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatal(err.Error())
	}

	if config.AppEnv == "production" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Cache.EndPoint, config.Cache.Port),
			Password: fmt.Sprintf("%s", config.Cache.Password),
			DB:       0,
		})
		status := redisClient.Ping(ctx)
		log.Println(status)
		return redisClient, nil

	} else {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Cache.EndPoint, config.Cache.Port),
			Password: "",
			DB:       0,
		})

		status := redisClient.Ping(ctx)
		log.Println(status)
		return redisClient, nil

	}
}
