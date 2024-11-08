package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/discord"
	"github.com/flunks-nft/discord-bot/pkg/worker"
)

func main() {
	// Connect to database & run migrations
	db.InitDB()

	// Create a wait group to coordinate shutdown
	wg := sync.WaitGroup{}
	wg.Add(2) // Number of services to wait for shutdown

	// Create a done channel to signal termination (CTRL + C)
	// Note: signal.Notify specifies which signals to send to the done channel
	// When one of the specified signals is received by the program, it will be sent on the done channel.
	// The program can then receive the signal from the done channel and handle it accordingly.
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Start worker service
	go worker.InitRaidWorker(&wg, done)

	// Start Discord text bot
	go discord.InitPureTextDiscord(&wg, done)

	// For Cloud Run: Start HTTP server on port 8080, only runs in production environment
	if os.Getenv("ENV") == "production" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello from Cloud Run!"))
		})

		go func() {
			if err := http.ListenAndServe(":8080", nil); err != nil {
				panic(err)
			}
		}()
	}

	// Wait for services to complete shutdown
	wg.Wait()

	// Signal termination to the done channel
	close(done)
}
