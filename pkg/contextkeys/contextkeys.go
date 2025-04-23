package contextkeys

import (
	"context"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/response"
)

type ContextKey string

const (
	UserIDContextKey    ContextKey = "userID"
	SignedOATContextKey ContextKey = "signedOAT"
)

func GetValueFromContext(ctx context.Context, key ContextKey) (string, error) {
	val := ctx.Value(key)
	if val == nil {
		return "", response.ErrInternal
	}

	return val.(string), nil
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(UserIDContextKey)
	if val == nil {
		return uuid.Nil, response.ErrUnauthorized
	}

	return uuid.MustParse(val.(string)), nil
}
