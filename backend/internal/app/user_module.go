package app

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/handler"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/routes"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

type UserModule struct {
	routes routes.Routes
}

func NewUserModule(addr string, ctx context.Context, rdb *redis.Client) *UserModule {
	// Load TLS credentials for gRPC client
	// Call by API Gateway, so use API Gateway's certs
	apiGatewayCertFile := utils.GetEnv("PATH_CERT_API_GATEWAY_CLIENT", "")
	apiGatewayKeyFile := utils.GetEnv("PATH_KEY_API_GATEWAY_CLIENT", "")

	// 1. Initialize repository
	user_client, err := client.NewUserClient(addr, apiGatewayCertFile, apiGatewayKeyFile)
	if err != nil {
		panic("Failed to initialize User client: " + err.Error())
	}

	// 2. Initialize handler
	user_handler := handler.NewUserHandler(user_client, rdb, connectGCS(ctx))

	// 3. Initialize routes
	user_routes := routes.NewUserRoutes(user_handler)

	return &UserModule{routes: user_routes}
}

func (us *UserModule) Routes() routes.Routes {
	return us.routes
}

func connectGCS(ctx context.Context) *storage.Client {
	serviceAccountKey := utils.GetEnv("GOOGLE_CLOUD_STORAGE_CREDENTIALS", "")

	var storage_client *storage.Client
	var err error
	if serviceAccountKey != "" {
		storage_client, err = storage.NewClient(ctx, option.WithAuthCredentialsFile(option.ServiceAccount, serviceAccountKey))
	} else {
		storage_client, err = storage.NewClient(ctx)
	}

	if err != nil {
		panic("Failed to initialize Google Cloud Storage client: " + err.Error())
	}

	return storage_client
}
