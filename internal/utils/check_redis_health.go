package utils

import (
	"context"
	"os"
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

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS"),
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,

		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   2,

		PoolTimeout: 5 * time.Second, // Thời gian chờ tối đa khi pool đã đầy
	})
}

// Hàm kiểm tra xem có redis có đang hoạt động hay không theo khoảng thời gian định kỳ
func RedisHealthChecker(ctx context.Context, rds *redis.Client, redisHealth *RedisHealth, interval time.Duration) {
	go func() {

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := rds.Ping(ctx).Result()
				if err != nil {
					redisHealth.SetAlive(false)
				} else {
					redisHealth.SetAlive(true)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func StartChecker(ctx context.Context) {
	redisHealth := NewRedisHealth()
	// 1. khởi tạo kết nối Redis
	rds := ConnectRedis()
	// 2. khởi động goroutine kiểm tra sức khỏe Redis
	startOnce.Do(func() {
		RedisHealthChecker(ctx, rds, redisHealth, 3*time.Second)
	})
}
