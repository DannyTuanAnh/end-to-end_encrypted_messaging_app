package rate_limiter_middleware

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rateLimitScript = redis.NewScript(`
	local current = redis.call("INCR", KEYS[1])

	if current == 1 then 
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end

	if current > tonumber(ARGV[2]) then
		return 0
	end

	return 1
	`)
	connectOnce sync.Once
)

type rateLimitRedis struct {
	rds *redis.Client

	timeTl int
	maxReq int
}

func NewRateLimitRedis(rds *redis.Client, timeTl int, maxReq int) *rateLimitRedis {
	return &rateLimitRedis{
		rds:    rds,
		timeTl: timeTl,
		maxReq: maxReq,
	}
}

func (r *rateLimitRedis) Allow(ctx context.Context, key string) bool {
	allowed, err := rateLimitScript.Run(ctx, r.rds, []string{key}, r.timeTl, r.maxReq).Int()

	if err != nil {
		return false
	}

	return allowed == 1
}

// sử dụng ApacheBench của Golang để test rate limiting
// ab -n 100 -c 1 -H "X-API-KEY: 4f4c48fb-665a-4e6b-a498-01e72e89db7c" localhost:8080/api/v1/users/1
// ab -n 110 -c 1 -H "X-API-KEY: 4f4c48fb-665a-4e6b-a498-01e72e89db7c" localhost:8080/api/v1/users/1
