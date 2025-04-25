package dto

import "github.com/google/uuid"

type CategoryDTO struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Description   *string       `json:"description"`
	CategoryOrder int           `json:"category_order"`
	IsDefault     bool          `json:"is_default"`
	MenuID        int           `json:"menu_id"`
	RestaurantID  uuid.UUID     `json:"restaurant_id"`
	Menu          MenuDTO       `json:"menu"`
	Restaurant    RestaurantDTO `json:"restaurant"`
	Articles      []ArticleDTO  `json:"articles"`
}
