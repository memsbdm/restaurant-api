package response

import (
	"errors"
	"net/http"
)

var (
	// General
	ErrBadRequest         = errors.New("bad request")
	ErrForbidden          = errors.New("forbidden request")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrInternal           = errors.New("internal error")
	ErrServiceUnavailable = errors.New("service unavailable")

	// Cache
	ErrCacheNotFound = errors.New("cache not found")

	// Conflict
	ErrEmailConflict = errors.New("email already taken")

	// Auth
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var ErrToHttpStatusCode = map[error]int{
	// General
	ErrBadRequest:         http.StatusBadRequest,
	ErrForbidden:          http.StatusForbidden,
	ErrUnauthorized:       http.StatusUnauthorized,
	ErrInternal:           http.StatusInternalServerError,
	ErrServiceUnavailable: http.StatusServiceUnavailable,

	// Conflict
	ErrEmailConflict: http.StatusConflict,

	// Auth
	ErrInvalidCredentials: http.StatusUnauthorized,
}
