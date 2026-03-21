package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/routes"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
)

type ModelHTTP interface {
	Routes() routes.Routes
}

type ModelService interface {
	Name() string
	Service() any
}

type ServiceRegistry struct {
	services map[string]any
}

type Application struct {
	config          *config.Config
	route           *gin.Engine
	modules         []ModelHTTP
	module_services []ModelService
}

func NewServiceRegistry(modules []ModelService) *ServiceRegistry {
	m := make(map[string]any)

	for _, module := range modules {
		m[module.Name()] = module.Service()
	}

	return &ServiceRegistry{
		services: m,
	}
}

func NewApplication(ctx context.Context, cfg *config.Config, db sqlc.Querier) (bool, *Application) {
	// 1. Initialize the Gin router
	r := gin.Default()

	// 3. Initialize custom validator
	err := validation.InitValidator()
	if err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
	}

	// 4. Initialize health check for Redis
	redisHealth := utils.NewRedisHealth()

	// 5. Initialize modules
	modules := []ModelHTTP{
		NewUserModule(),
	}

	module_services := []ModelService{
		NewAPIKeyModule(db),
	}

	serviceRegister := NewServiceRegistry(module_services)

	// Check for command-line arguments and execute corresponding commands before starting the application
	hasExecuteCmd := commandTool(ctx, serviceRegister)
	if hasExecuteCmd {
		return true, nil
	}

	// 6. Register all routes from modules by calling the getModuleRoutes helper function to extract the routes from each module
	// and then passing them to the routes.RegisterRoutes function to register them with the Gin router
	routes.RegisterRoutes(ctx, r, redisHealth, getModuleRoutes(modules)...)

	return false, &Application{
		config:  cfg,
		route:   r,
		modules: modules,
	}
}

func (ac *Application) Run(ctx context.Context) (string, error) {
	// 1. Start server with shut down gracefully
	srv := &http.Server{
		Addr:         ":" + ac.config.Server.Port,
		Handler:      ac.route,
		ReadTimeout:  ac.config.Server.ReadTimeout,
		WriteTimeout: ac.config.Server.WriteTimeout,
		IdleTimeout:  ac.config.Server.IdleTimeout,
	}

	utils.StartChecker(ctx)

	// 2. Create a channel to listen for server errors
	errChan := make(chan error, 1)

	// 3. Listen and serve in a goroutine
	go func() {
		log.Printf("Server is running on port %s...", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// 4. Wait for an error or a shutdown signal
	select {
	case err := <-errChan:
		return "Server error", err

	case <-ctx.Done():
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), ac.config.Server.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return "Server forced to shutdown", err
		}
		return "Server exiting gracefully!", nil
	}

}

// getModuleRoutes is a helper function that takes a slice of Model interfaces
// and returns a slice of routes.Routes by calling the Routes() method on each module
func getModuleRoutes(models []ModelHTTP) []routes.Routes {
	routeList := make([]routes.Routes, len(models))
	for i, model := range models {
		routeList[i] = model.Routes()
	}
	return routeList
}

func commandTool(ctx context.Context, sr *ServiceRegistry) bool {
	if len(os.Args) > 1 {
		apiKeyService := GetService[service.APIKeyService](sr, "api_key")

		switch os.Args[1] {
		case "generate-api-key":
			args, err := utils.ParseGenerateAPIKeyArgs(os.Args[2:])

			if err != nil {
				log.Printf("Error when parsing arguments: %v\n", err)
				return true
			}

			if err := apiKeyService.CreateAPIKey(ctx, args); err != nil {
				log.Printf("Error when generating API key: %v\n", err)
			}
			return true
		case "revoke-api-key":
			args, err := utils.ParseRevokeAPIKeyArgs(os.Args[2:])

			if err != nil {
				log.Printf("Error when parsing arguments: %v\n", err)
				return true
			}

			if args.RevokeAll {
				if err := apiKeyService.RevokeAll(ctx); err != nil {
					log.Printf("Error when revoke all API keys: %v\n", err)
				}
			}

			if err := apiKeyService.RevokeAPIKey(ctx, args.KeyID); err != nil {
				log.Printf("Error when revoke API key (%s): %v\n", args.KeyID, err)
			}
			return true
		}
	}

	return false
}

func GetService[T any](sr *ServiceRegistry, name string) T {
	service := sr.services[name].(T)

	return service
}
