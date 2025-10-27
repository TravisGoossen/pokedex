package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedex/internal/pokecache"
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

type locationAreaEncounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Config struct {
	NextUrl string
	PrevUrl string
}

func Map(cfg *Config, cache *pokecache.Cache, args ...string) error {
	var url string
	if cfg.NextUrl == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	if cfg.NextUrl != "" {
		url = cfg.NextUrl
	}

	body, entryFound := cache.Get(url)

	if !entryFound {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		cache.Add(url, body)
	}

	var locData locationArea
	err := json.Unmarshal(body, &locData)
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

func Mapb(cfg *Config, cache *pokecache.Cache, args ...string) error {
	if cfg.PrevUrl == "" {
		return fmt.Errorf("you're on the first page")
	}
	url := cfg.PrevUrl

	body, entryFound := cache.Get(url)

	if !entryFound {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		cache.Add(url, body)
	}

	var locData locationArea
	err := json.Unmarshal(body, &locData)
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

func Explore(cfg *Config, cache *pokecache.Cache, args ...string) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", args[0])
	body, entryFound := cache.Get(url)
	if !entryFound {
		res, err := http.Get(url)
		if res.StatusCode == 404 {
			return fmt.Errorf("'%s' is an invalid area-name", args[0])
		}
		if err != nil {
			return fmt.Errorf("failed to make http get request. error: %v", err)
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body. eror: %w", err)
		}
		cache.Add(url, body)
	}
	var encounters locationAreaEncounters
	err := json.Unmarshal(body, &encounters)
	if err != nil {
		return fmt.Errorf("failed to unmarshal area encounters. error: %w", err)
	}
	for _, pokemon := range encounters.PokemonEncounters {
		fmt.Printf("%v\n", pokemon.Pokemon.Name)
	}
	return nil
}
