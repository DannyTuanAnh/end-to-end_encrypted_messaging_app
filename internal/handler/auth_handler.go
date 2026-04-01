package handler

import (
	"net/http"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/dto"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/auth"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth_client *client.AuthClient
}

func NewAuthHandler(auth_client *client.AuthClient) *AuthHandler {
	return &AuthHandler{
		auth_client: auth_client,
	}
}

func (h *AuthHandler) LoginGoogle(ctx *gin.Context) {
	var input dto.RequestLoginGoogle

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	authReq := &auth_proto.LoginRequest{
		AuthorCode: input.AuthCode,
	}

	resp, err := h.auth_client.Client.LoginGoogle(ctx, authReq)
	if err != nil {
		utils.WriteGRPCErrorToGin(ctx, err)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    resp.Session,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Domain:   utils.GetEnv("COOKIE_DOMAIN", ""),
		Path:     "/",
		MaxAge:   utils.GetEnvInt("SESSION_ID_MAX_AGE", 168) * 3600,
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "device_id",
		Value:    resp.DeviceId,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Domain:   utils.GetEnv("COOKIE_DOMAIN", ""),
		Path:     "/",
		MaxAge:   utils.GetEnvInt("DEVICE_ID_MAX_AGE", 168) * 3600,
	})

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	sessionID, exist := ctx.Get("session_id")
	if !exist {
		utils.ResponseErrorAbort(ctx, utils.NewError("session_id not found in context", utils.ErrCodeUnauthorized))
		return
	}

	deviceID, exist := ctx.Get("device_id")
	if !exist {
		utils.ResponseErrorAbort(ctx, utils.NewError("device_id not found in context", utils.ErrCodeUnauthorized))
		return
	}

	req := &auth_proto.LogoutRequest{
		SessionId: sessionID.(string),
		DeviceId:  deviceID.(string),
	}

	_, err := h.auth_client.Client.Logout(ctx, req)
	if err != nil {
		utils.WriteGRPCErrorToGin(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

func (h *AuthHandler) LogoutAll(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		utils.ResponseErrorAbort(ctx, utils.NewError("user_id not found in context", utils.ErrCodeUnauthorized))
		return
	}

	id, ok := userID.(int64)
	if !ok {
		utils.ResponseErrorAbort(ctx, utils.NewError("user_id in context has invalid type", utils.ErrCodeInternal))
		return
	}

	req := &auth_proto.LogoutAllRequest{
		UserId: id,
	}

	_, err := h.auth_client.Client.LogoutAll(ctx, req)
	if err != nil {
		utils.ResponseErrorAbort(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
