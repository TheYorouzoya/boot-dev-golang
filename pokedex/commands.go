package main

import (
	"fmt"
	"os"

	"github.com/TheYorouzoya/boot-dev-golang/pokedex/internal/pokeAPIHandler"
	"github.com/TheYorouzoya/boot-dev-golang/pokedex/internal/pokecache"
)

type cliCommand struct {
	name 			string
	description 	string
	callback 		func(*config) error
}

type config struct {
	Next 		*string
	Previous 	*string
	Cache		*pokecache.Cache
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
		"mapb": {
			name: "mapb",
			description: "Goes to the previous 20 city locations",
			callback: commandMapb,
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

	apiResponse, err := pokeAPIHandler.QueryMap(apiURL, config.Cache)
	if err != nil {
		return err
	}

	config.Next = apiResponse.Next
	config.Previous = apiResponse.Previous

	for _, city := range apiResponse.Results {
		fmt.Println(city.Name)
	}

	return nil
}

func commandMapb(config *config) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	apiURL := config.Previous

	apiResponse, err := pokeAPIHandler.QueryMap(apiURL, config.Cache)
	if err != nil {
		return err
	}

	config.Next = apiResponse.Next
	config.Previous = apiResponse.Previous

	for _, city := range apiResponse.Results {
		fmt.Println(city.Name)
	}
	return nil
}
