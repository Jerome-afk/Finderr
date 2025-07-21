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
	// app.Get("/", handlers.AuthRequired, func(c *fiber.Ctx) error {
	// 	user := c.Locals("user").(*models.User)
	// 	return c.Render("pages/home", fiber.Map{
	// 		"Title": "Finderr - Home",
	// 		"User":  user,
	// 	}, "layout/default")
	// })

	app.Get("/", handlers.AuthRequired, handlers.HomeHandler)

	// Normal pages
	norm := app.Group("/pages")
	norm.Get("/anime", handlers.AuthRequired, handlers.AnimeHandler)

	norm.Get("/movies", handlers.AuthRequired, handlers.MoviesHandler)

	norm.Get("/tv-shows", handlers.AuthRequired, handlers.ShowsHandler)

	norm.Get("/downloads", handlers.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		return c.Render("pages/downloads", fiber.Map{
			"Title": "Downloads",
			"User": user,
		}, "layout/default")
	})

	norm.Get("/dashboard", handlers.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		return c.Render("pages/dashboard", fiber.Map{
			"Title": "User Dashboard",
			"User": user,
		}, "layout/default")
	})

	norm.Get("/music", handlers.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		return c.Render("pages/music", fiber.Map{
			"Title": "Finderr - Music",
			"User": user,
		}, "layout/default")
	})

	norm.Get("/discover", handlers.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		return c.Render("pages/discover", fiber.Map{
			"Title": "Finderr - Discover",
			"User": user,
		})
	})

	// Config pages
	config := app.Group("/config")
	config.Get("/users", handlers.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		return c.Render("changes/users", fiber.Map{
			"Title": "Config - Users",
			"User":  user,
		}, "layout/settings")
	})

	config.Get("/logs", handlers.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		return c.Render("changes/logs", fiber.Map{
			"Title": "Config - Logs",
			"User":  user,
		}, "layout/settings")
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
