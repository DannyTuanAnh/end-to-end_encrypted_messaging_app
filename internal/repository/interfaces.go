package repository

import "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"

type UserRepository interface {
	FindAll() (map[string]models.User, error)
	Create(user models.User) error
	FindByUUID(uuid string) (models.User, bool)
	Update(uuid string, user models.User) error
	Delete(uuid string) error
	FindByEmail(email string) (models.User, bool)
}
