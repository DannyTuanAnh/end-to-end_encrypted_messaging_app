package redis_memory

import (
	"context"
	"log"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	RDB       *redis.Client
	Redis_GCP *redis.Client
}

func InitRedis() (*Redis, error) {
	v_redis := config.NewConfigRedis().Redis

	rdb := redis.NewClient(&redis.Options{
		Addr:     v_redis.Addr,
		Password: v_redis.Password,
		DB:       v_redis.DB,

		PoolSize:     v_redis.Options.PoolSize,
		MinIdleConns: v_redis.Options.MinIdleConns,

		DialTimeout:  v_redis.Options.DialTimeout,
		ReadTimeout:  v_redis.Options.ReadTimeout,
		WriteTimeout: v_redis.Options.WriteTimeout,

		MaxRetries:      v_redis.Options.MaxRetries,
		MinRetryBackoff: v_redis.Options.MinRetryBackOff,
		MaxRetryBackoff: v_redis.Options.MaxRetryBackOff,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Connecting to redis successfully")

	v_redis_gcp := config.NewConfigGCPRedis().RedisGCP

	redis_gcp := redis.NewClient(&redis.Options{
		Addr:     v_redis_gcp.Addr,
		Password: v_redis_gcp.Password,
		DB:       v_redis_gcp.DB,

		PoolSize:     v_redis_gcp.Options.PoolSize,
		MinIdleConns: v_redis_gcp.Options.MinIdleConns,

		DialTimeout:  v_redis_gcp.Options.DialTimeout,
		ReadTimeout:  v_redis_gcp.Options.ReadTimeout,
		WriteTimeout: v_redis_gcp.Options.WriteTimeout,

		MaxRetries:      v_redis_gcp.Options.MaxRetries,
		MinRetryBackoff: v_redis_gcp.Options.MinRetryBackOff,
		MaxRetryBackoff: v_redis_gcp.Options.MaxRetryBackOff,
	})

	ctxGCP, cancelGCP := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelGCP()

	if err := redis_gcp.Ping(ctxGCP).Err(); err != nil {
		return nil, err
	}

	log.Println("Connecting to GCP redis successfully")

	return &Redis{
		RDB:       rdb,
		Redis_GCP: redis_gcp,
	}, nil
}

func (r *Redis) CloseRedis() {
	if r == nil {
		log.Println("Redis struct is nil, no need to close")
		return
	}

	if r.RDB != nil {
		if err := r.RDB.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v\n", err)
		}

		log.Println("Redis client closed")
	}

	if r.Redis_GCP != nil {
		if err := r.Redis_GCP.Close(); err != nil {
			log.Printf("Failed to close GCP Redis client: %v\n", err)
		}

		log.Println("GCP Redis client closed")
	}
}
