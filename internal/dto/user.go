package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/database/codegen"
)

type UserDTO struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Password        string    `json:"-"`
	IsEmailVerified bool      `json:"is_email_verified"`
	AvatarURL       *string   `json:"avatar_url"`
}

func NewUserDTO(user codegen.User) UserDTO {
	return UserDTO{
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

type CreateUserDto struct {
	Name     string
	Email    string
	Password string
}

func (u CreateUserDto) ToParams() codegen.CreateUserParams {
	return codegen.CreateUserParams{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
