package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// HotKeyDetector 热点检测器
type HotKeyDetector struct {
	rdb        *redis.Client
	threshold  int64         // 多少个访问后算热点
	window     time.Duration // 统计时间窗口
	localCache *LocalCache   // 热点数据存到本地缓存
}

// NewHotKeyDetector 创建热点检测器
func NewHotKeyDetector(rdb *redis.Client, threshold int64, window time.Duration, localCache *LocalCache) *HotKeyDetector {
	return &HotKeyDetector{
		rdb:        rdb,
		threshold:  threshold,
		window:     window,
		localCache: localCache,
	}
}

// Record 记录一次访问，返回当前计数和是否达到热点阈值
func (d *HotKeyDetector) Record(ctx context.Context, key string) (int64, bool) {
	counterKey := fmt.Sprintf("hotkey:counter:%s", key)

	count, err := d.rdb.Incr(ctx, counterKey).Result()
	if err != nil {
		return 0, false // Redis 出错，不标记热点
	}

	// 第一次访问时设置过期时间
	if count == 1 {
		d.rdb.Expire(ctx, counterKey, d.window)
	}

	isHot := count >= d.threshold
	return count, isHot
}
