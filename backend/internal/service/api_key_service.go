package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/redis/go-redis/v9"
)

type apiKeyService struct {
	apiKey_repo  repository.APIKeyRepository
	redis_memory *redis.Client
}

func NewAPIKeyService(apiKey_repo repository.APIKeyRepository, redis_memory *redis.Client) APIKeyService {
	return &apiKeyService{
		apiKey_repo:  apiKey_repo,
		redis_memory: redis_memory,
	}
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

	gen, err := s.redis_memory.Get(ctx, "generation_api_key").Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return fmt.Errorf("Failed to get generation number from Redis: %v", err)
		}
	}
	genNum, err := strconv.Atoi(gen)
	if err != nil {
		return fmt.Errorf("Invalid generation number: %v", err)
	}

	cacheKey := fmt.Sprintf("api_key:%d:%s", genNum, keyHash)

	if err := s.redis_memory.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to delete API key from Redis cache: %v", err)
	}

	log.Printf("API key (%s) revoked successfully\n", keyID)

	return nil
}

func (s *apiKeyService) RevokeAll(ctx context.Context) error {
	if err := s.apiKey_repo.RevokeAll(ctx); err != nil {
		return err
	}

	err := s.redis_memory.Incr(ctx, "generation_api_key").Err()
	if err != nil {
		return fmt.Errorf("Failed to increment generation number in Redis: %v", err)
	}

	log.Println("All API keys revoked successfully")

	return nil
}
