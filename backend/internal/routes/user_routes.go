package routes

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/handler"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	user_handler *handler.UserHandler
}

func NewUserRoutes(handler *handler.UserHandler) Routes {
	return &UserRoutes{
		user_handler: handler,
	}
}

func (ur *UserRoutes) Register(r *gin.RouterGroup) {
	user := r.Group("/user")
	{
		// Users.GET("", ur.User_handler.GetAllUsers)
		user.GET("/profile", ur.user_handler.GetProfile)
		user.GET("/profile/:uuid", ur.user_handler.GetProfileByUserUUID)
		user.DELETE("/disable", ur.user_handler.DisableUser)
		user.PUT("/profile", ur.user_handler.VerifyIDTokenOTP, ur.user_handler.UpdateProfile)
		user.POST("/report-avatar", ur.user_handler.ReportUserImage)
		// Users.DELETE("/:uuid", ur.User_handler.DeleteUser)
	}
}
