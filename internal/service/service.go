package service

import (
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/cache"
	"github.com/memsbdm/restaurant-api/internal/mailer"
	"github.com/memsbdm/restaurant-api/internal/repository"
)

type Services struct {
	AuthService   AuthService
	GoogleService GoogleService
	MailerService MailerService
	TokenService  TokenService
	UserService   UserService
}

func New(cfg *config.Container, repos *repository.Repositories, cache cache.Cache, mailer mailer.Mailer) *Services {
	googleSvc := NewGoogleService(cfg.Google)
	tokenSvc := NewTokenService(cfg.Security, cache)
	mailerSvc := NewMailerService(cfg.Mailer, mailer)
	userSvc := NewUserService(cfg.App, repos.UserRepository, tokenSvc, mailerSvc)
	authSvc := NewAuthService(cfg.Security, cache, userSvc, tokenSvc)

	return &Services{
		AuthService:   authSvc,
		GoogleService: googleSvc,
		MailerService: mailerSvc,
		TokenService:  tokenSvc,
		UserService:   userSvc,
	}
}
