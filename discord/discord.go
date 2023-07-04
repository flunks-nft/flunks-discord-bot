package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/db"
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

// RaidMessageDelete deletes the message sent by users in the channel.
func RaidMessageDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		log.Printf("Error deleting message: %v", err)
	}
}

// RaidMessageCreate creates an embedded message with buttons for users to interact with.
// Note it's only supposed to be used in the raid channel by admin users.
func RaidMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Content != "!raid-setup" {
		return
	}

	// Channel has to be the raid channel
	// TODO: send user an ephemeral message for visibility
	if m.ChannelID != RAID_CHANNEL_ID {
		return
	}
	defer RaidMessageDelete(s, m)

	// Check if the user is Alfred
	// TODO: maintain an admin list
	if ADMIN_AUTHOR_IDS.Contains(m.Author.ID) == false {
		response := fmt.Sprintf("Only admins can use this command, ask  <@%s>", ALFREDOO_ID)
		s.ChannelMessageSend(
			m.ChannelID,
			response,
		)
		return
	}

	// Create and admin message with buttons
	msg := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Flunks Raid Playground",
			Description: "Send Your Flunks to Daily Raids to Earn Rewards!",
			Image: &discordgo.MessageEmbedImage{
				URL: "https://storage.googleapis.com/zeero-public/arcade.png", // Replace with the actual image URL
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Raid",
						Style:    discordgo.PrimaryButton,
						CustomID: "start_raid",
					},
					discordgo.Button{
						Label:    "Manage Wallet",
						Style:    discordgo.SuccessButton,
						CustomID: "manage_wallet",
					},
					discordgo.Button{
						Label:    "Yearbook",
						Style:    discordgo.SecondaryButton,
						CustomID: "yearbook",
					},
					discordgo.Button{
						Label:    "🍀",
						Style:    discordgo.DangerButton,
						CustomID: "lottery",
					},
				},
			},
		},
	}

	// Send the message to the channel where the original command was received
	_, err := s.ChannelMessageSendComplex(m.ChannelID, msg)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

// ButtonInteractionCreate handles the button click on the Raid Playground button interaction from users.
// Note it's supposed to send Ephemeral to the user but broadcast in the raid public information log channel.
func ButtonInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check user profile in the first place
	user, err := db.UserProfile(i.Member.User.ID)
	if err != nil {
		respondeEphemeralMessage(s, i, "Use !dapper command to set up your Dapper wallet address")
		return
	}

	if i.Type == discordgo.InteractionMessageComponent {
		switch i.MessageComponentData().CustomID {
		case "start_raid":
			respondeEphemeralMessage(s, i, "You clicked the Raid button.")
		case "manage_wallet":
			respondeEphemeralMessage(s, i, "You clicked the Manage Wallet button.")
		case "yearbook":
			tokenIds := user.GetTokenIds()
			msg := fmt.Sprintf("You have %d Flunks.", len(tokenIds))
			respondeEphemeralMessage(s, i, msg)
		case "lottery":
			respondeEphemeralMessage(s, i, "You clicked the 🍀 button.")
		}
	}
}

// TODO: make this a Dapper sign - verify workflow
func FlowAddressHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.Contains(m.Content, "!dapper") {
		return
	}

	if m.ChannelID != RAID_CHANNEL_ID {
		return
	}

	defer RaidMessageDelete(s, m)

	flowAddress := strings.Split(m.Content, "!dapper ")[1]
	fmt.Println(flowAddress)

	// Create or update Flow wallet address for user
	// TODO: validate the address
	db.UpdateFlowAddress(m.Author.ID, flowAddress)
}

func respondeEphemeralMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	// Write a code to handle the button interaction
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   64, // Ephemeral
		},
	})
	if err != nil {
		log.Printf("Error handling start_raid: %v", err)
		return
	}
}
