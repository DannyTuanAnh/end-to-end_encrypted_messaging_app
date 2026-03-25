package routes

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(ctx context.Context, r *gin.Engine, rdb *redis.Client, db sqlc.Querier, routes ...Routes) {
	// Register middleware for all routes, including: logger, rate limiter, API key and authentication
	r.Use(middleware.LoggerMiddleware(),
		middleware.RateLimitMiddleware(ctx, rdb, 60, 100), // 100 requests per 60 seconds
		middleware.ApiKeyMiddleware(db, rdb),
		middleware.AuthMiddleware(),
	)

	// Register all routes under the '/api' group
	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
}
