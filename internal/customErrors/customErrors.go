package customErrors

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidData        = errors.New("invalid data")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidRequest     = errors.New("invalid request body")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrISE                = errors.New("internal server error")
)
