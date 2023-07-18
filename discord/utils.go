package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/zeero"
)

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

func respondeEphemeralMessageWithMedia(s *discordgo.Session, i *discordgo.InteractionCreate, item zeero.NftDtoWithActivity) {
	// Create the button component
	raidButton := discordgo.Button{
		Label:    "Raid",
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("start_raid_one_%v", item.TemplateID),
	}

	zeeroButton := discordgo.Button{
		Label:    "Check on Zeero",
		Style:    discordgo.PrimaryButton,
		CustomID: "redirect_zeero",
	}

	nextButton := discordgo.Button{
		Label:    "Next Flunk",
		Style:    discordgo.PrimaryButton,
		CustomID: "next_flunk",
	}

	// Write a code to handle the button interaction
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 64, // Ephemeral
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{&raidButton, &zeeroButton, &nextButton},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Image: &discordgo.MessageEmbedImage{
						URL: item.Metadata.URI,
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("📚 Flunks # %v", item.TemplateID),
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error handling start_raid: %v", err)
		return
	}
}

// DeleteMessage deletes the message sent by users in the channel.
func DeleteMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		log.Printf("Error deleting message: %v", err)
	}
}

func DefaultGuiID() *string {
	s := ""
	return &s
}
