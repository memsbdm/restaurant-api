package service

import (
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/cache"
	"github.com/memsbdm/restaurant-api/internal/repository"
)

type Services struct {
	UserService UserService
	AuthService AuthService
}

func New(cfg *config.Container, repos *repository.Repositories, cache cache.Cache) *Services {
	userSvc := NewUserService(repos.UserRepository)
	authSvc := NewAuthService(cfg.Security, cache, userSvc)

	return &Services{
		UserService: userSvc,
		AuthService: authSvc,
	}
}
