package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// global variables
var tokensFile = "tokens.txt"
var discordToken string
var commandPrefix string
var remindmes []string

// remindmes = append(remindmes, <thing to add>)

func setupTokens(fileName string) {
	file, err := os.Open(fileName)
	errCheck("Error opening \""+fileName+"\" - the file does not exist", err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// increment scanner to the first token
	if scanner.Scan() {
		discordToken = scanner.Text()
	}

	errCheck("Error reading \""+fileName+"\" - the file cannot be read", scanner.Err())

	return
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v\n", msg, err)
		panic(err)
	}
	return
}

func main() {
	setupTokens(tokensFile)
	discord, err := discordgo.New("Bot " + discordToken)
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	commandPrefix = "@" + user.Username + "#" + user.Discriminator
	fmt.Println(commandPrefix)
}
