package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/db"
)

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

// TODO: display Flunk stats in the message
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
