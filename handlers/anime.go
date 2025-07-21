package handlers

import (
	"finderr/services"

	"github.com/gofiber/fiber/v2"
)

func AnimeHandler(c *fiber.Ctx) error {
	anilist := services.NewAniListClient()

	// Movies API calls
	popularAnime, _ := anilist.GetPopularAnime()
	trendingAnime, _ := anilist.GetTrendingAnime(15)

	return c.Render("pages/movies", fiber.Map{
		"Title":          "Anime",
		"TrendingAnime":  trendingAnime,
		"PopularMovies":  popularAnime,
		"User":           c.Locals("user"),
	}, "layout/default")
}