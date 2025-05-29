package routes

import (
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

	// TODO: Add function to create a Admin and Password

	// Initialize database
	defaultRoute := "./forum.db"

	route, isDefault, err := database.AskForRoute(defaultRoute)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		os.Exit(1)
	}

	if isDefault {
		// Use default route
	} else {
		defaultRoute = route
	}
	db, err := database.InitDB(defaultRoute)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	err = database.RunMigrations(db)
	if err != nil {
		logger.LogEvent("Server", "Failed to run migration")
		log.Fatalf("Failed to run migrations: %v", err)
	}
	logger.LogEvent("Server", "Database set up correctly and running")

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()


	// Initialize routes
	InitializeRoutes(router, hub)

	// Log server start
	logger.LogEvent("server", "Server started successfully")

	return router
}
