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
	Logout(ctx context.Context, arg sqlc.RevokeSessionParams) error
	LogoutAll(ctx context.Context, userId int64) error
}

type UserRepository interface {
	CreateProfile(ctx context.Context, arg sqlc.CreateProfileParams) (sqlc.Profile, error)
}
