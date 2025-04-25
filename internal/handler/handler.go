package handler

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/service"
)

type Handlers struct {
	AuthHandler        *AuthHandler
	GoogleHandler      *GoogleHandler
	RestaurantHandler  *RestaurantHandler
	VerifyEmailHandler *VerifyEmailHandler
}

func New(cfg *config.Container, services *service.Services) *Handlers {
	return &Handlers{
		AuthHandler:        NewAuthHandler(cfg.App, services.AuthService),
		GoogleHandler:      NewGoogleHandler(services.GoogleService),
		RestaurantHandler:  NewRestaurantHandler(services.RestaurantService),
		VerifyEmailHandler: NewVerifyEmailHandler(services.UserService),
	}
}

func IsMobileRequest(r *http.Request) bool {
	return r.Header.Get("Client-Type") == "mobile"
}
