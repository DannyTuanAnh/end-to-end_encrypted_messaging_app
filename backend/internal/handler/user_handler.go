package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"cloud.google.com/go/storage"
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
	gcs_client   *storage.Client
}

func NewUserHandler(user_client *client.UserClient, rdb *redis.Client, gcs_client *storage.Client) *UserHandler {
	return &UserHandler{
		user_client:  user_client,
		redis_memory: rdb,
		gcs_client:   gcs_client,
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

func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	var req dto.UpdateUserRequest

	if err := ctx.ShouldBind(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	imageFile, err := ctx.FormFile("avatar_url")

	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			imageFile = nil
		} else {
			utils.ResponseErrorAbort(ctx, utils.NewError(fmt.Sprintf("Failed to get avatar file: %v", err), utils.ErrCodeBadRequest))
			return
		}
	}

	if imageFile != nil {
		file, objName, err := utils.ValidateAndReturnObjNameImage(imageFile)
		if err != nil {
			utils.ResponseErrorAbort(ctx, utils.NewError(fmt.Sprintf("Invalid avatar file: %v", err), utils.ErrCodeBadRequest))
			return
		}

		contentType := imageFile.Header.Get("Content-Type")

		err = h.uploadToGCS(ctx, file, objName, contentType)
		if err != nil {
			utils.ResponseErrorAbort(ctx, utils.NewError(fmt.Sprintf("Failed to upload avatar to GCS: %v", err), utils.ErrCodeInternal))
			return
		}
	}

	var verifiedPhoneStr string

	if v, exist := ctx.Get("verified_phone"); exist {
		if s, ok := v.(string); ok {
			verifiedPhoneStr = s
		}
	}

	if req.Phone != nil && verifiedPhoneStr == "" {
		utils.ResponseErrorAbort(ctx, utils.NewError("Phone number provided but not verified", utils.ErrCodeBadRequest))
		return
	}

	if req.Phone != nil && verifiedPhoneStr != *req.Phone {
		utils.ResponseErrorAbort(ctx, utils.NewError("Provided phone number does not match verified phone number", utils.ErrCodeBadRequest))
		return
	}

	if req.Phone == nil && verifiedPhoneStr != "" {
		utils.ResponseErrorAbort(ctx, utils.NewError("Verified phone number exists but no phone number provided in request", utils.ErrCodeBadRequest))
		return
	}

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

	resp, err := h.user_client.Client.UpdateProfile(c, &user_proto.UpdateProfileRequest{
		UserId:   userID,
		Name:     req.Name,
		Birthday: req.Birthday,
		Phone:    req.Phone,
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

func (h *UserHandler) VerifyIDTokenOTP(ctx *gin.Context) {
	var req dto.UpdateUserRequest

	if err := ctx.ShouldBind(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		ctx.Abort()
		return
	}

	if req.Phone == nil {
		ctx.Next()
		return
	}

	var reqOTP dto.VerifyIDTokenOTP

	if err := ctx.ShouldBind(&reqOTP); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		ctx.Abort()
		return
	}

	resp, err := h.user_client.Client.VerifyIDTokenOTP(ctx, &user_proto.VerifyIDTokenOTPRequest{IdToken: reqOTP.IDToken})
	if err != nil {
		utils.WriteGRPCErrorToGin(ctx, err)
		ctx.Abort()
		return
	}

	if !resp.Valid {
		utils.ResponseErrorAbort(ctx, utils.NewError(resp.GetMessage(), utils.ErrCodeUnauthorized))
		return
	}

	ctx.Set("verified_phone", resp.GetVerifiedPhone())
	ctx.Next()
}

func (h *UserHandler) ReportUserImage(ctx *gin.Context) {
	var req dto.ReportUserImageRequest

	if err := ctx.ShouldBind(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	resp, err := h.user_client.Client.ReportUserImage(ctx, &user_proto.ReportUserImageRequest{
		Uuid:   req.UUID,
		Bucket: req.Bucket,
		Name:   req.Name,
	})
	if err != nil {
		utils.WriteGRPCErrorToGin(ctx, err)
		return
	}

	utils.ResponseSuccessWithMessage(ctx, http.StatusOK, resp.GetMessage())
}

func (h *UserHandler) uploadToGCS(ctx context.Context, file multipart.File, objectName string, contentType string) error {
	bucketName := utils.GetEnv("GOOGLE_CLOUD_STORAGE_BUCKET_RAW", "")
	wc := h.gcs_client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	wc.ObjectAttrs.CacheControl = "no-store, max-age=0"
	wc.ContentType = contentType

	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("failed to copy file to GCS: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %v", err)
	}

	return nil
}
