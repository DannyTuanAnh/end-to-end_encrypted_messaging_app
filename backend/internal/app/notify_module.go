package app

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/handler"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/routes"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type NotifyModule struct {
	routes routes.Routes
}

func NewNotifyModule(addr string) *NotifyModule {
	// Load TLS credentials for gRPC client
	// Call by API Gateway, so use API Gateway's certs
	apiGatewayCertFile := utils.GetEnv("PATH_CERT_API_GATEWAY", "")
	apiGatewayKeyFile := utils.GetEnv("PATH_KEY_API_GATEWAY", "")

	// 1. Initialize repository
	notify_client, err := client.NewUserClient(addr, apiGatewayCertFile, apiGatewayKeyFile)
	if err != nil {
		panic("Failed to initialize notify client: " + err.Error())
	}

	// 2. Initialize handler
	notify_handler := handler.NewNotifyHandler(notify_client)

	// 3. Initialize routes
	notify_routes := routes.NewNotifyRoutes(notify_handler)

	return &NotifyModule{routes: notify_routes}
}

func (au *NotifyModule) Routes() routes.Routes {
	return au.routes
}
