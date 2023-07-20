package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/flunks-nft/discord-bot/db"
	"github.com/flunks-nft/discord-bot/discord"
	"github.com/flunks-nft/discord-bot/worker"
)

func main() {
	// Connect to database & run migrations
	db.InitDB()

	// Create a wait group to coordinate shutdown
	wg := sync.WaitGroup{}
	wg.Add(2) // Number of services to wait for shutdown

	// Create a done channel to signal termination (CTRL + C)
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Start Discord service
	go discord.InitDiscord(&wg, done)

	// Start worker service
	go worker.InitRaidWorker(&wg, done)

	// Wait for services to complete shutdown
	wg.Wait()

	// Signal termination to the done channel
	close(done)
}
