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

func (r *userRepository) CreateProfile(ctx context.Context, arg sqlc.CreateProfileParams) (sqlc.Profile, error) {
	return r.user_repo.CreateProfile(ctx, arg)
}
