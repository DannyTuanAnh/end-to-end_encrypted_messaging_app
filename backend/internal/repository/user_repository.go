package repository

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db"
	sqlc "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/user"
)

type userRepository struct {
	user_repo db.UserDB
}

func NewUserRepository(db db.UserDB) UserRepository {
	return &userRepository{user_repo: db}
}

func (ur *userRepository) CreateUser(ctx context.Context, displayName string) (int64, error) {
	userID, err := ur.user_repo.DB.CreateUser(ctx, displayName)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (ur *userRepository) DeleteUserByUserID(ctx context.Context, userId int64) error {
	result, err := ur.user_repo.DB.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCannotDeleteUser
	}

	return nil
}

func (ur *userRepository) ActiveUser(ctx context.Context, userId int64) error {
	result, err := ur.user_repo.DB.ActiveUser(ctx, userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCannotRestoreUser
	}

	return nil
}

func (ur *userRepository) IsExistProfile(ctx context.Context, userId int64) (bool, error) {
	exists, err := ur.user_repo.DB.IsExistProfile(ctx, userId)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ur *userRepository) GetProfileByUserID(ctx context.Context, userId int64) (sqlc.GetProfileByUserIdRow, error) {
	return ur.user_repo.DB.GetProfileByUserId(ctx, userId)
}

func (ur *userRepository) GetProfileByUserUUID(ctx context.Context, arg sqlc.GetProfileByUserUUIDParams) (sqlc.GetProfileByUserUUIDRow, error) {
	return ur.user_repo.DB.GetProfileByUserUUID(ctx, arg)
}

func (ur *userRepository) CreateProfile(ctx context.Context, arg sqlc.CreateProfileParams) (sqlc.Profile, error) {
	return ur.user_repo.DB.CreateProfile(ctx, arg)
}

func (ur *userRepository) UpdateProfile(ctx context.Context, arg sqlc.UpdateProfileByUserIdParams) (sqlc.UpdateProfileByUserIdRow, error) {
	return ur.user_repo.DB.UpdateProfileByUserId(ctx, arg)
}

func (ur *userRepository) DisableUserByUserID(ctx context.Context, userId int64) error {
	return ur.user_repo.DB.DisableUser(ctx, userId)
}
