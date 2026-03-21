package service

import (
	"context"
	"log"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type apiKeyService struct {
	apiKey_repo repository.APIKeyRepository
}

func NewAPIKeyService(apiKey_repo repository.APIKeyRepository) APIKeyService {
	return &apiKeyService{apiKey_repo: apiKey_repo}
}

func (s *apiKeyService) CreateAPIKey(ctx context.Context, args *models.GenerateAPIKeyArgs) error {
	plaintext := utils.GenerateAPIKey()
	keyHash := utils.HashAPIKey(plaintext)

	utils.SaveKeyToEnv(args.ClientType, plaintext)

	if err := s.apiKey_repo.CreateAPIKey(ctx, keyHash); err != nil {
		return err
	}

	log.Println("API key created and saved to environment variable successfully")

	return nil
}

func (s *apiKeyService) RevokeAPIKey(ctx context.Context, keyID string) error {
	keyHash := utils.HashAPIKey(keyID)
	if err := s.apiKey_repo.RevokeAPIKey(ctx, keyHash); err != nil {
		return err
	}

	log.Printf("API key (%s) revoked successfully\n", keyID)

	return nil
}

func (s *apiKeyService) RevokeAll(ctx context.Context) error {
	if err := s.apiKey_repo.RevokeAll(ctx); err != nil {
		return err
	}

	log.Println("All API keys revoked successfully")

	return nil
}
