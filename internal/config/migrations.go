package config

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cnf Config) {

	m, err := migrate.New("file://migrations", cnf.DB_DSN)

	if err != nil {
		log.Fatalf("failed to initialize migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}