package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/memsbdm/restaurant-api/internal/handler"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
	"github.com/memsbdm/restaurant-api/pkg/keys"
	"github.com/memsbdm/restaurant-api/pkg/security"
)

func AuthMiddleware(appEnv string, tokenSvc service.TokenService, authSvc service.AuthService) Middle {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			oat, err := extractAuthOATFromRequest(r)
			if err != nil {
				response.HandleError(w, err)
				return
			}

			userID, err := tokenSvc.VerifyOAT(r.Context(), keys.AuthToken, oat)
			if err != nil {
				response.HandleError(w, response.ErrUnauthorized)
				return
			}

			// Reset cookie expiration time
			handler.SetAuthCookie(w, oat, appEnv)
			// Enrich context with user ID and OAT
			ctx := context.WithValue(r.Context(), keys.UserIDContextKey, userID)
			decodedOAT, _ := security.DecodeTokenURLSafe(oat)
			oat = strings.Split(decodedOAT, ".")[0]
			ctx = context.WithValue(ctx, keys.AuthOATContextKey, oat)

			// Reset the OAT cache TTL
			err = authSvc.ResetAuthOATCacheTTL(r.Context(), oat, userID)
			if err != nil {
				response.HandleError(w, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GuestMiddleware(tokenSvc service.TokenService) Middle {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			oat, err := extractAuthOATFromRequest(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			_, err = tokenSvc.VerifyOAT(r.Context(), keys.AuthToken, oat)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			response.HandleError(w, response.ErrForbidden)
		})
	}
}

func extractAuthOATFromRequest(r *http.Request) (string, error) {
	var oat string
	var err error
	if handler.IsMobileRequest(r) {
		oat, err = getAuthOATFromHeader(r)
	} else {
		oat, err = getAuthOATFromCookie(r)
	}

	if err != nil {
		return "", err
	}
	return oat, nil
}

func getAuthOATFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(keys.AuthOATCookieName)
	if err != nil {
		return "", response.ErrUnauthorized
	}

	return cookie.Value, nil
}

func getAuthOATFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get(keys.AuthorizationHeaderName)
	if authHeader == "" {
		return "", response.ErrUnauthorized
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", response.ErrUnauthorized
	}

	token := strings.TrimPrefix(authHeader, prefix)

	return token, nil
}
