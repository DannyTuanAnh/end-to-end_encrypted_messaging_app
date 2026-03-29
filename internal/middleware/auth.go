package middleware

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(db sqlc.Querier, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := ctx.Cookie("session_id")
		if err != nil {
			utils.ResponseError(ctx, utils.WrapError(err, "session_id cookie not found", utils.ErrCodeUnauthorized))
			return
		}

		if sessionId == "" {
			utils.ResponseError(ctx, utils.NewError("session_id cookie not found", utils.ErrCodeUnauthorized))
			return
		}

		data, err := rdb.Get(ctx, fmt.Sprintf("session:%s", sessionId)).Bytes()
		if err != nil {
			utils.ResponseError(ctx, utils.WrapError(err, "Failed to get session from Redis", utils.ErrCodeUnauthorized))
			return
		}

		var valueSession models.SessionRedis
		if err := json.Unmarshal(data, &valueSession); err != nil {
			utils.ResponseError(ctx, utils.WrapError(err, "Failed to decode session from Redis", utils.ErrCodeUnauthorized))
			return
		}

		log.Println("Session found in Redis: ", valueSession)

		ctx.Next()

	}
}
