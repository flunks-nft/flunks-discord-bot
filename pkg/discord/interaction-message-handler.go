package discord

// This modules is responsible for handling button interactions

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/utils"
)

// ButtonInteractionCreate handles the button click on the Raid Playground button interaction from users.
// Note it's supposed to send Ephemeral to the user but broadcast in the raid public information log channel.
func ButtonInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// return if not a InteractionMessageComponent incase conflicts with other interactions
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	if i.Type == discordgo.InteractionMessageComponent {
		switch i.MessageComponentData().CustomID {
		case "start_raid_all":
			handlesRaidAll(s, i)
			return
		case "yearbook":
			user, err := ValidateUser(i)
			if err != nil {
				respondeEphemeralMessage(s, i, err.Error())
				return
			}
			handlesYearbook(s, i, user)
			return
		case "leaderboard":
			handlesLeaderBoard(s, i)
			return
		case "next_flunk":
			user, err := ValidateUser(i)
			if err != nil {
				respondeEphemeralMessage(s, i, err.Error())
				return
			}
			handlesYearbook(s, i, user)
			return
		}
	}
}

func ValidateUser(i *discordgo.InteractionCreate) (db.User, error) {
	// Check user profile in the first place
	user, err := db.UserProfile(i.Member.User.ID)
	if err != nil {
		return user, errors.New("⚠️ Please use /dapper command to set up / update your Dapper wallet address.")
	}

	return user, nil
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

	if strings.Contains(customID, "raid_history") {
		handlesRaidHistory(s, i)
		return
	}
}

func handlesRaidAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Create a defer interaction message
	if err := deferEphemeralResponse(s, i); err != nil {
		return
	}

	// Check user profile in the first place
	user, err := ValidateUser(i)
	if err != nil {
		return
	}

	// Queue all Flunks for the raid
	msg, err := QueueForRaidAll(s, i, user)
	if err != nil {
		editTextResponse(s, i, err.Error())
		return
	}

	// Edit the original deferred interaction response with the new message
	editTextResponse(s, i, msg)
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
				templateIDInt, _ := utils.StringToUInt(templateID)
				if msg, err := QueueForRaidOne(s, i, templateIDInt); err != nil {
					respondeEphemeralMessage(s, i, err.Error())
				} else {
					respondeEphemeralMessage(s, i, msg)
				}
			}
		}
	}
}

// handlesYearbook is a handler for the Yearbook button to Zeero
func handlesYearbook(s *discordgo.Session, i *discordgo.InteractionCreate, user db.User) {
	// Defer interaction with placeholder Ephemeral msg to we have 15 minutes to respond to the original interaction
	if err := deferEphemeralResponse(s, i); err != nil {
		return
	}

	items, err := user.GetFlunks()
	if err != nil {
		editTextResponse(s, i, "⚠️ Failed to get your Flunks from Dapper.")
		return
	}

	if len(items) == 0 {
		editTextResponse(s, i, "⚠️ You don't have any Flunks in your Dapper wallet.")
		return
	}

	// returns next Flunk based on last fetched index
	totalCount := len(items)
	nextIndex := user.GetNextTokenIndex(totalCount)
	item := items[nextIndex]

	nft, err := db.GetNftByTemplateID(uint(item.TemplateID))

	if err != nil {
		editTextResponse(s, i, "⚠️ Failed to get your Flunks from Dapper.")
		return
	}

	// Edit the original deferred interaction response with the new message
	editEphemeralMessageWithMedia(s, i, nft)
}

// handlesLeaderBoard displays the leaderboard information of top 10 Flunks
func handlesLeaderBoard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferEphemeralResponse(s, i); err != nil {
		return
	}

	// Get top 10 Flunks from the database
	nfts := db.LeaderBoard()
	if len(nfts) == 0 {
		content := "⚠️ Failed to get the leaderboard information."
		editTextResponse(s, i, content)
		return
	}

	// Edit original ephemeral message with the leaderboard information
	if err := respondeEditFlunkLeaderBoard(s, i, nfts); err != nil {
		return
	}
}

