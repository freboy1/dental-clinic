package main

import (
	"net/http"
	"log"

	"dental_clinic/internal/config"
	"dental_clinic/internal/database"
	"dental_clinic/internal/router"
)

// @title Dental Clinic API
// @version 1.0
// @description API for managing users, authentication, and clinic operations.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /api
func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg.DB_DSN)
	defer db.Close()

	r := router.NewRouter(cfg, db)

	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
