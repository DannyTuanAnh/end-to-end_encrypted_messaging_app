package service

import (
	"strings"

	"github.com/google/uuid"
	"github.com/user-manage/internal/models"
	"github.com/user-manage/internal/repository"
	"github.com/user-manage/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	user_repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		user_repo: repo,
	}
}

// GetAllUsers is a service function that retrieves all users from the repository.
func (us *userService) GetAllUsers(search string, page, limit int) (map[string]models.User, error) {
	userMap, err := us.user_repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(err, "failed to get all users", utils.ErrCodeInternal)
	}

	return userMap, nil
}

// GetUserByUUID is a service function that retrieves a user by their UUID from the repository.
func (us *userService) GetUserByUUID(uuid string) (models.User, error) {
	if uuid == "" {
		return models.User{}, utils.NewError("UUID cannot be empty", utils.ErrCodeBadRequest)
	}

	user, exist := us.user_repo.FindByUUID(uuid)
	if !exist {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	return user, nil
}

// CreateUser is a service function that handles the business logic for creating a new user.
func (us *userService) CreateUser(user models.User) (models.User, error) {
	// 1. Normalize email before saving to database
	user.Email = utils.NormalizeString(user.Email)

	// 2. Check if email already exists in database
	if _, exist := us.user_repo.FindByEmail(user.Email); exist {
		return models.User{}, utils.NewError("email already exists", utils.ErrCodeConflict)
	}

	// 3. Generate UUID for new user
	user.UUID = uuid.New().String()
	// 4. Hash password before saving to database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
	}

	// 5. Save user to database
	user.Password = string(hashedPassword)
	if err := us.user_repo.Create(user); err != nil {
		return models.User{}, utils.WrapError(err, "failed to create user", utils.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) UpdateUser(uuid string, user models.User) (models.User, error) {
	currentUser, exist := us.user_repo.FindByUUID(uuid)
	if !exist {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	user.Email = utils.NormalizeString(user.Email)
	if currentUser, exist := us.user_repo.FindByEmail(user.Email); exist && currentUser.UUID != uuid {
		return models.User{}, utils.NewError("email already exists", utils.ErrCodeConflict)
	}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
		}
		currentUser.Password = string(hashedPassword)
	}

	if user.Name != "" {
		currentUser.Name = strings.TrimSpace(user.Name)
	}

	if user.Age != 0 {
		currentUser.Age = user.Age
	}

	if user.Status != 0 {
		currentUser.Status = user.Status
	}

	if user.Level != 0 {
		currentUser.Level = user.Level
	}

	if err := us.user_repo.Update(uuid, currentUser); err != nil {
		return models.User{}, utils.WrapError(err, "failed to update user", utils.ErrCodeInternal)
	}

	return currentUser, nil
}
func (us *userService) DeleteUser(uuid string) error {
	if _, exist := us.user_repo.FindByUUID(uuid); !exist {
		return utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	if err := us.user_repo.Delete(uuid); err != nil {
		return utils.WrapError(err, "failed to delete user", utils.ErrCodeInternal)
	}

	return nil
}
