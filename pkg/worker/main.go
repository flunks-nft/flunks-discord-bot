package worker

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/discord"
)

var (
	Tiker *time.Ticker
)

func init() {
	Tiker = time.NewTicker(5 * time.Second)
}

func InitRaidWorker(wg *sync.WaitGroup, done chan os.Signal) {
	log.Println("🌱 Worker is now running. Press CTRL-C to exit.")

	// Run the worker loop until the done signal is received
	for {
		select {
		case <-Tiker.C:
			createMatchedChallenge()
			concludeRaid()

		case <-done:
			// Stop the worker
			Tiker.Stop()
			wg.Done() // Signal worker has completed shutdown
			log.Println("🐒 Worker is Gracefully shut down.")
			return
		}
	}
}

func createMatchedChallenge() error {
	raid, nfts, err := db.QueueNextTokenPairForRaiding()
	if err != nil {
		return err
	}

	// TODO: get token ownership from the database and mention@ the owners
	fromNft := nfts[0]
	toNft := nfts[1]

	fromNftOwnerDiscordID := fromNft.Owner.DiscordID
	toNftOwnerDiscordID := toNft.Owner.DiscordID

	msg := fmt.Sprintf(
		"<@%s> #%d is challenging <@%s> #%d, it's a %s **%v** game!",
		fromNftOwnerDiscordID,
		fromNft.TemplateID,
		toNftOwnerDiscordID,
		toNft.TemplateID,
		raid.ChallengeTypeEmoji(),
		raid.ChallengeType,
	)
	discord.SendMessageToRaidLogChannel(msg, fromNft, toNft)

	return nil
}

func concludeRaid() {
	raid, err := db.ConcludeOneRaid()
	if err != nil {
		log.Println(err)
		return
	}
	emoji := fmt.Sprintf("<:emoji:%s>", db.RAID_WON_EMOJI_ID)
	msgs := make([]string, 0)
	msgs = append(msgs, fmt.Sprintf(
		"%s %s game concluded:\n",
		raid.ChallengeTypeEmoji(),
		raid.ChallengeType,
	))
	msgs = append(msgs, fmt.Sprintf(
		"Flunk #%d ⚔️ Flunk #%d\n",
		raid.FromNft.TemplateID,
		raid.ToNft.TemplateID,
	))
	msgs = append(msgs, fmt.Sprintf("%sWinner: Flunk #%d", emoji, raid.WinnerNft.TemplateID))

	var winnerThread, loserThread string
	winnerThread = fmt.Sprintf("<@%s>", raid.WinnerNft.Owner.DiscordID)
	if raid.WinnerNft.ID == raid.FromNftID {
		loserThread = fmt.Sprintf("<@%s>", raid.ToNft.Owner.DiscordID)
	} else {
		loserThread = fmt.Sprintf("<@%s>", raid.FromNft.Owner.DiscordID)
	}

	discord.SendRaidConcludedMessageToRaidLogChannel(msgs, raid.WinnerNft, winnerThread, loserThread)
}