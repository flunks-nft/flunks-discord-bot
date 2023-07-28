package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/db"
)

// SendRaidConcludedMessageToRaidLogChannel sends a message when a raid is created
func SendMessageToRaidLogChannel(msg string, nft1 db.Nft, nft2 db.Nft) {
	sendMessageToRaidLogChannel(msg)
	sendFlunksStatsMessageToRaidLogChannel(nft1, nft2)
}

func sendMessageToRaidLogChannel(message string) {
	_, err := dg.ChannelMessageSend(RAID_LOG_CHANNEL_ID, message)
	if err != nil {
		fmt.Println("Error sending message to channel:", err)
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

func sendFlunksStatsMessageToRaidLogChannel(nft1 db.Nft, nft2 db.Nft) {
	var fields []*discordgo.MessageEmbedField

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Challenge Accepted!",
		Inline: false,
	})

	traits1 := nft1.GetTraits()
	traits2 := nft2.GetTraits()

	for _, trait := range traits1 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   trait.Name,
			Value:  trait.Value,
			Inline: false,
		})
	}

	embed1 := &discordgo.MessageEmbed{
		Fields: fields,
		Image: &discordgo.MessageEmbedImage{
			URL: nft1.Uri,
		},
	}

	fields = []*discordgo.MessageEmbedField{} // reset fields for second embed

	for _, trait := range traits2 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   trait.Name,
			Value:  trait.Value,
			Inline: false,
		})
	}

	embed2 := &discordgo.MessageEmbed{
		Fields: fields,
		Image: &discordgo.MessageEmbedImage{
			URL: nft2.Uri,
		},
	}

	_, err := dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed1)
	if err != nil {
		fmt.Println("Error sending first embedded message to channel:", err)
	}

	_, err = dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed2)
	if err != nil {
		fmt.Println("Error sending second embedded message to channel:", err)
	}
}
