package app

import (
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/cache"
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/handler"
	"github.com/memsbdm/restaurant-api/internal/middleware"
	"github.com/memsbdm/restaurant-api/internal/repository"
	"github.com/memsbdm/restaurant-api/internal/server"
	"github.com/memsbdm/restaurant-api/internal/service"
)

type App struct {
	DB     *database.DB
	Cache  cache.Cache
	Server *server.Server
}

func New() *App {
	cfg := config.New()
	db := database.NewPostgres(cfg.DB)
	cache := cache.NewRedis(cfg.Cache)

	repos := repository.New(db)
	services := service.New(cfg, repos, cache)
	middle := middleware.New(cfg, services)
	handlers := handler.New(cfg, services)

	server := server.New(cfg, handlers, middle)

	return &App{
		Cache:  cache,
		DB:     db,
		Server: server,
	}
}

func (a *App) Cleanup() {
	a.DB.Close()
	a.Cache.Close()
}
