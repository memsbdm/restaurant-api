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
	Create(ctx context.Context, placeID string, userID uuid.UUID) (dto.RestaurantDTO, error)
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

func (s *restaurantService) Create(ctx context.Context, placeID string, userID uuid.UUID) (dto.RestaurantDTO, error) {
	createRestaurantDTO, err := s.googleSvc.GetDetails(ctx, placeID)
	if err != nil {
		return dto.RestaurantDTO{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return dto.RestaurantDTO{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.db.Queries.WithTx(tx)

	taken, err := qtx.IsRestaurantAlreadyTaken(ctx, repository.IsRestaurantAlreadyTakenParams{
		PlaceID: placeID,
		UserID:  userID,
	})
	if err != nil {
		return dto.RestaurantDTO{}, err
	}
	if taken {
		return dto.RestaurantDTO{}, ErrRestaurantAlreadyTaken
	}

	restaurant, err := qtx.CreateRestaurant(ctx, createRestaurantDTO.ToParams())
	if err != nil {
		return dto.RestaurantDTO{}, err
	}

	err = qtx.AddRestaurantUser(ctx, repository.AddRestaurantUserParams{
		RestaurantID: restaurant.ID,
		UserID:       userID,
		RoleID:       int16(enum.RoleAdmin),
	})
	if err != nil {
		return dto.RestaurantDTO{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.RestaurantDTO{}, err
	}

	return dto.NewRestaurantDTO(restaurant), nil
}