func handlesRaidHistory(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferEphemeralResponse(s, i); err != nil {
		return
	}

	customIDParts := strings.Split(i.MessageComponentData().CustomID, "_")

	if i.Type == discordgo.InteractionMessageComponent {
		switch customIDParts[0] {
		case "raid":
			if customIDParts[1] == "history" {
				templateID := customIDParts[2]
				templateIDUInt, _ := utils.StringToUInt(templateID)
				records := db.GetRaidHistoryByTemplateID(templateIDUInt)
				if len(records) == 0 {
					msg := fmt.Sprintf("⚠️ No raid history found for Flunk #%s", templateID)
					editTextResponse(s, i, msg)
					return
				}

				// Create a string to store the concatenated records
				var recordsString string
				for _, record := range records {
					recordsString += record
				}
				editTextResponse(s, i, recordsString)
			}
		}
	}
}

// QueueForRaidAll adds all Flunks that belongs to a user to a raid queue for the challenge match
func QueueForRaidAll(s *discordgo.Session, i *discordgo.InteractionCreate, user db.User) (string, error) {
	// Retrieve Flunks for the user
	items, err := user.GetFlunks()
	if err != nil {
		return "", fmt.Errorf("⚠️ Failed to get your Flunks from Dapper: %v", err)
	}

	if len(items) == 0 {
		return "", fmt.Errorf("⚠️ You don't have any Flunks in your Dapper wallet.")
	}

	successArray := []string{}
	failedArray := []string{}

	// Queue them all
	for _, item := range items {
		msg, err := QueueForRaidOne(s, i, uint(item.TemplateID))
		if err != nil {
			failedArray = append(failedArray, msg)
			continue
		}
		successArray = append(successArray, msg)
	}

	succeedCnt := len(successArray)
	failedCnt := len(failedArray)

	if succeedCnt == 0 {
		return "", errors.New("⚠️ Failed to add any Flunks to the raid queue.")
	} else {
		if failedCnt == 0 {
			return fmt.Sprintf("✅ %v Flunks have been added to the raid queue.", succeedCnt), nil
		} else {
			return fmt.Sprintf("✅ %v Flunks have been added to the raid queue.\n⚠️ Failed to add %v Flunks.", succeedCnt, failedCnt), nil
		}
	}
}

// QueueForRaidOne adds Flunk to a raid queue for the challenge match
func QueueForRaidOne(s *discordgo.Session, i *discordgo.InteractionCreate, templateID uint) (string, error) {
	// Get NFT instance from database
	nft, err := db.GetNftByTemplateID(templateID)
	if err != nil {
		msg := fmt.Sprintf("⚠️ Syncing, please try later...")
		return "", errors.New(msg)
	}
	// TODO: Check if token is owned by the Discord user
	// Check if token has raided in the last 24 hours
	if isReady, nextValidRaidTime := nft.IsReadyForRaidQueue(); !isReady {
		msg := fmt.Sprintf("⚠️ NFT with tokenID %v is not ready for raid queue. Still %s hours remaining", templateID, nextValidRaidTime)
		return "", errors.New(msg)
	}
	// check if token is already in the raid queue
	if isInRaidQueue := nft.IsInRaidQueue(); isInRaidQueue {
		msg := fmt.Sprintf("⚠️ FLunk #%v is already in the raid queue.", templateID)
		return "", errors.New(msg)
	}
	// TODO: check if token is already in a raid
	if isRaiding := nft.IsRaiding(); isRaiding {
		msg := fmt.Sprintf("⚠️ FLunk #%v is already in a raid.", templateID)
		return "", errors.New(msg)
	}

	// Add Flunk to the match queue
	if err := nft.AddToRaidQueue(); err != nil {
		msg := fmt.Sprintf("⚠️ Failed to add Flunk #%v to the raid queue", templateID)
		return "", errors.New(msg)
	} else {
		msg := fmt.Sprintf("✅ Flunk %v has been added to the raid queue.", templateID)
		return msg, nil
	}
}
