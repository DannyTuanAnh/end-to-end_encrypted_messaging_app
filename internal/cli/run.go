package cli

import (
	"context"
	"log"
	"os"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/app"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/redis/go-redis/v9"
)

type CLI struct {
	registry *ServiceRegistry
}

func NewCLI(db sqlc.Querier, rdb *redis.Client) *CLI {
	cli_services = []ModelCLIService{
		app.NewAPIKeyModule(db, rdb),
	}

	registry := NewServiceRegistry(cli_services)

	return &CLI{
		registry: registry,
	}
}

func (c *CLI) Run(ctx context.Context) bool {
	if len(os.Args) < 2 {
		return false
	}

	switch os.Args[1] {
	case "generate-api-key":
		return c.handleGenerateAPIKey(ctx, os.Args[2:])
	case "revoke-api-key":
		return c.handleRevokeAPIKey(ctx, os.Args[2:])
	}

	return false
}

func (c *CLI) handleGenerateAPIKey(ctx context.Context, args []string) bool {
	apiKeyService, ok := GetService[service.APIKeyService](c.registry, "api_key")
	if !ok {
		log.Println("APIKeyService not found in registry")
		return false
	}

	cmdArgs, err := utils.ParseGenerateAPIKeyArgs(args)
	if err != nil {
		log.Printf("Error when parsing arguments: %v\n", err)
		return false
	}

	if err := apiKeyService.CreateAPIKey(ctx, cmdArgs); err != nil {
		log.Printf("Error when generating API key: %v\n", err)
	}

	return true
}

func (c *CLI) handleRevokeAPIKey(ctx context.Context, args []string) bool {
	apiKeyService, ok := GetService[service.APIKeyService](c.registry, "api_key")
	if !ok {
		log.Println("APIKeyService not found in registry")
		return false
	}

	cmdArgs, err := utils.ParseRevokeAPIKeyArgs(args)
	if err != nil {
		log.Printf("Error when parsing arguments: %v\n", err)
		return false
	}

	if cmdArgs.RevokeAll {
		if err := apiKeyService.RevokeAll(ctx); err != nil {
			log.Printf("Error when revoke all API keys: %v\n", err)
		}
	}

	if cmdArgs.KeyID != "" {
		if err := apiKeyService.RevokeAPIKey(ctx, cmdArgs.KeyID); err != nil {
			log.Printf("Error when revoke API key (%s): %v\n", cmdArgs.KeyID, err)
		}
	}
	return true
}

func GetService[T any](sr *ServiceRegistry, name string) (T, bool) {
	var zero T

	svc, exists := sr.services[name]
	if !exists || svc == nil {
		log.Printf("Service '%s' not found in registry\n", name)
		return zero, false
	}

	typed, ok := svc.(T)
	if !ok {
		log.Printf("Service '%s' found but has unexpected type\n", name)
		return zero, false
	}

	return typed, true
}
