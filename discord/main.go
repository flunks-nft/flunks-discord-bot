package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/utils"
)

var (
	ADMIN_AUTHOR_IDS utils.StringArray = []string{
		"594334378746707980", // Alfredoo
	}

	ALFREDOO_ID = "594334378746707980"

	DISCORD_TOKEN       string
	RAID_CHANNEL_ID     string
	RAID_LOG_CHANNEL_ID string

	GUI_ID *string

	dg                  *discordgo.Session
	messageHandlers     []func(s *discordgo.Session, i *discordgo.MessageCreate)
	interactionHandlers []func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

func init() {
	utils.LoadEnv()

	DISCORD_TOKEN = os.Getenv("DISCORD_TOKEN")
	RAID_CHANNEL_ID = os.Getenv("RAID_CHANNLE_ID")
	RAID_LOG_CHANNEL_ID = os.Getenv("RAID_LOG_CHANNEL_ID")
	DISCORD_TOKEN = os.Getenv("DISCORD_TOKEN")

	GUI_ID = DefaultGuiID()

	messageHandlers = []func(s *discordgo.Session, i *discordgo.MessageCreate){
		FlowAddressHandler,
		RaidMessageCreate,
	}

	interactionHandlers = []func(s *discordgo.Session, i *discordgo.InteractionCreate){
		ButtonInteractionCreate,
	}
}

func InitDiscord() {
	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + DISCORD_TOKEN)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg = s

	// Register message handlers
	for _, handler := range messageHandlers {
		dg.AddHandler(handler)
	}

	// Register interaction handlers
	for _, handler := range interactionHandlers {
		dg.AddHandler(handler)
	}

	// we only care about receiving message events
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
