package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kilimanjaro-api/config"
	"time"
)

var (
	REDIS *redis.Client
	Ctx   = context.Background()
)

func getCacheUri(cfg *config.Config) string {
	var cacheHost = cfg.RedisHost
	var cachePort = cfg.RedisPort
	return fmt.Sprintf("%s:%s", cacheHost, cachePort)
}

func connectRedisCache() {
	cfg := config.GetConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     getCacheUri(cfg),
		Password: cfg.RedisPass,
		DB:       0,
	})
	if _, redisErr := redisClient.Ping(Ctx).Result(); redisErr != nil {
		fmt.Println(redisErr.Error())
		panic("Error: Unable to connect to Redis")
	}
	REDIS = redisClient
	fmt.Println("Redis cache init was completed")
}

func SetInCache(c *redis.Client, key string, value interface{}, expiration time.Duration) bool {
	marshalledValue, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to set element in cache")
		return false
	}
	c.Set(Ctx, key, marshalledValue, expiration)
	return true
}

func GetFromCache(c *redis.Client, key string) ([]byte, error) {
	value, err := c.Get(Ctx, key).Bytes()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return value, nil
}

func DeleteFromCache(c *redis.Client, key string) {
	c.Del(Ctx, key)
}

func InitialMigration() {
	connectRedisCache()
}
