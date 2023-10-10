package discord

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/utils"
)

// RaidMessageCreate creates an embedded message with buttons for users to interact with.
// Note it's only supposed to be used in the raid channel by admin users.
func RaidMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!raid-setup" {
		// Channel has to be the raid channel
		// TODO: send user an ephemeral message for visibility
		if m.ChannelID != RAID_CHANNEL_ID {
			return
		}

		defer DeleteMessage(s, m)

		// Check if the user is Alfred
		if ADMIN_AUTHOR_IDS.Contains(m.Author.ID) == false {
			// response := fmt.Sprintf("Only admins can use this command, ask  <@%s>", ALFREDOO_ID)
			// s.ChannelMessageSend(
			// 	m.ChannelID,
			// 	response,
			// )

			// No response, ignore and return
			return
		}

		// Create and admin message with buttons
		msg := &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "School Yard Battles BETA",
				Description: "Send Your Flunks to Daily Battle to Earn Rewards!",
				Image: &discordgo.MessageEmbedImage{
					URL: "https://storage.googleapis.com/zeero-public/raid_bot_face.png", // Replace with the actual image URL
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "BETA v1.0",
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Raid All",
							Style:    discordgo.PrimaryButton,
							CustomID: "start_raid_all",
						},
						discordgo.Button{
							Label:    "Yearbook",
							Style:    discordgo.SuccessButton,
							CustomID: "yearbook",
						},
						discordgo.Button{
							Label: "Manage Wallet",
							Style: discordgo.LinkButton,
							URL:   DISCORD_DAPPER_VERIFY_URL,
						},
						discordgo.Button{
							Label:    "🏆Leaderboard",
							Style:    discordgo.DangerButton,
							CustomID: "leaderboard",
						},
					},
				},
			},
		}

		// Send the message to the channel where the original command was received
		_, err := s.ChannelMessageSendComplex(m.ChannelID, msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}

	if strings.HasPrefix(m.Content, "!pixel") {
		re := regexp.MustCompile(`\d+`)
		matches := re.FindStringSubmatch(m.Content)

		if len(matches) > 0 {
			fmt.Println("Extracted number:", matches[0])
		} else {
			fmt.Println("No numeric field found.")
			return
		}

		pixelUri, exists := utils.PixelTemplateIdToUri[matches[0]]
		if !exists {
			// Silently ignore if pixelUri is not found
			return
		}

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("Flunk #%s", matches[0]),
			Image: &discordgo.MessageEmbedImage{URL: pixelUri},
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			fmt.Println("Error sending message: ", err)
		}

		return
	}

	if strings.HasPrefix(m.Content, "!battle") {
		defer DeleteMessage(s, m)

		re := regexp.MustCompile(`\d+`)
		matches := re.FindStringSubmatch(m.Content)

		if len(matches) > 0 {
			fmt.Println("Extracted number:", matches[0])
		} else {
			fmt.Println("No numeric field found.")
			return
		}

		battleIDStr := matches[0]
		battleIDInt, _ := utils.StringToInt(battleIDStr)

		raid := db.GetRaidByID(battleIDInt)
		if raid.ID == 0 {
			// Silently ignore if raid is not found
			return
		}
		PostRaidDetailsMsgUpdate(&raid, m.ChannelID)

		return
	}
}

// TODO: make this a Dapper sign - verify workflow
func FlowAddressHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.Contains(m.Content, "!dapper") {
		return
	}

	if m.ChannelID != RAID_CHANNEL_ID {
		return
	}

	defer DeleteMessage(s, m)

	flowAddress := strings.Split(m.Content, "!dapper ")[1]

	// Create or update Flow wallet address for user
	// TODO: validate the address
	db.CreateOrUpdateFlowAddress(m.Author.ID, flowAddress)
}
