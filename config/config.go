package config

import (
	"log"
	"time"

	"github.com/memsbdm/restaurant-api/pkg/env"
)

const (
	EnvProduction = "production"
	EnvStaging    = "staging"
	EnvDev        = "dev"
)

type (
	Container struct {
		App      *App
		Cache    *Cache
		DB       *DB
		Security *Security
		Server   *Server
	}

	App struct {
		Env string
	}

	Cache struct {
		Addr     string
		Password string
	}

	DB struct {
		Host        string
		Port        int
		User        string
		Password    string
		Database    string
		Schema      string
		MaxIdleConn int
		MaxOpenConn int
		MaxIdleTime time.Duration
	}

	Security struct {
		OATSignature []byte
	}

	Server struct {
		Port int
	}
)

func New() *Container {
	app := &App{
		Env: env.GetString("ENVIRONMENT"),
	}

	cache := &Cache{
		Addr:     env.GetString("CACHE_ADDR"),
		Password: env.GetOptionalString("CACHE_PASSWORD", ""),
	}

	db := &DB{
		Host:        env.GetString("DB_HOST"),
		Port:        env.GetInt("DB_PORT"),
		User:        env.GetString("DB_USER"),
		Password:    env.GetOptionalString("DB_PASSWORD", ""),
		Database:    env.GetString("DB_DATABASE"),
		Schema:      env.GetOptionalString("DB_SCHEMA", "public"),
		MaxIdleConn: env.GetOptionalInt("DB_MAX_IDLE_CONN", 30),
		MaxOpenConn: env.GetOptionalInt("DB_MAX_OPEN_CONN", 30),
		MaxIdleTime: env.GetOptionalDuration("MAX_IDLE_TIME", 15*time.Minute),
	}

	security := &Security{
		OATSignature: env.GetBytes("OAT_SIGNATURE"),
	}

	server := &Server{
		Port: env.GetOptionalInt("PORT", 8080),
	}

	c := &Container{
		App:      app,
		Cache:    cache,
		DB:       db,
		Security: security,
		Server:   server,
	}

	return c.Validate()
}

func (c *Container) Validate() *Container {
	if c.App.Env != EnvProduction && c.App.Env != EnvStaging && c.App.Env != EnvDev {
		log.Fatalf("env variable ENVIRONMENT is incorrect, got: %s", c.App.Env)
	}

	return c
}
