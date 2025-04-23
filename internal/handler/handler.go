package handler

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/service"
)

type Handlers struct {
	AuthHandler *AuthHandler
}

func New(cfg *config.Container, services *service.Services) *Handlers {
	return &Handlers{
		AuthHandler: NewAuthHandler(cfg.App, services.AuthService),
	}
}

func IsMobileRequest(r *http.Request) bool {
	return r.Header.Get("Client-Type") == "mobile"
}
