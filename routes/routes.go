package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Set up static file serving
	app.Static("/", "./public")

	// Home page
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("pages/home", fiber.Map{
			"Title": "Finderr - Home",
		}, "layout/default")
	})

	// Auth routes (without layout)
	auth := app.Group("/auth")
	auth.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("auth/login", fiber.Map{
			"Title": "Login",
		})
	})
	auth.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("auth/signup", fiber.Map{
			"Title": "Sign Up",
		})
	})
}