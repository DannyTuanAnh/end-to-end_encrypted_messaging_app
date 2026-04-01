package app

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/redis/go-redis/v9"
)

type APIKeyModule struct {
	service service.APIKeyService
}

func NewAPIKeyModule(db sqlc.Querier, rdb *redis.Client) *APIKeyModule {
	// 1. Initialize repository
	apiKey_repo := repository.NewAPIKeyRepository(db)
	apiKey_service := service.NewAPIKeyService(apiKey_repo, rdb)

	return &APIKeyModule{service: apiKey_service}
}

func (m *APIKeyModule) Name() string {
	return "api_key"
}

func (m *APIKeyModule) Service() any {
	return m.service
}
