package routes

import (
	"context"
	"net/http"

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
	r.Use(middleware.CORSMiddleware(),
		middleware.RateLimitMiddleware(ctx, rdb, 60, 100), // 100 requests per 60 seconds
		middleware.LoggerMiddleware(),
	)

	// Public health check (NO API KEY)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	public := r.Group("/api/v1")
	public.Use(middleware.ApiKeyMiddleware(db, rdb))

	for _, route := range routes {
		if publicRoute, ok := route.(interface{ RegisterPublic(r *gin.RouterGroup) }); ok {
			publicRoute.RegisterPublic(public)
		}
	}

	protect := r.Group("/api/v1")
	protect.Use(middleware.AuthMiddleware(db, rdb))
	for _, route := range routes {
		route.Register(protect)
	}
}
