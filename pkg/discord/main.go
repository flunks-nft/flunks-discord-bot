package discord

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/utils"
)

var (
	ADMIN_AUTHOR_IDS utils.StringArray = []string{
		"594334378746707980", // Alfredoo
	}

	ALFREDOO_ID = "594334378746707980"

	DISCORD_TOKEN             string
	RAID_CHANNEL_ID           string
	RAID_LOG_CHANNEL_ID       string
	DISCORD_DAPPER_VERIFY_URL string

	dg                                  *discordgo.Session
	messageHandlers                     []func(s *discordgo.Session, i *discordgo.MessageCreate)
	interactionMessageComponentHandlers []func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

func init() {
	utils.LoadEnv()

	DISCORD_TOKEN = os.Getenv("DISCORD_TOKEN")
	RAID_CHANNEL_ID = os.Getenv("RAID_CHANNLE_ID")
	RAID_LOG_CHANNEL_ID = os.Getenv("RAID_LOG_CHANNEL_ID")
	DISCORD_DAPPER_VERIFY_URL = os.Getenv("DISCORD_DAPPER_VERIFY_URL")

	GuildID = DefaultGuiID()

	messageHandlers = []func(s *discordgo.Session, i *discordgo.MessageCreate){
		FlowAddressHandler,
		RaidMessageCreate,
	}

	interactionMessageComponentHandlers = []func(s *discordgo.Session, i *discordgo.InteractionCreate){
		ButtonInteractionCreate,    // for handling generic button interactions
		ButtonInteractionCreateOne, // for handling button interactions of individual Flunk in the yearbook
	}
}

// InitPureTextDiscord inits a Discord bot that only sends text messages
func InitPureTextDiscord(wg *sync.WaitGroup, done chan os.Signal) {
	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + DISCORD_TOKEN)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg = s

	// we only care about sending messages
	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Cleanly close down the Discord session.
	defer s.Close()

	// Wait here until CTRL-C or other term signal is received. (with <-done)
	log.Println("🌱 Discord Text Bot is now running. Press CTRL-C to exit.")
	<-done

	log.Println("🐒 Discord Text Bot Server is Gracefully shut down.")

	// Signal the WaitGroup that Discord service has finished
	wg.Done()
}

func InitDiscord(wg *sync.WaitGroup, done chan os.Signal) {
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
	for _, handler := range interactionMessageComponentHandlers {
		dg.AddHandler(handler)
	}

	// Add slash command handler
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// return if not a InteractionApplicationCommand command in case conflicts with other interactions
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("➕ Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// we only care about receiving message events
	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	registerSlashCommands()

	// Cleanly close down the Discord session.
	defer s.Close()

	// Wait here until CTRL-C or other term signal is received. (with <-done)
	log.Println("🌱 Bot is now running. Press CTRL-C to exit.")
	<-done

	log.Println("🐒 Discord Server is Gracefully shut down.")

	// Signal the WaitGroup that Discord service has finished
	wg.Done()
}
