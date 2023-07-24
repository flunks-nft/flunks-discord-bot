package discord

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/utils"
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

	// Note that Zeero redirect is using tokenID as it's how the url is rendered
	zeeroButton := discordgo.Button{
		Label:    "Check on Zeero",
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("redirect_zeero_%v", item.TokenID),
	}

	nextButton := discordgo.Button{
		Label:    "Next Flunk",
		Style:    discordgo.PrimaryButton,
		CustomID: "next_flunk",
	}

	traits := item.Metadata.Traits()
	// Create a string to store the concatenated traits
	var traitsString string
	for _, trait := range traits {
		emoji := utils.TraitNameToEmoji[trait.Name]
		traitsString += fmt.Sprintf("%v %s: %s\n", emoji, trait.Name, trait.Value)
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
						Text: fmt.Sprintf("📚 Flunks # %v\n%s", item.TemplateID, traitsString),
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

func StringToInt(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Conversion error:", err)
		return -1, err
	}

	return num, nil
}

func StringToUInt(str string) (uint, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Conversion error:", err)
		return 0, err
	}

	return uint(num), nil
}
