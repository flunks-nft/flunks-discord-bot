package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendMessageToRaidLogChannel(message, image1URL, image2URL string) {
	embed := &discordgo.MessageEmbed{
		Description: message,
		Image: &discordgo.MessageEmbedImage{
			URL: image1URL,
		},
	}

	// Add a field for the second image
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "VS",
	})

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "VS",
		Value: image2URL,
	})

	_, err := dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed)
	if err != nil {
		fmt.Println("Error sending embedded message to channel:", err)
	}

}
