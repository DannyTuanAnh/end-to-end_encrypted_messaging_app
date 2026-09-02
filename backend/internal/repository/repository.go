package repository

import (
	"context"

	sqlc_auth "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/auth"
	sqlc_friend "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/friend"
	sqlc_user "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/user"
	"github.com/google/uuid"
)

type APIKeyRepository interface {
	CreateAPIKey(ctx context.Context, keyHash string) error
	RevokeAPIKey(ctx context.Context, keyHash string) error
	RevokeAll(ctx context.Context) error
}

type AuthRepository interface {
	IsExistingIdentityID(ctx context.Context, arg sqlc_auth.FindExistingIdentityParams) (sqlc_auth.FindExistingIdentityRow, error)
	CreateIdentity(ctx context.Context, arg sqlc_auth.CreateIdentityParams) error
	ActiveIdentity(ctx context.Context, arg sqlc_auth.ActiveIdentityParams) error
	DisableIdentity(ctx context.Context, arg sqlc_auth.DisableIdentityParams) error
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
	GetProfile(ctx context.Context, userId int64) (sqlc_user.GetProfileRow, error)
	GetProfileByUserID(ctx context.Context, arg sqlc_user.GetProfileByUserIdParams) (sqlc_user.GetProfileByUserIdRow, error)
	GetUserByUUID(ctx context.Context, targetUserUUID uuid.UUID) (sqlc_user.GetUserByUUIDRow, error)
	CreateProfile(ctx context.Context, arg sqlc_user.CreateProfileParams) (sqlc_user.Profile, error)
	DisableUserByUserID(ctx context.Context, userId int64) error
	UpdateProfile(ctx context.Context, arg sqlc_user.UpdateProfileByUserIdParams) (sqlc_user.UpdateProfileByUserIdRow, error)
}

type FriendRepository interface {
	GetInfoRelationship(ctx context.Context, params sqlc_friend.GetInfoRelationshipParams) (sqlc_friend.GetInfoRelationshipRow, error)
}

type NotifyRepository interface {
}
