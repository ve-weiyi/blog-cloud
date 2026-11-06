package svc

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/config"
)

// ConnectRedis 建立 Redis 连接并做一次连通性探测。
func ConnectRedis(c config.RedisConf) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Host + ":" + c.Port,
		Password: c.Password,
		DB:       c.DB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("initialize: ping redis: %w", err)
	}

	return client, nil
}
