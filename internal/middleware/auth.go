package middleware

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(db sqlc.Querier, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userId int64

		sessionId, err, errCode := ValidateSession(ctx)
		if err != nil {
			utils.ResponseErrorAbort(ctx, utils.WrapError(err, "Failed to validate session", errCode))
			return
		}

		deviceId, err, errCode := ValidateDeviceID(ctx)
		if err != nil {
			utils.ResponseErrorAbort(ctx, utils.WrapError(err, "Failed to validate device ID", errCode))
			return
		}

		data, err := rdb.Get(ctx, fmt.Sprintf("session:%s", sessionId)).Bytes()
		if err == nil && len(data) > 0 {
			log.Println("Session found in Redis")
			var valueSession models.SessionRedis
			if err := json.Unmarshal(data, &valueSession); err != nil {
				utils.ResponseErrorAbort(ctx, utils.WrapError(err, "Failed to decode session from Redis", utils.ErrCodeInternal))
				return
			}

			version, err := utils.GetKeyRedisAndConvertToInt(ctx, fmt.Sprintf("user:%d:session_version", valueSession.UserID), rdb)
			if err != nil {
				utils.ResponseErrorAbort(ctx, utils.WrapError(err, "Failed to get session version from Redis", utils.ErrCodeInternal))
				return
			}

			if !valueSession.Valid || version != valueSession.SessionVersion {
				utils.ResponseErrorAbort(ctx, utils.NewError("Invalid session", utils.ErrCodeUnauthorized))
				return
			}

			if valueSession.UserID == 0 {
				utils.ResponseErrorAbort(ctx, utils.NewError("Invalid session: user_id is 0", utils.ErrCodeUnauthorized))
				return
			}

			if valueSession.DeviceID != deviceId {
				utils.ResponseErrorAbort(ctx, utils.NewError("Invalid session: device mismatch", utils.ErrCodeUnauthorized))
				return
			}

			userId = valueSession.UserID

		} else if errors.Is(err, redis.Nil) {
			log.Println("Session not found in Redis, checking database...")
			params := sqlc.CheckSessionParams{
				SessionID: sessionId,
				DeviceID:  deviceId,
			}

			result, err := db.CheckSession(ctx, params)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					utils.ResponseErrorAbort(ctx, utils.NewError("Session not found", utils.ErrCodeUnauthorized))
					return
				}
				utils.ResponseErrorAbort(ctx, utils.WrapError(err, "Failed to check session in database", utils.ErrCodeInternal))
				return
			}

			if result.Revoked {
				utils.ResponseErrorAbort(ctx, utils.NewError("Session revoked", utils.ErrCodeUnauthorized))
				return
			}

			userId = result.UserID
		} else {
			utils.ResponseErrorAbort(ctx, utils.WrapError(err, "Failed to get session from Redis", utils.ErrCodeInternal))
			return
		}

		ctx.Set("user_id", userId)
		ctx.Set("device_id", deviceId.String())
		ctx.Set("session_id", sessionId.String())

		ctx.Next()

	}
}

func ValidateSession(ctx *gin.Context) (uuid.UUID, error, utils.ErrCode) {
	sessionId, err := ctx.Cookie("session_id")
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to get session from cookie: %v", err), utils.ErrCodeInternal
	}

	if sessionId == "" {
		return uuid.Nil, fmt.Errorf("session_id cookie is empty"), utils.ErrCodeUnauthorized
	}

	if !utils.CheckUUID(sessionId) {
		return uuid.Nil, fmt.Errorf("invalid session_id format"), utils.ErrCodeUnauthorized
	}

	sessionUUID, err := uuid.Parse(sessionId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to parse session_id: %v", err), utils.ErrCodeInternal
	}

	return sessionUUID, nil, ""
}

func ValidateDeviceID(ctx *gin.Context) (uuid.UUID, error, utils.ErrCode) {
	deviceId, err := ctx.Cookie("device_id")
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to get device_id from cookie: %v", err), utils.ErrCodeInternal
	}

	if deviceId == "" {
		return uuid.Nil, fmt.Errorf("device_id cookie is empty"), utils.ErrCodeUnauthorized
	}

	if !utils.CheckUUID(deviceId) {
		return uuid.Nil, fmt.Errorf("invalid device_id format"), utils.ErrCodeUnauthorized
	}

	deviceUUID, err := uuid.Parse(deviceId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to parse device_id: %v", err), utils.ErrCodeInternal
	}

	return deviceUUID, nil, ""
}
