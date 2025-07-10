package routes

import (
	"finderr/handlers"
	"finderr/models"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Set up static file serving
	app.Static("/", "./public")

	// Home page
	app.Get("/", handlers.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		return c.Render("pages/home", fiber.Map{
			"Title": "Finderr - Home",
			"User":  user,
		}, "layout/default")
	})

	// Auth routes (without layout)
	auth := app.Group("/auth")
	auth.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("auth/login", fiber.Map{
			"Title": "Login",
		})
	})
	auth.Post("/login", handlers.Login)

	auth.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("auth/signup", fiber.Map{
			"Title": "Sign Up",
		})
	})
	auth.Post("/signup", handlers.Register)

	auth.Get("/logout", handlers.Logout)
}
