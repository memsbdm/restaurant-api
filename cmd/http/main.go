package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/memsbdm/restaurant-api/internal/app"
)

func main() {
	app := app.New()
	defer app.Cleanup()

	done := make(chan bool)
	go gracefulShutdown(app.Server.Server, done)
	app.Server.ListenAndServe()

	<-done
	log.Println("Graceful shutdown complete")
}

func gracefulShutdown(srv *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Server exiting")

	done <- true
}
