package service

import (
	"context"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
)

type APIKeyService interface {
	CreateAPIKey(ctx context.Context, args *models.GenerateAPIKeyArgs) error
	RevokeAPIKey(ctx context.Context, keyID string) error
	RevokeAll(ctx context.Context) error
}

type UserService interface {
	GetAllUsers(search string, page, limit int) (map[string]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByUUID(uuid string) (models.User, error)
	UpdateUser(uuid string, user models.User) (models.User, error)
	DeleteUser(uuid string) error
}
