package service

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
)

type UserService interface {
	GetAllUsers(search string, page, limit int) (map[string]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByUUID(uuid string) (models.User, error)
	UpdateUser(uuid string, user models.User) (models.User, error)
	DeleteUser(uuid string) error
}
