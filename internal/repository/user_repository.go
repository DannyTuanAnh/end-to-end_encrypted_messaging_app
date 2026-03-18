package repository

import (
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type InMemoryUserRepository struct {
	users map[string]models.User
}

func NewInMemoryUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]models.User, 0),
	}
}

// FindAll is a repository function that retrieves all users from the in-memory store.
func (repo *InMemoryUserRepository) FindAll() (map[string]models.User, error) {
	if repo.users == nil {
		return nil, utils.NewError("user repository is not initialized", utils.ErrCodeInternal)
	}
	return repo.users, nil
}

// FindByUUID is a repository function that retrieves a user by their UUID from the in-memory store.
func (repo *InMemoryUserRepository) FindByUUID(uuid string) (models.User, bool) {
	user, exist := repo.users[uuid]
	if !exist {
		return models.User{}, false
	}

	return user, true
}

// FindByEmail is a repository function that retrieves a user by their email from the in-memory store.
func (repo *InMemoryUserRepository) FindByEmail(email string) (models.User, bool) {
	for _, user := range repo.users {
		if user.Email == email {
			return user, true
		}
	}

	return models.User{}, false
}

// Create is a repository function that adds a new user to the in-memory store.
func (repo *InMemoryUserRepository) Create(user models.User) error {
	if _, exist := repo.users[user.UUID]; exist {
		return utils.NewError("user with this UUID already exists", utils.ErrCodeConflict)
	}

	repo.users[user.UUID] = user

	return nil
}

func (repo *InMemoryUserRepository) Update(uuid string, user models.User) error {
	if _, exist := repo.users[uuid]; !exist {
		return utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	repo.users[uuid] = user

	return nil
}
func (repo *InMemoryUserRepository) Delete(uuid string) error {
	if _, exist := repo.users[uuid]; !exist {
		return utils.NewError("user not found", utils.ErrCodeNotFound)
	}
	delete(repo.users, uuid)

	return nil
}
