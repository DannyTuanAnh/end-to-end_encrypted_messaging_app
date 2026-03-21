package repository

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
)

type APIKeyRepository interface {
	CreateAPIKey(ctx context.Context, keyHash string) error
	RevokeAPIKey(ctx context.Context, keyHash string) error
	RevokeAll(ctx context.Context) error
}

type UserRepository interface {
	FindAll() (map[string]models.User, error)
	Create(user models.User) error
	FindByUUID(uuid string) (models.User, bool)
	Update(uuid string, user models.User) error
	Delete(uuid string) error
	FindByEmail(email string) (models.User, bool)
}
