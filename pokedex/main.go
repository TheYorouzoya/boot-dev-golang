package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TheYorouzoya/boot-dev-golang/pokedex/internal/pokecache"
)

func main() {
	const prompt string = "Pokedex > "
	const DEFAULT_CACHE_TIMEOUT = time.Duration(5000 * time.Millisecond)
	var APICache = pokecache.NewCache(DEFAULT_CACHE_TIMEOUT)

	initCommandRegistry()
	defaultAPIMapPath := "https://pokeapi.co/api/v2/location-area/"
	apiConfig := config{
		Next: &defaultAPIMapPath,
		Previous: nil,
		Cache: &APICache,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for ;; {
		fmt.Print(prompt)
		scanner.Scan()
		userInput := scanner.Text()
		userInput = strings.ToLower(userInput)
		splitInput := cleanInput(userInput)
		command := splitInput[0]

		register, ok := commandRegistry[command]

		if !ok {
			fmt.Println("Unknown command")
		} else {
			err := register.callback(&apiConfig)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}


func cleanInput(text string) []string {

	if len(text) <= 0 {
		return []string{}
	}

	cleaned := strings.TrimSpace(text)
	split := strings.Fields(cleaned)
	return split

}
