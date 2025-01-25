package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func main() {
	const prompt string = "Pokedex > "

	initCommandRegistry()
	apiConfig := config{
		Next: "https://pokeapi.co/api/v2/location-area/",
		Previous: "null",
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
