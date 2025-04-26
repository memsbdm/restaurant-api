package response

import (
	"errors"
	"net/http"

	"github.com/memsbdm/restaurant-api/internal/service"
)

var (
	// General
	ErrBadRequest         = errors.New("bad request")
	ErrForbidden          = errors.New("forbidden request")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrInternal           = errors.New("internal error")
	ErrServiceUnavailable = errors.New("service unavailable")

	// Middleware
	ErrNoRestaurantFoundForUser = errors.New("no restaurant found for user")
)

var ErrToHttpStatusCode = map[error]int{
	// General
	ErrBadRequest:         http.StatusBadRequest,
	ErrForbidden:          http.StatusForbidden,
	ErrUnauthorized:       http.StatusUnauthorized,
	ErrInternal:           http.StatusInternalServerError,
	ErrServiceUnavailable: http.StatusServiceUnavailable,

	// Conflict
	service.ErrEmailConflict:          http.StatusConflict,
	service.ErrEmailAlreadyVerified:   http.StatusForbidden,
	service.ErrRestaurantAlreadyTaken: http.StatusConflict,

	// Token
	service.ErrInvalidToken: http.StatusBadRequest,

	// Auth
	service.ErrInvalidCredentials: http.StatusUnauthorized,

	// Restaurant
	service.ErrNoRestaurantFoundForUser: http.StatusForbidden,

	// Mailer
	service.ErrMailerUnavailable: http.StatusServiceUnavailable,

	// Google
	service.ErrGoogleServiceUnavailable:      http.StatusServiceUnavailable,
	service.ErrGoogleAutocompleteQueryLength: http.StatusBadRequest,
	service.ErrGoogleInvalidPlaceID:          http.StatusBadRequest,
}
