package admin

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/helper"
)

var (
	ADMIN_AUTHOR_IDS = []string{
		"594334378746707980", // Alfredoo
	}

	ALFREDOO_ID = "594334378746707980"
)

func RaidMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Content != "!raid-setup" {
		return
	}

	// Check if the user is Alfred
	// TODO: maintain an admin list
	if helper.Contains(ADMIN_AUTHOR_IDS, m.Author.ID) == false {
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
						Label:    "Check Flunks",
						Style:    discordgo.SecondaryButton,
						CustomID: "check_flunks",
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
