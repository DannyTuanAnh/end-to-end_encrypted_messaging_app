package handler

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user_client *client.UserClient
}

func NewUserHandler(user_client *client.UserClient) *UserHandler {
	return &UserHandler{
		user_client: user_client,
	}
}

func (h *UserHandler) UpdateProfile(ctx *gin.Context) {}
