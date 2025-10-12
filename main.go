package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commands = make(map[string]cliCommand)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		textInput := scanner.Text()
		cleanedText := cleanInput(textInput)
		switch cleanedText[0] {
		case "help":
			err := commands["help"].callback()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "exit":
			err := commands["exit"].callback()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
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

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, cmd := range commands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
