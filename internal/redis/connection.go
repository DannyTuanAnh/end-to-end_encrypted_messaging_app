package redis_memory

import (
	"context"
	"log"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	RDB *redis.Client
}

func InitRedis() (*Redis, error) {
	v_redis := config.NewConfigRedis().Redis

	rdb := redis.NewClient(&redis.Options{
		Addr:     v_redis.Addr,
		Password: v_redis.Password,
		DB:       v_redis.DB,

		PoolSize:     v_redis.PoolSize,
		MinIdleConns: v_redis.MinIdleConns,

		DialTimeout:  v_redis.DialTimeout,
		ReadTimeout:  v_redis.ReadTimeout,
		WriteTimeout: v_redis.WriteTimeout,

		MaxRetries:      v_redis.MaxRetries,
		MinRetryBackoff: v_redis.MinRetryBackOff,
		MaxRetryBackoff: v_redis.MaxRetryBackOff,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Connecting to redis successfully")

	return &Redis{RDB: rdb}, nil
}

func (r *Redis) CloseRedis() {
	if r.RDB != nil {
		if err := r.RDB.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v\n", err)
		}

		log.Println("Redis client closed")
		return
	}

	log.Println("Redis client was not initialized, no need to close")
}
