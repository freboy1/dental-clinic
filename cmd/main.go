package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"dental_clinic/internal/config"
	"dental_clinic/internal/database"
	"dental_clinic/internal/jobs"
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

	jobs.StartAppointmentStatusCron(context.Background(), db, time.Minute)

	r := router.NewRouter(cfg, db)

	log.Printf("Server running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Println(err)
	}
}
