package routes

import (
	"database/sql"
	"log"
	"os"

	"finderr/database"
	"finderr/logger"
	"finderr/torrentHandler/websocket"

	"github.com/gin-gonic/gin"
)

// Initialize server
func InitServer() *gin.Engine {
	router := gin.Default()

	// Initialize the logger
	logger.InitLogger()

	// Create directories
	dirs := []string{"./movies", "./tv_shows", "./anime", "./anime_m", "./novels", "./music"}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Initialize database
	defaultRoute := "./forum.db"

	// Ask for database configuration
	dbConfig, err := database.AskForDatabaseConfig(defaultRoute)
	if err != nil {
		log.Fatalf("Error configuring database: %v\n", err)
		os.Exit(1)
	}

	// Save the configuration for future use
	err = database.SaveConfig(dbConfig)
	if err != nil {
		log.Printf("Warning: Failed to save database configuration: %v\n", err)
	}

	// Initialize database with the configuration
	var db *sql.DB
	if dbConfig.IsLocal && dbConfig.DBType == "sqlite" {
		// For SQLite, use the file path
		db, err = database.InitDB(dbConfig.FilePath)
	} else {
		// For cloud databases, use the configuration
		db, err = database.InitDBWithConfig(dbConfig)
	}

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	err = database.RunMigrations(db)
	if err != nil {
		logger.LogEvent("Server", "Failed to run migration")
		log.Fatalf("Failed to run migrations: %v", err)
	}
	logger.LogEvent("Server", "Database set up correctly and running")

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// TODO: Add function to create a Admin and Password

	// Initialize routes
	InitializeRoutes(router, hub)

	// Log server start
	logger.LogEvent("server", "Server started successfully")

	return router
}
