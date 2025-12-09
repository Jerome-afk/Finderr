package main

import (
	"finderr/database"
	"finderr/handlers"
	"finderr/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

// Load environment files
func prev() error {
	err := godotenv.Load()
	return err
}

func main() {
	err := prev()
	if err != nil {
		log.Fatalf("Error loading the env file: %v", err)
	}

	err = handlers.LogHandler()
	if err != nil {
		log.Fatalf("Log init failed: %v", err)
	}

	handlers.WriteLog("INFO", "Server started", map[string]interface{}{"port": 8080})



	// Initialize the database
	if err := database.InitDB(); err != nil {
		handlers.WriteLog("ERROR", "Database connection failed", nil)
		log.Fatalf("failed to initialize database: %v", err)
	}
	handlers.WriteLog("PASS", "Database connection success", nil)


	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	routes.SetRoutes(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))

}
