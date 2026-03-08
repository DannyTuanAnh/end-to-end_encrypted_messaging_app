package app

import (
	"github.com/user-manage/internal/handler"
	"github.com/user-manage/internal/repository"
	"github.com/user-manage/internal/routes"
	"github.com/user-manage/internal/service"
)

type UserModule struct {
	routes routes.Routes
}

func NewUserModule() *UserModule {
	// 1. Initialize repository
	user_repo := repository.NewInMemoryUserRepository()

	// 2. Initialize service
	user_service := service.NewUserService(user_repo)

	// 3. Initialize handler
	user_handler := handler.NewUserHandler(user_service)

	// 4. Initialize routes
	user_routes := routes.NewUserRoutes(user_handler)

	return &UserModule{routes: user_routes}
}

func (um *UserModule) Routes() routes.Routes {
	return um.routes
}
