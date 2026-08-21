package repository

import (
	"context"

	sqlc_auth "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/auth"
	sqlc_user "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/user"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
)

type APIKeyRepository interface {
	CreateAPIKey(ctx context.Context, keyHash string) error
	RevokeAPIKey(ctx context.Context, keyHash string) error
	RevokeAll(ctx context.Context) error
}

type AuthRepository interface {
	Login(ctx context.Context, arg sqlc_auth.OAuthLoginParams) (models.GoogleLoginResponse, error)
	Logout(ctx context.Context, arg sqlc_auth.RevokeSessionParams) error
	LogoutAll(ctx context.Context, userId int64) error
}

type UserRepository interface {
	GetProfileByUserID(ctx context.Context, userId int64) (sqlc_user.GetProfileByUserIdRow, error)
	GetProfileByUserUUID(ctx context.Context, arg sqlc_user.GetProfileByUserUUIDParams) (sqlc_user.GetProfileByUserUUIDRow, error)
	CreateProfile(ctx context.Context, arg sqlc_user.CreateProfileParams) (sqlc_user.Profile, error)
	DisableUserByUserID(ctx context.Context, userId int64) error
	UpdateProfile(ctx context.Context, arg sqlc_user.UpdateProfileByUserIdParams) (sqlc_user.UpdateProfileByUserIdRow, error)
}

type NotifyRepository interface {
}
