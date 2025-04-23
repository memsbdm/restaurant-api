package repository

import (
	"context"
	"log"

	"github.com/memsbdm/restaurant-api/internal/database/codegen"
	"github.com/memsbdm/restaurant-api/internal/dto"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (codegen.User, error)
	Create(ctx context.Context, userDto *dto.CreateUserDto) (codegen.User, error)
	EmailTaken(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	queries *codegen.Queries
}

func NewUserRepository(queries *codegen.Queries) *userRepository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (codegen.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("error: GetByEmail (user repository): %v", err)
		return codegen.User{}, err
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, userDto *dto.CreateUserDto) (codegen.User, error) {
	createdUser, err := r.queries.CreateUser(ctx, userDto.ToParams())
	if err != nil {
		log.Printf("error: Create (user repository): %v", err)
		return codegen.User{}, err
	}

	return createdUser, nil
}

func (r *userRepository) EmailTaken(ctx context.Context, email string) (bool, error) {
	taken, err := r.queries.UserEmailTaken(ctx, email)
	if err != nil {
		log.Printf("error: EmailTaken (user repository): %v", err)
		return false, err
	}

	return taken, nil
}
