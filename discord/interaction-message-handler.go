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
			QueueForRaidAll(s, i, user)
			return
		case "manage_wallet":
			respondeEphemeralMessage(s, i, "⚠️ Please use /dapper command to set up / update your Dapper wallet address.")
		case "yearbook":
			handlesYearbook(s, i, user)
			return
		case "lottery":
			respondeEphemeralMessage(s, i, "You clicked the 🍀 button.")
		case "start_raid_one":
			respondeEphemeralMessage(s, i, "You clicked the start_raid_one button.")
		case "next_flunk":
			handlesYearbook(s, i, user)
			return
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

	customID := i.MessageComponentData().CustomID
	if strings.Contains(customID, "start_raid_one") {
		handlesRaidOne(s, i)
		return
	}

	if strings.Contains(customID, "redirect_zeero") {
		handlesZeeroRedirect(s, i)
		return
	}
}

// Handles start raid command for individual Flunks
func handlesRaidOne(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				msg, _ := QueueForRaidOne(s, i, templateIDInt)
				respondeEphemeralMessage(s, i, msg)
			}
		}
	}
}

// handlesZeeroRedirect is a handler for the Yearbook button to Zeero
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
}

// handlesZeeroRedirect is a handler for the "Check on Zeero" button to Zeero
func handlesZeeroRedirect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customIDParts := strings.Split(i.MessageComponentData().CustomID, "_")

	if i.Type == discordgo.InteractionMessageComponent {
		switch customIDParts[0] {
		case "redirect":
			if customIDParts[1] == "zeero" {
				tokenID := customIDParts[2]
				msg := fmt.Sprintf("https://zeero.art/collection/flunks/%v", tokenID)
				respondeEphemeralMessage(s, i, msg)
			}
		}
	}
}

// QueueForRaidAll adds all Flunks that belongs to a user to a raid queue for the challenge match
func QueueForRaidAll(s *discordgo.Session, i *discordgo.InteractionCreate, user db.User) {
	// Retrieve Flunks for the user
	items, err := user.GetFlunks()
	if err != nil {
		respondeEphemeralMessage(s, i, "⚠️ Failed to get your Flunks from Dapper.")
		return
	}

	if len(items) == 0 {
		respondeEphemeralMessage(s, i, "⚠️ You don't have any Flunks in your Dapper wallet.")
		return
	}

	msgArray := []string{}

	// Queue them all
	for _, item := range items {
		msg, err := QueueForRaidOne(s, i, uint(item.TemplateID))
		if err == nil {
			msgArray = append(msgArray, msg)
		}
	}

	msg := fmt.Sprintf("✅ %v Flunks have been added to the raid queue.", len(msgArray))

	// Send a message to the user with all the Flunks that are queued
	respondeEphemeralMessage(s, i, msg)
}

// QueueForRaidOne adds Flunk to a raid queue for the challenge match
func QueueForRaidOne(s *discordgo.Session, i *discordgo.InteractionCreate, templateID uint) (string, error) {
	// Get NFT instance from database
	nft, err := db.GetNftByTemplateID(templateID)
	if err != nil {
		msg := fmt.Sprintf("⚠️ Syncing, please try later...")
		return msg, err
	}
	// TODO: Check if token is owned by the Discord user
	// Check if token has raided in the last 24 hours
	if isReady, nextValidRaidTime := nft.IsReadyForRaidQueue(); !isReady {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is not ready for raid queue. Still %s hours remaining", templateID, nextValidRaidTime)
		return msg, err
	}
	// check if token is already in the raid queue
	if isInRaidQueue := nft.IsInRaidQueue(); isInRaidQueue {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is already in the raid queue.", templateID)
		return msg, err
	}
	// TODO: check if token is already in a raid
	if isRaiding := nft.IsRaiding(); isRaiding {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is already in a raid.", templateID)
		return msg, err
	}

	// Add Flunk to the match queue
	if err := nft.AddToRaidQueue(); err != nil {
		msg := fmt.Sprintf("⚠️ Failed to add Flunk #%v to the raid queue", templateID)
		return msg, err
	} else {
		msg := fmt.Sprintf("%v", templateID)
		return msg, nil
	}
}
