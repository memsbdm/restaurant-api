package dto

import "github.com/google/uuid"

type MenuDTO struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	MenuOrder    int           `json:"menu_order"`
	IsActive     bool          `json:"is_active"`
	RestaurantID uuid.UUID     `json:"restaurant_id"`
	Restaurant   RestaurantDTO `json:"restaurant"`
}
