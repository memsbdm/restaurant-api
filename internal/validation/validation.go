package validation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

const (
	// MaxRequestSize defines the maximum allowed size for request bodies (1MB)
	MaxRequestSize = 1 << 20
)

var (
	Validate *validator.Validate
	once     sync.Once
)

func init() {
	once.Do(func() {
		Validate = validator.New(validator.WithRequiredStructEnabled())
		if err := Validate.RegisterValidation("notblank", notBlank); err != nil {
			log.Printf("failed to register notblank validation: %v", err)
		}
	})
}

// notBlank validates that the string length is greater than 0 after trimming whitespace.
func notBlank(fl validator.FieldLevel) bool {
	return len(strings.TrimSpace(fl.Field().String())) > 0
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateRequest takes a payload from an HTTP request and verifies it.
func ValidateRequest(w http.ResponseWriter, r *http.Request, payload any) ([]ValidationError, error) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}()

	decoder := json.NewDecoder(r.Body)

	// Validate JSON format
	if err := decoder.Decode(&payload); err != nil {
		return nil, err
	}

	// Validate payload
	var errs []ValidationError
	if err := Validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := fmt.Sprintf("%s.%s", err.StructNamespace(), err.Tag())

			message, ok := errorMessages[field]
			if !ok {
				message = fmt.Errorf("validation failed on field '%s' for condition '%s'", err.Field(), err.Tag())
			}
			validationErr := ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: message.Error(),
			}

			errs = append(errs, validationErr)
		}
	}
	return errs, nil
}
