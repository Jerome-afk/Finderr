package db

import (
	"fmt"
	"log"

	"finderr/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error

	// SQLite connection
	DB, err = gorm.Open(sqlite.Open("finderr.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// Auto migrate models
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Create a default user
	if err := createDafaultAdmin(); err != nil {
		return fmt.Errorf("failed to create default admin user: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func createDafaultAdmin() error {
	var count int64
	if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	// Only create admin if no users exist
	if count == 0 {
		admin := &models.User{
			Username:     "admin",
			Email:        "admin@finderr.com",
			Password:     "admin123", // In production, use environment variable for this
			Role:         models.AdminRole,
			ProfileImage: "/images/admin.webp",
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin.Password = string(hashedPassword)

		if err := DB.Create(admin).Error; err != nil {
			return err
		}
		log.Println("Default admin user created")
	}
	return nil
}
