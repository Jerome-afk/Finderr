package server

import "github.com/gin-gonic/gin"


//Initialize server
func InitServer() *gin.Engine {
	router := gin.Default()

	// Initialize the logger
	initLogger()

	// check user directory
	checkAndCreateAdminFile()

	//Initialize routes
	InitializeRoutes(router)

	//Log server start
	logEvent("server", "Server started successfully")

	return router
}