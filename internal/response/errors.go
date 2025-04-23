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

	// Token
	ErrInvalidToken = errors.New("invalid or expired token")

	// Conflict
	ErrEmailConflict = errors.New("email already taken")

	// Auth
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Not returned to the client
var (
	// Cache
	ErrCacheNotFound = errors.New("cache not found")
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

	// Token
	ErrInvalidToken: http.StatusBadRequest,

	// Auth
	ErrInvalidCredentials: http.StatusUnauthorized,
}
