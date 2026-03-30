package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ApiKeyMiddleware(db sqlc.Querier, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			utils.ResponseErrorAbort(ctx, utils.NewError("API key is required", utils.ErrCodeUnauthorized))
			return
		}

		hash := utils.HashAPIKey(apiKey)

		genNum, err := utils.GetKeyRedisAndConvertToInt(ctx, "generation_api_key", rdb)
		if err != nil {
			log.Println("Error in get generation number of API key from Redis (in middleware layer): ", err)
		}

		cacheKey := fmt.Sprintf("api_key:%d:%s", genNum, hash)
		val, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil && val == "1" {
			ctx.Next()
			return
		} else if err != nil && !errors.Is(err, redis.Nil) {
			log.Println("Error in check API key in Redis (in middleware layer): ", err)
		}

		// Check if the API key exists in the database
		isActive, err := db.ValidateAPIKey(ctx, hash)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
				return
			}

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check API key"})
			return
		}

		if !isActive {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Inactive API key"})
			return
		}

		// Cache the valid API key in Redis
		if err := rdb.Set(ctx, cacheKey, "1", 24*7*time.Hour).Err(); err != nil {
			log.Printf("redis set error: %v", err)
		}

		ctx.Next()
	}
}
