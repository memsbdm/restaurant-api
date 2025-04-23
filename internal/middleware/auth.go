package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/memsbdm/restaurant-api/internal/handler"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
	"github.com/memsbdm/restaurant-api/pkg/contextkeys"
)

func AuthMiddleware(authSvc service.AuthService) Middle {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signedOAT, err := extractSignedOATFromRequest(r)
			if err != nil {
				response.HandleError(w, err)
				return
			}

			userID, err := authSvc.GetUserIDFromSignedOAT(r.Context(), signedOAT)
			if err != nil {
				response.HandleError(w, response.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.UserIDContextKey, userID)
			ctx = context.WithValue(ctx, contextkeys.SignedOATContextKey, signedOAT)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GuestMiddleware(authSvc service.AuthService) Middle {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signedOAT, err := extractSignedOATFromRequest(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			_, err = authSvc.GetUserIDFromSignedOAT(r.Context(), signedOAT)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			response.HandleError(w, response.ErrForbidden)
		})
	}
}

func extractSignedOATFromRequest(r *http.Request) (string, error) {
	var signedOAT string
	var err error
	if handler.IsMobileRequest(r) {
		signedOAT, err = getSignedOATFromHeader(r)
	} else {
		signedOAT, err = getSignedOATFromCookie(r)
	}

	if err != nil {
		return "", err
	}
	return signedOAT, nil
}

func getSignedOATFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("go-session")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func getSignedOATFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
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
