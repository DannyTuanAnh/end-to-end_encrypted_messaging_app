package app

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/handler"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/routes"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type AuthModule struct {
	routes routes.Routes
}

func NewAuthModule(addr string) *AuthModule {
	// Load TLS credentials for gRPC client
	// Call by API Gateway, so use API Gateway's certs
	apiGatewayCertFile := utils.GetEnv("PATH_CERT_API_GATEWAY", "")
	apiGatewayKeyFile := utils.GetEnv("PATH_KEY_API_GATEWAY", "")

	// 1. Initialize repository
	auth_client, err := client.NewAuthClient(addr, apiGatewayCertFile, apiGatewayKeyFile)
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
