package discord

// This modules is responsible for handling button interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/db"
)

// ButtonInteractionCreate handles the button click on the Raid Playground button interaction from users.
// Note it's supposed to send Ephemeral to the user but broadcast in the raid public information log channel.
func ButtonInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// return if not a InteractionMessageComponent incase conflicts with other interactions
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// Check user profile in the first place
	user, err := db.UserProfile(i.Member.User.ID)
	if err != nil {
		respondeEphemeralMessage(s, i, "⚠️ Please use /dapper command to set up / update your Dapper wallet address.")
		return
	}

	if i.Type == discordgo.InteractionMessageComponent {
		switch i.MessageComponentData().CustomID {
		case "start_raid":
			respondeEphemeralMessage(s, i, "You clicked the Raid button.")
		case "manage_wallet":
			respondeEphemeralMessage(s, i, "⚠️ Please use /dapper command to set up / update your Dapper wallet address.")
		case "yearbook":
			tokenIds := user.GetTokenIds()
			msg := fmt.Sprintf("You have %d Flunks.", len(tokenIds))
			respondeEphemeralMessage(s, i, msg)
		case "lottery":
			respondeEphemeralMessage(s, i, "You clicked the 🍀 button.")
		}
	}
}
