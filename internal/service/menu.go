package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
	"github.com/memsbdm/restaurant-api/internal/dto"
)

type MenuService interface {
	Create(ctx context.Context, name string, restaurantID uuid.UUID) (*dto.Menu, error)
}

type menuService struct {
	db *database.DB
}

func NewMenuService(db *database.DB) *menuService {
	return &menuService{
		db: db,
	}
}

func (s *menuService) Create(ctx context.Context, name string, restaurantID uuid.UUID) (*dto.Menu, error) {
	menuAlreadyExists, err := s.db.Queries.MenuExistsForRestaurantID(ctx, restaurantID)
	if err != nil {
		log.Printf("Error checking if menu exists for restaurant ID %s: %v", restaurantID, err)
		return nil, err
	}

	dbCreatedMenu, err := s.db.Queries.CreateMenu(ctx, repository.CreateMenuParams{
		Name:         name,
		IsActive:     !menuAlreadyExists,
		RestaurantID: restaurantID,
	})
	if err != nil {
		log.Printf("Error creating menu for restaurant ID %s: %v", restaurantID, err)
		return nil, err
	}
	return dto.NewMenu(&dbCreatedMenu), nil
}
