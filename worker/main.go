package worker

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	Tiker *time.Ticker
)

func init() {
	Tiker = time.NewTicker(3 * time.Second)
}

func InitRaidWorker(wg *sync.WaitGroup, done chan os.Signal) {
	log.Println("🌱 Worker is now running. Press CTRL-C to exit.")

	// Run the worker loop until the done signal is received
	for {
		select {
		case <-Tiker.C:
			// Worker logic here
			fmt.Println("Worker is running...")

		case <-done:
			// Stop the worker
			Tiker.Stop()
			wg.Done() // Signal worker has completed shutdown
			log.Println("🐒 Worker is Gracefully shut down.")
			return
		}
	}
}
