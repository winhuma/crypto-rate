package myconnect

import (
	"time"

	"github.com/go-redis/redis"
)

var REDIS *redis.Client

func RedisConnect(params ...string) *redis.Client {
	REDIS = redis.NewClient(&redis.Options{
		Addr: params[0],
		DB:   0,
	})
	err := REDIS.Set("key", "value", time.Second*5).Err()
	if err != nil {
		panic(err)
	}

	_, err = REDIS.Get("key").Result()
	if err != nil {
		panic(err)
	}
	return REDIS
}

func RedisGetInstance() *redis.Client {
	return REDIS
}
