package db

import (
	"github.com/FlowerWrong/pusher/log"
	"github.com/go-redis/redis"
)

func initClusterRedisClient(redisURLs []string) error {
	var err error
	var nodes []string
	var redisOptions *redis.Options
	for _, redisURL := range redisURLs {
		redisOptions, err = redis.ParseURL(redisURL)
		if err != nil {
			log.Panic(err)
		}
		nodes = append(nodes, redisOptions.Addr)
	}

	universalOptions := &redis.UniversalOptions{
		Addrs:       nodes,
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
