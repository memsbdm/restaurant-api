package keys

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ContextKey string

const (
	UserIDContextKey       ContextKey = "userID"
	AuthOATContextKey      ContextKey = "authOAT"
	RestaurantIDContextKey ContextKey = "restaurantID"
	RestaurantContextKey   ContextKey = "restaurant"
	UserRoleIDContextKey   ContextKey = "userRoleID"
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

func GetRestaurantIDFromContext(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(RestaurantIDContextKey)
	if val == nil {
		return uuid.Nil, errors.New("restaurant ID not found in context")
	}

	return val.(uuid.UUID), nil
}

func GetUserRoleIDFromContext(ctx context.Context) (int16, error) {
	val := ctx.Value(UserRoleIDContextKey)
	if val == nil {
		return 0, errors.New("user role ID not found in context")
	}

	return val.(int16), nil
}
