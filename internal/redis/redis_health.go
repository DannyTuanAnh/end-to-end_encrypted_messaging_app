package redis_memory

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var startOnce sync.Once

type RedisHealth struct {
	alive atomic.Bool
}

func NewRedisHealth() *RedisHealth {
	return &RedisHealth{}
}

func (r *RedisHealth) SetAlive(v bool) {
	r.alive.Store(v)
}

func (r *RedisHealth) IsAlive() bool {
	return r.alive.Load()
}

// Hàm kiểm tra xem có redis có đang hoạt động hay không theo khoảng thời gian định kỳ
func (r *RedisHealth) RedisHealthChecker(ctx context.Context, rds *redis.Client, interval time.Duration) {
	go func() {

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := rds.Ping(ctx).Result()
				if err != nil {
					r.SetAlive(false)
				} else {
					r.SetAlive(true)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
