package dto

import "github.com/google/uuid"

type Category struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	CategoryOrder int        `json:"category_order"`
	IsDefault     bool       `json:"is_default"`
	MenuID        int        `json:"menu_id"`
	RestaurantID  uuid.UUID  `json:"restaurant_id"`
	Menu          Menu       `json:"menu"`
	Restaurant    Restaurant `json:"restaurant"`
	Articles      []Article  `json:"articles"`
}
