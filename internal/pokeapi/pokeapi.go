package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"pokedex/internal/pokecache"
	"time"
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

type Pokedex struct {
	PokemonCaught map[string]Pokemon
}

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		Base_stat int `json:"base_stat"`
		Stat      struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func Map(cfg *Config, cache *pokecache.Cache, pokedex *Pokedex, args ...string) error {
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

func Mapb(cfg *Config, cache *pokecache.Cache, pokedex *Pokedex, args ...string) error {
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

func Explore(cfg *Config, cache *pokecache.Cache, pokedex *Pokedex, args ...string) error {
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
			return fmt.Errorf("failed to read response body. error: %w", err)
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

func Catch(cfg *Config, cache *pokecache.Cache, pokedex *Pokedex, args ...string) error {
	pokemonName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)
	body, entryFound := cache.Get(url)
	if !entryFound {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to get info on pokemon named: %v. error: %w", pokemonName, err)
		}
		body, err = io.ReadAll(res.Body)
		cache.Add(url, body)
		if err != nil {
			return fmt.Errorf("failed to read response body. error: %w", err)
		}
	}

	var pokemon Pokemon
	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		return fmt.Errorf("failed to unmarshal pokemon data. error: %w", err)
	}

	// Calculate the chance to catch
	var l float64 = 50         //	lowest base exp
	var h float64 = 306        // highest base exp
	var lChance float64 = 0.75 // base % chance of catching at L base exp
	var hChance float64 = 0.25 // base % chance of catching at H base exp
	var baseExp float64 = float64(pokemon.BaseExperience)
	t := (baseExp - l) / (h - l)
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	catchChance := lChance + t*(hChance-lChance)
	rolled := rand.Float64()

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	time.Sleep(100 * time.Millisecond)

	if rolled <= catchChance {
		fmt.Printf("%v was caught!\n", pokemon.Name)
		pokedex.Add(pokemon)
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}

	return nil
}

func (P *Pokedex) Add(newPokemon Pokemon) {
	_, ok := P.PokemonCaught[newPokemon.Name]
	if !ok {
		P.PokemonCaught[newPokemon.Name] = newPokemon
		fmt.Printf("%v added to pokedex! It can now be inspected.\n", newPokemon.Name)
	}
}

func ShowPokedex(cfg *Config, cache *pokecache.Cache, pokedex *Pokedex, args ...string) error {
	if len(pokedex.PokemonCaught) == 0 {
		fmt.Println("You have not yet caught any Pokemon!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range pokedex.PokemonCaught {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

func Inspect(cfg *Config, cache *pokecache.Cache, pokedex *Pokedex, args ...string) error {
	pokemonName := args[0]
	pokemon, ok := pokedex.PokemonCaught[pokemonName]
	if !ok {
		fmt.Printf("You don't yet have %v in your Pokedex.\n", pokemonName)
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats: \n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.Base_stat)
	}
	fmt.Printf("Types: \n")
	for _, pokeType := range pokemon.Types {
		fmt.Printf("  -%s\n", pokeType.Type.Name)
	}
	return nil
}
