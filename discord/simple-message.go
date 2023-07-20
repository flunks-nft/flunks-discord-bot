package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendMessageToRaidLogChannel(msg_1, msg_2, image1URL, image2URL string) {
	sendMessageToRaidLogChannel(msg_1)
	sendFlunksStatsMessageToRaidLogChannel(image1URL)
	sendFlunksStatsMessageToRaidLogChannel(image2URL)

}

func sendMessageToRaidLogChannel(message string) {
	_, err := dg.ChannelMessageSend(RAID_LOG_CHANNEL_ID, message)
	if err != nil {
		fmt.Println("Error sending message to channel:", err)
	}
}

// TODO: display Flunk stats in the message
func sendFlunksStatsMessageToRaidLogChannel(imgUrl string) {
	embed := &discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Challenge Accepted!",
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: imgUrl,
		},
	}

	_, err := dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed)
	if err != nil {
		fmt.Println("Error sending embedded message to channel:", err)
	}
}
