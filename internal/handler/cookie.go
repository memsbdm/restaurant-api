package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/pkg/keys"
)

func SetActiveRestaurantCookie(w http.ResponseWriter, restaurantID uuid.UUID, appEnv string) {
	cookie := &http.Cookie{
		Name:     keys.ActiveRestaurantCookieName,
		Value:    restaurantID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == config.EnvProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24, // 1 day
	}
	http.SetCookie(w, cookie)
}

func SetAuthCookie(w http.ResponseWriter, oat, appEnv string) {
	cookie := &http.Cookie{
		Name:     keys.AuthOATCookieName,
		Value:    oat,
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == config.EnvProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60, // 1 hour
	}

	http.SetCookie(w, cookie)
}

func clearAuthCookie(w http.ResponseWriter, appEnv string) {
	cookie := &http.Cookie{
		Name:     keys.AuthOATCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == config.EnvProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}
