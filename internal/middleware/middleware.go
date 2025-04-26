package middleware

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/service"
)

type Middleware struct {
	Auth       MiddlewareFunc
	Guest      MiddlewareFunc
	Logging    Middle
	Restaurant MiddlewareFunc
}

type MiddlewareFunc func(handler func(http.ResponseWriter, *http.Request)) http.Handler

func New(cfg *config.Container, s *service.Services) *Middleware {
	return &Middleware{
		Auth:       newHandlerMiddleware(AuthMiddleware(s.TokenService)),
		Guest:      newHandlerMiddleware(GuestMiddleware(s.TokenService)),
		Logging:    LoggingMiddleware,
		Restaurant: newHandlerMiddleware(RestaurantMiddleware(cfg.App.Env, s.RestaurantService, s.RestaurantUserService)),
	}
}

func newHandlerMiddleware(m Middle) MiddlewareFunc {
	return func(handler func(http.ResponseWriter, *http.Request)) http.Handler {
		return m(http.HandlerFunc(handler))
	}
}

type Middle func(http.Handler) http.Handler

func CreateStack(xs ...Middle) Middle {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

type HandlerMiddleware func(handler http.Handler) http.Handler

func Chain(f http.HandlerFunc, middleware ...MiddlewareFunc) http.HandlerFunc {
	for _, m := range middleware {
		h := m(f)
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	return f
}

func ChainHandlerFunc(h http.Handler, middleware ...HandlerMiddleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}
