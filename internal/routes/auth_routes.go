package routes

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/handler"
	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	auth_handler *handler.AuthHandler
}

func NewAuthRoutes(handler *handler.AuthHandler) Routes {
	return &AuthRoutes{
		auth_handler: handler,
	}
}

func (ar *AuthRoutes) RegisterPublic(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/google/login", ar.auth_handler.LoginGoogle)
	}
}

func (ar *AuthRoutes) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.DELETE("/logout", ar.auth_handler.Logout)
		auth.DELETE("/logout/all", ar.auth_handler.LogoutAll)
	}
}
