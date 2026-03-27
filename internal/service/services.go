package service

import (
	"context"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
)

type APIKeyService interface {
	CreateAPIKey(ctx context.Context, args *models.GenerateAPIKeyArgs) error
	RevokeAPIKey(ctx context.Context, keyID string) error
	RevokeAll(ctx context.Context) error
}
