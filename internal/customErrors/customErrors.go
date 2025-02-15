package customErrors

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidData        = errors.New("invalid data")
	ErrNotEnoughCoins     = errors.New("not enough coins")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidRequest     = errors.New("invalid request body")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrISE                = errors.New("internal server error")
)
