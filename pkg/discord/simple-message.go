package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
)

func PostRaidAcceptedMsg(raid *db.Raid, nfts []db.Nft) {
	nft1 := nfts[0]
	nft2 := nfts[1]

	var fields []*discordgo.MessageEmbedField

	// Attach fist line challenge accepted message
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Challenge Accepted!", raid.ChallengeTypeEmoji()),
		Inline: false,
	})

	// Attach second line descriptions
	fields = append(fields, &discordgo.MessageEmbedField{
		Value: fmt.Sprintf(
			"Flunk #%d has accepted Flunk #%d's challenge to a %v battle",
			nft2.TemplateID,
			nft1.TemplateID,
			raid.ChallengeType,
		),
		Inline: false,
	})

	// Attach the from and to classes
	var challengeClassMsg, defenderClassMsg string
	challengeClassMsg = fmt.Sprintf("<@%s>", nft1.Owner.DiscordID)
	defenderClassMsg = fmt.Sprintf("<@%s>", nft2.Owner.DiscordID)

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Challenger Class",
		Value:  challengeClassMsg,
		Inline: false,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Defender Class",
		Value:  defenderClassMsg,
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

func PostRaidDetailsMsg(raid *db.Raid) {
	// nft := raid.FromNft

	var fields []*discordgo.MessageEmbedField

	// battleDesc := fmt.Sprintf(
	// 	"%s | Look at those gains! Flunk #%d is getting prepped with a **Protein Shake**. \n"+
	// 		"%s | BOOM! **PROTEIN SHAKE** TO THE DOME, Flunk #%d just lobbed that shake straight into their face. \n"+
	// 		"%s | Flunk #%d won!",
	// 	raid.ChallengeTypeEmoji(), nft.TemplateID,
	// 	raid.ChallengeTypeEmoji(), nft.TemplateID,
	// 	raid.ChallengeTypeEmoji(), nft.TemplateID,
	// )

	battleLog := fmt.Sprintf(
		"%s | %s \n"+
			"%s | %s \n"+
			"%s | %s",
		raid.ChallengeTypeEmoji(), raid.BattleLog.Action,
		raid.ChallengeTypeEmoji(), raid.BattleLog.ActionOutcome,
		raid.ChallengeTypeEmoji(), raid.BattleLog.BattleOutcome,
	)

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Battle Log",
		Value:  battleLog,
		Inline: false,
	})

	embed := &discordgo.MessageEmbed{
		Fields: fields,
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
