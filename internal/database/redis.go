package database

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

var redisCache *cache.Cache

func ConnectRedis(REDIS_URL, REDIS_PASSWORD string, REDIS_DATABASE int) {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     REDIS_URL,
		Password: REDIS_PASSWORD,
		DB:       REDIS_DATABASE,
	})
	redisCache = cache.New(&cache.Options{Redis: redisDB})
}

func SetCache(key string, value interface{}) (error) {
	err := redisCache.Set(&cache.Item{
		Ctx:	context.Background(),
		Key:	key,
		Value:	value,
		TTL:	time.Duration(time.Duration.Minutes(15)),
	})
	return err
}

func GetCache(key string, value interface{}) (exists bool) {
	err := redisCache.Get(context.Background(), key, value)
	return err == nil
}