package db

import (
	"context"
	"fmt"
	"mianshiya-go-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化 Redis 客户端并测试连接，如果连接成功返回 Redis 客户端实例，否则返回错误
func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
