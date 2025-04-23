package middleware

import (
	"net/http"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/service"
)

type Middleware struct {
	Auth    MiddlewareFunc
	Guest   MiddlewareFunc
	Logging Middle
}

type MiddlewareFunc func(handler func(http.ResponseWriter, *http.Request)) http.Handler

func New(cfg *config.Container, s *service.Services) *Middleware {
	return &Middleware{
		Auth:    newHandlerMiddleware(AuthMiddleware(s.AuthService)),
		Guest:   newHandlerMiddleware(GuestMiddleware(s.AuthService)),
		Logging: LoggingMiddleware,
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
