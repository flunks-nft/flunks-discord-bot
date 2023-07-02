package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/raid"
	"github.com/joho/godotenv"
)

var discordToken string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load DISCORD_TOKEN from .env file
	discordToken = os.Getenv("DISCORD_TOKEN")
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate functions as a callback for MessageCreate events.
	dg.AddHandler(raid.PingPongMessageCreate)
	dg.AddHandler(raid.RaidMessageCreate)
	dg.AddHandler(raid.ButtonInteractionCreate)

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
