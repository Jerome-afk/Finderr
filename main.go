package main

import (
	"log"
	"os"

	"finderr/db"
	"finderr/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	engine := html.New("./views", ".html")
	engine.AddFunc("formatYear", func(date string) string {
		if len(date) >= 4 {
			return date[:4]
		}
		return ""
	})

	engine.AddFunc("formatType", func(mediaType string) string {
		switch mediaType {
		case "movie":
			return "Movie"
		case "tv":
			return "TV Show"
		case "anime":
			return "Anime"
		default:
			return "Unknown"
		}
	})
	
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SECRET"),
	}))

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}