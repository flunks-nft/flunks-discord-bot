package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/utils"
)

func PostRaidAcceptedMsg(raid *db.Raid) {
	var fields []*discordgo.MessageEmbedField

	// Attach fist line challenge accepted message
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Challenge Accepted! (ID: %d)", raid.ChallengeTypeEmoji(), raid.ID),
		Inline: false,
	})

	// Attach second line descriptions
	fields = append(fields, &discordgo.MessageEmbedField{
		Value: fmt.Sprintf(
			"**Flunk #%d** has accepted **Flunk #%d**'s challenge to a %v battle in the **%s**.",
			raid.FromTemplateID,
			raid.ToTemplateID,
			raid.ChallengeType,
			raid.BattleLocation,
		),
		Inline: false,
	})

	// Attach the from and to classes
	challengerClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.FromNft.Owner.DiscordID, raid.FromNft.TemplateID)
	defenderClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.ToNft.Owner.DiscordID, raid.ToNft.TemplateID)
	fields = append(fields, &discordgo.MessageEmbedField{
		Value:  fmt.Sprintf("**Challenger Class**: %s", challengerClass),
		Inline: false,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Value:  fmt.Sprintf("**Defender Class**: %s", defenderClass),
		Inline: false,
	})

	battleBgImgMap, _ := utils.BattleBgImages[raid.ChallengeType.String()]
	battleBgImgUrl := battleBgImgMap[raid.BattleLocation]

	embed := &discordgo.MessageEmbed{
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: battleBgImgUrl,
		},
	}

	msg, err := dg.ChannelMessageSendEmbed(RAID_LOG_CHANNEL_ID, embed)
	if err != nil {
		fmt.Println("Error sending embedded message to channel:", err)
	}

	// Update the raid with the message ID
	raid.SetMsgID(msg.ID)
}

func PostRaidDetailsMsgUpdate(raid *db.Raid, channelID string) string {
	var fields []*discordgo.MessageEmbedField

	// -------------------------Challenge Accepted Start-------------------------
	// Attach fist line challenge accepted message
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Challenge Accepted! (ID: %d)", raid.ChallengeTypeEmoji(), raid.ID),
		Inline: false,
	})

	// Attach second line descriptions
	fields = append(fields, &discordgo.MessageEmbedField{
		Value: fmt.Sprintf(
			"**Flunk #%d** has accepted **Flunk #%d**'s challenge to a %v battle in the **%s**.",
			raid.FromTemplateID,
			raid.ToTemplateID,
			raid.ChallengeType,
			raid.BattleLocation,
		),
		Inline: false,
	})

	challengerClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.FromNft.Owner.DiscordID, raid.FromNft.TemplateID)
	defenderClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.ToNft.Owner.DiscordID, raid.ToNft.TemplateID)
	fields = append(fields, &discordgo.MessageEmbedField{
		Value:  fmt.Sprintf("**Challenger Class**: %s", challengerClass),
		Inline: false,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Value:  fmt.Sprintf("**Defender Class**: %s", defenderClass),
		Inline: false,
	})

	battleBgImgMap, _ := utils.BattleBgImages[raid.ChallengeType.String()]
	battleBgImgUrl := battleBgImgMap[raid.BattleLocation]

	// -------------------------Challenge Accepted End-------------------------

	// -------------------------Battle Log Start-------------------------
	// Battle concluded, display the winner class
	if raid.BattleLogNounce >= 3 {
		winnerClass := fmt.Sprintf("<@%s> **Flunk #%d**", raid.WinnerNft.Owner.DiscordID, raid.WinnerTemplateID)
		fields = append(fields, &discordgo.MessageEmbedField{
			Value:  fmt.Sprintf("**Winner Class**: %s", winnerClass),
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
			URL: battleBgImgUrl,
		},
	}

	// Edit message
	_, err := dg.ChannelMessageEditEmbed(
		channelID,
		raid.BattleLogMessageID,
		embed,
	)
	if err != nil {
		fmt.Println("Error editing embedded message to channel:", err)
	}

	// -------------------------Battle Log End-------------------------

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
