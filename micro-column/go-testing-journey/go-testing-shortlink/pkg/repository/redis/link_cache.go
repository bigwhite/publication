package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisLinkCache struct {
	client *redis.Client
}

func NewRedisLinkCache(client *redis.Client) *RedisLinkCache {
	return &RedisLinkCache{client: client}
}

func (r *RedisLinkCache) key(code string) string {
	return fmt.Sprintf("link:visits:%s", code)
}

func (r *RedisLinkCache) IncrementVisitCount(ctx context.Context, code string) error {
	return r.client.Incr(ctx, r.key(code)).Err()
}

func (r *RedisLinkCache) GetVisitCount(ctx context.Context, code string) (int64, error) {
	count, err := r.client.Get(ctx, r.key(code)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}
