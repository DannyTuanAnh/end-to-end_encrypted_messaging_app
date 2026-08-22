package repository

import "errors"

var (
	ErrIdentityAlreadyExist = errors.New("identity already exists")
	ErrNotFoundIdentityID   = errors.New("identity not found")
	ErrCannotRestore        = errors.New("cannot restore user, user has been disabled for more than 30 days")
	ErrCannotDeleteUser     = errors.New("cannot delete user, user does not exist or has already been deleted")
	ErrNotFoundSessionID    = errors.New("session not found")
)
