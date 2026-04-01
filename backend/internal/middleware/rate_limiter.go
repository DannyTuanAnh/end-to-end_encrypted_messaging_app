package middleware

import (
	"context"
	"strconv"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/gin-gonic/gin"
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
)

func Allow(ctx context.Context, key string, rdb *redis.Client, timeTl int, maxReq int) (bool, error) {
	allowed, err := rateLimitScript.Run(ctx, rdb, []string{key}, timeTl, maxReq).Int()

	if err != nil {
		return false, err
	}

	return allowed == 1, nil
}

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	// lấy IP đã được mã hóa thông qua proxy (nếu có)
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func GetRateLimitKey(ctx *gin.Context) string {
	if uid, exists := ctx.Get("user_id"); exists {
		return "rate:user:" + strconv.FormatInt(uid.(int64), 10)
	}

	return "rate:ip:" + getClientIP(ctx)
}

func RateLimitMiddleware(ctx context.Context, rdb *redis.Client, timeTl int, maxReq int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := GetRateLimitKey(c)

		allow, err := Allow(ctx, key, rdb, timeTl, maxReq)
		if err != nil {
			utils.ResponseErrorAbort(c, utils.WrapError(err, "Failed to check rate limit", utils.ErrCodeInternal))
			return
		}

		if !allow {
			utils.ResponseErrorAbort(c, utils.NewError("Rate limit exceeded", utils.ErrCodeTooManyRequests))
			return
		}

		c.Next()
	}
}

// sử dụng ApacheBench của Golang để test rate limiting
// ab -n 100 -c 1 -H "X-API-KEY: 4f4c48fb-665a-4e6b-a498-01e72e89db7c" localhost:8080/api/v1/users/1
// ab -n 110 -c 1 -H "X-API-KEY: 4f4c48fb-665a-4e6b-a498-01e72e89db7c" localhost:8080/api/v1/users/1
