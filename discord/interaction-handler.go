package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/db"
)

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
