package validation

import (
	"errors"
	"fmt"
)

const (
	UserPasswordMinLength = 8
	UserNameMaxLength     = 50
)

// Required
var (
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
)

// Min
var (
	ErrPasswordTooShort = fmt.Errorf("password should contain at least %d characters", UserPasswordMinLength)
	ErrUserNameTooLong  = fmt.Errorf("name should contain at most %d characters", UserNameMaxLength)
)

// errorMessages holds custom error messages for specific validation failures.
var errorMessages = map[string]error{
	// Required
	"registerUserRequest.Name.notblank":     ErrNameRequired,
	"registerUserRequest.Email.notblank":    ErrEmailRequired,
	"registerUserRequest.Password.notblank": ErrPasswordRequired,
	"loginUserRequest.Email.notblank":       ErrEmailRequired,
	"loginUserRequest.Password.notblank":    ErrEmailRequired,

	// Min
	"registerUserRequest.Password.min": ErrPasswordTooShort,

	// Max
	"registerUserRequest.Name.max": ErrUserNameTooLong,

	// Email
	"registerUserRequest.Email.email": ErrInvalidEmail,
	"loginUserRequest.Email.email":    ErrInvalidEmail,
}
