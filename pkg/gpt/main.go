package gpt

import (
	"fmt"
	"log"
	"os"

	"github.com/ayush6624/go-chatgpt"
	"github.com/flunks-nft/discord-bot/pkg/utils"
)

var (
	ChapGPTClient *chatgpt.Client
)

func init() {
	utils.LoadEnv()

	key := os.Getenv("OPENAI_KEY")

	var err error
	ChapGPTClient, err = chatgpt.NewClient(key)
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateBattlePrompt(clique string, challenger, defender uint) string {
	return fmt.Sprintf(
		`
	Suppose you are creating a 3-action battle between challenger Flunk #%d and defender Flunk #%d.
	1st item is weapon, which describes the weapons that challenger Flunk has picked,
	2nd item is the action that challenger Flunk has taken.
	3rd item is actionOutcome, it can either be a positive or negative action for the challenger with equal probability.
	4th item is battleOutcome, it depends on the action is positive or negative in the action generated in the actionOutcome.
	Put the winner in a winner variable in the json object, 0 when challenger's actionOutcome is positive and 1 when challenger's actionOutcome is negative.
	Make sure the weapon picked fits the clique of a %s battle in a high school setup.
	Make sure the wording fits the style in the example provided below.
	Generate a JSON representation for a 'BattleLog' that includes fields 'Action', 'ActionOutcome', and 'BattleOutcome'.
	Example:
{
	"weapon": "Look at those gains! Flunk #%d is getting prepped with a Protein Shake.",
	"action": "Flunk #%d finished the protein shake.",
	"actionOutcome": "That should not be allowed in school, Flunk #%d has turned green after taking the protein shake and SMASHES Flunk #%d to smithereens.",
	"battleOutcome": "Flunk #%d nailed it, we have a winner!",
	"winner": 0,
}
	`, challenger, defender, clique, challenger, challenger, challenger, defender, challenger)
}
