package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/user-manage/internal/handler"
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
		users.GET("", ur.user_handler.GetAllUsers)
		users.GET("/:uuid", ur.user_handler.GetUserByUUID)
		users.POST("", ur.user_handler.CreateUser)
		users.PUT("/:uuid", ur.user_handler.UpdateUser)
		users.DELETE("/:uuid", ur.user_handler.DeleteUser)
	}
}
