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

func GenerateBattlePrompt(clique string, challenger, defender uint, location string, isPositiveOutcome bool) string {
	return fmt.Sprintf(
		`
	Suppose you are creating a 3-action battle between challenger Flunk #%d and defender Flunk #%d.
	1st item is weapon, which describes the weapons that challenger Flunk has picked,
	2nd item is the action that challenger Flunk has taken with the weapon picked.
	3rd item is actionOutcome, it will be a %v action for the challenger.
	4th item is battleOutcome, challenger wins if challenger's actionOutcome is positive, otherwise defender wins.
	Put the winner in a winner variable in the json object, 0 when challenger's actionOutcome is positive and 1 when challenger's actionOutcome is negative.
	Make sure the weapon picked fits the clique of a %s battle in a high school setup, happening in %s.
	Make sure the wording fits the style in the example provided below.
	Make sure to bold all the Flunk #%d and Flunk #%d in the example provided below with **<content>**.
	Generate a JSON representation for a 'BattleLog' that includes fields 'Action', 'ActionOutcome', and 'BattleOutcome'.
	Example:
	Positive payload:
{
	"weapon": "Look at those gains! **Flunk #%d** is getting prepped with a Protein Shake.",
	"action": "**Flunk #%d** finished the protein shake.",
	"actionOutcome": "That should not be allowed in school, **Flunk #%d** has turned green after taking the protein shake and SMASHES **Flunk #%d** to smithereens.",
	"battleOutcome": "**Flunk #%d** emerges victorious, crushing **Flunk #%d** in the battle!",
	"winner": 0,
}
Negative payload:
{
	"weapon": "Look at those gains! **Flunk #%d** is getting prepped with a Protein Shake.",
	"action": "**Flunk #%d** finished the protein shake.",
	"actionOutcome": "**Flunk #%d** downed the lumpy-ass protein shake in one, then flexed his muscles so much they passed out",
	"battleOutcome": "What a shame! **Flunk #%d** won the battle with no effort!",
	"winner": 0,
}

	`, challenger, defender, isPositiveOutcome, clique, location, challenger, defender, challenger, challenger, challenger, defender, challenger, defender, challenger, challenger, challenger, defender)
}
