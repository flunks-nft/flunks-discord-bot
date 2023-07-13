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
			// TODO: figure out how to get the tokenID from the message
			QueueForRaid(s, i, 1)
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

func QueueForRaid(s *discordgo.Session, i *discordgo.InteractionCreate, tokenID uint) {
	// Get NFT instance from database
	nft, err := db.GetNft(tokenID)
	if err != nil {
		msg := fmt.Sprintf("⚠️ Flunk #%v not found.", tokenID)
		respondeEphemeralMessage(s, i, msg)
		return
	}
	// TODO: Check if token is owned by the Discord user
	// Check if token has raided in the last 24 hours
	if isReady, nextValidRaidTime := nft.IsReadyForRaidQueue(); !isReady {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is not ready for raid queue. Still %s hours remaining", tokenID, nextValidRaidTime)
		respondeEphemeralMessage(s, i, msg)
		return
	}
	// check if token is already in the raid queue
	if isInRaidQueue := nft.IsInRaidQueue(); isInRaidQueue {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is already in the raid queue.", tokenID)
		respondeEphemeralMessage(s, i, msg)
		return
	}
	// TODO: check if token is already in a raid
	if isRaiding := nft.IsRaiding(); isRaiding {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is already in a raid.", tokenID)
		respondeEphemeralMessage(s, i, msg)
		return
	}

	// Try to find a match first
	// TODO: add this to the worker
	retryCounter := 0
	for retryCounter < 10 {
		err := nft.RaidMatch()
		retryCounter += 1
		if err == nil {
			msg := fmt.Sprintf("Flunk #%v is in the raid queue", tokenID)
			respondeEphemeralMessage(s, i, msg)
			break
		} else {
			// TODO: log error
		}
	}

	// Otherwise just add to the match queue
	if err := nft.AddToRaidQueue(); err != nil {
		msg := fmt.Sprintf("⚠️ Failed to add Flunk #%v to the raid queue", tokenID)
		respondeEphemeralMessage(s, i, msg)
		return
	} else {
		msg := fmt.Sprintf("Flunk #%v is in the raid queue", tokenID)
		respondeEphemeralMessage(s, i, msg)
		return
	}
}
