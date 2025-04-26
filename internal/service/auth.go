package service

import (
	"context"
	"errors"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/cache"
	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/pkg/keys"
	"github.com/memsbdm/restaurant-api/pkg/security"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService interface {
	Register(ctx context.Context, user *dto.CreateUser) (*dto.User, string, error)
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, string, error)
	Logout(ctx context.Context, oat string) error
}

type authService struct {
	cfg               *config.Security
	cache             cache.Cache
	restaurantService RestaurantService
	tokenService      TokenService
	userService       UserService
}

func NewAuthService(cfg *config.Security, cache cache.Cache, userService UserService, tokenService TokenService, restaurantService RestaurantService) *authService {
	return &authService{
		cfg:               cfg,
		cache:             cache,
		restaurantService: restaurantService,
		tokenService:      tokenService,
		userService:       userService,
	}
}

func (s *authService) Register(ctx context.Context, user *dto.CreateUser) (*dto.User, string, error) {
	createdUser, err := s.userService.Create(ctx, user)
	if err != nil {
		return nil, "", err
	}

	if err := s.userService.SendVerificationEmail(ctx, createdUser); err != nil {
		return nil, "", err
	}

	oat, err := s.tokenService.GenerateOAT(ctx, keys.AuthToken, createdUser.ID.String(), keys.AuthTokenDuration)
	if err != nil {
		return nil, "", err
	}

	return createdUser, oat, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*dto.LoginResponse, string, error) {
	fetchedUser, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	err = security.ComparePassword(fetchedUser.Password, password)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	oat, err := s.tokenService.GenerateOAT(ctx, keys.AuthToken, fetchedUser.ID.String(), keys.AuthTokenDuration)
	if err != nil {
		return nil, "", err
	}

	restaurants, err := s.restaurantService.GetRestaurantsByUserID(ctx, fetchedUser.ID)
	if err != nil {
		if !errors.Is(err, ErrNoRestaurantFoundForUser) {
			return nil, "", err
		}
	}

	loginResponse := &dto.LoginResponse{
		User:        fetchedUser,
		Restaurants: restaurants,
	}

	return loginResponse, oat, nil
}

func (s *authService) Logout(ctx context.Context, oat string) error {
	return s.cache.Delete(ctx, cache.GenerateKey(string(keys.AuthToken), oat))
}
