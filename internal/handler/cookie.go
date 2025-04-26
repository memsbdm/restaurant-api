package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/config"
)

func SetActiveRestaurantCookie(w http.ResponseWriter, restaurantID uuid.UUID, appEnv string) {
	cookie := &http.Cookie{
		Name:     "active_restaurant",
		Value:    restaurantID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == config.EnvProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7,
	}
	http.SetCookie(w, cookie)
}

func setAuthCookie(w http.ResponseWriter, oat, appEnv string) {
	cookie := &http.Cookie{
		Name:     "go-session",
		Value:    oat,
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == config.EnvProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7,
	}

	http.SetCookie(w, cookie)
}

func clearAuthCookie(w http.ResponseWriter, appEnv string) {
	cookie := &http.Cookie{
		Name:     "go-session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == config.EnvProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}
