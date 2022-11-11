package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(addr, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	_, err := rdb.Ping(context.TODO()).Result()

	if err != nil {
		return &redis.Client{}, err
	}

	return rdb, nil
}
