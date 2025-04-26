package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
)

var ErrRestaurantOrUserNotFound = errors.New("restaurant or user not found")

type RestaurantUserService interface {
	GetRestaurantUserRoleID(ctx context.Context, restaurantID, userID uuid.UUID) (int16, error)
	GetAnyRestaurantUserLinkByUserID(ctx context.Context, userID uuid.UUID) (repository.RestaurantUser, error)
}

type restaurantUserService struct {
	db *database.DB
}

func NewRestaurantUserService(db *database.DB) *restaurantUserService {
	return &restaurantUserService{
		db: db,
	}
}

func (s *restaurantUserService) GetRestaurantUserRoleID(ctx context.Context, restaurantID, userID uuid.UUID) (int16, error) {
	role, err := s.db.Queries.GetRestaurantUserRoleID(ctx, repository.GetRestaurantUserRoleIDParams{
		RestaurantID: restaurantID,
		UserID:       userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrRestaurantOrUserNotFound
		}
		return 0, err
	}

	return role, nil
}

func (s *restaurantUserService) GetAnyRestaurantUserLinkByUserID(ctx context.Context, userID uuid.UUID) (repository.RestaurantUser, error) {
	restaurantUser, err := s.db.Queries.GetAnyRestaurantUserLinkByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.RestaurantUser{}, ErrRestaurantOrUserNotFound
		}
		return repository.RestaurantUser{}, err
	}

	return restaurantUser, nil
}
