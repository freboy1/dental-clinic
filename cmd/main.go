package main

import (
	"dental_clinic/internal/database"
	"dental_clinic/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg.DB_DSN)
	defer db.Close()

}
