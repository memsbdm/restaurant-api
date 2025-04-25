package server

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/internal/handler"
	"github.com/memsbdm/restaurant-api/internal/middleware"
)

func registerRoutes(h *handler.Handlers, m *middleware.Middleware) http.Handler {
	r := http.NewServeMux()

	// Auth
	r.Handle("POST /auth/register", m.Guest(h.AuthHandler.Register))
	r.Handle("POST /auth/login", m.Guest(h.AuthHandler.Login))
	r.Handle("DELETE /auth/logout", m.Auth(h.AuthHandler.Logout))

	// Users
	r.HandleFunc("GET /users/verify-email", h.VerifyEmailHandler.VerifyEmail)
	r.Handle("POST /users/verify-email/resend", m.Auth(h.VerifyEmailHandler.ResendVerificationEmail))

	// Restaurants
	r.Handle("POST /restaurants", m.Auth(h.RestaurantHandler.Create))

	// Google
	r.Handle("GET /google/autocomplete", m.Auth(h.GoogleHandler.Autocomplete))

	// Sub-routes
	apiV1 := http.NewServeMux()
	apiV1.Handle("/api/v1/", http.StripPrefix("/api/v1", r))

	return apiV1
}
