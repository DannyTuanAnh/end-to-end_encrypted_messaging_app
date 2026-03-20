package routes

import (
	"context"

	api_key_middleware "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/middleware/api_key"
	auth_middleware "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/middleware/auth"
	logger_middleware "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/middleware/logger"
	rate_limiter_middleware "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/middleware/rate_limiter"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/gin-gonic/gin"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(ctx context.Context, r *gin.Engine, redisHealth *utils.RedisHealth, routes ...Routes) {
	// Register middleware for all routes, including: logger, rate limiter, API key and authentication
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
