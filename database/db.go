package database

import (
	"finderr/models"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() error {
	var err error

	// Using SQLite for development
	DB, err := gorm.Open(sqlite.Open("finderr.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// PostgressDB


	// Auto migrate
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialisation successful.")
	return nil
}