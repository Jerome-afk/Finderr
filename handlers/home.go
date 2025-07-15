package handlers

import (
	"finderr/models"
	"finderr/services"

	"github.com/gofiber/fiber/v2"
)

func HomeHandler(c *fiber.Ctx) error {
	tmdb := services.NewTMDBClient()
	anilist := services.NewAniListClient()

	// Get trending content
	trendingMovies, _ := tmdb.GetTrendingMovies()
	trendingTV, _ := tmdb.GetTrendingTVShows()
	trendingAnime, _ := anilist.GetTrendingAnime()

	// Get Popular content
	popularMovies, _ := tmdb.GetPopularMovies()
	popularTV, _ := tmdb.GetPopularTVShows()
	popularAnime, _ := anilist.GetPopularAnime()

	return c.Render("pages/home", fiber.Map{
		"Title": "Finderr - Home",
		"TrendingSections": []models.MediaSection{
			{
				Title: "Trending Movies",
				Items: trendingMovies,
			},
			{
				Title: "Trending TV Shows",
				Items: trendingTV,
			},
			{
				Title: "Trending Anime",
				Items: trendingAnime,
			},
		},
		"PopularSections": []models.MediaSection{
			{
				Title: "Popular Movies",
				Items: popularMovies,
				MoreURL: "/pages/discover/movies?sort=popular",
			},
			{
				Title: "Popular TV Shows",
				Items: popularTV,
				MoreURL: "/pages/discover/tv-shows?sort=popular",
			},
			{
				Title: "Popular Anime",
				Items: popularAnime,
				MoreURL: "/pages/discover/anime?sort=popular",
			},
		},
		"User": c.Locals("user"),
	}, "layout/default")
}
