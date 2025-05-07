package apperr

import (
	"errors"
)

var (
	ErrUnAuthenticated    = errors.New("unauthenticated request")
	ErrUnAuthorized       = errors.New("unauthorized request")
	ErrInternalError      = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource conflict")
	ErrForbidden          = errors.New("forbidden request")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidUser        = errors.New("invalid user")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrMissingCredentials = errors.New("missing credentials for authentication")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAPIKey      = errors.New("invalid api key")
	ErrInvalidAPISecret   = errors.New("invalid api secret")
)
