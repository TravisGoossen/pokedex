package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
)

var commands = make(map[string]cliCommand)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config, *pokecache.Cache, *pokeapi.Pokedex, ...string) error
}

func main() {
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the pokedex",
		callback:    commandExit,
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "List the next 20 map location areas",
		callback:    pokeapi.Map,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "List the previous 20 map location areas",
		callback:    pokeapi.Mapb,
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "List the pokemon that can be found in an area. Proper use: 'explore area-name'",
		callback:    pokeapi.Explore,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Attempt to catch the named pokemon",
		callback:    pokeapi.Catch,
	}

	var cfg pokeapi.Config
	cache := pokecache.NewCache(5 * time.Second)
	var pokedex pokeapi.Pokedex
	pokedex.PokemonCaught = make(map[string]pokeapi.Pokemon)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		textInput := scanner.Text()
		cleanedText := cleanInput(textInput)
		switch cleanedText[0] {
		case "help":
			err := commands["help"].callback(&cfg, cache, &pokedex)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		case "map":
			err := commands["map"].callback(&cfg, cache, &pokedex)
			if err != nil {
				fmt.Println(err)
			}
		case "mapb":
			err := commands["mapb"].callback(&cfg, cache, &pokedex)
			if err != nil {
				fmt.Println(err)
			}
		case "explore":
			if len(cleanedText) < 2 {
				fmt.Println("No area name provided. Proper use: 'explore area-name'")
			} else if len(cleanedText) >= 3 {
				fmt.Println("The area name cannot contain spaces. Proper use: 'explore area-name'")
			} else {
				err := commands["explore"].callback(&cfg, cache, &pokedex, cleanedText[1])
				if err != nil {
					fmt.Println(err)
				}
			}
		case "catch":
			if len(cleanedText) < 2 {
				fmt.Println("No Pokemon entered. Proper use: 'catch pokemon-name'")
			} else if len(cleanedText) >= 3 {
				fmt.Println("Only one Pokemon can be caught at a time. Proper use: 'catch pokemon-name'")
			} else {
				err := commands["catch"].callback(&cfg, cache, &pokedex, cleanedText[1])
				if err != nil {
					fmt.Println(err)
				}
			}
		case "exit":
			err := commands["exit"].callback(&cfg, cache, &pokedex)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		default:
			fmt.Print("Unknown command\n")
		}
	}
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	finalText := strings.Fields(lower)
	return finalText
}

func commandExit(cfg *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, args ...string) error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, cmd := range commands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
