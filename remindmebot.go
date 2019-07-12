package main

import (
	"bufio"
	"fmt"
	"os"
)

// global variables
var discordToken string
var tokensFile = "tokens.txt"

func setupTokens(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening \"%v\" - the file does not exist\n", fileName)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// increment scanner to the first token
	if scanner.Scan() {
		discordToken = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading \"%v\" - the file cannot be read\n", fileName)
		return false
	}

	return true
}

func main() {
	if setupTokens(tokensFile) {
		fmt.Printf("Here is your discord token - %v\n", discordToken)
	}
}
