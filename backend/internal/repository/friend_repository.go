package repository

import (
	"context"
	"errors"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db"
	sqlc "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/friend"
	"github.com/jackc/pgx/v5"
)

type friendRepository struct {
	friend_repo db.FriendDB
}

func NewFriendRepository(db db.FriendDB) FriendRepository {
	return &friendRepository{friend_repo: db}
}

func (fr *friendRepository) GetInfoRelationship(ctx context.Context, params sqlc.GetInfoRelationshipParams) (sqlc.GetInfoRelationshipRow, error) {
	info, err := fr.friend_repo.DB.GetInfoRelationship(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.GetInfoRelationshipRow{}, ErrSameUser
		}

		return sqlc.GetInfoRelationshipRow{}, err
	}

	return info, nil
}
