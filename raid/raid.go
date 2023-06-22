package raid

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// RaidMessageCreate will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
//
// It is called whenever a "!raid" message is created but only when it's sent through a
// server as we did not request IntentsDirectMessages.
func RaidMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "!raid" {
		return
	}

	response := fmt.Sprintf("Yo <@%s>, you flunk started a raid!", m.Author.ID)
	_, err := s.ChannelMessageSend(m.ChannelID, response)
	if err != nil {
		fmt.Println("error sending message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send the message!",
		)
	}

	// TODO: implement actual raid functionality
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
