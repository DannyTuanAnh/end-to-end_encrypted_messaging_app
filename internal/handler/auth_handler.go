package handler

import (
	"net/http"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/dto"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/auth"
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
		utils.ResponseError(ctx, err)
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

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	sessionID, exist := ctx.Get("session_id")
	if !exist {
		utils.ResponseError(ctx, utils.NewError("session_id not found in context", utils.ErrCodeUnauthorized))
		return
	}

	if !utils.CheckUUID(sessionID.(string)) {
		utils.ResponseError(ctx, utils.NewError("invalid session_id format", utils.ErrCodeUnauthorized))
		return
	}

	deviceID, exist := ctx.Get("device_id")
	if !exist {
		utils.ResponseError(ctx, utils.NewError("device_id not found in context", utils.ErrCodeUnauthorized))
		return
	}

	if !utils.CheckUUID(deviceID.(string)) {
		utils.ResponseError(ctx, utils.NewError("invalid device_id format", utils.ErrCodeUnauthorized))
		return
	}

	req := &auth_proto.LogoutRequest{
		SessionId: sessionID.(string),
		DeviceId:  deviceID.(string),
	}

	_, err := h.auth_client.Client.Logout(ctx, req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
