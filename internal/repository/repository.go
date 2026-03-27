package repository

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
)

type APIKeyRepository interface {
	CreateAPIKey(ctx context.Context, keyHash string) error
	RevokeAPIKey(ctx context.Context, keyHash string) error
	RevokeAll(ctx context.Context) error
}

type AuthRepository interface {
	Login(ctx context.Context, arg sqlc.OAuthLoginParams) (models.GoogleLoginResponse, error)
}

type UserRepository interface {
	CreateProfile(ctx context.Context, arg sqlc.CreateProfileParams) (sqlc.Profile, error)
}
