package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/db"
	"github.com/flunks-nft/discord-bot/raid"
	"github.com/joho/godotenv"
)

var (
	discordToken   string
	DISCORD_GUI_ID string

	GuildID = ""
	s       *discordgo.Session

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
	}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load DISCORD_TOKEN from .env file
	discordToken = os.Getenv("DISCORD_TOKEN")

	DISCORD_GUI_ID = os.Getenv("DISCORD_GUI_ID")
}

func main() {
	// Connect to database & run migrations
	db.InitDB()

	// Create a new Discord session using the provided bot token.
	var err error
	s, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate functions as a callback for MessageCreate events.
	s.AddHandler(raid.FlowAddressHandler)
	s.AddHandler(raid.RaidMessageCreate)
	s.AddHandler(raid.ButtonInteractionCreate)

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
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
	s.Close()
}
