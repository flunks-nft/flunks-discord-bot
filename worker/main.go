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
	Tiker = time.NewTicker(5 * time.Second)
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

	msg_1 := fmt.Sprintf("Flunk #%v has accepted #%v's challenge!", raid.ToTemplateID, raid.FromTemplateID)
	msg_2 := fmt.Sprintf("Challenge type <%v>", raid.ChallengeID)
	discord.SendMessageToRaidLogChannel(msg_1, msg_2, nfts[0].Uri, nfts[1].Uri)

	return nil
}
