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

func (ur *AuthRoutes) Register(r *gin.RouterGroup) {
	auths := r.Group("/auth")
	{
		// auths.GET("", ur.Auth_handler.GetAllAuths)
		// auths.GET("/:uuid", ur.Auth_handler.GetAuthByUUID)
		auths.GET("/login", ur.auth_handler.Login)
		auths.POST("/google/callback", ur.auth_handler.Login)
		// auths.PUT("/:uuid", ur.Auth_handler.UpdateAuth)
		// auths.DELETE("/:uuid", ur.Auth_handler.DeleteAuth)
	}
}
