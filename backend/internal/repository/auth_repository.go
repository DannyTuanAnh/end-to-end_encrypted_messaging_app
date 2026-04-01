package repository

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
)

type authRepository struct {
	auth_repo sqlc.Querier
}

func NewAuthRepository(db sqlc.Querier) AuthRepository {
	return &authRepository{auth_repo: db}
}

func (ar *authRepository) Login(ctx context.Context, arg sqlc.OAuthLoginParams) (models.GoogleLoginResponse, error) {
	result, err := ar.auth_repo.OAuthLogin(ctx, arg)
	if err != nil {
		return models.GoogleLoginResponse{}, err
	}

	return models.GoogleLoginResponse{
		SessionId:     result.FSessionID,
		UserId:        result.FUserID,
		ProfileExists: result.FProfileExists,
	}, nil

}

func (ar *authRepository) Logout(ctx context.Context, arg sqlc.RevokeSessionParams) error {
	err := ar.auth_repo.RevokeSession(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}

func (ar *authRepository) LogoutAll(ctx context.Context, userId int64) error {
	err := ar.auth_repo.RevokeAllSessions(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}
