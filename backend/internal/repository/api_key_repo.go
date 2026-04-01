package repository

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
)

type apiKeyRepository struct {
	db sqlc.Querier
}

func NewAPIKeyRepository(db sqlc.Querier) APIKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) CreateAPIKey(ctx context.Context, keyHash string) error {
	if err := r.db.CreateAPIKey(ctx, keyHash); err != nil {
		return err
	}

	return nil
}

func (r *apiKeyRepository) RevokeAPIKey(ctx context.Context, keyHash string) error {
	if err := r.db.RevokeAPIKeyByKey(ctx, keyHash); err != nil {
		return err
	}

	return nil
}

func (r *apiKeyRepository) RevokeAll(ctx context.Context) error {
	if err := r.db.RevokeAllAPIKeys(ctx); err != nil {
		return err
	}

	return nil
}
