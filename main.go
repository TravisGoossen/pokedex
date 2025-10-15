package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pokeapi "pokedex/internal/pokeapi"
)

var commands = make(map[string]cliCommand)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
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

	var cfg pokeapi.Config

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		textInput := scanner.Text()
		cleanedText := cleanInput(textInput)
		switch cleanedText[0] {
		case "help":
			err := commands["help"].callback(&cfg)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		case "map":
			err := commands["map"].callback(&cfg)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		case "mapb":
			err := commands["mapb"].callback(&cfg)
			if err != nil {
				fmt.Println(err)
			}
		case "exit":
			err := commands["exit"].callback(&cfg)
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

func commandExit(cfg *pokeapi.Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *pokeapi.Config) error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, cmd := range commands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
