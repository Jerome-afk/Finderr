package handlers

import (
	"math/rand/v2"

	"finderr/models"
	"finderr/services"

	"github.com/gofiber/fiber/v2"
)

func HomeHandler(c *fiber.Ctx) error {
	tmdb := services.NewTMDBClient()
	anilist := services.NewAniListClient()

	// Get trending content
	trendingMovies, _ := tmdb.GetTrendingMovies(5)
	trendingTV, _ := tmdb.GetTrendingTVShows(5)
	trendingAnime, _ := anilist.GetTrendingAnime(5)

	// Combine trending sections
	var trendingItems []models.MediaItem
	trendingItems = append(trendingItems, trendingMovies...)
	trendingItems = append(trendingItems, trendingTV...)
	trendingItems = append(trendingItems, trendingAnime...)

	// Shuffle the trending items
	rand.Shuffle(len(trendingItems), func(i, j int) {
		trendingItems[i], trendingItems[j] = trendingItems[j], trendingItems[i]
	})

	// Get Popular content
	popularMovies, _ := tmdb.GetPopularMovies()
	popularTV, _ := tmdb.GetPopularTVShows()
	popularAnime, _ := anilist.GetPopularAnime()

	return c.Render("pages/home", fiber.Map{
		"Title": "Home",
		"TrendingItems": trendingItems,
		"PopularSections": []models.MediaSection{
			{
				Title:   "Popular Movies",
				Items:   popularMovies,
				MoreURL: "/pages/discover/movies?sort=popular",
			},
			{
				Title:   "Popular TV Shows",
				Items:   popularTV,
				MoreURL: "/pages/discover/tv-shows?sort=popular",
			},
			{
				Title:   "Popular Anime",
				Items:   popularAnime,
				MoreURL: "/pages/discover/anime?sort=popular",
			},
		},
		"User": c.Locals("user"),
	}, "layout/default")
}
