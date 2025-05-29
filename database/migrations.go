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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS movies (
			movie_id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			state TEXT NOT NULL,
			url TEXT NOT NULL,
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS shows (
			season_id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			season TEXT NOT NULL,
			state TEXT NOT NULL,
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tv_episodes (
			episode_id INTEGER PRIMARY KEY AUTOINCREMENT,
			episode_title VARCHAR(255) NOT NULL,
			state TEXT NOT NULL,
			url TEXT NOT NULL,
			season_id INTEGER NOT NULL,
			FOREIGN KEY (season_id) REFERENCES shows(season_id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS anime (
			anime_id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			season TEXT NOT NULL,
			state TEXT NOT NULL,
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS anime_episodes (
			episode_id INTEGER PRIMARY KEY AUTOINCREMENT,
			episode_title VARCHAR(255) NOT NULL,
			state TEXT NOT NULL,
			url TEXT NOT NULL,
			season_id INTEGER NOT NULL,
			FOREIGN KEY (season_id) REFERENCES anime(season_id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS anime_movies (
			animeM_id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			state TEXT NOT NULL,
			url TEXT NOT NULL,
		)
	`)
	if err != nil {
		return err
	}
	
	log.Println("Database migrations completed")
	return nil
}