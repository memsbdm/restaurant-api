package service

import (
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/cache"
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/mailer"
)

type Services struct {
	AuthService       AuthService
	GoogleService     GoogleService
	MailerService     MailerService
	RestaurantService RestaurantService
	TokenService      TokenService
	UserService       UserService
}

func New(cfg *config.Container, db *database.DB, cache cache.Cache, mailer mailer.Mailer) *Services {
	googleSvc := NewGoogleService(cfg.Google)
	tokenSvc := NewTokenService(cfg.Security, cache)
	mailerSvc := NewMailerService(cfg.Mailer, mailer)
	userSvc := NewUserService(cfg.App, db, tokenSvc, mailerSvc)
	authSvc := NewAuthService(cfg.Security, cache, userSvc, tokenSvc)
	restaurantSvc := NewRestaurantService(db, googleSvc)

	return &Services{
		AuthService:       authSvc,
		GoogleService:     googleSvc,
		MailerService:     mailerSvc,
		RestaurantService: restaurantSvc,
		TokenService:      tokenSvc,
		UserService:       userSvc,
	}
}
