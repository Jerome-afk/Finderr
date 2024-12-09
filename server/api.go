package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const tmdbBaseURL = "https://api.themoviedb.org/3"
const aniListBaseURL = "https://graphql.anilist.co"

func GetPopularMovies(apiKey string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/movie/popular?api_key=%s", tmdbBaseURL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []map[string]interface{} `json:"results"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

func GetPopularSeries(apiKey string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/tv/popular?api_key=%s", tmdbBaseURL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []map[string]interface{} `json:"results"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

func GetMovieDetails(apiKey string, movieID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/movie/%d?api_key=%s", tmdbBaseURL, movieID, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetTVSHOWDetails(apiKey string, tvID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/tv/%d?api_key=%s", tmdbBaseURL, tvID, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetPopularAnime(apiKey string) ([]map[string]interface{}, error) {
	query := `{
	    Media(type: ANIME, sort: POPULARITY){
		    id
			title {
			    romanji
			}
			description
			bannerImage
		}
	}`

	client := &http.Client{}
	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"query": "%s"}`, query)))
	req, err := http.NewRequest("POST", aniListBaseURL, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Media []map[string]interface{} `json:"Media"`
		} `json:"data"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Media, nil
}