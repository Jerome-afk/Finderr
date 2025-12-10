package routes

import "github.com/gofiber/fiber/v2"

func SetRoutes(app *fiber.App) {
	// Serve the static files
	app.Static("/static", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("pages/index", fiber.Map{
			"Title": "Home",
		}, "layout/default")
	})

	// Authentication pages
	auth := app.Group("/auth")
	auth.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("auth/login", fiber.Map{
			"Title": "Login",
		}, "layout/auth")
	})

	// Normal pages
	// norm := app.Group("/pages")

	// Configuration pages
	// conf := app.Group("/config")
}