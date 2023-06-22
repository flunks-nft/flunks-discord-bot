package raid

import "github.com/bwmarrin/discordgo"

type RaidBot struct {
	discord *discordgo.Session
}

func NewRaidBot(token string) (*RaidBot, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &RaidBot{discord}, nil
}
