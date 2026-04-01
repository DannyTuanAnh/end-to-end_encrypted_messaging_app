package app

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/handler"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/routes"
)

type AuthModule struct {
	routes routes.Routes
}

func NewAuthModule(addr string) *AuthModule {
	// 1. Initialize repository
	auth_client, err := client.NewAuthClient(addr)
	if err != nil {
		panic("Failed to initialize auth client: " + err.Error())
	}

	// 2. Initialize handler
	auth_handler := handler.NewAuthHandler(auth_client)

	// 3. Initialize routes
	auth_routes := routes.NewAuthRoutes(auth_handler)

	return &AuthModule{routes: auth_routes}
}

func (au *AuthModule) Routes() routes.Routes {
	return au.routes
}
