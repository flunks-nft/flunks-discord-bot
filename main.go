package main

import (
	"github.com/flunks-nft/discord-bot/db"
	"github.com/flunks-nft/discord-bot/discord"
)

func main() {
	// Connect to database & run migrations
	db.InitDB()
	discord.InitDiscord()
}
