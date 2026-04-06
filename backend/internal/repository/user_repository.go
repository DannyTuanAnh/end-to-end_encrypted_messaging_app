package repository

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
)

type userRepository struct {
	user_repo sqlc.Querier
}

func NewUserRepository(db sqlc.Querier) UserRepository {
	return &userRepository{user_repo: db}
}

func (r *userRepository) GetProfileByUserID(ctx context.Context, userId int64) (sqlc.GetProfileByUserIdRow, error) {
	return r.user_repo.GetProfileByUserId(ctx, userId)
}

func (r *userRepository) GetProfileByUserUUID(ctx context.Context, arg sqlc.GetProfileByUserUUIDParams) (sqlc.GetProfileByUserUUIDRow, error) {
	return r.user_repo.GetProfileByUserUUID(ctx, arg)
}

func (r *userRepository) CreateProfile(ctx context.Context, arg sqlc.CreateProfileParams) (sqlc.Profile, error) {
	return r.user_repo.CreateProfile(ctx, arg)
}

func (r *userRepository) DisableUserByUserID(ctx context.Context, userId int64) error {
	return r.user_repo.DisableUser(ctx, userId)
}
