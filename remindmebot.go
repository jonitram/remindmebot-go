package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// global variables
var tokensFile = "tokens.txt"
var discordToken string
var commandPrefix string

// remindmes = list of structs w/ author, time message will execute (post-converted), remind message, goroutine
// could also do a map of author -> list of remindmes
// remindmes = append(remindmes, <thing to add>)

// right now the work flow is going to look like this:
// on receive message:
// check for command
// check for formatting of command
// spin up background goroutine
// 	- goroutine should take in:
// 	author
// 	time request
// 	desired message
// add goroutine to master list
// send remindme confirmation message
// react to remindme confirmation message with options

// on reaction:
// check content for remindme confirmation update + author = the bot
// check reactor for original message author
// handle reaction

// goroutine will schedule job based on params passed in
// check to make sure message still exists -> then remindmes?

// could use standardized date / time and sleep by subtracting current time from it
// that way could resume saved and restored jobs on startup

// here are some potential options:
// delete remindme (will delete the confirmation and original message)
// delete command (just delete original message command)

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

func messageHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if strings.HasPrefix(message.Content, commandPrefix) {
		fmt.Printf("ChannelID: %s Username: %s MessageID: %s Content: %s\n", message.ChannelID, message.Author.Username, message.ID, message.Content)
	}
}

func main() {
	setupTokens(tokensFile)
	discord, err := discordgo.New("Bot " + discordToken)
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	commandPrefix = "<@" + user.ID + ">"
	fmt.Printf("Command Prefix: %s\n", commandPrefix)
	discord.AddHandler(messageHandler)
	errCheck("Error opening connection to Discord", discord.Open())
	defer discord.Close()
	<-make(chan struct{})
}
