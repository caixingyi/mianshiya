package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

// NewESBreaker 创建一个 Elasticsearch 搜索熔断器
func NewESBreaker(name string) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,                // 半开状态最多允许 3 个请求探测
		Interval:    30 * time.Second, // 统计窗口，30秒后清空失败计数
		Timeout:     10 * time.Second, // 熔断打开 10 秒后进入半开
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			// 当失败率超过 60% 时触发熔断
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}
