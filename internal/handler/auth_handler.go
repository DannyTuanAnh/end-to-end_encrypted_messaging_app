package handler

import (
	"log"
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
		log.Printf("Failed to bind JSON input: %v", err)
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	authReq := &auth_proto.LoginRequest{
		AuthorCode: input.AuthCode,
	}

	resp, err := h.auth_client.Client.LoginGoogle(ctx, authReq)
	if err != nil {
		log.Printf("Failed to call LoginGoogle on AuthService: %v", err)
		utils.ResponseError(ctx, err)
		return
	}

	log.Println("Received response from AuthService LoginGoogle:", resp)

	respAuthDTO := dto.MapLoginGoogleResponseToDTO(resp)

	log.Println("Mapped AuthService response to DTO:", respAuthDTO)
	utils.ResponseSuccess(ctx, http.StatusOK, respAuthDTO)
}
