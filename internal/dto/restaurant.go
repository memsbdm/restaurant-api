package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
)

type RestaurantDTO struct {
	ID          uuid.UUID             `json:"id"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Name        string                `json:"name"`
	Alias       string                `json:"alias"`
	Description *string               `json:"description"`
	Address     string                `json:"address"`
	Lat         *float64              `json:"lat"`
	Lng         *float64              `json:"lng"`
	Phone       *string               `json:"phone"`
	ImageURL    *string               `json:"image_url"`
	IsVerified  bool                  `json:"is_verified"`
	PlaceID     string                `json:"place_id"`
	Menus       []MenuDTO             `json:"menus,omitempty"`
	Categories  []CategoryDTO         `json:"categories,omitempty"`
	Articles    []ArticleDTO          `json:"articles,omitempty"`
	Invites     []RestaurantInviteDTO `json:"invites,omitempty"`
}

func NewRestaurantDTO(restaurant repository.Restaurant) RestaurantDTO {
	return RestaurantDTO{
		ID:          restaurant.ID,
		CreatedAt:   restaurant.CreatedAt,
		UpdatedAt:   restaurant.UpdatedAt,
		Name:        restaurant.Name,
		Alias:       restaurant.Alias,
		Description: restaurant.Description,
		Address:     restaurant.Address,
		Phone:       restaurant.Phone,
		ImageURL:    restaurant.ImageUrl,
		IsVerified:  restaurant.IsVerified,
		PlaceID:     restaurant.PlaceID,
	}
}

type CreateRestaurantDTO struct {
	Name    string
	Alias   string
	Address string
	Lat     *float64
	Lng     *float64
	Phone   *string
	PlaceID string
}

func (r CreateRestaurantDTO) ToParams() repository.CreateRestaurantParams {
	return repository.CreateRestaurantParams{
		Name:    r.Name,
		Alias:   r.Alias,
		Address: r.Address,
		Lat:     r.Lat,
		Lng:     r.Lng,
		Phone:   r.Phone,
		PlaceID: r.PlaceID,
	}
}
