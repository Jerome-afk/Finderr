package handlers

import (
	"finderr/services"

	"github.com/gofiber/fiber/v2"
)

func MoviesHandler(c *fiber.Ctx) error {
	tmdb := services.NewTMDBClient()

	// Movies API calls
	popularMovies, _ := tmdb.GetPopularMovies()
	trendingMovies, _ := tmdb.GetTrendingMovies(15)

	return c.Render("pages/movies", fiber.Map{
		"Title":          "Movies",
		"TrendingMovies": trendingMovies,
		"PopularMovies":  popularMovies,
		"User":           c.Locals("user"),
	}, "layout/default")
}