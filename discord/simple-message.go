package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/db"
)

// SendRaidConcludedMessageToRaidLogChannel sends a message when a raid is concluded
func SendRaidConcludedMessageToRaidLogChannel(msg string, nft db.Nft) {
	sendMessageToRaidLogChannel(msg)
	sendPureFlunksStatsMessageToRaidLogChannel(nft)
}

// SendRaidConcludedMessageToRaidLogChannel sends a message when a raid is created
func SendMessageToRaidLogChannel(msg string, nft1 db.Nft, nft2 db.Nft) {
	sendMessageToRaidLogChannel(msg)
	sendFlunksStatsMessageToRaidLogChannel(nft1)
	sendFlunksStatsMessageToRaidLogChannel(nft2)
}

func sendMessageToRaidLogChannel(message string) {
	_, err := dg.ChannelMessageSend(RAID_LOG_CHANNEL_ID, message)
	if err != nil {
		fmt.Println("Error sending message to channel:", err)
	}
}

func sendPureFlunksStatsMessageToRaidLogChannel(nft db.Nft) {
	var fields []*discordgo.MessageEmbedField

	traits := nft.GetTraits()
	for _, trait := range traits {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   trait.Name,
			Value:  trait.Value,
			Inline: false,
		})
	}

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

func sendFlunksStatsMessageToRaidLogChannel(nft db.Nft) {
	var fields []*discordgo.MessageEmbedField

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Challenge Accepted!",
		Inline: false,
	})

	traits := nft.GetTraits()
	for _, trait := range traits {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   trait.Name,
			Value:  trait.Value,
			Inline: false,
		})
	}

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
