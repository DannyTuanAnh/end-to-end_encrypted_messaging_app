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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing X-API-KEY"})
			return
		}

		hash := utils.HashAPIKey(apiKey)

		genNum, err := utils.GetKeyRedisAndConvertToInt(ctx, "generation_api_key", rdb)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check API key generation"})
			return
		}

		cacheKey := fmt.Sprintf("api_key:%d:%s", genNum, hash)
		val, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil && val == "1" {
			log.Println("API key found in cache")
			ctx.Next()
			return
		} else if err != nil && !errors.Is(err, redis.Nil) {
			log.Printf("redis get error: %v", err)
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

		log.Println("API key found in database and is active")

		// Cache the valid API key in Redis
		if err := rdb.Set(ctx, cacheKey, "1", 24*7*time.Hour).Err(); err != nil {
			log.Printf("redis set error: %v", err)
		}

		ctx.Next()
	}
}
