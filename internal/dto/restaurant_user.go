package dto

import (
	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
)

type RestaurantUser struct {
	ID           int       `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	UserID       uuid.UUID `json:"user_id"`
	RoleID       int       `json:"role_id"`
}

func NewRestaurantUser(restaurantUser *repository.RestaurantUser) *RestaurantUser {
	return &RestaurantUser{
		ID:           int(restaurantUser.ID),
		RestaurantID: restaurantUser.RestaurantID,
		UserID:       restaurantUser.UserID,
		RoleID:       int(restaurantUser.RoleID),
	}
}
