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

func initRedisClient() error {
	redisOptions, err := redis.ParseURL(viper.GetString("REDIS_URL"))
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
				err := initClusterRedisClient()
				if err != nil {
					log.Panic(err)
				}
			} else {
				err := initRedisClient()
				if err != nil {
					log.Panic(err)
				}
			}
		})
	}
	return redisClient
}

// AsksKey ...
func AsksKey(symbol string) string {
	return "exchange:" + symbol + ":depth:asks"
}

// BidsKey ...
func BidsKey(symbol string) string {
	return "exchange:" + symbol + ":depth:bids"
}

// DepthKey ...
func DepthKey(symbol string) string {
	return "exchange:" + symbol + ":depth"
}

// OrderBookKey ...
func OrderBookKey(symbol string) string {
	return "exchange:" + symbol + ":order_book"
}
