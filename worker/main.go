package worker

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/flunks-nft/discord-bot/db"
	"github.com/flunks-nft/discord-bot/discord"
	"github.com/flunks-nft/discord-bot/utils"
)

var (
	Tiker *time.Ticker
)

func init() {
	Tiker = time.NewTicker(2 * time.Second)
}

func InitRaidWorker(wg *sync.WaitGroup, done chan os.Signal) {
	log.Println("🌱 Worker is now running. Press CTRL-C to exit.")

	// Run the worker loop until the done signal is received
	for {
		select {
		case <-Tiker.C:
			createMatchedChallenge()

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

	cliqueEmoji := utils.CliqueEmojis[raid.ChallengeType.String()]

	msg := fmt.Sprintf("<@%s> Your Flunk has accepted <@%s>'s challenge! It's a %s **%v** game!", fromNftOwnerDiscordID, toNftOwnerDiscordID, cliqueEmoji, raid.ChallengeType)
	discord.SendMessageToRaidLogChannel(msg, fromNft, toNft)

	return nil
}
