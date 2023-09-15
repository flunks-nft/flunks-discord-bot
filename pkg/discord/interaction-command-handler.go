package discord

// This modules is responsible for handling slash commands

import (
	"flag"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/utils"
)

var (
	GuildID *string

	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "dapper",
			Description: "Set up your Dapper wallet address",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "address",
					Description: "Dapper wallet address",
					Required:    true,
				},
			},
		},
		{
			Name:        "flunk",
			Description: "Check Flunk by ID",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "edition",
					Description: "Flunks Edition Number",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"dapper": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			margs := make([]interface{}, 0, len(options))
			msgformat := "Your Dapper addressed is updated to: \n"

			// Get the value from the option map.
			// When the option exists, ok = true
			option, ok := optionMap["address"]
			if ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				margs = append(margs, option.StringValue())
				msgformat += "> address: %s\n"
			}

			msg := fmt.Sprintf(
				msgformat,
				margs...,
			)

			// Update user dapper wallet address
			db.CreateOrUpdateFlowAddress(i.Member.User.ID, option.StringValue())

			respondeEphemeralMessage(s, i, msg)
		},
		"flunk": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user, err := ValidateUser(i)
			if err != nil {
				respondeEphemeralMessage(s, i, err.Error())
				return
			}

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			editionStr := optionMap["edition"].StringValue()
			editionInt, _ := utils.StringToInt(editionStr)

			// TODO: add edition pass-in
			HandlesYearbook(s, i, user, editionInt)
		},
	}

	registeredCommands []*discordgo.ApplicationCommand
)

func init() {
	flag.Parse()
	GuildID = DefaultGuiID()
}

func registerSlashCommands() {
	log.Println("🚧 Adding commands...")
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func removeSlashCommands() {
	log.Println("🚧 Removing commands...")
	// We need to fetch the commands, since deleting requires the command ID.
	// We are doing this from the returned commands from `registeredCommands`, because using
	// this will delete all the commands, which might not be desirable, so we
	// are deleting only the commands that we added.

	for _, v := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
