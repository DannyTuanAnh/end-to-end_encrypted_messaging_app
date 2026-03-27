package app

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/handler"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/routes"
)

type UserModule struct {
	routes routes.Routes
}

func NewUserModule(addr string) *UserModule {
	// 1. Initialize repository
	user_client, err := client.NewUserClient(addr)
	if err != nil {
		panic("Failed to initialize User client: " + err.Error())
	}

	// 2. Initialize handler
	user_handler := handler.NewUserHandler(user_client)

	// 3. Initialize routes
	user_routes := routes.NewUserRoutes(user_handler)

	return &UserModule{routes: user_routes}
}

func (us *UserModule) Routes() routes.Routes {
	return us.routes
}
