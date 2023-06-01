package util

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

func NewRedisInstance() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     8,
		MinIdleConns: 5,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
