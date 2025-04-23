package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/cache"
	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/pkg/security"
)

const (
	OATCachePrefix   = "oat"
	OATCacheDuration = 7 * 24 * time.Hour
)

type AuthService interface {
	Register(ctx context.Context, user *dto.CreateUserDto) (dto.UserDTO, string, error)
	Login(ctx context.Context, email, password string) (dto.UserDTO, string, error)
	GetUserIDFromSignedOAT(ctx context.Context, signedOAT string) (string, error)
	Logout(ctx context.Context, signedOAT string) error
}

type authService struct {
	cfg         *config.Security
	cache       cache.Cache
	userService UserService
}

func NewAuthService(cfg *config.Security, cache cache.Cache, userService UserService) *authService {
	return &authService{
		cfg:         cfg,
		cache:       cache,
		userService: userService,
	}
}

func (s *authService) Register(ctx context.Context, user *dto.CreateUserDto) (dto.UserDTO, string, error) {
	createdUser, err := s.userService.Create(ctx, user)
	if err != nil {
		return dto.UserDTO{}, "", err
	}

	signedOAT, err := s.generateSignedOAT(ctx, createdUser.ID)
	if err != nil {
		return dto.UserDTO{}, "", err
	}

	return createdUser, signedOAT, nil
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

	signedOAT, err := s.generateSignedOAT(ctx, fetchedUser.ID)
	if err != nil {
		return dto.UserDTO{}, "", err
	}

	return fetchedUser, signedOAT, nil
}

func (s *authService) Logout(ctx context.Context, signedOAT string) error {
	parts := strings.Split(signedOAT, ".")
	if len(parts) != 2 {
		return response.ErrUnauthorized
	}

	return s.cache.Delete(ctx, cache.GenerateKey(OATCachePrefix, parts[0]))
}

func (s *authService) GetUserIDFromSignedOAT(ctx context.Context, signedOAT string) (string, error) {
	parts := strings.Split(signedOAT, ".")
	if len(parts) != 2 {
		return "", response.ErrUnauthorized
	}

	oat, signature := parts[0], parts[1]
	hasValidSignature := security.VerifySignature(oat, signature, s.cfg.OATSignature)
	if !hasValidSignature {
		return "", response.ErrUnauthorized
	}

	userID, err := s.cache.Get(ctx, cache.GenerateKey(OATCachePrefix, oat))
	if err != nil {
		return "", err
	}

	return string(userID), nil
}

func (s *authService) generateSignedOAT(ctx context.Context, userID uuid.UUID) (string, error) {
	oat, err := security.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	err = s.cache.Set(ctx, cache.GenerateKey(OATCachePrefix, oat), []byte(userID.String()), OATCacheDuration)
	if err != nil {
		return "", err
	}

	signature := security.SignString(oat, s.cfg.OATSignature)
	signedOAT := string(oat) + "." + signature

	return signedOAT, nil
}
