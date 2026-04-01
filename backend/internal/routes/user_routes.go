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
	users := r.Group("/users")
	{
		// Users.GET("", ur.User_handler.GetAllUsers)
		// Users.GET("/:uuid", ur.User_handler.GetUserByUUID)
		users.PUT("profile", ur.user_handler.UpdateProfile)
		// Users.PUT("/:uuid", ur.User_handler.UpdateUser)
		// Users.DELETE("/:uuid", ur.User_handler.DeleteUser)
	}
}
