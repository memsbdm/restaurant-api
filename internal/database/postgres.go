package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/database/repository"
)

type DB struct {
	*pgx.Conn
	Queries *repository.Queries
}

func NewPostgres(cfg *config.DB) *DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Schema)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("Connected to database")
	queries := repository.New(conn)

	return &DB{conn, queries}
}

func (db *DB) Close() {
	if db.Conn == nil {
		return
	}

	if err := db.Conn.Close(context.Background()); err != nil {
		log.Printf("failed to close database connection: %v", err)
	}

	log.Println("Database connection closed")
}
