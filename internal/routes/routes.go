package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	api_key_middleware "github.com/user-manage/internal/middleware/api_key"
	auth_middleware "github.com/user-manage/internal/middleware/auth"
	logger_middleware "github.com/user-manage/internal/middleware/logger"
	rate_limiter_middleware "github.com/user-manage/internal/middleware/rate_limiter"
	"github.com/user-manage/internal/utils"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(ctx context.Context, r *gin.Engine, redisHealth *utils.RedisHealth, routes ...Routes) {
	// Đăng ký middleware cho tất cả các route, bao gồm: logger, rate limiter, API key và authentication
	r.Use(logger_middleware.LoggerMiddleware(),
		rate_limiter_middleware.RateLimitMiddleware(ctx, redisHealth),
		api_key_middleware.ApiKeyMiddleware(),
		auth_middleware.AuthMiddleware(),
	)

	// Register all routes under the '/api' group
	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
}
