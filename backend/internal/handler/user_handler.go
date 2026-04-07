package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/dto"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/interceptor"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
)

type UserHandler struct {
	user_client  *client.UserClient
	redis_memory *redis.Client
}

func NewUserHandler(user_client *client.UserClient, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		user_client:  user_client,
		redis_memory: rdb,
	}
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId, exist := ctx.Get("user_id")
	if !exist {
		utils.ResponseErrorAbort(ctx, utils.NewError("User ID not found in context", utils.ErrCodeNotFound))
		return
	}

	userID, ok := userId.(int64)
	if !ok {
		utils.ResponseErrorAbort(ctx, utils.NewError("User ID in context has invalid type", utils.ErrCodeInternal))
		return
	}

	if userID <= 0 {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(errors.New("UserID must greater than 0")))
		return
	}

	profileData, err := h.redis_memory.Get(ctx, fmt.Sprintf("user_profile:%d", userID)).Bytes()
	if err == nil && len(profileData) > 0 {
		var cachedProfile models.ProfileRedis
		if err := json.Unmarshal(profileData, &cachedProfile); err == nil {
			log.Println("Profile data found in Redis in api-gateway for user ID: ", userID)
			utils.ResponseSuccessWithData(ctx, http.StatusOK, &user_proto.GetProfileByUserIDResponse{
				Uuid:      cachedProfile.UserUUID.String(),
				Name:      cachedProfile.Name,
				Email:     cachedProfile.Email,
				Phone:     cachedProfile.Phone,
				AvatarUrl: cachedProfile.AvatarUrl,
				Birthday:  cachedProfile.Birthday,
			})
			return
		}
	}

	baseCtx := ctx.Request.Context()

	c := context.WithValue(baseCtx, interceptor.CtxCallerKey, "api-gateway")
	c = context.WithValue(c, interceptor.CtxUserIDKey, userID)
	c = context.WithValue(c, interceptor.CtxAudKey, "user-service")

	resp, err := h.user_client.Client.GetProfileByUserID(c, &user_proto.GetProfileByUserIDRequest{UserId: userID})
	if err != nil {
		utils.WriteGRPCErrorToGin(ctx, err)
		return
	}

	utils.ResponseSuccessWithData(ctx, http.StatusOK, resp)
}

func (h *UserHandler) GetProfileByUserUUID(ctx *gin.Context) {
	var params dto.GetProfileByUserUUID
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userId, exist := ctx.Get("user_id")
	if !exist {
		utils.ResponseErrorAbort(ctx, utils.NewError("User ID not found in context", utils.ErrCodeNotFound))
		return
	}

	userID, ok := userId.(int64)
	if !ok {
		utils.ResponseErrorAbort(ctx, utils.NewError("User ID in context has invalid type", utils.ErrCodeInternal))
		return
	}

	if userID <= 0 {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(errors.New("UserID must greater than 0")))
		return
	}

	user_uuid, exist := ctx.Get("user_uuid")
	if exist {
		if userUUIDStr, ok := user_uuid.(string); ok {
			if userUUIDStr == params.UUID {
				utils.ResponseErrorAbort(ctx, utils.NewError("You can't find your profile", utils.ErrCodeNotFound))
				return
			}
		}
	}

	baseCtx := ctx.Request.Context()

	c := context.WithValue(baseCtx, interceptor.CtxCallerKey, "api-gateway")
	c = context.WithValue(c, interceptor.CtxUserIDKey, userID)
	c = context.WithValue(c, interceptor.CtxAudKey, "user-service")

	resp, err := h.user_client.Client.GetProfileByUserUUID(c, &user_proto.GetProfileByUserUUIDRequest{
		Uuid:   params.UUID,
		UserId: userID,
	})
	if err != nil {
		utils.WriteGRPCErrorToGin(ctx, err)
		return
	}

	utils.ResponseSuccessWithData(ctx, http.StatusOK, resp)
}

func (h *UserHandler) DisableUser(ctx *gin.Context) {
	userId, exist := ctx.Get("user_id")
	if !exist {
		utils.ResponseErrorAbort(ctx, utils.NewError("User ID not found in context", utils.ErrCodeNotFound))
	}

	userID, ok := userId.(int64)
	if !ok {
		utils.ResponseErrorAbort(ctx, utils.NewError("User ID in context has invalid type", utils.ErrCodeInternal))
	}

	if userID <= 0 {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(errors.New("UserID must greater than 0")))
	}

	baseCtx := ctx.Request.Context()

	c := context.WithValue(baseCtx, interceptor.CtxCallerKey, "api-gateway")
	c = context.WithValue(c, interceptor.CtxUserIDKey, userID)
	c = context.WithValue(c, interceptor.CtxAudKey, "user-service")

	_, err := h.user_client.Client.DisableUserByUserID(c, &user_proto.DisableUserRequest{UserId: userID})
	if err != nil {
		utils.WriteGRPCErrorToGin(ctx, err)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Domain:   utils.GetEnv("COOKIE_DOMAIN", ""),
		Path:     "/",
		MaxAge:   -1,
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "device_id",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Domain:   utils.GetEnv("COOKIE_DOMAIN", ""),
		Path:     "/",
		MaxAge:   -1,
	})

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
