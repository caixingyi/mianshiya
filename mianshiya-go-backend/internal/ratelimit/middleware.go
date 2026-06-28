package ratelimit

import (
	"context"
	"fmt"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// fixedWindowScript 是一个 Lua 脚本，用于实现固定窗口计数器的速率限制逻辑。
// 它会在 Redis 中对指定的键进行递增操作，并在第一次递增时设置过期时间，从而实
// 现对请求数量的限制。
var fixedWindowScript = redis.NewScript(`
	local current = redis.call("INCR", KEYS[1])
	if current == 1 then
	    redis.call("EXPIRE", KEYS[1], ARGV[1])
	end
	return current
`)

// FixedWindowMiddleware 固定窗口限流中间件
// name: 限流资源名，比如 "question:list"
// limit: 窗口内允许的最大请求数
// window: 时间窗口，比如 time.Second
func FixedWindowMiddleware(rdb *redis.Client, name string, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 生成 Redis key
		currentWindow := time.Now().UnixNano() / int64(window)
		key := fmt.Sprintf("rate_limit:%s:%d", name, currentWindow)

		// 2. 执行 Lua 脚本
		count, err := fixedWindowScript.Run(
			context.Background(),
			rdb,
			[]string{key},
			int(window.Seconds())+1,
		).Int64()

		// 3. Redis 出错：为了保护系统，直接拒绝请求
		if err != nil {
			c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, "系统繁忙，请稍后再试"))
			c.Abort()
			return
		}

		// 4. 超过限制：返回错误
		if count > limit {
			c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, "请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		// 5. 没超过限制：放行
		c.Next()
	}
}
