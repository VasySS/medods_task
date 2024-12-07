package main

import (
	appHTTP "auth_service/internal/app/http"
	"auth_service/internal/config"
	"auth_service/internal/repository/postgres"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Routes for authentication services and Swagger documentation.
// @title           Auth Service API
// @version         0.1
// @host            localhost:8080
// @BasePath        /v1/auth
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pg, _ := postgres.New(context.Background(), config.PostgresDSN)
	defer pg.Close()

	r := appHTTP.NewRouter(pg)
	startHTTP(ctx, r)
}

func startHTTP(ctx context.Context, r http.Handler) {
	srv := &http.Server{
		Addr:         config.BackendAddr,
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
	}

	go func() {
		log.Printf("запуск сервера на %s", config.BackendAddr)

		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ошибка сервера: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("остановка сервера...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("ошибка при остановке сервера: %v", err)
	}
}
