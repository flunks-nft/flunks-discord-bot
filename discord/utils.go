package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
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

// DeleteMessage deletes the message sent by users in the channel.
func DeleteMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		log.Printf("Error deleting message: %v", err)
	}
}
