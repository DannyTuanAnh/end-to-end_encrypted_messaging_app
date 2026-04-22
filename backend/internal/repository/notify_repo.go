package repository

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
)

type notifyRepository struct {
}

func NewNotifyRepository(db sqlc.Querier) NotifyRepository {
	return &notifyRepository{}
}
