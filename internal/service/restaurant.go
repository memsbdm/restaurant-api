package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/database/enum"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
	"github.com/memsbdm/restaurant-api/internal/dto"
)

var ErrRestaurantAlreadyTaken = errors.New("restaurant already taken")

type RestaurantService interface {
	Create(ctx context.Context, placeID string, userID uuid.UUID) (dto.Restaurant, error)
}

type restaurantService struct {
	db        *database.DB
	googleSvc GoogleService
}

func NewRestaurantService(db *database.DB, googleSvc GoogleService) RestaurantService {
	return &restaurantService{
		db:        db,
		googleSvc: googleSvc,
	}
}

func (s *restaurantService) Create(ctx context.Context, placeID string, userID uuid.UUID) (dto.Restaurant, error) {
	createRestaurantDTO, err := s.googleSvc.GetDetails(ctx, placeID)
	if err != nil {
		return dto.Restaurant{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return dto.Restaurant{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.db.Queries.WithTx(tx)

	taken, err := qtx.IsRestaurantAlreadyTaken(ctx, repository.IsRestaurantAlreadyTakenParams{
		PlaceID: placeID,
		UserID:  userID,
	})
	if err != nil {
		return dto.Restaurant{}, err
	}
	if taken {
		return dto.Restaurant{}, ErrRestaurantAlreadyTaken
	}

	restaurant, err := qtx.CreateRestaurant(ctx, createRestaurantDTO.ToParams())
	if err != nil {
		return dto.Restaurant{}, err
	}

	err = qtx.AddRestaurantUser(ctx, repository.AddRestaurantUserParams{
		RestaurantID: restaurant.ID,
		UserID:       userID,
		RoleID:       int16(enum.RoleAdmin),
	})
	if err != nil {
		return dto.Restaurant{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.Restaurant{}, err
	}

	return dto.NewRestaurant(restaurant), nil
}
