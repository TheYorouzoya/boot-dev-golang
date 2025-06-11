package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"

	"github.com/TheYorouzoya/boot-dev-golang/pokedex/internal/pokeAPIHandler"
	"github.com/TheYorouzoya/boot-dev-golang/pokedex/internal/pokecache"
)

type cliCommand struct {
	name 			string
	description 	string
	callback 		func(*config, []string) error
}

type config struct {
	Next 		*string
	Previous 	*string
	Cache		*pokecache.Cache
	Pokedex		map[string]pokeAPIHandler.PokeAPIPokemonResponse
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
		"explore": {
			name: "explore",
			description: "Displays all the Pokemon in a map location\nUsage: explore <area-name>",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "Attempt to catch a pokemon\nUsage: catch <pokemon>",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "Displays the details of the pokemon from your pokedex\nUsage: inspect <pokemon>",
			callback: commandInspect,
		},
	}
}


func commandExit(config *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func commandHelp(config *config, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}


func commandMap(config *config, args []string) error {
	apiURL := config.Next		// get the next map page url

	apiResponse, err := pokeAPIHandler.QueryMap(apiURL, config.Cache)
	if err != nil {
		return err
	}

	config.Next = apiResponse.Next
	config.Previous = apiResponse.Previous

	for _, city := range apiResponse.Results {
		fmt.Println(*city.Name)
	}

	return nil
}

func commandMapb(config *config, args []string) error {
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
		fmt.Println(*city.Name)
	}
	return nil
}

func commandExplore(config *config, args []string) error {
	if args == nil || len(args) != 1 {
		return fmt.Errorf("Invalid number of arguments.\nUsage: explore <area_name>")
	}

	apiURL := "https://pokeapi.co/api/v2/location-area/" + args[0]

	apiResponse, err := pokeAPIHandler.QueryExplore(&apiURL, config.Cache)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\nFound Pokemon:\n", args[0])

	for _, pokemon := range apiResponse.PokemonEncounters {
		fmt.Println(*pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *config, args []string) error {
	if args == nil || len(args) != 1 {
		return fmt.Errorf("Invalid number of arguments./nUsage: catch <pokemon>")
	}

	pokemon := args[0]

	apiURL := "https://pokeapi.co/api/v2/pokemon/" + pokemon + "/"

	apiResponse, err := pokeAPIHandler.QueryPokemon(&apiURL, config.Cache)

	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	baseXP := apiResponse.BaseExperience

	chance := 1.0 / (1.0 + math.Log(float64(baseXP)))

	if rand.Float64() < chance {
		config.Pokedex[pokemon] = apiResponse
		fmt.Printf("%s was caught!\n", pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}

	return nil
}


func commandInspect(config *config, args []string) error {
	if args == nil || len(args) != 1 {
		return fmt.Errorf("Invalid number of arguments.\nUsage: inspect <pokemon>")
	}

	pokemon := args[0]

	pokedexEntry, exists := config.Pokedex[pokemon]

	if !exists {
		return fmt.Errorf("you have not caught that pokemon")
	}

	stats := ""

	for _, stat := range pokedexEntry.Stats {
		stats += fmt.Sprintf("   -%s: %d\n", *stat.Stat.Name, stat.BaseValue)
	}

	types := ""

	for _, typ := range pokedexEntry.Types {
		types += fmt.Sprintf("   -%s\n", *typ.Type.Name)
	}

	entry := fmt.Sprintf(
		`Name: %s
Height: %d
Weight: %d
Stats:
%s
Types:
%s`,
		pokemon,
		pokedexEntry.Height,
		pokedexEntry.Weight,
		stats,
		types)

	fmt.Print(entry)
	return nil
}
