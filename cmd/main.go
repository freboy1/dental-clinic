package main

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/database"
	"dental_clinic/internal/handlers"
	"dental_clinic/internal/repository"
	"dental_clinic/internal/services"
	"net/http"

	"log"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg.DB_DSN)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, *cfg)
	userHandler := handlers.NewUserHandler(userService)

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	{
		api.HandleFunc("/register", userHandler.Register).Methods("POST")
		api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
		api.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
		api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
		api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
		api.HandleFunc("/verify", userHandler.VerifyAccountByLink).Methods("GET")
	}

	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))

}
