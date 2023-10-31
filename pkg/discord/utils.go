package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/utils"
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

func editEphemeralMessageWithMedia(s *discordgo.Session, i *discordgo.InteractionCreate, nft db.Nft) {
	// Create the button component
	raidButton := discordgo.Button{
		Label:    "Battle",
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("start_raid_one_%v", nft.TemplateID),
	}

	// Note that Zeero redirect is using tokenID as it's how the url is rendered
	zeeroButton := discordgo.Button{
		Label: "Check on Zeero",
		Style: discordgo.LinkButton,
		URL:   fmt.Sprintf("https://zeero.art/collection/flunks/%v", nft.TokenID),
	}

	raidHistoryButton := discordgo.Button{
		Label:    "History",
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("raid_history_%v", nft.TemplateID),
	}

	nextButton := discordgo.Button{
		Label:    "Next Flunk",
		Style:    discordgo.PrimaryButton,
		CustomID: "next_flunk",
	}

	traits := nft.Traits

	// Create a string to store the concatenated traits
	var traitsString string
	for _, trait := range traits {
		emoji := utils.TraitNameToEmoji[trait.Name]
		traitsString += fmt.Sprintf("%v %s: %s\n", emoji, trait.Name, trait.Value)
	}

	// Add grade
	traitsString += fmt.Sprintf("%s %s: %d\n", "🍎", "Grade", nft.Grade())

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: nil,
		Components: &[]discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{&raidButton, &raidHistoryButton, &zeeroButton, &nextButton},
			},
		},
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Image: &discordgo.MessageEmbedImage{
					URL: nft.Uri,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("📚 Flunks # %v\n%s", nft.TemplateID, traitsString),
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

func respondeEditFlunkLeaderBoard(s *discordgo.Session, i *discordgo.InteractionCreate, nfts []db.Nft) error {
	// Parse leaderboard information into a string
	var msg string
	for idx, nft := range nfts {
		wins, losses := nft.GetWinsAndLosses()

		msg += fmt.Sprintf("[%d] Flunk **#%v**: %v Points (%d wins/ %d losses) \n", idx+1, nft.TemplateID, nft.Points, wins, wins+losses)
	}

	// Create an embed for the message
	embed := &discordgo.MessageEmbed{
		Title: "🏆 Leaderboard",
		Image: &discordgo.MessageEmbedImage{
			URL: "https://storage.googleapis.com/zeero-public/trophy.png",
		},
		Description: msg,
		Color:       0x0099ff, // light blue, in hexadecimal
	}

	// Edit the message to contain the embed
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			embed,
		},
	})
	if err != nil {
		log.Printf("Error handling leaderboard: %v", err)
		return err
	}
	return nil
}

func respondeEditFlunkRaidHistory(s *discordgo.Session, i *discordgo.InteractionCreate, msg string, uri string) error {
	// Create an embed for the message
	embed := &discordgo.MessageEmbed{
		Title: "🎮 Battle History",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: uri,
		},
		Description: msg,
		Color:       0x0099ff, // light blue, in hexadecimal
	}

	// Edit the message to contain the embed
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			embed,
		},
	})
	if err != nil {
		log.Printf("Error handling Raid History: %v", err)
		return err
	}
	return nil
}

func deferEphemeralResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Create a defer interaction message
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 64, // Ephemeral
		},
	})
	if err != nil {
		fmt.Println("Failed to defer interaction:", err)
		return err
	}
	return nil
}

func editTextResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		log.Printf("Error editing msg: %v", err)
	}
}
