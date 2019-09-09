package db

import (
	"sync"

	"github.com/FlowerWrong/pusher/log"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	redisClient redis.UniversalClient
	redisOnce   sync.Once
)

func initRedisClient(redisURL string) error {
	redisOptions, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}
	universalOptions := &redis.UniversalOptions{
		Addrs:       []string{redisOptions.Addr},
		DB:          redisOptions.DB,
		Password:    redisOptions.Password,
		PoolSize:    redisOptions.PoolSize,
		PoolTimeout: redisOptions.PoolTimeout,
	}
	redisClient = redis.NewUniversalClient(universalOptions)

	_, err = redisClient.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// Redis return redis client
func Redis() redis.UniversalClient {
	if redisClient == nil {
		redisOnce.Do(func() {
			redisURLs := viper.GetStringSlice("REDIS_URL")
			if len(redisURLs) > 1 {
				err := initClusterRedisClient(redisURLs)
				if err != nil {
					log.Panic(err)
				}
			} else {
				err := initRedisClient(redisURLs[0])
				if err != nil {
					log.Panic(err)
				}
			}
		})
	}
	return redisClient
}
