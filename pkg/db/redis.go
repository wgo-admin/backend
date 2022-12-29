package db

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisOptions struct {
	Addr     string
	Password string
	DB       int
}

// NewRedis 创建一个redis客户端
func NewRedis(opts *RedisOptions) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})

	ping := rdb.Ping(context.Background())
	if _, err := ping.Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
