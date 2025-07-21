package handlers

import (
	"finderr/services"

	"github.com/gofiber/fiber/v2"
)

func ShowsHandler(c *fiber.Ctx) error {
	tmdb := services.NewTMDBClient()

	// Movies API calls
	popularTvShows, _ := tmdb.GetPopularTVShows()
	trendingTvShows, _ := tmdb.GetTrendingTVShows(15)

	return c.Render("pages/movies", fiber.Map{
		"Title":          "TV Shows",
		"TrendingTVShows": trendingTvShows,
		"PopularMovies":  popularTvShows,
		"User":           c.Locals("user"),
	}, "layout/default")
}