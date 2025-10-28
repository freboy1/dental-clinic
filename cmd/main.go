package main

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/database"
	"dental_clinic/internal/handlers"
	"dental_clinic/internal/middleware"
	"dental_clinic/internal/repository"
	"dental_clinic/internal/services"
	"net/http"

	"log"

	_ "dental_clinic/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Dental Clinic API
// @version 1.0
// @description API for managing users, authentication, and clinic operations.
// @host localhost:8080
// @BasePath /api
func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg.DB_DSN)
	defer db.Close()

	config.RunMigrations(*cfg)

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, *cfg)
	userHandler := handlers.NewUserHandler(userService, *cfg)

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := router.PathPrefix("/api").Subrouter()

	public := api.NewRoute().Subrouter()
	{
		public.HandleFunc("/register", userHandler.Register).Methods("POST")
		public.HandleFunc("/login", userHandler.Login).Methods("POST")
		public.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
		// public.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
		public.HandleFunc("/verify", userHandler.VerifyAccountByLink).Methods("GET")
	}
	private := api.NewRoute().Subrouter()
	private.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		private.HandleFunc("/users/verify-email", userHandler.VerifyNewEmail).Methods("GET")
		private.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
		private.HandleFunc("/users/update-password", userHandler.UpdatePassword).Methods("POST")
		private.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
		private.HandleFunc("/users/update-email", userHandler.UpdateEmail).Methods("POST")
		private.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	}

	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))

}
