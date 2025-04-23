package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/handler"
	"github.com/memsbdm/restaurant-api/internal/middleware"
)

type Server struct {
	*http.Server
}

func New(cfg *config.Container, h *handler.Handlers, m *middleware.Middleware) *Server {
	router := registerRoutes(h, m)

	stack := middleware.CreateStack(
		m.Logging,
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{Server: srv}
}

func (s *Server) ListenAndServe() {
	err := s.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown with error: %v", err)
	}
}
