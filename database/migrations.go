package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	// Users table
	_, err := db.Exec(`
	    CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Database migrations completed")
	return nil
}