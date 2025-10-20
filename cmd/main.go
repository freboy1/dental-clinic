package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/freboy1/dental-clinic/internal/config"
	"github.com/freboy1/dental-clinic/internal/database"
	"github.com/freboy1/dental-clinic/internal/handlers"
	"github.com/freboy1/dental-clinic/internal/repository"
	"github.com/freboy1/dental-clinic/internal/server"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	ctx := context.Background()

	dbpool, err := database.NewPgPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer dbpool.Close(ctx)

	userRepo := repository.NewUserRepository(dbpool)
	loginRepo := repository.NewLoginHistoryRepository(dbpool)

	h := handlers.NewHandler(cfg, userRepo, loginRepo)

	router := server.NewRouter(h, cfg.JWTSecret)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      loggingMiddleware(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server started on port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
