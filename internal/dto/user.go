package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Password        string    `json:"-"`
	IsEmailVerified bool      `json:"is_email_verified"`
	AvatarURL       *string   `json:"avatar_url"`
}

func (u *User) ToUpdateParams() repository.UpdateUserParams {
	return repository.UpdateUserParams{
		ID:              u.ID,
		Email:           u.Email,
		IsEmailVerified: u.IsEmailVerified,
		AvatarUrl:       u.AvatarURL,
	}
}

func NewUser(user repository.User) User {
	return User{
		ID:              user.ID,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		IsEmailVerified: user.IsEmailVerified,
		AvatarURL:       user.AvatarUrl,
	}
}

type CreateUser struct {
	Name     string
	Email    string
	Password string
}

func (u CreateUser) ToParams() repository.CreateUserParams {
	return repository.CreateUserParams{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
