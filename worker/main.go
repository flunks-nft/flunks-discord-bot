package worker

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/flunks-nft/discord-bot/db"
	"github.com/flunks-nft/discord-bot/discord"
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
	msg := fmt.Sprintf("Flunk #%v has accepted #%v's challenge! It's a <%v> game!", raid.ToTemplateID, raid.FromTemplateID, raid.ChallengeID)
	discord.SendMessageToRaidLogChannel(msg, nfts[0], nfts[1])

	return nil
}
