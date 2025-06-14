package dto

import "github.com/google/uuid"

type Article struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Price        float64    `json:"price"`
	ImageURL     *string    `json:"image_url"`
	ArticleOrder int        `json:"article_order"`
	CategoryID   string     `json:"category_id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	Category     Category   `json:"category"`
	Restaurant   Restaurant `json:"restaurant"`
}
