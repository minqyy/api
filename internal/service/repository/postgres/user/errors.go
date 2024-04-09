package user

import "errors"

var (
	ErrUserAlreadyExists = errors.New("repository.User: user already exists")
	ErrUserNotFound      = errors.New("repository.User: user not found")
)
