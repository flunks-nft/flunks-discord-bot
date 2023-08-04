package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
)

// RaidMessageCreate creates an embedded message with buttons for users to interact with.
// Note it's only supposed to be used in the raid channel by admin users.
func RaidMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	fmt.Println(m.ChannelID)

	if m.Author.ID == s.State.User.ID || m.Content != "!raid-setup" {
		return
	}

	// Channel has to be the raid channel
	// TODO: send user an ephemeral message for visibility
	if m.ChannelID != RAID_CHANNEL_ID {
		return
	}
	defer DeleteMessage(s, m)

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
						Label:    "Raid All",
						Style:    discordgo.PrimaryButton,
						CustomID: "start_raid_all",
					},
					discordgo.Button{
						Label:    "Yearbook",
						Style:    discordgo.SecondaryButton,
						CustomID: "yearbook",
					},
					discordgo.Button{
						Label:    "Manage Wallet",
						Style:    discordgo.SuccessButton,
						CustomID: "manage_wallet",
					},
					discordgo.Button{
						Label:    "🏆Leaderboard",
						Style:    discordgo.DangerButton,
						CustomID: "leaderboard",
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

	defer DeleteMessage(s, m)

	flowAddress := strings.Split(m.Content, "!dapper ")[1]

	// Create or update Flow wallet address for user
	// TODO: validate the address
	db.CreateOrUpdateFlowAddress(m.Author.ID, flowAddress)
}
