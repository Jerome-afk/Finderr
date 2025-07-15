package models

type MediaItem struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	Poster         string     `json:"poster_path"`
	CoverImage     string     `json:"cover_image"`
	Rating         float64    `json:"vote_average"`
	AverageScore   int        `json:"average_score"`
	Type           string     `json:"type"`
}

type MediaSection struct {
	Title    string
	Items    []MediaItem
	MoreURL  string
}