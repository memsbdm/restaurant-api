package keys

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ContextKey string

const (
	UserIDContextKey  ContextKey = "userID"
	AuthOATContextKey ContextKey = "authOAT"
)

func GetValueFromContext(ctx context.Context, key ContextKey) (string, error) {
	val := ctx.Value(key)
	if val == nil {
		return "", errors.New("value not found in context")
	}

	return val.(string), nil
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(UserIDContextKey)
	if val == nil {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	return uuid.MustParse(val.(string)), nil
}
