package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
	"github.com/memsbdm/restaurant-api/internal/dto"
)

var ErrRestaurantOrUserNotFound = errors.New("restaurant or user not found")

type RestaurantUserService interface {
	GetRestaurantUserRoleID(ctx context.Context, restaurantID, userID uuid.UUID) (int, error)
	GetAnyRestaurantUserLinkByUserID(ctx context.Context, userID uuid.UUID) (*dto.RestaurantUser, error)
}

type restaurantUserService struct {
	db *database.DB
}

func NewRestaurantUserService(db *database.DB) *restaurantUserService {
	return &restaurantUserService{
		db: db,
	}
}

func (s *restaurantUserService) GetRestaurantUserRoleID(ctx context.Context, restaurantID, userID uuid.UUID) (int, error) {
	role, err := s.db.Queries.GetRestaurantUserRoleID(ctx, repository.GetRestaurantUserRoleIDParams{
		RestaurantID: restaurantID,
		UserID:       userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrRestaurantOrUserNotFound
		}
		return 0, fmt.Errorf("error fetching restaurant user role ID for restaurant ID %s and user ID %s: %w", restaurantID, userID, err)
	}

	return int(role), nil
}

func (s *restaurantUserService) GetAnyRestaurantUserLinkByUserID(ctx context.Context, userID uuid.UUID) (*dto.RestaurantUser, error) {
	dbRestaurantUser, err := s.db.Queries.GetAnyRestaurantUserLinkByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRestaurantOrUserNotFound
		}
		return nil, fmt.Errorf("error fetching restaurant user link by user ID %s: %w", userID, err)
	}

	return dto.NewRestaurantUser(&dbRestaurantUser), nil
}
