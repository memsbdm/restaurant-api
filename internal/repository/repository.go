package repository

import (
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/database/codegen"
)

type Repositories struct {
	UserRepository UserRepository
}

func New(db *database.DB) *Repositories {
	queries := codegen.New(db)

	return &Repositories{
		UserRepository: NewUserRepository(queries),
	}
}
