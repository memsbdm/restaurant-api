package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/memsbdm/restaurant-api/internal/validation"
)

func HandleSuccess(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data == nil {
		return
	}
	json.NewEncoder(w).Encode(data)
}

func HandleValidationError(w http.ResponseWriter, errs []validation.ValidationError, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid JSON")
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	resp := struct {
		Errors []validation.ValidationError `json:"errors"`
	}{Errors: errs}

	json.NewEncoder(w).Encode(resp)
}

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	statusCode, ok := ErrToHttpStatusCode[err]
	if !ok {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusServiceUnavailable
			err = ErrServiceUnavailable
		case errors.Is(err, context.Canceled):
			statusCode = http.StatusBadRequest
			err = ErrBadRequest
		default:
			statusCode = http.StatusInternalServerError
			err = ErrInternal
		}
	}

	resp := struct {
		Errors []string `json:"errors"`
	}{
		[]string{err.Error()},
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
