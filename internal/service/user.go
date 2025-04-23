package service

import (
	"context"

	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/internal/repository"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/pkg/security"
)

type UserService interface {
	GetByEmail(ctx context.Context, email string) (dto.UserDTO, error)
	Create(ctx context.Context, user *dto.CreateUserDto) (dto.UserDTO, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByEmail(ctx context.Context, email string) (dto.UserDTO, error) {
	dbUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.NewUserDTO(dbUser), nil
}

func (s *userService) Create(ctx context.Context, user *dto.CreateUserDto) (dto.UserDTO, error) {
	emailTaken, err := s.userRepo.EmailTaken(ctx, user.Email)
	if err != nil {
		return dto.UserDTO{}, err
	}

	if emailTaken {
		return dto.UserDTO{}, response.ErrEmailConflict
	}

	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return dto.UserDTO{}, err
	}

	user.Password = hashedPassword

	dbUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.NewUserDTO(dbUser), nil
}
