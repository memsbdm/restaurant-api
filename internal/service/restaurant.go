package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/database/enum"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
	"github.com/memsbdm/restaurant-api/internal/dto"
)

var (
	ErrRestaurantAlreadyTaken   = errors.New("restaurant already taken")
	ErrRestaurantNotFound       = errors.New("restaurant not found")
	ErrNoRestaurantFoundForUser = errors.New("no restaurant found for user")
)

type RestaurantService interface {
	Create(ctx context.Context, placeID string, userID uuid.UUID) (*dto.Restaurant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.Restaurant, error)
	GetRestaurantsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.Restaurant, error)
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

func (s *restaurantService) GetByID(ctx context.Context, id uuid.UUID) (*dto.Restaurant, error) {
	dbRestaurant, err := s.db.Queries.GetRestaurantByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRestaurantNotFound
		}
		return nil, fmt.Errorf("error fetching restaurant by ID %s: %w", id, err)
	}

	return dto.NewRestaurant(&dbRestaurant), nil
}

func (s *restaurantService) GetRestaurantsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.Restaurant, error) {
	dbRestaurants, err := s.db.Queries.GetRestaurantsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*dto.Restaurant{}, ErrNoRestaurantFoundForUser
		}
		return nil, fmt.Errorf("error fetching restaurants for user ID %s: %w", userID, err)
	}

	restaurants := make([]*dto.Restaurant, len(dbRestaurants))
	for i := range dbRestaurants {
		restaurants[i] = dto.NewRestaurant(&dbRestaurants[i])
	}
	return restaurants, nil
}

func (s *restaurantService) Create(ctx context.Context, placeID string, userID uuid.UUID) (*dto.Restaurant, error) {
	createRestaurantDTO, err := s.googleSvc.GetDetails(ctx, placeID)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.db.Queries.WithTx(tx)

	taken, err := qtx.IsRestaurantAlreadyTaken(ctx, repository.IsRestaurantAlreadyTakenParams{
		PlaceID: placeID,
		UserID:  userID,
	})
	if err != nil {
		return nil, fmt.Errorf("error checking if restaurant is already taken: %w", err)
	}
	if taken {
		return nil, ErrRestaurantAlreadyTaken
	}

	restaurant, err := qtx.CreateRestaurant(ctx, createRestaurantDTO.ToParams())
	if err != nil {
		return nil, fmt.Errorf("error creating restaurant: %w", err)
	}

	err = qtx.AddRestaurantUser(ctx, repository.AddRestaurantUserParams{
		RestaurantID: restaurant.ID,
		UserID:       userID,
		RoleID:       int16(enum.RoleOwner),
	})
	if err != nil {
		return nil, fmt.Errorf("error adding restaurant user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return dto.NewRestaurant(&restaurant), nil
}
