package handler

import (
	"log"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
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
}
