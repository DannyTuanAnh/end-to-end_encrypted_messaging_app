package repository

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
)

type authRepository struct {
	auth_repo sqlc.Querier
}

func NewAuthRepository(db sqlc.Querier) AuthRepository {
	return &authRepository{auth_repo: db}
}

func (r *authRepository) Login(ctx context.Context, arg sqlc.OAuthLoginParams) (string, string, error) {
	return "", "", nil
}
