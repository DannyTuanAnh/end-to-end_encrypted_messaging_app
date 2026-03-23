package handler

import "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"

type AuthHandler struct {
	auth_client *client.AuthClient
}

func NewAuthHandler(auth_client *client.AuthClient) *AuthHandler {
	return &AuthHandler{
		auth_client: auth_client,
	}
}
