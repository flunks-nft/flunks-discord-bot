package worker

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/discord"
	"github.com/flunks-nft/discord-bot/pkg/utils"
)

var (
	Tiker *time.Ticker
)

func init() {
	RAID_MATCH_INTERVAL_IN_SECONDS := os.Getenv("RAID_MATCH_INTERVAL_IN_SECONDS")
	RAID_MATCH_INTERVAL_IN_SECONDS_INT, _ := utils.StringToInt(RAID_MATCH_INTERVAL_IN_SECONDS)
	RAID_MATCH_INTERVAL := time.Duration(RAID_MATCH_INTERVAL_IN_SECONDS_INT) * time.Second
	Tiker = time.NewTicker(RAID_MATCH_INTERVAL)
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

	discord.PostRaidAcceptedMsg(raid, nfts)

	return nil
}

func concludeRaid() {
	raid, err := db.ConcludeOneRaid()
	if err != nil {
		log.Println(err)
		return
	}
	discord.PostRaidDetailsMsg(&raid, discord.RAID_LOG_CHANNEL_ID)
}
