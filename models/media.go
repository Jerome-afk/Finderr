package models

type MediaItem struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	Poster         string     `json:"poster_path"`
	Name		   string     `json:"name"`
	Backdrop       string     `json:"backdrop_path"`
	Rating         float64    `json:"vote_average"`
	Overview       string     `json:"overview"`
	ReleaseDate    string     `json:"release_date"`
	FirstAirDate   string     `json:"first_air_date"`
	Type           string     `json:"media_type"`
}

type MediaSection struct {
	Title    string
	Items    []MediaItem
	MoreURL  string
}

func (m *MediaItem) GetYear() string {
	dateStr := ""
	if m.Type == "tv" {
		dateStr = m.FirstAirDate
	} else {
		dateStr = m.ReleaseDate
	}

	if len(dateStr) >= 4 {
		return dateStr[:4]
	}

	return ""
}