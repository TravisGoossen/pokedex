package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	NextUrl string
	PrevUrl string
}

func Map(cfg *Config) error {
	var url string
	if cfg.NextUrl == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	if cfg.NextUrl != "" {
		url = cfg.NextUrl
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	var locData locationArea
	err = json.Unmarshal(body, &locData)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	cfg.NextUrl = locData.Next
	cfg.PrevUrl = locData.Previous

	for _, result := range locData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func Mapb(cfg *Config) error {
	if cfg.PrevUrl == "" {
		return fmt.Errorf("you're on the first page")
	}
	url := cfg.PrevUrl

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	var locData locationArea
	err = json.Unmarshal(body, &locData)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	cfg.NextUrl = locData.Next
	cfg.PrevUrl = locData.Previous

	for _, location := range locData.Results {
		fmt.Println(location.Name)
	}
	return nil
}
