package discord

// This modules is responsible for handling button interactions

import (
	"fmt"
	"log"
	"strings"

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
		case "start_raid_all":
			respondeEphemeralMessage(s, i, "⚠️ TODO: queue all of your Flunks to Raid.")
		case "manage_wallet":
			respondeEphemeralMessage(s, i, "⚠️ Please use /dapper command to set up / update your Dapper wallet address.")
		case "yearbook":
			handlesYearbook(s, i, user)
		case "lottery":
			respondeEphemeralMessage(s, i, "You clicked the 🍀 button.")
		case "start_raid_one":
			respondeEphemeralMessage(s, i, "You clicked the start_raid_one button.")
		case "next_flunk":
			handlesYearbook(s, i, user)
		}
	}
}

func ButtonInteractionCreateOne(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// return if not a InteractionMessageComponent incase conflicts with other interactions
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// Check user profile in the first place
	_, err := db.UserProfile(i.Member.User.ID)
	if err != nil {
		respondeEphemeralMessage(s, i, "⚠️ Please use /dapper command to set up / update your Dapper wallet address.")
		return
	}

	customIDParts := strings.Split(i.MessageComponentData().CustomID, "_")
	if len(customIDParts) < 4 {
		log.Printf("Invalid custom ID: %s", i.MessageComponentData().CustomID)
		return
	}

	if i.Type == discordgo.InteractionMessageComponent {
		switch customIDParts[0] {
		case "start":
			if customIDParts[1] == "raid" && customIDParts[2] == "one" {
				// You can use customIdParts[3] which should be the tokenID
				// Replace the line below with your desired function handling the specific tokenID
				templateID := customIDParts[3]
				templateIDInt, _ := StringToUInt(templateID)
				respondeEphemeralMessage(s, i, fmt.Sprintf("Starting Raid for Flunk: %s", templateID))
				QueueForRaid(s, i, templateIDInt)
			}
		}
	}
}

func handlesYearbook(s *discordgo.Session, i *discordgo.InteractionCreate, user db.User) {
	items, err := user.GetFlunks()
	if err != nil {
		respondeEphemeralMessage(s, i, "⚠️ Failed to get your Flunks from Dapper.")
		return
	}

	if len(items) == 0 {
		respondeEphemeralMessage(s, i, "⚠️ You don't have any Flunks in your Dapper wallet.")
		return
	}

	// returns next Flunk based on last fetched index
	totalCount := len(items)
	nextIndex := user.GetNextTokenIndex(totalCount)
	item := items[nextIndex]
	respondeEphemeralMessageWithMedia(s, i, item)
	return
}

func QueueForRaid(s *discordgo.Session, i *discordgo.InteractionCreate, templateID uint) {
	// Get NFT instance from database
	nft, err := db.GetNft(templateID)
	if err != nil {
		msg := fmt.Sprintf("⚠️ Flunk #%v not found.", templateID)
		respondeEphemeralMessage(s, i, msg)
		return
	}
	// TODO: Check if token is owned by the Discord user
	// Check if token has raided in the last 24 hours
	if isReady, nextValidRaidTime := nft.IsReadyForRaidQueue(); !isReady {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is not ready for raid queue. Still %s hours remaining", templateID, nextValidRaidTime)
		respondeEphemeralMessage(s, i, msg)
		return
	}
	// check if token is already in the raid queue
	if isInRaidQueue := nft.IsInRaidQueue(); isInRaidQueue {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is already in the raid queue.", templateID)
		respondeEphemeralMessage(s, i, msg)
		return
	}
	// TODO: check if token is already in a raid
	if isRaiding := nft.IsRaiding(); isRaiding {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is already in a raid.", templateID)
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
			msg := fmt.Sprintf("Flunk #%v is in the raid queue", templateID)
			respondeEphemeralMessage(s, i, msg)
			break
		} else {
			// TODO: log error
		}
	}

	// Otherwise just add to the match queue
	if err := nft.AddToRaidQueue(); err != nil {
		msg := fmt.Sprintf("⚠️ Failed to add Flunk #%v to the raid queue", templateID)
		respondeEphemeralMessage(s, i, msg)
		return
	} else {
		msg := fmt.Sprintf("Flunk #%v is in the raid queue", templateID)
		respondeEphemeralMessage(s, i, msg)
		return
	}
}
