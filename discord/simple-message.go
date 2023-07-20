package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendMessageToRaidLogChannel(msg_1, msg_2, image1URL, image2URL string) {
	embed := &discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Challenge Accepted!",
				Value:  msg_1,
				Inline: false,
			},
			{
				Value:  msg_2,
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: image1URL,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: image2URL,
		},
	}

	_, err := dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed)
	if err != nil {
		fmt.Println("Error sending embedded message to channel:", err)
	}

}
