package session

import "errors"

var (
	ErrSessionNotExists          = errors.New("repository.Session: session doesn't exists")
	ErrRefreshTokenAlreadyExists = errors.New("repository.Session: this refresh token already exists")
)

func IsErrSessionNotExists(err error) bool {
	return errors.Is(err, ErrSessionNotExists)
}

func IsErrRefreshTokenAlreadyExists(err error) bool {
	return errors.Is(err, ErrRefreshTokenAlreadyExists)
}
