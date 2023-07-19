package worker

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func InitRaidWorker(wg *sync.WaitGroup, done chan struct{}) {
	// Create a ticker for the worker
	ticker := time.NewTicker(3 * time.Second)

	// Start worker service
	go runWorker(ticker, done, wg)

	// Signal the WaitGroup that Worker service has started
	wg.Done()

	// Wait for termination signal
	log.Println("🌱 Worker is now running. Press CTRL-C to exit.")
	<-done

	log.Println("🐒 Worker is shutdown.")
}

func runWorker(ticker *time.Ticker, done chan struct{}, wg *sync.WaitGroup) {
	// Run the worker loop until the done signal is received
	for {
		select {
		case <-ticker.C:
			// Worker logic here
			fmt.Println("Worker is running...")

		case <-done:
			// Stop the worker
			ticker.Stop()
			wg.Done() // Signal worker has completed shutdown
			return
		}
	}
}
