package rate_limiter_middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/user-manage/internal/utils"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) bool
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

func RateLimitMiddleware(ctx context.Context, redisHealth *utils.RedisHealth) gin.HandlerFunc {

	// 1. khởi tạo kết nối Redis
	rds := utils.ConnectRedis()
	// 2. khởi tạo rate limiting
	memRateLimiter := NewRateLimitMemory()
	redisRateLimiter := NewRateLimitRedis(rds, 60, 100)

	return func(c *gin.Context) {
		key := GetRateLimitKey(c)

		var allowed bool
		if redisHealth.IsAlive() {
			allowed = redisRateLimiter.Allow(ctx, key)
		} else {
			allowed = memRateLimiter.Allow(ctx, key)
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Quá nhiều yêu cầu, vui lòng thử lại sau"})
			return
		}

		c.Next()
	}
}
