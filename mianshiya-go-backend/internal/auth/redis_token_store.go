package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const tokenKeyPrefix = "mianshiya:login:token:"

// RedisTokenStore 实现了 TokenStore 接口，使用 Redis 存储 token 和 userID 的映射关系
type RedisTokenStore struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewRedisTokenStore 创建一个新的 RedisTokenStore 实例，连接到 Redis 数据库，并设置 token 的过期时间
func NewRedisTokenStore(rdb *redis.Client, ttl time.Duration) *RedisTokenStore {
	return &RedisTokenStore{
		rdb: rdb,
		ttl: ttl,
	}
}

// tokenKey 生成 Redis 中存储 token 的 key，使用统一的前缀，方便后续管理和查询
func tokenKey(token string) string {
	return tokenKeyPrefix + token
}

// CreateToken 生成一个新的 token，并将 token 和 userID 的映射关系存储在 Redis 中，设置过期时间
func (r *RedisTokenStore) CreateToken(ctx context.Context, userID int64) (string, error) {
	// 生成随机 token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)
	err = r.rdb.Set(ctx, tokenKey(token), strconv.FormatInt(userID, 10), r.ttl).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetUserID 根据 token 从 Redis 中获取对应的 userID，如果 token 不存在或已过期，返回 false
func (r *RedisTokenStore) GetUserID(ctx context.Context, token string) (int64, bool, error) {
	userID, err := r.rdb.Get(ctx, tokenKey(token)).Result()
	if errors.Is(err, redis.Nil) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return 0, false, err
	}
	return userIDInt, true, nil
}

// DeleteToken 从 Redis 中删除 token 和对应的 userID 的映射关系，通常在用户注销时调用
func (r *RedisTokenStore) DeleteToken(ctx context.Context, token string) error {
	return r.rdb.Del(ctx, tokenKey(token)).Err()
}
