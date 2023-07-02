package raid

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/helper"
)

var (
	ADMIN_AUTHOR_IDS = []string{
		"594334378746707980", // Alfredoo
	}

	ALFREDOO_ID = "594334378746707980"

	RAID_CHANNEL_ID     = os.Getenv("RAID_CHANNEL_ID")
	RAID_LOG_CHANNEL_ID = os.Getenv("RAID_LOG_CHANNEL_ID")
)

func RaidMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Content != "!raid-setup" {
		return
	}

	// Channel has to be the raid channel
	// TODO: send user an ephemeral message for visibility
	if m.ChannelID != RAID_CHANNEL_ID {
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
	if i.Type == discordgo.InteractionMessageComponent {
		switch i.MessageComponentData().CustomID {
		case "start_raid":
			// Write a code to handle the "start_raid" button interaction
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("<@%s> clicked the Raid button.", i.Member.User.ID),
					Flags:   64, // Ephemeral
				},
			})
			if err != nil {
				log.Printf("Error handling start_raid: %v", err)
				return
			}
		case "manage_wallet":
			// Write a code to handle the "manage_wallet" button interaction
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("<@%s> clicked the Manage Wallet.", i.Member.User.ID),
					Flags:   64, // Ephemeral
				},
			})
			if err != nil {
				log.Printf("Error handling manage_wallet: %v", err)
				return
			}
		case "yearbook":
			// Write a code to handle the "yearbook" button interaction
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("<@%s> clicked the Yearbook button.", i.Member.User.ID),
					Flags:   64, // Ephemeral
				},
			})
			if err != nil {
				log.Printf("Error handling yearbook: %v", err)
				return
			}
		case "lottery":
			// Write a code to handle the "lottery" button interaction
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("<@%s> clicked the 🍀 button.", i.Member.User.ID),
					Flags:   64, // Ephemeral
				},
			})
			if err != nil {
				log.Printf("Error handling lottery: %v", err)
				return
			}
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
//
// It is called whenever a message is created but only when it's sent through a
// server as we did not request IntentsDirectMessages.
func PingPongMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// In this example, we only care about messages that are "ping".
	// Note that Message.Content is only available when the intent is enabled on Discord Developer Portal.
	// Ref: https://github.com/bwmarrin/discordgo/issues/961#issuecomment-1565032340
	if m.Content != "ping" {
		return
	}

	// Reply with "Pong!" in the same channel where the user posted "ping".
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		// If an error occurred, we failed to send the message.
		fmt.Println("error sending message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send the message!",
		)
	}
}
