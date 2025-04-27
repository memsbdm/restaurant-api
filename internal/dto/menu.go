package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
)

type Menu struct {
	ID           int         `json:"id"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Name         string      `json:"name"`
	IsActive     bool        `json:"is_active"`
	RestaurantID uuid.UUID   `json:"restaurant_id"`
	Restaurant   *Restaurant `json:"restaurant,omitempty"`
}

func NewMenu(menu *repository.Menu) *Menu {
	return &Menu{
		ID:           int(menu.ID),
		CreatedAt:    menu.CreatedAt,
		UpdatedAt:    menu.UpdatedAt,
		Name:         menu.Name,
		IsActive:     menu.IsActive,
		RestaurantID: menu.RestaurantID,
	}
}
