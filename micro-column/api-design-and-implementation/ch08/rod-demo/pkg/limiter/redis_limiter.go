package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

// LimitDefinition 定义限流规则
type LimitDefinition struct {
	Rate   int           // 周期内允许的次数
	Period time.Duration // 周期
	Burst  int           // 允许的突发量
}

type Limiter struct {
	core *redis_rate.Limiter
}

func NewLimiter(rdb *redis.Client) *Limiter {
	return &Limiter{
		core: redis_rate.NewLimiter(rdb),
	}
}

// Allow 检查是否允许通过
// key: 限流的标识（如 IP、UserID）
// limit: 限流规则
func (l *Limiter) Allow(ctx context.Context, key string, limit LimitDefinition) (*redis_rate.Result, error) {
	// redis_rate 的 Limit 结构体转换
	rrLimit := redis_rate.Limit{
		Rate:   limit.Rate,
		Period: limit.Period,
		Burst:  limit.Burst,
	}

	return l.core.Allow(ctx, key, rrLimit)
}
