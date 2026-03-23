package app

// import (
// 	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/handler"
// 	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
// 	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/routes"
// 	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
// )

// type AuthModule struct {
// 	service service.AuthService
// }

// func NewAuthModule() *AuthModule {
// 	// 1. Initialize repository
// 	user_repo := repository.NewInMemoryUserRepository()

// 	// 2. Initialize service
// 	user_service := service.NewUserService(user_repo)

// 	// 3. Initialize handler
// 	user_handler := handler.NewUserHandler(user_service)

// 	// 4. Initialize routes
// 	user_routes := routes.NewUserRoutes(user_handler)

// 	return &UserModule{routes: user_routes}
// }

// func (um *UserModule) Routes() routes.Routes {
// 	return um.routes
// }
