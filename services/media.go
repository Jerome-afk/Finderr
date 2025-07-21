package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"finderr/models"

	"github.com/Khan/genqlient/graphql"
)

type TMDBClient struct {
	BaseURL      string
	APIKey       string
	ImageBaseURL string
	client       *http.Client
}

type AniListClient struct {
	client graphql.Client
}

func NewTMDBClient() *TMDBClient {
	return &TMDBClient{
		BaseURL:      os.Getenv("TMDB_BASE_URL"),
		APIKey:       os.Getenv("TMDB_API_KEY"),
		ImageBaseURL: os.Getenv("TMDB_IMAGE_BASE_URL"),
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func NewAniListClient() *AniListClient {
	return &AniListClient{
		client: graphql.NewClient(os.Getenv("ANILIST_BASE_URL"), nil),
	}
}

func (t *TMDBClient) GetTrendingMovies(no int) ([]models.MediaItem, error) {
	url := fmt.Sprintf("%s/trending/movie/day?api_key=%s", t.BaseURL, t.APIKey)
	resp, err := t.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trending movies: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Results []models.MediaItem `json:"results"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Limit to 5 items
	if len(result.Results) > no {
		result.Results = result.Results[:no]
	}
	return result.Results, nil
}

func (t *TMDBClient) GetPopularMovies() ([]models.MediaItem, error) {
	url := fmt.Sprintf("%s/movie/popular?api_key=%s", t.BaseURL, t.APIKey)
	resp, err := t.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular movies: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Results []models.MediaItem `json:"results"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Limit to 15 items
	if len(result.Results) > 15 {
		result.Results = result.Results[:15]
	}
	return result.Results, nil
}

func (t *TMDBClient) GetTrendingTVShows(no int) ([]models.MediaItem, error) {
	url := fmt.Sprintf("%s/trending/tv/day?api_key=%s", t.BaseURL, t.APIKey)
	resp, err := t.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trending TV shows: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Results []models.MediaItem `json:"results"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	for i := range result.Results {
		result.Results[i].Type = "tv"
		if result.Results[i].Name == "" && result.Results[i].Title != "" {
			result.Results[i].Name = result.Results[i].Title
		}
	}

	// Limit to 5 items
	if len(result.Results) > no {
		result.Results = result.Results[:no]
	}
	return result.Results, nil
}

func (t *TMDBClient) GetPopularTVShows() ([]models.MediaItem, error) {
	url := fmt.Sprintf("%s/tv/popular?api_key=%s", t.BaseURL, t.APIKey)
	resp, err := t.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular TV shows: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Results []models.MediaItem `json:"results"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	for i := range result.Results {
		result.Results[i].Type = "tv"
		if result.Results[i].Name == "" && result.Results[i].Title != "" {
			result.Results[i].Name = result.Results[i].Title
		}
	}

	// Limit to 15 items
	if len(result.Results) > 15 {
		result.Results = result.Results[:15]
	}
	return result.Results, nil
}

func (a *AniListClient) GetTrendingAnime(no int) ([]models.MediaItem, error) {
	// GraphQL query to fetch trending anime
	query := `query {
		Page(page: 1, perPage: ` + fmt.Sprintf("%d", no) + `) {
			media(type: ANIME, sort: TRENDING_DESC) {
				id
				title {
					romaji
					english
				}
				coverImage {
				    large
				}
				averageScore
				description
			}
		}
	}`

	var response struct {
		Page struct {
			Media []struct {
				ID          int    `json:"id"`
				Title       struct {
					Romaji  string `json:"romaji"`
					English string `json:"english"`
				} `json:"title"`
				CoverImage struct {
					Large string `json:"large"`
				} `json:"coverImage"`
				AverageScore int `json:"averageScore"`
				Description string `json:"description"`
			} `json:"media"`
		} `json:"Page"`
	}

	err := a.client.MakeRequest(context.Background(), &graphql.Request{
		Query: query,
	}, &graphql.Response{
		Data: &response,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch trending anime: %w", err)
	}

	// Convert response to MediaItem format
	var items []models.MediaItem
	for _, m := range response.Page.Media {
		title := m.Title.English
		if title == "" {
			title = m.Title.Romaji
		}
		items = append(items, models.MediaItem{
			ID:           m.ID,
			Title:        title,
			Backdrop:   m.CoverImage.Large,
			Poster:       m.CoverImage.Large, // Assuming CoverImage is used as Poster
			Overview:    m.Description,
			Rating: 	 float64(m.AverageScore) / 10.0, // Convert to float64 for consistency
			Type:         "anime",
		})
	}

	return items, nil
}

func (a *AniListClient) GetPopularAnime() ([]models.MediaItem, error) {
	// GraphQL query to fetch popular anime
	query := `query {
		Page(page: 1, perPage: 15) {
			media(type: ANIME, sort: POPULARITY_DESC) {
				id
				title {
					romaji
					english
				}
				coverImage {
				    large
				}
				averageScore
				description
			}
		}
	}`

	var response struct {
		Page struct {
			Media []struct {
				ID          int    `json:"id"`
				Title       struct {
					Romaji  string `json:"romaji"`
					English string `json:"english"`
				} `json:"title"`
				CoverImage struct {
					Large string `json:"large"`
				} `json:"coverImage"`
				AverageScore int `json:"averageScore"`
				Description string `json:"description"`
			} `json:"media"`
		} `json:"Page"`
	}

	err := a.client.MakeRequest(context.Background(), &graphql.Request{
		Query: query,
	}, &graphql.Response{
		Data: &response,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular anime: %w", err)
	}

	var items []models.MediaItem
	for _, m := range response.Page.Media {
		title := m.Title.English
		if title == "" {
			title = m.Title.Romaji
		}
		items = append(items, models.MediaItem{
			ID:           m.ID,
			Title:        title,
			Backdrop:     m.CoverImage.Large,
			Poster:       m.CoverImage.Large, // Assuming CoverImage is used as Poster
			Overview:     m.Description,
			Rating: 	 float64(m.AverageScore) / 10.0, // Convert to float64 for consistency
			Type:         "anime",
		})
	}

	return items, nil
}
