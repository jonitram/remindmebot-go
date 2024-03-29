package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/araddon/dateparse"
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
// link to message is:
// https://discordapp.com/channels/<guildID>/<channelID>/<messageID>
// guildID == "" for PMs -> "@me" instead of guildID

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
		fmt.Printf("GuildID: %s ChannelID: %s Timestamp: %s Username: %s MessageID: %s Content: %s Current Time: %s\n", message.GuildID, message.ChannelID, message.Timestamp, message.Author.Username, message.ID, message.Content, time.Now().UTC().String())
		// var time string
		var reminder string
		var timestring string
		parameters := strings.TrimSpace(message.Content[len(commandPrefix):])
		if strings.Contains(parameters, "\"") {
			if strings.IndexByte(parameters, '"') == strings.LastIndexByte(parameters, '"') {
				// send update saying that only one " exists
				return
			}
			reminder = parameters[strings.IndexByte(parameters, '"') : strings.LastIndexByte(parameters, '"')+1]
			fmt.Printf(reminder)
		}
		if len(reminder) > 0 {
			parameters = strings.Trim(parameters, reminder)
		}
		timestring = strings.TrimSpace(parameters)
		fmt.Printf("\nTimestamp: " + timestring)
		if len(timestring) > 0 {
			parsedtime, err := dateparse.ParseAny(timestring)
			if err != nil {
				// parsedtime was successful
				fmt.Printf("\nParsed Timestring: " + parsedtime.String())
			} else {
				// check for other forms of parsed time
				// if not successful throw error
				// timenums := regexp.MustCompile("\\d+").FindStringSubmatch(timestring)
				// numberInt, err := strconv.Atoi(timenums)
				if err != nil {
					// send error message abt not including a number
				}
			}
		}
		// spin up goroutine with sleep using extracted parameters
		// on goroutine completion send message id through channel

		// write function for extracting all of this stuff
		// write function for determining time to sleep
		// have goroutine go to sleep for time
		// on startup loop through loaded json of messages and extract all of the data
		// only if the message exists tho
		// save a list of every message -
		// loop through every message and clear them out if they are done
		// on shutdown + startup
		// delete messages.txt file
		// spin up a goroutine for each message
		// goroutine checks for if remindme is still active or not
		// if so appends to message.txt
		// SORT THE LIST OF MESSAGES BY COMPLETION DATE
		// UPON COMPLETION OF A GOROUTINE - maybe kill a message?
	}
}

// func extractDigitIndices(input string) (int, int) {
// 	for i := 0; i < len(input); i++ {
// 		if unicode.IsNumber(input[i]) {
// 			for j := i; j < len(input); j++ {
// 				if !unicode.isNumber(input[j]) {
// 					return i, j
// 				}
// 			}
// 		}
// 	}
// 	return -1, -1
// }

func sortDuration(unit string) int {
	unit = strings.ToLower(unit)
	if strings.HasPrefix(unit, "s") {
		return 1
	} else if strings.HasPrefix(unit, "mi") {
		return 60
	} else if strings.HasPrefix(unit, "h") {
		return 3600
	} else if strings.HasPrefix(unit, "d") {
		return 86400
	} else if strings.HasPrefix(unit, "w") {
		return 604800
	} else if strings.HasPrefix(unit, "mo") {
		return 2592000
	} else if strings.HasPrefix(unit, "y") {
		return 31104000
	}
	return -1
}

func remindMe(message *discordgo.Message) {
	// parameters := regexp.MustCompile("\\s+").Split(message.Content, 4)
	// time, err := strconv.Atoi(parameters[1])
	// multiplier := sortDuration(parameters[2])

}

func main() {
	setupTokens(tokensFile)
	discord, err := discordgo.New("Bot " + discordToken)
	errCheck("error creating discord session", err)
	// user, err := discord.User(s"@me")
	errCheck("error retrieving user from discord", err)
	// commandPrefix = "<@" + user.ID + ">"
	commandPrefix = "rm!"
	fmt.Printf("Command Prefix: %s\n", commandPrefix)
	discord.AddHandler(messageHandler)
	errCheck("Error opening connection to Discord", discord.Open())
	defer discord.Close()
	<-make(chan struct{})
}
