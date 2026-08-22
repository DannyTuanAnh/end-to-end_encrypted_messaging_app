package repository

import (
	"context"

	sqlc_auth "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/auth"
	sqlc_user "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/user"
	"github.com/google/uuid"
)

type APIKeyRepository interface {
	CreateAPIKey(ctx context.Context, keyHash string) error
	RevokeAPIKey(ctx context.Context, keyHash string) error
	RevokeAll(ctx context.Context) error
}

type AuthRepository interface {
	IsExistingIdentityID(ctx context.Context, arg sqlc_auth.FindExistingIdentityParams) (int64, error)
	CreateIdentity(ctx context.Context, arg sqlc_auth.CreateIdentityParams) error
	CreateSession(ctx context.Context, userID int64) (uuid.UUID, error)
	CheckSession(ctx context.Context, sessionID uuid.UUID) (sqlc_auth.CheckSessionRow, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	LogoutAll(ctx context.Context, userId int64) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, displayName string) (int64, error)
	DeleteUserByUserID(ctx context.Context, userId int64) error
	ActiveUser(ctx context.Context, userId int64) error
	IsExistProfile(ctx context.Context, userId int64) (bool, error)
	GetProfileByUserID(ctx context.Context, userId int64) (sqlc_user.GetProfileByUserIdRow, error)
	GetProfileByUserUUID(ctx context.Context, arg sqlc_user.GetProfileByUserUUIDParams) (sqlc_user.GetProfileByUserUUIDRow, error)
	CreateProfile(ctx context.Context, arg sqlc_user.CreateProfileParams) (sqlc_user.Profile, error)
	DisableUserByUserID(ctx context.Context, userId int64) error
	UpdateProfile(ctx context.Context, arg sqlc_user.UpdateProfileByUserIdParams) (sqlc_user.UpdateProfileByUserIdRow, error)
}

type NotifyRepository interface {
}
