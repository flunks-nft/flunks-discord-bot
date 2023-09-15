package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
)

// SendRaidConcludedMessageToRaidLogChannel sends a message when a raid is created
func SendMessageToRaidLogChannel(raid *db.Raid, nft1 db.Nft, nft2 db.Nft) {
	var fields []*discordgo.MessageEmbedField

	// Attach fist line challenge accepted message
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Challenge Accepted!", raid.ChallengeTypeEmoji()),
		Inline: false,
	})

	// Attach second line descriptions
	fields = append(fields, &discordgo.MessageEmbedField{
		Value:  fmt.Sprintf("<@%s> Flunk #%d ⚔️ <@%s> Flunk #%d", nft1.Owner.DiscordID, nft1.TemplateID, nft2.Owner.DiscordID, nft2.TemplateID),
		Inline: false,
	})

	embed := &discordgo.MessageEmbed{
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: nft1.Uri,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: nft2.Uri,
		},
	}

	_, err := dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed)
	if err != nil {
		fmt.Println("Error sending embedded message to channel:", err)
	}
}

func SendRaidConcludedMessageToRaidLogChannel(msgs []string, nft db.Nft, winnerThread, loserThread string) {
	var fields []*discordgo.MessageEmbedField

	// Attach the raid result text messages
	for _, msg := range msgs {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   msg,
			Inline: false,
		})
	}

	// Attach the winner and loser threads
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Winning Class",
		Value:  winnerThread,
		Inline: false,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Losing Class",
		Value:  loserThread,
		Inline: false,
	})

	// Attach the winner image
	embed := &discordgo.MessageEmbed{
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: nft.Uri,
		},
	}

	_, err := dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed)
	if err != nil {
		fmt.Println("Error sending embedded message to channel:", err)
	}
}
