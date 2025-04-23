package service

import (
	"context"
	"strings"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/cache"
	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/pkg/keys"
	"github.com/memsbdm/restaurant-api/pkg/security"
)

type AuthService interface {
	Register(ctx context.Context, user *dto.CreateUserDto) (dto.UserDTO, string, error)
	Login(ctx context.Context, email, password string) (dto.UserDTO, string, error)
	Logout(ctx context.Context, oat string) error
}

type authService struct {
	cfg          *config.Security
	cache        cache.Cache
	tokenService TokenService
	userService  UserService
}

func NewAuthService(cfg *config.Security, cache cache.Cache, userService UserService, tokenService TokenService) *authService {
	return &authService{
		cfg:          cfg,
		cache:        cache,
		tokenService: tokenService,
		userService:  userService,
	}
}

func (s *authService) Register(ctx context.Context, user *dto.CreateUserDto) (dto.UserDTO, string, error) {
	createdUser, err := s.userService.Create(ctx, user)
	if err != nil {
		return dto.UserDTO{}, "", err
	}

	if err := s.userService.SendVerificationEmail(ctx, createdUser); err != nil {
		return dto.UserDTO{}, "", err
	}

	oat, err := s.tokenService.GenerateOAT(ctx, keys.AuthToken, createdUser.ID.String(), keys.AuthTokenDuration)
	if err != nil {
		return dto.UserDTO{}, "", err
	}

	return createdUser, oat, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (dto.UserDTO, string, error) {
	fetchedUser, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return dto.UserDTO{}, "", response.ErrInvalidCredentials
	}

	err = security.ComparePassword(fetchedUser.Password, password)
	if err != nil {
		return dto.UserDTO{}, "", response.ErrInvalidCredentials
	}

	oat, err := s.tokenService.GenerateOAT(ctx, keys.AuthToken, fetchedUser.ID.String(), keys.AuthTokenDuration)
	if err != nil {
		return dto.UserDTO{}, "", err
	}

	return fetchedUser, oat, nil
}

func (s *authService) Logout(ctx context.Context, oat string) error {
	parts := strings.Split(oat, ".")
	if len(parts) != 2 {
		return response.ErrUnauthorized
	}

	return s.cache.Delete(ctx, cache.GenerateKey(string(keys.AuthToken), parts[0]))
}
