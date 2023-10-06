package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/utils"
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
			"**Flunk #%d** has accepted **Flunk #%d**'s challenge to a %v battle in the **%s**.",
			nft2.TemplateID,
			nft1.TemplateID,
			raid.ChallengeType,
			raid.BattleLocation,
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
	fields = append(fields, &discordgo.MessageEmbedField{
		Value:  fmt.Sprintf("**Battle ID: %d**", raid.ID),
		Inline: false,
	})

	battleBgImgMap, _ := utils.BattleBgImages[raid.ChallengeType.String()]
	battleBgImgUrl := battleBgImgMap[raid.BattleLocation]

	embed := &discordgo.MessageEmbed{
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: nft1.Uri,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: battleBgImgUrl,
		},
	}

	_, err := dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed)
	if err != nil {
		fmt.Println("Error sending embedded message to channel:", err)
	}
}

func PostRaidDetailsMsgUpdate(raid *db.Raid, channelID string) string {
	var fields []*discordgo.MessageEmbedField

	fields = append(fields, &discordgo.MessageEmbedField{
		Value:  fmt.Sprintf("**Battle ID: %d**", raid.ID),
		Inline: false,
	})

	if raid.BattleLogNounce < 3 {
		// Battle not concluded yet, use challenger and defender
		challengerClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.FromNft.Owner.DiscordID, raid.FromNft.TemplateID)
		defenderClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.ToNft.Owner.DiscordID, raid.ToNft.TemplateID)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Challenger Class",
			Value:  challengerClass,
			Inline: false,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Defender Class",
			Value:  defenderClass,
			Inline: false,
		})
	} else {
		winnerClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.WinnerNft.Owner.DiscordID, raid.WinnerTemplateID)
		loserClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.LoserNft.Owner.DiscordID, raid.LoserTemplateID)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Winner Class",
			Value:  winnerClass,
			Inline: false,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Loser Class",
			Value:  loserClass,
			Inline: false,
		})
	}

	battleLog := parseBattleLogFromNounce(raid)

	fields = append(fields, &discordgo.MessageEmbedField{
		Value:  battleLog,
		Inline: false,
	})

	embed := &discordgo.MessageEmbed{
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: raid.FromNft.Uri,
		},
	}

	if raid.BattleLogMessageID == "" {
		// First time sending msg
		msg, err := dg.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			fmt.Println("Error sending embedded message to channel:", err)
		}

		// return ID of the msg and new nounce
		return msg.ID
	} else {
		// Later time editing msg
		_, err := dg.ChannelMessageEditEmbed(
			channelID,
			raid.BattleLogMessageID,
			embed,
		)
		if err != nil {
			fmt.Println("Error editing embedded message to channel:", err)
		}
	}

	return ""
}

func parseBattleLogFromNounce(raid *db.Raid) string {
	if raid.BattleLogNounce == 0 {
		return fmt.Sprintf(
			"%s | Flunk **#%d** is deciding on the weapon. \n"+
				"......",
			raid.ChallengeTypeEmoji(), raid.FromTemplateID,
		)
	}

	if raid.BattleLogNounce == 1 {
		return fmt.Sprintf(
			"%s | %s \n"+
				"......",
			raid.ChallengeTypeEmoji(), raid.BattleLog.Weapon,
		)
	}

	if raid.BattleLogNounce == 2 {
		return fmt.Sprintf(
			"%s | %s \n"+
				"%s | %s \n"+
				"......",
			raid.ChallengeTypeEmoji(), raid.BattleLog.Weapon,
			raid.ChallengeTypeEmoji(), raid.BattleLog.Action,
		)
	}

	// Default return the full battle log
	return fmt.Sprintf(
		"%s | %s \n"+
			"%s | %s \n"+
			"%s | %s \n"+
			"%s | %s",
		raid.ChallengeTypeEmoji(), raid.BattleLog.Weapon,
		raid.ChallengeTypeEmoji(), raid.BattleLog.Action,
		raid.ChallengeTypeEmoji(), raid.BattleLog.ActionOutcome,
		raid.ChallengeTypeEmoji(), raid.BattleLog.BattleOutcome,
	)
}
