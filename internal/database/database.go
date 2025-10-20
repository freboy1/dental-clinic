package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB(dsn string) *sql.DB {
	if dsn == "" {
		return nil
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil
	}

	if err := db.Ping(); err != nil {
		return nil
	}

	fmt.Println("Database connected")
	return db
}
