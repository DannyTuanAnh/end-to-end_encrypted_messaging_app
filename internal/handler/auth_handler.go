package handler

import (
	"log"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/auth"
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

func (h *AuthHandler) Login(ctx *gin.Context) {
	log.Println("Login endpoint called")

	authReq := &auth_proto.LoginRequest{
		AuthorCode: "fake_auth_code",
	}
	_, err := h.auth_client.Client.Login(ctx, authReq)
	if err == nil {
		log.Println("Login request sent successfully to AuthService")
	} else {
		log.Printf("Failed to send login request to AuthService: %v", err)
		return
	}

}
