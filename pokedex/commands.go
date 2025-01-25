package main

import (
	"fmt"
	"os"
	"github.com/TheYorouzoya/boot-dev-golang/internal/pokeAPIHandler"
)

type cliCommand struct {
	name 			string
	description 	string
	callback 		func() error
}

type config struct {
	Next 		string
	Previous 	string
}


var commandRegistry map[string]cliCommand

func initCommandRegistry() {
	commandRegistry = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "Displays the next 20 city locations",
			callback: commandMap,
		},
	}
}


func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func commandHelp(config *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}


func commandMap(config *config) error {
	apiURL := config.Next		// get the next map page url

	apiResponse, err := pokeAPIHandler.QueryMap(apiURL)
	if err != nil {
		return err
	}

	config.Next = apiResponse.Next
	config.Previous = apiResponse.Previous

	for city := range apiResponse.Results {
		fmt.Println(city.Name)
	}

	return nil
}
